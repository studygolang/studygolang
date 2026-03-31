import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicList } from "@/components/topic-list"
import { NodeNavigation } from "@/components/node-navigation"
import Link from "next/link"
import { PenSquare } from "lucide-react"
import { AuthLink } from "@/components/auth-link"
import type { TopicListData, TopicNode } from "@/lib/types"

export const metadata: Metadata = {
  title: "主题讨论 - Go语言中文网",
  description: "Go语言中文社区主题讨论，分享技术经验，交流开发心得",
}

import { fetchAPI } from "@/lib/api"

async function getTopicsData(tab: string, page: number) {
  const [topicsResult, nodesResult] = await Promise.allSettled([
    fetchAPI<TopicListData>(
      `/topics?tab=${tab}&p=${page}`,
      { cache: 'no-store' }
    ),
    fetchAPI<TopicNode[]>('/nodes', { cache: 'no-store' }),
  ])

  return {
    topicsData: topicsResult.status === 'fulfilled'
      ? topicsResult.value
      : { list: [], page, total: 0, has_more: false, tab, tab_list: [] },
    nodes: nodesResult.status === 'fulfilled' ? nodesResult.value : [],
  }
}

const TAB_LABELS: Record<string, string> = {
  all: '全部',
  go: 'Go语言',
  ask: '问与答',
  share: '分享',
  news: '新闻',
}

interface TopicsPageProps {
  searchParams: Promise<{ tab?: string; p?: string }>
}

export default async function TopicsPage({ searchParams }: TopicsPageProps) {
  const { tab: tabParam, p: pageStr } = await searchParams
  const tab = tabParam ?? 'all'
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const { topicsData, nodes } = await getTopicsData(tab, page)
  // 后端返回字段为 "list"，与 TopicListData.list 对应
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

  const tabs = [
    { id: 'all', label: '全部' },
    { id: 'go', label: 'Go语言' },
    { id: 'ask', label: '问与答' },
    { id: 'share', label: '分享' },
    { id: 'news', label: '新闻' },
  ]

  return (
    <PageLayout>
      <PageHeader
        title="主题讨论"
        description="分享你的技术见解，和 Gopher 们一起交流成长"
        breadcrumbs={[{ label: "主题讨论" }]}
        actions={
          
          <AuthLink href="/publish" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
            <PenSquare className="h-3.5 w-3.5" />
            发布主题
          </AuthLink>
        }
      />

      <div className="rounded-lg border border-border bg-card p-4">
        {/* Tab Filter */}
        <div className="space-y-3">
          <div className="flex items-center gap-1 border-b border-border">
            {tabs.map((t) => (
              <Link
                key={t.id}
                href={`/topics?tab=${t.id}`}
                className={`relative px-4 py-2.5 text-sm font-medium transition-colors ${
                  tab === t.id
                    ? "text-primary"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {t.label}
                {tab === t.id && (
                  <span className="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary" />
                )}
              </Link>
            ))}
          </div>
        </div>

        <div className="mt-4">
          <TopicList topics={topics} />
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            {page > 1 ? (
              <Link
                href={`/topics?tab=${tab}&p=${page - 1}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                {"上一页"}
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                {"上一页"}
              </button>
            )}
            {pageNumbers.map((p) => (
              <Link
                key={p}
                href={`/topics?tab=${tab}&p=${p}`}
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
                  href={`/topics?tab=${tab}&p=${totalPages}`}
                  className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                >
                  {totalPages}
                </Link>
              </>
            )}
            {(hasMore || page < totalPages) ? (
              <Link
                href={`/topics?tab=${tab}&p=${page + 1}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                {"下一页"}
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                {"下一页"}
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
