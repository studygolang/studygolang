# StudyGolang 前端项目

基于 Next.js 16 App Router 的 Go 语言中文网前端重构项目。

## 技术栈

- **框架**: Next.js 16 (App Router)
- **UI**: React 19 + Tailwind CSS v4 + shadcn/ui
- **渲染策略**: SSR (列表页) + ISR (详情页)
- **后端**: Go + Echo v4 (位于父目录)

## 环境变量配置

### 开发环境

复制 `.env.local` 文件并根据需要修改：

```bash
# Go 后端 API 地址（客户端浏览器调用）
NEXT_PUBLIC_API_BASE_URL=http://localhost:8090

# 内部 API 调用地址（SSR 服务端到服务端）
API_BASE_URL=http://localhost:8090

# 站点域名
NEXT_PUBLIC_SITE_URL=http://localhost:3000
```

### 生产环境

**⚠️ 重要**: 生产环境必须设置以下环境变量：

```bash
# 后端实际域名（用于客户端 API 调用和 /admin 跳转）
NEXT_PUBLIC_API_BASE_URL=https://studygolang.com

# 服务端 API 调用地址（可以是内网地址）
API_BASE_URL=http://backend-service:8090

# 站点域名
NEXT_PUBLIC_SITE_URL=https://studygolang.com
```

**说明**:
- `NEXT_PUBLIC_API_BASE_URL` 会暴露给浏览器，用于：
  - 客户端 API 调用（发布、点赞、评论等）
  - `/admin` 页面跳转到 Go 后端管理后台
  - 必须设置为后端实际可访问域名
- `API_BASE_URL` 仅用于服务端 SSR，可以使用内网地址

## 本地开发

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build

# 启动生产服务器
pnpm start
```

## 项目结构

```
frontend/
├── app/                    # Next.js App Router 页面
│   ├── (auth)/            # 认证相关页面
│   ├── articles/          # 文章列表和详情
│   ├── topics/            # 话题列表和详情
│   ├── projects/          # 项目列表
│   ├── interview/         # 面试题列表和详情
│   ├── publish/           # 发布页面
│   └── admin/             # 管理后台跳转
├── components/            # UI 组件
│   ├── ui/               # shadcn/ui 基础组件
│   └── ...               # 业务组件
├── lib/                   # 工具函数和类型定义
│   ├── api.ts            # API 客户端
│   └── types.ts          # TypeScript 类型定义
└── public/               # 静态资源
```

## SEO 策略

| 页面类型 | 渲染方式 | 说明 |
|---------|---------|------|
| 列表页 | SSR (no-store) | 实时内容，需要索引 |
| 详情页 | ISR (revalidate: 60) | 静态优先，SEO 最佳 |
| 用户页 | SSR | 动态内容 |
| 管理后台 | 跳转到 Go 后端 | 不需要 SEO |

## 认证机制

- **前端**: Token 认证（localStorage）
- **管理后台**: Go Session 认证（保持原有体系）
- `/admin` 路由直接跳转到 Go 后端管理页面

## API 代理

`next.config.mjs` 配置了以下代理规则：

```javascript
rewrites: [
  { source: '/api/v1/:path*', destination: 'http://localhost:8090/api/v1/:path*' },
  { source: '/static/:path*', destination: 'http://localhost:8090/static/:path*' }
]
```

## 部署注意事项

1. 确保 `NEXT_PUBLIC_API_BASE_URL` 设置为后端实际域名
2. 配置 Nginx 反向代理，将 `/api/v1` 和 `/static` 路径转发到 Go 后端
3. 管理后台 `/admin` 路由会直接跳转到 Go 后端，确保后端服务可访问
4. 生产环境建议使用 CDN 加速静态资源

## 相关文档

- [项目 AI 上下文](../CLAUDE.md)
- [Next.js 文档](https://nextjs.org/docs)
- [shadcn/ui 文档](https://ui.shadcn.com)
