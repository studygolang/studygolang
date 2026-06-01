import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "管理后台 - Go语言中文网",
  robots: { index: false, follow: false },
}

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return children
}
