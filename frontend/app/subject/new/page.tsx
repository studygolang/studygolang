"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"

export default function NewSubjectPage() {
  const router = useRouter()
  const [name, setName] = useState("")
  const [intro, setIntro] = useState("")
  const [cover, setCover] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim() || submitting) return

    setSubmitting(true)
    setError("")

    try {
      const params = new URLSearchParams()
      params.set("name", name)
      if (intro) params.set("desc", intro)
      if (cover) params.set("cover", cover)

      const res = await fetch("/api/v1/subject/new", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: params.toString(),
        credentials: "include",
      })

      const data = await res.json()
      if (data.code === 0 && data.data?.sid) {
        router.push(`/subject/${data.data.sid}`)
      } else {
        setError(data.msg || "创建失败")
      }
    } catch {
      setError("网络错误，请重试")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="创建专栏"
        breadcrumbs={[
          { label: "专栏", href: "/subject" },
          { label: "创建专栏" },
        ]}
      />

      <Card>
        <CardContent className="p-6">
          {error && (
            <div className="mb-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium">
                专栏名称 <span className="text-destructive">*</span>
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="输入专栏名称"
                className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                required
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium">简介</label>
              <textarea
                value={intro}
                onChange={(e) => setIntro(e.target.value)}
                placeholder="专栏简介"
                className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                rows={4}
              />
            </div>

            <div>
              <label className="mb-1 block text-sm font-medium">封面 URL</label>
              <input
                type="text"
                value={cover}
                onChange={(e) => setCover(e.target.value)}
                placeholder="封面图片 URL（可选）"
                className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
              />
            </div>

            <div className="flex gap-2 pt-2">
              <button
                type="submit"
                disabled={submitting || !name.trim()}
                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
              >
                {submitting ? "创建中..." : "创建专栏"}
              </button>
              <Link
                href="/subject"
                className="rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                取消
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </PageLayout>
  )
}
