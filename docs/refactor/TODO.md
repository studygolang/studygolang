# StudyGolang 前后端分离重构 - 待办事项清单

> 创建时间: 2026-03-14
> 目标: 完成前后端分离重构，保持所有功能不丢失，界面风格统一

## 开发原则

1. **界面风格**: 严格参考 `frontend/` 项目现有实现，保持风格一致
2. **功能完整性**: 完全参考旧版本，确保功能不丢失
3. **SEO 优先**: 列表页 SSR (no-store)，详情页 ISR (revalidate: 60)
4. **测试驱动**: 所有功能开发完成后必须通过 E2E 测试

---

## 🔴 P0 - 核心用户功能 (必须实现)

### 1. 用户个人中心模块

#### 1.1 用户内容列表页面
- [ ] **用户话题列表** `/user/:username/topics`
  - 后端 API: `GET /api/v1/users/:username/topics`
  - 前端页面: `frontend/app/user/[username]/topics/page.tsx`
  - 功能: 展示用户发布的所有话题，支持分页
  - 参考: `internal/http/controller/user.go` Topics 方法

- [ ] **用户文章列表** `/user/:username/articles`
  - 后端 API: `GET /api/v1/users/:username/articles`
  - 前端页面: `frontend/app/user/[username]/articles/page.tsx`
  - 功能: 展示用户发布的所有文章，支持分页
  - 参考: `internal/http/controller/user.go` Articles 方法

- [ ] **用户资源列表** `/user/:username/resources`
  - 后端 API: `GET /api/v1/users/:username/resources`
  - 前端页面: `frontend/app/user/[username]/resources/page.tsx`
  - 功能: 展示用户分享的所有资源，支持分页
  - 参考: `internal/http/controller/user.go` Resources 方法

- [ ] **用户项目列表** `/user/:username/projects`
  - 后端 API: `GET /api/v1/users/:username/projects`
  - 前端页面: `frontend/app/user/[username]/projects/page.tsx`
  - 功能: 展示用户的所有项目，支持分页
  - 参考: `internal/http/controller/user.go` Projects 方法

- [ ] **用户评论列表** `/user/:username/comments`
  - 后端 API: `GET /api/v1/users/:username/comments`
  - 前端页面: `frontend/app/user/[username]/comments/page.tsx`
  - 功能: 展示用户的评论历史，支持分页
  - 参考: `internal/http/controller/user.go` Comments 方法

#### 1.2 会员列表页面
- [ ] **会员列表** `/users`
  - 后端 API: `GET /api/v1/users` (已有，需验证)
  - 前端页面: `frontend/app/users/page.tsx`
  - 功能: 展示活跃会员 + 新加入会员，显示会员总数
  - 参考: `internal/http/controller/user.go` ReadList 方法

### 2. 个人设置模块

- [ ] **个人设置页面** `/account/edit`
  - 后端 API:
    - `GET /api/v1/account/profile` (获取个人信息)
    - `POST /api/v1/account/profile` (更新个人信息)
  - 前端页面: `frontend/app/account/edit/page.tsx`
  - 功能: 修改昵称、邮箱、个人简介、城市等
  - 参考: `internal/http/controller/account.go` Edit 方法

- [ ] **修改密码** `/account/changepwd`
  - 后端 API: `POST /api/v1/account/password`
  - 前端页面: `frontend/app/account/changepwd/page.tsx`
  - 功能: 修改登录密码，需验证旧密码
  - 参考: `internal/http/controller/account.go` ChangePwd 方法

- [ ] **修改头像** `/account/change_avatar`
  - 后端 API: `POST /api/v1/account/avatar`
  - 前端页面: 集成到 `/account/edit` 页面
  - 功能: 上传头像图片
  - 参考: `internal/http/controller/account.go` ChangeAvatar 方法

---

## 🟠 P1 - 核心交互功能

### 3. 消息系统模块

- [ ] **消息列表页面** `/message/:msgtype`
  - 后端 API: `GET /api/v1/messages?type=system|inbox|outbox`
  - 前端页面: `frontend/app/message/[msgtype]/page.tsx`
  - 功能: 系统消息、收件箱、发件箱，支持分页
  - 参考: `internal/http/controller/message.go` ReadList 方法

- [ ] **发送私信** `/message/send`
  - 后端 API: `POST /api/v1/messages`
  - 前端页面: `frontend/app/message/send/page.tsx`
  - 功能: 发送私信给指定用户
  - 参考: `internal/http/controller/message.go` Send 方法

- [ ] **删除消息**
  - 后端 API: `DELETE /api/v1/messages/:id`
  - 功能: 删除消息
  - 参考: `internal/http/controller/message.go` Delete 方法

### 4. 收藏功能

- [ ] **收藏列表** `/favorites/:username`
  - 后端 API: `GET /api/v1/users/:username/favorites`
  - 前端页面: `frontend/app/favorites/[username]/page.tsx`
  - 功能: 展示用户收藏的内容（文章、话题、资源等）
  - 参考: `internal/http/controller/favorite.go` ReadList 方法

- [ ] **收藏/取消收藏**
  - 后端 API: `POST /api/v1/favorites/:objid`
  - 功能: 收藏或取消收藏内容
  - 参考: `internal/http/controller/favorite.go` Create 方法

### 5. 点赞功能

- [ ] **点赞/取消点赞**
  - 后端 API: `POST /api/v1/likes/:objid`
  - 功能: 点赞或取消点赞内容
  - 参考: `internal/http/controller/like.go`

### 6. 评论功能完善

- [ ] **评论组件优化**
  - 前端组件: `frontend/components/comment-section.tsx`
  - 功能:
    - 评论列表展示
    - 发表评论
    - 回复评论
    - 点赞评论
  - 参考: 各详情页的评论区

### 7. 内容编辑功能

- [ ] **文章修改** `/articles/modify/:id`
  - 后端 API:
    - `GET /api/v1/articles/:id/edit` (获取编辑数据)
    - `PUT /api/v1/articles/:id` (更新文章)
  - 前端页面: `frontend/app/articles/modify/[id]/page.tsx`
  - 功能: 编辑已发布的文章
  - 参考: `internal/http/controller/article.go` Modify 方法

- [ ] **话题修改** `/topics/modify/:id`
  - 后端 API:
    - `GET /api/v1/topics/:id/edit` (获取编辑数据)
    - `PUT /api/v1/topics/:id` (更新话题)
  - 前端页面: `frontend/app/topics/modify/[id]/page.tsx`
  - 功能: 编辑已发布的话题
  - 参考: `internal/http/controller/topic.go` Modify 方法

---

## 🟡 P2 - 次要功能

### 8. 其他功能

- [ ] **邮件退订** `/user/email/unsubscribe`
  - 后端 API: `GET/POST /api/v1/account/email/unsubscribe`
  - 前端页面: `frontend/app/user/email/unsubscribe/page.tsx`
  - 参考: `internal/http/controller/user.go` EmailUnsub 方法

- [ ] **社交账号解绑** `/account/social/unbind`
  - 后端 API: `POST /api/v1/account/social/unbind`
  - 参考: `internal/http/controller/account.go` Unbind 方法

- [ ] **任务系统** `/mission/*`
  - 参考: `internal/http/controller/mission.go`

- [ ] **积分商城** `/gift/*`
  - 参考: `internal/http/controller/gift.go`

- [ ] **排行榜** `/top/*`
  - 参考: `internal/http/controller/top.go`

- [ ] **专题页** `/subject/*`
  - 参考: `internal/http/controller/subject.go`

---

## 开发分工建议

### 阶段一: 用户中心 (本周)
- **Agent 1**: 用户内容列表页面 (topics, articles, resources, projects, comments)
- **Agent 2**: 个人设置模块 (edit, changepwd, avatar)
- **Agent 3**: 会员列表页面

### 阶段二: 交互功能 (下周)
- **Agent 4**: 消息系统 (列表、发送、删除)
- **Agent 5**: 收藏功能 + 点赞功能
- **Agent 6**: 评论UI完善 + 内容编辑功能

### 阶段三: 次要功能 (后续)
- **Agent 7**: 其他次要功能

---

## 测试要求

### E2E 测试覆盖

每个功能模块开发完成后，必须通过以下测试：

1. **用户中心测试**
   - 访问用户个人主页
   - 查看用户各类内容列表
   - 修改个人信息
   - 修改密码
   - 上传头像

2. **消息系统测试**
   - 查看系统消息
   - 发送私信
   - 接收私信
   - 删除消息

3. **交互功能测试**
   - 收藏内容
   - 取消收藏
   - 点赞内容
   - 取消点赞
   - 发表评论
   - 回复评论

4. **内容编辑测试**
   - 编辑文章
   - 编辑话题
   - 保存修改

### 测试通过标准

- ✅ 所有页面正常渲染
- ✅ 所有 API 调用成功
- ✅ 数据正确展示
- ✅ 交互功能正常
- ✅ 错误处理完善
- ✅ 界面风格一致

---

## 进度跟踪

- [ ] 阶段一完成
- [ ] 阶段二完成
- [ ] 阶段三完成
- [ ] E2E 测试通过
- [ ] 上线部署

---

## 注意事项

1. **API 响应格式统一**: 使用 `{ code, data, message }` 格式
2. **错误处理**: 所有 API 必须有完善的错误处理
3. **认证鉴权**: 需要登录的功能必须验证 Cookie 认证
4. **SEO 优化**: 所有页面必须有 `generateMetadata`
5. **样式一致**: 使用 shadcn/ui 组件，保持与现有页面风格一致
6. **响应式设计**: 所有页面必须支持移动端
7. **性能优化**: 列表页支持分页，避免一次加载过多数据
