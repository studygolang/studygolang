import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "每日任务 - Go语言中文网",
  description: "完成每日任务，获取积分奖励",
  robots: { index: false, follow: false },
}

export default function MissionLayout({ children }: { children: React.ReactNode }) {
  return children
}
