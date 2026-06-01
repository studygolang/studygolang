import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑话题 - Go语言中文网",
  description: "编辑已发布的话题",
  robots: { index: false, follow: false },
}

export default function ModifyTopicLayout({ children }: { children: React.ReactNode }) {
  return children
}
