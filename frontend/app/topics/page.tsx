import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TabFilter } from "@/components/tab-filter"
import { TopicList } from "@/components/topic-list"
import { PenSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

export const metadata = {
  title: "主题讨论 - Go语言中文网",
  description: "Go语言中文社区主题讨论，分享技术经验，交流开发心得",
}

export default function TopicsPage() {
  return (
    <PageLayout>
      <PageHeader
        title="主题讨论"
        description="分享你的技术见解，和 Gopher 们一起交流成长"
        breadcrumbs={[{ label: "主题讨论" }]}
        actions={
          <Link href="/topics/new">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <PenSquare className="h-3.5 w-3.5" />
              {"发布主题"}
            </Button>
          </Link>
        }
      />

      <div className="rounded-lg border border-border bg-card p-4">
        <TabFilter />
        <div className="mt-4">
          <TopicList />
        </div>

        {/* Pagination */}
        <div className="mt-6 flex items-center justify-center gap-2">
          <button
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
            disabled
          >
            {"上一页"}
          </button>
          {[1, 2, 3, 4, 5].map((p) => (
            <button
              key={p}
              className={`rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                p === 1
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
            >
              {p}
            </button>
          ))}
          <span className="px-1 text-sm text-muted-foreground">...</span>
          <button className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80">
            {"下一页"}
          </button>
        </div>
      </div>
    </PageLayout>
  )
}
