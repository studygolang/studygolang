"use client"

import "../../../publish/blocknote.css"
import React, { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
import dynamic from "next/dynamic"
import { Loader2, Info } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"

const BlockNoteEditor = dynamic(
  () => import("@/components/blocknote-editor").then((mod) => ({ default: mod.BlockNoteEditor })),
  {
    ssr: false,
    loading: () => <div className="p-4 text-sm text-muted-foreground">加载编辑器...</div>,
  }
)

const API_BASE = ""

interface Article {
  id: number
  title: string
  content: string
  txt: string
  tags: string
}

export default function ArticleModifyPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string

  const [loading, setLoading] = useState(true)
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    const uid = localStorage.getItem("uid")
    if (!uid) {
      router.replace(`/account/login?redirect=/articles/modify/${id}`)
      return
    }
    fetchArticle()
  }, [id])

  async function fetchArticle() {
    try {
      const res = await fetch(`${API_BASE}/api/v1/articles/${id}/edit`, {
        credentials: "include",
      })
      const json = await res.json()
      if (json.code !== 0) {
        setError(json.msg || "加载失败")
        setLoading(false)
        return
      }
      const article: Article = json.data.article
      setTitle(article.title || "")
      setContent(article.txt || article.content || "")
      if (article.tags) {
        setTags(article.tags.split(",").filter(Boolean))
      }
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  function handleTagKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if ((e.key === "Enter" || e.key === ",") && tagInput.trim()) {
      e.preventDefault()
      const tag = tagInput.trim()
      if (tags.length < 5 && !tags.includes(tag)) {
        setTags([...tags, tag])
      }
      setTagInput("")
    } else if (e.key === "Backspace" && !tagInput && tags.length > 0) {
      setTags(tags.slice(0, -1))
    }
  }

  async function handleSave() {
    setError("")
    if (!title.trim()) { setError("请填写标题"); return }
    if (!content.trim()) { setError("请填写内容"); return }

    setSubmitting(true)
    try {
      const form = new URLSearchParams()
      form.set("title", title.trim())
      form.set("content", content.trim())
      if (tags.length > 0) form.set("tags", tags.join(","))

      const res = await fetch(`${API_BASE}/api/v1/articles/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        credentials: "include",
        body: form.toString(),
      })
      const json = await res.json()
      if (json.code !== 0) {
        setError(json.msg || "保存失败")
        return
      }
      setSuccess(true)
      setTimeout(() => router.push(`/articles/${id}`), 1200)
    } catch {
      setError("网络错误，请稍后重试")
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
          <h1 className="text-xl font-semibold text-foreground">编辑文章</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改文章内容后点击保存</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                标题 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder="请输入文章标题"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="tags" className="text-sm font-medium">
                标签
                <span className="ml-1 font-normal text-muted-foreground">（最多5个，回车添加）</span>
              </Label>
              <div className="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm focus-within:ring-2 focus-within:ring-ring">
                {tags.map((tag) => (
                  <Badge
                    key={tag}
                    variant="secondary"
                    className="gap-1 py-0.5 text-xs cursor-pointer"
                    onClick={() => setTags(tags.filter((t) => t !== tag))}
                  >
                    {tag} ×
                  </Badge>
                ))}
                {tags.length < 5 && (
                  <input
                    id="tags"
                    className="min-w-16 flex-1 bg-transparent outline-none placeholder:text-muted-foreground"
                    placeholder={tags.length === 0 ? "输入标签" : ""}
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    onKeyDown={handleTagKeyDown}
                  />
                )}
              </div>
            </div>

            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                内容 <span className="text-destructive">*</span>
              </Label>
              <BlockNoteEditor
                content={content}
                onChange={setContent}
                placeholder="请输入文章内容，支持 Markdown 语法..."
              />
            </div>

            <div className="flex items-start gap-2 rounded-md bg-muted/50 p-3 text-xs text-muted-foreground">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
              <span>
                内容支持 <strong className="text-foreground">Markdown</strong> 语法，保存后可以继续编辑。
              </span>
            </div>

            {error && <p className="text-sm text-destructive">{error}</p>}
            {success && <p className="text-sm text-green-600">保存成功，正在跳转...</p>}
          </div>

          <div className="flex items-center justify-between border-t border-border px-5 py-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              取消
            </button>
            <Button size="sm" disabled={submitting} onClick={handleSave} className="gap-2">
              {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
              {submitting ? "保存中..." : "保存"}
            </Button>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
