"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { Pencil, Share2 } from "lucide-react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { LikeButton } from "@/components/like-button"
import { FavoriteButton } from "@/components/favorite-button"
import { TopicSetTopButton } from "@/components/topic-set-top-button"
import { useAuth } from "@/lib/auth-context"

interface ContentActionsProps {
  objid: number
  objtype: number
  likeCount?: number
  // 服务端已判定的点赞/收藏状态（详情页需 force-dynamic 透传 Cookie）
  initialLiked?: boolean
  initialFavorited?: boolean
  ownerUid?: number
  // 文章/项目的属主判定用用户名（后端 logic.CanEdit 比对的是 Username）
  ownerUsername?: string
  editHref?: string
  // 仅话题支持置顶；后端只提供置顶接口，取消由定时任务按 top_time 自动处理
  showTop?: boolean
  isTop?: boolean
  className?: string
}

export function ContentActions({
  objid,
  objtype,
  likeCount = 0,
  initialLiked = false,
  initialFavorited = false,
  ownerUid,
  ownerUsername,
  editHref,
  showTop = false,
  isTop = false,
  className,
}: ContentActionsProps) {
  const { user } = useAuth()
  const router = useRouter()

  const isOwner =
    !!user &&
    ((ownerUid !== undefined && user.uid === ownerUid) ||
      (!!ownerUsername && user.username === ownerUsername))
  const isStaff = !!user && (user.is_root || user.is_admin)
  const canEdit = !!editHref && (isOwner || isStaff)
  const canTop = showTop && !isTop && (isOwner || isStaff)

  async function handleShare() {
    try {
      await navigator.clipboard.writeText(window.location.href)
      toast.success("链接已复制到剪贴板")
    } catch {
      toast.error("复制失败，请手动复制地址栏链接")
    }
  }

  return (
    <div className={`flex flex-wrap items-center gap-2 ${className ?? ""}`}>
      <LikeButton
        objid={objid}
        objtype={objtype}
        initialLiked={initialLiked}
        initialCount={likeCount}
      />
      <FavoriteButton
        objid={objid}
        objtype={objtype}
        initialFavorited={initialFavorited}
      />
      <Button
        variant="ghost"
        size="sm"
        onClick={handleShare}
        className="gap-1 text-muted-foreground"
      >
        <Share2 className="h-4 w-4" />
        分享
      </Button>
      {canEdit && editHref && (
        <Button asChild variant="ghost" size="sm" className="gap-1 text-muted-foreground">
          <Link href={editHref}>
            <Pencil className="h-4 w-4" />
            编辑
          </Link>
        </Button>
      )}
      {canTop && (
        <TopicSetTopButton tid={objid} isTop={isTop} onSuccess={() => router.refresh()} />
      )}
    </div>
  )
}