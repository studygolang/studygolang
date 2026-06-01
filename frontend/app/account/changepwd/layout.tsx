import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "修改密码 - Go语言中文网",
  description: "修改你的账户登录密码",
  robots: { index: false, follow: false },
}

export default function ChangePwdLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
