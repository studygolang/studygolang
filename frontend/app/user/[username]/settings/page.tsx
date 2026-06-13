'use client'

import { useEffect, useState, useCallback } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { PageLayout } from '@/components/page-layout'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { userAPI, accountAPI } from '@/lib/api'
import type { User, BindUser } from '@/lib/types'
import { User as UserIcon, Save, Loader2, Link2, Unlink } from 'lucide-react'

// 平台图标映射
function PlatformIcon({ type }: { type: string }) {
  if (type === 'github') return (
    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
  )
  return <Link2 className="h-4 w-4" />
}

export default function UserSettingsPage() {
  const params = useParams<{ username: string }>()
  const router = useRouter()
  const username = params.username

  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  // 社交账号绑定状态
  const [bindUsers, setBindUsers] = useState<BindUser[]>([])
  const [bindLoading, setBindLoading] = useState(false)
  const [unbinding, setUnbinding] = useState<number | null>(null)

  // 表单字段
  const [name, setName] = useState('')
  const [city, setCity] = useState('')
  const [company, setCompany] = useState('')
  const [github, setGithub] = useState('')
  const [website, setWebsite] = useState('')
  const [introduce, setIntroduce] = useState('')
  const [open, setOpen] = useState('0')

  // 密码修改
  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [curPassword, setCurPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [changingPassword, setChangingPassword] = useState(false)

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

  // 加载社交账号绑定
  const loadBindUsers = useCallback(async () => {
    setBindLoading(true)
    try {
      const data = await accountAPI.getBindUsers()
      setBindUsers(data?.bind_users ?? [])
    } catch {
      // 未登录时不阻断页面，但提示用户
      setBindUsers([])
    } finally {
      setBindLoading(false)
    }
  }, [])

  useEffect(() => {
    loadProfile()
    loadBindUsers()
  }, [loadProfile, loadBindUsers])

  // 解绑社交账号
  const handleUnbind = async (bu: BindUser) => {
    if (!window.confirm(`确定要解绑 ${bu.platform} 账号「${bu.username}」吗？`)) return
    setMessage(null)
    setUnbinding(bu.id)
    try {
      await accountAPI.socialUnbind(bu.id, bu.platform)
      setBindUsers((prev) => prev.filter((b) => b.id !== bu.id))
    } catch (err) {
      setMessage({
        type: 'error',
        text: err instanceof Error ? err.message : '解绑失败，请重试',
      })
    } finally {
      setUnbinding(null)
    }
  }

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

  // 修改密码
  const handleChangePassword = async () => {
    setMessage(null)
    if (!curPassword) {
      setMessage({ type: 'error', text: '请输入当前密码' })
      return
    }
    if (!newPassword || newPassword.length < 6) {
      setMessage({ type: 'error', text: '新密码长度至少6个字符' })
      return
    }
    if (curPassword === newPassword) {
      setMessage({ type: 'error', text: '新密码不能与当前密码相同' })
      return
    }
    setChangingPassword(true)
    try {
      await userAPI.changePassword(curPassword, newPassword)
      setMessage({ type: 'success', text: '密码修改成功' })
      setCurPassword('')
      setNewPassword('')
      setShowPasswordForm(false)
    } catch (err) {
      setMessage({ type: 'error', text: err instanceof Error ? err.message : '密码修改失败' })
    } finally {
      setChangingPassword(false)
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
        <div className="flex justify-end mb-6">
          <Button onClick={handleSave} disabled={saving} className="gap-2">
            {saving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            保存修改
          </Button>
        </div>

        {/* 社交账号绑定 */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">社交账号绑定</CardTitle>
          </CardHeader>
          <CardContent>
            {bindLoading ? (
              <div className="flex items-center justify-center py-6">
                <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
              </div>
            ) : bindUsers.length === 0 ? (
              <p className="text-sm text-muted-foreground py-2">暂未绑定任何社交账号</p>
            ) : (
              <ul className="space-y-3">
                {bindUsers.map((bu) => (
                  <li key={bu.id} className="flex items-center justify-between gap-3">
                    <div className="flex items-center gap-3">
                      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted">
                        <PlatformIcon type={bu.platform} />
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-medium">{bu.username}</span>
                          <Badge variant="outline" className="text-xs capitalize">
                            {bu.platform}
                          </Badge>
                        </div>
                        {bu.name && bu.name !== bu.username && (
                          <p className="text-xs text-muted-foreground">{bu.name}</p>
                        )}
                      </div>
                    </div>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="gap-1.5 text-destructive hover:text-destructive hover:bg-destructive/10"
                      disabled={unbinding === bu.id}
                      onClick={() => handleUnbind(bu)}
                    >
                      {unbinding === bu.id ? (
                        <Loader2 className="h-3.5 w-3.5 animate-spin" />
                      ) : (
                        <Unlink className="h-3.5 w-3.5" />
                      )}
                      解绑
                    </Button>
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>

        {/* 修改密码 */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Lock className="h-4 w-4" />
              修改密码
            </CardTitle>
          </CardHeader>
          <CardContent>
            {!showPasswordForm ? (
              <Button variant="outline" onClick={() => setShowPasswordForm(true)}>
                修改密码
              </Button>
            ) : (
              <div className="space-y-3">
                <div className="space-y-2">
                  <Label htmlFor="curPassword">当前密码</Label>
                  <Input
                    id="curPassword"
                    type="password"
                    value={curPassword}
                    onChange={(e) => setCurPassword(e.target.value)}
                    placeholder="请输入当前密码"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="newPassword">新密码</Label>
                  <Input
                    id="newPassword"
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder="请输入新密码（至少6个字符）"
                  />
                </div>
                <div className="flex gap-2">
                  <Button onClick={handleChangePassword} disabled={changingPassword} className="gap-2">
                    {changingPassword ? <Loader2 className="h-4 w-4 animate-spin" /> : '确认修改'}
                  </Button>
                  <Button variant="ghost" onClick={() => { setShowPasswordForm(false); setCurPassword(''); setNewPassword('') }}>
                    取消
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </PageLayout>
  )
}
import { Lock } from 'lucide-react'
