'use client'

import { useEffect, useState, useCallback } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { PageLayout } from '@/components/page-layout'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { userAPI } from '@/lib/api'
import type { User } from '@/lib/types'
import { User as UserIcon, Save, Loader2 } from 'lucide-react'

export default function UserSettingsPage() {
  const params = useParams<{ username: string }>()
  const router = useRouter()
  const username = params.username

  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  // 表单字段
  const [name, setName] = useState('')
  const [city, setCity] = useState('')
  const [company, setCompany] = useState('')
  const [github, setGithub] = useState('')
  const [website, setWebsite] = useState('')
  const [introduce, setIntroduce] = useState('')
  const [open, setOpen] = useState('0')

  // 加载用户资料
  const loadProfile = useCallback(async () => {
    try {
      const data = await userAPI.getMyProfile()
      if (data?.user) {
        const u = data.user
        setUser(u)
        setName(u.name || '')
        setCity(u.city || '')
        setCompany(u.company || '')
        setGithub(u.github || '')
        setWebsite(u.website || '')
        setIntroduce(u.introduce || '')
        setOpen(u.open === 1 ? '1' : '0')

        // 验证当前登录用户是否与 URL 中的 username 匹配
        if (u.username !== username) {
          setMessage({ type: 'error', text: '只能编辑自己的资料' })
        }
      }
    } catch {
      setMessage({ type: 'error', text: '请先登录后再编辑资料' })
    } finally {
      setLoading(false)
    }
  }, [username])

  useEffect(() => {
    loadProfile()
  }, [loadProfile])

  // 保存资料
  const handleSave = async () => {
    setMessage(null)

    // 客户端验证
    if (name.length > 50) {
      setMessage({ type: 'error', text: '昵称不能超过50个字符' })
      return
    }
    if (name.length === 0) {
      setMessage({ type: 'error', text: '昵称不能为空' })
      return
    }
    if (website && !website.match(/^https?:\/\/.+/)) {
      setMessage({ type: 'error', text: '网站地址格式不正确，需以 http:// 或 https:// 开头' })
      return
    }
    if (github && !github.match(/^[a-zA-Z0-9-]+$/)) {
      setMessage({ type: 'error', text: 'GitHub 用户名格式不正确，只允许字母、数字和连字符' })
      return
    }
    if (introduce.length > 500) {
      setMessage({ type: 'error', text: '个人简介不能超过500个字符' })
      return
    }

    setSaving(true)
    try {
      await userAPI.updateProfile({
        name,
        city,
        company,
        github,
        website,
        introduce,
        open,
      })
      setMessage({ type: 'success', text: '资料更新成功' })
    } catch (err) {
      setMessage({
        type: 'error',
        text: err instanceof Error ? err.message : '更新失败，请重试',
      })
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <PageLayout sidebar={false}>
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      </PageLayout>
    )
  }

  // 未登录或非本人
  if (!user || user.username !== username) {
    return (
      <PageLayout sidebar={false}>
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <UserIcon className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm font-medium text-foreground">无权访问</p>
          <p className="mt-1 text-sm text-muted-foreground">只能编辑自己的资料</p>
        </div>
      </PageLayout>
    )
  }

  return (
    <PageLayout sidebar={false}>
      <div className="mx-auto max-w-2xl">
        <h1 className="mb-6 text-xl font-bold text-foreground">个人设置</h1>

        {message && (
          <div
            className={`mb-4 rounded-md px-4 py-3 text-sm ${
              message.type === 'success'
                ? 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400'
                : 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400'
            }`}
          >
            {message.text}
          </div>
        )}

        {/* 基本信息 */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-base">基本信息</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">昵称</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="输入昵称"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="introduce">个人签名</Label>
              <Textarea
                id="introduce"
                value={introduce}
                onChange={(e) => setIntroduce(e.target.value)}
                placeholder="介绍一下自己吧"
                rows={3}
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="city">城市</Label>
                <Input
                  id="city"
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                  placeholder="所在城市"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="company">公司</Label>
                <Input
                  id="company"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                  placeholder="所在公司"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="website">个人网站</Label>
              <Input
                id="website"
                value={website}
                onChange={(e) => setWebsite(e.target.value)}
                placeholder="https://example.com"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="github">GitHub</Label>
              <Input
                id="github"
                value={github}
                onChange={(e) => setGithub(e.target.value)}
                placeholder="GitHub 用户名"
              />
            </div>

            <div className="space-y-2">
              <Label>信息公开</Label>
              <div className="flex items-center gap-4">
                <label className="flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    name="open"
                    value="1"
                    checked={open === '1'}
                    onChange={(e) => setOpen(e.target.value)}
                    className="accent-primary"
                  />
                  公开
                </label>
                <label className="flex items-center gap-2 text-sm">
                  <input
                    type="radio"
                    name="open"
                    value="0"
                    checked={open === '0'}
                    onChange={(e) => setOpen(e.target.value)}
                    className="accent-primary"
                  />
                  不公开
                </label>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* 保存按钮 */}
        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={saving} className="gap-2">
            {saving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            保存修改
          </Button>
        </div>
      </div>
    </PageLayout>
  )
}
