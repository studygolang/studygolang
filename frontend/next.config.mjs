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
	      // RSS/Atom Feed 代理
	      {
	        source: '/feed.xml',
	        destination: `${backend}/api/v1/feed`,
	      },
	      {
	        source: '/feed.html',
	        destination: `${backend}/api/v1/feed`,
	      },
	    ]
	  },
}

export default nextConfig
