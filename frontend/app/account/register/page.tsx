"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"

export default function RegisterPage() {
  const router = useRouter()
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  })
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const [loading, setLoading] = useState(false)

  function handleChange(field: keyof typeof form) {
    return (e: React.ChangeEvent<HTMLInputElement>) => {
      setForm((prev) => ({ ...prev, [field]: e.target.value }))
    }
  }

  function validate(): string | null {
    if (!form.username.trim()) return "请输入用户名"
    if (form.username.trim().length < 3) return "用户名至少 3 个字符"
    if (!form.email.trim()) return "请输入邮箱"
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) return "邮箱格式不正确"
    if (!form.password) return "请输入密码"
    if (form.password.length < 6) return "密码至少 6 位"
    if (form.password !== form.confirmPassword) return "两次输入的密码不一致"
    return null
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    const validationError = validate()
    if (validationError) {
      setError(validationError)
      return
    }

    setLoading(true)
    try {
      // 客户端组件使用相对路径，通过 next.config.mjs 中的 rewrites 代理到后端
      const res = await fetch(`/api/v1/user/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: form.username.trim(),
          email: form.email.trim(),
          passwd: form.password,
        }),
      })
      const json = await res.json()
      if (json.code !== 0) {
        setError(json.msg || "注册失败，请稍后重试")
        return
      }
      setSuccess(true)
      setTimeout(() => {
        router.push("/account/login")
      }, 2000)
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
        <div className="w-full max-w-md">
          <Card>
            <CardContent className="flex flex-col items-center gap-4 py-12">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
                <svg
                  className="h-8 w-8 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h2 className="text-xl font-bold text-foreground">注册成功！</h2>
              <p className="text-center text-sm text-muted-foreground">
                账号已创建，正在跳转到登录页...
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
      <div className="w-full max-w-md">
        {/* Logo / Title */}
        <div className="mb-8 text-center">
          <Link href="/" className="text-2xl font-bold text-primary">
            Go语言中文网
          </Link>
          <p className="mt-2 text-sm text-muted-foreground">加入 Go 开发者社区</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-xl">创建账号</CardTitle>
            <CardDescription>填写以下信息完成注册</CardDescription>
          </CardHeader>

          <form onSubmit={handleSubmit} noValidate>
            <CardContent className="space-y-4">
              {error && (
                <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {error}
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="username">用户名</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="至少 3 个字符"
                  value={form.username}
                  onChange={handleChange("username")}
                  autoComplete="username"
                  disabled={loading}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">邮箱</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="your@email.com"
                  value={form.email}
                  onChange={handleChange("email")}
                  autoComplete="email"
                  disabled={loading}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">密码</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="至少 6 位"
                  value={form.password}
                  onChange={handleChange("password")}
                  autoComplete="new-password"
                  disabled={loading}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirmPassword">确认密码</Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="再次输入密码"
                  value={form.confirmPassword}
                  onChange={handleChange("confirmPassword")}
                  autoComplete="new-password"
                  disabled={loading}
                />
              </div>
            </CardContent>

            <CardFooter className="flex flex-col gap-3">
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "注册中..." : "注册"}
              </Button>
              <p className="text-center text-sm text-muted-foreground">
                已有账号？{" "}
                <Link
                  href="/account/login"
                  className="font-medium text-primary transition-colors hover:text-primary/80"
                >
                  立即登录
                </Link>
              </p>
            </CardFooter>
          </form>
        </Card>
      </div>
    </div>
  )
}
