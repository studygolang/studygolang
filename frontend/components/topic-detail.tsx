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
  Flag,
  ChevronUp,
  MoreHorizontal,
} from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { TopicAppendForm } from "@/components/topic-append-form"
import { CommentForm } from "@/components/comment-form"
import { MarkdownContent } from "@/components/markdown-content"
import { useAuth } from "@/lib/auth-context"
import { likeAPI, favoriteAPI } from "@/lib/api"
import type { Topic, TopicReply, TopicAppend } from "@/lib/types"

interface TopicDetailProps {
  id: string
  topic?: Topic
  replies?: TopicReply[]
  appends?: TopicAppend[]
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

export function TopicDetail({ id, topic, replies = [], appends = [] }: TopicDetailProps) {
  const router = useRouter()
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)
  const [likeCount, setLikeCount] = useState(topic?.likenum || 0)

  // 加载用户的点赞和收藏状态
  useEffect(() => {
    const loadUserStatus = async () => {
      try {
        const likeRes = await fetch(`${API_BASE}/api/v1/likes/${id}/status?objtype=0`, {
          credentials: "include",
        })
        if (likeRes.ok) {
          const likeData = await likeRes.json()
          if (likeData.code === 0 && likeData.data?.has_like) {
            setLiked(true)
          }
        }

        // 检查收藏状态
        const favRes = await fetch(`${API_BASE}/api/v1/favorites/${id}/status?objtype=0`, {
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
      const res = await fetch(`${API_BASE}/api/v1/likes/${id}?objtype=0&flag=${liked ? 0 : 1}`, {
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
      const res = await fetch(`${API_BASE}/api/v1/favorites/${id}?objtype=0&collect=${bookmarked ? 0 : 1}`, {
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

  if (!topic) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-muted-foreground">
          话题不存在或已被删除
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      {/* Topic Card */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          {/* Title */}
          <h1 className="text-balance text-xl font-bold leading-snug text-foreground sm:text-2xl">
            {topic.title}
          </h1>

          {/* Meta */}
          <div className="mt-3 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-7 w-7">
                <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                  {topic.user?.username ? topic.user.username.charAt(0).toUpperCase() : (topic.name ? topic.name.charAt(0).toUpperCase() : "?")}
                </AvatarFallback>
              </Avatar>
              <Link
                href={`/user/${topic.user?.username || topic.name}`}
                className="text-sm font-medium text-foreground hover:text-primary"
              >
                {topic.user?.username || topic.name}
              </Link>
            </div>
            {topic.node && (
              <>
                <Separator orientation="vertical" className="h-4" />
                <Badge variant="secondary" className="bg-primary/10 text-xs text-primary">
                  {topic.node.name}
                </Badge>
              </>
            )}
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Clock className="h-3 w-3" />
              {formatTime(topic.ctime)}
            </span>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Eye className="h-3 w-3" />
              {topic.view || topic.viewnum || 0}
            </span>
          </div>

          {/* Content */}
          <MarkdownContent content={topic.content ?? ""} className="mt-6" />

          {/* Actions */}
          <Separator className="my-5" />
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant={liked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                liked
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground"
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
                bookmarked
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground"
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
            <Button variant="ghost" size="sm" className="ml-auto gap-1.5 text-muted-foreground">
              <Flag className="h-3.5 w-3.5" />
              {"举报"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* 追加内容展示 */}
      {appends.length > 0 && (
        <Card>
          <CardContent className="p-5 sm:p-6">
            <h2 className="text-base font-semibold text-foreground">附言</h2>
            <div className="mt-4 space-y-4 divide-y divide-border">
              {appends.map((append, index) => (
                <div key={append.id} className={index > 0 ? "pt-4" : ""}>
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <span className="font-medium">附言 {index + 1}</span>
                    <span>{formatTime(append.created_at)}</span>
                  </div>
                  <MarkdownContent content={append.content} className="mt-2" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 追加表单（仅作者可见，客户端判断） */}
      <TopicAppendSection tid={parseInt(id)} topicUid={topic?.uid} />

      {/* Comments Section */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <MessageSquare className="h-4 w-4 text-primary" />
            {"评论"} ({replies.length})
          </h2>

          {/* Comment Input */}
          <div className="mt-4">
            <CommentForm
              objid={parseInt(id)}
              objtype={0}
              onSuccess={() => router.refresh()}
            />
          </div>

          {/* Comments List */}
          <div className="mt-6 space-y-0 divide-y divide-border">
            {replies.map((reply, index) => (
              <div key={reply.id} className="py-4 first:pt-0">
                <div className="flex gap-3">
                  <Avatar className="mt-0.5 h-8 w-8 shrink-0">
                    <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                      {reply.name ? reply.name.charAt(0).toUpperCase() : "?"}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <Link
                        href={`/user/${reply.name}`}
                        className="text-sm font-semibold text-foreground hover:text-primary"
                      >
                        {reply.name}
                      </Link>
                      {reply.uid === topic.uid && (
                        <Badge variant="secondary" className="bg-primary/10 text-[10px] text-primary">
                          {"楼主"}
                        </Badge>
                      )}
                      <span className="text-xs text-muted-foreground">
                        {formatTime(reply.ctime)}
                      </span>
                      <span className="ml-auto text-xs text-muted-foreground">
                        #{index + 1}
                      </span>
                    </div>
                    <MarkdownContent content={reply.content} className="mt-2" />
                    <div className="mt-2 flex items-center gap-3">
                      <button className="flex cursor-pointer items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-primary">
                        <ChevronUp className="h-3.5 w-3.5" />
                        0
                      </button>
                      <button className="cursor-pointer text-xs text-muted-foreground transition-colors hover:text-primary">
                        {"回复"}
                      </button>
                      <button className="ml-auto cursor-pointer text-muted-foreground transition-colors hover:text-primary">
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

// 追加表单区域：客户端判断当前用户是否为作者
function TopicAppendSection({ tid, topicUid }: { tid: number; topicUid?: number }) {
  const { user } = useAuth()

  if (topicUid === undefined || !user || user.uid !== topicUid) return null
  return <TopicAppendForm tid={tid} />
}
