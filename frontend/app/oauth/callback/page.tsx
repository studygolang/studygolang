"use client"

import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { CheckCircle2, Loader2 } from "lucide-react"

// 合法的 OAuth code/state 字符白名单
const isValidOAuthParam = (value: string): boolean => /^[a-zA-Z0-9_-]+$/.test(value)

function OAuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading")

  useEffect(() => {
    const code = searchParams.get("code")
    const state = searchParams.get("state")

    if (!code || !state || !isValidOAuthParam(code) || !isValidOAuthParam(state)) {
      setStatus("error")
      return
    }

    // OAuth 回调由后端处理，这里只做跳转
    if (state === "bind") {
      router.push("/account/edit#connection")
      return
    }

    // 登录成功，跳转首页
    setStatus("success")
    const timer = setTimeout(() => {
      router.push("/")
    }, 2000)

    return () => clearTimeout(timer)
  }, [router, searchParams])

  return (
    <>
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
            <p className="text-lg font-medium text-destructive">登录失败</p>
            <p className="mt-2 text-sm text-muted-foreground">请重新尝试</p>
          </>
        )}
      </div>
    </>
  )
}

export default function OAuthCallbackPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader title="OAuth 登录" breadcrumbs={[{ label: "账号" }]} />
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
