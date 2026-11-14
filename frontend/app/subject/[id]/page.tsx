import Link from "next/link"
import Image from "next/image"
import { cookies } from "next/headers"
import type { Metadata } from "next"
import { notFound } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { SubjectActions } from "@/components/subject-actions"
import { fetchAPINullable } from "@/lib/api"
import { cn } from "@/lib/utils"

// followed / 投稿 / 编辑入口都取决于当前登录用户，必须动态渲染并透传 Cookie，
// 否则 ISR 会把某个用户的关注状态缓存给所有人。
export const dynamic = "force-dynamic"

async function authHeaders(): Promise<Record<string, string>> {
  const store = await cookies()
  const cookie = store.getAll().map((c) => `${c.name}=${c.value}`).join("; ")
  return cookie ? { Cookie: cookie } : {}
}

interface Subject {
  id: number
  name: string
  cover: string
  description: string
  article_num: number
  uid: number
}

// 后端返回 model.SubjectArticles（内嵌 Article），字段沿用 Article 的 json tag
interface SubjectArticle {
  id: number
  title: string
  author: string
  viewnum: number
  ctime: string
}

interface SubjectDetailData {
  subject: Subject
  articles: SubjectArticle[]
  article_num: number
  follower_num: number
  followed: boolean
  order_by: string
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
    cache: "no-store",
    headers: await authHeaders(),
  })

  const subjectName = data?.subject?.name
  return {
    title: subjectName ? `${subjectName} - Go语言中文网` : "专题 - Go语言中文网",
    description: data?.subject?.description || "Go语言中文网专题栏目，汇聚优质内容",
  }
}

export default async function SubjectPage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>
  searchParams: Promise<{ p?: string; order_by?: string }>
}) {
  const { id: idStr } = await params
  const id = parseInt(idStr, 10)

  if (!id || isNaN(id)) {
    return notFound()
  }

  const resolvedSearchParams = await searchParams
  const currentPage = parseInt(resolvedSearchParams.p || "1", 10) || 1
  const orderBy = resolvedSearchParams.order_by === "commented_at" ? "commented_at" : "added_at"

  const listHref = (p: number, order: string) => {
    const q = new URLSearchParams()
    if (p > 1) q.set("p", String(p))
    if (order !== "added_at") q.set("order_by", order)
    const qs = q.toString()
    return `/subject/${id}${qs ? `?${qs}` : ""}`
  }

  const data = await fetchAPINullable<SubjectDetailData>(
    `/subject/${id}?p=${currentPage}&order_by=${orderBy}`,
    { cache: "no-store", headers: await authHeaders() }
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
            <div className="min-w-0 flex-1">
              <h1 className="text-2xl font-bold">{subject.name}</h1>
              <p className="mt-2 text-muted-foreground">{subject.description}</p>
              <div className="mt-3 flex items-center gap-4 text-sm text-muted-foreground">
                <span>{data.article_num} 篇文章</span>
                <span>{data.follower_num} 人关注</span>
              </div>
            </div>
            <SubjectActions
              sid={subject.id}
              ownerUid={subject.uid}
              followed={data.followed}
            />
          </div>
        </CardContent>
      </Card>

      {/* 排序：与旧版 template/subject/index.html 的两个 tab 一致 */}
      <div className="mb-4 flex items-center gap-2">
        {[
          { key: "added_at", label: "最新收录" },
          { key: "commented_at", label: "最新评论" },
        ].map((tab) => (
          <Link
            key={tab.key}
            href={listHref(1, tab.key)}
            className={cn(
              "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
              orderBy === tab.key
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
          >
            {tab.label}
          </Link>
        ))}
      </div>

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
                      {article.ctime}
                    </span>
                  </div>
                  <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                    <span>{article.author}</span>
                    <span>{article.viewnum} 浏览</span>
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
            href={listHref(currentPage + 1, orderBy)}
            className="text-sm text-primary hover:underline"
          >
            查看更多文章
          </Link>
        </div>
      )}

      {currentPage > 1 && (
        <div className="mt-2 text-center">
          <Link
            href={listHref(currentPage - 1, orderBy)}
            className="text-sm text-muted-foreground hover:underline"
          >
            上一页
          </Link>
        </div>
      )}
    </PageLayout>
  )
}
