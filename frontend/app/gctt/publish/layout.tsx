import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "发布 GCTT 译文 - Go语言中文网",
  robots: { index: false, follow: false },
}

export default function GcttPublishLayout({ children }: { children: React.ReactNode }) {
  return children
}
