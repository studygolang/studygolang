import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "申请 GCTT - Go语言中文网",
  description: "申请成为 GCTT（Go Chinese Translation Team）译者",
  robots: { index: false, follow: false },
}

export default function GcttApplyLayout({ children }: { children: React.ReactNode }) {
  return children
}
