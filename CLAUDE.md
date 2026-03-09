# StudyGolang 项目 AI 上下文

## 项目概述

Go 语言中文网社区 (https://studygolang.com)，十多年前开发，正在进行**前后端分离重构**。

- **后端**: Go + Echo v4 框架，已有 JSON API 层（`internal/http/controller/app/`）
- **前端**: Next.js 16 (App Router) + React 19 + Tailwind CSS v4 + shadcn/ui，位于 `frontend/` 目录

## 架构决策（重构核心）

### SEO 策略（关键）
> 社区内容站点，SEO 至关重要，不能牺牲。

| 页面类型 | 渲染方式 | 原因 |
|---------|---------|------|
| 列表页（文章/话题/项目） | Next.js SSR (`fetch` + no-store) | 实时内容，需索引 |
| 详情页（文章/话题详情） | Next.js SSG + ISR（revalidate=60） | 静态优先，SEO 最佳 |
| 用户个人页 | SSR | 动态内容 |
| 管理后台 | CSR（SPA） | 不需要 SEO |
| 交互操作（点赞/评论/发帖） | Client 调用 API | 无需 SEO |

### URL 结构（保持不变，保护 SEO）
```
/articles        → 文章列表
/articles/:id    → 文章详情
/topics          → 话题列表
/topics/:tid     → 话题详情
/projects        → 项目列表
/resources       → 资源列表
/readings        → 晨读列表
/books           → 书籍列表
/wiki            → Wiki
/search          → 搜索
```

### API 结构
- **现有 App API**: `GET /app/v1/...`（手机端 API，可复用）
- **新增 Web API**: `GET /api/v1/...`（专为 Next.js 前端设计，支持 SSR）
- **认证**: 现有 Token 机制（`GenToken`/`ValidateToken` in `internal/http/http.go`），Web 端用 Cookie Session

## 关键文件路径

### 后端
```
internal/http/controller/        # Web 控制器（渲染模板，逐步迁移）
internal/http/controller/app/    # App JSON API（可复用逻辑）
internal/http/controller/admin/  # 管理后台
internal/http/http.go            # 基础 HTTP 工具（Token、模板、Session）
internal/logic/                  # 业务逻辑层（核心，保持不变）
internal/model/                  # 数据模型
middleware/                      # 中间件
```

### 前端
```
frontend/                        # Next.js 16 项目根目录
frontend/app/                    # App Router 页面
frontend/components/             # UI 组件（v0.dev 生成）
frontend/lib/                    # 工具函数、API 客户端
```

## 重构分阶段计划

### Phase 1：API 层完善（后端）
- 新增 `internal/http/controller/api/` 目录
- 为 Next.js SSR 提供完整 REST API
- 复用 `app/` 层已有逻辑
- 统一 API 响应格式：`{ code, data, message }`

### Phase 2：前端 SSR 页面（Next.js）
- 列表页：文章、话题、项目、资源、晨读
- 详情页：文章详情、话题详情（ISR）
- SEO 组件：动态 `<head>` meta 标签

### Phase 3：前端交互功能（Client）
- 登录/注册表单
- 发帖/评论
- 点赞/收藏
- 用户个人中心

### Phase 4：切流与下线旧模板
- Nginx 反向代理配置（Next.js + Go 后端共存）
- 灰度切流
- 逐步下线 Go 模板渲染

## 开发规范

### API 响应格式
```go
// 统一使用此格式
type APIResponse struct {
    Code    int         `json:"code"`    // 0=成功, 非0=错误
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
}
```

### Next.js 数据获取规范
```typescript
// SSR 列表页
async function getArticles(page: number) {
  const res = await fetch(`${API_BASE}/api/v1/articles?p=${page}`, {
    cache: 'no-store'
  })
  return res.json()
}

// ISR 详情页
async function getArticleDetail(id: string) {
  const res = await fetch(`${API_BASE}/api/v1/articles/${id}`, {
    next: { revalidate: 60 }
  })
  return res.json()
}
```

### SEO 组件规范
```typescript
// 每个页面必须有动态 metadata
export async function generateMetadata({ params }) {
  const article = await getArticleDetail(params.id)
  return {
    title: `${article.title} - Go语言中文网`,
    description: article.summary,
    openGraph: { ... }
  }
}
```

## 当前状态
- 后端：`internal/http/controller/app/` 已有部分 JSON API（articles, topics, projects, resources, user, comments）
- 前端：`frontend/` 已有基础 UI 组件和页面骨架（v0.dev 生成）
- 需要：完善 API 层 → 连接前后端 → SSR 实现
