import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "登录 - Go语言中文网",
  description: "登录 Go语言中文网，加入中国最大的 Go 语言开发者社区",
}

export default function LoginLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
