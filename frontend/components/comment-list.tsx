import Link from "next/link"
import { MessageSquare, Clock } from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { formatNum } from "@/lib/utils"
import type { UserComment } from "@/lib/types"

interface CommentListProps {
  comments?: UserComment[]
}

function formatTime(ctime: string): string {
  try {
    const date = new Date(ctime)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    if (minutes < 60) return `${minutes} 分钟前`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours} 小时前`
    const days = Math.floor(hours / 24)
    if (days < 30) return `${days} 天前`
    return ctime.slice(0, 10)
  } catch {
    return ctime
  }
}

function getObjTypeName(objtype: number): string {
  switch (objtype) {
    case 1: return "话题"
    case 2: return "文章"
    case 3: return "资源"
    case 4: return "项目"
    case 5: return "书籍"
    default: return "内容"
  }
}

function getObjLink(comment: UserComment): string {
  if (!comment.objinfo?.uri) {
    switch (comment.objtype) {
      case 1: return `/topics/${comment.objid}`
      case 2: return `/articles/${comment.objid}`
      case 3: return `/resources/${comment.objid}`
      case 4: return `/projects/${comment.objid}`
      case 5: return `/books/${comment.objid}`
      default: return "#"
    }
  }
  return comment.objinfo.uri
}

export function CommentList({ comments = [] }: CommentListProps) {
  if (comments.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无评论
      </div>
    )
  }

  return (
    <div className="divide-y divide-border">
      {comments.map((comment) => (
        <div
          key={comment.cid}
          className="py-4 first:pt-0"
        >
          {/* 评论对象链接 */}
          <div className="mb-2 flex items-center gap-2">
            <MessageSquare className="h-4 w-4 text-muted-foreground" />
            <span className="text-xs text-muted-foreground">
              回复了
            </span>
            <Link
              href={getObjLink(comment)}
              className="text-sm font-medium text-primary transition-colors hover:text-primary/80"
            >
              {comment.objinfo?.title || `${getObjTypeName(comment.objtype)} #${comment.objid}`}
            </Link>
          </div>

          {/* 评论内容 */}
          <div className="rounded-lg bg-secondary/30 p-3">
            <p className="text-sm leading-relaxed text-foreground whitespace-pre-wrap">
              {comment.content}
            </p>
          </div>

          {/* 评论元信息 */}
          <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
            <span className="flex items-center gap-1">
              <Clock className="h-3 w-3" />
              {formatTime(comment.ctime)}
            </span>
            {comment.floor > 0 && (
              <span>第 {comment.floor} 楼</span>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}
