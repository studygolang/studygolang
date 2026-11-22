import Link from "next/link"
import { cn } from "@/lib/utils"
import type { IndexNav } from "@/lib/types"

interface TabFilterProps {
  tabs: IndexNav[]
  activeTab: string
  currentNav?: IndexNav | null
}

export function TabFilter({ tabs, activeTab, currentNav }: TabFilterProps) {
  return (
    <div className="space-y-3">
      <nav aria-label="首页频道" className="flex items-center gap-1 overflow-x-auto border-b border-border">
        {tabs.map((tab) => (
          <Link
            key={tab.tab}
            href={`/?tab=${encodeURIComponent(tab.tab)}`}
            aria-current={activeTab === tab.tab ? "page" : undefined}
            className={cn(
              "relative shrink-0 cursor-pointer px-4 py-2.5 text-sm font-medium transition-colors",
              activeTab === tab.tab
                ? "text-primary"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            {tab.name}
            {activeTab === tab.tab && (
              <span className="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary" />
            )}
          </Link>
        ))}
      </nav>
      {currentNav?.children && currentNav.children.length > 0 && (
        <nav aria-label="首页子频道" className="flex flex-wrap gap-x-5 gap-y-2 rounded-md bg-secondary/40 px-4 py-3">
          {currentNav.children.map((child) => (
            <Link
              key={`${child.uri}-${child.name}`}
              href={child.uri}
              className="text-sm text-muted-foreground transition-colors hover:text-primary"
            >
              {child.name}
            </Link>
          ))}
        </nav>
      )}
    </div>
  )
}
