"use client"

import { useState } from "react"
import Link from "next/link"
import {
  ExternalLink,
  ThumbsUp,
  MessageSquare,
  Clock,
  Tag,
  Download,
  Globe,
  BookOpen,
  Video,
  Wrench,
  GraduationCap,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

const categories = [
  { label: "全部", icon: Tag, count: 256 },
  { label: "教程指南", icon: GraduationCap, count: 68 },
  { label: "开发工具", icon: Wrench, count: 45 },
  { label: "视频课程", icon: Video, count: 32 },
  { label: "网站推荐", icon: Globe, count: 51 },
  { label: "文档翻译", icon: BookOpen, count: 29 },
  { label: "下载资源", icon: Download, count: 31 },
]

interface Resource {
  id: number
  title: string
  description: string
  url: string
  category: string
  tags: string[]
  author: string
  time: string
  likes: number
  comments: number
  featured?: boolean
}

const resources: Resource[] = [
  {
    id: 1,
    title: "Go by Example - Go 语言实例教程中文版",
    description:
      "通过注释丰富的示例代码来学习 Go 语言，覆盖从基础语法到高级特性的方方面面。适合初学者快速上手。",
    url: "https://gobyexample.com",
    category: "教程指南",
    tags: ["入门", "教程", "实例"],
    author: "polaris",
    time: "3 天前",
    likes: 342,
    comments: 28,
    featured: true,
  },
  {
    id: 2,
    title: "GoLand - JetBrains 出品的 Go IDE",
    description:
      "功能强大的 Go 语言集成开发环境，提供智能代码补全、重构、调试、测试等一站式开发体验。",
    url: "https://www.jetbrains.com/go/",
    category: "开发工具",
    tags: ["IDE", "开发工具", "JetBrains"],
    author: "alice_go",
    time: "5 天前",
    likes: 256,
    comments: 19,
  },
  {
    id: 3,
    title: "Go 语言高级编程 - 开源电子书",
    description:
      "涵盖 CGO、汇编语言、RPC、Web 框架、分布式系统等高级话题的开源图书。适合有一定基础的 Go 开发者深入学习。",
    url: "https://github.com/chai2010/advanced-go-programming-book",
    category: "文档翻译",
    tags: ["高级", "电子书", "开源"],
    author: "chai2010",
    time: "1 周前",
    likes: 489,
    comments: 45,
    featured: true,
  },
  {
    id: 4,
    title: "GopherChina 技术大会视频合集",
    description:
      "历年 GopherChina 大会精彩演讲视频回放，来自各大厂的 Go 实战经验分享。",
    url: "https://www.bilibili.com",
    category: "视频课程",
    tags: ["大会", "视频", "实战"],
    author: "bob_dev",
    time: "1 周前",
    likes: 198,
    comments: 12,
  },
  {
    id: 5,
    title: "Go 语言官方标准库中文翻译",
    description:
      "Go 标准库完整中文翻译项目，由社区协作完成。持续更新到最新 Go 版本，方便中文开发者查阅。",
    url: "https://studygolang.com/pkgdoc",
    category: "文档翻译",
    tags: ["标准库", "中文", "翻译"],
    author: "polaris",
    time: "2 周前",
    likes: 567,
    comments: 38,
    featured: true,
  },
  {
    id: 6,
    title: "Delve - Go 语言调试器",
    description:
      "Go 语言的专用调试器，支持断点、变量查看、Goroutine 调试等功能。与 VS Code 和 GoLand 深度集成。",
    url: "https://github.com/go-delve/delve",
    category: "开发工具",
    tags: ["调试", "Debug", "工具"],
    author: "charlie",
    time: "2 周前",
    likes: 145,
    comments: 8,
  },
  {
    id: 7,
    title: "Effective Go 中文版",
    description:
      "Go 官方推荐的编码指南中文版，介绍如何编写清晰、地道的 Go 代码，是每位 Gopher 必读的文档。",
    url: "https://go.dev/doc/effective_go",
    category: "教程指南",
    tags: ["官方", "最佳实践", "编码规范"],
    author: "polaris",
    time: "3 周前",
    likes: 623,
    comments: 52,
    featured: true,
  },
  {
    id: 8,
    title: "Go 语言入门到精通 - B站视频教程",
    description:
      "系统化的 Go 语言视频教程，从零基础到实战项目，配套源码和练习，适合自学入门。",
    url: "https://www.bilibili.com",
    category: "视频课程",
    tags: ["入门", "视频", "零基础"],
    author: "guo_hongzhi",
    time: "3 周前",
    likes: 312,
    comments: 67,
  },
  {
    id: 9,
    title: "awesome-go 中文版 - Go 资源集合",
    description:
      "精心收集的 Go 框架、库和软件列表，按分类整理，持续更新。发现 Go 生态的优秀项目。",
    url: "https://awesome-go.com",
    category: "网站推荐",
    tags: ["集合", "awesome", "生态"],
    author: "Minnie",
    time: "1 个月前",
    likes: 456,
    comments: 23,
  },
  {
    id: 10,
    title: "Go 各版本安装包下载镜像",
    description:
      "提供 Go 语言所有版本的国内高速下载镜像，包括 Windows、macOS、Linux 全平台安装包。",
    url: "https://studygolang.com/dl",
    category: "下载资源",
    tags: ["下载", "安装", "镜像"],
    author: "polaris",
    time: "1 个月前",
    likes: 789,
    comments: 31,
  },
]

function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + "k"
  }
  return num.toString()
}

export function ResourceList() {
  const [activeCategory, setActiveCategory] = useState("全部")

  const filtered =
    activeCategory === "全部"
      ? resources
      : resources.filter((r) => r.category === activeCategory)

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      {/* Sidebar categories */}
      <aside className="w-full shrink-0 lg:w-56">
        <div className="rounded-lg border border-border bg-card p-3">
          <h3 className="mb-2 px-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            {"资源分类"}
          </h3>
          <nav className="space-y-0.5">
            {categories.map((cat) => (
              <button
                key={cat.label}
                onClick={() => setActiveCategory(cat.label)}
                className={cn(
                  "flex w-full items-center gap-2.5 rounded-md px-2.5 py-2 text-left text-sm font-medium transition-colors",
                  activeCategory === cat.label
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                )}
              >
                <cat.icon className="h-4 w-4 shrink-0" />
                <span className="flex-1">{cat.label}</span>
                <span
                  className={cn(
                    "text-xs",
                    activeCategory === cat.label
                      ? "text-primary"
                      : "text-muted-foreground/60"
                  )}
                >
                  {cat.count}
                </span>
              </button>
            ))}
          </nav>
        </div>
      </aside>

      {/* Resource list */}
      <div className="min-w-0 flex-1 space-y-3">
        {/* Result header */}
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            {"共 "}<span className="font-semibold text-foreground">{filtered.length}</span>{" 个资源"}
          </p>
          <div className="flex items-center gap-2">
            {["最新", "最热", "推荐"].map((sort, i) => (
              <button
                key={sort}
                className={cn(
                  "rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
                  i === 0
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-secondary"
                )}
              >
                {sort}
              </button>
            ))}
          </div>
        </div>

        {/* Cards */}
        {filtered.map((resource) => (
          <Card
            key={resource.id}
            className={cn(
              "group transition-all hover:border-primary/20 hover:shadow-sm",
              resource.featured && "border-primary/10 bg-primary/[0.02]"
            )}
          >
            <CardContent className="p-4 sm:p-5">
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0 flex-1">
                  {/* Title + external link */}
                  <div className="flex items-start gap-2">
                    {resource.featured && (
                      <Badge className="mt-0.5 shrink-0 bg-primary/10 text-[10px] font-semibold text-primary">
                        {"精选"}
                      </Badge>
                    )}
                    <Link
                      href={`/resources/${resource.id}`}
                      className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                    >
                      {resource.title}
                    </Link>
                  </div>

                  {/* Description */}
                  <p className="mt-1.5 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                    {resource.description}
                  </p>

                  {/* Tags + Meta */}
                  <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                    <div className="flex items-center gap-1.5">
                      {resource.tags.map((tag) => (
                        <Badge
                          key={tag}
                          variant="secondary"
                          className="text-[10px] font-medium"
                        >
                          {tag}
                        </Badge>
                      ))}
                    </div>
                    <span className="text-xs text-muted-foreground">
                      {resource.author}
                    </span>
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Clock className="h-3 w-3" />
                      {resource.time}
                    </span>
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <ThumbsUp className="h-3 w-3" />
                      {formatNumber(resource.likes)}
                    </span>
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <MessageSquare className="h-3 w-3" />
                      {resource.comments}
                    </span>
                  </div>
                </div>

                {/* External link button */}
                <a
                  href={resource.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
                  aria-label={"访问 " + resource.title}
                >
                  <ExternalLink className="h-4 w-4" />
                </a>
              </div>
            </CardContent>
          </Card>
        ))}

        {/* Load more */}
        <div className="pt-4 text-center">
          <button className="rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80">
            {"加载更多资源"}
          </button>
        </div>
      </div>
    </div>
  )
}
