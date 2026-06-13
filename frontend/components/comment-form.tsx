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
  // 可选的回复目标：填写后发表时会以 `#floor楼 @username ` 前缀，后端据此解析 reply_floor
  replyTo?: { floor: number; username: string } | null
}

export function CommentForm({ objid, objtype, onSuccess, replyTo }: CommentFormProps) {
  const [content, setContent] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const router = useRouter()

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (submitting) return

    // 回复某楼层时自动加上 `#N楼 @username ` 前缀（与后端 decodeCmtContentForShow 正则一致）
    let finalContent = content
    if (replyTo && replyTo.floor > 0) {
      finalContent = `#${replyTo.floor}楼 @${replyTo.username} ${content}`
    }

    if (!finalContent.trim()) return

    setSubmitting(true)

    try {
      await commentAPI.create(objid, objtype, finalContent)
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
      {replyTo && replyTo.floor > 0 && (
        <div className="rounded-md border border-border bg-muted/40 px-3 py-1.5 text-xs text-muted-foreground">
          回复 #{replyTo.floor}楼 @{replyTo.username}
        </div>
      )}
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder={replyTo && replyTo.floor > 0 ? `回复 @${replyTo.username}...` : "发表评论..."}
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
