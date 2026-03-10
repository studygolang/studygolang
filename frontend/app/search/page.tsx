import { Suspense } from "react"
import Link from "next/link"
import { Search, Clock, User } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { SearchBox } from "./search-box"
import type { SearchData } from "@/lib/types"

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

const typeColorMap: Record<string, string> = {
  topic: "bg-blue-100 text-blue-700",
  article: "bg-green-100 text-green-700",
  project: "bg-purple-100 text-purple-700",
  resource: "bg-orange-100 text-orange-700",
  book: "bg-pink-100 text-pink-700",
}

const typeLabelMap: Record<string, string> = {
  topic: "话题",
  article: "文章",
  project: "项目",
  resource: "资源",
  book: "书籍",
}

interface SearchPageProps {
  searchParams: Promise<{ q?: string; page?: string }>
}

export async function generateMetadata({ searchParams }: SearchPageProps): Promise<Metadata> {
  const params = await searchParams
  const q = params.q?.trim() || ""
  if (!q) {
    return {
      title: "搜索 - Go语言中文网",
      description: "搜索 Go 语言相关内容",
    }
  }
  return {
    title: `搜索"${q}" - Go语言中文网`,
    description: `在 Go语言中文网搜索"${q}"的相关内容`,
  }
}

async function SearchResults({ q, page }: { q: string; page: number }) {
  if (!q) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <Search className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">输入关键词开始搜索</p>
      </div>
    )
  }

  const data = await fetchAPI<SearchData>(
    `/search?q=${encodeURIComponent(q)}&p=${page}`,
    { cache: "no-store" }
  )

  if (!data || !data.results || data.results.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <Search className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm font-medium text-foreground">
          没有找到与 <span className="text-primary">&ldquo;{q}&rdquo;</span> 相关的内容
        </p>
        <p className="mt-1 text-sm text-muted-foreground">
          试试换个关键词，或者减少搜索条件
        </p>
      </div>
    )
  }

  return (
    <div>
      {/* Result count */}
      <p className="mb-4 text-sm text-muted-foreground">
        找到{" "}
        <span className="font-semibold text-foreground">{data.total}</span>{" "}
        个与{" "}
        <span className="font-semibold text-primary">&ldquo;{data.keyword}&rdquo;</span>{" "}
        相关的结果
      </p>

      {/* Results */}
      <div className="space-y-3">
        {data.results.map((result) => (
          <Card
            key={result.id}
            className="group transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="p-4 sm:p-5">
              <div className="flex items-start gap-3">
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    {result.type && (
                      <span
                        className={`rounded px-1.5 py-0.5 text-[10px] font-semibold ${
                          typeColorMap[result.type] ?? "bg-secondary text-secondary-foreground"
                        }`}
                      >
                        {typeLabelMap[result.type] ?? result.type}
                      </span>
                    )}
                    <a
                      href={result.url}
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
                    {result.ctime && (
                      <span className="flex items-center gap-1 text-xs text-muted-foreground">
                        <Clock className="h-3 w-3" />
                        {result.ctime}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Pagination */}
      <div className="mt-6 flex items-center justify-center gap-2">
        {page > 1 && (
          <Link
            href={`/search?q=${encodeURIComponent(q)}&page=${page - 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            上一页
          </Link>
        )}
        <span className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground">
          {page}
        </span>
        {data.results.length >= 20 && (
          <Link
            href={`/search?q=${encodeURIComponent(q)}&page=${page + 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </div>
  )
}

export default async function SearchPage({ searchParams }: SearchPageProps) {
  const params = await searchParams
  const q = params.q?.trim() || ""
  const page = Math.max(1, parseInt(params.page || "1", 10))

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="搜索"
        breadcrumbs={[{ label: "搜索" }]}
      />

      {/* Search box */}
      <div className="mb-6 flex justify-center">
        <Suspense
          fallback={
            <div className="h-12 w-full max-w-2xl animate-pulse rounded-md bg-muted" />
          }
        >
          <SearchBox defaultValue={q} />
        </Suspense>
      </div>

      {/* Results */}
      <Suspense
        key={`${q}-${page}`}
        fallback={
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-24 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <SearchResults q={q} page={page} />
      </Suspense>
    </PageLayout>
  )
}
