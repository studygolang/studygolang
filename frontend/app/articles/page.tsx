import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ArticleList } from "@/components/article-list"
import { ArticleFilter } from "@/components/article-filter"
import { PenSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import type { ArticleListData } from "@/lib/types"

export const metadata = {
  title: "Go语言文章 - Go语言中文网",
  description: "Go语言高质量技术文章，涵盖基础教程、实战经验、源码解析等，助力 Gopher 成长",
}

async function fetchFromAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const base = process.env.API_BASE_URL || 'http://localhost:8090'
  const res = await fetch(`${base}/api/v1${path}`, options)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  return json.data
}

async function getArticles(page: number) {
  try {
    return await fetchFromAPI<ArticleListData>(
      `/articles?p=${page}`,
      { cache: 'no-store' }
    )
  } catch {
    // 后端返回字段名为 list（非 articles）
    return { list: [], page: page, total: 0, has_more: false }
  }
}

interface ArticlesPageProps {
  searchParams: Promise<{ p?: string }>
}

export default async function ArticlesPage({ searchParams }: ArticlesPageProps) {
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const data = await getArticles(page)
  // 后端文章列表返回字段名为 list（非 articles）
  const articles = data.list ?? []
  const total = data.total ?? 0
  const hasMore = data.has_more ?? false

  // 计算总页数（假设每页20条）
  const pageSize = 20
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : (hasMore ? page + 1 : page)

  // 生成页码列表
  const pageNumbers: number[] = []
  const start = Math.max(1, page - 2)
  const end = Math.min(totalPages, page + 2)
  for (let i = start; i <= end; i++) {
    pageNumbers.push(i)
  }

  return (
    <PageLayout>
      <PageHeader
        title="技术文章"
        description="精选 Go 语言技术文章，深度学习与实践"
        breadcrumbs={[{ label: "技术文章" }]}
        actions={
          <Link href="/publish">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <PenSquare className="h-3.5 w-3.5" />
              {"投稿"}
            </Button>
          </Link>
        }
      />

      <ArticleFilter />
      <div className="mt-4">
        <ArticleList articles={articles} />
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          {page > 1 ? (
            <Link
              href={`/articles?p=${page - 1}`}
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
              href={`/articles?p=${p}`}
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
                href={`/articles?p=${totalPages}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                {totalPages}
              </Link>
            </>
          )}
          {(hasMore || page < totalPages) ? (
            <Link
              href={`/articles?p=${page + 1}`}
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
    </PageLayout>
  )
}
