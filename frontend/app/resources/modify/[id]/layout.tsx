import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑资源 - Go语言中文网",
  description: "编辑你分享的 Go 学习资源",
  robots: { index: false, follow: false },
}

export default function ModifyResourceLayout({ children }: { children: React.ReactNode }) {
  return children
}
