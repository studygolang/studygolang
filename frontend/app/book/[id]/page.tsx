import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { BookOpen, ExternalLink, Eye, MessageSquare, ThumbsUp } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import type { Book } from "@/lib/types"

export const revalidate = 60

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const base = process.env.API_BASE_URL || "http://localhost:8088"
    const res = await fetch(`${base}/api/v1${path}`, options)
    if (!res.ok) return null
    const json = await res.json()
    return json.code === 0 ? json.data : null
  } catch {
    return null
  }
}

interface BookDetailPageProps {
  params: Promise<{ id: string }>
}

export async function generateMetadata({ params }: BookDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const book = await fetchAPI<Book>(`/books/${id}`)
  if (!book) {
    return { title: "书籍详情 - Go语言中文网" }
  }
  return {
    title: `${book.name} - Go语言中文网`,
    description: book.desc || `${book.name}，作者：${book.author}`,
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

export default async function BookDetailPage({ params }: BookDetailPageProps) {
  const { id } = await params
  const book = await fetchAPI<Book>(`/books/${id}`)

  if (!book) {
    notFound()
  }

  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "Book",
    name: book.name,
    author: {
      "@type": "Person",
      name: book.author,
    },
    ...(book.translator && {
      translator: {
        "@type": "Person",
        name: book.translator,
      },
    }),
    description: book.desc,
    image: book.cover,
    url: book.url,
    inLanguage: book.lang || "zh-CN",
  }

  return (
    <PageLayout sidebar={false}>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      <PageHeader
        title={book.name}
        breadcrumbs={[
          { label: "Go 图书", href: "/books" },
          { label: book.name },
        ]}
      />

      <div className="mx-auto max-w-4xl">
        <Card>
          <CardContent className="p-6 sm:p-8">
            <div className="flex flex-col gap-8 sm:flex-row">
              {/* Book Cover */}
              <div className="shrink-0">
                {book.cover ? (
                  <div className="overflow-hidden rounded-lg border border-border shadow-md">
                    <img
                      src={book.cover}
                      alt={book.name}
                      className="h-64 w-44 object-cover"
                    />
                  </div>
                ) : (
                  <div
                    className={cn(
                      "flex h-64 w-44 items-center justify-center rounded-lg bg-gradient-to-br shadow-md",
                      getBookColor(book.id)
                    )}
                  >
                    <div className="flex flex-col items-center gap-2 px-3">
                      <BookOpen className="h-12 w-12 text-primary-foreground/80" />
                      <span className="text-center text-sm font-bold leading-snug text-primary-foreground/90 line-clamp-3">
                        {book.name}
                      </span>
                    </div>
                  </div>
                )}
              </div>

              {/* Book Info */}
              <div className="flex min-w-0 flex-1 flex-col gap-4">
                <div>
                  <h1 className="text-2xl font-bold text-foreground">{book.name}</h1>

                  <div className="mt-2 flex flex-wrap gap-2">
                    {book.lang && (
                      <Badge variant="outline" className="text-xs">
                        {book.lang}
                      </Badge>
                    )}
                    {book.buyer_show === 0 && (
                      <Badge className="bg-accent/10 text-xs font-semibold text-accent">
                        免费
                      </Badge>
                    )}
                  </div>
                </div>

                {/* Meta Info */}
                <dl className="grid grid-cols-1 gap-2 text-sm sm:grid-cols-2">
                  <div className="flex gap-2">
                    <dt className="text-muted-foreground">作者</dt>
                    <dd className="font-medium text-foreground">{book.author || "未知"}</dd>
                  </div>
                  {book.translator && (
                    <div className="flex gap-2">
                      <dt className="text-muted-foreground">译者</dt>
                      <dd className="font-medium text-foreground">{book.translator}</dd>
                    </div>
                  )}
                </dl>

                {/* Stats */}
                <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Eye className="h-4 w-4" />
                    {book.viewnum} 次浏览
                  </span>
                  <span className="flex items-center gap-1">
                    <ThumbsUp className="h-4 w-4" />
                    {book.likenum} 人喜欢
                  </span>
                  <span className="flex items-center gap-1">
                    <MessageSquare className="h-4 w-4" />
                    {book.cmtnum} 条评论
                  </span>
                </div>

                {/* Buy / View Link */}
                {book.url && (
                  <div className="mt-2">
                    <a
                      href={book.url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <Button className="gap-2">
                        {book.buyer_show === 0 ? "免费获取" : "购买此书"}
                        <ExternalLink className="h-4 w-4" />
                      </Button>
                    </a>
                  </div>
                )}
              </div>
            </div>

            {/* Description */}
            {book.desc && (
              <div className="mt-8 border-t border-border pt-6">
                <h2 className="mb-3 text-base font-semibold text-foreground">内容简介</h2>
                <p className="text-sm leading-relaxed text-muted-foreground whitespace-pre-line">
                  {book.desc}
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        <div className="mt-4 text-center">
          <Link
            href="/books"
            className="text-sm text-muted-foreground transition-colors hover:text-primary"
          >
            ← 返回书籍列表
          </Link>
        </div>
      </div>
    </PageLayout>
  )
}
