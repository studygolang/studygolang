import Link from "next/link"
import { Card, CardContent } from "@/components/ui/card"
import type { Comment, User } from "@/lib/types"

interface CommentDetailViewProps {
  comment: Comment
  nearbyComments: Comment[]
  users: Record<string, User>
  backLabel: string
  backHref: string
}

export function CommentDetailView({
  comment,
  nearbyComments,
  users,
  backLabel,
  backHref,
}: CommentDetailViewProps) {
  const commentUser = users[comment.uid]

  return (
    <>
      {/* 返回链接 */}
      <div className="mb-4">
        <Link
          href={backHref}
          className="text-sm text-muted-foreground hover:text-primary"
        >
          &larr; 返回{backLabel}
        </Link>
      </div>

      {/* 主评论 */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="mb-4 flex items-center gap-3">
            {commentUser && (
              <Link
                href={`/user/${commentUser.username}`}
                className="flex items-center gap-2"
              >
                <img
                  src={commentUser.avatar}
                  alt={commentUser.username}
                  className="h-8 w-8 rounded-full"
                />
                <span className="font-medium">{commentUser.username}</span>
              </Link>
            )}
            <span className="text-sm text-muted-foreground">
              #{comment.floor}
            </span>
            {comment.ctime && (
              <span className="text-sm text-muted-foreground">
                {comment.ctime}
              </span>
            )}
          </div>
          <div
            className="prose prose-sm dark:prose-invert max-w-none"
            dangerouslySetInnerHTML={{ __html: comment.content }}
          />
        </CardContent>
      </Card>

      {/* 附近评论 */}
      {nearbyComments.length > 0 && (
        <Card>
          <CardContent className="p-4">
            <h3 className="mb-3 text-sm font-semibold text-muted-foreground">
              相邻评论
            </h3>
            <div className="divide-y">
              {nearbyComments.map((c) => {
                const cUser = users[c.uid]
                return (
                  <div key={c.id || c.floor} className="py-3">
                    <div className="mb-2 flex items-center gap-2">
                      {cUser && (
                        <Link
                          href={`/user/${cUser.username}`}
                          className="text-sm font-medium"
                        >
                          {cUser.username}
                        </Link>
                      )}
                      <span className="text-xs text-muted-foreground">
                        #{c.floor}
                      </span>
                    </div>
                    <div
                      className="prose prose-sm dark:prose-invert max-w-none text-sm"
                      dangerouslySetInnerHTML={{ __html: c.content }}
                    />
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      )}
    </>
  )
}
