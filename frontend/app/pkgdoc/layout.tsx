import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Go语言标准库文档中文版 | Go语言中文网",
}

export default function PkgDocLayout({ children }: { children: React.ReactNode }) {
  return children
}
