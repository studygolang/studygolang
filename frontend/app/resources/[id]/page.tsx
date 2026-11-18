import { notFound } from "next/navigation"
import Link from "next/link"
import { cookies } from "next/headers"
import { ExternalLink, ThumbsUp, MessageSquare, Clock, Eye } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { CommentSection } from "@/components/comment-section"
import { ContentActions } from "@/components/content-actions"
import { MarkdownContent } from "@/components/markdown-content"
import { OBJTYPE_RESOURCE, type Resource, type Comment } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { formatNum } from "@/lib/utils"

// 点赞/收藏状态与评论 liked 都是当前用户维度的数据，不能进 ISR 缓存
export const dynamic = "force-dynamic"

async function authHeaders(): Promise<Record<string, string>> {
  const store = await cookies()
  const cookie = store
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")
  return cookie ? { Cookie: cookie } : {}
}

interface ResourceDetailData {
  resource: Resource
  comments?: Comment[]
  likeflag?: number
  hadcollect?: number
}

interface ResourceDetailPageProps {
  params: Promise<{ id: string }>
}

async function getResourceDetail(id: string): Promise<ResourceDetailData | null> {
  return fetchAPINullable<ResourceDetailData>(`/resources/${id}`, {
    cache: "no-store",
    headers: await authHeaders(),
  })
}

export async function generateMetadata({ params }: ResourceDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const data = await getResourceDetail(id)

  if (!data?.resource) {
    return { title: "资源未找到 - Go语言中文网" }
  }

  const { resource } = data
  // 后端 Resource 无 desc 字段（用 content），无 cover 字段
  return {
    title: `${resource.title} - Go 资源 - Go语言中文网`,
    description: resource.content || `${resource.title} - Go语言中文网资源`,
    openGraph: {
      title: resource.title,
      description: resource.content || "",
    },
  }
}

export default async function ResourceDetailPage({ params }: ResourceDetailPageProps) {
  const { id } = await params
  const data = await getResourceDetail(id)

  if (!data?.resource) {
    return notFound()
  }

  const { resource, comments } = data
  const rid = parseInt(id)
  const tags = (resource.tags ?? "").split(",").map((t) => t.trim()).filter(Boolean)

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
            {/* 后端 Resource 无 cover 字段，统一显示图标 */}
            <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-lg bg-primary/10">
              <ExternalLink className="h-8 w-8 text-primary" />
            </div>

            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-xl font-bold text-foreground">{resource.title}</h1>
                {resource.catname && (
                  <Badge variant="secondary">{resource.catname}</Badge>
                )}
              </div>
              {/* 后端 Resource 无 author 字段，用 user 子对象或 uid */}
              <p className="mt-1 text-sm text-muted-foreground">
                by{" "}
                {resource.user ? (
                  <Link
                    href={`/user/${resource.user.username}`}
                    className="font-medium text-foreground hover:text-primary"
                  >
                    {resource.user.name || resource.user.username}
                  </Link>
                ) : (
                  <span className="font-medium text-foreground">uid:{resource.uid}</span>
                )}
              </p>
              {resource.content && (
                <MarkdownContent
                  content={resource.content}
                  className="mt-3 text-muted-foreground"
                />
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

          {tags.length > 0 && (
            <div className="mt-5 flex flex-wrap items-center gap-1.5">
              {tags.map((tag) => (
                <Link key={tag} href={`/tag/${encodeURIComponent(tag)}`}>
                  <Badge variant="outline" className="text-xs">
                    {tag}
                  </Badge>
                </Link>
              ))}
            </div>
          )}

          <ContentActions
            className="mt-5 border-t border-border pt-5"
            objid={rid}
            objtype={OBJTYPE_RESOURCE}
            likeCount={resource.likenum || 0}
            initialLiked={Boolean(data.likeflag)}
            initialFavorited={Boolean(data.hadcollect)}
            ownerUid={resource.uid}
            editHref={`/resources/modify/${rid}`}
          />
        </CardContent>
      </Card>

      <CommentSection
        objid={rid}
        objtype={OBJTYPE_RESOURCE}
        comments={comments}
        ownerUid={resource.uid}
      />
    </PageLayout>
  )
}
