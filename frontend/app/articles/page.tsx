import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ArticleList } from "@/components/article-list"
import { ArticleFilter } from "@/components/article-filter"
import { PenSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

export const metadata = {
  title: "技术文章 - Go语言中文网",
  description: "Go语言高质量技术文章，涵盖基础教程、实战经验、源码解析等",
}

export default function ArticlesPage() {
  return (
    <PageLayout>
      <PageHeader
        title="技术文章"
        description="精选 Go 语言技术文章，深度学习与实践"
        breadcrumbs={[{ label: "技术文章" }]}
        actions={
          <Link href="/articles/new">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <PenSquare className="h-3.5 w-3.5" />
              {"投稿"}
            </Button>
          </Link>
        }
      />

      <ArticleFilter />
      <div className="mt-4">
        <ArticleList />
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
    </PageLayout>
  )
}
