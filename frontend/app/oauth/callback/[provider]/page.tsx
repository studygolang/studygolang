"use client"

import { useEffect, useState, Suspense } from "react"
import { useRouter, useSearchParams, useParams } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { Card, CardContent } from "@/components/ui/card"
import { Loader2, CheckCircle2, XCircle } from "lucide-react"

function OAuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const params = useParams<{ provider: string }>()
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading")
  const [errorMsg, setErrorMsg] = useState("")

  useEffect(() => {
    // 后端 OAuth 回调会直接重定向回前端页面，此时 URL 中可能带有 error 参数
    const error = searchParams.get("error")
    const oauth = searchParams.get("oauth")

    if (error) {
      setStatus("error")
      setErrorMsg(error === "oauth_invalid" ? "授权验证失败，请重试"
        : error === "oauth_failed" ? "授权码获取失败"
        : error === "bind_failed" ? "绑定失败，请稍后重试"
        : error === "oauth_login_failed" ? "登录失败，请稍后重试"
        : "登录失败")
      return
    }

    if (oauth === "bind_success") {
      setStatus("success")
      setTimeout(() => router.push("/"), 2000)
      return
    }

    // 如果没有 error 和 oauth 参数，说明后端回调成功并重定向回来了
    // 此时 Cookie 已经设好，直接视为登录成功
    setStatus("success")
    setTimeout(() => router.push("/"), 1500)
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
              <p className="font-medium">登录成功</p>
              <p className="mt-2 text-xs text-muted-foreground">即将跳转...</p>
            </>
          )}
          {status === "error" && (
            <>
              <XCircle className="mb-4 h-12 w-12 text-destructive" />
              <p className="font-medium text-destructive">登录失败</p>
              <p className="mt-1 text-sm text-muted-foreground">{errorMsg}</p>
              <a href="/account/login" className="mt-4 text-sm text-primary hover:underline">
                返回登录
              </a>
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
