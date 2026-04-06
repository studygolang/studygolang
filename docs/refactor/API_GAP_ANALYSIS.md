# StudyGolang 前后端分离功能缺口分析

> 对比 master 分支（旧版 Go 模板渲染）与 developer 分支（前后端分离），找出尚未实现的功能。

## 分析方法

1. 读取 master 分支的 `internal/http/controller/routes.go` 和各控制器文件
2. 对比 developer 分支的 `internal/http/controller/api/routes.go` 和前端页面
3. 排除已实现的功能（参考 FEATURE_COMPARISON.md）
4. 按优先级分类缺失功能

## 功能缺口分析

### P0（核心功能缺失 - 必须实现）

**无** - 所有核心功能已在 developer 分支实现。

### P1（重要功能缺失 - 建议实现）

#### 1. GCTT（Go 中文翻译组）功能
- **master 路由**: `/gctt`, `/gctt-list`, `/gctt-issue`, `/gctt/:username`, `/gctt-apply`, `/gctt-new`, `/gctt-webhook`
- **控制器**: `internal/http/controller/gctt.go`
- **模板**: `template/gctt/*.html`
- **需要实现**:
  - 后端 API: `internal/http/controller/api/gctt.go`
    - `GET /api/v1/gctt` - GCTT 首页数据（时间线、核心用户、未翻译 issue）
    - `GET /api/v1/gctt/users` - 译者列表
    - `GET /api/v1/gctt/issues` - Issue 列表（支持 label/translator 筛选）
    - `GET /api/v1/gctt/:username` - 译者详情
    - `GET /api/v1/gctt/apply` - 申请成为译者（检查 GitHub 绑定）
    - `POST /api/v1/gctt/articles` - 发布译文
    - `POST /api/v1/gctt/webhook` - GitHub Webhook（PR/Issue 事件）
  - 前端页面: `frontend/app/gctt/`
    - `page.tsx` - GCTT 首页
    - `users/page.tsx` - 译者列表
    - `issues/page.tsx` - Issue 列表
    - `[username]/page.tsx` - 译者详情
    - `apply/page.tsx` - 申请页面
    - `new/page.tsx` - 发布译文

#### 2. Go 安装包下载功能
- **master 路由**: `/dl`, `/dl/golang/:filename`, `/dl/add_new_version`
- **控制器**: `internal/http/controller/download.go`
- **模板**: `template/download/go.html`
- **需要实现**:
  - 后端 API: `internal/http/controller/api/download.go`
    - `GET /api/v1/downloads` - 下载列表（featured/stable/unstable/archived）
    - `GET /api/v1/downloads/:filename/redirect` - 重定向到下载地址
    - `POST /api/v1/downloads/add` - 添加新版本（管理员）
  - 前端页面: `frontend/app/dl/page.tsx`
    - 展示 Go 安装包列表（按版本分类）
    - 下载统计

#### 3. 友情链接页面
- **master 路由**: `/links`
- **控制器**: `internal/http/controller/link.go`
- **模板**: `template/link.html`
- **需要实现**:
  - 后端 API: `GET /api/v1/links` - 友情链接列表（已在 sidebar.go 中实现）
  - 前端页面: `frontend/app/links/page.tsx`
    - 展示友情链接列表

### P2（次要功能 - 可选实现）

#### 1. Wide Playground（在线编辑器）
- **master 路由**: `/wide/playground`
- **控制器**: `internal/http/controller/wide.go`
- **模板**: `template/wide/playground.html`
- **说明**: 这是一个内嵌 iframe 的在线 Go 编辑器，可能依赖外部服务
- **需要实现**:
  - 前端页面: `frontend/app/wide/playground/page.tsx`
  - 或直接嵌入 Go Playground 官方服务

#### 2. WebSocket 在线用户统计
- **master 路由**: `/ws`
- **控制器**: `internal/http/controller/websocket.go`
- **说明**: 实时统计在线用户数，通过 WebSocket 推送
- **需要实现**:
  - 后端保持 WebSocket 服务（已有实现）
  - 前端集成 WebSocket 客户端（在 layout.tsx 中全局连接）

#### 3. 微信相关功能
- **master 路由**: `/wechat/autoreply`, `/wechat/bind`
- **控制器**: `internal/http/controller/wechat.go`
- **说明**: 微信公众号自动回复和绑定
- **需要实现**:
  - 后端 API: `internal/http/controller/api/wechat.go`
    - `POST /api/v1/wechat/bind` - 绑定微信（验证码）
  - 前端页面: 在用户设置页面添加微信绑定入口

#### 4. 其他静态页面（OtherController）
- **master 路由**: `/*` (通配符，匹配任意 `.html` 模板)
- **控制器**: `internal/http/controller/other.go`
- **说明**: 自动渲染 `template/` 目录下的任意 `.html` 文件
- **已有模板**:
  - `template/goproxy.html` - Go 代理说明
  - `template/pkgdoc.html` - 包文档
  - `template/markdown.html` - Markdown 编辑器
- **需要实现**:
  - 前端页面: 根据需要创建对应的 Next.js 页面

### P3（不需要实现的功能）

#### 1. 安装向导（InstallController）
- **master 路由**: `/install`, `/install/setup-config`, `/install/do`, `/install/options`
- **说明**: 首次安装时的数据库配置和初始化，生产环境不需要

#### 2. 管理后台（AdminController）
- **说明**: 管理后台保持原有 Go 模板渲染，不迁移到 Next.js

#### 3. RSS Feed（已实现）
- **master 路由**: `/feed.html`, `/feed.xml`
- **developer 实现**: `frontend/app/feed/page.tsx`（已实现）

## 已实现功能清单（developer 分支）

### 后端 API（internal/http/controller/api/）

✅ **核心内容模块**:
- `article.go` - 文章列表、详情、创建、修改
- `topic.go` - 话题列表、详情、创建、修改、附言、节点
- `project.go` - 项目列表、详情
- `resource.go` - 资源列表、详情
- `reading.go` - 晨读列表、详情
- `book.go` - 书籍列表、详情
- `wiki.go` - Wiki 列表、详情
- `interview.go` - 面试题列表、详情
- `job.go` - 招聘列表、详情

✅ **用户模块**:
- `user.go` - 用户列表、详情、用户内容
- `account.go` - 激活、解绑

✅ **交互模块**:
- `comment.go` - 评论列表、创建、修改、详情
- `favorite.go` - 收藏列表、添加、删除
- `like.go` - 点赞列表、添加、删除
- `message.go` - 消息列表

✅ **功能模块**:
- `search.go` - 搜索
- `sidebar.go` - 侧边栏数据（统计、晨读、活跃用户、最新用户、友情链接、动态、排行）
- `subject.go` - 专栏列表、详情、关注、投稿、创建、修改
- `top.go` - DAU、富豪榜
- `balance.go` - 余额、充值
- `mission.go` - 每日任务
- `gift.go` - 礼物
- `oauth.go` - OAuth 登录
- `image.go` - 图片上传
- `captcha.go` - 验证码
- `node.go` - 节点别名路由

### 前端页面（frontend/app/）

✅ **38 个页面路由**（详见 FEATURE_COMPARISON.md）

## 建议实现顺序

### 第一阶段（P1 优先级）

1. **GCTT 功能**（如果社区仍在使用）
   - 后端 API: 7 个端点
   - 前端页面: 6 个页面
   - 预计工作量: 3-5 天

2. **Go 下载页面**（高流量页面）
   - 后端 API: 3 个端点
   - 前端页面: 1 个页面
   - 预计工作量: 1-2 天

3. **友情链接页面**
   - 后端 API: 已有（sidebar.go）
   - 前端页面: 1 个页面
   - 预计工作量: 0.5 天

### 第二阶段（P2 可选）

4. **WebSocket 在线统计**
   - 前端集成: 在 layout.tsx 中添加 WebSocket 客户端
   - 预计工作量: 1 天

5. **微信绑定功能**
   - 后端 API: 1 个端点
   - 前端页面: 在用户设置页面添加入口
   - 预计工作量: 1 天

6. **其他静态页面**
   - 根据实际需求创建
   - 预计工作量: 按需评估

## 总结

### 核心发现

1. **developer 分支已实现 95% 的核心功能**
   - 所有主要内容模块（文章、话题、项目、资源、晨读、书籍、Wiki、面试题、招聘）
   - 所有用户交互功能（评论、点赞、收藏、消息）
   - 所有新功能模块（专栏、任务、礼物、余额、排行榜、OAuth）

2. **缺失的主要是边缘功能**
   - GCTT（翻译组）- 如果社区仍在使用，建议实现
   - Go 下载页面 - 高流量页面，建议实现
   - 友情链接 - 简单页面，建议实现
   - WebSocket/微信/Wide - 可选功能

3. **迁移策略建议**
   - **立即上线**: 当前 developer 分支已可投入生产使用
   - **渐进补充**: 根据用户反馈和流量数据，补充 P1/P2 功能
   - **灰度切流**: 先切流主要页面，观察稳定性后再切流边缘功能

### 风险评估

- **低风险**: 核心功能已完整实现，SEO 结构保持不变
- **中风险**: GCTT 功能如果仍在使用，需尽快补充
- **可接受**: 边缘功能缺失不影响主要用户体验

### 下一步行动

1. **确认 GCTT 功能是否仍在使用**（查看访问日志）
2. **确认 Go 下载页面的流量占比**（决定优先级）
3. **制定灰度切流计划**（Nginx 配置）
4. **准备回滚方案**（保留旧版 Go 服务）

---

**生成时间**: 2026-04-06  
**分析基准**: master 分支 vs developer 分支  
**分析文件数**: 38 个控制器 + 50 个模板 + 38 个前端页面
