import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "礼品兑换 - Go语言中文网",
  description: "用积分兑换精美礼品",
  robots: { index: false, follow: false },
}

export default function GiftLayout({ children }: { children: React.ReactNode }) {
  return children
}
