/**
 * ProjectGrid - 项目网格展示组件（Client Component）
 *
 * 注意：此组件使用的是 Mock 数据，因为该组件在客户端运行（带有分类筛选交互）。
 * 实际项目列表数据由 app/projects/page.tsx（SSR）通过 /api/v1/projects 接口获取。
 * 此组件目前仅用于首页等无需 SEO 的展示场景，后续如需接入 API 请改为服务端组件或
 * 将数据作为 props 传入。
 *
 * TODO: 将 mock 数据替换为从父组件传入的真实数据，或通过 SWR/TanStack Query 客户端获取。
 */
"use client"

import { useState } from "react"
import { Star, GitFork, ExternalLink, Github } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

const categories = [
  "全部",
  "Web框架",
  "微服务",
  "数据库",
  "CLI工具",
  "DevOps",
  "AI/ML",
  "中间件",
  "测试",
  "其他",
]

interface Project {
  id: number
  name: string
  description: string
  author: string
  stars: string
  forks: string
  category: string
  tags: string[]
  language: string
  languageColor: string
  url: string
}

const projects: Project[] = [
  {
    id: 1,
    name: "gin",
    description:
      "Gin 是一个用 Go 编写的 HTTP Web 框架。它提供了类似 Martini 的 API，但性能更好，速度最高可提升 40 倍。",
    author: "gin-gonic",
    stars: "79.8k",
    forks: "8.1k",
    category: "Web框架",
    tags: ["HTTP", "Web", "Router", "Middleware"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/gin-gonic/gin",
  },
  {
    id: 2,
    name: "go-zero",
    description:
      "go-zero 是一个集成了各种工程实践的 Web 和 RPC 框架，内建级联超时控制、限流、自适应熔断等微服务治理能力。",
    author: "zeromicro",
    stars: "29.5k",
    forks: "3.9k",
    category: "微服务",
    tags: ["RPC", "微服务", "API Gateway"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/zeromicro/go-zero",
  },
  {
    id: 3,
    name: "gorm",
    description:
      "Go 语言最流行的 ORM 框架，支持 MySQL、PostgreSQL、SQLite、SQL Server 等多种数据库。提供完善的关联关系和迁移功能。",
    author: "go-gorm",
    stars: "37.2k",
    forks: "3.9k",
    category: "数据库",
    tags: ["ORM", "Database", "MySQL", "PostgreSQL"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/go-gorm/gorm",
  },
  {
    id: 4,
    name: "cobra",
    description:
      "Cobra 是一个用于创建命令行应用的库，被 Kubernetes、Docker、Hugo 等知名项目采用。",
    author: "spf13",
    stars: "38.5k",
    forks: "2.8k",
    category: "CLI工具",
    tags: ["CLI", "Command Line", "Flag"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/spf13/cobra",
  },
  {
    id: 5,
    name: "kratos",
    description:
      "Kratos 是 bilibili 开源的一套 Go 微服务框架，包含大量微服务相关的框架和工具集。",
    author: "go-kratos",
    stars: "23.8k",
    forks: "4.0k",
    category: "微服务",
    tags: ["gRPC", "微服务", "Protobuf"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/go-kratos/kratos",
  },
  {
    id: 6,
    name: "eino",
    description:
      "字节跳动开源的 Go 语言 LLM 应用开发框架，提供 Chain、Agent、RAG 等核心抽象和丰富的组件实现。",
    author: "cloudwego",
    stars: "12.3k",
    forks: "1.2k",
    category: "AI/ML",
    tags: ["LLM", "AI", "Agent", "RAG"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/cloudwego/eino",
  },
  {
    id: 7,
    name: "etcd",
    description:
      "分布式可靠的键值存储，用于分布式系统最关键的数据。是 Kubernetes 的核心组件之一。",
    author: "etcd-io",
    stars: "48.1k",
    forks: "9.8k",
    category: "中间件",
    tags: ["KV Store", "分布式", "Raft"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/etcd-io/etcd",
  },
  {
    id: 8,
    name: "fiber",
    description:
      "一个受 Express.js 启发的 Go Web 框架，构建在 Fasthttp 之上，提供极致的性能和零内存分配。",
    author: "gofiber",
    stars: "34.6k",
    forks: "1.7k",
    category: "Web框架",
    tags: ["HTTP", "Web", "Fasthttp"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/gofiber/fiber",
  },
  {
    id: 9,
    name: "testify",
    description:
      "Go 测试工具套件，提供断言、mock 和测试套件功能，使测试代码更加清晰和可读。",
    author: "stretchr",
    stars: "23.9k",
    forks: "1.6k",
    category: "测试",
    tags: ["Testing", "Assert", "Mock"],
    language: "Go",
    languageColor: "bg-primary",
    url: "https://github.com/stretchr/testify",
  },
]

export function ProjectGrid() {
  const [activeCategory, setActiveCategory] = useState("全部")

  const filtered =
    activeCategory === "全部"
      ? projects
      : projects.filter((p) => p.category === activeCategory)

  return (
    <div>
      {/* Category filter */}
      <div className="mb-5 flex flex-wrap items-center gap-2">
        {categories.map((cat) => (
          <button
            key={cat}
            onClick={() => setActiveCategory(cat)}
            className={cn(
              "cursor-pointer rounded-full px-3.5 py-1.5 text-sm font-medium transition-all",
              activeCategory === cat
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
          >
            {cat}
          </button>
        ))}
      </div>

      {/* Project cards grid */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {filtered.map((project) => (
          <Card
            key={project.id}
            className="group transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="flex h-full flex-col p-5">
              {/* Header */}
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10">
                    <Github className="h-4.5 w-4.5 text-primary" />
                  </div>
                  <div>
                    <h3 className="text-sm font-bold text-foreground group-hover:text-primary">
                      {project.name}
                    </h3>
                    <p className="text-xs text-muted-foreground">{project.author}</p>
                  </div>
                </div>
                <a
                  href={project.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-secondary hover:text-primary"
                  aria-label={"在 GitHub 上查看 " + project.name}
                >
                  <ExternalLink className="h-4 w-4" />
                </a>
              </div>

              {/* Description */}
              <p className="mt-3 flex-1 line-clamp-3 text-sm leading-relaxed text-muted-foreground">
                {project.description}
              </p>

              {/* Tags */}
              <div className="mt-3 flex flex-wrap gap-1.5">
                {project.tags.slice(0, 3).map((tag) => (
                  <Badge
                    key={tag}
                    variant="secondary"
                    className="text-[10px] font-medium"
                  >
                    {tag}
                  </Badge>
                ))}
              </div>

              {/* Stats */}
              <div className="mt-3 flex items-center gap-4 border-t border-border pt-3">
                <div className="flex items-center gap-1">
                  <span className={`h-2.5 w-2.5 rounded-full ${project.languageColor}`} />
                  <span className="text-xs text-muted-foreground">{project.language}</span>
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Star className="h-3.5 w-3.5" />
                  {project.stars}
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <GitFork className="h-3.5 w-3.5" />
                  {project.forks}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* TODO: 加载更多 - 当前为静态 mock，接入真实 API 后需实现分页逻辑 */}
      <div className="mt-8 text-center">
        <button
          className="rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
          disabled
        >
          {"加载更多项目"}
        </button>
      </div>
    </div>
  )
}
