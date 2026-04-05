import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicList } from "@/components/topic-list"
import Link from "next/link"
import { PenSquare, Hash } from "lucide-react"
import { AuthLink } from "@/components/auth-link"
import type { Topic, TopicNode } from "@/lib/types"

import { fetchAPI } from "@/lib/api"

interface NodeTopicsData {
  list: Topic[]
  total: number
  page: number
  has_more: boolean
}

async function getNodeDetail(node: string, page?: number): Promise<{ nodeInfo: TopicNode | null; topicsData: NodeTopicsData | null }> {
  try {
    const data = await fetchAPI<{
      node: TopicNode
      list: Topic[]
      total: number
      page: number
      has_more: boolean
    }>(`/go/${node}${page && page > 1 ? `?p=${page}` : ''}`, { cache: 'no-store' })

    return {
      nodeInfo: data.node,
      topicsData: {
        list: data.list ?? [],
        total: data.total ?? 0,
        page: data.page ?? 1,
        has_more: data.has_more ?? false,
      },
    }
  } catch {
    return { nodeInfo: null, topicsData: null }
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ node: string }>
}): Promise<Metadata> {
  const { node } = await params
  const { nodeInfo } = await getNodeDetail(node)

  if (!nodeInfo) {
    return {
      title: "节点不存在 - Go语言中文网",
    }
  }

  const description = nodeInfo.intro || `${nodeInfo.name}节点的技术讨论`

  return {
    title: `${nodeInfo.name} - Go语言中文网`,
    description,
    openGraph: {
      title: nodeInfo.name,
      description,
      type: "website",
    },
  }
}

interface GoNodePageProps {
  params: Promise<{ node: string }>
  searchParams: Promise<{ p?: string }>
}

export default async function GoNodePage({ params, searchParams }: GoNodePageProps) {
  const { node } = await params
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const { nodeInfo, topicsData } = await getNodeDetail(node, page)

  if (!nodeInfo) {
    return (
      <PageLayout>
        <PageHeader
          title="节点不存在"
          breadcrumbs={[
            { label: "主题讨论", href: "/topics" },
            { label: `/${node}` },
          ]}
        />
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <Hash className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">该节点可能已被删除或不存在</p>
        </div>
      </PageLayout>
    )
  }

  const topics = topicsData?.list ?? []
  const total = topicsData?.total ?? 0
  const hasMore = topicsData?.has_more ?? false

  // 计算分页
  const pageSize = 20
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : (hasMore ? page + 1 : page)
  const pageNumbers: number[] = []
  const start = Math.max(1, page - 2)
  const end = Math.min(totalPages, page + 2)
  for (let i = start; i <= end; i++) {
    pageNumbers.push(i)
  }

  return (
    <PageLayout>
      <PageHeader
        title={nodeInfo.name}
        description={nodeInfo.intro || `${nodeInfo.name}节点的技术讨论`}
        breadcrumbs={[
          { label: "主题讨论", href: "/topics" },
          { label: nodeInfo.name },
        ]}
        actions={
          <AuthLink href="/publish" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
            <PenSquare className="h-3.5 w-3.5" />
            发布主题
          </AuthLink>
        }
      />

      {/* 节点信息卡片 */}
      {nodeInfo.logo && (
        <div className="mb-6 rounded-lg border border-border bg-card p-4">
          <div className="flex items-center gap-4">
            <img
              src={nodeInfo.logo}
              alt={nodeInfo.name}
              className="h-16 w-16 rounded-lg object-cover"
            />
            <div>
              <h2 className="text-lg font-semibold">{nodeInfo.name}</h2>
              {nodeInfo.intro && (
                <p className="text-sm text-muted-foreground">{nodeInfo.intro}</p>
              )}
            </div>
          </div>
        </div>
      )}

      <div className="rounded-lg border border-border bg-card p-4">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-semibold">话题列表</h2>
          <span className="text-sm text-muted-foreground">共 {total} 个话题</span>
        </div>

        <TopicList topics={topics} />

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            {page > 1 ? (
              <Link
                href={`/go/${node}?p=${page - 1}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                上一页
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                上一页
              </button>
            )}
            {pageNumbers.map((p) => (
              <Link
                key={p}
                href={`/go/${node}?p=${p}`}
                className={`rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                  p === page
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
                }`}
              >
                {p}
              </Link>
            ))}
            {end < totalPages && (
              <>
                <span className="px-1 text-sm text-muted-foreground">...</span>
                <Link
                  href={`/go/${node}?p=${totalPages}`}
                  className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                >
                  {totalPages}
                </Link>
              </>
            )}
            {(hasMore || page < totalPages) ? (
              <Link
                href={`/go/${node}?p=${page + 1}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                下一页
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                下一页
              </button>
            )}
          </div>
        )}
      </div>
    </PageLayout>
  )
}
