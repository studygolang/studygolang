# StudyGolang 环境变量配置指南

## 概述

StudyGolang 支持通过环境变量覆盖配置文件中的敏感信息。环境变量的优先级高于配置文件。

## 使用方法

### 1. 创建 .env 文件

```bash
cp .env.example .env
# 编辑 .env 文件，填入实际值
vim .env
```

### 2. 加载环境变量

```bash
# 方式 1: 使用 source
source .env

# 方式 2: 使用 export
export $(cat .env | xargs)

# 方式 3: 在 systemd service 文件中
# [Service]
# EnvironmentFile=/path/to/.env
```

## 环境变量列表

### 数据库配置 (MySQL)

| 环境变量 | 配置文件路径 | 说明 | 示例 |
|---------|------------|------|------|
| `DB_HOST` | `[mysql] host` | MySQL 主机地址 | `127.0.0.1` |
| `DB_PORT` | `[mysql] port` | MySQL 端口 | `3306` |
| `DB_USER` | `[mysql] user` | MySQL 用户名 | `root` |
| `DB_PASSWORD` | `[mysql] password` | MySQL 密码 | `your_password` |
| `DB_NAME` | `[mysql] dbname` | 数据库名 | `studygolang` |

### Redis 配置

| 环境变量 | 配置文件路径 | 说明 | 示例 |
|---------|------------|------|------|
| `REDIS_HOST` | `[redis] host` | Redis 主机地址 | `127.0.0.1` |
| `REDIS_PORT` | `[redis] port` | Redis 端口 | `6379` |
| `REDIS_PASSWORD` | `[redis] password` | Redis 密码 | `your_redis_password` |

### 安全配置

| 环境变量 | 配置文件路径 | 说明 | 要求 |
|---------|------------|------|------|
| `TOKEN_SALT` | `[security] token_salt` | Token 签名盐值 | **必须**，64位随机字符串 |
| `UNSUBSCRIBE_TOKEN_KEY` | `[security] unsubscribe_token_key` | 退订邮件 token key | 20+ 字符 |
| `ACTIVATE_SIGN_SALT` | `[security] activate_sign_salt` | 激活邮件签名盐值 | 20+ 字符 |
| `COOKIE_SECRET` | `[global] cookie_secret` | Cookie 加密密钥 | 11+ 字符 |

### 邮件配置

| 环境变量 | 配置文件路径 | 说明 | 示例 |
|---------|------------|------|------|
| `SMTP_USERNAME` | `[email] smtp_username` | SMTP 用户名 | `noreply@example.com` |
| `SMTP_PASSWORD` | `[email] smtp_password` | SMTP 密码 | `your_smtp_password` |
| `SMTP_HOST` | `[email] smtp_host` | SMTP 主机 | `smtp.example.com` |
| `SMTP_PORT` | `[email] smtp_port` | SMTP 端口 | `587` |
| `SMTP_FROM_EMAIL` | `[email] from_email` | 发件人邮箱 | `noreply@example.com` |

### 认证邮件配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `AUTH_SMTP_USERNAME` | `[email.auth] smtp_username` | 认证邮件 SMTP 用户名 |
| `AUTH_SMTP_PASSWORD` | `[email.auth] smtp_password` | 认证邮件 SMTP 密码 |
| `AUTH_SMTP_HOST` | `[email.auth] smtp_host` | 认证邮件 SMTP 主机 |
| `AUTH_SMTP_PORT` | `[email.auth] smtp_port` | 认证邮件 SMTP 端口 |
| `AUTH_SMTP_FROM_EMAIL` | `[email.auth] from_email` | 认证邮件发件人 |

### 七牛云存储配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `QINIU_ACCESS_KEY` | `[qiniu] access_key` | 七牛云 Access Key |
| `QINIU_SECRET_KEY` | `[qiniu] secret_key` | 七牛云 Secret Key |
| `QINIU_BUCKET_NAME` | `[qiniu] bucket_name` | 七牛云存储空间名 |

### 微信公众号配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `WECHAT_APPID` | `[wechat] appid` | 微信公众号 AppID |
| `WECHAT_APPSECRET` | `[wechat] appsecret` | 微信公众号 AppSecret |

### GitHub OAuth 配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `GITHUB_CLIENT_ID` | `[github] client_id` | GitHub OAuth Client ID |
| `GITHUB_CLIENT_SECRET` | `[github] client_secret` | GitHub OAuth Client Secret |

### Gitea OAuth 配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `GITEA_CLIENT_ID` | `[gitea] client_id` | Gitea OAuth Client ID |
| `GITEA_CLIENT_SECRET` | `[gitea] client_secret` | Gitea OAuth Client Secret |

### GCTT 配置

| 环境变量 | 配置文件路径 | 说明 |
|---------|------------|------|
| `GCTT_TOKEN_SECRET` | `[gctt] token_secret` | GCTT Token Secret |

## 生产环境部署示例

### systemd service 文件

```ini
[Unit]
Description=StudyGolang Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=studygolang
WorkingDirectory=/opt/studygolang
EnvironmentFile=/opt/studygolang/.env
ExecStart=/opt/studygolang/studygolang
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Docker Compose

```yaml
version: '3.8'

services:
  studygolang:
    image: studygolang:latest
    env_file:
      - .env
    ports:
      - "8090:8090"
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
```

### Kubernetes Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: studygolang-secrets
type: Opaque
stringData:
  DB_PASSWORD: "your_password"
  REDIS_PASSWORD: "your_redis_password"
  TOKEN_SALT: "your_64_char_token_salt"
  SMTP_PASSWORD: "your_smtp_password"
```

## 安全建议

1. **永远不要**将 `.env` 文件提交到版本控制系统
2. **生产环境**必须使用环境变量，不要在配置文件中存储敏感信息
3. **定期轮换**密钥和密码（建议每 90 天）
4. **最小权限原则**：数据库用户只授予必要的权限
5. **网络隔离**：敏感服务（MySQL、Redis）不要暴露到公网

## 验证配置

```bash
# 检查环境变量是否正确加载
echo $TOKEN_SALT
echo $DB_PASSWORD

# 检查配置文件是否正确
cat config/env.ini | grep -v "^;" | grep -v "^$"
```

## 常见问题

### Q: 环境变量和配置文件哪个优先级高？
**A:** 环境变量优先级更高。如果设置了环境变量，会忽略配置文件中的值。

### Q: 必须的环境变量有哪些？
**A:** 生产环境必须设置：
- `TOKEN_SALT`（Token 签名）
- `DB_PASSWORD`（数据库密码）
- `REDIS_PASSWORD`（Redis 密码）

### Q: 如何生成安全的随机字符串？
**A:** 使用以下命令：
```bash
# 64位随机字符串（用于 TOKEN_SALT）
openssl rand -hex 32

# 或者
head -c 32 /dev/urandom | base64
```

## 迁移指南

### 从硬编码迁移到环境变量

1. 导出现有配置到环境变量：
```bash
# 从配置文件提取敏感信息
grep "password\|secret\|token_salt" config/env.ini
```

2. 创建 `.env` 文件
3. 更新配置文件，移除敏感信息（用占位符替代）
4. 重启服务验证

## 更新历史

- 2024-03-27: 初始版本，支持核心配置的环境变量覆盖
