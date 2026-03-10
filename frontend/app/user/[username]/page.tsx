import Link from "next/link"
import {
  User as UserIcon,
  MapPin,
  Globe,
  Github,
  Building2,
  Users,
  BookOpen,
  MessageSquare,
  Star,
  Clock,
  Eye,
} from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import type { User, Topic, Article } from "@/lib/types"

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const base = process.env.API_BASE_URL || "http://localhost:8090"
    const res = await fetch(`${base}/api/v1${path}`, options)
    if (!res.ok) return null
    const json = await res.json()
    return json.code === 0 ? json.data : null
  } catch {
    return null
  }
}

interface UserDetailData {
  user: User
  topics?: Topic[]
  articles?: Article[]
}

function formatNum(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + "k"
  return String(n)
}

interface UserPageProps {
  params: Promise<{ username: string }>
}

export async function generateMetadata({ params }: UserPageProps): Promise<Metadata> {
  const { username } = await params
  const data = await fetchAPI<UserDetailData>(
    `/user/${username}`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return { title: "用户未找到 - Go语言中文网" }
  }

  const { user } = data
  const displayName = user.name || user.username
  return {
    title: `${displayName}的主页 - Go语言中文网`,
    description: user.tagline || `${displayName} 的个人主页 - Go语言中文网`,
    openGraph: {
      title: `${displayName}的主页`,
      description: user.tagline || "",
      images: user.avatar ? [{ url: user.avatar }] : [],
    },
  }
}

export default async function UserPage({ params }: UserPageProps) {
  const { username } = await params
  const data = await fetchAPI<UserDetailData>(
    `/user/${username}`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return (
      <PageLayout sidebar={false}>
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <UserIcon className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm font-medium text-foreground">用户不存在</p>
          <p className="mt-1 text-sm text-muted-foreground">该用户不存在或已注销</p>
        </div>
      </PageLayout>
    )
  }

  const { user, topics, articles } = data
  const displayName = user.name || user.username

  return (
    <PageLayout>
      {/* Profile card */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex flex-col gap-5 sm:flex-row sm:items-start">
            {/* Avatar */}
            {user.avatar ? (
              <img
                src={user.avatar}
                alt={displayName}
                className="h-20 w-20 shrink-0 rounded-full object-cover ring-2 ring-border"
              />
            ) : (
              <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-primary/10 ring-2 ring-border">
                <UserIcon className="h-10 w-10 text-primary" />
              </div>
            )}

            {/* Info */}
            <div className="min-w-0 flex-1">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-xl font-bold text-foreground">{displayName}</h1>
                {user.username !== displayName && (
                  <span className="text-sm text-muted-foreground">@{user.username}</span>
                )}
                {user.is_vip && (
                  <Badge className="bg-chart-3/20 text-chart-3 text-[10px]">VIP</Badge>
                )}
                {user.is_root && (
                  <Badge variant="destructive" className="text-[10px]">管理员</Badge>
                )}
                {user.level > 0 && (
                  <Badge variant="secondary" className="text-[10px]">
                    Lv.{user.level}
                  </Badge>
                )}
              </div>

              {user.tagline && (
                <p className="mt-1.5 text-sm text-muted-foreground">{user.tagline}</p>
              )}

              {/* Meta info */}
              <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1.5">
                {user.city && (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <MapPin className="h-3 w-3" />
                    {user.city}
                  </span>
                )}
                {user.company && (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Building2 className="h-3 w-3" />
                    {user.company}
                  </span>
                )}
                {user.website && (
                  <a
                    href={user.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
                  >
                    <Globe className="h-3 w-3" />
                    {user.website}
                  </a>
                )}
                {user.github && (
                  <a
                    href={`https://github.com/${user.github}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
                  >
                    <Github className="h-3 w-3" />
                    {user.github}
                  </a>
                )}
              </div>

              {/* Stats */}
              <div className="mt-4 flex flex-wrap gap-6 border-t border-border pt-4">
                <div className="flex items-center gap-1.5">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-semibold">{formatNum(user.follow_count)}</span>
                  <span className="text-xs text-muted-foreground">关注</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-semibold">{formatNum(user.fans_count)}</span>
                  <span className="text-xs text-muted-foreground">粉丝</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">加入于 {user.ctime}</span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Recent topics */}
      {topics && topics.length > 0 && (
        <div className="mb-6 space-y-3">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <MessageSquare className="h-4 w-4" />
            最近话题
          </h2>
          <div className="space-y-2">
            {topics.map((topic) => (
              <Card key={topic.tid} className="group transition-all hover:border-primary/20">
                <CardContent className="p-4">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 flex-1">
                      <Link
                        href={`/topics/${topic.tid}`}
                        className="text-sm font-medium text-foreground transition-colors group-hover:text-primary"
                      >
                        {topic.title}
                      </Link>
                      <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1">
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {topic.ctime}
                        </span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Eye className="h-3 w-3" />
                          {formatNum(topic.viewnum)}
                        </span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <MessageSquare className="h-3 w-3" />
                          {topic.replynum}
                        </span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Recent articles */}
      {articles && articles.length > 0 && (
        <div className="space-y-3">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <BookOpen className="h-4 w-4" />
            最近文章
          </h2>
          <div className="space-y-2">
            {articles.map((article) => (
              <Card key={article.id} className="group transition-all hover:border-primary/20">
                <CardContent className="p-4">
                  <div className="flex items-start gap-3">
                    <div className="min-w-0 flex-1">
                      <Link
                        href={`/articles/${article.id}`}
                        className="text-sm font-medium text-foreground transition-colors group-hover:text-primary"
                      >
                        {article.title}
                      </Link>
                      {/* 后端 Article 无 summary 字段，截取 content 前 100 字 */}
                      {article.content && (
                        <p className="mt-1 line-clamp-1 text-xs text-muted-foreground">
                          {article.content.slice(0, 100)}
                        </p>
                      )}
                      <div className="mt-1.5 flex flex-wrap items-center gap-x-3 gap-y-1">
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {/* 后端 Article 用 pub_date（非 pubdate）字段 */}
                          {article.pub_date || article.ctime}
                        </span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Eye className="h-3 w-3" />
                          {formatNum(article.viewnum)}
                        </span>
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Star className="h-3 w-3" />
                          {formatNum(article.likenum)}
                        </span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Empty state when no content */}
      {(!topics || topics.length === 0) && (!articles || articles.length === 0) && (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <UserIcon className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">该用户暂无内容</p>
        </div>
      )}
    </PageLayout>
  )
}
