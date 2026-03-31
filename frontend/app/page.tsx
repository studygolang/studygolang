import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { TabFilter } from "@/components/tab-filter"
import { TopicList } from "@/components/topic-list"
import { FeedList } from "@/components/feed-list"
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
  Feed,
  FeedListData,
  TopicListData,
  SiteStats,
  Reading,
  Comment,
  User,
  FriendLink,
  TopicNode,
} from "@/lib/types"
import { fetchAPI } from "@/lib/api"

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

// 首页数据响应类型
interface HomeDataResult {
  feeds: Feed[]
  topics: Topic[]
  hasMore: boolean
  useFeeds: boolean  // 是否使用 feed 模式
  trendingTopics: Topic[]
  stats?: SiteStats
  readings: Reading[]
  hotNodes: TopicNode[]
  friendLinks: FriendLink[]
  recentComments: Comment[]
  activeUsers: User[]
}

async function getHomeData(tab: string = 'all'): Promise<HomeDataResult> {
  const [
    homeData,
    stats,
    readingsData,
    hotNodesData,
    friendLinksData,
    recentCommentsData,
    activeUsersData,
  ] = await Promise.allSettled([
    fetchAPI<FeedListData>(`/home?tab=${tab}&p=1`, { cache: 'no-store' }),
    fetchAPI<SiteStats>('/stat/site', { cache: 'no-store' }),
    fetchAPI<{ readings: Reading[] }>('/sidebar/readings/recent?limit=7', { cache: 'no-store' }),
    fetchAPI<{ nodes: TopicNode[] }>('/sidebar/nodes/hot', { cache: 'no-store' }),
    fetchAPI<{ links: FriendLink[] }>('/sidebar/friend/links', { cache: 'no-store' }),
    fetchAPI<{ comments: Comment[] }>('/sidebar/comments/recent', { cache: 'no-store' }),
    fetchAPI<{ users: User[] }>('/sidebar/users/active', { cache: 'no-store' }),
  ])

  const homeValue = homeData.status === 'fulfilled' ? homeData.value : null
  // 判断返回的是 feeds 还是 topics
  const useFeeds = homeValue?.feeds !== undefined
  const feeds = homeValue?.feeds ?? ([] as Feed[])
  const topics = (homeValue as TopicListData | null)?.topics ?? ([] as Topic[])

  // 从 feeds 中提取话题用于趋势显示
  const trendingFromFeeds = feeds
    .filter(f => f.Objtype === 1) // 只取话题类型
    .slice(0, 5)

  return {
    feeds,
    topics,
    hasMore: homeValue?.has_more ?? false,
    useFeeds,
    trendingTopics: useFeeds ? trendingFromFeeds.map(f => ({
      tid: f.Objid,
      title: f.Title,
      view: 0,
    } as Topic)) : (topics ? [...topics].sort((a, b) => (b.view || 0) - (a.view || 0)).slice(0, 5) : []),
    stats: stats.status === 'fulfilled' ? stats.value : undefined,
    readings: readingsData.status === 'fulfilled' ? (readingsData.value.readings ?? []) : [] as Reading[],
    hotNodes: hotNodesData.status === 'fulfilled' ? (hotNodesData.value.nodes ?? []) : [] as TopicNode[],
    friendLinks: friendLinksData.status === 'fulfilled' ? (friendLinksData.value.links ?? []) : [] as FriendLink[],
    recentComments: recentCommentsData.status === 'fulfilled' ? (recentCommentsData.value.comments ?? []) : [] as Comment[],
    activeUsers: activeUsersData.status === 'fulfilled' ? (activeUsersData.value.users ?? []) : [] as User[],
  }
}

interface HomePageProps {
  searchParams: Promise<{ tab?: string }>
}

export default async function HomePage({ searchParams }: HomePageProps) {
  const { tab: tabParam } = await searchParams
  // tab 参数驱动首页内容筛选，有效值：all/recommend/no_reply/节点名
  const tab = tabParam ?? 'all'

  const {
    feeds,
    topics,
    hasMore,
    useFeeds,
    trendingTopics,
    stats,
    readings,
    hotNodes,
    friendLinks,
    recentComments,
    activeUsers,
  } = await getHomeData(tab)

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
                {useFeeds ? (
                  <FeedList feeds={feeds} />
                ) : (
                  <TopicList topics={topics} />
                )}
              </div>
              {/* Load More - 只在有更多数据时显示 */}
              {hasMore && (
                <div className="mt-6 text-center">
                  <button className="cursor-pointer rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80">
                    {"加载更多"}
                  </button>
                </div>
              )}
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
