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
  // 开发时将 /api/v1/* 代理到 Go 后端（避免浏览器 CORS）
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: `${process.env.API_BASE_URL || 'http://localhost:8088'}/api/v1/:path*`,
      },
    ]
  },
}

export default nextConfig
