import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "新建 Wiki - Go语言中文网",
  description: "为 Go 社区贡献 Wiki 内容",
  robots: { index: false, follow: false },
}

export default function NewWikiLayout({ children }: { children: React.ReactNode }) {
  return children
}
