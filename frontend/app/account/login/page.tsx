"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import type { LoginData } from "@/lib/types"

export default function LoginPage() {
  const router = useRouter()
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [checking, setChecking] = useState(true)

  // 检查是否已登录，如果已登录则跳转到首页
  useEffect(() => {
    const uid = localStorage.getItem("uid")
    if (uid) {
      // 验证 token 是否有效
      fetch("/api/v1/user/me", { credentials: "include" })
        .then((res) => res.json())
        .then((data) => {
          if (data.code === 0) {
            // 已登录，跳转到首页
            router.replace("/")
          } else {
            // token 无效，清除本地存储
            localStorage.removeItem("uid")
            localStorage.removeItem("username")
            setChecking(false)
          }
        })
        .catch(() => {
          setChecking(false)
        })
    } else {
      setChecking(false)
    }
  }, [router])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)

    if (!username.trim()) {
      setError("请输入用户名或邮箱")
      return
    }
    if (!password) {
      setError("请输入密码")
      return
    }

    setLoading(true)
    try {
      // 客户端组件使用相对路径，通过 next.config.mjs 中的 rewrites 代理到后端
      const res = await fetch(`/api/v1/user/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: username.trim(), password }),
      })
      const json = await res.json()
      if (json.code !== 0) {
        setError(json.msg || "登录失败，请检查用户名和密码")
        return
      }
      const data: LoginData = json.data
      // token 已通过 HttpOnly Cookie 存储，此处仅保存非敏感 UI 状态
      localStorage.setItem("uid", String(data.uid))
      localStorage.setItem("username", data.username)
      router.push("/")
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  // 正在检查登录状态，显示加载中
  if (checking) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center">
          <div className="mb-4 text-sm text-muted-foreground">检查登录状态...</div>
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
          <p className="mt-2 text-sm text-muted-foreground">欢迎回来，请登录你的账号</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-xl">登录</CardTitle>
            <CardDescription>使用用户名/邮箱和密码登录</CardDescription>
          </CardHeader>

          <form onSubmit={handleSubmit}>
            <CardContent className="space-y-3">
              {error && (
                <div className="rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {error}
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="username">用户名 / 邮箱</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  autoComplete="username"
                  disabled={loading}
                />
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="password">密码</Label>
                  <Link
                    href="/account/forgot-password"
                    className="text-xs text-muted-foreground transition-colors hover:text-primary"
                  >
                    忘记密码？
                  </Link>
                </div>
                <Input
                  id="password"
                  type="password"
                  placeholder="请输入密码"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  disabled={loading}
                />
              </div>

              <div className="pt-2">
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? "登录中..." : "登录"}
                </Button>
              </div>

              <p className="text-center text-sm text-muted-foreground">
                还没有账号？{" "}
                <Link
                  href="/account/register"
                  className="font-medium text-primary transition-colors hover:text-primary/80"
                >
                  立即注册
                </Link>
              </p>
            </CardContent>
          </form>
        </Card>
      </div>
    </div>
  )
}
