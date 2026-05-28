import Link from "next/link"
import { BookOpen, ExternalLink, Clock, Eye, ChevronLeft, ChevronRight } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { ReadingListData } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

export const metadata: Metadata = {
  title: "技术晨读 - Go语言中文网",
  description: "每日精选技术文章，开启你的技术早读时间",
  openGraph: {
    title: "技术晨读 - Go语言中文网",
    description: "每日精选技术文章，开启你的技术早读时间",
    type: "website",
  },
}

// 后端 MorningReading 无 lang 字段
const rtypeLabelMap: Record<number, string> = {
  1: "文章",
  2: "视频",
  3: "工具",
  4: "资讯",
}

interface ReadingsPageProps {
  searchParams: Promise<{ id?: string }>
}

export default async function ReadingsPage({ searchParams }: ReadingsPageProps) {
  const params = await searchParams
  const idParam = params.id

  // 后端参数名为 lastid（非 id）
  const apiPath = idParam ? `/readings?lastid=${idParam}` : "/readings"
  const data = await fetchAPINullable<ReadingListData>(apiPath, { cache: "no-store" })

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="技术晨读"
        description="每日精选技术文章，开启你的技术早读时间"
        breadcrumbs={[{ label: "技术晨读" }]}
      />

      {!data || !data.readings || data.readings.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">暂无晨读内容</p>
        </div>
      ) : (
        <div>
          {/* Reading cards */}
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {data.readings.map((reading) => (
              <Card
                key={reading.id}
                className="group overflow-hidden transition-all hover:border-primary/20 hover:shadow-sm"
              >
                <CardContent className="flex h-full flex-col p-0">
                  {/* 后端 MorningReading 无 cover 字段，统一显示占位图标 */}
                  <div className="flex h-40 items-center justify-center bg-gradient-to-br from-primary/10 to-primary/5">
                    <BookOpen className="h-10 w-10 text-primary/40" />
                  </div>

                  {/* Content */}
                  <div className="flex flex-1 flex-col p-4">
                    {/* Badges：后端无 lang 字段，只显示 rtype */}
                    <div className="mb-2 flex flex-wrap gap-1.5">
                      {reading.rtype && rtypeLabelMap[reading.rtype] && (
                        <Badge variant="secondary" className="text-[10px]">
                          {rtypeLabelMap[reading.rtype]}
                        </Badge>
                      )}
                      {reading.rdate && (
                        <Badge variant="outline" className="text-[10px]">
                          {reading.rdate}
                        </Badge>
                      )}
                    </div>

                    {/* 标题：后端用 content 字段存标题 */}
                    <a
                      href={reading.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                    >
                      {reading.content}
                    </a>

                    {/* Meta */}
                    <div className="mt-3 flex items-center justify-between border-t border-border pt-3">
                      <div className="flex items-center gap-3">
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {reading.rdate || reading.ctime}
                        </span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Eye className="h-3 w-3" />
                          {reading.clicknum}
                        </span>
                      </div>
                      <a
                        href={reading.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1 text-xs font-medium text-primary transition-colors hover:text-primary/80"
                        aria-label={"阅读 " + reading.content}
                      >
                        阅读
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Prev / Next navigation
              后端分页信息嵌套在 data.page 对象中：{ has_prev, prev_id, has_next, next_id }
              翻页参数使用 lastid（后端参数名）
          */}
          <div className="mt-8 flex items-center justify-center gap-4">
            {data.page?.has_prev && (
              <Link
                href={`/readings?id=${data.page.prev_id}`}
                className="flex items-center gap-1.5 rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
              >
                <ChevronLeft className="h-4 w-4" />
                上一期
              </Link>
            )}
            {data.page?.has_next && (
              <Link
                href={`/readings?id=${data.page.next_id}`}
                className="flex items-center gap-1.5 rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
              >
                下一期
                <ChevronRight className="h-4 w-4" />
              </Link>
            )}
          </div>
        </div>
      )}
    </PageLayout>
  )
}
