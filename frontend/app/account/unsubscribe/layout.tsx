import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "退订邮件 - Go语言中文网",
  description: "取消订阅邮件通知",
  robots: { index: false, follow: false },
}

export default function UnsubscribeLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
