"use client"

import { useEffect, useRef, useState } from "react"
import { useRouter } from "next/navigation"
import {
  MessageSquare,
  FileText,
  FolderGit2,
  Heading2,
  Bold,
  Italic,
  Code,
  Quote,
  List,
  ListOrdered,
  LinkIcon,
  ImageIcon,
  Eye,
  Loader2,
  Info,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import type { TopicNode } from "@/lib/types"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

type ContentType = "topic" | "article" | "project"

const CONTENT_TYPES: { type: ContentType; icon: React.ReactNode; label: string; desc: string }[] = [
  { type: "topic", icon: <MessageSquare className="h-5 w-5" />, label: "主题", desc: "发起一个讨论话题" },
  { type: "article", icon: <FileText className="h-5 w-5" />, label: "文章", desc: "发表原创或翻译文章" },
  { type: "project", icon: <FolderGit2 className="h-5 w-5" />, label: "项目", desc: "分享开源项目" },
]

// Markdown 工具栏按钮配置
const TOOLBAR = [
  { icon: <Heading2 className="h-3.5 w-3.5" />, title: "标题", syntax: (s: string) => `## ${s || "标题"}` },
  { icon: <Bold className="h-3.5 w-3.5" />, title: "粗体", syntax: (s: string) => `**${s || "粗体"}**` },
  { icon: <Italic className="h-3.5 w-3.5" />, title: "斜体", syntax: (s: string) => `*${s || "斜体"}*` },
  { icon: <Code className="h-3.5 w-3.5" />, title: "代码", syntax: (s: string) => s ? `\`${s}\`` : "```go\n\n```" },
  { icon: <Quote className="h-3.5 w-3.5" />, title: "引用", syntax: (s: string) => `> ${s || "引用"}` },
  { icon: <List className="h-3.5 w-3.5" />, title: "无序列表", syntax: (_: string) => "- 列表项" },
  { icon: <ListOrdered className="h-3.5 w-3.5" />, title: "有序列表", syntax: (_: string) => "1. 列表项" },
  { icon: <LinkIcon className="h-3.5 w-3.5" />, title: "链接", syntax: (s: string) => `[${s || "链接文字"}](url)` },
  { icon: <ImageIcon className="h-3.5 w-3.5" />, title: "图片", syntax: (_: string) => "![alt](image-url)" },
]

async function fetchNodes(): Promise<TopicNode[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/nodes`)
    if (!res.ok) return []
    const json = await res.json()
    return json.code === 0 ? (json.data as TopicNode[]) : []
  } catch {
    return []
  }
}

export default function PublishPage() {
  const router = useRouter()

  const [contentType, setContentType] = useState<ContentType>("topic")
  const [nodes, setNodes] = useState<TopicNode[]>([])
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [nid, setNid] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [preview, setPreview] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // 未登录重定向
  useEffect(() => {
    const token = localStorage.getItem("token")
    if (!token) {
      router.replace("/account/login?redirect=/publish")
    }
  }, [router])

  // 加载节点
  useEffect(() => {
    fetchNodes().then(setNodes)
  }, [])

  // Ctrl/Cmd+Enter 提交
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        handlePublish()
      }
    }
    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  })

  // 工具栏插入 Markdown
  function insertMarkdown(syntaxFn: (sel: string) => string) {
    const ta = textareaRef.current
    if (!ta) return
    const start = ta.selectionStart
    const end = ta.selectionEnd
    const selected = content.slice(start, end)
    const inserted = syntaxFn(selected)
    const next = content.slice(0, start) + inserted + content.slice(end)
    setContent(next)
    requestAnimationFrame(() => {
      ta.focus()
      ta.setSelectionRange(start + inserted.length, start + inserted.length)
    })
  }

  // 标签输入
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

  async function handlePublish() {
    setError("")
    if (!title.trim()) { setError("请填写标题"); return }
    if (contentType === "topic" && !nid) { setError("请选择节点"); return }

    const token = localStorage.getItem("token")
    if (!token) { router.replace("/account/login?redirect=/publish"); return }

    setSubmitting(true)
    try {
      if (contentType === "topic") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("content", content.trim())
        form.set("nid", nid)
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/topics`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
            "X-Token": token,
          },
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.message || "发布失败"); return }
        router.push(json.data?.tid ? `/topics/${json.data.tid}` : "/topics")
      } else {
        // 文章/项目暂未实现 API，提示用户
        setError(`${contentType === "article" ? "文章" : "项目"}发布功能即将上线`)
      }
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setSubmitting(false)
    }
  }

  function handleCancel() {
    router.back()
  }

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6">
        {/* 页面标题 */}
        <div className="mb-6">
          <h1 className="text-xl font-semibold text-foreground">发布内容</h1>
          <p className="mt-1 text-sm text-muted-foreground">选择内容类型，与社区分享你的知识和发现</p>
        </div>

        {/* 内容类型选择 */}
        <div className="mb-6 grid grid-cols-3 gap-3">
          {CONTENT_TYPES.map(({ type, icon, label, desc }) => (
            <button
              key={type}
              type="button"
              onClick={() => setContentType(type)}
              className={`flex flex-col items-center gap-2 rounded-lg border-2 px-4 py-5 text-center transition-all ${
                contentType === type
                  ? "border-primary bg-primary/5 text-primary"
                  : "border-border bg-card text-muted-foreground hover:border-primary/40 hover:text-foreground"
              }`}
            >
              {icon}
              <span className="text-sm font-medium">{label}</span>
              <span className="text-xs text-muted-foreground">{desc}</span>
            </button>
          ))}
        </div>

        {/* 表单卡片 */}
        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            {/* 标题 */}
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                标题 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder={
                  contentType === "topic"
                    ? "请输入主题标题，简明扼要地描述你的问题或话题"
                    : contentType === "article"
                    ? "请输入文章标题"
                    : "请输入项目名称"
                }
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            {/* 节点 + 标签（并排） */}
            <div className="flex gap-4">
              {/* 节点 */}
              {contentType === "topic" && (
                <div className="w-44 space-y-1.5">
                  <Label htmlFor="nid" className="text-sm font-medium">
                    节点 <span className="text-destructive">*</span>
                  </Label>
                  <Select value={nid} onValueChange={setNid}>
                    <SelectTrigger id="nid" className="h-10">
                      <SelectValue placeholder="选择节点" />
                    </SelectTrigger>
                    <SelectContent>
                      {nodes.map((node) => (
                        <SelectItem key={node.id} value={String(node.id)}>
                          {node.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}

              {/* 标签 */}
              <div className="min-w-0 flex-1 space-y-1.5">
                <Label htmlFor="tags" className="text-sm font-medium">
                  标签
                  <span className="ml-1 font-normal text-muted-foreground">（最多5个，回车添加）</span>
                </Label>
                <div className="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-0">
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
            </div>

            {/* Markdown 工具栏 + 内容区 */}
            <div className="space-y-0">
              {/* 工具栏 */}
              <div className="flex items-center justify-between rounded-t-md border border-b-0 border-border bg-muted/40 px-3 py-1.5">
                <div className="flex items-center gap-0.5">
                  {TOOLBAR.map((btn, i) => (
                    <button
                      key={i}
                      type="button"
                      title={btn.title}
                      onClick={() => insertMarkdown(btn.syntax)}
                      className="rounded p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                    >
                      {btn.icon}
                    </button>
                  ))}
                </div>
                <button
                  type="button"
                  onClick={() => setPreview(!preview)}
                  className={`flex items-center gap-1 rounded px-2 py-1 text-xs transition-colors ${
                    preview
                      ? "bg-primary/10 text-primary"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Eye className="h-3.5 w-3.5" />
                  预览
                </button>
              </div>

              {/* 编辑 / 预览 */}
              {preview ? (
                <div className="min-h-[200px] rounded-b-md border border-border bg-background p-4">
                  {content ? (
                    <div
                      className="prose prose-sm max-w-none text-sm text-foreground"
                      dangerouslySetInnerHTML={{
                        __html: content
                          .replace(/&/g, "&amp;")
                          .replace(/</g, "&lt;")
                          .replace(/>/g, "&gt;")
                          .replace(/\n/g, "<br/>"),
                      }}
                    />
                  ) : (
                    <p className="text-sm text-muted-foreground">暂无内容，切换到编辑模式输入</p>
                  )}
                </div>
              ) : (
                <Textarea
                  ref={textareaRef}
                  placeholder={"请详细描述你的问题或话题……\n\n支持 Markdown 语法，可使用上方工具栏快捷插入格式。"}
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  rows={12}
                  className="rounded-t-none border-t-0 font-mono text-sm focus-visible:ring-0 focus-visible:ring-offset-0 resize-y"
                />
              )}
            </div>

            {/* 提示信息 */}
            <div className="flex items-start gap-2 rounded-md bg-muted/50 px-3 py-2.5 text-xs text-muted-foreground">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-primary" />
              <span>
                内容支持 <strong className="text-foreground">Markdown</strong> 语法，发布前请确认标题准确、内容完整、节点正确，发布后可以继续编辑。
              </span>
            </div>

            {/* 错误提示 */}
            {error && <p className="text-sm text-destructive">{error}</p>}
          </div>

          {/* 底部操作栏 */}
          <div className="flex items-center justify-between border-t border-border px-5 py-4">
            <button
              type="button"
              onClick={handleCancel}
              className="text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              取消
            </button>
            <div className="flex items-center gap-3">
              <Button variant="outline" size="sm" disabled={submitting}>
                保存草稿
              </Button>
              <Button size="sm" disabled={submitting} onClick={handlePublish} className="gap-2">
                {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
                {submitting ? "发布中..." : "发布"}
              </Button>
            </div>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
