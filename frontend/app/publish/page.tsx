"use client"

import "./blocknote.css"
import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import dynamic from "next/dynamic"
import {
  MessageSquare,
  FileText,
  FolderGit2,
  Loader2,
  Info,
  Check,
  ChevronsUpDown,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"

// 动态导入 BlockNote 编辑器,避免 SSR 问题
const BlockNoteEditor = dynamic(
  () => import("@/components/blocknote-editor").then((mod) => ({ default: mod.BlockNoteEditor })),
  {
    ssr: false,
    loading: () => <div className="p-4 text-sm text-muted-foreground">加载编辑器...</div>
  }
)
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { cn } from "@/lib/utils"
import type { TopicNode } from "@/lib/types"

// 客户端组件使用相对路径，通过 next.config.mjs 中的 rewrites 代理到后端
const API_BASE = ""

type ContentType = "topic" | "article" | "project"

const CONTENT_TYPES: { type: ContentType; icon: React.ReactNode; label: string; desc: string }[] = [
  { type: "topic", icon: <MessageSquare className="h-5 w-5" />, label: "主题", desc: "发起一个讨论话题" },
  { type: "article", icon: <FileText className="h-5 w-5" />, label: "文章", desc: "发表原创或翻译文章" },
  { type: "project", icon: <FolderGit2 className="h-5 w-5" />, label: "项目", desc: "分享开源项目" },
]

// 节点分组类型
type NodeGroup = {
  category: string
  nodes: TopicNode[]
}

async function fetchNodes(): Promise<NodeGroup[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/nodes`)
    if (!res.ok) return []
    const json = await res.json()

    if (json.code !== 0 || !Array.isArray(json.data)) return []

    // 后端返回的是分组结构：[{ "Go语言": [...], "StudyGolang": [...] }]
    // 保留分组结构，方便在 UI 中显示层级关系
    const groups: NodeGroup[] = []
    for (const group of json.data) {
      for (const [category, categoryNodes] of Object.entries(group)) {
        if (Array.isArray(categoryNodes)) {
          const nodes: TopicNode[] = categoryNodes.map((node: any) => ({
            id: node.nid,  // 后端使用 nid，前端类型定义使用 id
            name: node.name,
            ename: node.ename,
            parent_id: node.pid,
            seq: node.seq || 0,
            pid: node.pid,
            intro: node.intro || '',
            logo: node.logo || '',
            style: node.style || '',
          }))
          groups.push({ category, nodes })
        }
      }
    }
    return groups
  } catch {
    return []
  }
}

export default function PublishPage() {
  const router = useRouter()

  const [contentType, setContentType] = useState<ContentType>("topic")
  const [nodeGroups, setNodeGroups] = useState<NodeGroup[]>([])
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [nid, setNid] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [nodeOpen, setNodeOpen] = useState(false)

  // 未登录重定向
  useEffect(() => {
    const uid = localStorage.getItem("uid")
    if (!uid) {
      router.replace("/account/login?redirect=/publish")
    }
  }, [router])

  // 加载节点
  useEffect(() => {
    fetchNodes().then(setNodeGroups)
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

    const uid = localStorage.getItem("uid")
    if (!uid) { router.replace("/account/login?redirect=/publish"); return }

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
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
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
        {/* 页面标题和内容类型选择 - 合并为一行 */}
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-semibold text-foreground">发布内容</h1>
            <p className="mt-1 text-sm text-muted-foreground">与社区分享你的知识和发现</p>
          </div>

          {/* 内容类型选择 - 紧凑的标签页样式 */}
          <div className="flex items-center gap-1 rounded-lg border border-border bg-muted/50 p-1">
            {CONTENT_TYPES.map(({ type, icon, label }) => (
              <button
                key={type}
                type="button"
                onClick={() => setContentType(type)}
                className={cn(
                  "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-all",
                  contentType === type
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                )}
              >
                <span className="h-4 w-4">{icon}</span>
                {label}
              </button>
            ))}
          </div>
        </div>

        {/* 表单卡片 */}
        <div className="rounded-lg border border-border bg-card">
          {/* 功能未上线提示 */}
          {contentType !== "topic" && (
            <div className="border-b border-border bg-muted/50 px-5 py-3">
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Info className="h-4 w-4" />
                <span>{contentType === "article" ? "文章" : "项目"}发布功能即将上线，敬请期待</span>
              </div>
            </div>
          )}

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
                <div className="w-56 space-y-1.5">
                  <Label htmlFor="nid" className="text-sm font-medium">
                    节点 <span className="text-destructive">*</span>
                  </Label>
                  <Popover open={nodeOpen} onOpenChange={setNodeOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={nodeOpen}
                        className="w-full justify-between h-10"
                      >
                        {nid
                          ? nodeGroups
                              .flatMap((g) => g.nodes)
                              .find((node) => String(node.id) === nid)?.name
                          : "选择节点"}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[300px] p-0">
                      <Command>
                        <CommandInput placeholder="搜索节点..." />
                        <CommandList>
                          <CommandEmpty>未找到节点</CommandEmpty>
                          {nodeGroups.map((group) => (
                            <CommandGroup key={group.category} heading={group.category}>
                              {group.nodes.map((node) => (
                                <CommandItem
                                  key={node.id}
                                  value={`${node.name} ${node.ename}`}
                                  onSelect={() => {
                                    setNid(String(node.id))
                                    setNodeOpen(false)
                                  }}
                                >
                                  <Check
                                    className={cn(
                                      "mr-2 h-4 w-4",
                                      nid === String(node.id) ? "opacity-100" : "opacity-0"
                                    )}
                                  />
                                  {node.name}
                                </CommandItem>
                              ))}
                            </CommandGroup>
                          ))}
                        </CommandList>
                      </Command>
                    </PopoverContent>
                  </Popover>
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

            {/* BlockNote 编辑器 */}
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">
                内容 <span className="text-destructive">*</span>
              </Label>
              <div className="rounded-md border border-border overflow-hidden">
                <BlockNoteEditor
                  content={content}
                  onChange={setContent}
                  placeholder="请详细描述你的问题或话题……"
                />
              </div>
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
