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

interface Topic {
  id: number
  title: string
  author: string
  authorInitial: string
  tag: string
  tagColor: string
  time: string
  views: number
  comments: number
  likes: number
  pinned?: boolean
  excerpt?: string
}

const topics: Topic[] = [
  {
    id: 1,
    title: "Go 1.26 \u84c4\u52bf\u5f85\u53d1\uff0c\u7b2c\u4e00\u4e2a RC \u7248\u672c\u5df2\u53d1\u5e03",
    author: "polaris",
    authorInitial: "P",
    tag: "Go\u52a8\u6001",
    tagColor: "bg-primary/10 text-primary",
    time: "3 \u5c0f\u65f6\u524d",
    views: 1850,
    comments: 42,
    likes: 128,
    pinned: true,
    excerpt: "Go 1.26 \u5c06\u5f15\u5165\u591a\u9879\u91cd\u8981\u66f4\u65b0\uff0c\u5305\u62ec\u6cdb\u578b\u589e\u5f3a\u3001\u6027\u80fd\u4f18\u5316\u548c\u65b0\u7684\u6807\u51c6\u5e93\u529f\u80fd\u3002\u672c\u6587\u8be6\u7ec6\u89e3\u6790\u4e86 RC1 \u7248\u672c\u7684\u4e3b\u8981\u53d8\u5316...",
  },
  {
    id: 2,
    title: "Golang \u7528\u4e0d\u52302000\u884c\u4ee3\u7801\u5b9e\u73b0\u4e00\u4e2a\u96f6\u62f7\u8d1d\u7684\u6d88\u606f\u961f\u5217",
    author: "chowyu12",
    authorInitial: "C",
    tag: "Go\u5468\u520a",
    tagColor: "bg-accent/10 text-accent",
    time: "5 \u5c0f\u65f6\u524d",
    views: 1490,
    comments: 23,
    likes: 87,
    excerpt: "\u901a\u8fc7\u5de7\u5999\u5229\u7528 mmap \u548c\u5171\u4eab\u5185\u5b58\uff0c\u6211\u4eec\u53ef\u4ee5\u6784\u5efa\u4e00\u4e2a\u9ad8\u6027\u80fd\u7684\u96f6\u62f7\u8d1d\u6d88\u606f\u961f\u5217\u7cfb\u7edf...",
  },
  {
    id: 3,
    title: "\u5b57\u8282\u5927\u6a21\u578b\u6846\u67b6 Eino\uff0c\u57fa\u4e8e Golang \u7684\u5927\u6a21\u578b\u6846\u67b6\u5168\u89e3\u6790",
    author: "guo_hongzhi",
    authorInitial: "G",
    tag: "\u4eba\u5de5\u667a\u80fd",
    tagColor: "bg-chart-3/10 text-chart-3",
    time: "8 \u5c0f\u65f6\u524d",
    views: 965,
    comments: 15,
    likes: 56,
    excerpt: "\u5b57\u8282\u8df3\u52a8\u5f00\u6e90\u7684 Eino \u6846\u67b6\u4e3a Go \u5f00\u53d1\u8005\u63d0\u4f9b\u4e86\u4fbf\u6377\u7684 LLM \u5e94\u7528\u5f00\u53d1\u80fd\u529b...",
  },
  {
    id: 4,
    title: "Go-Zero \u5fae\u670d\u52a1\u6846\u67b6\u5168\u6d41\u7a0b\u5b9e\u6218\u5373\u65f6\u901a\u8baf\u7cfb\u7edf",
    author: "Minnie",
    authorInitial: "M",
    tag: "\u5fae\u670d\u52a1",
    tagColor: "bg-chart-4/10 text-chart-4",
    time: "12 \u5c0f\u65f6\u524d",
    views: 732,
    comments: 8,
    likes: 34,
    excerpt: "\u672c\u6587\u8be6\u7ec6\u4ecb\u7ecd\u4e86\u5982\u4f55\u4f7f\u7528 Go-Zero \u6846\u67b6\u6784\u5efa\u9ad8\u53ef\u7528\u7684\u5373\u65f6\u901a\u8baf\u7cfb\u7edf...",
  },
  {
    id: 5,
    title: "\u4e3a\u4ec0\u4e48\u963f\u91cc\u3001\u5b57\u8282\u3001\u817e\u8baf\u90fd\u5728\u7528 Go\uff1f\u63ed\u79d8\u8fd9\u95e8\u8bed\u8a00\u7684\u4f18\u52bf",
    author: "guo_hongzhi",
    authorInitial: "G",
    tag: "Go\u57fa\u7840",
    tagColor: "bg-primary/10 text-primary",
    time: "1 \u5929\u524d",
    views: 2899,
    comments: 67,
    likes: 215,
    excerpt: "\u4ece\u5e76\u53d1\u6a21\u578b\u5230\u7f16\u8bd1\u901f\u5ea6\uff0c\u4ece\u90e8\u7f72\u4fbf\u5229\u6027\u5230\u751f\u6001\u5708\u5b8c\u5584\u5ea6\uff0c\u89e3\u8bfb\u5927\u5382\u9009\u62e9 Go \u7684\u6838\u5fc3\u539f\u56e0...",
  },
  {
    id: 6,
    title: "LLM \u5f00\u53d1\u5de5\u7a0b\u5e08\u5165\u884c\u5b9e\u6218 -- \u4ece 0 \u5230 1 \u5f00\u53d1\u8f7b\u91cf\u5316\u79c1\u6709\u5927\u6a21\u578b",
    author: "umansyds",
    authorInitial: "U",
    tag: "\u4eba\u5de5\u667a\u80fd",
    tagColor: "bg-chart-3/10 text-chart-3",
    time: "1 \u5929\u524d",
    views: 1109,
    comments: 12,
    likes: 43,
  },
  {
    id: 7,
    title: "\u6df1\u5165\u7406\u89e3 Go \u8bed\u8a00\u4e2d\u7684 sync.Pool\uff1a\u63d0\u5347\u767e\u4e07\u7ea7\u5e76\u53d1\u670d\u52a1\u7a33\u5b9a\u6027",
    author: "guo_hongzhi",
    authorInitial: "G",
    tag: "Go\u6807\u51c6\u5e93",
    tagColor: "bg-accent/10 text-accent",
    time: "2 \u5929\u524d",
    views: 2559,
    comments: 31,
    likes: 98,
  },
  {
    id: 8,
    title: "Web \u5f00\u53d1\u5b9e\u6218\uff1aGin + GORM \u6784\u5efa\u4f01\u4e1a\u7ea7 API \u9879\u76ee\u5168\u89e3\u6790",
    author: "guo_hongzhi",
    authorInitial: "G",
    tag: "Go Web\u5f00\u53d1",
    tagColor: "bg-chart-5/10 text-chart-5",
    time: "2 \u5929\u524d",
    views: 2880,
    comments: 45,
    likes: 156,
  },
]

function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + "k"
  }
  return num.toString()
}

export function TopicList() {
  return (
    <div className="divide-y divide-border">
      {topics.map((topic) => (
        <article
          key={topic.id}
          className="group relative flex gap-4 px-1 py-4 transition-colors first:pt-0 hover:bg-secondary/30 sm:px-3 sm:py-5"
        >
          {/* Avatar */}
          <Avatar className="mt-0.5 hidden h-9 w-9 shrink-0 sm:flex">
            <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
              {topic.authorInitial}
            </AvatarFallback>
          </Avatar>

          {/* Content */}
          <div className="min-w-0 flex-1 space-y-1.5">
            <div className="flex items-start gap-2">
              {topic.pinned && (
                <Pin className="mt-0.5 h-4 w-4 shrink-0 rotate-45 text-primary" />
              )}
              <Link
                href={`/topics/${topic.id}`}
                className="text-[15px] font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
              >
                {topic.title}
              </Link>
            </div>

            {topic.excerpt && (
              <p className="line-clamp-1 text-sm leading-relaxed text-muted-foreground">
                {topic.excerpt}
              </p>
            )}

            <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
              <Badge
                variant="secondary"
                className={`text-xs font-medium ${topic.tagColor}`}
              >
                {topic.tag}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {topic.author}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Clock className="h-3 w-3" />
                {topic.time}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Eye className="h-3 w-3" />
                {formatNumber(topic.views)}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <MessageSquare className="h-3 w-3" />
                {topic.comments}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <ThumbsUp className="h-3 w-3" />
                {topic.likes}
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
