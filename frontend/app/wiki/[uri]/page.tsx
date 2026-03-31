import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { Clock, User } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import type { Wiki } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"
import { formatDate } from "@/lib/utils"

export const revalidate = 60

interface WikiDetailPageProps {
  params: Promise<{ uri: string }>
}

export async function generateMetadata({ params }: WikiDetailPageProps): Promise<Metadata> {
  const { uri } = await params
  // 后端返回 { wiki: {...} }，需要取 wiki 字段
  const data = await fetchAPINullable<{ wiki: Wiki }>(`/wiki/${uri}`)
  const wiki = data?.wiki
  if (!wiki) {
    return { title: "Wiki - Go语言中文网" }
  }
  return {
    title: `${wiki.title} - Wiki - Go语言中文网`,
    description: wiki.content?.slice(0, 160).replace(/<[^>]+>/g, "") || wiki.title,
  }
}

export default async function WikiDetailPage({ params }: WikiDetailPageProps) {
  const { uri } = await params
  // 后端返回 { wiki: {...} }，需要取 wiki 字段
  const data = await fetchAPINullable<{ wiki: Wiki }>(`/wiki/${uri}`)
  const wiki = data?.wiki

  if (!wiki) {
    notFound()
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title={wiki.title}
        breadcrumbs={[
          { label: "Wiki", href: "/wiki" },
          { label: wiki.title },
        ]}
      />

      <div className="mx-auto max-w-4xl">
        <Card>
          <CardContent className="p-6 sm:p-8">
            {/* Meta info */}
            <div className="mb-6 flex flex-wrap items-center gap-4 border-b border-border pb-4 text-sm text-muted-foreground">
              {wiki.author && (
                <span className="flex items-center gap-1.5">
                  <User className="h-4 w-4" />
                  {wiki.author}
                </span>
              )}
              <span className="flex items-center gap-1.5">
                <Clock className="h-4 w-4" />
                最后更新：{formatDate(wiki.mtime || wiki.ctime)}
              </span>
            </div>

            {/* Wiki Content */}
            {wiki.content ? (
              <div
                className="prose prose-sm max-w-none dark:prose-invert prose-headings:font-semibold prose-headings:text-foreground prose-p:text-muted-foreground prose-a:text-primary prose-code:text-primary prose-pre:bg-muted"
                dangerouslySetInnerHTML={{ __html: sanitizeHtml(wiki.content) }}
              />
            ) : (
              <p className="text-sm text-muted-foreground">暂无内容</p>
            )}
          </CardContent>
        </Card>

        <div className="mt-4 text-center">
          <Link
            href="/wiki"
            className="text-sm text-muted-foreground transition-colors hover:text-primary"
          >
            ← 返回 Wiki 列表
          </Link>
        </div>
      </div>
    </PageLayout>
  )
}
