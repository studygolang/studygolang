"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import {
  ThumbsUp,
  MessageSquare,
  Bookmark,
  Share2,
  Eye,
  Clock,
  ExternalLink,
  ChevronUp,
  MoreHorizontal,
} from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { MarkdownContent } from "@/components/markdown-content"
import { CommentForm } from "@/components/comment-form"
import type { Article, ArticleComment } from "@/lib/types"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

interface ArticleDetailProps {
  id: string
  article?: Article
  prev?: Article
  next?: Article
  // 后端文章详情评论字段名为 replies，类型为 ArticleComment
  comments?: ArticleComment[]
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

export function ArticleDetail({ id, article, prev, next, comments = [] }: ArticleDetailProps) {
  const router = useRouter()
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)
  const [likeCount, setLikeCount] = useState(article?.likenum || 0)

  // 加载用户的点赞和收藏状态
  useEffect(() => {
    const loadUserStatus = async () => {
      try {
        // 检查点赞状态
        const likeRes = await fetch(`${API_BASE}/api/v1/likes/${id}/status?objtype=1`, {
          credentials: "include",
        })
        if (likeRes.ok) {
          const likeData = await likeRes.json()
          if (likeData.code === 0 && likeData.data?.has_like) {
            setLiked(true)
          }
        }

        // 检查收藏状态
        const favRes = await fetch(`${API_BASE}/api/v1/favorites/${id}/status?objtype=1`, {
          credentials: "include",
        })
        if (favRes.ok) {
          const favData = await favRes.json()
          if (favData.code === 0 && favData.data?.has_favorite) {
            setBookmarked(true)
          }
        }
      } catch {
        // 未登录或网络错误，忽略
      }
    }
    loadUserStatus()
  }, [id])

  const handleLike = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/v1/likes/${id}?objtype=1&flag=${liked ? 0 : 1}`, {
        method: "POST",
        credentials: "include",
      })
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0) {
          setLiked(!liked)
          setLikeCount(liked ? likeCount - 1 : likeCount + 1)
        }
      }
    } catch {
      // 网络错误
    }
  }

  const handleBookmark = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/v1/favorites/${id}?objtype=1&collect=${bookmarked ? 0 : 1}`, {
        method: "POST",
        credentials: "include",
      })
      if (res.ok) {
        const data = await res.json()
        if (data.code === 0) {
          setBookmarked(!bookmarked)
        }
      }
    } catch {
      // 网络错误
    }
  }

  if (!article) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-muted-foreground">
          文章不存在或已被删除
        </CardContent>
      </Card>
    )
  }

  const relatedArticles = [
    ...(prev ? [{ title: prev.title, href: `/articles/${prev.id}` }] : []),
    ...(next ? [{ title: next.title, href: `/articles/${next.id}` }] : []),
  ]

  return (
    <div className="space-y-4">
      {/* Article Card */}
      <Card>
        <CardContent className="p-5 sm:p-8">
          {/* Title */}
          <h1 className="text-balance text-2xl font-bold leading-tight text-foreground sm:text-[28px]">
            {article.title}
          </h1>

          {/* Meta */}
          <div className="mt-4 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
                  {article.author ? article.author.charAt(0).toUpperCase() : "?"}
                </AvatarFallback>
              </Avatar>
              <div>
                {/* author 为原文作者名，站内用 author_txt 对应 username，API 暂未返回 author_uid，用 author 作路由参数 */}
                <Link
                  href={`/user/${article.author || ""}`}
                  className="text-sm font-medium text-foreground hover:text-primary"
                >
                  {article.author || "未知作者"}
                </Link>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Clock className="h-3 w-3" />
                    {formatTime(article.ctime)}
                  </span>
                </div>
              </div>
            </div>
            <div className="ml-auto hidden items-center gap-2 sm:flex">
              {article.tags && (
                <Badge variant="secondary" className="bg-primary/10 text-xs text-primary">
                  {article.tags.split(",")[0]}
                </Badge>
              )}
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Eye className="h-3 w-3" />
                {article.viewnum}
              </span>
            </div>
          </div>

          <Separator className="my-6" />

          {/* Content */}
          <MarkdownContent content={article.content ?? ""} />

          {/* Actions */}
          <Separator className="my-6" />
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant={liked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                liked ? "bg-primary text-primary-foreground" : "text-muted-foreground"
              }`}
              onClick={handleLike}
            >
              <ThumbsUp className="h-3.5 w-3.5" />
              {likeCount}
            </Button>
            <Button
              variant={bookmarked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                bookmarked ? "bg-primary text-primary-foreground" : "text-muted-foreground"
              }`}
              onClick={handleBookmark}
            >
              <Bookmark className="h-3.5 w-3.5" />
              {bookmarked ? "已收藏" : "收藏"}
            </Button>
            <Button variant="outline" size="sm" className="gap-1.5 text-muted-foreground">
              <Share2 className="h-3.5 w-3.5" />
              {"分享"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Related Articles */}
      {relatedArticles.length > 0 && (
        <Card>
          <CardContent className="p-5">
            <h3 className="mb-3 text-sm font-semibold text-foreground">{"相关推荐"}</h3>
            <div className="space-y-2">
              {relatedArticles.map((a, i) => (
                <Link
                  key={i}
                  href={a.href}
                  className="flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
                >
                  <ExternalLink className="h-3 w-3 shrink-0" />
                  {a.title}
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Comments */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <MessageSquare className="h-4 w-4 text-primary" />
            {"评论"} ({comments.length})
          </h2>

          <div className="mt-4">
            <CommentForm
              objid={parseInt(id)}
              objtype={1}
              onSuccess={() => router.refresh()}
            />
          </div>

          <div className="mt-6 space-y-0 divide-y divide-border">
            {comments.map((comment) => (
              <div key={comment.id} className="py-4 first:pt-0">
                <div className="flex gap-3">
                  <Avatar className="mt-0.5 h-8 w-8 shrink-0">
                    <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                      {comment.name ? comment.name.charAt(0).toUpperCase() : "?"}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">{comment.name}</span>
                      <span className="text-xs text-muted-foreground">{formatTime(comment.ctime)}</span>
                    </div>
                    <MarkdownContent content={comment.content} className="mt-1.5" />
                    <div className="mt-2 flex items-center gap-3">
                      <button className="flex cursor-pointer items-center gap-1 text-xs text-muted-foreground hover:text-primary">
                        <ChevronUp className="h-3.5 w-3.5" />
                        0
                      </button>
                      <button className="cursor-pointer text-xs text-muted-foreground hover:text-primary">{"回复"}</button>
                      <button className="ml-auto cursor-pointer text-muted-foreground hover:text-primary">
                        <MoreHorizontal className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
