"use client"

import React, { useState } from "react"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { commentAPI } from "@/lib/api"
import { toast } from "sonner"

interface CommentEditFormProps {
  cid: number
  initialContent: string
  onCancel: () => void
  onSuccess: () => void
}

export function CommentEditForm({ cid, initialContent, onCancel, onSuccess }: CommentEditFormProps) {
  const [content, setContent] = useState(initialContent)
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim() || submitting) return

    setSubmitting(true)

    try {
      await commentAPI.update(cid, content.trim())
      toast.success("评论修改成功")
      onSuccess()
    } catch (err: any) {
      toast.error(err.message || "评论修改失败")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
        rows={3}
      />
      <div className="flex gap-2">
        <Button
          type="submit"
          size="sm"
          disabled={submitting || !content.trim()}
          className="gap-1.5"
        >
          {submitting && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
          {submitting ? "保存中..." : "保存"}
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={onCancel}
          disabled={submitting}
        >
          取消
        </Button>
      </div>
    </form>
  )
}
