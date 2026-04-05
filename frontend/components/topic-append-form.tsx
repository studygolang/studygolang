"use client"

import React, { useState } from "react"
import { useRouter } from "next/navigation"

interface TopicAppendFormProps {
  tid: number
}

export function TopicAppendForm({ tid }: TopicAppendFormProps) {
  const [expanded, setExpanded] = useState(false)
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const router = useRouter()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim() || submitting) return

    setSubmitting(true)
    setError("")

    try {
      const res = await fetch(`/api/v1/topics/${tid}/append`, {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `content=${encodeURIComponent(content)}`,
        credentials: "include",
      })

      const data = await res.json()
      if (data.code === 0) {
        setContent("")
        setExpanded(false)
        router.refresh()
      } else {
        setError(data.msg || "附言失败")
      }
    } catch {
      setError("网络错误，请重试")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="mt-4">
      {!expanded ? (
        <button
          onClick={() => setExpanded(true)}
          className="text-sm text-muted-foreground hover:text-primary"
        >
          + 添加附言
        </button>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-3">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="附言内容（最多 3 条附言）"
            className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
            rows={4}
          />
          {error && <p className="text-xs text-destructive">{error}</p>}
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={submitting || !content.trim()}
              className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              {submitting ? "提交中..." : "提交附言"}
            </button>
            <button
              type="button"
              onClick={() => {
                setExpanded(false)
                setError("")
              }}
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
            >
              取消
            </button>
          </div>
        </form>
      )}
    </div>
  )
}
