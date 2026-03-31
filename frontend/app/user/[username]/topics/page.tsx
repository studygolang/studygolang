import type { Metadata } from "next"
import { notFound } from "next/navigation"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { TopicList } from "@/components/topic-list"
import { User as UserIcon } from "lucide-react"
import type { Topic, User } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

interface UserTopicsData {
  user: User
  topics: Topic[]
  total: number
  page: number
  has_more: boolean
}

interface UserTopicsPageProps {
  params: Promise<{ username: string }>
  searchParams: Promise<{ p?: string }>
}

export async function generateMetadata({ params }: UserTopicsPageProps): Promise<Metadata> {
  const { username } = await params
  const data = await fetchAPINullable<UserTopicsData>(
    `/users/${username}/topics`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return { title: "用户话题 - Go语言中文网" }
  }

  const displayName = data.user.name || data.user.username
  return {
    title: `${displayName}的话题 - Go语言中文网`,
    description: `${displayName}在Go语言中文网发布的话题列表`,
  }
}

export default async function UserTopicsPage({ params, searchParams }: UserTopicsPageProps) {
  const { username } = await params
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const data = await fetchAPINullable<UserTopicsData>(
    `/users/${username}/topics?p=${page}`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    notFound()
  }

  const { user, topics = [], total, has_more } = data
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
              发布了 {total} 个话题
            </p>
          </div>
        </div>
      </div>

      {/* 话题列表 */}
      <div className="rounded-lg border border-border bg-card">
        <div className="border-b border-border px-4 py-3">
          <h2 className="text-sm font-medium text-muted-foreground">话题列表</h2>
        </div>
        <div className="p-4">
          <TopicList topics={topics} />
          
          {/* 分页 */}
          {topics.length === 0 && (
            <div className="py-8 text-center text-sm text-muted-foreground">
              暂无话题
            </div>
          )}
          
          {/* 简单分页导航 */}
          {has_more && (
            <div className="mt-4 flex justify-center">
              <Link
                href={`/user/${username}/topics?p=${page + 1}`}
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
