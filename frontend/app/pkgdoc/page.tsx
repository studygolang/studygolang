import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ExternalLink } from "lucide-react"

export const metadata: Metadata = {
  title: "Go 包文档 - Go语言中文网",
  description: "浏览 Go 标准库及第三方包文档",
}

export default function PkgDocPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go 包文档"
        description="浏览 Go 标准库及第三方包文档"
        breadcrumbs={[{ label: "包文档" }]}
      />

      {/* 降级提示：iframe 无法加载时显示跳转链接 */}
      <noscript>
        <div className="mb-4 rounded-lg border border-border bg-muted/50 p-4 text-sm text-muted-foreground">
          无法加载内嵌文档，请{" "}
          <a
            href="https://pkg.go.dev"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-primary hover:underline"
          >
            直接访问 pkg.go.dev
            <ExternalLink className="h-3.5 w-3.5" />
          </a>
        </div>
      </noscript>

      <div className="rounded-lg border border-border overflow-hidden">
        <iframe
          src="https://pkg.go.dev"
          width="100%"
          height="700"
          className="border-0"
          title="Go 包文档"
          loading="lazy"
          sandbox="allow-scripts allow-same-origin allow-forms allow-popups allow-top-navigation"
        />
      </div>

      {/* 始终显示的直接访问链接（降级处理） */}
      <div className="mt-3 text-center text-sm text-muted-foreground">
        如果文档无法加载，请{" "}
        <a
          href="https://pkg.go.dev"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1 text-primary hover:underline"
        >
          直接访问 pkg.go.dev
          <ExternalLink className="h-3.5 w-3.5" />
        </a>
      </div>
    </PageLayout>
  )
}
