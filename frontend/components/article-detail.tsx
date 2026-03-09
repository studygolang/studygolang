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

const articleData = {
  title: "深入理解 Go 语言 GMP 调度模型",
  author: "polaris",
  authorInitial: "P",
  tag: "Go基础",
  tagColor: "bg-primary/10 text-primary",
  createdAt: "2026-02-25 10:00",
  views: 3210,
  likes: 156,
  readTime: "15 分钟",
}

const contentParagraphs = [
  "GMP 模型是 Go 语言运行时调度的核心。理解它对于编写高性能的 Go 程序至关重要。",
  "## G (Goroutine)",
  "Goroutine 是 Go 中最基本的执行单元。每个 Goroutine 初始只需 2KB 栈空间，远小于线程的 1MB，这使得 Go 可以轻松创建数百万个并发任务。",
  "## M (Machine/Thread)",
  "M 代表操作系统线程。Go 运行时会创建 M 来执行 G，M 的数量由 GOMAXPROCS 和系统资源共同决定。当 M 执行阻塞的系统调用时，运行时会创建新的 M 来保证 P 的利用率。",
  "## P (Processor)",
  "P 是逻辑处理器，它包含运行 Go 代码所需的资源。P 的数量默认等于 CPU 核数，可以通过 runtime.GOMAXPROCS() 设置。每个 P 维护一个本地运行队列（local run queue）。",
  "## 调度流程",
  "1. 当创建一个新的 Goroutine 时，它会被放入当前 P 的本地队列",
  "2. 当本地队列满了，会把一半的 G 转移到全局队列",
  "3. M 从关联的 P 的本地队列获取 G 执行",
  "4. 如果本地队列为空，M 会尝试从全局队列或其他 P 的队列「偷取」G",
  "## Work Stealing 机制",
  "当一个 P 的本地队列为空时，它会尝试从其他 P 的队列中「偷取」一半的 Goroutine。这种机制确保了所有 CPU 核心都能被充分利用，避免了负载不均的问题。",
  "## 总结",
  "GMP 模型是 Go 高并发的基石。通过 Goroutine 的轻量级、M:N 线程模型和 Work Stealing 调度算法，Go 实现了高效的并发调度。理解这些底层原理，有助于我们编写更高效的并发程序。",
]

const relatedArticles = [
  { title: "Go 并发编程最佳实践 2026 版", href: "/articles/2" },
  { title: "深入理解 Go Channel 的实现原理", href: "/articles/3" },
  { title: "Go 内存模型与 Happens-Before 原则", href: "/articles/4" },
]

const comments = [
  {
    id: 1,
    author: "guo_hongzhi",
    initial: "G",
    content: "写得非常清晰！GMP 模型之前一直似懂非懂，看完这篇终于通透了。特别是 Work Stealing 那部分的解释很到位。",
    time: "1 小时前",
    likes: 18,
  },
  {
    id: 2,
    author: "bob_dev",
    initial: "B",
    content: "请问在实际生产环境中，GOMAXPROCS 应该设置成什么值比较好？是直接等于 CPU 核数吗？",
    time: "45 分钟前",
    likes: 5,
  },
]

export function ArticleDetail({ id }: { id: string }) {
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)

  return (
    <div className="space-y-4">
      {/* Article Card */}
      <Card>
        <CardContent className="p-5 sm:p-8">
          {/* Title */}
          <h1 className="text-balance text-2xl font-bold leading-tight text-foreground sm:text-[28px]">
            {articleData.title}
          </h1>

          {/* Meta */}
          <div className="mt-4 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
                  {articleData.authorInitial}
                </AvatarFallback>
              </Avatar>
              <div>
                <Link
                  href={`/user/${articleData.author}`}
                  className="text-sm font-medium text-foreground hover:text-primary"
                >
                  {articleData.author}
                </Link>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Clock className="h-3 w-3" />
                    {articleData.createdAt}
                  </span>
                  <span>{"阅读 " + articleData.readTime}</span>
                </div>
              </div>
            </div>
            <div className="ml-auto hidden items-center gap-2 sm:flex">
              <Badge variant="secondary" className={`text-xs ${articleData.tagColor}`}>
                {articleData.tag}
              </Badge>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Eye className="h-3 w-3" />
                {articleData.views}
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
              if (/^\d+\.\s/.test(para)) {
                return (
                  <div key={i} className="flex gap-2 py-0.5 pl-4 text-sm leading-relaxed text-foreground">
                    <span className="shrink-0 font-semibold text-primary">{para.match(/^\d+/)?.[0]}.</span>
                    <span>{para.replace(/^\d+\.\s/, "")}</span>
                  </div>
                )
              }
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
              {liked ? articleData.likes + 1 : articleData.likes}
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
      <Card>
        <CardContent className="p-5">
          <h3 className="mb-3 text-sm font-semibold text-foreground">{"相关推荐"}</h3>
          <div className="space-y-2">
            {relatedArticles.map((article, i) => (
              <Link
                key={i}
                href={article.href}
                className="flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
              >
                <ExternalLink className="h-3 w-3 shrink-0" />
                {article.title}
              </Link>
            ))}
          </div>
        </CardContent>
      </Card>

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
                      {comment.initial}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">{comment.author}</span>
                      <span className="text-xs text-muted-foreground">{comment.time}</span>
                    </div>
                    <p className="mt-1.5 text-sm leading-relaxed text-foreground">
                      {comment.content}
                    </p>
                    <div className="mt-2 flex items-center gap-3">
                      <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary">
                        <ChevronUp className="h-3.5 w-3.5" />
                        {comment.likes}
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
