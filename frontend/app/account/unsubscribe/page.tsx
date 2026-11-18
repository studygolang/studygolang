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
  // 邮件链接格式：?u={token}&email={email}（兼容旧参数名 token）
  const token = searchParams.get("u") ?? searchParams.get("token")
  const emailParam = searchParams.get("email")

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [unsubscribed, setUnsubscribed] = useState(false)
  const [completed, setCompleted] = useState(false)
  const [email, setEmail] = useState("")
  const [error, setError] = useState("")

  // 加载退订信息
  useEffect(() => {
    if (!token || !emailParam) {
      setError("缺少 token 或 email 参数")
      setLoading(false)
      return
    }

    accountAPI.unsubscribePage(token, emailParam)
      .then((data) => {
        setEmail(data.email || emailParam)
        setUnsubscribed(data.unsubscribe === 1)
      })
      .catch((err) => {
        setError(err.message || "加载失败")
      })
      .finally(() => setLoading(false))
  }, [token, emailParam])

  const handleSubscriptionChange = async () => {
    if (!token || !emailParam) return

    const nextUnsubscribed = !unsubscribed
    setSubmitting(true)
    try {
      await accountAPI.updateSubscription(token, emailParam, nextUnsubscribed ? 1 : 0)
      setUnsubscribed(nextUnsubscribed)
      setCompleted(true)
      toast.success(nextUnsubscribed ? "已成功退订邮件" : "已重新订阅每周精选")
    } catch (err: any) {
      toast.error(err.message || "操作失败")
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

  if (completed) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <CheckCircle2 className="h-12 w-12 text-green-600" />
            <h2 className="text-xl font-semibold">{unsubscribed ? "退订成功" : "订阅成功"}</h2>
            <p className="text-muted-foreground text-center">
              {unsubscribed
                ? "你将不再收到来自 Go语言中文网 的每周精选邮件"
                : "你将恢复接收来自 Go语言中文网 的每周精选邮件"}
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>邮件订阅设置</CardTitle>
        <CardDescription>
          {unsubscribed
            ? "你当前已退订每周精选，可以重新订阅"
            : "确认退订后，你将不再收到来自 Go语言中文网 的每周精选邮件"}
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
              onClick={handleSubscriptionChange}
              disabled={submitting}
              variant={unsubscribed ? "default" : "destructive"}
              size="lg"
            >
              {submitting ? "处理中..." : unsubscribed ? "重新订阅" : "确认退订"}
            </Button>
            <p className="text-xs text-muted-foreground text-center">
              {unsubscribed
                ? "重新订阅后，将恢复接收每周精选邮件"
                : "退订后，仍可通过邮件中的此链接重新订阅"}
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
