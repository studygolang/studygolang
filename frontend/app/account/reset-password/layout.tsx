import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "重置密码 - Go语言中文网",
  description: "使用重置链接设置新的登录密码",
  robots: { index: false, follow: false },
}

export default function ResetPasswordLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
