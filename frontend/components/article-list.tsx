import Link from "next/link"
import {
  Eye,
  MessageSquare,
  ThumbsUp,
  Clock,
  ExternalLink,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"

interface Article {
  id: number
  title: string
  author: string
  authorInitial: string
  source?: string
  sourceUrl?: string
  tag: string
  tagColor: string
  time: string
  views: number
  comments: number
  likes: number
  excerpt: string
  cover?: boolean
}

const articles: Article[] = [
  {
    id: 1,
    title: "深入理解 Go 语言 GMP 调度模型",
    author: "polaris",
    authorInitial: "P",
    tag: "Go基础",
    tagColor: "bg-primary/10 text-primary",
    time: "2 小时前",
    views: 3210,
    comments: 28,
    likes: 156,
    excerpt:
      "GMP 模型是 Go 语言运行时的核心调度机制。本文从源码层面深入分析 Goroutine、Machine 和 Processor 三者的协作关系，帮助你理解 Go 高并发的底层原理。",
    cover: true,
  },
  {
    id: 2,
    title: "使用 Go 构建高性能 WebSocket 服务器",
    author: "alice_go",
    authorInitial: "A",
    tag: "Go Web开发",
    tagColor: "bg-chart-5/10 text-chart-5",
    time: "5 小时前",
    views: 1890,
    comments: 15,
    likes: 89,
    excerpt:
      "WebSocket 在实时应用中扮演着重要角色。本文将带你从零构建一个支持百万连接的 WebSocket 服务器，涵盖连接管理、消息广播和心跳检测等核心功能。",
  },
  {
    id: 3,
    title: "Go 错误处理最佳实践：从 errors 到 errgroup",
    author: "bob_dev",
    authorInitial: "B",
    source: "个人博客",
    sourceUrl: "#",
    tag: "Go标准库",
    tagColor: "bg-accent/10 text-accent",
    time: "8 小时前",
    views: 2450,
    comments: 32,
    likes: 112,
    excerpt:
      "Go 的错误处理一直是社区讨论的热点。本文总结了从基础的 error 接口到 errors.Is/As、fmt.Errorf %w、以及 errgroup 并发错误处理的完整实践指南。",
    cover: true,
  },
  {
    id: 4,
    title: "Kubernetes Operator 开发入门：使用 Go 扩展 K8s",
    author: "charlie",
    authorInitial: "C",
    tag: "Kubernetes",
    tagColor: "bg-chart-4/10 text-chart-4",
    time: "1 天前",
    views: 1560,
    comments: 11,
    likes: 67,
    excerpt:
      "Operator 是 Kubernetes 生态的核心扩展模式。本文通过一个完整的实战案例，手把手教你如何用 Go 开发自定义 Operator 来管理有状态应用。",
  },
  {
    id: 5,
    title: "Go 语言性能调优实战：pprof 从入门到精通",
    author: "guo_hongzhi",
    authorInitial: "G",
    tag: "Go实战",
    tagColor: "bg-primary/10 text-primary",
    time: "1 天前",
    views: 3890,
    comments: 45,
    likes: 198,
    excerpt:
      "性能优化是每个 Go 开发者的必修课。本文详细介绍如何使用 pprof 工具进行 CPU、内存、Goroutine 泄漏等问题的定位和优化。",
    cover: true,
  },
  {
    id: 6,
    title: "从零实现一个 Go 语言的依赖注入框架",
    author: "Minnie",
    authorInitial: "M",
    tag: "Go实战",
    tagColor: "bg-primary/10 text-primary",
    time: "2 天前",
    views: 1230,
    comments: 9,
    likes: 54,
    excerpt:
      "依赖注入是构建可测试、可维护代码的关键模式。本文将展示如何利用 Go 的反射机制实现一个轻量级的 DI 容器。",
  },
]

function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + "k"
  }
  return num.toString()
}

export function ArticleList() {
  return (
    <div className="space-y-4">
      {articles.map((article) => (
        <article
          key={article.id}
          className="group rounded-lg border border-border bg-card p-4 transition-all hover:border-primary/20 hover:shadow-sm sm:p-5"
        >
          <div className="flex gap-4">
            <div className="min-w-0 flex-1 space-y-2">
              {/* Title */}
              <div className="flex items-start gap-2">
                <Link
                  href={`/articles/${article.id}`}
                  className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary sm:text-[17px]"
                >
                  {article.title}
                </Link>
                {article.source && (
                  <a
                    href={article.sourceUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="mt-0.5 shrink-0"
                  >
                    <Badge variant="outline" className="gap-1 text-[10px] text-muted-foreground">
                      {article.source}
                      <ExternalLink className="h-2.5 w-2.5" />
                    </Badge>
                  </a>
                )}
              </div>

              {/* Excerpt */}
              <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                {article.excerpt}
              </p>

              {/* Meta */}
              <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                <div className="flex items-center gap-1.5">
                  <Avatar className="h-5 w-5">
                    <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                      {article.authorInitial}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-xs font-medium text-foreground">
                    {article.author}
                  </span>
                </div>
                <Badge
                  variant="secondary"
                  className={`text-[11px] font-medium ${article.tagColor}`}
                >
                  {article.tag}
                </Badge>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Clock className="h-3 w-3" />
                  {article.time}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Eye className="h-3 w-3" />
                  {formatNumber(article.views)}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <MessageSquare className="h-3 w-3" />
                  {article.comments}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <ThumbsUp className="h-3 w-3" />
                  {article.likes}
                </span>
              </div>
            </div>

            {/* Cover placeholder - only on large articles */}
            {article.cover && (
              <div className="hidden h-24 w-36 shrink-0 overflow-hidden rounded-md bg-secondary/60 sm:block">
                <div className="flex h-full w-full items-center justify-center">
                  <span className="font-mono text-2xl font-bold text-primary/30">Go</span>
                </div>
              </div>
            )}
          </div>
        </article>
      ))}
    </div>
  )
}
