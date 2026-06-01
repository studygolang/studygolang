import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "编辑 Wiki - Go语言中文网",
  description: "编辑 Wiki 词条",
  robots: { index: false, follow: false },
}

export default function EditWikiLayout({ children }: { children: React.ReactNode }) {
  return children
}
