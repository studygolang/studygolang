import type { Metadata } from "next"
import Link from "next/link"
import { PenSquare } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicList } from "@/components/topic-list"
import { NodeNavigation } from "@/components/node-navigation"
import { Button } from "@/components/ui/button"
import type { TopicListData, TopicNode } from "@/lib/types"

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

interface NodeTopicsPageProps {
  params: Promise<{ nid: string }>
  searchParams: Promise<{ p?: string }>
}

export async function generateMetadata({ params }: NodeTopicsPageProps): Promise<Metadata> {
  const { nid } = await params
  const nodes = await fetchAPI<TopicNode[]>("/nodes", { cache: "no-store" })
  const node = nodes?.find((n) => String(n.id) === nid)
  const nodeName = node?.name ?? "节点"
  return {
    title: `${nodeName} - 话题 - Go语言中文网`,
    description: node?.intro ?? `浏览 ${nodeName} 节点下的所有话题`,
  }
}

export default async function NodeTopicsPage({ params, searchParams }: NodeTopicsPageProps) {
  const { nid } = await params
  const { p } = await searchParams
  const page = Math.max(1, parseInt(p ?? "1", 10) || 1)

  const [topicsData, nodes] = await Promise.all([
    fetchAPI<TopicListData>(`/topics/node/${nid}?p=${page}`, { cache: "no-store" }),
    fetchAPI<TopicNode[]>("/nodes", { cache: "no-store" }),
  ])

  const nodeList = nodes ?? []
  const currentNode = nodeList.find((n) => String(n.id) === nid)
  const nodeName = currentNode?.name ?? "节点话题"

  // 后端返回字段为 "list"，与 TopicListData.list 对应
  const topics = topicsData?.list ?? []
  const total = topicsData?.total ?? 0
  const hasMore = topicsData?.has_more ?? false

  const pageSize = 20
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : (hasMore ? page + 1 : page)
  const pageNumbers: number[] = []
  const start = Math.max(1, page - 2)
  const end = Math.min(totalPages, page + 2)
  for (let i = start; i <= end; i++) {
    pageNumbers.push(i)
  }

  return (
    <PageLayout sidebarData={{ nodes: nodeList }}>
      <PageHeader
        title={nodeName}
        description={currentNode?.intro ?? `浏览 ${nodeName} 节点下的所有话题`}
        breadcrumbs={[
          { label: "主题讨论", href: "/topics" },
          { label: "话题节点", href: "/nodes" },
          { label: nodeName },
        ]}
        actions={
          // TODO: /topics/new 页面尚未实现，未登录时指向登录页，登录后应跳转至发帖页
          <Link href="/account/login">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <PenSquare className="h-3.5 w-3.5" />
              发布主题
            </Button>
          </Link>
        }
      />

      <div className="rounded-lg border border-border bg-card p-4">
        <div className="mt-2">
          <TopicList topics={topics} />
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            {page > 1 ? (
              <Link
                href={`/topics/node/${nid}?p=${page - 1}`}
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

            {pageNumbers.map((num) => (
              <Link
                key={num}
                href={`/topics/node/${nid}?p=${num}`}
                className={`rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                  num === page
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
                }`}
              >
                {num}
              </Link>
            ))}

            {end < totalPages && (
              <>
                <span className="px-1 text-sm text-muted-foreground">...</span>
                <Link
                  href={`/topics/node/${nid}?p=${totalPages}`}
                  className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                >
                  {totalPages}
                </Link>
              </>
            )}

            {(hasMore || page < totalPages) ? (
              <Link
                href={`/topics/node/${nid}?p=${page + 1}`}
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
        <NodeNavigation nodes={nodeList} />
      </div>
    </PageLayout>
  )
}
