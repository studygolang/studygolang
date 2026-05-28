import { Suspense } from "react"
import Link from "next/link"
import { Plus, ExternalLink, ThumbsUp, MessageSquare, Clock } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { ResourceListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "Go 资源 - Go语言中文网",
  description: "Go语言学习资源大全，教程、工具、视频、文档一站式索引",
  openGraph: {
    title: "Go 资源 - Go语言中文网",
    description: "Go语言学习资源大全，教程、工具、视频、文档一站式索引",
    type: "website",
  },
}

import { fetchAPINullable } from "@/lib/api"
import { formatNum } from "@/lib/utils"

async function ResourceItems({ page, catid = 0, sort = "new" }: { page: number; catid?: number; sort?: string }) {
  // 支持 catid 分类过滤和 sort 排序，0 表示全部，sort 默认最新
  const params = new URLSearchParams({ p: String(page), sort })
  if (catid > 0) {
    params.set("catid", String(catid))
  }
  const qs = `/resources?${params.toString()}`
  const data = await fetchAPINullable<ResourceListData>(qs, { cache: "no-store" })

  if (!data || !data.resources || data.resources.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <ExternalLink className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">暂无资源数据</p>
      </div>
    )
  }

  return (
    <div className="min-w-0 flex-1 space-y-3">
      {/* Result header */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {/* 后端资源列表 API 不返回 total，只返回 has_more */}
          共{" "}
          <span className="font-semibold text-foreground">{data.resources.length}</span>{" "}
          个资源（本页）
        </p>
      </div>

      {/* Cards */}
      {data.resources.map((resource) => (
        <Card
          key={resource.id}
          className="group transition-all hover:border-primary/20 hover:shadow-sm"
        >
          <CardContent className="p-4 sm:p-5">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                {/* Title */}
                <div className="flex items-start gap-2">
                  {resource.catname && (
                    <Badge className="mt-0.5 shrink-0 bg-primary/10 text-[10px] font-semibold text-primary">
                      {resource.catname}
                    </Badge>
                  )}
                  <Link
                    href={`/resources/${resource.id}`}
                    className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                  >
                    {resource.title}
                  </Link>
                </div>

                {/* Description：后端 Resource 用 content 字段存描述，无 desc 字段 */}
                {resource.content && (
                  <p className="mt-1.5 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                    {resource.content}
                  </p>
                )}

                {/* Meta */}
                <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                  {/* 后端 Resource 无 author 字段，用 user 子对象或 uid */}
                  <span className="text-xs text-muted-foreground">
                    {resource.user?.name || resource.user?.username || `uid:${resource.uid}`}
                  </span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Clock className="h-3 w-3" />
                    {resource.ctime}
                  </span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <ThumbsUp className="h-3 w-3" />
                    {formatNum(resource.likenum)}
                  </span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <MessageSquare className="h-3 w-3" />
                    {resource.cmtnum}
                  </span>
                </div>
              </div>

              {/* External link button */}
              {resource.url && (
                <a
                  href={resource.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
                  aria-label={"访问 " + resource.title}
                >
                  <ExternalLink className="h-4 w-4" />
                </a>
              )}
            </div>
          </CardContent>
        </Card>
      ))}

      {/* Pagination */}
      <div className="flex items-center justify-center gap-2 pt-4">
        {page > 1 && (
          <Link
            href={`/resources?p=${page - 1}&sort=${sort}${catid > 0 ? `&catid=${catid}` : ""}`}
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
            href={`/resources?p=${page + 1}&sort=${sort}${catid > 0 ? `&catid=${catid}` : ""}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </div>
  )
}

interface ResourcesPageProps {
  searchParams: Promise<{ p?: string; catid?: string; sort?: string }>
}

export default async function ResourcesPage({ searchParams }: ResourcesPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.p || "1", 10))
  const catid = params.catid ? parseInt(params.catid, 10) : 0
  const sort = params.sort || "new"

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="资源索引"
        description="精选 Go 语言学习资源，教程、工具、视频、文档一网打尽"
        breadcrumbs={[{ label: "资源索引" }]}
        actions={
          <Link href="/publish?type=resource">
            <Button size="sm" className="gap-1.5">
              <Plus className="h-3.5 w-3.5" />
              分享资源
            </Button>
          </Link>
        }
      />
      <Suspense
        fallback={
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-24 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <ResourceItems page={page} catid={catid} sort={sort} />
      </Suspense>
    </PageLayout>
  )
}
