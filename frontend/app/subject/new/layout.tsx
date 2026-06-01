import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "创建主题 - Go语言中文网",
  description: "创建一个新的主题，召集志同道合的 Go 开发者一起讨论",
  robots: { index: false, follow: false },
}

export default function NewSubjectLayout({ children }: { children: React.ReactNode }) {
  return children
}
