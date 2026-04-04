"use client"

import { useEffect, useState, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { Card, CardContent } from "@/components/ui/card"
import { Loader2, CheckCircle2, XCircle } from "lucide-react"

interface OAuthResult {
  action: "login" | "bind"
  username?: string
  balance?: number
  message?: string
}

function OAuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading")
  const [result, setResult] = useState<OAuthResult | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const code = searchParams.get("code")
    const provider = window.location.pathname.includes("github") ? "github" : "gitea"

    if (!code) {
      setStatus("error")
      setError("授权码缺失")
      return
    }

    // 调用后端回调 API
    const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"
    fetch(`${apiBase}/api/v1/oauth/${provider}/callback?code=${code}`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((json) => {
        if (json.code === 0) {
          setStatus("success")
          setResult(json.data)
          // 2秒后跳转
          setTimeout(() => {
            const redirect = searchParams.get("redirect_url") || "/"
            router.push(redirect)
          }, 2000)
        } else {
          setStatus("error")
          setError(json.msg || "登录失败")
        }
      })
      .catch((e) => {
        setStatus("error")
        setError(e.message || "网络错误")
      })
  }, [router, searchParams])

  return (
    <PageLayout sidebar={false}>
      <Card className="mx-auto max-w-md">
        <CardContent className="flex flex-col items-center py-12">
          {status === "loading" && (
            <>
              <Loader2 className="mb-4 h-12 w-12 animate-spin text-primary" />
              <p className="text-muted-foreground">正在处理登录...</p>
            </>
          )}
          {status === "success" && (
            <>
              <CheckCircle2 className="mb-4 h-12 w-12 text-green-500" />
              <p className="font-medium">
                {result?.action === "bind" ? "绑定成功" : "登录成功"}
              </p>
              {result?.username && (
                <p className="mt-1 text-sm text-muted-foreground">
                  欢迎回来，{result.username}
                </p>
              )}
              <p className="mt-2 text-xs text-muted-foreground">
                即将跳转...
              </p>
            </>
          )}
          {status === "error" && (
            <>
              <XCircle className="mb-4 h-12 w-12 text-destructive" />
              <p className="font-medium text-destructive">登录失败</p>
              <p className="mt-1 text-sm text-muted-foreground">{error}</p>
            </>
          )}
        </CardContent>
      </Card>
    </PageLayout>
  )
}

export default function OAuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <PageLayout sidebar={false}>
          <Card className="mx-auto max-w-md">
            <CardContent className="flex flex-col items-center py-12">
              <Loader2 className="mb-4 h-12 w-12 animate-spin text-primary" />
              <p className="text-muted-foreground">加载中...</p>
            </CardContent>
          </Card>
        </PageLayout>
      }
    >
      <OAuthCallbackContent />
    </Suspense>
  )
}