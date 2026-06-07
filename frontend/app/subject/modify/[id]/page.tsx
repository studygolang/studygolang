"use client"

import { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { subjectAPI, fetchAPI } from "@/lib/api"
import { toast } from "sonner"

interface Subject {
  id: number
  name: string
  description: string
  cover: string
  tags: string
}

export default function SubjectModifyPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [cover, setCover] = useState("")
  const [tags, setTags] = useState("")

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/subject/modify/${id}`)
      return
    }
    loadSubject()
  }, [id, isLoggedIn, authLoading, router])

  async function loadSubject() {
    try {
      const data = await fetchAPI<{ subject: Subject }>(`/subject/${id}`, {
        credentials: "include",
      })
      const subject: Subject = data.subject || (data as unknown as Subject)
      setName(subject.name || "")
      setDescription(subject.description || "")
      setCover(subject.cover || "")
      setTags(subject.tags || "")
    } catch (error) {
      toast.error("加载失败")
      router.push("/")
    } finally {
      setLoading(false)
    }
  }

  async function handleSubmit() {
    if (!name.trim()) {
      toast.error("请填写专栏名称")
      return
    }

    setSubmitting(true)
    try {
      await subjectAPI.modify({
        id,
        name: name.trim(),
        description: description.trim(),
        cover: cover.trim(),
        tags: tags.trim(),
      })
      toast.success("修改成功")
      router.push(`/subject/${id}`)
    } catch (error) {
      const message = error instanceof Error ? error.message : "修改失败"
      toast.error(message)
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
          <h1 className="text-xl font-semibold text-foreground">编辑专栏</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            修改专栏信息后点击保存
          </p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="space-y-4 p-5">
            <div className="space-y-1.5">
              <Label htmlFor="name" className="text-sm font-medium">
                专栏名称 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="name"
                placeholder="请输入专栏名称"
                value={name}
                onChange={(e) => setName(e.target.value)}
                maxLength={100}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="description" className="text-sm font-medium">
                专栏描述
              </Label>
              <Input
                id="description"
                placeholder="请输入专栏描述"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="cover" className="text-sm font-medium">
                封面图片
              </Label>
              <Input
                id="cover"
                placeholder="请输入封面图片URL"
                value={cover}
                onChange={(e) => setCover(e.target.value)}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="tags" className="text-sm font-medium">
                标签
              </Label>
              <Input
                id="tags"
                placeholder="多个标签用逗号分隔"
                value={tags}
                onChange={(e) => setTags(e.target.value)}
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
            <Button onClick={handleSubmit} disabled={submitting} className="gap-2">
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
