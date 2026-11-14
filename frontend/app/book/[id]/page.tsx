import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { cookies } from "next/headers"
import { BookOpen, ExternalLink, Eye, MessageSquare, ThumbsUp } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { CommentSection } from "@/components/comment-section"
import { ContentActions } from "@/components/content-actions"
import { SafeHtml } from "@/components/safe-html"
import { cn } from "@/lib/utils"
import { OBJTYPE_BOOK, type BookDetailData } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"

// 点赞/收藏状态与评论 liked 都是当前用户维度的数据，不能进 ISR 缓存
export const dynamic = "force-dynamic"

interface BookDetailPageProps {
  params: Promise<{ id: string }>
}

async function getBookDetail(id: string): Promise<BookDetailData | null> {
  const store = await cookies()
  const cookie = store
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")
  return fetchAPINullable<BookDetailData>(`/books/${id}`, {
    cache: "no-store",
    headers: cookie ? { Cookie: cookie } : {},
  })
}

export async function generateMetadata({ params }: BookDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const data = await getBookDetail(id)
  const book = data?.book
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
  const data = await getBookDetail(id)
  const book = data?.book

  if (!book) {
    return notFound()
  }

  const bid = parseInt(id)
  const tags = (book.tags ?? "").split(",").map((t) => t.trim()).filter(Boolean)

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
    url: book.online_url || book.buy_url || book.download_url,
    inLanguage: book.lang === 1 ? "en" : "zh-CN",
  }

  return (
    <PageLayout sidebar={false}>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: sanitizeHtml(JSON.stringify(jsonLd)) }}
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
                    <Badge variant="outline" className="text-xs">
                      {book.lang === 1 ? "英文版" : "中文版"}
                    </Badge>
                    {book.is_free && (
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

                {/* 在线阅读 / 下载 / 去购买，对齐旧版 template/books/detail.html 的三个入口 */}
                {(book.online_url || book.download_url || book.buy_url) && (
                  <div className="mt-2 flex flex-wrap gap-2">
                    {book.online_url && (
                      <a href={book.online_url} target="_blank" rel="noopener noreferrer">
                        <Button className="gap-2">
                          在线阅读
                          <ExternalLink className="h-4 w-4" />
                        </Button>
                      </a>
                    )}
                    {book.download_url && (
                      <a href={book.download_url} target="_blank" rel="noopener noreferrer">
                        <Button variant="outline" className="gap-2">
                          下载
                          <ExternalLink className="h-4 w-4" />
                        </Button>
                      </a>
                    )}
                    {book.buy_url && (
                      <a href={book.buy_url} target="_blank" rel="noopener noreferrer">
                        <Button variant="outline" className="gap-2">
                          去购买
                          <ExternalLink className="h-4 w-4" />
                        </Button>
                      </a>
                    )}
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

            {/* Catalogue：后端存的是富文本 HTML（旧版 noescape 输出） */}
            {book.catalogue && (
              <div className="mt-8 border-t border-border pt-6">
                <h2 className="mb-3 text-base font-semibold text-foreground">目录</h2>
                <SafeHtml
                  html={book.catalogue}
                  className="prose prose-sm max-w-none text-muted-foreground"
                />
              </div>
            )}

            {tags.length > 0 && (
              <div className="mt-8 flex flex-wrap items-center gap-1.5 border-t border-border pt-6">
                {tags.map((tag) => (
                  <Link key={tag} href={`/tag/${encodeURIComponent(tag)}`}>
                    <Badge variant="outline" className="text-xs">
                      {tag}
                    </Badge>
                  </Link>
                ))}
              </div>
            )}

            <ContentActions
              className="mt-8 border-t border-border pt-6"
              objid={bid}
              objtype={OBJTYPE_BOOK}
              likeCount={book.likenum || 0}
              initialLiked={Boolean(data.likeflag)}
              initialFavorited={Boolean(data.hadcollect)}
              ownerUid={book.uid}
              editHref={`/books/modify/${bid}`}
            />
          </CardContent>
        </Card>

        <div className="mt-4">
          <CommentSection
            objid={bid}
            objtype={OBJTYPE_BOOK}
            comments={data.replies}
            ownerUid={book.uid}
          />
        </div>

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
