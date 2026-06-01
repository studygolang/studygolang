import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "忘记密码 - Go语言中文网",
  description: "通过邮箱重置你的 Go语言中文网密码",
  robots: { index: false, follow: false },
}

export default function ForgotPasswordLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
