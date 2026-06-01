import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "我的积分 - Go语言中文网",
  description: "查看积分余额、收支明细和兑换记录",
  robots: { index: false, follow: false },
}

export default function BalanceLayout({ children }: { children: React.ReactNode }) {
  return children
}
