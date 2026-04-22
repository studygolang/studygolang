import type { Metadata, Viewport } from 'next'
import { Analytics } from '@vercel/analytics/next'
import { Toaster } from '@/components/ui/sonner'
import { Providers } from '@/components/providers'
import './globals.css'

const CDN_DOMAIN = process.env.NEXT_PUBLIC_CDN_DOMAIN || "https://static.golangjob.cn"

export const metadata: Metadata = {
  title: 'Go语言中文网 - Golang中文社区',
  description: '中国最大的 Go 语言社区，探讨主题、阅读文章、分享项目与资源，助力 Gopher 成长',
  icons: {
    icon: `${CDN_DOMAIN}/static/img/favicon.ico`,
    shortcut: `${CDN_DOMAIN}/static/img/favicon.ico`,
  },
}

export const viewport: Viewport = {
  themeColor: '#1a9ca0',
  width: 'device-width',
  initialScale: 1,
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="zh-CN">
      <body className="font-sans antialiased">
        <Providers>
          {children}
          <Toaster position="top-center" richColors />
          <Analytics />
        </Providers>
      </body>
    </html>
  )
}
