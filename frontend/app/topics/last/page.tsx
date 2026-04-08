import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicList } from "@/components/topic-list"
import { NodeNavigation } from "@/components/node-navigation"
import Link from "next/link"
import { PenSquare } from "lucide-react"
import { AuthLink } from "@/components/auth-link"
import type { TopicListData, TopicNode } from "@/lib/types"
import { fetchAPI } from "@/lib/api"

export const metadata: Metadata = {
  title: "最新话题 - Go语言中文网",
  description: "查看最新的 Go 语言讨论话题，了解社区最新动态和技术交流",
}

async function getLastTopicsData(page: number) {
  const [topicsResult, nodesResult] = await Promise.allSettled([
    fetchAPI<TopicListData>(`/topics/last?p=${page}`, { cache: 'no-store' }),
    fetchAPI<TopicNode[]>('/nodes', { cache: 'no-store' }),
  ])

  return {
    topicsData: topicsResult.status === 'fulfilled'
      ? topicsResult.value
      : { list: [], page, total: 0, has_more: false, tab: 'all', tab_list: [] },
    nodes: nodesResult.status === 'fulfilled' ? nodesResult.value : [],
  }
}

interface LastTopicsPageProps {
  searchParams: Promise<{ p?: string }>
}

export default async function LastTopicsPage({ searchParams }: LastTopicsPageProps) {
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const { topicsData, nodes } = await getLastTopicsData(page)
  const topics = topicsData.list ?? []
  const total = topicsData.total ?? 0
  const hasMore = topicsData.has_more ?? false

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
        title="最新话题"
        description="查看社区最新发布的讨论话题"
        breadcrumbs={[{ label: "主题讨论", href: "/topics" }, { label: "最新话题" }]}
        actions={
          <AuthLink href="/publish" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
            <PenSquare className="h-3.5 w-3.5" />
            发布主题
          </AuthLink>
        }
      />

      <div className="rounded-lg border border-border bg-card p-4">
        <div className="mt-4">
          <TopicList topics={topics} />
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            {page > 1 ? (
              <Link
                href={`/topics/last?p=${page - 1}`}
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
                href={`/topics/last?p=${p}`}
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
                  href={`/topics/last?p=${totalPages}`}
                  className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                >
                  {totalPages}
                </Link>
              </>
            )}
            {(hasMore || page < totalPages) ? (
              <Link
                href={`/topics/last?p=${page + 1}`}
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

      {/* Node Navigation - mobile only */}
      <div className="mt-6 lg:hidden">
        <NodeNavigation nodes={nodes} />
      </div>
    </PageLayout>
  )
}
