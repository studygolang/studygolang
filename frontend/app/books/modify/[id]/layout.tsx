import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑图书 - Go语言中文网",
  description: "编辑 Go 相关图书信息",
  robots: { index: false, follow: false },
}

export default function ModifyBookLayout({ children }: { children: React.ReactNode }) {
  return children
}
