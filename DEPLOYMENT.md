# StudyGolang 部署指南

## 生产环境部署清单

### 1. 环境变量配置

**必须配置的环境变量（CRITICAL）：**

```bash
# Token 签名盐值（64位随机字符串）
export TOKEN_SALT="your_64_char_random_string_here"

# 数据库配置
export DB_HOST="your_db_host"
export DB_PORT="3306"
export DB_USER="your_db_user"
export DB_PASSWORD="your_secure_db_password"
export DB_NAME="studygolang"

# Redis 配置
export REDIS_HOST="your_redis_host"
export REDIS_PORT="6379"
export REDIS_PASSWORD="your_redis_password"
```

**推荐配置的环境变量：**

```bash
# JWT Token 签名密钥（生产环境必须配置，至少 32 字节）
export JWT_SECRET="your-super-secret-jwt-key-min-32-chars-long"

# CORS 允许的来源（生产环境必须配置，逗号分隔）
export ALLOWED_ORIGINS="https://studygolang.com,https://www.studygolang.com"

# Cookie 安全配置
export COOKIE_SECRET="your_cookie_secret"

# 邮件服务配置
export SMTP_USERNAME="your_smtp_username"
export SMTP_PASSWORD="your_smtp_password"
export SMTP_HOST="smtp.example.com"
export SMTP_PORT="587"

# 其他安全配置
export UNSUBSCRIBE_TOKEN_KEY="your_unsubscribe_key"
export ACTIVATE_SIGN_SALT="your_activate_salt"
```

### 2. 配置文件安全

**配置文件 `config/env.ini` 最佳实践：**

1. ✅ **敏感信息使用占位符**：
   ```ini
   [mysql]
   password = USE_ENV_DB_PASSWORD

   [security]
   token_salt = USE_ENV_TOKEN_SALT
   ```

2. ✅ **文件权限控制**：
   ```bash
   # 设置配置文件权限（仅所有者可读写）
   chmod 600 config/env.ini

   # 设置 .env 文件权限
   chmod 600 .env
   ```

3. ✅ **禁止提交敏感文件**：
   `.gitignore` 中已包含：
   ```
   .env
   config/env.ini
   ```

### 3. HTTPS 配置

**强制 HTTPS（生产环境）：**

1. Cookie 已配置：
   - `Secure=true` - 仅通过 HTTPS 传输
   - `SameSite=Strict` - 严格的 CSRF 保护
   - `HttpOnly=true` - 防止 XSS 攻击

2. Nginx 配置示例：
   ```nginx
   server {
       listen 443 ssl http2;
       server_name studygolang.com;

       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;

       # 强制 HTTPS
       add_header Strict-Transport-Security "max-age=31536000" always;

       location / {
           proxy_pass http://127.0.0.1:8090;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }

   server {
       listen 80;
       server_name studygolang.com;
       return 301 https://$server_name$request_uri;
   }
   ```

### 4. 数据库安全

1. ✅ **最小权限原则**：
   - 应用程序使用的数据库用户只授予必要权限
   - 禁止使用 root 用户

2. ✅ **SQL 注入防护**：
   - 使用参数化查询
   - 字段白名单验证

3. ✅ **连接池配置**：
   ```ini
   [mysql]
   max_idle = 2
   max_conn = 10
   ```

### 5. Redis 安全

1. ✅ **密码保护**：
   - 生产环境必须设置 Redis 密码
   - 使用环境变量 `REDIS_PASSWORD`

2. ✅ **网络隔离**：
   - Redis 不暴露到公网
   - 使用防火墙限制访问

3. ✅ **原子操作**：
   - 使用 Redis INCR 原子计数器
   - 避免竞态条件

### 6. Session 和 Cookie 安全

1. ✅ **Cookie 安全属性**：
   ```go
   HttpOnly = true    // 防止 XSS
   Secure = true       // 仅 HTTPS
   SameSite = Strict   // 防止 CSRF
   ```

2. ✅ **Session 配置**：
   - 使用 Redis 存储 session（支持多机部署）
   - 设置合理的过期时间

### 7. 监控和日志

1. **日志级别**：
   ```ini
   [global]
   env = prod
   log_level = INFO  # 生产环境使用 INFO 或 WARN
   ```

2. **关键操作日志**：
   - 用户登录/登出
   - 敏感操作（修改密码、删除内容）
   - 管理员操作

### 8. 定期安全检查

**每月检查：**

- [ ] 更换 TokenSalt（如果有泄露风险）
- [ ] 检查数据库用户权限
- [ ] 审查管理员账户列表
- [ ] 检查异常登录日志
- [ ] 更新依赖包版本（安全补丁）

**每季度检查：**

- [ ] 全面安全审计
- [ ] 渗透测试
- [ ] 代码审查（安全相关）

### 9. 备份策略

1. **数据库备份**：
   ```bash
   # 每日自动备份
   0 2 * * * /usr/bin/mysqldump -u backup_user -p studygolang > /backup/db_$(date +\%Y\%m\%d).sql
   ```

2. **Redis 持久化**：
   ```ini
   # redis.conf
   save 900 1  # 15分钟内至少 1 个 key 变更则保存
   appendonly yes  # AOF 持久化
   ```

### 10. 应急响应

**安全事件响应流程：**

1. **立即行动**：
   - 更换所有密钥和密码
   - 检查日志定位问题
   - 通知相关用户

2. **事后分析**：
   - 记录事件详情
   - 分析根本原因
   - 制定改进措施

3. **预防措施**：
   - 更新安全策略
   - 加强监控
   - 培训团队

---

## 快速部署脚本

```bash
#!/bin/bash
# deploy.sh - StudyGolang 部署脚本

set -e

echo "🚀 开始部署 StudyGolang..."

# 1. 检查环境变量
echo "✅ 检查环境变量..."
required_vars=("TOKEN_SALT" "DB_PASSWORD" "REDIS_PASSWORD")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ 错误: 环境变量 $var 未设置"
        exit 1
    fi
done

# 2. 构建应用
echo "🔨 构建应用..."
go build -o studygolang cmd/studygolang/main.go

# 3. 运行数据库迁移
echo "📊 运行数据库迁移..."
./studygolang migrate

# 4. 启动服务
echo "🚀 启动服务..."
systemctl restart studygolang

echo "✅ 部署完成！"
```

---

**重要提示：** 部署前请确保所有 CRITICAL 和 HIGH 优先级的安全修复已完成！

---

## 安全升级说明（2026-03-30）

### 密码存储升级

**变更内容：**
- 新用户：使用 bcrypt 哈希存储密码
- 现有用户：MD5 密码仍然有效，登录成功后自动升级到 bcrypt
- 影响表：`user_login` 新增 `passwd_type` 字段（`varchar(10) default('md5')`）

**迁移步骤：**
```bash
# 1. 执行数据库迁移
mysql -u username -p database_name < migrations/20260329_add_passwd_type_column.sql

# 2. 部署新代码（自动迁移，无需停机）
systemctl restart studygolang

# 3. 监控迁移进度（可选）
# 查看日志中 "自动升级用户密码到 bcrypt" 消息
```

**监控指标：**
- 登录失败率（应保持不变）
- bcrypt 升级成功数
- 升级失败数（仅记录错误，不影响登录）

**注意事项：**
- ✅ 无停机时间，用户无感知
- ✅ 现有用户无需重置密码
- ✅ 升级失败不影响登录（降级到 MD5 验证）
- ⚠️ 建议观察 6 个月迁移进度，考虑强制未升级用户重置密码

---

### Token 机制升级

**变更内容：**
- 新登录：使用 JWT Token（格式：`eyJ...`）
- 现有 Token：MD5 Token（格式：`{timestamp}{md5}uid{uid}`）30 天内仍有效
- 自动识别：`ValidateTokenAuto()` 自动识别 Token 类型

**部署要求：**
- ✅ `TOKEN_SALT` 必须保留（兼容旧 Token）
- ⚠️ `JWT_SECRET` 生产环境必须配置（至少 32 字节）

**监控建议：**
```bash
# 查看日志统计 JWT vs MD5 Token 使用比例
grep "使用.*Token" /var/log/studygolang/app.log | \
  awk '{print $NF}' | sort | uniq -c
```

**预期结果：**
- 初期：MD5 Token 占主导（90%+）
- 30 天后：JWT Token 占主导（逐步替换）
- 60 天后：MD5 Token 接近 0（自然过期）

---

### Redis 连接池优化

**变更内容：**
- 新增全局 Redis 连接池（`internal/logic/redis_pool.go`）
- `uuid_redis.go` 和 `user_cache.go` 使用连接池
- 避免每次操作创建新连接

**性能提升：**
- 减少连接创建开销
- 降低 Redis 服务器负载
- 提升并发处理能力

**注意事项：**
- ✅ 向后兼容，行为不变
- ⚠️ 程序退出时应调用 `logic.CloseRedisClient()`（可选）

---

### CSRF 保护增强

**变更内容：**
- 消息删除接口：从 Host 比较改为 `ALLOWED_ORIGINS` 白名单验证
- 支持通配符子域名（如 `*.studygolang.com`）
- 更严格的来源检查

**部署要求：**
- ⚠️ 生产环境必须配置 `ALLOWED_ORIGINS`（否则程序启动时 panic）
- 示例：`ALLOWED_ORIGINS=https://studygolang.com,https://www.studygolang.com`

**测试验证：**
```bash
# 测试 CSRF 保护
curl -X DELETE https://api.studygolang.com/api/v1/messages/123 \
  -H "Origin: https://evil.com" \
  -H "Cookie: sg_token=xxx"
# 预期：403 Forbidden

curl -X DELETE https://api.studygolang.com/api/v1/messages/123 \
  -H "Origin: https://studygolang.com" \
  -H "Cookie: sg_token=xxx"
# 预期：200 OK
```

---

### CORS 生产环境强制校验

**变更内容：**
- 生产环境（`env=prod`）必须配置 `ALLOWED_ORIGINS`
- 未配置时程序启动会 panic
- 开发环境（`env=dev`）使用默认值 `http://localhost:3000`

**配置示例：**
```bash
# .env 文件（生产环境）
ALLOWED_ORIGINS=https://studygolang.com,https://www.studygolang.com

# config/env.ini（备用）
[global]
env = prod
```

---

### 评论楼层计数器竞态修复

**变更内容：**
- 使用 Redis Lua 脚本保证原子性
- Redis 不可用时使用分布式锁回退
- 添加并发测试（100 goroutine）

**影响：**
- ✅ 完全消除楼层号重复问题
- ✅ 高并发场景稳定性提升
- ✅ Redis 故障时优雅降级

**验证方法：**
```bash
# 运行并发测试（需要 Redis + MySQL 环境）
go test -v -tags=integration -run TestGetNextCommentFloorConcurrency \
  ./internal/logic/
```

---

### Goroutine 泄漏修复

**变更内容：**
- 评论发布、用户登录、用户更新中的 goroutine 传递 context
- 使用 `errgroup.WithContext` 管理 goroutine 生命周期
- 程序关闭时可正确取消未完成的异步任务

**影响：**
- ✅ 避免 goroutine 泄漏
- ✅ 优雅关闭支持
- ✅ 更好的错误传播

**注意事项：**
- ⚠️ `context.Background()` 用于后台任务（不影响主流程）
- ⚠️ 异步任务失败仅记录日志，不影响响应

---

## 部署检查清单

**部署前：**
- [ ] 配置 `JWT_SECRET`（至少 32 字节随机字符串）
- [ ] 配置 `ALLOWED_ORIGINS`（生产环境必须）
- [ ] 保留 `TOKEN_SALT`（兼容旧 Token）
- [ ] 运行数据库迁移（`20260329_add_passwd_type_column.sql`）
- [ ] 检查 `.env` 文件权限（`chmod 600 .env`）

**部署后：**
- [ ] 查看日志确认无 panic
- [ ] 测试登录功能（新/旧用户）
- [ ] 测试消息删除功能（CSRF 保护）
- [ ] 监控 Redis 连接数（应保持稳定）
- [ ] 观察 bcrypt 迁移进度（日志中搜索 "自动升级"）

**监控指标：**
```bash
# 1. 密码迁移进度
grep -c "自动升级用户密码到 bcrypt" /var/log/studygolang/app.log

# 2. Token 使用分布
grep "使用.*Token" /var/log/studygolang/app.log | wc -l

# 3. CSRF 拦截统计
grep -c "非法的跨域请求" /var/log/studygolang/app.log

# 4. Redis 连接错误
grep -c "Redis.*error" /var/log/studygolang/app.log
```

---

## 回滚计划

**如果出现问题，回滚步骤：**

1. **回滚代码**
   ```bash
   git revert <commit-hash>
   systemctl restart studygolang
   ```

2. **数据库回滚**（可选）
   ```sql
   -- 如果需要删除 passwd_type 字段（通常不需要）
   ALTER TABLE user_login DROP COLUMN passwd_type;
   ```

3. **配置回滚**
   - 保留 `TOKEN_SALT`（旧版本需要）
   - 可删除 `JWT_SECRET` 和 `ALLOWED_ORIGINS`

**风险评估：**
- 低风险：所有修改都保持向后兼容
- 现有用户：无感知，功能正常
- 新用户：使用更安全的机制

---

## 常见问题

**Q1: 部署后旧用户无法登录？**
A1: 检查 `TOKEN_SALT` 是否保留。旧 Token 验证依赖此配置。

**Q2: 生产环境启动 panic "ALLOWED_ORIGINS must be set"？**
A2: 配置环境变量：`export ALLOWED_ORIGINS=https://your-domain.com`

**Q3: 日志中出现 "自动升级用户密码到 bcrypt" 失败？**
A3: 检查数据库连接和权限。升级失败不影响登录，仅记录错误。

**Q4: Redis 连接错误增加？**
A4: 检查 Redis 服务状态和连接池配置。连接池会自动重连。

**Q5: 评论楼层号仍然重复？**
A5: 检查 Redis 服务可用性。确保 Lua 脚本正常执行。
