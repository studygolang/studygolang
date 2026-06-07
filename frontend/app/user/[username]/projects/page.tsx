import type { Metadata } from "next"
import { notFound } from "next/navigation"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { ProjectList } from "@/components/project-list"
import { User as UserIcon } from "lucide-react"
import type { Project, User } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

interface UserProjectsData {
  user: User
  projects: Project[]
  total: number
  page: number
  has_more: boolean
}

interface UserProjectsPageProps {
  params: Promise<{ username: string }>
  searchParams: Promise<{ p?: string }>
}

export async function generateMetadata({ params }: UserProjectsPageProps): Promise<Metadata> {
  const { username } = await params
  const data = await fetchAPINullable<UserProjectsData>(
    `/users/${username}/projects`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return { title: "用户项目 - Go语言中文网" }
  }

  const displayName = data.user.name || data.user.username
  return {
    title: `${displayName}的开源项目 - Go语言中文网`,
    description: `${displayName}在Go语言中文网分享的开源项目列表`,
  }
}

export default async function UserProjectsPage({ params, searchParams }: UserProjectsPageProps) {
  const { username } = await params
  const { p: pageStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? '1', 10) || 1)

  const data = await fetchAPINullable<UserProjectsData>(
    `/users/${username}/projects?p=${page}`,
    { cache: "no-store" }
  )

  if (!data?.user) {
    return notFound()
  }

  const { user, projects = [], total, has_more } = data
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
              发布了 {total} 个开源项目
            </p>
          </div>
        </div>
      </div>

      {/* 项目列表 */}
      <div className="rounded-lg border border-border bg-card">
        <div className="border-b border-border px-4 py-3">
          <h2 className="text-sm font-medium text-muted-foreground">开源项目列表</h2>
        </div>
        <div className="p-4">
          <ProjectList projects={projects} />
          
          {/* 分页 */}
          {projects.length === 0 && (
            <div className="py-8 text-center text-sm text-muted-foreground">
              暂无项目
            </div>
          )}
          
          {/* 简单分页导航 */}
          {has_more && (
            <div className="mt-4 flex justify-center">
              <Link
                href={`/user/${username}/projects?p=${page + 1}`}
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
