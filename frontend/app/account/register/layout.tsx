import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "注册 - Go语言中文网",
  description: "注册成为 Go语言中文网会员，加入中国最大的 Go 语言开发者社区",
}

export default function RegisterLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
