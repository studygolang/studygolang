import { notFound } from "next/navigation"
import Link from "next/link"
import { cookies } from "next/headers"
import {
  Eye,
  ExternalLink,
  Github,
  Globe,
  Download,
  MessageSquare,
  ThumbsUp,
} from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { CommentSection } from "@/components/comment-section"
import { ContentActions } from "@/components/content-actions"
import { OBJTYPE_PROJECT, type ProjectDetailData } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { encodeProjectUri, formatNum } from "@/lib/utils"
import { sanitizeHtml } from "@/lib/sanitize"

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

interface ProjectPageProps {
  // catch-all：历史项目的 uri 由名称派生，可能含 /（如 gin-gonic/gin）
  params: Promise<{ uri: string[] }>
}

async function getProjectDetail(uri: string): Promise<ProjectDetailData | null> {
  return fetchAPINullable<ProjectDetailData>(`/projects/${encodeURIComponent(uri)}`, {
    cache: "no-store",
    headers: await authHeaders(),
  })
}

export async function generateMetadata({ params }: ProjectPageProps): Promise<Metadata> {
  const uri = (await params).uri.join("/")
  const data = await getProjectDetail(uri)

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
  const uri = (await params).uri.join("/")
  const data = await getProjectDetail(uri)

  if (!data?.project) {
    return notFound()
  }

  const { project, replies } = data

  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: project.name,
    description: project.desc,
    url: project.home || project.src,
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
        dangerouslySetInnerHTML={{ __html: sanitizeHtml(JSON.stringify(jsonLd)) }}
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
                  href={`/user/${project.username || project.author}`}
                  className="font-medium text-foreground hover:text-primary"
                >
                  {project.username || project.author}
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
              <Eye className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.viewnum)}</span>
              <span className="text-xs text-muted-foreground">浏览</span>
            </div>
            <div className="flex items-center gap-2">
              <MessageSquare className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.cmtnum)}</span>
              <span className="text-xs text-muted-foreground">评论</span>
            </div>
            <div className="flex items-center gap-2">
              <ThumbsUp className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm font-semibold">{formatNum(project.likenum)}</span>
              <span className="text-xs text-muted-foreground">点赞</span>
            </div>
          </div>

          {/* Links */}
          <div className="mt-5 flex flex-wrap gap-3">
            {project.src && (
              <a
                href={project.src}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <Github className="h-4 w-4" />
                源码仓库
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.home && (
              <a
                href={project.home}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <Globe className="h-4 w-4" />
                项目主页
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.doc && (
              <a
                href={project.doc}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-md border border-border px-3 py-2 text-sm font-medium text-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
              >
                <ExternalLink className="h-4 w-4" />
                文档
                <ExternalLink className="h-3 w-3" />
              </a>
            )}
            {project.download && (
              <a
                href={project.download}
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

          <ContentActions
            className="mt-5 border-t border-border pt-5"
            objid={project.id}
            objtype={OBJTYPE_PROJECT}
            likeCount={project.likenum || 0}
            initialLiked={Boolean(data.likeflag)}
            initialFavorited={Boolean(data.hadcollect)}
            ownerUsername={project.username}
            editHref={`/projects/modify/${encodeProjectUri(uri)}`}
          />
        </CardContent>
      </Card>

      <CommentSection
        objid={project.id}
        objtype={OBJTYPE_PROJECT}
        comments={replies}
        ownerUid={project.user?.uid}
      />
    </PageLayout>
  )
}
