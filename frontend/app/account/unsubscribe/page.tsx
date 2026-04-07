"use client"

import { useState, useEffect, Suspense } from "react"
import { useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { accountAPI } from "@/lib/api"
import { toast } from "sonner"
import { Loader2, CheckCircle2, Mail } from "lucide-react"

function UnsubscribeContent() {
  const searchParams = useSearchParams()
  const token = searchParams.get("token")

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [unsubscribed, setUnsubscribed] = useState(false)
  const [email, setEmail] = useState("")
  const [error, setError] = useState("")

  // 加载退订信息
  useEffect(() => {
    if (!token) {
      setError("缺少 token 参数")
      setLoading(false)
      return
    }

    accountAPI.unsubscribePage(token)
      .then((data) => {
        setEmail(data.email)
      })
      .catch((err) => {
        setError(err.message || "加载失败")
      })
      .finally(() => setLoading(false))
  }, [token])

  const handleUnsubscribe = async () => {
    if (!token) return

    setSubmitting(true)
    try {
      await accountAPI.unsubscribe(token)
      setUnsubscribed(true)
      toast.success("已成功退订邮件")
    } catch (err: any) {
      toast.error(err.message || "退订失败")
    } finally {
      setSubmitting(false)
    }
  }

  // 加载中
  if (loading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <Loader2 className="h-12 w-12 animate-spin text-primary" />
            <p className="text-muted-foreground">加载中...</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 错误
  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <p className="text-destructive">{error}</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 退订成功
  if (unsubscribed) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <CheckCircle2 className="h-12 w-12 text-green-600" />
            <h2 className="text-xl font-semibold">退订成功</h2>
            <p className="text-muted-foreground text-center">
              你将不再收到来自 Go语言中文网 的邮件通知
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 退订确认
  return (
    <Card>
      <CardHeader>
        <CardTitle>邮件退订</CardTitle>
        <CardDescription>
          确认退订后，你将不再收到来自 Go语言中文网 的邮件通知
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          <div className="flex items-center gap-3 rounded-lg bg-muted p-4">
            <Mail className="h-5 w-5 text-muted-foreground" />
            <div className="flex-1">
              <p className="text-sm font-medium">邮箱地址</p>
              <p className="text-sm text-muted-foreground">{email}</p>
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <Button
              onClick={handleUnsubscribe}
              disabled={submitting}
              variant="destructive"
              size="lg"
            >
              {submitting ? "处理中..." : "确认退订"}
            </Button>
            <p className="text-xs text-muted-foreground text-center">
              退订后，你仍可以在个人设置中重新订阅
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default function UnsubscribePage() {
  return (
    <div className="min-h-screen flex flex-col bg-muted/30">
      <SiteHeader />
      <main className="flex-1 flex items-center justify-center py-12 px-4">
        <div className="w-full max-w-md">
          <Suspense fallback={
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
              </CardContent>
            </Card>
          }>
            <UnsubscribeContent />
          </Suspense>
        </div>
      </main>
      <SiteFooter />
    </div>
  )
}
