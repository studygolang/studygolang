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

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <HeroBanner />

      <main className="mx-auto max-w-7xl px-4 py-6 lg:px-6">
        <div className="flex flex-col gap-6 lg:flex-row">
          {/* Main Content */}
          <div className="min-w-0 flex-1">
            <div className="rounded-lg border border-border bg-card p-4">
              <TabFilter />
              <div className="mt-4">
                <TopicList />
              </div>
              {/* Load More */}
              <div className="mt-6 text-center">
                <button className="rounded-md bg-secondary px-6 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80">
                  {"\u52a0\u8f7d\u66f4\u591a"}
                </button>
              </div>
            </div>

            {/* Node Navigation - below main on mobile, visible always */}
            <div className="mt-6 lg:hidden">
              <NodeNavigation />
            </div>
          </div>

          {/* Sidebar */}
          <aside className="w-full shrink-0 space-y-4 lg:w-80">
            <LoginCard />
            <DailyQuestion />
            <MorningReading />
            <TrendingTopics />
            <div className="hidden lg:block">
              <NodeNavigation />
            </div>
            <LatestComments />
            <ActiveMembers />
            <StatsCard />
            <FriendLinks />
          </aside>
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
