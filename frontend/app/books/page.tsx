import { Suspense } from "react"
import Link from "next/link"
import { BookOpen, ExternalLink } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import type { BookListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "Go 书籍 - Go语言中文网",
  description: "Go语言经典图书推荐，从入门到精通",
}

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

function getBookColor(id: number): string {
  const colors = [
    "from-primary/80 to-primary",
    "from-chart-2/80 to-chart-2",
    "from-chart-3/80 to-chart-3",
    "from-chart-4/80 to-chart-4",
    "from-chart-5/80 to-chart-5",
  ]
  return colors[id % colors.length]
}

async function BookItems({ page }: { page: number }) {
  const data = await fetchAPI<BookListData>(
    `/books?p=${page}`,
    { cache: "no-store" }
  )

  if (!data || !data.books || data.books.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">暂无书籍数据</p>
      </div>
    )
  }

  return (
    <div>
      {/* Books grid */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {data.books.map((book) => (
          <Card
            key={book.id}
            className="group overflow-hidden transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="p-0">
              {/* Book cover */}
              {book.cover ? (
                <div className="relative h-44 overflow-hidden bg-muted">
                  <img
                    src={book.cover}
                    alt={book.name}
                    className="h-full w-full object-cover"
                  />
                </div>
              ) : (
                <div
                  className={cn(
                    "flex h-44 items-center justify-center bg-gradient-to-br",
                    getBookColor(book.id)
                  )}
                >
                  <div className="flex flex-col items-center gap-2">
                    <BookOpen className="h-10 w-10 text-primary-foreground/80" />
                    <span className="max-w-[120px] text-center text-sm font-bold text-primary-foreground/90 line-clamp-2">
                      {book.name}
                    </span>
                  </div>
                </div>
              )}

              {/* Info */}
              <div className="p-4">
                <div className="flex items-start justify-between gap-2">
                  <h3 className="text-sm font-bold leading-snug text-foreground group-hover:text-primary">
                    {book.name}
                  </h3>
                  {book.buyer_show === 0 && (
                    <Badge className="shrink-0 bg-accent/10 text-[10px] font-semibold text-accent">
                      免费
                    </Badge>
                  )}
                </div>

                <p className="mt-1 text-xs text-muted-foreground">
                  {book.author}
                  {book.translator && ` / 译：${book.translator}`}
                </p>

                {book.lang && (
                  <Badge variant="outline" className="mt-1.5 text-[10px]">
                    {book.lang}
                  </Badge>
                )}

                {/* Description */}
                {book.desc && (
                  <p className="mt-2 line-clamp-2 text-xs leading-relaxed text-muted-foreground">
                    {book.desc}
                  </p>
                )}

                {/* Link */}
                <div className="mt-3 flex items-center justify-end border-t border-border pt-3">
                  {book.url ? (
                    <a
                      href={book.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs font-medium text-primary transition-colors hover:text-primary/80"
                      aria-label={"查看 " + book.name}
                    >
                      查看
                      <ExternalLink className="h-3 w-3" />
                    </a>
                  ) : (
                    // 书籍详情页路由是 /book/[id]（注意是单数），非 /books/[id]
                    <Link
                      href={`/book/${book.id}`}
                      className="flex items-center gap-1 text-xs font-medium text-primary transition-colors hover:text-primary/80"
                    >
                      详情
                      <ExternalLink className="h-3 w-3" />
                    </Link>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Pagination */}
      <div className="mt-8 flex items-center justify-center gap-2">
        {page > 1 && (
          <Link
            href={`/books?page=${page - 1}`}
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
            href={`/books?page=${page + 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </div>
  )
}

interface BooksPageProps {
  searchParams: Promise<{ page?: string }>
}

export default async function BooksPage({ searchParams }: BooksPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.page || "1", 10))

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go 图书"
        description="精选 Go 语言经典图书，助你系统学习"
        breadcrumbs={[{ label: "Go 图书" }]}
      />
      <Suspense
        fallback={
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-72 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <BookItems page={page} />
      </Suspense>
    </PageLayout>
  )
}
