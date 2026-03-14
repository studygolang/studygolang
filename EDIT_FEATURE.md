# 内容编辑功能开发完成

## 已实现功能

### 后端 API

#### 文章编辑
- `GET /api/v1/articles/:id/edit` - 获取文章编辑数据（需登录，验证权限）
- `PUT /api/v1/articles/:id` - 更新文章（需登录，验证权限）

#### 话题编辑
- `GET /api/v1/topics/:tid/edit` - 获取话题编辑数据（需登录，验证权限）
- `PUT /api/v1/topics/:tid` - 更新话题（需登录，验证权限）

### 前端页面

#### 文章编辑页面
- 路由: `/articles/modify/[id]`
- 文件: `frontend/app/articles/modify/[id]/page.tsx`
- 功能:
  - 加载现有文章内容
  - 编辑标题、内容、标签
  - 使用 BlockNote 富文本编辑器
  - 保存修改后跳转到文章详情页

#### 话题编辑页面
- 路由: `/topics/modify/[id]`
- 文件: `frontend/app/topics/modify/[id]/page.tsx`
- 功能:
  - 加载现有话题内容
  - 编辑标题、节点、内容、标签
  - 使用 BlockNote 富文本编辑器
  - 保存修改后跳转到话题详情页

## 权限验证

- 只能编辑自己的内容
- 管理员可以编辑所有内容
- 使用 Cookie 认证（sg_token）
- 未登录自动跳转到登录页

## 技术实现

### 后端
- 新增 `requireAuth()` 辅助函数统一处理认证
- 复用 `logic.CanEdit()` 进行权限验证
- 复用 `logic.DefaultArticle.Modify()` 和 `logic.DefaultTopic.Modify()` 进行更新
- 统一 API 响应格式 `{code, msg, data}`

### 前端
- 复用 `/publish` 页面的 BlockNote 编辑器组件
- 使用 Next.js App Router 动态路由
- 客户端渲染（CSR）
- 表单验证和错误提示
- 保存成功后自动跳转

## 测试建议

1. 登录后访问自己的文章/话题详情页
2. 点击"编辑"按钮（需要在详情页添加编辑按钮）
3. 修改内容后保存
4. 验证权限：尝试编辑他人的内容（应该被拒绝）
5. 验证未登录：退出登录后访问编辑页面（应该跳转到登录页）

## 后续工作

- 在文章/话题详情页添加"编辑"按钮（仅作者和管理员可见）
- 添加草稿保存功能
- 添加编辑历史记录
