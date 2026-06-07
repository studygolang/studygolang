import { Suspense } from "react"
import Link from "next/link"
import { Tag, Clock, User } from "lucide-react"
import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import type { SearchData } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

// objtype 对应后端 model/comment.go iota 常量
// TypeTopic=0, TypeArticle=1, TypeResource=2, TypeWiki=3, TypeProject=4, TypeBook=5, TypeInterview=6
const typeColorMap: Record<number, string> = {
  0: "bg-blue-100 text-blue-700",
  1: "bg-green-100 text-green-700",
  2: "bg-orange-100 text-orange-700",
  3: "bg-cyan-100 text-cyan-700",
  4: "bg-purple-100 text-purple-700",
  5: "bg-amber-100 text-amber-700",
  6: "bg-pink-100 text-pink-700",
}

const typeLabelMap: Record<number, string> = {
  0: "话题",
  1: "文章",
  2: "资源",
  3: "Wiki",
  4: "项目",
  5: "图书",
  6: "面试题",
}

function buildResultUrl(objtype: number, objid: number): string {
  switch (objtype) {
    case 0: return `/topics/${objid}`
    case 1: return `/articles/${objid}`
    case 2: return `/resources/${objid}`
    case 3: return `/wiki/${objid}`
    case 4: return `/p/${objid}`
    case 5: return `/book/${objid}`
    default: return `/topics/${objid}`
  }
}

interface TagPageProps {
  params: Promise<{ name: string }>
  searchParams: Promise<{ page?: string }>
}

export async function generateMetadata({ params }: TagPageProps): Promise<Metadata> {
  const { name } = await params
  const tagName = decodeURIComponent(name)
  return {
    title: `标签: ${tagName} - Go语言中文网`,
    description: `浏览 Go语言中文网 标签"${tagName}"相关内容`,
  }
}

async function TagResults({ name, page }: { name: string; page: number }) {
  const data = await fetchAPINullable<SearchData>(
    `/tag/${encodeURIComponent(name)}?p=${page}`,
    { cache: "no-store" }
  )

  if (!data || !data.results || data.results.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <Tag className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm font-medium text-foreground">
          暂无标签 <span className="text-primary">&ldquo;{name}&rdquo;</span> 相关内容
        </p>
      </div>
    )
  }

  return (
    <div>
      <p className="mb-4 text-sm text-muted-foreground">
        找到{" "}
        <span className="font-semibold text-foreground">{data.total}</span>{" "}
        个标签{" "}
        <span className="font-semibold text-primary">&ldquo;{name}&rdquo;</span>{" "}
        相关的内容
      </p>

      <div className="space-y-3">
        {data.results.map((result) => {
          const resultUrl = buildResultUrl(result.objtype, result.objid)
          return (
            <Card
              key={result.id}
              className="group transition-all hover:border-primary/20 hover:shadow-sm"
            >
              <CardContent className="p-4 sm:p-5">
                <div className="flex items-start gap-3">
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      {result.objtype && (
                        <span
                          className={`rounded px-1.5 py-0.5 text-[10px] font-semibold ${
                            typeColorMap[result.objtype] ?? "bg-secondary text-secondary-foreground"
                          }`}
                        >
                          {typeLabelMap[result.objtype] ?? String(result.objtype)}
                        </span>
                      )}
                      <a
                        href={resultUrl}
                        className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                      >
                        {result.title}
                      </a>
                    </div>

                    {result.content && (
                      <p className="mt-1.5 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                        {result.content}
                      </p>
                    )}

                    <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1">
                      {result.author && (
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <User className="h-3 w-3" />
                          {result.author}
                        </span>
                      )}
                      {result.pub_time && (
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {result.pub_time}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* 分页 */}
      <div className="mt-6 flex items-center justify-center gap-2">
        {page > 1 && (
          <Link
            href={`/tag/${encodeURIComponent(name)}?page=${page - 1}`}
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
            href={`/tag/${encodeURIComponent(name)}?page=${page + 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </div>
  )
}

export default async function TagPage({ params, searchParams }: TagPageProps) {
  const { name } = await params
  const sp = await searchParams
  const tagName = decodeURIComponent(name)
  const page = Math.max(1, parseInt(sp.page || "1", 10))

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title={`标签: ${tagName}`}
        breadcrumbs={[{ label: "标签", href: "/search" }, { label: tagName }]}
      />
      <Suspense fallback={<div className="space-y-3">{Array.from({ length: 5 }).map((_, i) => (<div key={i} className="h-20 animate-pulse rounded-lg bg-muted" />))}</div>}>
        <TagResults name={tagName} page={page} />
      </Suspense>
    </PageLayout>
  )
}
