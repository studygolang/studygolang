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

interface PageLayoutProps {
  children: React.ReactNode
  /** Show sidebar or full-width */
  sidebar?: boolean
  /** Custom sidebar content instead of default widgets */
  sidebarContent?: React.ReactNode
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
                <>
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
