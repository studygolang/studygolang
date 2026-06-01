import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "发送消息 - Go语言中文网",
  robots: { index: false, follow: false },
}

export default function SendMessageLayout({ children }: { children: React.ReactNode }) {
  return children
}
