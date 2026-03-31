-- Migration: Add passwd_type column for bcrypt migration
-- Date: 2026-03-29
-- Description: Add passwd_type field to support gradual migration from MD5 to bcrypt password hashing

-- Add passwd_type column (default 'md5' for backward compatibility)
ALTER TABLE user_login
ADD COLUMN passwd_type VARCHAR(10) DEFAULT 'md5'
COMMENT 'password type: md5(legacy) or bcrypt(new)';

-- Mark all existing users as md5 (for safety)
UPDATE user_login SET passwd_type = 'md5' WHERE passwd_type IS NULL OR passwd_type = '';
