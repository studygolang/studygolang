import type { Metadata } from "next"
import Link from "next/link"
import Image from "next/image"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPI } from "@/lib/api"
import type { SubjectListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "专栏 - Go语言中文网",
  description: "Go语言中文网专栏，汇聚优质专题内容",
}

async function getSubjects(
  page?: number
): Promise<SubjectListData | null> {
  try {
    return await fetchAPI<SubjectListData>(
      `/subjects${page && page > 1 ? `?p=${page}` : ""}`,
      { cache: "no-store" }
    )
  } catch {
    return null
  }
}

export default async function SubjectListPage({
  searchParams,
}: {
  searchParams: Promise<{ p?: string }>
}) {
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? "1", 10) || 1)

  const data = await getSubjects(page)
  const subjects = data?.subjects ?? []
  const hasMore = data?.has_more ?? false

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="专栏"
        description="汇聚优质专题内容"
        breadcrumbs={[{ label: "专栏" }]}
      />

      {subjects.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {subjects.map((subject) => (
            <Link key={subject.id} href={`/subject/${subject.id}`}>
              <Card className="h-full transition-shadow hover:shadow-md">
                <CardContent className="p-4">
                  <div className="flex gap-3">
                    {subject.cover && (
                      <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-lg">
                        <Image
                          src={subject.cover}
                          alt={subject.name}
                          fill
                          className="object-cover"
                        />
                      </div>
                    )}
                    <div className="min-w-0 flex-1">
                      <h3 className="truncate font-medium">{subject.name}</h3>
                      {subject.intro && (
                        <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">
                          {subject.intro}
                        </p>
                      )}
                      <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
                        <span>{subject.article_num} 篇</span>
                        <span>{subject.follower_num} 关注</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <p className="text-sm text-muted-foreground">暂无专栏</p>
        </div>
      )}

      {/* 分页 */}
      {page > 1 && (
        <div className="mt-6 flex justify-center gap-2">
          <Link
            href={page > 2 ? `/subject?p=${page - 1}` : "/subject"}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            上一页
          </Link>
        </div>
      )}
      {hasMore && (
        <div className="mt-2 flex justify-center">
          <Link
            href={`/subject?p=${page + 1}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        </div>
      )}
    </PageLayout>
  )
}
