import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { BookOpen, ExternalLink, Clock, Eye } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { Reading } from "@/lib/types"

import { fetchAPI } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"

// 晨读详情数据结构（后端返回 { reading: Reading, comments: Comment[] }）
interface ReadingDetailData {
  reading: Reading
  comments?: Array<{
    id: number
    uid: number
    name: string
    avatar: string
    content: string
    ctime: string
  }>
}

async function getReadingDetail(id: string): Promise<ReadingDetailData | null> {
  try {
    return await fetchAPI<ReadingDetailData>(
      `/readings/${id}`,
      { next: { revalidate: 60 } }
    )
  } catch {
    return null
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>
}): Promise<Metadata> {
  const { id } = await params
  const data = await getReadingDetail(id)

  if (!data?.reading) {
    return {
      title: "晨读不存在 - Go语言中文网",
    }
  }

  const { reading } = data
  const description = reading.content?.slice(0, 150) || "Go语言技术晨读"

  return {
    title: `${reading.content.slice(0, 50)} - Go语言中文网技术晨读`,
    description,
    openGraph: {
      title: reading.content.slice(0, 50),
      description,
      type: "article",
      publishedTime: reading.ctime,
    },
  }
}

const rtypeLabelMap: Record<number, string> = {
  1: "文章",
  2: "视频",
  3: "工具",
  4: "资讯",
}

export default async function ReadingDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const data = await getReadingDetail(id)

  const reading = data?.reading
  const comments = data?.comments ?? []

  // JSON-LD 结构化数据
  const jsonLd = reading
    ? {
        "@context": "https://schema.org",
        "@type": "Article",
        headline: reading.content.slice(0, 100),
        description: reading.content.slice(0, 150),
        author: {
          "@type": "Person",
          name: reading.username || "匿名",
        },
        datePublished: reading.ctime,
        publisher: {
          "@type": "Organization",
          name: "Go语言中文网",
          url: "https://studygolang.com",
        },
      }
    : null

  if (!reading) {
    return (
      <PageLayout>
        <PageHeader
          title="晨读不存在"
          breadcrumbs={[
            { label: "技术晨读", href: "/readings" },
            { label: `#${id}` },
          ]}
        />
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">该晨读内容可能已被删除</p>
        </div>
      </PageLayout>
    )
  }

  return (
    <PageLayout>
      {jsonLd && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: sanitizeHtml(JSON.stringify(jsonLd)) }}
        />
      )}
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "技术晨读", href: "/readings" },
          { label: reading.rdate || reading.ctime.slice(0, 10) },
        ]}
      />

      <div className="space-y-6">
        {/* 晨读内容卡片 */}
        <Card className="overflow-hidden">
          <CardContent className="p-0">
            {/* 顶部图标区域 */}
            <div className="flex h-40 items-center justify-center bg-gradient-to-br from-primary/10 to-primary/5">
              <BookOpen className="h-10 w-10 text-primary/40" />
            </div>

            {/* 晨读详情 */}
            <div className="p-6">
              {/* 标签 */}
              <div className="mb-3 flex flex-wrap gap-1.5">
                {reading.rtype && rtypeLabelMap[reading.rtype] && (
                  <Badge variant="secondary" className="text-xs">
                    {rtypeLabelMap[reading.rtype]}
                  </Badge>
                )}
                {reading.rdate && (
                  <Badge variant="outline" className="text-xs">
                    {reading.rdate}
                  </Badge>
                )}
                {reading.inner === 1 && (
                  <Badge variant="default" className="text-xs">
                    站内内容
                  </Badge>
                )}
              </div>

              {/* 标题/描述 */}
              <h1 className="text-xl font-semibold leading-snug text-foreground mb-4">
                {reading.content}
              </h1>

              {/* 主链接 */}
              <a
                href={reading.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
              >
                阅读原文
                <ExternalLink className="h-4 w-4" />
              </a>

              {/* 多链接支持 */}
              {reading.urls && reading.urls.length > 0 && (
                <div className="mt-4 space-y-2">
                  <p className="text-sm font-medium text-foreground">相关链接：</p>
                  <div className="flex flex-wrap gap-2">
                    {reading.urls.map((url, index) => (
                      <a
                        key={index}
                        href={url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1.5 rounded-md bg-secondary px-3 py-1.5 text-xs font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
                      >
                        链接 {index + 1}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    ))}
                  </div>
                </div>
              )}

              {/* 元信息 */}
              <div className="mt-6 flex items-center gap-4 border-t border-border pt-4 text-sm text-muted-foreground">
                <span>发布者：{reading.username || "匿名"}</span>
                <span className="flex items-center gap-1">
                  <Clock className="h-3.5 w-3.5" />
                  {reading.ctime}
                </span>
                <span className="flex items-center gap-1">
                  <Eye className="h-3.5 w-3.5" />
                  {reading.clicknum} 次浏览
                </span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* 评论区（如果有评论） */}
        {comments.length > 0 && (
          <Card>
            <CardContent className="p-6">
              <h2 className="text-lg font-semibold mb-4">评论 ({comments.length})</h2>
              <div className="space-y-4">
                {comments.map((comment) => (
                  <div key={comment.id} className="border-b border-border pb-4 last:border-0">
                    <div className="flex items-start gap-3">
                      <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-semibold text-primary">
                        {comment.name?.charAt(0).toUpperCase() || "?"}
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-medium">{comment.name}</span>
                          <span className="text-xs text-muted-foreground">{comment.ctime}</span>
                        </div>
                        <p className="text-sm text-foreground">{comment.content}</p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </PageLayout>
  )
}
