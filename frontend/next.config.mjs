/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    unoptimized: true,
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'studygolang.com',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
      },
    ],
  },
  // 旧版 URL 兼容：邮件、书签、搜索引擎里仍是老路径。
  // 灰度期用 302，切流稳定后再改 permanent。
  async redirects() {
    return [
      // 账号链路：找回密码邮件里的旧链接
      { source: '/account/forgetpwd', destination: '/account/forgot-password', permanent: false },
      { source: '/account/resetpwd', destination: '/account/reset-password', permanent: false },
      { source: '/user/email/unsubscribe', destination: '/account/unsubscribe', permanent: false },

      // 工具页
      { source: '/wide/playground', destination: '/wide', permanent: false },

      // GCTT
      { source: '/gctt-list', destination: '/gctt/users', permanent: false },
      { source: '/gctt-issue', destination: '/gctt/issues', permanent: false },
      { source: '/gctt-apply', destination: '/gctt/apply', permanent: false },

      // 签到任务：旧版分散在三个页面，新版合并到 /mission
      { source: '/mission/daily', destination: '/mission', permanent: false },
      { source: '/mission/daily/redeem', destination: '/mission', permanent: false },
      // 旧版 GET /mission/complete/1 直接完成「初始资本」任务；新版领取入口在 /balance
      { source: '/mission/complete/:id', destination: '/balance', permanent: false },

      // 面试题：旧版「每日一题」入口
      { source: '/interview/question', destination: '/interview', permanent: false },

      { source: '/gift/mine', destination: '/gift', permanent: false },
      { source: '/subject/my_articles', destination: '/subject/mine', permanent: false },
      { source: '/friend/links', destination: '/links', permanent: false },

      // 旧版独立的「最新/最热」列表页，新版统一为列表页
      { source: '/articles/recent', destination: '/articles', permanent: false },
      { source: '/projects/recent', destination: '/projects', permanent: false },
      { source: '/resources/recent', destination: '/resources', permanent: false },
      { source: '/readings/recent', destination: '/readings', permanent: false },
      { source: '/topics/recent', destination: '/topics', permanent: false },
      { source: '/topics/:nid/others', destination: '/topics/node/:nid', permanent: false },
      { source: '/nodes/hot', destination: '/nodes', permanent: false },
      { source: '/users/active', destination: '/users', permanent: false },
      { source: '/users/newest', destination: '/users', permanent: false },
    ]
  },
  // 将特定路径代理到 Go 后端
  async rewrites() {
    const backend = process.env.API_BASE_URL || 'http://localhost:8090'
    return [
      // API 代理（避免浏览器 CORS）
      {
        source: '/api/v1/:path*',
        destination: `${backend}/api/v1/:path*`,
      },
      // 管理后台静态资源（管理后台 HTML 中引用的 /static/ 路径）
      {
        source: '/static/:path*',
        destination: `${backend}/static/:path*`,
      },
      // Sitemap 代理（复用后端 logic/sitemap.go 生成的静态文件）
      {
        source: '/sitemap',
        destination: `${backend}/sitemap`,
      },
      {
        source: '/sitemap/:path*',
        destination: `${backend}/sitemap/:path*`,
      },
      // RSS Feed 代理；旧版订阅说明页交给新版页面渲染
      {
        source: '/feed.xml',
        destination: `${backend}/api/v1/feed`,
      },
      {
        source: '/feed.html',
        destination: '/feed',
      },
      // OAuth 回调：提供商后台登记的 redirect_uri 可能还是旧路径，
      // 用 rewrite 而非 redirect，避免多一跳导致 state/code 丢失。
      {
        source: '/oauth/:provider(github|gitea)/callback',
        destination: `${backend}/api/v1/oauth/:provider/callback`,
      },
      // 旧版 GET 注销：必须落到后端才能清掉 HttpOnly Cookie，
      // 前端 redirect 到登录页只会让用户以为自己退出了。
      {
        source: '/account/logout',
        destination: `${backend}/api/v1/user/logout`,
      },
      // 外链包装跳转（/link 加 utm、/wr 内嵌阅读），只有 Go 端有实现
      {
        source: '/link',
        destination: `${backend}/link`,
      },
      {
        source: '/wr',
        destination: `${backend}/wr`,
      },
    ]
  },
}

export default nextConfig