import type { Metadata } from "next"
import { notFound } from "next/navigation"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { User as UserIcon, ExternalLink, ThumbsUp, MessageSquare, Clock } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { formatNum } from "@/lib/utils"
import type { Resource, User } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

interface UserResourcesData {
  user: User
  resources: Resource[]
  total: number
  page: number
  has_more: boolean
}

interface UserResourcesPageProps {
  params: Promise<{ username: string }>
  searchParams: Promise<{ p?: string }>
}

export async function generateMetadata({ params }: UserResourcesPageProps): Promise<Metadata> {
  const { username } = await params
  const data = await fetchAPINullable<UserResourcesData>(
    `/users/${username}/resources`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return { title: "用户资源 - Go语言中文网" }
  }

  const displayName = data.user.name || data.user.username
  return {
    title: `${displayName}的资源 - Go语言中文网`,
    description: `${displayName}在Go语言中文网分享的资源列表`,
  }
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

export default async function UserResourcesPage({ params, searchParams }: UserResourcesPageProps) {
  const { username } = await params
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const data = await fetchAPINullable<UserResourcesData>(
    `/users/${username}/resources?p=${page}`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    notFound()
  }

  const { user, resources = [], total, has_more } = data
  const displayName = user.name || user.username

  return (
    <PageLayout>
      {/* 用户信息卡片 */}
      <div className="mb-6 rounded-lg border border-border bg-card p-4">
        <div className="flex items-center gap-3">
          {user.avatar ? (
            <img
              src={user.avatar}
              alt={displayName}
              className="h-12 w-12 rounded-full object-cover ring-2 ring-border"
            />
          ) : (
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 ring-2 ring-border">
              <UserIcon className="h-6 w-6 text-primary" />
            </div>
          )}
          <div>
            <h1 className="text-lg font-semibold text-foreground">{displayName}</h1>
            <p className="text-sm text-muted-foreground">
              分享了 {total} 个资源
            </p>
          </div>
        </div>
      </div>

      {/* 资源列表 */}
      <div className="rounded-lg border border-border bg-card">
        <div className="border-b border-border px-4 py-3">
          <h2 className="text-sm font-medium text-muted-foreground">资源列表</h2>
        </div>
        <div className="p-4">
          {resources.length === 0 ? (
            <div className="py-8 text-center text-sm text-muted-foreground">
              暂无资源
            </div>
          ) : (
            <div className="space-y-3">
              {resources.map((resource) => (
                <Card
                  key={resource.id}
                  className="group transition-all hover:border-primary/20 hover:shadow-sm"
                >
                  <CardContent className="p-4 sm:p-5">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0 flex-1">
                        {/* Title */}
                        <div className="flex items-start gap-2">
                          {resource.catname && (
                            <Badge className="mt-0.5 shrink-0 bg-primary/10 text-[10px] font-semibold text-primary">
                              {resource.catname}
                            </Badge>
                          )}
                          <Link
                            href={`/resources/${resource.id}`}
                            className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                          >
                            {resource.title}
                          </Link>
                        </div>

                        {/* Description */}
                        {resource.content && (
                          <p className="mt-1.5 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                            {resource.content}
                          </p>
                        )}

                        {/* Tags + Meta */}
                        <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                          {resource.tags && (
                            <div className="flex items-center gap-1.5">
                              {resource.tags.split(",").slice(0, 3).map((tag) => (
                                <Badge
                                  key={tag}
                                  variant="secondary"
                                  className="text-[10px] font-medium"
                                >
                                  {tag.trim()}
                                </Badge>
                              ))}
                            </div>
                          )}
                          <span className="text-xs text-muted-foreground">
                            {resource.user?.name || resource.user?.username || `uid:${resource.uid}`}
                          </span>
                          <span className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Clock className="h-3 w-3" />
                            {resource.ctime}
                          </span>
                          <span className="flex items-center gap-1 text-xs text-muted-foreground">
                            <ThumbsUp className="h-3 w-3" />
                            {formatNum(resource.likenum)}
                          </span>
                          <span className="flex items-center gap-1 text-xs text-muted-foreground">
                            <MessageSquare className="h-3 w-3" />
                            {resource.cmtnum}
                          </span>
                        </div>
                      </div>

                      {/* External link button */}
                      {resource.url && (
                        <a
                          href={resource.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
                          aria-label={"访问 " + resource.title}
                        >
                          <ExternalLink className="h-4 w-4" />
                        </a>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
          
          {/* 简单分页导航 */}
          {has_more && (
            <div className="mt-4 flex justify-center">
              <Link
                href={`/user/${username}/resources?p=${page + 1}`}
                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
              >
                加载更多
              </Link>
            </div>
          )}
        </div>
      </div>
    </PageLayout>
  )
}
