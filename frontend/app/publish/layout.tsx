import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "发布内容 - Go语言中文网",
  description: "在 Go语言中文网发布话题、文章、项目或资源，分享你的 Go 语言经验",
}

export default function PublishLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
