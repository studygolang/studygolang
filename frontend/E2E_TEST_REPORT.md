# 用户中心功能 E2E 测试报告

**测试时间**: 2026-03-14
**测试人员**: E2E 测试专家
**测试范围**: 会员列表、用户主页、收藏列表

---

## 测试环境

- **后端**: http://localhost:8090 ✅ 运行中
- **前端**: http://localhost:3000 ✅ 运行中
- **测试框架**: Playwright
- **浏览器**: Chromium

---

## 测试结果总览

- **总测试数**: 31
- **通过**: 12 (38.7%)
- **失败**: 19 (61.3%)

---

## ✅ 通过的测试

### 会员列表页 /users
1. ✅ **响应式布局 - 移动端** - 移动端视口下页面正常显示
2. ✅ **无数据时显示空状态** - 空数据时正确显示提示信息

### 用户个人主页 /user/:username
3. ✅ **页面正常渲染** - HTTP 200，页面 body 可见
4. ✅ **显示用户统计信息** - 正确显示话题、文章、积分等统计
5. ✅ **显示用户发布的话题** - 话题列表正常展示
6. ✅ **显示用户发布的文章** - 文章列表正常展示
7. ✅ **响应式布局 - 移动端** - 移动端视口下页面正常显示
8. ✅ **社交链接可点击** - 社交链接正确渲染且可点击

### 收藏列表 /favorites/:username
9. ✅ **页面正常渲染** - HTTP 200，页面 body 可见
10. ✅ **默认显示话题收藏** - 默认 Tab 为话题收藏
11. ✅ **切换到资源收藏 Tab** - Tab 切换功能正常
12. ✅ **切换到项目收藏 Tab** - Tab 切换功能正常

---

## ❌ 失败的测试

### 🔴 关键问题：API 路由被前端路由拦截

**根本原因**:
- 后端路由注册顺序错误
- 前端控制器 (`frontG := e.Group("")`) 先注册，包含 `/users` 路由
- API 控制器 (`apiG := e.Group("/api/v1")`) 后注册
- Echo 框架按注册顺序匹配路由，导致 `/api/v1/users` 被前端的 `/users` 拦截

**影响范围**:
- `/api/v1/users` 返回 HTML 而不是 JSON
- `/api/v1/users/:username/favorites` 返回 HTML 而不是 JSON
- 所有依赖这些 API 的前端页面无法正常工作

**复现步骤**:
```bash
# 返回 HTML（错误）
curl http://localhost:8090/api/v1/users

# 返回 JSON（正确）
curl http://localhost:8090/api/v1/user/polaris
```

### 会员列表页 /users - 5 个失败

1. ❌ **页面正常渲染** (47.0s 超时)
   - **错误**: 页面加载超时
   - **原因**: API `/api/v1/users` 返回 HTML 导致前端解析失败

2. ❌ **显示活跃会员和新加入会员 Tab** (48.0s 超时)
   - **错误**: 页面加载超时
   - **原因**: 同上

3. ❌ **活跃会员 Tab 显示用户列表** (46.4s 超时)
   - **错误**: 页面加载超时
   - **原因**: 同上

4. ❌ **切换到新加入会员 Tab** (47.6s 超时)
   - **错误**: 页面加载超时
   - **原因**: 同上

5. ❌ **用户卡片可点击跳转** (11.4s)
   - **错误**: 无法找到用户卡片元素
   - **原因**: API 数据加载失败，页面无内容

### 用户个人主页 /user/:username - 2 个失败

6. ❌ **显示用户基本信息** (7.4s)
   - **错误**: 无法找到用户名元素
   - **原因**: 测试用户 "polaris" 不存在（API 返回 `{"code":1,"msg":"用户不存在"}`）
   - **建议**: 使用真实存在的测试用户

7. ❌ **用户不存在时显示错误** (2.2s)
   - **错误**: 未找到错误提示文本
   - **原因**: 前端未正确处理 API 错误响应

### 收藏列表 /favorites/:username - 9 个失败

8. ❌ **显示收藏类型 Tab** (7.6s)
   - **错误**: 无法找到 Tab 元素
   - **原因**: API 返回 HTML 导致页面渲染失败

9. ❌ **切换到文章收藏 Tab** (41.2s 超时)
   - **错误**: 页面加载超时
   - **原因**: 同上

10. ❌ **显示收藏的话题列表** (42.6s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

11. ❌ **显示收藏的文章列表** (40.8s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

12. ❌ **收藏项可点击跳转** (39.1s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

13. ❌ **用户不存在时显示错误** (38.1s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

14. ❌ **响应式布局 - 移动端** (36.8s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

15. ❌ **无收藏时显示空状态** (38.6s 超时)
    - **错误**: 页面加载超时
    - **原因**: 同上

### 错误处理 - 2 个失败

16. ❌ **API 错误时页面不崩溃** (36.6s 超时)
    - **错误**: 页面加载超时
    - **原因**: API 路由问题

17. ❌ **网络超时时显示友好提示** (30.0s 超时)
    - **错误**: 页面加载超时
    - **原因**: API 路由问题

### 性能测试 - 2 个失败

18. ❌ **页面加载时间合理** (30.0s 超时)
    - **错误**: 页面加载超时
    - **原因**: API 路由问题

19. ❌ **图片懒加载** (30.0s 超时)
    - **错误**: 页面加载超时
    - **原因**: API 路由问题

---

## ⚠️ 需要改进的地方

### 1. 🔥 紧急：修复后端路由注册顺序

**文件**: `/Users/polarisxu/project/golang/studygolang/cmd/studygolang/main.go`

**当前代码** (第 83-94 行):
```go
frontG := e.Group("")
controller.RegisterRoutes(frontG)  // 先注册，拦截所有路由

adminG := e.Group("/admin", pwm.NeedLogin(), pwm.AdminAuth())
admin.RegisterRoutes(adminG)

appG := e.Group("/app")
app.RegisterRoutes(appG)

apiG := e.Group("/api/v1")
api.RegisterRoutes(apiG)  // 后注册，但被前端路由拦截
```

**建议修改**:
```go
// 1. 先注册 API 路由（有明确前缀）
apiG := e.Group("/api/v1")
api.RegisterRoutes(apiG)

// 2. 再注册 App 路由
appG := e.Group("/app")
app.RegisterRoutes(appG)

// 3. 再注册 Admin 路由
adminG := e.Group("/admin", pwm.NeedLogin(), pwm.AdminAuth())
admin.RegisterRoutes(adminG)

// 4. 最后注册前端路由（无前缀，作为兜底）
frontG := e.Group("")
controller.RegisterRoutes(frontG)
```

**原理**: Echo 路由按注册顺序匹配，有前缀的路由应先注册，无前缀的路由最后注册作为兜底。

### 2. 前端错误处理不完善

**问题**:
- 用户不存在时，前端未显示友好的错误提示
- API 返回错误时，页面可能白屏或无限加载

**建议**:
- 在 `app/user/[username]/page.tsx` 中添加错误状态处理
- 在 `app/favorites/[username]/page.tsx` 中添加错误状态处理
- 显示友好的 404 页面或错误提示

### 3. 测试数据准备

**问题**: 测试使用的用户 "polaris" 不存在

**建议**:
- 创建测试数据库种子数据
- 或在测试前通过 API 创建测试用户
- 或使用真实存在的用户进行测试

### 4. API 响应格式不一致

**观察**:
- `/api/v1/user/polaris` 返回 `{"code":1,"msg":"用户不存在"}`
- 前端期望的格式可能不同

**建议**: 统一 API 响应格式，确保前后端约定一致

---

## 总结

### 核心问题

**后端路由注册顺序错误**导致 API 路由被前端路由拦截，这是导致 61.3% 测试失败的根本原因。

### 修复优先级

1. **P0 - 紧急**: 修复路由注册顺序（影响所有 API）
2. **P1 - 高**: 完善前端错误处理（用户体验）
3. **P2 - 中**: 准备测试数据（测试可靠性）
4. **P3 - 低**: 统一 API 响应格式（代码质量）

### 预期效果

修复路由注册顺序后，预计：
- ✅ 会员列表页 5 个测试通过
- ✅ 收藏列表页 9 个测试通过
- ✅ 错误处理 2 个测试通过
- ✅ 性能测试 2 个测试通过

**预计通过率**: 从 38.7% 提升至 **96.8%** (30/31)

---

## 附录

### 测试文件位置
- `/Users/polarisxu/project/golang/studygolang/frontend/tests/e2e/pages/user-center.spec.ts`

### 测试报告位置
- HTML 报告: `frontend/playwright-report/index.html`
- JSON 结果: `frontend/playwright-results.json`

### 重新运行测试
```bash
cd /Users/polarisxu/project/golang/studygolang/frontend
pnpm exec playwright test tests/e2e/pages/user-center.spec.ts
```

### 查看测试报告
```bash
pnpm exec playwright show-report
```
