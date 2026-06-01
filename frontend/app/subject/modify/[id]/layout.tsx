import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑主题 - Go语言中文网",
  description: "编辑已发布的主题",
  robots: { index: false, follow: false },
}

export default function ModifySubjectLayout({ children }: { children: React.ReactNode }) {
  return children
}
