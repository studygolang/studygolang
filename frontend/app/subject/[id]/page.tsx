import Link from "next/link"
import Image from "next/image"
import type { Metadata } from "next"
import { notFound } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { fetchAPINullable } from "@/lib/api"

interface Subject {
  id: number
  name: string
  cover: string
  intro: string
  article_num: number
  follower_num: number
}

interface SubjectArticle {
  id: number
  title: string
  author: string
  views: number
  created_at: string
}

interface SubjectDetailData {
  subject: Subject
  articles: SubjectArticle[]
  article_num: number
  follower_num: number
  followed: boolean
  page: number
  has_more: boolean
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>
}): Promise<Metadata> {
  const { id: idStr } = await params
  const id = parseInt(idStr, 10)

  if (!id || isNaN(id)) {
    return { title: "专题 - Go语言中文网" }
  }

  const data = await fetchAPINullable<SubjectDetailData>(`/subject/${id}`, {
    next: { revalidate: 60 },
  })

  const subjectName = data?.subject?.name
  return {
    title: subjectName ? `${subjectName} - Go语言中文网` : "专题 - Go语言中文网",
    description: data?.subject?.intro || "Go语言中文网专题栏目，汇聚优质内容",
  }
}

export default async function SubjectPage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>
  searchParams: Promise<{ p?: string }>
}) {
  const { id: idStr } = await params
  const id = parseInt(idStr, 10)

  if (!id || isNaN(id)) {
    return notFound()
  }

  const resolvedSearchParams = await searchParams
  const currentPage = parseInt(resolvedSearchParams.p || "1", 10) || 1

  const data = await fetchAPINullable<SubjectDetailData>(
    `/subject/${id}${currentPage > 1 ? `?p=${currentPage}` : ""}`,
    { next: { revalidate: 60 } }
  )

  if (!data || !data.subject || data.subject.id === 0) {
    return notFound()
  }

  const subject = data.subject

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title={subject.name}
        breadcrumbs={[{ label: "专题" }, { label: subject.name }]}
      />

      {/* 专题信息 */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex gap-6">
            {subject.cover && (
              <div className="relative h-24 w-24 shrink-0 overflow-hidden rounded-lg">
                <Image
                  src={subject.cover}
                  alt={subject.name}
                  fill
                  className="object-cover"
                />
              </div>
            )}
            <div className="flex-1">
              <h1 className="text-2xl font-bold">{subject.name}</h1>
              <p className="mt-2 text-muted-foreground">{subject.intro}</p>
              <div className="mt-3 flex items-center gap-4 text-sm text-muted-foreground">
                <span>{subject.article_num} 篇文章</span>
                <span>{subject.follower_num} 人关注</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 文章列表 */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">专题文章</CardTitle>
        </CardHeader>
        <CardContent>
          {data.articles && data.articles.length > 0 ? (
            <div className="divide-y">
              {data.articles.map((article) => (
                <Link
                  key={article.id}
                  href={`/articles/${article.id}`}
                  className="block py-3 hover:bg-muted/50"
                >
                  <div className="flex items-center justify-between">
                    <h3 className="font-medium">{article.title}</h3>
                    <span className="text-xs text-muted-foreground">
                      {article.created_at}
                    </span>
                  </div>
                  <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                    <span>{article.author}</span>
                    <span>{article.views} 浏览</span>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <p className="text-sm text-muted-foreground">暂无文章</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* 分页 */}
      {data.has_more && (
        <div className="mt-4 text-center">
          <Link
            href={`/subject/${id}?p=${currentPage + 1}`}
            className="text-sm text-primary hover:underline"
          >
            查看更多文章
          </Link>
        </div>
      )}

      {currentPage > 1 && (
        <div className="mt-2 text-center">
          <Link
            href={currentPage > 2 ? `/subject/${id}?p=${currentPage - 1}` : `/subject/${id}`}
            className="text-sm text-muted-foreground hover:underline"
          >
            上一页
          </Link>
        </div>
      )}
    </PageLayout>
  )
}
