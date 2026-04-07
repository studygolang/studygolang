"use client"

import React, { useState } from "react"
import { useRouter } from "next/navigation"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { commentAPI } from "@/lib/api"
import { toast } from "sonner"

interface CommentFormProps {
  objid: number
  objtype: number
  onSuccess: () => void
}

export function CommentForm({ objid, objtype, onSuccess }: CommentFormProps) {
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const router = useRouter()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim() || submitting) return

    setSubmitting(true)

    try {
      await commentAPI.create(objid, objtype, content)
      toast.success("评论发表成功")
      setContent("")
      onSuccess()
    } catch (err: any) {
      if (err.message?.includes("未登录")) {
        router.push("/account/login")
      } else {
        toast.error(err.message || "评论发表失败")
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="发表评论..."
        className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
        rows={3}
      />
      <div className="flex justify-end">
        <Button
          type="submit"
          size="sm"
          disabled={submitting || !content.trim()}
          className="gap-1.5"
        >
          {submitting && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
          {submitting ? "提交中..." : "发表评论"}
        </Button>
      </div>
    </form>
  )
}
