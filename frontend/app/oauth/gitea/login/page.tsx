"use client"

import { useEffect, useState, Suspense } from "react"
import { useSearchParams } from "next/navigation"
import { Loader2 } from "lucide-react"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

function GiteaLoginContent() {
  const searchParams = useSearchParams()
  const [error, setError] = useState("")

  useEffect(() => {
    const redirect = searchParams.get("redirect") || "/"
    fetch(`${API_BASE}/api/v1/oauth/gitea/url?uri=${encodeURIComponent(redirect)}`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.code === 0 && data.data?.url) {
          window.location.href = data.data.url
        } else {
          setError(data.msg || "获取 Gitea 授权地址失败")
        }
      })
      .catch(() => {
        setError("网络错误，请稍后重试")
      })
  }, [searchParams])

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-destructive">{error}</p>
          <a href="/account/login" className="mt-4 inline-block text-sm text-primary hover:underline">
            返回登录
          </a>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <Loader2 className="mx-auto mb-4 h-8 w-8 animate-spin text-muted-foreground" />
        <p className="text-sm text-muted-foreground">正在跳转到 Gitea 登录...</p>
      </div>
    </div>
  )
}

export default function GiteaLoginPage() {
  return (
    <Suspense fallback={<div className="flex min-h-screen items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /></div>}>
      <GiteaLoginContent />
    </Suspense>
  )
}
