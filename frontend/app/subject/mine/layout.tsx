import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "我的主题 - Go语言中文网",
  robots: { index: false, follow: false },
}

export default function MySubjectsLayout({ children }: { children: React.ReactNode }) {
  return children
}
