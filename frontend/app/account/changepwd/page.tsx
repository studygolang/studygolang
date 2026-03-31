"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { userAPI } from "@/lib/api"

export default function ChangePasswordPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [hasPasswd, setHasPasswd] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)

  const [formData, setFormData] = useState({
    curPasswd: "",
    newPasswd: "",
    confirmPasswd: "",
  })

  // 检查登录状态并加载用户信息
  useEffect(() => {
    const loadProfile = async () => {
      try {
        const data = await userAPI.getMyProfile()
        if (!data) {
          router.replace("/account/login")
          return
        }
        setHasPasswd(data.has_passwd)
      } catch {
        router.replace("/account/login")
      } finally {
        setLoading(false)
      }
    }
    loadProfile()
  }, [router])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSuccess(false)

    // 验证表单
    if (!hasPasswd && !formData.curPasswd) {
      setError("请输入当前密码")
      return
    }
    if (!formData.newPasswd) {
      setError("请输入新密码")
      return
    }
    if (formData.newPasswd.length < 6 || formData.newPasswd.length > 32) {
      setError("新密码长度必须在 6-32 个字符之间")
      return
    }
    if (formData.newPasswd !== formData.confirmPasswd) {
      setError("两次输入的新密码不一致")
      return
    }

    setSubmitting(true)

    try {
      await userAPI.changePassword(formData.curPasswd, formData.newPasswd)
      setSuccess(true)
      setFormData({ curPasswd: "", newPasswd: "", confirmPasswd: "" })
    } catch (err: any) {
      setError(err.message || "修改密码失败")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-muted/30">
        <div className="text-muted-foreground">加载中...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-muted/30 py-8">
      <div className="container mx-auto max-w-4xl px-4">
        <div className="mb-6">
          <Link href="/" className="text-2xl font-bold text-primary hover:opacity-80 transition-opacity">
            Go语言中文网
          </Link>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          {/* 左侧菜单 */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">个人设置</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <Link
                href="/account/edit"
                className="block text-sm text-muted-foreground hover:text-primary hover:underline"
              >
                基本信息
              </Link>
              <Link
                href="/account/changepwd"
                className="block text-sm text-primary font-medium hover:underline"
              >
                修改密码
              </Link>
            </CardContent>
          </Card>

          {/* 右侧内容 */}
          <div className="md:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>修改密码</CardTitle>
                <CardDescription>
                  {hasPasswd ? "修改你的登录密码" : "设置你的登录密码"}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  {error && (
                    <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
                      {error}
                    </div>
                  )}
                  {success && (
                    <div className="bg-green-500/10 text-green-600 px-4 py-3 rounded-md text-sm">
                      密码修改成功！
                    </div>
                  )}
                  {hasPasswd && (
                    <div className="space-y-2">
                      <Label htmlFor="curPasswd">当前密码</Label>
                      <Input
                        id="curPasswd"
                        name="curPasswd"
                        type="password"
                        placeholder="请输入当前密码"
                        value={formData.curPasswd}
                        onChange={handleChange}
                      />
                    </div>
                  )}
                  <div className="space-y-2">
                    <Label htmlFor="newPasswd">新密码</Label>
                    <Input
                      id="newPasswd"
                      name="newPasswd"
                      type="password"
                      placeholder="请输入新密码"
                      value={formData.newPasswd}
                      onChange={handleChange}
                    />
                    <p className="text-xs text-muted-foreground">6-32个字符</p>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="confirmPasswd">确认新密码</Label>
                    <Input
                      id="confirmPasswd"
                      name="confirmPasswd"
                      type="password"
                      placeholder="请再次输入新密码"
                      value={formData.confirmPasswd}
                      onChange={handleChange}
                    />
                  </div>
                  <Button type="submit" disabled={submitting}>
                    {submitting ? "提交中..." : "提交"}
                  </Button>
                </form>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
