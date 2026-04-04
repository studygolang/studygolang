import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "OAuth 登录 - Go语言中文网",
  robots: { index: false, follow: false },
}

export default function OAuthCallbackLayout({ children }: { children: React.ReactNode }) {
  return children
}
