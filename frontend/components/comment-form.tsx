"use client"

import React, { useRef, useState } from "react"
import { useRouter } from "next/navigation"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { AtUserAutocomplete } from "@/components/at-user-autocomplete"
import { commentAPI, needRelogin } from "@/lib/api"
import { extractMentionedUsernames } from "@/lib/utils"
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
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const router = useRouter()

  const replyTarget = replyTo && replyTo.floor > 0 ? replyTo : null

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (submitting) return

    // 回复某楼层时自动加上 `#N楼 @username ` 前缀（与后端 decodeCmtContentForShow 正则一致）
    const finalContent = replyTarget
      ? `#${replyTarget.floor}楼 @${replyTarget.username} ${content}`
      : content

    if (!finalContent.trim()) return

    setSubmitting(true)

    try {
      await commentAPI.create(objid, objtype, finalContent, extractMentionedUsernames(finalContent))
      toast.success("评论发表成功")
      setContent("")
      onSuccess()
    } catch (err) {
      if (needRelogin(err)) {
        router.push("/account/login")
      } else {
        toast.error(err instanceof Error ? err.message : "评论发表失败")
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      {replyTarget && (
        <div className="rounded-md border border-border bg-muted/40 px-3 py-1.5 text-xs text-muted-foreground">
          回复 #{replyTarget.floor}楼 @{replyTarget.username}
        </div>
      )}
      <AtUserAutocomplete
        textareaRef={textareaRef}
        value={content}
        onChange={setContent}
        placeholder={
          replyTarget ? `回复 @${replyTarget.username}...` : "发表评论，输入 @ 可提及用户"
        }
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