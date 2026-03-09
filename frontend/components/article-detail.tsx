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
  ExternalLink,
  ChevronUp,
  MoreHorizontal,
} from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import type { Article, Comment } from "@/lib/types"

interface ArticleDetailProps {
  id: string
  article?: Article
  prev?: Article
  next?: Article
  comments?: Comment[]
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
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)

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

  const contentParagraphs = article.content ? article.content.split("\n") : []

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
                <Link
                  href={`/user/${article.author}`}
                  className="text-sm font-medium text-foreground hover:text-primary"
                >
                  {article.author}
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
          <div className="space-y-4">
            {contentParagraphs.map((para, i) => {
              if (para.startsWith("## ")) {
                return (
                  <h2
                    key={i}
                    className="mt-8 mb-3 text-xl font-bold text-foreground first:mt-0"
                  >
                    {para.replace("## ", "")}
                  </h2>
                )
              }
              if (para.startsWith("### ")) {
                return (
                  <h3 key={i} className="mb-2 mt-4 text-base font-semibold text-foreground">
                    {para.replace("### ", "")}
                  </h3>
                )
              }
              if (/^\d+\.\s/.test(para)) {
                return (
                  <div key={i} className="flex gap-2 py-0.5 pl-4 text-sm leading-relaxed text-foreground">
                    <span className="shrink-0 font-semibold text-primary">{para.match(/^\d+/)?.[0]}.</span>
                    <span>{para.replace(/^\d+\.\s/, "")}</span>
                  </div>
                )
              }
              if (para.trim() === "") return <div key={i} className="h-2" />
              return (
                <p key={i} className="text-[15px] leading-relaxed text-foreground">
                  {para}
                </p>
              )
            })}
          </div>

          {/* Actions */}
          <Separator className="my-6" />
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant={liked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                liked ? "bg-primary text-primary-foreground" : "text-muted-foreground"
              }`}
              onClick={() => setLiked(!liked)}
            >
              <ThumbsUp className="h-3.5 w-3.5" />
              {liked ? article.likenum + 1 : article.likenum}
            </Button>
            <Button
              variant={bookmarked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                bookmarked ? "bg-primary text-primary-foreground" : "text-muted-foreground"
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

          <div className="mt-4 rounded-lg border border-border p-3">
            <textarea
              placeholder="写下你的评论..."
              className="min-h-[80px] w-full resize-none bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
            />
            <div className="mt-2 flex items-center justify-between">
              <p className="text-xs text-muted-foreground">{"支持 Markdown 语法"}</p>
              <Button size="sm" className="bg-primary text-primary-foreground hover:bg-primary/90">
                {"发表评论"}
              </Button>
            </div>
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
                    <p className="mt-1.5 text-sm leading-relaxed text-foreground">
                      {comment.content}
                    </p>
                    <div className="mt-2 flex items-center gap-3">
                      <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary">
                        <ChevronUp className="h-3.5 w-3.5" />
                        0
                      </button>
                      <button className="text-xs text-muted-foreground hover:text-primary">{"回复"}</button>
                      <button className="ml-auto text-muted-foreground hover:text-primary">
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
