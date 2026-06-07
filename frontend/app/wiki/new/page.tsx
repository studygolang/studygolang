"use client"

import React, { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { Loader2 } from "lucide-react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { wikiAPI } from "@/lib/api"

export default function WikiNewPage() {
  const router = useRouter()
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [title, setTitle] = useState("")
  const [uri, setUri] = useState("")
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace("/account/login?redirect=/wiki/new")
    }
  }, [authLoading, isLoggedIn, router])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!title.trim()) {
      toast.error("请填写标题")
      return
    }
    if (!uri.trim()) {
      toast.error("请填写 URI（用于 URL 路径）")
      return
    }
    if (!/^[a-zA-Z0-9-]+$/.test(uri.trim())) {
      toast.error("URI 仅支持字母、数字和连字符")
      return
    }
    if (!content.trim()) {
      toast.error("请填写内容")
      return
    }

    setSubmitting(true)
    try {
      await wikiAPI.create({ title: title.trim(), content: content.trim(), uri: uri.trim() })
      toast.success("Wiki 创建成功")
      router.push("/wiki/" + uri.trim())
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "创建失败，请稍后重试"
      toast.error(msg)
    } finally {
      setSubmitting(false)
    }
  }

  if (authLoading) {
    return (
      <div className="min-h-screen bg-background">
        <SiteHeader />
        <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6 flex items-center justify-center">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </main>
        <SiteFooter />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6">
        <div className="mb-6">
          <h1 className="text-xl font-semibold text-foreground">创建 Wiki</h1>
          <p className="mt-1 text-sm text-muted-foreground">填写以下信息，创建新的 Wiki 页面</p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="rounded-lg border border-border bg-card">
            <div className="p-5 space-y-5">
              <div className="space-y-1.5">
                <Label htmlFor="title" className="text-sm font-medium">
                  标题 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="title"
                  placeholder="请输入 Wiki 标题"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  maxLength={100}
                  className="h-10"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="uri" className="text-sm font-medium">
                  URI <span className="text-destructive">*</span>
                  <span className="ml-1 font-normal text-muted-foreground">（用于 URL 路径，仅支持字母、数字和连字符）</span>
                </Label>
                <Input
                  id="uri"
                  placeholder="例如: go-basic-syntax"
                  value={uri}
                  onChange={(e) => setUri(e.target.value)}
                  maxLength={100}
                  className="h-10"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="content" className="text-sm font-medium">
                  内容 <span className="text-destructive">*</span>
                </Label>
                <Textarea
                  id="content"
                  placeholder="请输入 Wiki 内容，支持 Markdown 语法"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  rows={12}
                  className="resize-y font-mono text-sm"
                />
              </div>
            </div>

            <div className="flex items-center justify-between border-t border-border px-5 py-4">
              <button
                type="button"
                onClick={() => router.back()}
                className="text-sm text-muted-foreground transition-colors hover:text-foreground"
              >
                取消
              </button>
              <Button type="submit" size="sm" disabled={submitting} className="gap-2">
                {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
                {submitting ? "提交中..." : "创建 Wiki"}
              </Button>
            </div>
          </div>
        </form>
      </main>
      <SiteFooter />
    </div>
  )
}
