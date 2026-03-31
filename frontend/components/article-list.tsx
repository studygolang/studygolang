import Link from "next/link"
import {
  Eye,
  MessageSquare,
  ThumbsUp,
  Clock,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import type { Article } from "@/lib/types"
import { formatNum } from "@/lib/utils"

interface ArticleListProps {
  articles?: Article[]
}

function formatTime(ctime: string): string {
  try {
    const date = new Date(ctime)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    if (minutes < 60) return `${minutes} 分钟前`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours} 小时前`
    const days = Math.floor(hours / 24)
    if (days < 30) return `${days} 天前`
    return ctime.slice(0, 10)
  } catch {
    return ctime
  }
}

const tagColors = [
  "bg-primary/10 text-primary",
  "bg-accent/10 text-accent",
  "bg-chart-3/10 text-chart-3",
  "bg-chart-4/10 text-chart-4",
  "bg-chart-5/10 text-chart-5",
]

export function ArticleList({ articles = [] }: ArticleListProps) {
  if (articles.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无文章
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {articles.map((article, index) => (
        <article
          key={article.id}
          className="group rounded-lg border border-border bg-card p-4 transition-all hover:border-primary/20 hover:shadow-sm sm:p-5"
        >
          <div className="flex gap-4">
            <div className="min-w-0 flex-1 space-y-2">
              {/* Title */}
              <div className="flex items-start gap-2">
                <Link
                  href={`/articles/${article.id}`}
                  className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary sm:text-[17px]"
                >
                  {article.title}
                </Link>
              </div>

              {/* Excerpt：后端 Article 无 summary 字段，截取 content 前 150 字作摘要 */}
              {article.content && (
                <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                  {article.content.slice(0, 150)}
                </p>
              )}

              {/* Meta */}
              <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                <div className="flex items-center gap-1.5">
                  <Avatar className="h-5 w-5">
                    <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                      {article.author ? article.author.charAt(0).toUpperCase() : "?"}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-xs font-medium text-foreground">
                    {article.author}
                  </span>
                </div>
                {article.tags && (
                  <Badge
                    variant="secondary"
                    className={`text-[11px] font-medium ${tagColors[index % tagColors.length]}`}
                  >
                    {article.tags.split(",")[0]}
                  </Badge>
                )}
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Clock className="h-3 w-3" />
                  {formatTime(article.ctime)}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Eye className="h-3 w-3" />
                  {formatNum(article.viewnum)}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <MessageSquare className="h-3 w-3" />
                  {article.cmtnum}
                </span>
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <ThumbsUp className="h-3 w-3" />
                  {article.likenum}
                </span>
              </div>
            </div>
          </div>
        </article>
      ))}
    </div>
  )
}
