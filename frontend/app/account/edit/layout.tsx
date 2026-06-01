import type { Metadata } from "next"
import type { ReactNode } from "react"

export const metadata: Metadata = {
  title: "编辑个人资料 - Go语言中文网",
  description: "修改你的个人资料、头像和签名",
  robots: { index: false, follow: false },
}

export default function EditAccountLayout({ children }: { children: ReactNode }) {
  return <>{children}</>
}
