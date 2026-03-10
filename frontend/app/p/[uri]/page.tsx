import { notFound } from "next/navigation"
import Link from "next/link"
import {
  Star,
  GitFork,
  Eye,
  ExternalLink,
  Github,
  Globe,
  Download,
  MessageSquare,
  ThumbsUp,
  Clock,
  User,
} from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { ProjectDetailData } from "@/lib/types"

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const base = process.env.API_BASE_URL || "http://localhost:8090"
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

interface ProjectPageProps {
  params: Promise<{ uri: string }>
}

export async function generateMetadata({ params }: ProjectPageProps): Promise<Metadata> {
  const { uri } = await params
  const data = await fetchAPI<ProjectDetailData>(
    `/projects/${uri}`,
    { next: { revalidate: 60 } }
  )

  if (!data?.project) {
    return { title: "项目未找到 - Go语言中文网" }
  }

  const { project } = data
  return {
    title: `${project.name} - Go 开源项目 - Go语言中文网`,
    description: project.desc || `${project.name} 是一个 Go 语言开源项目`,
    openGraph: {
      title: project.name,
      description: project.desc || "",
      images: project.logo ? [{ url: project.logo }] : [],
    },
  }
}

export default async function ProjectDetailPage({ params }: ProjectPageProps) {
  const { uri } = await params
  const data = await fetchAPI<ProjectDetailData>(
    `/projects/${uri}`,
    { next: { revalidate: 60 } }
  )

  if (!data?.project) {
    notFound()
  }

  const { project, comments } = data

  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: project.name,
    description: project.desc,
    url: project.homepage || project.src_url,
    applicationCategory: project.category,
    programmingLanguage: project.lang || "Go",
    author: {
      "@type": "Person",
      name: project.author,
    },
  }

  return (
    <PageLayout>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      <PageHeader
        title={project.name}
        description={project.desc}
        breadcrumbs={[
          { label: "开源项目", href: "/projects" },
          { label: project.name },
        ]}
      />

      {/* Project main card */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex items-start gap-4">
            {project.logo ? (
              <img
                src={project.logo}
                alt={project.name}
                className="h-16 w-16 shrink-0 rounded-xl object-cover"
              />
            ) : (
              <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-primary/10">
                <Github className="h-8 w-8 text-primary" />
              </div>
            )}

            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-xl font-bold text-foreground">{project.name}</h1>
                {project.category && (
                  <Badge variant="secondary">{project.category}</Badge>
                )}
                {project.lang && (
                  <Badge variant="outline">{project.lang}</Badge>
                )}
              </div>
              <p className="mt-1 text-sm text-muted-foreground">
                by{" "}
                <Link
                  href={`/user/${project.author}`}
                  className="font-medium text-foreground hover:text-primary"
                >
                  {project.author}
                </Link>
              </p>
              {project.desc && (
                <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                  {project.desc}
                </p>
              )}
            </div>
          </div>

          {/* Stats row */}
          <div className="mt-6 flex flex-wrap gap-6 border-t border-border pt-5">
            <div className="flex items-center gap-2">
              <Star className="h-4 w-4 text-chart-3" />
              <span className="text-sm font-semibold">{formatNum(project.star)}</span>
              <span className="text-xs text-muted-foreground">Stars</span>
            </div>
            <div className="flex items-center gap-2">
              <GitFork className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.fork)}</span>
              <span className="text-xs text-muted-foreground">Forks</span>
            </div>
            <div className="flex items-center gap-2">
              <Eye className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.watch)}</span>
              <span className="text-xs text-muted-foreground">Watchers</span>
            </div>
            <div className="flex items-center gap-2">
              <Eye className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.viewnum)}</span>
              <span className="text-xs text-muted-foreground">浏览</span>
            </div>
            <div className="flex items-center gap-2">
              <ThumbsUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.likenum)}</span>
              <span className="text-xs text-muted-foreground">点赞</span>
            </div>
          </div>

          {/* Links */}
          <div className="mt-5 flex flex-wrap gap-3">
            {project.src_url && (
              <a
                href={project.src_url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <Github className="h-4 w-4" />
                源码仓库
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.homepage && (
              <a
                href={project.homepage}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <Globe className="h-4 w-4" />
                项目主页
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.doc_url && (
              <a
                href={project.doc_url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <ExternalLink className="h-4 w-4" />
                文档
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.download_url && (
              <a
                href={project.download_url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <Download className="h-4 w-4" />
                下载
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Comments */}
      <div className="space-y-4">
        <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
          <MessageSquare className="h-4 w-4" />
          评论
          {comments && comments.length > 0 && (
            <Badge variant="secondary" className="text-xs">
              {comments.length}
            </Badge>
          )}
        </h2>

        {!comments || comments.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-10 text-center">
            <MessageSquare className="mb-2 h-8 w-8 text-muted-foreground/40" />
            <p className="text-sm text-muted-foreground">暂无评论，来发表第一条评论吧</p>
          </div>
        ) : (
          <div className="space-y-3">
            {comments.map((comment) => (
              <Card key={comment.id}>
                <CardContent className="p-4">
                  <div className="flex items-start gap-3">
                    {comment.avatar ? (
                      <img
                        src={comment.avatar}
                        alt={comment.name}
                        className="h-8 w-8 shrink-0 rounded-full object-cover"
                      />
                    ) : (
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                        <User className="h-4 w-4 text-primary" />
                      </div>
                    )}
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-foreground">
                          {comment.name}
                        </span>
                        <span className="text-xs text-muted-foreground/60">#{comment.floor}</span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {comment.ctime}
                        </span>
                      </div>
                      <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                        {comment.content}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </PageLayout>
  )
}
