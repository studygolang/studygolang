import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { TabFilter } from "@/components/tab-filter"
import { TopicList } from "@/components/topic-list"
import { FeedList } from "@/components/feed-list"
import { ArticleList } from "@/components/article-list"
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
import Image from "next/image"
import Link from "next/link"
import type {
  Article,
  Topic,
  Feed,
  HomeDocument,
  IndexNav,
  SiteStats,
  Reading,
  Comment,
  User,
  FriendLink,
  TopicNode,
  Subject,
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

interface HomeResponse {
  feeds?: Feed[] | null
  topics?: Topic[] | null
  articles?: Article[] | null
  docs?: HomeDocument[] | null
  subjects?: Subject[] | null
  users?: Record<number, User> | null
  nodes?: Record<number, TopicNode> | null
  index_navs?: IndexNav[] | null
  cur_nav?: IndexNav | null
  tab: string
  has_more?: boolean
}

interface HomeDataResult {
  home: HomeResponse | null
  stats?: SiteStats
  readings: Reading[]
  hotNodes: TopicNode[]
  friendLinks: FriendLink[]
  recentComments: Comment[]
  activeUsers: User[]
  newUsers: User[]
}

const fallbackNavs: IndexNav[] = [
  { tab: "all", name: "全部", data_source: "feed" },
  { tab: "recommend", name: "推荐", data_source: "feed" },
  { tab: "no_reply", name: "未回复", data_source: "no_reply" },
]

async function getHomeData(tab: string | undefined, page: number): Promise<HomeDataResult> {
  const homeParams = new URLSearchParams({ p: String(page) })
  if (tab) homeParams.set("tab", tab)

  const [
    homeData,
    stats,
    readingsData,
    hotNodesData,
    friendLinksData,
    recentCommentsData,
    activeUsersData,
    newestUsersData,
  ] = await Promise.allSettled([
    fetchAPI<HomeResponse>(`/home?${homeParams.toString()}`, { cache: "no-store" }),
    fetchAPI<SiteStats>("/stat/site", { cache: "no-store" }),
    fetchAPI<{ readings: Reading[] }>("/sidebar/readings/recent?limit=7", { cache: "no-store" }),
    fetchAPI<{ nodes: TopicNode[] }>("/sidebar/nodes/hot", { cache: "no-store" }),
    fetchAPI<{ links: FriendLink[] }>("/sidebar/friend/links", { cache: "no-store" }),
    fetchAPI<{ comments: Comment[] }>("/sidebar/comments/recent", { cache: "no-store" }),
    fetchAPI<{ users: User[] }>("/sidebar/users/active", { cache: "no-store" }),
    fetchAPI<{ users: User[] }>("/sidebar/users/newest", { cache: "no-store" }),
  ])

  return {
    home: homeData.status === "fulfilled" ? homeData.value : null,
    stats: stats.status === "fulfilled" ? stats.value : undefined,
    readings: readingsData.status === "fulfilled" ? (readingsData.value.readings ?? []) : [],
    hotNodes: hotNodesData.status === "fulfilled" ? (hotNodesData.value.nodes ?? []) : [],
    friendLinks: friendLinksData.status === "fulfilled" ? (friendLinksData.value.links ?? []) : [],
    recentComments: recentCommentsData.status === "fulfilled" ? (recentCommentsData.value.comments ?? []) : [],
    activeUsers: activeUsersData.status === "fulfilled" ? (activeUsersData.value.users ?? []) : [],
    newUsers: newestUsersData.status === "fulfilled" ? (newestUsersData.value.users ?? []) : [],
  }
}

function getDocumentUri(doc: HomeDocument): string {
  switch (doc.objtype) {
    case 0:
      return `/topics/${doc.objid}`
    case 1:
      return `/articles/${doc.objid}`
    case 2:
      return `/resources/${doc.objid}`
    case 4:
      return `/p/${doc.objid}`
    default:
      return `/wiki/${doc.objid}`
  }
}

function documentsToFeeds(
  docs: HomeDocument[],
  users: Record<number, User>,
  nodes: Record<number, TopicNode>
): Feed[] {
  return docs.map((doc) => ({
    Id: Number(doc.id) || doc.objid,
    Title: doc.title,
    Objid: doc.objid,
    Objtype: doc.objtype,
    Uid: doc.uid,
    Author: doc.author,
    Nid: doc.nid,
    Lastreplyuid: doc.lastreplyuid,
    Lastreplytime: doc.lastreplytime,
    Tags: doc.tags,
    Cmtnum: doc.cmtnum,
    Likenum: doc.likenum,
    Top: doc.top,
    Seq: 0,
    State: 0,
    CreatedAt: doc.pub_time || doc.created_at,
    updated_at: doc.updated_at,
    User: users[doc.uid],
    Lastreplyuser: users[doc.lastreplyuid],
    Node: nodes[doc.nid],
    Uri: getDocumentUri(doc),
  }))
}

function HomeSubjectList({ subjects }: { subjects: Subject[] }) {
  return (
    <div className="divide-y divide-border">
      {subjects.map((subject) => (
        <article key={subject.id} className="py-4 first:pt-0">
          <Link href={`/subject/${subject.id}`} className="group flex items-center gap-3">
            {subject.cover && (
              <Image
                src={subject.cover}
                alt=""
                width={40}
                height={40}
                className="h-10 w-10 shrink-0 rounded-lg object-cover"
              />
            )}
            <div className="min-w-0 flex-1">
              <h3 className="truncate text-[15px] font-semibold transition-colors group-hover:text-primary">
                {subject.name}
              </h3>
              <div className="mt-1.5 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                <span>专栏</span>
                {subject.user?.username && <span>{subject.user.username}</span>}
                <span>{subject.article_num} 篇文章</span>
              </div>
            </div>
          </Link>
        </article>
      ))}
    </div>
  )
}

interface HomePageProps {
  searchParams: Promise<{ tab?: string; p?: string }>
}

export default async function HomePage({ searchParams }: HomePageProps) {
  const sp = await searchParams
  const currentPage = Math.max(1, Number(sp.p) || 1)
  const {
    home,
    stats,
    readings,
    hotNodes,
    friendLinks,
    recentComments,
    activeUsers,
    newUsers,
  } = await getHomeData(sp.tab, currentPage)

  const feeds = home?.feeds ?? []
  const topics = home?.topics ?? []
  const articles = home?.articles ?? []
  const subjects = home?.subjects ?? []
  const documentFeeds = documentsToFeeds(home?.docs ?? [], home?.users ?? {}, home?.nodes ?? {})
  const activeTab = home?.tab ?? sp.tab ?? "all"
  const indexNavs = home?.index_navs?.length ? home.index_navs : fallbackNavs
  const hasContent = feeds.length + topics.length + articles.length + documentFeeds.length + subjects.length > 0

  const feedSource = feeds.length > 0 ? feeds : documentFeeds
  const trendingFromFeeds = feedSource.filter((feed) => feed.Objtype === 0).slice(0, 5)
  const trendingTopics = trendingFromFeeds.length > 0
    ? trendingFromFeeds.map((feed) => ({
        tid: feed.Objid,
        title: feed.Title,
        view: 0,
      } as Topic))
    : [...topics].sort((a, b) => (b.view || 0) - (a.view || 0)).slice(0, 5)

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <HeroBanner stats={stats} />

      <main className="mx-auto max-w-7xl px-4 py-6 lg:px-6">
        <div className="flex flex-col gap-6 lg:flex-row">
          <div className="min-w-0 flex-1">
            <div className="rounded-lg border border-border bg-card p-4">
              <TabFilter tabs={indexNavs} activeTab={activeTab} currentNav={home?.cur_nav} />
              <div className="mt-4 space-y-4">
                {feeds.length > 0 && <FeedList feeds={feeds} />}
                {topics.length > 0 && <TopicList topics={topics} />}
                {articles.length > 0 && <ArticleList articles={articles} />}
                {documentFeeds.length > 0 && <FeedList feeds={documentFeeds} />}
                {subjects.length > 0 && <HomeSubjectList subjects={subjects} />}
                {!hasContent && (
                  <div className="py-8 text-center text-sm text-muted-foreground">暂无主题</div>
                )}
              </div>
              {home?.has_more && (
                <div className="mt-6 text-center">
                  <Link
                    href={`/?${new URLSearchParams({ tab: activeTab, p: String(currentPage + 1) }).toString()}`}
                    className="inline-block cursor-pointer rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
                  >
                    加载更多
                  </Link>
                </div>
              )}
            </div>

            <div className="mt-6 lg:hidden">
              <NodeNavigation nodes={hotNodes} />
            </div>
          </div>

          <aside className="w-full shrink-0 space-y-4 lg:w-80">
            <LoginCard />
            <DailyQuestion />
            <MorningReading readings={readings} />
            <TrendingTopics topics={trendingTopics} />
            <div className="hidden lg:block">
              <NodeNavigation nodes={hotNodes} />
            </div>
            <LatestComments comments={recentComments} />
            <ActiveMembers users={newUsers} title="新会员" />
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
