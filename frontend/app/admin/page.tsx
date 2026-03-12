"use client"

import { useEffect } from "react"

/**
 * 管理后台入口：直接跳转到 Go 后端管理页面（保留原有 session 认证体系）
 */
export default function AdminRedirectPage() {
  useEffect(() => {
    const backend = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"
    window.location.replace(`${backend}/admin`)
  }, [])

  return (
    <div className="flex min-h-screen items-center justify-center text-sm text-muted-foreground">
      正在跳转到管理后台…
    </div>
  )
}
