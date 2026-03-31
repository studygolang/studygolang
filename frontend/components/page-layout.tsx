import { Suspense } from "react"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import {
  LoginCard,
  DailyQuestion,
  MorningReading,
  TrendingTopics,
  LatestComments,
  ActiveMembers,
  StatsCard,
  FriendLinks,
} from "@/components/sidebar-widgets"
import { NodeNavigation } from "@/components/node-navigation"
import type { Reading, SiteStats, Comment, User, FriendLink, TopicNode } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

interface PageLayoutProps {
  children: React.ReactNode
  /** Show sidebar or full-width */
  sidebar?: boolean
  /** Custom sidebar content instead of default widgets */
  sidebarContent?: React.ReactNode
}

// 异步子组件：自动获取全部 sidebar 数据
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const AutoSidebar: () => any = async function AutoSidebar() {
  const [statsRes, readingsRes, commentsRes, activeUsersRes, friendLinksRes, nodesRes] =
    await Promise.allSettled([
      fetchAPINullable<SiteStats>("/stat/site", { cache: "no-store" }),
      fetchAPINullable<{ readings: Reading[] }>("/sidebar/readings/recent?limit=7", { cache: "no-store" }),
      fetchAPINullable<{ comments: Comment[] }>("/sidebar/comments/recent", { cache: "no-store" }),
      fetchAPINullable<{ users: User[] }>("/sidebar/users/active", { cache: "no-store" }),
      fetchAPINullable<{ links: FriendLink[] }>("/sidebar/friend/links", { cache: "no-store" }),
      fetchAPINullable<{ nodes: TopicNode[] }>("/sidebar/nodes/hot", { cache: "no-store" }),
    ])

  const stats = statsRes.status === "fulfilled" ? statsRes.value ?? undefined : undefined
  const readings = readingsRes.status === "fulfilled" ? (readingsRes.value?.readings ?? []) : []
  const comments = commentsRes.status === "fulfilled" ? (commentsRes.value?.comments ?? []) : []
  const activeUsers = activeUsersRes.status === "fulfilled" ? (activeUsersRes.value?.users ?? []) : []
  const friendLinks = friendLinksRes.status === "fulfilled" ? (friendLinksRes.value?.links ?? []) : []
  const nodes = nodesRes.status === "fulfilled" ? (nodesRes.value?.nodes ?? []) : []

  return (
    <>
      <LoginCard />
      <DailyQuestion />
      <MorningReading readings={readings} />
      {/* TrendingTopics 需要话题数据，其他页面不传（返回 null） */}
      <TrendingTopics topics={undefined} />
      <div className="hidden lg:block">
        <NodeNavigation nodes={nodes} />
      </div>
      <LatestComments comments={comments} />
      <ActiveMembers users={activeUsers} />
      <StatsCard stats={stats} />
      <FriendLinks links={friendLinks} />
    </>
  )
}

function SidebarSkeleton() {
  return (
    <div className="space-y-4">
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className="h-24 animate-pulse rounded-lg bg-muted" />
      ))}
    </div>
  )
}

export function PageLayout({
  children,
  sidebar = true,
  sidebarContent,
}: PageLayoutProps) {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-7xl px-4 py-6 lg:px-6">
        {sidebar ? (
          <div className="flex flex-col gap-6 lg:flex-row">
            {/* Main Content */}
            <div className="min-w-0 flex-1">{children}</div>

            {/* Sidebar */}
            <aside className="w-full shrink-0 space-y-4 lg:w-80">
              {sidebarContent ?? (
                <Suspense fallback={<SidebarSkeleton />}>
                  <AutoSidebar />
                </Suspense>
              )}
            </aside>
          </div>
        ) : (
          children
        )}
      </main>

      <SiteFooter />
    </div>
  )
}
