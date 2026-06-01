import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "账号设置 - Go语言中文网",
  description: "个人账号与隐私设置",
  robots: { index: false, follow: false },
}

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  return children
}
