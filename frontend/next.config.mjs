/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    ignoreBuildErrors: true,
  },
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
    ]
  },
}

export default nextConfig
