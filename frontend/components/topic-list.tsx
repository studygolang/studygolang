"use client"

import Link from "next/link"
import {
  MessageSquare,
  Eye,
  ThumbsUp,
  Pin,
  Clock,
  ArrowUpRight,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import type { Topic } from "@/lib/types"

interface TopicListProps {
  topics?: Topic[]
}

function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + "k"
  }
  return num.toString()
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

const fallbackTopics: Topic[] = []

export function TopicList({ topics = fallbackTopics }: TopicListProps) {
  if (topics.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无话题
      </div>
    )
  }

  return (
    <div className="divide-y divide-border">
      {topics.map((topic) => (
        <article
          key={topic.tid}
          className="group relative flex gap-4 px-1 py-4 transition-colors first:pt-0 hover:bg-secondary/30 sm:px-3 sm:py-5"
        >
          {/* Avatar */}
          <Avatar className="mt-0.5 hidden h-9 w-9 shrink-0 sm:flex">
            <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
              {(topic.user?.username || topic.name) ? (topic.user?.username || topic.name).charAt(0).toUpperCase() : "?"}
            </AvatarFallback>
          </Avatar>

          {/* Content */}
          <div className="min-w-0 flex-1 space-y-1.5">
            <div className="flex items-start gap-2">
              {topic.top > 0 && (
                <Pin className="mt-0.5 h-4 w-4 shrink-0 rotate-45 text-primary" />
              )}
              <Link
                href={`/topics/${topic.tid}`}
                className="text-[15px] font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
              >
                {topic.title}
              </Link>
            </div>

            <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
              {topic.node && (
                <Badge
                  variant="secondary"
                  className="bg-primary/10 text-xs font-medium text-primary"
                >
                  {topic.node.name}
                </Badge>
              )}
              <span className="text-xs text-muted-foreground">
                {topic.user?.username || topic.name || "匿名"}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Clock className="h-3 w-3" />
                {formatTime(topic.ctime)}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Eye className="h-3 w-3" />
                {formatNumber(topic.view || topic.viewnum || 0)}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <MessageSquare className="h-3 w-3" />
                {topic.reply || topic.replynum || 0}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <ThumbsUp className="h-3 w-3" />
                {topic.like || topic.likenum || 0}
              </span>
            </div>
          </div>

          {/* External link indicator */}
          <ArrowUpRight className="mt-1 hidden h-4 w-4 shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 sm:block" />
        </article>
      ))}
    </div>
  )
}
