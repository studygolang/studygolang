import type { Metadata, Viewport } from 'next'
import { Analytics } from '@vercel/analytics/next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Go\u8bed\u8a00\u4e2d\u6587\u7f51 - Golang\u4e2d\u6587\u793e\u533a',
  description: '\u4e2d\u56fd\u6700\u5927\u7684 Go \u8bed\u8a00\u793e\u533a\uff0c\u63a2\u8ba8\u4e3b\u9898\u3001\u9605\u8bfb\u6587\u7ae0\u3001\u5206\u4eab\u9879\u76ee\u4e0e\u8d44\u6e90\uff0c\u52a9\u529b Gopher \u6210\u957f',
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
        {children}
        <Analytics />
      </body>
    </html>
  )
}
