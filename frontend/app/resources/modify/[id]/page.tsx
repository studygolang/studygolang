"use client"

import React, { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
import { Loader2, Info } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { resourceWriteAPI, resourceAPI } from "@/lib/api"
import { toast } from "sonner"

export default function ResourceModifyPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  // 表单字段
  const [title, setTitle] = useState("")
  const [url, setUrl] = useState("")
  const [content, setContent] = useState("")
  const [catid, setCatid] = useState("")
  const [form, setForm] = useState("")

  // 分类列表
  const [categories, setCategories] = useState<Array<{ id: number; name: string }>>([])

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/resources/modify/${id}`)
      return
    }
    fetchCategories()
    fetchResource()
  }, [id, isLoggedIn, authLoading, router])

  async function fetchCategories() {
    try {
      const data = await resourceAPI.getCategories()
      setCategories(data.categories || [])
    } catch {
      // 忽略错误，使用默认分类
      setCategories([
        { id: 1, name: "Go语言" },
        { id: 2, name: "Web开发" },
        { id: 3, name: "工具" },
        { id: 4, name: "教程" },
        { id: 5, name: "其他" },
      ])
    }
  }

  async function fetchResource() {
    try {
      const data = await resourceWriteAPI.getEdit(id)
      const resource = data.resource
      setTitle(resource.title || "")
      setUrl(resource.url || "")
      setContent(resource.content || "")
      setCatid(String(resource.catid || ""))
      setForm(resource.form || "")
    } catch (err: any) {
      toast.error(err.message || "加载失败")
    } finally {
      setLoading(false)
    }
  }

  async function handleSave() {
    if (!title.trim()) {
      toast.error("请填写标题")
      return
    }
    if (!catid) {
      toast.error("请选择分类")
      return
    }

    setSubmitting(true)
    try {
      await resourceWriteAPI.update(id, {
        title: title.trim(),
        url: url.trim(),
        content: content.trim(),
        catid: catid,
        form: form,
      })
      toast.success("保存成功")
      setTimeout(() => router.push(`/resources/${id}`), 1200)
    } catch (err: any) {
      toast.error(err.message || "保存失败")
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
          <h1 className="text-xl font-semibold text-foreground">编辑资源</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改资源信息后点击保存</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                标题 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder="请输入资源标题"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="catid" className="text-sm font-medium">
                  分类 <span className="text-destructive">*</span>
                </Label>
                <Select value={catid} onValueChange={setCatid}>
                  <SelectTrigger className="h-10">
                    <SelectValue placeholder="选择分类" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((cat) => (
                      <SelectItem key={cat.id} value={String(cat.id)}>
                        {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="form" className="text-sm font-medium">
                  资源类型
                </Label>
                <Select value={form} onValueChange={setForm}>
                  <SelectTrigger className="h-10">
                    <SelectValue placeholder="选择类型" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="只是链接">只是链接</SelectItem>
                    <SelectItem value="包括内容">包括内容</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="url" className="text-sm font-medium">
                链接地址
              </Label>
              <Input
                id="url"
                placeholder="资源链接地址"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                className="h-10"
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="content" className="text-sm font-medium">
                内容描述
              </Label>
              <Textarea
                id="content"
                placeholder="请输入资源描述或内容"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={6}
                className="resize-none"
              />
            </div>

            <div className="flex items-start gap-2 rounded-md bg-muted/50 p-3 text-xs text-muted-foreground">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
              <span>
                带有 <strong className="text-foreground">*</strong> 的为必填项，保存后可以继续编辑。
              </span>
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
