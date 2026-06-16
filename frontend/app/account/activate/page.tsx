"use client"

import { useState, useEffect, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { accountAPI } from "@/lib/api"
import { toast } from "sonner"
import { CheckCircle2, XCircle, Loader2 } from "lucide-react"

function ActivateContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  // 邮件链接使用 `?param=xxx`（见 internal/logic/email.go:SendActivateMail）
  // 同时兼容历史 `?token=xxx` 调用
  const token = searchParams.get("param") || searchParams.get("token")

  const [loading, setLoading] = useState(false)
  const [activated, setActivated] = useState(false)
  const [error, setError] = useState("")
  const [sendingEmail, setSendingEmail] = useState(false)

  // 如果有 token，自动激活
  useEffect(() => {
    if (!token) return

    setLoading(true)
    accountAPI.activate(token)
      .then(() => {
        setActivated(true)
        toast.success("账户激活成功")
        setTimeout(() => router.push("/"), 2000)
      })
      .catch((err) => {
        setError(err.message || "激活失败")
        toast.error(err.message || "激活失败")
      })
      .finally(() => setLoading(false))
  }, [token, router])

  const handleSendEmail = async () => {
    setSendingEmail(true)
    try {
      await accountAPI.sendActivateEmail()
      toast.success("激活邮件已发送，请查收")
    } catch (err: any) {
      toast.error(err.message || "发送失败")
    } finally {
      setSendingEmail(false)
    }
  }

  // 激活中
  if (loading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <Loader2 className="h-12 w-12 animate-spin text-primary" />
            <p className="text-muted-foreground">正在激活账户...</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 激活成功
  if (activated) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <CheckCircle2 className="h-12 w-12 text-green-600" />
            <h2 className="text-xl font-semibold">账户激活成功</h2>
            <p className="text-muted-foreground text-center">
              你的账户已成功激活，即将跳转到首页...
            </p>
            <Button onClick={() => router.push("/")}>立即跳转</Button>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 有 token 但激活失败
  if (token && error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center gap-4 py-8">
            <XCircle className="h-12 w-12 text-destructive" />
            <h2 className="text-xl font-semibold">激活失败</h2>
            <p className="text-destructive text-center">{error}</p>
            <Button onClick={handleSendEmail} disabled={sendingEmail}>
              {sendingEmail ? "发送中..." : "重新发送激活邮件"}
            </Button>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 无 token，显示发送邮件表单
  return (
    <Card>
      <CardHeader>
        <CardTitle>激活账户</CardTitle>
        <CardDescription>
          点击下方按钮发送激活邮件到你的注册邮箱
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center gap-4 py-4">
          <p className="text-sm text-muted-foreground text-center">
            激活邮件将发送到你注册时使用的邮箱地址
          </p>
          <Button onClick={handleSendEmail} disabled={sendingEmail} size="lg">
            {sendingEmail ? "发送中..." : "发送激活邮件"}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

export default function ActivatePage() {
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
            <ActivateContent />
          </Suspense>
        </div>
      </main>
      <SiteFooter />
    </div>
  )
}
