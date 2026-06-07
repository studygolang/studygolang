"use client"

import React, { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
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

interface WikiEdit {
  id: number
  title: string
  content: string
  uri: string
}

export default function WikiEditPage() {
  const router = useRouter()
  const params = useParams()
  const uri = params.uri as string
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [loading, setLoading] = useState(true)
  const [wiki, setWiki] = useState<WikiEdit | null>(null)
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/wiki/${uri}/edit`)
      return
    }
    loadWiki()
  }, [uri, isLoggedIn, authLoading, router])

  async function loadWiki() {
    try {
      const data = await wikiAPI.getEdit(uri)
      const w = data.wiki
      setWiki(w)
      setTitle(w.title || "")
      setContent(w.content || "")
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "加载失败"
      toast.error(msg)
    } finally {
      setLoading(false)
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!wiki) return

    if (!title.trim()) {
      toast.error("请填写标题")
      return
    }
    if (!content.trim()) {
      toast.error("请填写内容")
      return
    }

    setSubmitting(true)
    try {
      await wikiAPI.update(wiki.id, { title: title.trim(), content: content.trim() })
      toast.success("Wiki 更新成功")
      router.push("/wiki/" + uri)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "更新失败，请稍后重试"
      toast.error(msg)
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
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
          <h1 className="text-xl font-semibold text-foreground">编辑 Wiki</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改 Wiki 内容后点击保存</p>
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
                <Label htmlFor="uri-display" className="text-sm font-medium">
                  URI
                  <span className="ml-1 font-normal text-muted-foreground">（不可修改）</span>
                </Label>
                <Input
                  id="uri-display"
                  value={wiki?.uri || ""}
                  disabled
                  className="h-10 bg-muted/50"
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
                {submitting ? "保存中..." : "保存"}
              </Button>
            </div>
          </div>
        </form>
      </main>
      <SiteFooter />
    </div>
  )
}
