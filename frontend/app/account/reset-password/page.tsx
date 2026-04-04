"use client"

import { useState, Suspense } from "react"
import Link from "next/link"
import { useRouter, useSearchParams } from "next/navigation"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Lock, ArrowLeft, Loader2, CheckCircle } from "lucide-react"

function ResetPasswordForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get("token") || ""

  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!token) {
    return (
      <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-center text-sm text-destructive">
        无效的重置链接，请重新申请{" "}
        <Link href="/account/forgot-password" className="underline">
          找回密码
        </Link>
      </div>
    )
  }

  async function handleSubmit(e: React.SyntheticEvent<HTMLFormElement>) {
    e.preventDefault()
    if (!password || password !== confirmPassword) return
    if (password.length < 6) {
      setError("密码长度不能少于6位")
      return
    }

    setLoading(true)
    setError(null)

    try {
      const res = await fetch("/api/v1/user/reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token, new_password: password }),
      })
      const json = await res.json()
      if (json.code === 0) {
        setSuccess(true)
        setTimeout(() => router.push("/account/login"), 2000)
      } else {
        setError(json.message || "重置失败，链接可能已过期")
      }
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="flex flex-col items-center gap-4 rounded-lg bg-green-50 p-6 text-center">
        <CheckCircle className="h-12 w-12 text-green-500" />
        <div className="space-y-1">
          <p className="font-medium text-foreground">密码重置成功</p>
          <p className="text-sm text-muted-foreground">
            正在跳转到登录页面...
          </p>
        </div>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="password">新密码</Label>
        <div className="relative">
          <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            id="password"
            type="password"
            placeholder="请输入新密码（至少6位）"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="pl-9"
            required
            minLength={6}
            disabled={loading}
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="confirmPassword">确认密码</Label>
        <div className="relative">
          <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            id="confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="pl-9"
            required
            disabled={loading}
          />
        </div>
        {confirmPassword && password !== confirmPassword && (
          <p className="text-xs text-destructive">两次输入的密码不一致</p>
        )}
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <Button
        type="submit"
        className="w-full"
        disabled={loading || !password || password !== confirmPassword}
      >
        {loading ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            重置中...
          </>
        ) : (
          "重置密码"
        )}
      </Button>
    </form>
  )
}

export default function ResetPasswordPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <Link href="/" className="text-2xl font-bold text-primary">
            Go语言中文网
          </Link>
          <p className="mt-2 text-sm text-muted-foreground">重置你的密码</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-xl">重置密码</CardTitle>
            <CardDescription>请输入你的新密码</CardDescription>
          </CardHeader>

          <CardContent>
            <Suspense fallback={<div className="h-40 animate-pulse rounded-lg bg-muted" />}>
              <ResetPasswordForm />
            </Suspense>
          </CardContent>

          <CardFooter className="flex justify-center">
            <Link
              href="/account/login"
              className="flex items-center gap-1 text-sm text-muted-foreground hover:text-primary"
            >
              <ArrowLeft className="h-3 w-3" />
              返回登录
            </Link>
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}
