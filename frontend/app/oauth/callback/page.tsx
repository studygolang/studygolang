"use client"

import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { Card, CardContent } from "@/components/ui/card"
import { Loader2, CheckCircle2, XCircle } from "lucide-react"

function OAuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading")
  const [errorMsg, setErrorMsg] = useState("")

  useEffect(() => {
    // 后端 OAuth 回调重定向回来后，检查是否有错误
    const error = searchParams.get("error")
    const oauth = searchParams.get("oauth")

    if (error) {
      setStatus("error")
      setErrorMsg("登录失败，请重试")
      return
    }

    if (oauth === "bind_success") {
      setStatus("success")
      setTimeout(() => router.push("/"), 2000)
      return
    }

    // 回调成功
    setStatus("success")
    setTimeout(() => router.push("/"), 1500)
  }, [router, searchParams])

  return (
    <div className="flex flex-col items-center justify-center py-8">
      {status === "loading" && (
        <div className="flex flex-col items-center">
          <Loader2 className="mb-4 h-8 w-8 animate-spin text-muted-foreground" />
          <p className="text-sm text-muted-foreground">正在登录...</p>
        </div>
      )}
      {status === "success" && (
        <>
          <CheckCircle2 className="mb-4 h-12 w-12 text-green-500" />
          <p className="text-lg font-medium">登录成功</p>
          <p className="text-sm text-muted-foreground">正在跳转...</p>
        </>
      )}
      {status === "error" && (
        <>
          <XCircle className="mb-4 h-12 w-12 text-destructive" />
          <p className="text-lg font-medium text-destructive">登录失败</p>
          <p className="mt-2 text-sm text-muted-foreground">{errorMsg}</p>
          <a href="/account/login" className="mt-4 text-sm text-primary hover:underline">
            返回登录
          </a>
        </>
      )}
    </div>
  )
}

export default function OAuthCallbackPage() {
  return (
    <PageLayout sidebar={false}>
      <Card>
        <CardContent>
          <Suspense
            fallback={
              <div className="flex flex-col items-center justify-center py-8">
                <Loader2 className="mb-4 h-8 w-8 animate-spin text-muted-foreground" />
                <p className="text-sm text-muted-foreground">正在登录...</p>
              </div>
            }
          >
            <OAuthCallbackContent />
          </Suspense>
        </CardContent>
      </Card>
    </PageLayout>
  )
}
