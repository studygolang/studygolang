import { Suspense } from "react"
import Link from "next/link"
import { Plus, Star, GitFork, Eye, ExternalLink, Github } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { ProjectListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "Go 开源项目 - Go语言中文网",
  description: "发现优秀的 Go 语言开源项目，分享你的作品",
}

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const base = process.env.API_BASE_URL || "http://localhost:8088"
    const res = await fetch(`${base}/api/v1${path}`, options)
    if (!res.ok) return null
    const json = await res.json()
    return json.code === 0 ? json.data : null
  } catch {
    return null
  }
}

function formatNum(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + "k"
  return String(n)
}

async function ProjectList({ page }: { page: number }) {
  const data = await fetchAPI<ProjectListData>(
    `/projects?p=${page}`,
    { cache: "no-store" }
  )

  if (!data || !data.projects || data.projects.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <Github className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">暂无项目数据</p>
      </div>
    )
  }

  return (
    <div>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {data.projects.map((project) => (
          <Card
            key={project.id}
            className="group transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="flex h-full flex-col p-5">
              {/* Header */}
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  {project.logo ? (
                    <img
                      src={project.logo}
                      alt={project.name}
                      className="h-9 w-9 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10">
                      <Github className="h-4.5 w-4.5 text-primary" />
                    </div>
                  )}
                  <div>
                    <Link href={`/p/${project.uri}`}>
                      <h3 className="text-sm font-bold text-foreground group-hover:text-primary">
                        {project.name}
                      </h3>
                    </Link>
                    <p className="text-xs text-muted-foreground">{project.author}</p>
                  </div>
                </div>
                {project.src_url && (
                  <a
                    href={project.src_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-secondary hover:text-primary"
                    aria-label={"查看 " + project.name + " 源码"}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </a>
                )}
              </div>

              {/* Description */}
              <p className="mt-3 flex-1 line-clamp-3 text-sm leading-relaxed text-muted-foreground">
                {project.desc || "暂无描述"}
              </p>

              {/* Category */}
              {project.category && (
                <div className="mt-3 flex flex-wrap gap-1.5">
                  <Badge variant="secondary" className="text-[10px] font-medium">
                    {project.category}
                  </Badge>
                  {project.lang && (
                    <Badge variant="outline" className="text-[10px] font-medium">
                      {project.lang}
                    </Badge>
                  )}
                </div>
              )}

              {/* Stats */}
              <div className="mt-3 flex items-center gap-4 border-t border-border pt-3">
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Star className="h-3.5 w-3.5" />
                  {formatNum(project.star)}
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <GitFork className="h-3.5 w-3.5" />
                  {formatNum(project.fork)}
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Eye className="h-3.5 w-3.5" />
                  {formatNum(project.watch)}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Pagination */}
      <div className="mt-8 flex items-center justify-center gap-2">
        {page > 1 && (
          <Link
            href={`/projects?page=${page - 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            上一页
          </Link>
        )}
        <span className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground">
          {page}
        </span>
        {data.has_more && (
          <Link
            href={`/projects?page=${page + 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </div>
  )
}

interface ProjectsPageProps {
  searchParams: Promise<{ page?: string }>
}

export default async function ProjectsPage({ searchParams }: ProjectsPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.page || "1", 10))

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="开源项目"
        description="发现和分享优秀的 Go 语言开源项目"
        breadcrumbs={[{ label: "开源项目" }]}
        actions={
          <Link href="/projects/new">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <Plus className="h-3.5 w-3.5" />
              提交项目
            </Button>
          </Link>
        }
      />
      <Suspense
        fallback={
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-48 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <ProjectList page={page} />
      </Suspense>
    </PageLayout>
  )
}
