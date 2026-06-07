"use client"

import Link from "next/link"
import {
  MessageSquare,
  ThumbsUp,
  Pin,
  Clock,
  ArrowUpRight,
  FileText,
  FolderGit2,
  Link2,
  BookOpen,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import type { Feed } from "@/lib/types"

interface FeedListProps {
  feeds?: Feed[]
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

// 根据 objtype 获取详情页链接
function getDetailUrl(feed: Feed): string {
  if (feed.Uri) return feed.Uri
  switch (feed.Objtype) {
    case 4: return `/p/${feed.Objid}`
    case 5: return `/book/${feed.Objid}`
    default: return `/${getTypePath(feed.Objtype)}/${feed.Objid}`
  }
}

// 根据 objtype 获取路径前缀
function getTypePath(objtype: number): string {
  switch (objtype) {
    case 0: return "topics"
    case 1: return "articles"
    case 2: return "resources"
    case 4: return "projects"
    case 5: return "books"
    default: return "topics"
  }
}

// 根据 objtype 获取类型标签
function getTypeBadge(objtype: number): { label: string; icon: React.ReactNode; className: string } {
  switch (objtype) {
    case 0:
      return { label: "话题", icon: <MessageSquare className="h-3 w-3" />, className: "bg-blue-500/10 text-blue-600" }
    case 1:
      return { label: "文章", icon: <FileText className="h-3 w-3" />, className: "bg-green-500/10 text-green-600" }
    case 2:
      return { label: "资源", icon: <Link2 className="h-3 w-3" />, className: "bg-orange-500/10 text-orange-600" }
    case 4:
      return { label: "项目", icon: <FolderGit2 className="h-3 w-3" />, className: "bg-purple-500/10 text-purple-600" }
    case 5:
      return { label: "书籍", icon: <BookOpen className="h-3 w-3" />, className: "bg-pink-500/10 text-pink-600" }
    default:
      return { label: "动态", icon: <MessageSquare className="h-3 w-3" />, className: "bg-gray-500/10 text-gray-600" }
  }
}

const fallbackFeeds: Feed[] = []

export function FeedList({ feeds = fallbackFeeds }: FeedListProps) {
  if (feeds.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无动态
      </div>
    )
  }

  return (
    <div className="divide-y divide-border">
      {feeds.map((feed) => {
        const typeBadge = getTypeBadge(feed.Objtype)
        const detailUrl = getDetailUrl(feed)
        const nodeName = feed.Node && 'name' in feed.Node ? feed.Node.name : undefined

        return (
          <article
            key={`feed-${feed.Objtype}-${feed.Objid}`}
            className="group relative flex gap-4 px-1 py-4 transition-colors first:pt-0 hover:bg-secondary/30 sm:px-3 sm:py-5"
          >
            {/* Avatar */}
            <Avatar className="mt-0.5 hidden h-9 w-9 shrink-0 sm:flex">
              <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
                {(feed.User?.username || feed.Author) ? (feed.User?.username || feed.Author).charAt(0).toUpperCase() : "?"}
              </AvatarFallback>
            </Avatar>

            {/* Content */}
            <div className="min-w-0 flex-1 space-y-1.5">
              <div className="flex items-start gap-2">
                {feed.Top > 0 && (
                  <Pin className="mt-0.5 h-4 w-4 shrink-0 rotate-45 text-primary" />
                )}
                <Link
                  href={detailUrl}
                  className="text-[15px] font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                >
                  {feed.Title}
                </Link>
              </div>

              <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                {/* 类型标签 */}
                <Badge
                  variant="secondary"
                  className={`gap-1 text-xs font-medium ${typeBadge.className}`}
                >
                  {typeBadge.icon}
                  {typeBadge.label}
                </Badge>
                {/* 节点标签（仅话题显示） */}
                {nodeName && feed.Objtype === 0 && (
                  <Badge
                    variant="secondary"
                    className="bg-primary/10 text-xs font-medium text-primary"
                  >
                    {nodeName}
                  </Badge>
                )}
                <span className="text-xs text-muted-foreground">
                  {feed.User?.username || feed.Author || "匿名"}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Clock className="h-3 w-3" />
                  {formatTime(feed.CreatedAt)}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <ThumbsUp className="h-3 w-3" />
                  {feed.Likenum || 0}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <MessageSquare className="h-3 w-3" />
                  {feed.Cmtnum || 0}
                </span>
              </div>
            </div>

            {/* External link indicator */}
            <ArrowUpRight className="mt-1 hidden h-4 w-4 shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 sm:block" />
          </article>
        )
      })}
    </div>
  )
}
