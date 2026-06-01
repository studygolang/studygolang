import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "激活账号 - Go语言中文网",
  description: "激活你的 Go语言中文网账号",
  robots: { index: false, follow: false },
}

export default function ActivateLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
