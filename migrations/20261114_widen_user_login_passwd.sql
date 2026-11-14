-- Migration: widen user_login.passwd so bcrypt hashes fit
-- Date: 2026-11-14
-- Description: passwd 原本是 char(32)（MD5 时代）。bcrypt 哈希固定 60 字符，写进去会：
--   * 非严格 sql_mode：静默截断成 32 字符，用户此后永久无法登录（注册、改密、
--     重置密码、md5 登录自动升级四条路径全部受影响）；
--   * 严格 sql_mode：INSERT/UPDATE 直接报 Data too long，注册与改密全部失败。
-- 本迁移加宽该列、补齐老 db.sql 缺失的 passwd_type，并清理已被截断的坏数据。
-- 可重复执行；已跑过 20260329_add_passwd_type_column.sql 的库同样适用。

-- 1. 补齐 passwd_type（用老版 config/db.sql 建的库没有这一列）
SET @ddl = (
  SELECT IF(COUNT(*) = 0,
    'ALTER TABLE user_login ADD COLUMN passwd_type VARCHAR(10) NOT NULL DEFAULT ''md5'' COMMENT ''密码类型：md5(旧) 或 bcrypt(新)''',
    'DO 0')
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'user_login'
    AND COLUMN_NAME = 'passwd_type'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2. 加宽 passwd，容纳 bcrypt（60 字符）及后续更长的哈希算法
ALTER TABLE user_login
  MODIFY COLUMN `passwd` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码哈希：md5 为 32 字符，bcrypt 为 60 字符';

-- 3. 清理已被截断的 bcrypt 哈希。哈希不可恢复，只能清空密码，
--    让这些用户走「忘记密码」邮件重置重新拿回账号。
--    合法 bcrypt 哈希恒为 60 字符，据此判定坏数据。
UPDATE user_login
SET passwd = '', passcode = '', passwd_type = 'md5'
WHERE passwd_type = 'bcrypt' AND CHAR_LENGTH(passwd) <> 60;