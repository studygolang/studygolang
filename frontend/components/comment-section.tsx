"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { ChevronUp, MessageSquare, Pencil } from "lucide-react"
import { toast } from "sonner"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { CommentForm } from "@/components/comment-form"
import { CommentEditForm } from "@/components/comment-edit-form"
import { SafeHtml } from "@/components/safe-html"
import { commentAPI, likeAPI, needRelogin } from "@/lib/api"
import { useAuth } from "@/lib/auth-context"
import { formatTime } from "@/lib/utils"
import { OBJTYPE_COMMENT, type Comment } from "@/lib/types"

interface CommentSectionProps {
  objid: number
  objtype: number
  // 服务端已取到评论时传入（首屏直出，利于 SEO）；缺省时客户端拉取
  comments?: Comment[]
  // 被评论对象的属主 uid，用于标记「楼主」
  ownerUid?: number
  // 作者关闭了回复（close_reply）：只读展示，隐藏评论框
  closed?: boolean
}

export function CommentSection({ objid, objtype, comments, ownerUid, closed = false }: CommentSectionProps) {
  const router = useRouter()
  const { user } = useAuth()
  const [list, setList] = useState<Comment[]>(comments ?? [])
  const [replyTo, setReplyTo] = useState<{ floor: number; username: string } | null>(null)
  const [editingCid, setEditingCid] = useState<number | null>(null)
  const [likingCid, setLikingCid] = useState<number | null>(null)
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    if (comments) {
      setList(comments)
      return
    }
    let cancelled = false
    commentAPI
      .getList(objid, objtype)
      .then((data) => {
        if (!cancelled) setList(data?.comments ?? [])
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [objid, objtype, comments, reloadKey])

  function refreshComments() {
    setReplyTo(null)
    setEditingCid(null)
    setReloadKey((k) => k + 1)
    router.refresh()
  }

  async function handleLike(cid: number) {
    if (likingCid !== null) return
    const target = list.find((c) => c.id === cid)
    if (!target) return

    const nextLiked = !target.liked
    setLikingCid(cid)
    try {
      await likeAPI.toggle(cid, OBJTYPE_COMMENT, nextLiked)
      setList((prev) =>
        prev.map((c) =>
          c.id === cid
            ? {
                ...c,
                liked: nextLiked,
                likenum: Math.max(0, (c.likenum ?? 0) + (nextLiked ? 1 : -1)),
              }
            : c
        )
      )
    } catch (err) {
      if (needRelogin(err)) {
        router.push("/account/login")
      } else {
        toast.error(err instanceof Error ? err.message : "操作失败")
      }
    } finally {
      setLikingCid(null)
    }
  }

  return (
    <Card>
      <CardContent className="p-5 sm:p-6">
        <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
          <MessageSquare className="h-4 w-4 text-primary" />
          评论 ({list.length})
        </h2>

        <div className="mt-4">
          {closed ? (
            <p className="rounded-md border border-border bg-muted/40 px-3 py-2 text-xs text-muted-foreground">
              作者关闭了回复功能
            </p>
          ) : (
            <CommentForm
              objid={objid}
              objtype={objtype}
              replyTo={replyTo}
              onSuccess={refreshComments}
            />
          )}
        </div>

        {list.length === 0 ? (
          <p className="mt-6 rounded-lg border border-dashed border-border py-8 text-center text-sm text-muted-foreground">
            暂无评论，来发表第一条评论吧
          </p>
        ) : (
          <div className="mt-6 divide-y divide-border">
            {list.map((comment, index) => {
              const floor = comment.floor || index + 1
              const canEdit =
                !!user && (user.is_root || user.is_admin || user.uid === comment.uid)
              const editing = editingCid === comment.id

              return (
                <div key={comment.id} id={`comment${comment.id}`} className="py-4 first:pt-0">
                  <div className="flex gap-3">
                    <Avatar className="mt-0.5 h-8 w-8 shrink-0">
                      {comment.avatar && <AvatarImage src={comment.avatar} alt={comment.name} />}
                      <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                        {comment.name ? comment.name.charAt(0).toUpperCase() : "?"}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <Link
                          href={`/user/${comment.username || comment.name}`}
                          className="text-sm font-semibold text-foreground hover:text-primary"
                        >
                          {comment.name}
                        </Link>
                        {ownerUid !== undefined && comment.uid === ownerUid && (
                          <Badge
                            variant="secondary"
                            className="bg-primary/10 text-[10px] text-primary"
                          >
                            楼主
                          </Badge>
                        )}
                        <span className="text-xs text-muted-foreground">
                          {formatTime(comment.ctime)}
                        </span>
                        <span className="ml-auto text-xs text-muted-foreground">#{floor}</span>
                      </div>

                      {editing ? (
                        <div className="mt-2">
                          <CommentEditForm
                            cid={comment.id}
                            initialContent={comment.raw_content ?? comment.content}
                            onCancel={() => setEditingCid(null)}
                            onSuccess={refreshComments}
                          />
                        </div>
                      ) : (
                        // 后端已把评论转义并链接化 @user / #N楼，这里不能再走 markdown
                        <SafeHtml
                          html={comment.content}
                          className="prose prose-sm mt-2 max-w-none whitespace-pre-wrap text-foreground"
                        />
                      )}

                      {!editing && (
                        <div className="mt-2 flex items-center gap-3">
                          <button
                            type="button"
                            onClick={() => handleLike(comment.id)}
                            disabled={likingCid === comment.id}
                            className={`flex cursor-pointer items-center gap-1 text-xs transition-colors hover:text-primary disabled:opacity-50 ${
                              comment.liked ? "text-primary" : "text-muted-foreground"
                            }`}
                          >
                            <ChevronUp className="h-3.5 w-3.5" />
                            {comment.likenum ?? 0}
                          </button>
                          {!closed && (
                            <button
                              type="button"
                              onClick={() =>
                                setReplyTo({ floor, username: comment.username || comment.name })
                              }
                              className="cursor-pointer text-xs text-muted-foreground transition-colors hover:text-primary"
                            >
                              回复
                            </button>
                          )}
                          {canEdit && (
                            <button
                              type="button"
                              onClick={() => setEditingCid(comment.id)}
                              className="ml-auto flex cursor-pointer items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-primary"
                            >
                              <Pencil className="h-3 w-3" />
                              编辑
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </CardContent>
    </Card>
  )
}