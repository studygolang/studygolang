"use client"

import { useState } from "react"
import Link from "next/link"
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
import type { Topic, TopicReply } from "@/lib/types"

interface TopicDetailProps {
  id: string
  topic?: Topic
  replies?: TopicReply[]
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

export function TopicDetail({ id, topic, replies = [] }: TopicDetailProps) {
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)

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
                  {topic.name ? topic.name.charAt(0).toUpperCase() : "?"}
                </AvatarFallback>
              </Avatar>
              <Link
                href={`/user/${topic.name}`}
                className="text-sm font-medium text-foreground hover:text-primary"
              >
                {topic.name}
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
              {topic.viewnum}
            </span>
          </div>

          {/* Content */}
          <div className="prose prose-sm mt-6 max-w-none text-foreground">
            {topic.content.split("\n").map((line, i) => {
              if (line.startsWith("## ")) {
                return (
                  <h2 key={i} className="mb-3 mt-6 text-lg font-bold text-foreground">
                    {line.replace("## ", "")}
                  </h2>
                )
              }
              if (line.startsWith("### ")) {
                return (
                  <h3 key={i} className="mb-2 mt-4 text-base font-semibold text-foreground">
                    {line.replace("### ", "")}
                  </h3>
                )
              }
              if (line.startsWith("```")) {
                return null
              }
              if (line.startsWith("- ")) {
                return (
                  <div key={i} className="flex gap-2 py-0.5 text-sm leading-relaxed text-foreground">
                    <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                    <span>{line.replace("- ", "")}</span>
                  </div>
                )
              }
              if (line.trim() === "") return <div key={i} className="h-2" />
              return (
                <p key={i} className="text-sm leading-relaxed text-foreground">
                  {line}
                </p>
              )
            })}
          </div>

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
              onClick={() => setLiked(!liked)}
            >
              <ThumbsUp className="h-3.5 w-3.5" />
              {liked ? topic.likenum + 1 : topic.likenum}
            </Button>
            <Button
              variant={bookmarked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                bookmarked
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground"
              }`}
              onClick={() => setBookmarked(!bookmarked)}
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

      {/* Comments Section */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <MessageSquare className="h-4 w-4 text-primary" />
            {"评论"} ({replies.length})
          </h2>

          {/* Comment Input */}
          <div className="mt-4 rounded-lg border border-border p-3">
            <textarea
              placeholder="写下你的评论..."
              className="min-h-[80px] w-full resize-none bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
            />
            <div className="mt-2 flex items-center justify-between">
              <p className="text-xs text-muted-foreground">
                {"支持 Markdown 语法"}
              </p>
              <Button size="sm" className="bg-primary text-primary-foreground hover:bg-primary/90">
                {"发表评论"}
              </Button>
            </div>
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
                    <div className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                      {reply.content}
                    </div>
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
