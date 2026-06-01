import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "消息中心 - Go语言中文网",
  description: "查看系统消息、私信和通知",
  robots: { index: false, follow: false },
}

export default function MessageLayout({ children }: { children: React.ReactNode }) {
  return children
}
