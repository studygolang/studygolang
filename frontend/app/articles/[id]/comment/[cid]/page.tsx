import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { Clock, MessageSquare } from "lucide-react"
import { formatTime } from "@/lib/utils"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { MarkdownContent } from "@/components/markdown-content"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { fetchAPI } from "@/lib/api"
import { OBJTYPE_ARTICLE } from "@/lib/types"
import type { ArticleDetailData, CommentDetailData, User } from "@/lib/types"

async function getArticleDetail(id: string): Promise<ArticleDetailData | null> {
  try {
    return await fetchAPI<ArticleDetailData>(
      `/articles/${id}`,
      { cache: 'no-store' }
    )
  } catch {
    return null
  }
}

async function getCommentDetail(cid: string, objid: string): Promise<CommentDetailData | null> {
  try {
    return await fetchAPI<CommentDetailData>(
      `/comments/${cid}/detail?objid=${objid}&objtype=${OBJTYPE_ARTICLE}`,
      { cache: 'no-store' }
    )
  } catch {
    return null
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string; cid: string }>
}): Promise<Metadata> {
  const { id, cid } = await params
  const [articleData, commentData] = await Promise.all([
    getArticleDetail(id),
    getCommentDetail(cid, id),
  ])

  if (!articleData?.article || !commentData?.comment) {
    return {
      title: "评论不存在 - Go语言中文网",
    }
  }

  const { article } = articleData
  const { comment } = commentData
  const preview = comment.content.slice(0, 50)

  return {
    title: `第 ${comment.floor} 楼评论 - ${article.title} - Go语言中文网`,
    description: preview,
  }
}

export default async function ArticleCommentPage({
  params,
}: {
  params: Promise<{ id: string; cid: string }>
}) {
  const { id, cid } = await params
  const [articleData, commentData] = await Promise.all([
    getArticleDetail(id),
    getCommentDetail(cid, id),
  ])

  if (!articleData?.article) {
    return notFound()
  }

  if (!commentData?.comment) {
    return notFound()
  }

  const { article } = articleData
  const { comment, nearby_comments, users } = commentData

  // 获取评论用户信息
  const commentUser = users[String(comment.uid)] as User | undefined
  const userDisplayName = commentUser?.name || `用户 #${comment.uid}`
  const userAvatar = commentUser?.avatar || ""
  const userInitial = userDisplayName.charAt(0).toUpperCase()

  return (
    <PageLayout>
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "技术文章", href: "/articles" },
          {
            label: article.title.length > 20 ? article.title.slice(0, 20) + "..." : article.title,
            href: `/articles/${id}`
          },
          { label: `第 ${comment.floor} 楼评论` },
        ]}
      />

      <div className="space-y-6">
        {/* 评论详情卡片 */}
        <div className="rounded-lg border bg-card p-6">
          <div className="mb-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src={userAvatar} alt={userDisplayName} />
                <AvatarFallback>{userInitial}</AvatarFallback>
              </Avatar>
              <div>
                <Link
                  href={`/user/${commentUser?.username || comment.uid}`}
                  className="font-medium text-primary hover:underline"
                >
                  {userDisplayName}
                </Link>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Clock className="h-3 w-3" />
                    {formatTime(comment.ctime)}
                  </span>
                  <span>·</span>
                  <span>第 {comment.floor} 楼</span>
                </div>
              </div>
            </div>
          </div>

          <div className="prose prose-sm max-w-none dark:prose-invert">
            <MarkdownContent content={comment.content} />
          </div>

          <div className="mt-4 flex items-center gap-4 text-sm text-muted-foreground">
            <Link
              href={`/articles/${id}`}
              className="flex items-center gap-1 hover:text-primary"
            >
              <MessageSquare className="h-4 w-4" />
              查看完整文章
            </Link>
          </div>
        </div>

        {/* 附近评论 */}
        {nearby_comments.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">附近评论</h2>
            <div className="space-y-4">
              {nearby_comments.map((nearbyComment) => {
                const nearbyUser = users[String(nearbyComment.uid)] as User | undefined
                const nearbyDisplayName = nearbyUser?.name || `用户 #${nearbyComment.uid}`
                const nearbyAvatar = nearbyUser?.avatar || ""
                const nearbyInitial = nearbyDisplayName.charAt(0).toUpperCase()

                return (
                  <div
                    key={nearbyComment.cid}
                    className="rounded-lg border bg-card p-4"
                  >
                    <div className="mb-3 flex items-center gap-3">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={nearbyAvatar} alt={nearbyDisplayName} />
                        <AvatarFallback>{nearbyInitial}</AvatarFallback>
                      </Avatar>
                      <div className="flex-1">
                        <Link
                          href={`/user/${nearbyUser?.username || nearbyComment.uid}`}
                          className="text-sm font-medium text-primary hover:underline"
                        >
                          {nearbyDisplayName}
                        </Link>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {formatTime(nearbyComment.ctime)}
                          </span>
                          <span>·</span>
                          <span>第 {nearbyComment.floor} 楼</span>
                        </div>
                      </div>
                    </div>

                    <MarkdownContent content={nearbyComment.content} />
                  </div>
                )
              })}
            </div>
          </div>
        )}
      </div>
    </PageLayout>
  )
}
