/**
 * BookGrid - 书籍网格展示组件（Client Component）
 *
 * 注意：此组件使用的是 Mock 数据，因为该组件在客户端运行（带有分类筛选交互）。
 * 实际书籍列表数据由 app/books/page.tsx（SSR）通过 /api/v1/books 接口获取。
 * 此组件目前仅用于首页等无需 SEO 的展示场景，后续如需接入 API 请改为服务端组件或
 * 将数据作为 props 传入。
 *
 * TODO: 将 mock 数据替换为从父组件传入的真实数据，或通过 SWR/TanStack Query 客户端获取。
 * TODO: mock 数据中 url 为 "#" 的书籍没有真实链接，接入 API 后需改为内部详情页链接 /book/[id]。
 */
"use client"

import { useState } from "react"
import { Star, ExternalLink, BookOpen, Users } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

const filters = ["全部", "入门", "进阶", "实战", "原理", "开源免费"]

interface Book {
  id: number
  title: string
  author: string
  cover: string
  rating: number
  ratingCount: number
  description: string
  tags: string[]
  category: string
  free?: boolean
  url: string
}

const books: Book[] = [
  {
    id: 1,
    title: "Go 程序设计语言",
    author: "Alan A. A. Donovan / Brian W. Kernighan",
    cover: "GOPL",
    rating: 9.2,
    ratingCount: 1256,
    description:
      "被誉为 Go 语言的「圣经」，由 K&R 的作者之一撰写，系统全面地介绍了 Go 语言的方方面面。",
    tags: ["经典", "全面"],
    category: "入门",
    url: "#",
  },
  {
    id: 2,
    title: "Go 语言实战",
    author: "William Kennedy / Brian Ketelsen / Erik St. Martin",
    cover: "GIA",
    rating: 8.4,
    ratingCount: 892,
    description:
      "面向有一定编程基础的开发者，通过实际项目带领读者深入理解 Go 的并发模型、类型系统和标准库。",
    tags: ["实战", "并发"],
    category: "实战",
    url: "#",
  },
  {
    id: 3,
    title: "Go Web 编程",
    author: "Sau Sheong Chang",
    cover: "GWP",
    rating: 7.8,
    ratingCount: 567,
    description:
      "专注于 Web 开发的 Go 语言教程，涵盖 HTTP 处理、模板引擎、数据持久化、测试等 Web 开发核心话题。",
    tags: ["Web", "HTTP"],
    category: "实战",
    url: "#",
  },
  {
    id: 4,
    title: "Go 语言高级编程",
    author: "柴树杉 / 曹春晖",
    cover: "AGP",
    rating: 8.6,
    ratingCount: 734,
    description:
      "深入探讨 CGO、汇编语言、RPC、Web 框架实现和分布式系统等高级话题。适合有经验的 Go 开发者。",
    tags: ["高级", "CGO", "RPC"],
    category: "进阶",
    free: true,
    url: "https://github.com/chai2010/advanced-go-programming-book",
  },
  {
    id: 5,
    title: "Go 语言设计与实现",
    author: "左书祺 (Draveness)",
    cover: "GDI",
    rating: 9.0,
    ratingCount: 1089,
    description:
      "从源码层面解析 Go 语言的设计原理和实现机制，涵盖编译器、运行时、调度器、内存管理和垃圾回收。",
    tags: ["源码", "原理", "运行时"],
    category: "原理",
    free: true,
    url: "https://draveness.me/golang/",
  },
  {
    id: 6,
    title: "Go 语言核心编程",
    author: "李文塔",
    cover: "GCP",
    rating: 7.6,
    ratingCount: 423,
    description:
      "从类型系统、接口、并发到反射，系统深入讲解 Go 核心特性。包含大量代码示例和实践建议。",
    tags: ["核心", "接口", "并发"],
    category: "进阶",
    url: "#",
  },
  {
    id: 7,
    title: "Go 语言学习笔记",
    author: "雨痕",
    cover: "GLN",
    rating: 8.8,
    ratingCount: 678,
    description:
      "以笔记形式深入剖析 Go 语言源码和运行时实现原理，文风简练，干货满满。",
    tags: ["源码", "笔记", "深入"],
    category: "原理",
    url: "#",
  },
  {
    id: 8,
    title: "Go 语言从入门到进阶实战",
    author: "徐波",
    cover: "GBP",
    rating: 7.4,
    ratingCount: 345,
    description:
      "面向零基础读者的 Go 语言教程，循序渐进地讲解语法基础、常用标准库和实战项目。",
    tags: ["入门", "零基础"],
    category: "入门",
    url: "#",
  },
  {
    id: 9,
    title: "Go 语言精进之路",
    author: "白明",
    cover: "GMA",
    rating: 8.5,
    ratingCount: 567,
    description:
      "总结了大量 Go 编程最佳实践和经验教训，帮助开发者写出更地道、更高效的 Go 代码。",
    tags: ["最佳实践", "进阶"],
    category: "进阶",
    url: "#",
  },
]

function StarRating({ rating }: { rating: number }) {
  const fullStars = Math.floor(rating / 2)
  const hasHalf = (rating / 2) % 1 >= 0.5

  return (
    <div className="flex items-center gap-0.5">
      {Array.from({ length: 5 }).map((_, i) => (
        <Star
          key={i}
          className={cn(
            "h-3 w-3",
            i < fullStars
              ? "fill-chart-3 text-chart-3"
              : i === fullStars && hasHalf
                ? "fill-chart-3/50 text-chart-3"
                : "fill-muted text-muted"
          )}
        />
      ))}
    </div>
  )
}

// Generate a consistent color for book covers
function getBookColor(id: number): string {
  const colors = [
    "from-primary/80 to-primary",
    "from-chart-2/80 to-chart-2",
    "from-chart-3/80 to-chart-3",
    "from-chart-4/80 to-chart-4",
    "from-chart-5/80 to-chart-5",
  ]
  return colors[id % colors.length]
}

export function BookGrid() {
  const [activeFilter, setActiveFilter] = useState("全部")

  const filtered =
    activeFilter === "全部"
      ? books
      : activeFilter === "开源免费"
        ? books.filter((b) => b.free)
        : books.filter((b) => b.category === activeFilter)

  return (
    <div>
      {/* Filter tabs */}
      <div className="mb-6 flex flex-wrap items-center gap-2">
        {filters.map((f) => (
          <button
            key={f}
            onClick={() => setActiveFilter(f)}
            className={cn(
              "cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-all",
              activeFilter === f
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
          >
            {f}
          </button>
        ))}
      </div>

      {/* Books grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {filtered.map((book) => (
          <Card
            key={book.id}
            className="group overflow-hidden transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="p-0">
              {/* Book cover */}
              <div
                className={cn(
                  "flex h-44 items-center justify-center bg-gradient-to-br",
                  getBookColor(book.id)
                )}
              >
                <div className="flex flex-col items-center gap-2">
                  <BookOpen className="h-10 w-10 text-primary-foreground/80" />
                  <span className="font-mono text-lg font-bold text-primary-foreground/90">
                    {book.cover}
                  </span>
                </div>
              </div>

              {/* Info */}
              <div className="p-4">
                <div className="flex items-start justify-between gap-2">
                  <h3 className="text-sm font-bold leading-snug text-foreground group-hover:text-primary">
                    {book.title}
                  </h3>
                  {book.free && (
                    <Badge className="shrink-0 bg-accent/10 text-[10px] font-semibold text-accent">
                      {"免费"}
                    </Badge>
                  )}
                </div>
                <p className="mt-1 text-xs text-muted-foreground">
                  {book.author}
                </p>

                {/* Rating */}
                <div className="mt-2 flex items-center gap-2">
                  <StarRating rating={book.rating} />
                  <span className="text-xs font-semibold text-foreground">
                    {book.rating}
                  </span>
                  <span className="flex items-center gap-0.5 text-[10px] text-muted-foreground">
                    <Users className="h-2.5 w-2.5" />
                    {book.ratingCount}
                  </span>
                </div>

                {/* Description */}
                <p className="mt-2 line-clamp-2 text-xs leading-relaxed text-muted-foreground">
                  {book.description}
                </p>

                {/* Tags + Link */}
                <div className="mt-3 flex items-center justify-between">
                  <div className="flex gap-1">
                    {book.tags.slice(0, 2).map((tag) => (
                      <Badge
                        key={tag}
                        variant="secondary"
                        className="text-[10px] font-medium"
                      >
                        {tag}
                      </Badge>
                    ))}
                  </div>
                  {/* url 为 "#" 表示暂无真实链接（mock 数据），接入 API 后替换为 /book/[id] */}
                  {book.url && book.url !== "#" ? (
                    <a
                      href={book.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs font-medium text-primary transition-colors hover:text-primary/80"
                      aria-label={"查看 " + book.title}
                    >
                      {"查看"}
                      <ExternalLink className="h-3 w-3" />
                    </a>
                  ) : (
                    <span className="text-xs text-muted-foreground">暂无链接</span>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
