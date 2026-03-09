import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { TabFilter } from "@/components/tab-filter"
import { TopicList } from "@/components/topic-list"
import { NodeNavigation } from "@/components/node-navigation"
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
import { HeroBanner } from "@/components/hero-banner"
import type { Metadata } from "next"
import type {
  Topic,
  TopicListData,
  SiteStats,
  Reading,
  Comment,
  User,
  FriendLink,
  TopicNode,
} from "@/lib/types"

export const metadata: Metadata = {
  title: "Go语言中文网 - 中国最大的Go语言社区",
  description: "Go语言中文网是中国最大的Go语言学习社区，提供话题讨论、技术文章、开源项目、资源分享等内容，助力Gopher成长。",
  keywords: "Go语言,Golang,Go中文社区,Go学习,Go教程",
  openGraph: {
    title: "Go语言中文网 - 中国最大的Go语言社区",
    description: "中国最大的Go语言学习社区",
    type: "website",
    url: "https://studygolang.com",
  },
}

async function fetchFromAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const base = process.env.API_BASE_URL || 'http://localhost:8088'
  const res = await fetch(`${base}/api/v1${path}`, options)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  return json.data
}

async function getHomeData() {
  const [
    topicData,
    stats,
    readingsData,
    hotNodesData,
    friendLinksData,
    recentCommentsData,
    activeUsersData,
  ] = await Promise.allSettled([
    fetchFromAPI<TopicListData>('/home?tab=all&p=1', { cache: 'no-store' }),
    fetchFromAPI<SiteStats>('/stat/site', { cache: 'no-store' }),
    fetchFromAPI<{ readings: Reading[] }>('/sidebar/readings/recent?limit=7', { cache: 'no-store' }),
    fetchFromAPI<{ nodes: TopicNode[] }>('/sidebar/nodes/hot', { cache: 'no-store' }),
    fetchFromAPI<{ links: FriendLink[] }>('/sidebar/friend/links', { cache: 'no-store' }),
    fetchFromAPI<{ comments: Comment[] }>('/sidebar/comments/recent', { cache: 'no-store' }),
    fetchFromAPI<{ users: User[] }>('/sidebar/users/active', { cache: 'no-store' }),
  ])

  const topicsValue = topicData.status === 'fulfilled' ? topicData.value : null

  return {
    topics: topicsValue?.topics ?? [] as Topic[],
    trendingTopics: topicsValue?.topics
      ? [...topicsValue.topics].sort((a, b) => b.viewnum - a.viewnum).slice(0, 5)
      : [] as Topic[],
    stats: stats.status === 'fulfilled' ? stats.value : undefined,
    readings: readingsData.status === 'fulfilled' ? (readingsData.value.readings ?? []) : [] as Reading[],
    hotNodes: hotNodesData.status === 'fulfilled' ? (hotNodesData.value.nodes ?? []) : [] as TopicNode[],
    friendLinks: friendLinksData.status === 'fulfilled' ? (friendLinksData.value.links ?? []) : [] as FriendLink[],
    recentComments: recentCommentsData.status === 'fulfilled' ? (recentCommentsData.value.comments ?? []) : [] as Comment[],
    activeUsers: activeUsersData.status === 'fulfilled' ? (activeUsersData.value.users ?? []) : [] as User[],
  }
}

export default async function HomePage() {
  const {
    topics,
    trendingTopics,
    stats,
    readings,
    hotNodes,
    friendLinks,
    recentComments,
    activeUsers,
  } = await getHomeData()

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <HeroBanner stats={stats} />

      <main className="mx-auto max-w-7xl px-4 py-6 lg:px-6">
        <div className="flex flex-col gap-6 lg:flex-row">
          {/* Main Content */}
          <div className="min-w-0 flex-1">
            <div className="rounded-lg border border-border bg-card p-4">
              <TabFilter />
              <div className="mt-4">
                <TopicList topics={topics} />
              </div>
              {/* Load More */}
              <div className="mt-6 text-center">
                <button className="rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80">
                  {"加载更多"}
                </button>
              </div>
            </div>

            {/* Node Navigation - below main on mobile, visible always */}
            <div className="mt-6 lg:hidden">
              <NodeNavigation nodes={hotNodes} />
            </div>
          </div>

          {/* Sidebar */}
          <aside className="w-full shrink-0 space-y-4 lg:w-80">
            <LoginCard />
            <DailyQuestion />
            <MorningReading readings={readings} />
            <TrendingTopics topics={trendingTopics} />
            <div className="hidden lg:block">
              <NodeNavigation nodes={hotNodes} />
            </div>
            <LatestComments comments={recentComments} />
            <ActiveMembers users={activeUsers} />
            <StatsCard stats={stats} />
            <FriendLinks links={friendLinks} />
          </aside>
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
