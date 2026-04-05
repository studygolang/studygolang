import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { Rss } from "lucide-react"

export const metadata: Metadata = {
  title: "RSS/Atom 订阅 - Go语言中文网",
  description: "订阅 Go语言中文网的最新内容，支持 RSS 和 Atom 格式",
}

const feedList = [
  {
    title: "Go语言中文网全站动态",
    url: "/feed.xml",
    format: "Atom",
    description: "包含话题、文章、资源等全站最新动态",
  },
]

export default function FeedPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="RSS/Atom 订阅"
        description="通过订阅获取 Go语言中文网最新内容"
        breadcrumbs={[{ label: "RSS 订阅" }]}
      />

      <Card>
        <CardContent className="p-6">
          <div className="mb-6 flex items-center gap-3">
            <Rss className="h-8 w-8 text-orange-500" />
            <div>
              <h2 className="text-lg font-semibold">订阅源</h2>
              <p className="text-sm text-muted-foreground">
                使用 RSS 阅读器订阅以下链接
              </p>
            </div>
          </div>

          <div className="space-y-4">
            {feedList.map((feed) => (
              <div
                key={feed.url}
                className="rounded-lg border border-border p-4"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">{feed.title}</h3>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {feed.description}
                    </p>
                  </div>
                  <span className="rounded-full bg-orange-100 px-2.5 py-0.5 text-xs font-medium text-orange-700 dark:bg-orange-900/30 dark:text-orange-400">
                    {feed.format}
                  </span>
                </div>
                <div className="mt-3 flex items-center gap-2">
                  <code className="flex-1 rounded bg-muted px-3 py-1.5 text-sm">
                    {feed.url}
                  </code>
                  <a
                    href={feed.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/90"
                  >
                    订阅
                  </a>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-8 rounded-lg bg-muted/50 p-4">
            <h3 className="font-medium">如何使用 RSS 订阅？</h3>
            <ul className="mt-2 space-y-1 text-sm text-muted-foreground">
              <li>1. 安装一个 RSS 阅读器（如 Feedly、Inoreader 等）</li>
              <li>2. 复制上方的订阅链接</li>
              <li>3. 在阅读器中添加订阅链接</li>
              <li>4. 即可自动获取最新内容更新</li>
            </ul>
          </div>
        </CardContent>
      </Card>
    </PageLayout>
  )
}
