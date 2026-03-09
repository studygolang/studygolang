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
import type { Reading, SiteStats, Topic, Comment, User, FriendLink, TopicNode } from "@/lib/types"

interface SidebarData {
  stats?: SiteStats
  readings?: Reading[]
  trendingTopics?: Topic[]
  comments?: Comment[]
  activeUsers?: User[]
  friendLinks?: FriendLink[]
  nodes?: TopicNode[]
}

interface PageLayoutProps {
  children: React.ReactNode
  /** Show sidebar or full-width */
  sidebar?: boolean
  /** Custom sidebar content instead of default widgets */
  sidebarContent?: React.ReactNode
  /** Data to pass to default sidebar widgets */
  sidebarData?: SidebarData
}

export function PageLayout({
  children,
  sidebar = true,
  sidebarContent,
  sidebarData,
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
                <>
                  <LoginCard />
                  <DailyQuestion />
                  <MorningReading readings={sidebarData?.readings} />
                  <TrendingTopics topics={sidebarData?.trendingTopics} />
                  <div className="hidden lg:block">
                    <NodeNavigation nodes={sidebarData?.nodes} />
                  </div>
                  <LatestComments comments={sidebarData?.comments} />
                  <ActiveMembers users={sidebarData?.activeUsers} />
                  <StatsCard stats={sidebarData?.stats} />
                  <FriendLinks links={sidebarData?.friendLinks} />
                </>
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
