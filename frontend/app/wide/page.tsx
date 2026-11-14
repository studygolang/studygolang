import type { Metadata } from "next"
import { ExternalLink } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"

export const metadata: Metadata = {
  title: "Go Playground - Go语言中文网",
  description: "在线运行 Go 代码，随时随地编写和测试 Go 程序",
}

// 旧版 /wide/playground 内嵌的是 Wide（b3log 的在线 Go IDE）。
// go.dev/play 设置了 frame-ancestors 'self'，无法被第三方内嵌。
const PLAYGROUND_URL = "https://wide.b3log.org/playground"

export default function WidePage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go Playground"
        description="在线编写和运行 Go 代码"
        breadcrumbs={[{ label: "Playground" }]}
      />

      <div className="rounded-lg border border-border overflow-hidden">
        <iframe
          src={PLAYGROUND_URL}
          width="100%"
          height="700"
          className="border-0"
          title="Go Playground"
          sandbox="allow-scripts allow-same-origin allow-forms allow-popups allow-modals allow-downloads"
        />
      </div>

      <div className="mt-3 text-center text-sm text-muted-foreground">
        如果 Playground 无法加载，请{" "}
        <a
          href={PLAYGROUND_URL}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1 text-primary hover:underline"
        >
          直接访问 Wide
          <ExternalLink className="h-3.5 w-3.5" />
        </a>
      </div>
    </PageLayout>
  )
}
