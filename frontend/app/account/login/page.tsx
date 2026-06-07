"use client"

import { useState, useEffect, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Github } from "lucide-react"
import { useAuth } from "@/lib/auth-context"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

function LoginForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const rawRedirect = searchParams.get("redirect") || "/"
  const redirect = (rawRedirect.startsWith("/") && !rawRedirect.startsWith("//")) ? rawRedirect : "/"
  const { login } = useAuth()

  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [rememberMe, setRememberMe] = useState(true)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  // 检查是否已登录
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(`${API_BASE}/api/v1/user/me`, { credentials: "include" })
        if (res.ok) {
          const data = await res.json()
          if (data.code === 0 && data.data?.user) {
            // 已登录,跳转到首页
            router.replace(redirect)
          }
        }
      } catch {
        // 未登录,继续显示登录页面
      }
    }
    checkAuth()
  }, [router, redirect])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    try {
      const res = await fetch(`${API_BASE}/api/v1/user/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          passwd: password,
        }),
        credentials: "include",
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 0) {
          if (data.data) {
            login({
              uid: data.data.uid,
              username: data.data.username,
              name: data.data.username,
              avatar: "",
              is_root: false,
              is_vip: false,
            })
          }
          window.location.href = redirect
        } else {
          setError(data.msg || "登录失败")
        }
      } else {
        setError("登录失败,请检查用户名和密码")
      }
    } catch {
      setError("网络错误,请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-muted/30 p-4">
      <Link href="/" className="mb-6 text-2xl font-bold text-primary hover:opacity-80 transition-opacity">
        Go语言中文网
      </Link>
      <div className="w-full max-w-4xl grid md:grid-cols-3 gap-6">
        <Card className="md:col-span-2">

          <CardHeader>
            <CardTitle>登录</CardTitle>
            <CardDescription>登录 Go语言中文网</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <Label htmlFor="username">用户名或邮箱</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="password">密码</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="请输入密码"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="remember"
                  checked={rememberMe}
                  onCheckedChange={(checked) => setRememberMe(checked as boolean)}
                />
                <Label htmlFor="remember" className="text-sm font-normal cursor-pointer">
                  记住登录状态
                </Label>
              </div>
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "登录中..." : "登录"}
              </Button>
            </form>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">第三方账号登录</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button variant="outline" className="w-full" asChild>
                <a href="/oauth/github/login">
                  <Github className="mr-2 h-4 w-4" />
                  GitHub
                </a>
              </Button>
              <Button variant="outline" className="w-full" asChild>
                <a href="/oauth/gitea/login">
                  <Github className="mr-2 h-4 w-4" />
                  Gitea
                </a>
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base">还没有帐号？</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <Link href="/account/register" className="block text-sm text-primary hover:underline">
                注册
              </Link>
              <Link href="/account/forgot-password" className="block text-sm text-primary hover:underline">
                忘记了密码？
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">加载中...</div>}>
      <LoginForm />
    </Suspense>
  )
}
