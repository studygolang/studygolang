import { Suspense } from "react"
import Link from "next/link"
import { BookOpen, Eye, ThumbsUp, MessageSquare, ChevronLeft, ChevronRight } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { InterviewListData } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"

export const metadata: Metadata = {
  title: "Go 面试题 - Go语言中文网",
  description: "精选 Go 语言面试题，助你轻松通过 Go 岗位面试",
}

const levelLabels = ["初级", "中级", "高级"]
const levelColors = [
  "bg-green-100 text-green-700",
  "bg-yellow-100 text-yellow-700",
  "bg-red-100 text-red-700",
]

async function InterviewItems({ page }: { page: number }) {
  const data = await fetchAPINullable<InterviewListData>(
    `/interviews?p=${page}`,
    { cache: "no-store" }
  )
  const questions = data?.questions ?? []
  const totalPages = data?.total_pages ?? 1

  if (questions.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">暂无面试题</p>
      </div>
    )
  }

  return (
    <>
      <div className="space-y-2">
        {questions.map((q) => (
          <Card key={q.id} className="group transition-all hover:border-primary/20">
            <CardContent className="p-4">
              <div className="flex items-start gap-3">
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    {q.level >= 0 && q.level <= 2 && (
                      <Badge className={`shrink-0 text-[10px] font-semibold ${levelColors[q.level]}`}>
                        {levelLabels[q.level]}
                      </Badge>
                    )}
                    <Link
                      href={`/interview/question/${q.show_sn}`}
                      className="text-sm font-medium text-foreground transition-colors group-hover:text-primary line-clamp-2"
                    >
                      <span dangerouslySetInnerHTML={{ __html: sanitizeHtml(q.question) }} />
                    </Link>
                  </div>
                  <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <Eye className="h-3 w-3" />
                      {q.viewnum}
                    </span>
                    <span className="flex items-center gap-1">
                      <ThumbsUp className="h-3 w-3" />
                      {q.likenum}
                    </span>
                    <span className="flex items-center gap-1">
                      <MessageSquare className="h-3 w-3" />
                      {q.cmtnum}
                    </span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          {page > 1 && (
            <Link
              href={`/interview?page=${page - 1}`}
              className="flex items-center gap-1 rounded-md border border-border px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted"
            >
              <ChevronLeft className="h-4 w-4" />
              上一页
            </Link>
          )}
          <span className="text-sm text-muted-foreground">
            第 {page} / {totalPages} 页
          </span>
          {page < totalPages && (
            <Link
              href={`/interview?page=${page + 1}`}
              className="flex items-center gap-1 rounded-md border border-border px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted"
            >
              下一页
              <ChevronRight className="h-4 w-4" />
            </Link>
          )}
        </div>
      )}
    </>
  )
}

interface InterviewPageProps {
  searchParams: Promise<{ page?: string }>
}

export default async function InterviewPage({ searchParams }: InterviewPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.page || "1", 10))

  return (
    <PageLayout>
      <PageHeader
        title="Go 面试题"
        description="精选 Go 语言面试题，助你轻松通过技术面试"
        breadcrumbs={[{ label: "Go 面试题" }]}
      />
      <Suspense
        fallback={
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-20 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <InterviewItems page={page} />
      </Suspense>
    </PageLayout>
  )
}
