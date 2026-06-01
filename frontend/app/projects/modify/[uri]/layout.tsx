import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑项目 - Go语言中文网",
  description: "编辑你发布的开源项目信息",
  robots: { index: false, follow: false },
}

export default function ModifyProjectLayout({ children }: { children: React.ReactNode }) {
  return children
}
