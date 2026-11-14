import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { Clock, Eye, User } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import type { Wiki } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { MarkdownContent } from "@/components/markdown-content"
import { WikiEditButton } from "@/components/wiki-edit-button"
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
    return notFound()
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
                <Link
                  href={`/user/${wiki.author}`}
                  className="flex items-center gap-1.5 hover:text-primary"
                >
                  <User className="h-4 w-4" />
                  {wiki.author}
                </Link>
              )}
              <span className="flex items-center gap-1.5">
                <Clock className="h-4 w-4" />
                最后更新：{formatDate(wiki.mtime || wiki.ctime)}
              </span>
              <span className="flex items-center gap-1.5">
                <Eye className="h-4 w-4" />
                {wiki.viewnum ?? 0} 次点击
              </span>
              <span className="ml-auto">
                <WikiEditButton uri={wiki.uri} ownerUid={wiki.uid} />
              </span>
            </div>

            {/* Wiki Content：编辑器产出 Markdown，走 markdown 渲染（兼容历史 HTML 内容） */}
            {wiki.content ? (
              <MarkdownContent
                content={wiki.content}
                className="prose-headings:font-semibold prose-headings:text-foreground prose-p:text-muted-foreground prose-a:text-primary prose-code:text-primary prose-pre:bg-muted"
              />
            ) : (
              <p className="text-sm text-muted-foreground">暂无内容</p>
            )}

            {/* 本页贡献者 */}
            {!!wiki.contributors?.length && (
              <div className="mt-8 border-t border-border pt-4">
                <h5 className="mb-3 text-sm font-semibold text-foreground">本页贡献者</h5>
                <div className="flex flex-wrap gap-3">
                  {wiki.contributors.map((c) => (
                    <Link
                      key={c.uid}
                      href={`/user/${c.username}`}
                      title={c.username}
                      className="flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
                    >
                      <Avatar className="h-8 w-8">
                        {c.avatar && <AvatarImage src={c.avatar} alt={c.username} />}
                        <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                          {c.username.charAt(0).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      {c.username}
                    </Link>
                  ))}
                </div>
              </div>
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
