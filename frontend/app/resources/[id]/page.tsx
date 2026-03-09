import { notFound } from "next/navigation"
import Link from "next/link"
import {
  ExternalLink,
  ThumbsUp,
  MessageSquare,
  Clock,
  Eye,
  User,
} from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { Resource, Comment } from "@/lib/types"

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

interface ResourceDetailData {
  resource: Resource
  comments?: Comment[]
}

function formatNum(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + "k"
  return String(n)
}

interface ResourceDetailPageProps {
  params: Promise<{ id: string }>
}

export async function generateMetadata({ params }: ResourceDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const data = await fetchAPI<ResourceDetailData>(
    `/resources/${id}`,
    { next: { revalidate: 60 } }
  )

  if (!data?.resource) {
    return { title: "资源未找到 - Go语言中文网" }
  }

  const { resource } = data
  return {
    title: `${resource.title} - Go 资源 - Go语言中文网`,
    description: resource.desc || `${resource.title} - Go语言中文网资源`,
    openGraph: {
      title: resource.title,
      description: resource.desc || "",
      images: resource.cover ? [{ url: resource.cover }] : [],
    },
  }
}

export default async function ResourceDetailPage({ params }: ResourceDetailPageProps) {
  const { id } = await params
  const data = await fetchAPI<ResourceDetailData>(
    `/resources/${id}`,
    { next: { revalidate: 60 } }
  )

  if (!data?.resource) {
    notFound()
  }

  const { resource, comments } = data

  return (
    <PageLayout>
      <PageHeader
        title={resource.title}
        breadcrumbs={[
          { label: "资源索引", href: "/resources" },
          { label: resource.title },
        ]}
      />

      {/* Resource main card */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex items-start gap-4">
            {resource.cover ? (
              <img
                src={resource.cover}
                alt={resource.title}
                className="h-20 w-20 shrink-0 rounded-lg object-cover"
              />
            ) : (
              <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                <ExternalLink className="h-8 w-8 text-primary" />
              </div>
            )}

            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-xl font-bold text-foreground">{resource.title}</h1>
                {resource.catname && (
                  <Badge variant="secondary">{resource.catname}</Badge>
                )}
              </div>
              <p className="mt-1 text-sm text-muted-foreground">
                by{" "}
                <Link
                  href={`/user/${resource.author}`}
                  className="font-medium text-foreground hover:text-primary"
                >
                  {resource.author}
                </Link>
              </p>
              {resource.desc && (
                <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                  {resource.desc}
                </p>
              )}
            </div>
          </div>

          {/* Stats row */}
          <div className="mt-6 flex flex-wrap gap-6 border-t border-border pt-5">
            <div className="flex items-center gap-2">
              <Eye className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(resource.viewnum)}</span>
              <span className="text-xs text-muted-foreground">浏览</span>
            </div>
            <div className="flex items-center gap-2">
              <ThumbsUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(resource.likenum)}</span>
              <span className="text-xs text-muted-foreground">点赞</span>
            </div>
            <div className="flex items-center gap-2">
              <MessageSquare className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{resource.cmtnum}</span>
              <span className="text-xs text-muted-foreground">评论</span>
            </div>
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span className="text-xs text-muted-foreground">{resource.ctime}</span>
            </div>
          </div>

          {/* External link */}
          {resource.url && (
            <div className="mt-5">
              <a
                href={resource.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 rounded-md border border-border px-4 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <ExternalLink className="h-4 w-4" />
                访问资源
              </a>
            </div>
          )}
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
