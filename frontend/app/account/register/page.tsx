"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

export default function RegisterPage() {
  const router = useRouter()

  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  })
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)
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
            router.replace("/")
          }
        }
      } catch {
        // 未登录,继续显示注册页面
      }
    }
    checkAuth()
  }, [router])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")

    if (formData.username.length < 4 || formData.username.length > 20) {
      setError("用户名长度必须在 4-20 个字符之间")
      return
    }
    if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
      setError("用户名只能包含大小写字母、数字和下划线")
      return
    }
    if (formData.password.length < 6 || formData.password.length > 32) {
      setError("密码长度必须在 6-32 个字符之间")
      return
    }
    if (formData.password !== formData.confirmPassword) {
      setError("两次输入的密码不一致")
      return
    }

    setLoading(true)

    try {
      const res = await fetch(`${API_BASE}/api/v1/user/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: formData.username,
          email: formData.email,
          passwd: formData.password,
        }),
        credentials: "include",
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 0) {
          setSuccess(true)
        } else {
          setError(data.msg || "注册失败")
        }
      } else {
        setError("注册失败,请稍后重试")
      }
    } catch {
      setError("网络错误,请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-muted/30 p-4">
        <Link href="/" className="mb-6 text-2xl font-bold text-primary hover:opacity-80 transition-opacity">
          Go语言中文网
        </Link>
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>注册成功</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              注册成功！请查收激活邮件。
            </p>
            <p className="text-sm text-muted-foreground">
              如果没有收到激活邮件,可以关注站长公众号,回复{" "}
              <span className="text-destructive font-semibold">{formData.username}</span>{" "}
              获取验证码来激活。
            </p>
            <div className="flex justify-center">
              <Button onClick={() => router.push("/account/login")}>
                前往登录
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-muted/30 p-4">
      <Link href="/" className="mb-6 text-2xl font-bold text-primary hover:opacity-80 transition-opacity">
        Go语言中文网
      </Link>
      <div className="w-full max-w-4xl grid md:grid-cols-3 gap-6">
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>注册新用户</CardTitle>
            <CardDescription>加入 Go语言中文网社区</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <Label htmlFor="username">
                  用户名 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="username"
                  name="username"
                  type="text"
                  placeholder="请输入用户名"
                  value={formData.username}
                  onChange={handleChange}
                  required
                />
                <p className="text-xs text-muted-foreground">
                  只能包含大小写字母、数字和下划线(4-20字符)
                </p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">
                  Email <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  placeholder="请输入Email"
                  value={formData.email}
                  onChange={handleChange}
                  required
                />
                <p className="text-xs text-muted-foreground">
                  可以在个人资料设置中更改
                </p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="password">
                  密码 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  placeholder="请输入密码"
                  value={formData.password}
                  onChange={handleChange}
                  required
                />
                <p className="text-xs text-muted-foreground">6-32个字符</p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirmPassword">
                  确认密码 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  placeholder="请再次输入密码"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  required
                />
              </div>
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "注册中..." : "注册"}
              </Button>
            </form>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">已有帐号？</CardTitle>
            </CardHeader>
            <CardContent>
              <Link href="/account/login" className="block text-sm text-primary hover:underline">
                立即登录
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
