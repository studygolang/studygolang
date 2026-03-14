"use client"

import { useEffect, useState } from "react"

/**
 * 管理后台入口：先同步 session 到后端，再跳转到 Go 后端管理页面
 */
export default function AdminRedirectPage() {
  const [error, setError] = useState<string>("")

  useEffect(() => {
    const backend = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

    // 先调用 sync-session API 在后端建立 session
    fetch(`/api/v1/user/sync-session`, {
      method: 'POST',
      credentials: 'include',
    })
      .then(res => res.json())
      .then(data => {
        if (data.code === 0) {
          // Session 同步成功，跳转到管理后台
          window.location.replace(`${backend}/admin`)
        } else {
          // 未登录或 token 过期，跳转到登录页
          window.location.replace(`/account/login?redirect=/admin`)
        }
      })
      .catch(err => {
        console.error('同步 session 失败:', err)
        setError('同步登录状态失败，请稍后重试')
      })
  }, [])

  return (
    <div className="flex min-h-screen items-center justify-center text-sm text-muted-foreground">
      {error || "正在跳转到管理后台…"}
    </div>
  )
}
