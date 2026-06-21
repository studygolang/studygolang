"use client"

import "../../../publish/blocknote.css"
import React, { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
import dynamic from "next/dynamic"
import { Loader2, Info, Check, ChevronsUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
import { useAuth } from "@/lib/auth-context"
import type { TopicNode } from "@/lib/types"

const BlockNoteEditor = dynamic(
  () => import("@/components/blocknote-editor").then((mod) => ({ default: mod.BlockNoteEditor })),
  {
    ssr: false,
    loading: () => <div className="p-4 text-sm text-muted-foreground">加载编辑器...</div>,
  }
)

const API_BASE = ""

type NodeGroup = {
  category: string
  nodes: TopicNode[]
}

interface Topic {
  tid: number
  title: string
  content: string
  nid: number
  tags: string
}

export default function TopicModifyPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [loading, setLoading] = useState(true)
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [nid, setNid] = useState("")
  const [nodeGroups, setNodeGroups] = useState<NodeGroup[]>([])
  const [nodeOpen, setNodeOpen] = useState(false)
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/topics/modify/${id}`)
      return
    }
    fetchTopic()
  }, [id, isLoggedIn, authLoading, router])

  async function fetchTopic() {
    try {
      const res = await fetch(`${API_BASE}/api/v1/topics/${id}/edit`, {
        credentials: "include",
      })
      const json = await res.json()
      if (json.code !== 0) {
        setError(json.msg || "加载失败")
        setLoading(false)
        return
      }

      const topic: Topic = json.data.topic
      setTitle(topic.title || "")
      setContent(topic.content || "")
      setNid(String(topic.nid || ""))
      if (topic.tags) {
        setTags(topic.tags.split(",").filter(Boolean))
      }

      // 解析节点列表
      const rawNodes = json.data.nodes
      if (Array.isArray(rawNodes)) {
        const groups: NodeGroup[] = []
        for (const group of rawNodes) {
          for (const [category, categoryNodes] of Object.entries(group)) {
            if (Array.isArray(categoryNodes)) {
              const nodes: TopicNode[] = (categoryNodes as any[]).map((node) => ({
                nid: node.nid,
                name: node.name,
                ename: node.ename,
                parent: node.pid ?? node.parent ?? 0,
                seq: node.seq || 0,
                intro: node.intro || "",
                logo: node.logo || "",
                show_index: node.show_index ?? false,
              }))
              groups.push({ category, nodes })
            }
          }
        }
        setNodeGroups(groups)
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

  function removeTag(tagToRemove: string) {
    setTags(tags.filter((t) => t !== tagToRemove))
  }

  async function handleSave() {
    setError("")
    if (!title.trim()) { setError("请填写标题"); return }
    if (!nid) { setError("请选择节点"); return }

    setSubmitting(true)
    let succeeded = false
    try {
      const form = new URLSearchParams()
      form.set("title", title.trim())
      form.set("content", content.trim())
      form.set("nid", nid)
      if (tags.length > 0) form.set("tags", tags.join(","))

      const res = await fetch(`${API_BASE}/api/v1/topics/${id}`, {
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
      succeeded = true
      setSuccess(true)
      setTimeout(() => router.push(`/topics/${id}`), 1200)
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      if (!succeeded) {
        setSubmitting(false)
      }
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
          <h1 className="text-xl font-semibold text-foreground">编辑话题</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改话题内容后点击保存</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                标题 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder="请输入话题标题"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            <div className="flex gap-4">
              <div className="w-56 space-y-1.5">
                <Label className="text-sm font-medium">
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
                            .find((node) => String(node.nid) === nid)?.name
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
                                key={node.nid}
                                value={`${node.name} ${node.ename}`}
                                onSelect={() => {
                                  setNid(String(node.nid))
                                  setNodeOpen(false)
                                }}
                              >
                                <Check
                                  className={cn(
                                    "mr-2 h-4 w-4",
                                    nid === String(node.nid) ? "opacity-100" : "opacity-0"
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

              <div className="min-w-0 flex-1 space-y-1.5">
                <Label htmlFor="tag-input" className="text-sm font-medium">
                  标签
                  <span className="ml-1 font-normal text-muted-foreground">（最多5个，回车添加）</span>
                </Label>
                <div className="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm focus-within:ring-2 focus-within:ring-ring">
                  {tags.map((tag) => (
                    <button
                      key={tag}
                      type="button"
                      aria-label={`删除标签 ${tag}`}
                      className="inline-flex items-center gap-1 rounded-md bg-secondary px-2 py-0.5 text-xs text-secondary-foreground cursor-pointer hover:bg-secondary/80 transition-colors"
                      onClick={() => removeTag(tag)}
                    >
                      {tag}
                      <span aria-hidden="true">×</span>
                    </button>
                  ))}
                  {tags.length < 5 && (
                    <input
                      id="tag-input"
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

            <div className="flex items-start gap-2 rounded-md bg-muted/50 px-3 py-2.5 text-xs text-muted-foreground">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-primary" />
              <span>内容支持 <strong className="text-foreground">Markdown</strong> 语法</span>
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
