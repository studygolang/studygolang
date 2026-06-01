import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑文章 - Go语言中文网",
  description: "编辑已发布的 Go 技术文章",
  robots: { index: false, follow: false },
}

export default function ModifyArticleLayout({ children }: { children: React.ReactNode }) {
  return children
}
