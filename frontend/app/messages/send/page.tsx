"use client"

import React, { Suspense, useState, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { userAPI, commentAPI, messageAPI } from "@/lib/api"
import { toast } from "sonner"

interface UserOption {
  username: string
  avatar: string
}

interface ReplyMessage {
  id: number
  content: string
  from: number
  ctime: string
}

function MessageSendLoading() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <main className="mx-auto flex max-w-3xl items-center justify-center px-4 py-8 lg:px-6">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </main>
      <SiteFooter />
    </div>
  )
}

export default function MessageSendPage() {
  return (
    <Suspense fallback={<MessageSendLoading />}>
      <MessageSendForm />
    </Suspense>
  )
}

function MessageSendForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const usernameParam = searchParams.get("username")?.trim() ?? ""
  const replyId = searchParams.get("id") ?? ""
  const { isLoggedIn, isLoading: authLoading } = useAuth()
  const [loading, setLoading] = useState(true)
  const [toUsername, setToUsername] = useState(usernameParam)
  const [replyMessage, setReplyMessage] = useState<ReplyMessage | null>(null)
  const [users, setUsers] = useState<UserOption[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      const query = searchParams.toString()
      const redirect = `/messages/send${query ? `?${query}` : ""}`
      router.replace(`/account/login?redirect=${encodeURIComponent(redirect)}`)
      return
    }
    setLoading(false)
  }, [authLoading, isLoggedIn, router, searchParams])

  useEffect(() => {
    setReplyMessage(null)
    if (!isLoggedIn || !usernameParam || !replyId) return

    let cancelled = false
    async function loadReplyMessage() {
      try {
        const [profile, data] = await Promise.all([
          userAPI.getProfile(usernameParam),
          messageAPI.detail(replyId),
        ])
        if (!cancelled && profile.user.uid === data.message.from) {
          setReplyMessage(data.message)
        }
      } catch {
        if (!cancelled) setReplyMessage(null)
      }
    }

    void loadReplyMessage()
    return () => {
      cancelled = true
    }
  }, [isLoggedIn, replyId, usernameParam])

  async function handleSearchUser(term: string) {
    if (!term.trim()) {
      setUsers([])
      setShowSuggestions(false)
      return
    }

    try {
      const results = await commentAPI.getAtUsers(term)
      setUsers(results)
      setShowSuggestions(results.length > 0)
    } catch {
      setUsers([])
      setShowSuggestions(false)
    }
  }

  function selectUser(user: UserOption) {
    setToUsername(user.username)
    setShowSuggestions(false)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    const username = toUsername.trim()
    if (!username) {
      toast.error("请选择收件人")
      return
    }
    if (!content.trim()) {
      toast.error("请输入私信内容")
      return
    }

    setSubmitting(true)

    try {
      const data = await userAPI.getProfile(username)
      if (!data?.user?.uid) {
        toast.error("收件人不存在")
        return
      }
      await messageAPI.send(data.user.uid, content.trim())
      toast.success("私信发送成功")
      router.back()
    } catch (err: any) {
      toast.error(err.message || "私信发送失败")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return <MessageSendLoading />
  }

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6">
        <div className="mb-6">
          <h1 className="text-xl font-semibold text-foreground">发送私信</h1>
          <p className="mt-1 text-sm text-muted-foreground">向其他用户发送私信</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <form onSubmit={handleSubmit} className="p-5 space-y-5">
            <div className="space-y-1.5 relative">
              <Label htmlFor="toUsername" className="text-sm font-medium">
                收件人 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="toUsername"
                placeholder="输入用户名搜索"
                value={toUsername}
                onChange={(e) => {
                  setToUsername(e.target.value)
                  handleSearchUser(e.target.value)
                }}
                onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
                className="h-10"
              />
              {showSuggestions && users.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-card border border-border rounded-md shadow-lg max-h-60 overflow-y-auto">
                  {users.map((user) => (
                    <button
                      key={user.username}
                      type="button"
                      className="w-full px-4 py-2 text-left text-sm hover:bg-secondary flex items-center gap-2"
                      onClick={() => selectUser(user)}
                    >
                      <img
                        src={user.avatar || "/static/img/avatar.png"}
                        alt={user.username}
                        className="w-6 h-6 rounded-full"
                      />
                      <span>{user.username}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>

            {replyMessage && toUsername.trim() === usernameParam && (
              <div className="rounded-md border border-border bg-muted/40 p-3">
                <p className="text-xs font-medium text-muted-foreground">
                  回复 {usernameParam} 的私信
                </p>
                <p className="mt-1 whitespace-pre-wrap text-sm text-foreground">
                  {replyMessage.content}
                </p>
              </div>
            )}

            <div className="space-y-1.5">
              <Label htmlFor="content" className="text-sm font-medium">
                内容 <span className="text-destructive">*</span>
              </Label>
              <textarea
                id="content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="请输入私信内容..."
                className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                rows={6}
              />
            </div>

            <div className="flex items-center justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => router.back()}
                disabled={submitting}
              >
                取消
              </Button>
              <Button
                type="submit"
                size="sm"
                disabled={submitting || !toUsername.trim() || !content.trim()}
                className="gap-1.5"
              >
                {submitting && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
                {submitting ? "发送中..." : "发送"}
              </Button>
            </div>
          </form>
        </div>
      </main>
      <SiteFooter />
    </div>
  )
}
