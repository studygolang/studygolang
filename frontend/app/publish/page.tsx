"use client"

import "./blocknote.css"
import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import dynamic from "next/dynamic"
import {
  MessageSquare,
  FileText,
  FolderGit2,
  Loader2,
  Info,
  Check,
  ChevronsUpDown,
  Link2,
  BookOpen,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

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

type ContentType = "topic" | "article" | "project" | "resource" | "book"

const CONTENT_TYPES: { type: ContentType; icon: React.ReactNode; label: string; desc: string }[] = [
  { type: "topic", icon: <MessageSquare className="h-5 w-5" />, label: "主题", desc: "发起一个讨论话题" },
  { type: "article", icon: <FileText className="h-5 w-5" />, label: "文章", desc: "发表原创或翻译文章" },
  { type: "project", icon: <FolderGit2 className="h-5 w-5" />, label: "项目", desc: "分享开源项目" },
  { type: "resource", icon: <Link2 className="h-5 w-5" />, label: "资源", desc: "分享学习资源" },
  { type: "book", icon: <BookOpen className="h-5 w-5" />, label: "图书", desc: "推荐 Go 相关图书" },
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
  return (
    <Suspense fallback={<div className="min-h-screen bg-background flex items-center justify-center"><p>加载中...</p></div>}>
      <PublishContent />
    </Suspense>
  )
}

function PublishContent() {
  const router = useRouter()
  const searchParams = useSearchParams()

  // 从 URL 参数读取初始内容类型
  const initialType = (searchParams.get("type") || "topic") as ContentType
  const validTypes: ContentType[] = ["topic", "article", "project", "resource", "book"]

  const [contentType, setContentType] = useState<ContentType>(
    validTypes.includes(initialType) ? initialType : "topic"
  )
  const [nodeGroups, setNodeGroups] = useState<NodeGroup[]>([])
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [nid, setNid] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [nodeOpen, setNodeOpen] = useState(false)

  // 项目特定字段
  const [src, setSrc] = useState("")          // 源码地址（必填）
  const [category, setCategory] = useState("") // 项目分类
  const [home, setHome] = useState("")         // 项目主页
  const [doc, setDoc] = useState("")           // 文档地址
  const [download, setDownload] = useState("") // 下载地址
  const [licence, setLicence] = useState("")   // 开源协议
  const [lang, setLang] = useState("Go")       // 开发语言，默认 Go
  const [os, setOs] = useState("")             // 操作系统

  // 资源特定字段
  const [resourceForm, setResourceForm] = useState("link") // link 或 content
  const [resourceUrl, setResourceUrl] = useState("")       // 资源链接
  const [resourceCatid, setResourceCatid] = useState("")   // 资源分类

  // 图书特定字段
  const [bookAuthor, setBookAuthor] = useState("")       // 作者
  const [bookTranslator, setBookTranslator] = useState("") // 译者
  const [bookCover, setBookCover] = useState("")         // 封面
  const [bookPubDate, setBookPubDate] = useState("")     // 出版日期
  const [bookLang, setBookLang] = useState("中文")        // 语言
  const [bookIsFree, setBookIsFree] = useState(false)    // 是否免费
  const [bookOnlineUrl, setBookOnlineUrl] = useState("") // 在线阅读
  const [bookDownloadUrl, setBookDownloadUrl] = useState("") // 下载地址
  const [bookBuyUrl, setBookBuyUrl] = useState("")       // 购买地址
  const [bookPrice, setBookPrice] = useState("")         // 价格

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
    if (contentType === "project" && !src.trim()) { setError("请填写源码地址"); return }
    if (contentType === "resource" && resourceForm === "link" && !resourceUrl.trim()) { setError("请填写资源链接"); return }
    if (contentType === "resource" && resourceForm === "content" && !content.trim()) { setError("请填写资源内容"); return }
    if (contentType === "book" && !bookAuthor.trim()) { setError("请填写作者"); return }

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
      } else if (contentType === "article") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("content", content.trim())
        form.set("txt", content.trim())
        form.set("cover", "")
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/articles`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/articles/${json.data.id}` : "/articles")
      } else if (contentType === "project") {
        const form = new URLSearchParams()
        form.set("name", title.trim())
        form.set("src", src.trim())
        form.set("desc", content.trim())
        if (category.trim()) form.set("category", category.trim())
        if (home.trim()) form.set("home", home.trim())
        if (doc.trim()) form.set("doc", doc.trim())
        if (download.trim()) form.set("download", download.trim())
        if (licence.trim()) form.set("licence", licence.trim())
        if (lang.trim()) form.set("lang", lang.trim())
        if (os.trim()) form.set("os", os.trim())
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/projects`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.uri ? `/projects/${json.data.uri}` : "/projects")
      } else if (contentType === "resource") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("form", resourceForm === "link" ? "只是链接" : "包括内容")
        if (resourceForm === "link") {
          form.set("url", resourceUrl.trim())
        } else {
          form.set("content", content.trim())
        }
        if (resourceCatid) form.set("catid", resourceCatid)
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/resources`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/resources/${json.data.id}` : "/resources")
      } else if (contentType === "book") {
        const form = new URLSearchParams()
        form.set("name", title.trim())
        form.set("author", bookAuthor.trim())
        form.set("desc", content.trim())
        if (bookTranslator.trim()) form.set("translator", bookTranslator.trim())
        if (bookCover.trim()) form.set("cover", bookCover.trim())
        if (bookPubDate.trim()) form.set("pub_date", bookPubDate.trim())
        if (bookLang.trim()) form.set("lang", bookLang.trim())
        form.set("is_free", bookIsFree ? "1" : "0")
        if (bookOnlineUrl.trim()) form.set("online_url", bookOnlineUrl.trim())
        if (bookDownloadUrl.trim()) form.set("download_url", bookDownloadUrl.trim())
        if (bookBuyUrl.trim()) form.set("buy_url", bookBuyUrl.trim())
        if (bookPrice.trim()) form.set("price", bookPrice.trim())
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/books`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/books/${json.data.id}` : "/books")
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
          <div className="p-5 space-y-5">
            {/* 标题 */}
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                {contentType === "project" ? "项目名称" : contentType === "book" ? "书名" : "标题"} <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder={
                  contentType === "topic"
                    ? "请输入主题标题，简明扼要地描述你的问题或话题"
                    : contentType === "article"
                    ? "请输入文章标题"
                    : contentType === "resource"
                    ? "请输入资源标题"
                    : contentType === "book"
                    ? "请输入图书名称"
                    : "请输入项目名称"
                }
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            {/* 项目特定字段 */}
            {contentType === "project" && (
              <>
                {/* 源码地址 */}
                <div className="space-y-1.5">
                  <Label htmlFor="src" className="text-sm font-medium">
                    源码地址 <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="src"
                    placeholder="例如：https://github.com/username/repo"
                    value={src}
                    onChange={(e) => setSrc(e.target.value)}
                    className="h-10"
                  />
                </div>

                {/* 项目分类 + 开发语言 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="category" className="text-sm font-medium">项目分类</Label>
                    <Input
                      id="category"
                      placeholder="例如：Web框架、CLI工具"
                      value={category}
                      onChange={(e) => setCategory(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="lang" className="text-sm font-medium">开发语言</Label>
                    <Select value={lang} onValueChange={setLang}>
                      <SelectTrigger className="h-10">
                        <SelectValue placeholder="选择语言" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Go">Go</SelectItem>
                        <SelectItem value="Rust">Rust</SelectItem>
                        <SelectItem value="Python">Python</SelectItem>
                        <SelectItem value="JavaScript">JavaScript</SelectItem>
                        <SelectItem value="TypeScript">TypeScript</SelectItem>
                        <SelectItem value="C">C</SelectItem>
                        <SelectItem value="C++">C++</SelectItem>
                        <SelectItem value="Java">Java</SelectItem>
                        <SelectItem value="其他">其他</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                {/* 项目主页 + 文档地址 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="home" className="text-sm font-medium">项目主页</Label>
                    <Input
                      id="home"
                      placeholder="https://example.com"
                      value={home}
                      onChange={(e) => setHome(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="doc" className="text-sm font-medium">文档地址</Label>
                    <Input
                      id="doc"
                      placeholder="https://docs.example.com"
                      value={doc}
                      onChange={(e) => setDoc(e.target.value)}
                      className="h-10"
                    />
                  </div>
                </div>

                {/* 开源协议 + 操作系统 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="licence" className="text-sm font-medium">开源协议</Label>
                    <Select value={licence} onValueChange={setLicence}>
                      <SelectTrigger className="h-10">
                        <SelectValue placeholder="选择协议" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="MIT">MIT</SelectItem>
                        <SelectItem value="Apache-2.0">Apache-2.0</SelectItem>
                        <SelectItem value="GPL-3.0">GPL-3.0</SelectItem>
                        <SelectItem value="BSD-3-Clause">BSD-3-Clause</SelectItem>
                        <SelectItem value="MPL-2.0">MPL-2.0</SelectItem>
                        <SelectItem value="其他">其他</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="os" className="text-sm font-medium">操作系统</Label>
                    <Select value={os} onValueChange={setOs}>
                      <SelectTrigger className="h-10">
                        <SelectValue placeholder="选择操作系统" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="跨平台">跨平台</SelectItem>
                        <SelectItem value="Linux">Linux</SelectItem>
                        <SelectItem value="macOS">macOS</SelectItem>
                        <SelectItem value="Windows">Windows</SelectItem>
                        <SelectItem value="其他">其他</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </>
            )}

            {/* 资源特定字段 */}
            {contentType === "resource" && (
              <>
                {/* 资源类型选择 */}
                <div className="space-y-1.5">
                  <Label className="text-sm font-medium">资源类型 <span className="text-destructive">*</span></Label>
                  <div className="flex gap-4">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="resourceForm"
                        value="link"
                        checked={resourceForm === "link"}
                        onChange={() => setResourceForm("link")}
                        className="h-4 w-4"
                      />
                      <span className="text-sm">只是链接</span>
                    </label>
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="radio"
                        name="resourceForm"
                        value="content"
                        checked={resourceForm === "content"}
                        onChange={() => setResourceForm("content")}
                        className="h-4 w-4"
                      />
                      <span className="text-sm">包括内容</span>
                    </label>
                  </div>
                </div>

                {/* 链接地址（只是链接类型时显示） */}
                {resourceForm === "link" && (
                  <div className="space-y-1.5">
                    <Label htmlFor="resourceUrl" className="text-sm font-medium">
                      资源链接 <span className="text-destructive">*</span>
                    </Label>
                    <Input
                      id="resourceUrl"
                      placeholder="https://example.com/resource"
                      value={resourceUrl}
                      onChange={(e) => setResourceUrl(e.target.value)}
                      className="h-10"
                    />
                  </div>
                )}
              </>
            )}

            {/* 图书特定字段 */}
            {contentType === "book" && (
              <>
                {/* 作者 + 译者 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookAuthor" className="text-sm font-medium">
                      作者 <span className="text-destructive">*</span>
                    </Label>
                    <Input
                      id="bookAuthor"
                      placeholder="作者姓名"
                      value={bookAuthor}
                      onChange={(e) => setBookAuthor(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookTranslator" className="text-sm font-medium">译者</Label>
                    <Input
                      id="bookTranslator"
                      placeholder="译者姓名（如有）"
                      value={bookTranslator}
                      onChange={(e) => setBookTranslator(e.target.value)}
                      className="h-10"
                    />
                  </div>
                </div>

                {/* 封面 + 出版日期 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookCover" className="text-sm font-medium">封面图片</Label>
                    <Input
                      id="bookCover"
                      placeholder="封面图片URL"
                      value={bookCover}
                      onChange={(e) => setBookCover(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookPubDate" className="text-sm font-medium">出版日期</Label>
                    <Input
                      id="bookPubDate"
                      placeholder="例如：2024-01"
                      value={bookPubDate}
                      onChange={(e) => setBookPubDate(e.target.value)}
                      className="h-10"
                    />
                  </div>
                </div>

                {/* 语言 + 是否免费 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookLang" className="text-sm font-medium">语言</Label>
                    <Select value={bookLang} onValueChange={setBookLang}>
                      <SelectTrigger className="h-10">
                        <SelectValue placeholder="选择语言" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="中文">中文</SelectItem>
                        <SelectItem value="英文">英文</SelectItem>
                        <SelectItem value="中英双语">中英双语</SelectItem>
                        <SelectItem value="其他">其他</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label className="text-sm font-medium">是否免费</Label>
                    <label className="flex items-center gap-2 h-10 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={bookIsFree}
                        onChange={(e) => setBookIsFree(e.target.checked)}
                        className="h-4 w-4"
                      />
                      <span className="text-sm">免费资源</span>
                    </label>
                  </div>
                </div>

                {/* 在线阅读 + 下载地址 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookOnlineUrl" className="text-sm font-medium">在线阅读</Label>
                    <Input
                      id="bookOnlineUrl"
                      placeholder="在线阅读地址"
                      value={bookOnlineUrl}
                      onChange={(e) => setBookOnlineUrl(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookDownloadUrl" className="text-sm font-medium">下载地址</Label>
                    <Input
                      id="bookDownloadUrl"
                      placeholder="电子版下载地址"
                      value={bookDownloadUrl}
                      onChange={(e) => setBookDownloadUrl(e.target.value)}
                      className="h-10"
                    />
                  </div>
                </div>

                {/* 购买地址 + 价格 */}
                <div className="flex gap-4">
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookBuyUrl" className="text-sm font-medium">购买地址</Label>
                    <Input
                      id="bookBuyUrl"
                      placeholder="购买链接"
                      value={bookBuyUrl}
                      onChange={(e) => setBookBuyUrl(e.target.value)}
                      className="h-10"
                    />
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <Label htmlFor="bookPrice" className="text-sm font-medium">价格</Label>
                    <Input
                      id="bookPrice"
                      placeholder="例如：99.00"
                      value={bookPrice}
                      onChange={(e) => setBookPrice(e.target.value)}
                      className="h-10"
                    />
                  </div>
                </div>
              </>
            )}

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
            {/* 资源类型只在"包括内容"时显示，图书总是显示 */}
            {(contentType !== "resource" || resourceForm === "content") && contentType !== "book" && (
              <div className="space-y-1.5">
                <Label className="text-sm font-medium">
                  {contentType === "project" ? "项目描述" : contentType === "resource" ? "资源内容" : "内容"} <span className="text-destructive">*</span>
                </Label>
                <div className="rounded-md border border-border overflow-hidden">
                  <BlockNoteEditor
                    content={content}
                    onChange={setContent}
                    placeholder={
                      contentType === "project"
                        ? "请详细描述项目的功能特点、使用场景……"
                        : contentType === "resource"
                        ? "请详细描述资源内容……"
                        : "请详细描述你的问题或话题……"
                    }
                  />
                </div>
              </div>
            )}

            {/* 图书简介 */}
            {contentType === "book" && (
              <div className="space-y-1.5">
                <Label className="text-sm font-medium">
                  图书简介 <span className="text-destructive">*</span>
                </Label>
                <div className="rounded-md border border-border overflow-hidden">
                  <BlockNoteEditor
                    content={content}
                    onChange={setContent}
                    placeholder="请简要介绍这本图书的主要内容、适合的读者群体……"
                  />
                </div>
              </div>
            )}

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
