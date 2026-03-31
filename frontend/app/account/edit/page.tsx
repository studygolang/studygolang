"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Textarea } from "@/components/ui/textarea"
import { userAPI } from "@/lib/api"
import { z } from "zod"

const profileSchema = z.object({
  name: z.string().min(1, "昵称不能为空").max(50, "昵称最多50个字符"),
  email: z.string().email("邮箱格式不正确").optional().or(z.literal("")),
  website: z.string().url("网站格式不正确").optional().or(z.literal("")),
  introduce: z.string().max(500, "简介最多500个字符").optional().or(z.literal("")),
})


export default function EditProfilePage() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  const [formData, setFormData] = useState({
    name: "",
    email: "",
    city: "",
    company: "",
    github: "",
    website: "",
    introduce: "",
    open: false,
  })
  const [avatarUrl, setAvatarUrl] = useState("")
  const [avatarInput, setAvatarInput] = useState("")

  // 检查登录状态并加载用户信息
  useEffect(() => {
    const loadProfile = async () => {
      try {
        const data = await userAPI.getMyProfile()
        if (!data) {
          router.replace("/account/login")
          return
        }
        setFormData({
          name: data.user.name || "",
          email: data.user.email || "",
          city: data.user.city || "",
          company: data.user.company || "",
          github: data.user.github || "",
          website: data.user.website || "",
          introduce: data.user.introduce || "",
          open: data.user.open === 1,
        })
        setAvatarUrl(data.user.avatar || "")
      } catch {
        router.replace("/account/login")
      } finally {
        setLoading(false)
      }
    }
    loadProfile()
  }, [router])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleCheckboxChange = (checked: boolean) => {
    setFormData((prev) => ({ ...prev, open: checked }))
  }

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setAvatarInput(e.target.value)
  }

  const handleAvatarSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSubmitting(true)

    try {
      await userAPI.updateAvatar(avatarInput)
      setAvatarUrl(avatarInput)
      setAvatarInput("")
      setSuccess(true)
      setTimeout(() => setSuccess(false), 3000)
    } catch (err: any) {
      setError(err.message || "更新头像失败")
    } finally {
      setSubmitting(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setValidationErrors({})
    setSubmitting(true)

    // Validate form data with Zod
    const result = profileSchema.safeParse({
      name: formData.name,
      email: formData.email,
      website: formData.website,
      introduce: formData.introduce,
    })

    if (!result.success) {
      const errors: Record<string, string> = {}
      result.error.issues.forEach((issue) => {
        const path = issue.path[0] as string
        if (!errors[path]) {
          errors[path] = issue.message
        }
      })
      setValidationErrors(errors)
      setSubmitting(false)
      return
    }

    try {
      await userAPI.updateProfile({
        name: formData.name,
        email: formData.email,
        city: formData.city,
        company: formData.company,
        github: formData.github,
        website: formData.website,
        introduce: formData.introduce,
        open: formData.open ? "1" : "0",
      })
      setSuccess(true)
      setTimeout(() => setSuccess(false), 3000)
    } catch (err: any) {
      setError(err.message || "更新失败")
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
                className="block text-sm text-primary font-medium hover:underline"
              >
                基本信息
              </Link>
              <Link
                href="/account/changepwd"
                className="block text-sm text-muted-foreground hover:text-primary hover:underline"
              >
                修改密码
              </Link>
            </CardContent>
          </Card>

          {/* 右侧内容 */}
          <div className="md:col-span-2 space-y-6">
            {/* 头像设置 */}
            <Card>
              <CardHeader>
                <CardTitle>头像设置</CardTitle>
                <CardDescription>设置你的头像图片地址</CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleAvatarSubmit} className="space-y-4">
                  {error && (
                    <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
                      {error}
                    </div>
                  )}
                  {success && (
                    <div className="bg-green-500/10 text-green-600 px-4 py-3 rounded-md text-sm">
                      头像更新成功！
                    </div>
                  )}
                  {avatarUrl && (
                    <div className="flex justify-center mb-4">
                      <img
                        src={avatarUrl}
                        alt="当前头像"
                        className="w-24 h-24 rounded-full object-cover border-2 border-muted"
                        onError={(e) => {
                          (e.target as HTMLImageElement).src = "/default-avatar.png"
                        }}
                      />
                    </div>
                  )}
                  <div className="space-y-2">
                    <Label htmlFor="avatar">头像 URL</Label>
                    <Input
                      id="avatar"
                      type="url"
                      placeholder="请输入头像图片地址"
                      value={avatarInput}
                      onChange={handleAvatarChange}
                    />
                    <p className="text-xs text-muted-foreground">
                      支持 http:// 或 https:// 开头的图片地址
                    </p>
                  </div>
                  <Button type="submit" disabled={submitting || !avatarInput}>
                    {submitting ? "更新中..." : "更新头像"}
                  </Button>
                </form>
              </CardContent>
            </Card>

            {/* 基本信息 */}
            <Card>
              <CardHeader>
                <CardTitle>基本信息</CardTitle>
                <CardDescription>修改你的个人资料信息</CardDescription>
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
                      保存成功！
                    </div>
                  )}
                  <div className="space-y-2">
                    <Label htmlFor="name">昵称</Label>
                    <Input
                      id="name"
                      name="name"
                      type="text"
                      placeholder="请输入昵称"
                      value={formData.name}
                      onChange={handleChange}
                      maxLength={50}
                    />
                    {validationErrors.name && (
                      <p className="text-xs text-destructive">{validationErrors.name}</p>
                    )}
                    <p className="text-xs text-muted-foreground">最多50个字符</p>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email">邮箱</Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      placeholder="请输入邮箱"
                      value={formData.email}
                      onChange={handleChange}
                    />
                    {validationErrors.email && (
                      <p className="text-xs text-destructive">{validationErrors.email}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="city">城市</Label>
                    <Input
                      id="city"
                      name="city"
                      type="text"
                      placeholder="请输入所在城市"
                      value={formData.city}
                      onChange={handleChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="company">公司</Label>
                    <Input
                      id="company"
                      name="company"
                      type="text"
                      placeholder="请输入公司名称"
                      value={formData.company}
                      onChange={handleChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="github">GitHub</Label>
                    <Input
                      id="github"
                      name="github"
                      type="text"
                      placeholder="请输入 GitHub 用户名"
                      value={formData.github}
                      onChange={handleChange}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="website">个人网站</Label>
                    <Input
                      id="website"
                      name="website"
                      type="url"
                      placeholder="请输入个人网站地址"
                      value={formData.website}
                      onChange={handleChange}
                      maxLength={200}
                    />
                    {validationErrors.website && (
                      <p className="text-xs text-destructive">{validationErrors.website}</p>
                    )}
                    <p className="text-xs text-muted-foreground">
                      必须以 http:// 或 https:// 开头
                    </p>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="introduce">个人简介</Label>
                    <Textarea
                      id="introduce"
                      name="introduce"
                      placeholder="请输入个人简介"
                      value={formData.introduce}
                      onChange={handleChange}
                      rows={4}
                      maxLength={500}
                    />
                    <p className="text-xs text-muted-foreground">
                      最多500个字符
                    </p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="open"
                      checked={formData.open}
                      onCheckedChange={handleCheckboxChange}
                    />
                    <Label htmlFor="open" className="text-sm font-normal cursor-pointer">
                      公开我的个人资料
                    </Label>
                  </div>
                  <Button type="submit" disabled={submitting}>
                    {submitting ? "保存中..." : "保存"}
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
