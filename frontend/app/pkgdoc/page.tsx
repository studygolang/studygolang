"use client"

import { useEffect, useState } from "react"
import { ExternalLink, Loader2 } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"

// 旧版 /pkgdoc 用 frameset 内嵌本地部署的「Go 标准库文档中文版」
// （static/pkgdoc 为部署资产，不在仓库中）。pkg.go.dev 设置了
// x-frame-options: deny，无法内嵌，因此恢复旧版的本地双栏文档。
const INDEX_URL = "/static/pkgdoc/i.html"
const MAIN_URL = "/static/pkgdoc/main.html"

export default function PkgDocPage() {
  // null = 检测中；true/false = 本地文档资产是否可用
  const [available, setAvailable] = useState<boolean | null>(null)

  useEffect(() => {
    let cancelled = false
    // 代理层不支持 HEAD（405），用 GET 探测；文件仅 ~10KB
    fetch(INDEX_URL, { cache: "no-store" })
      .then((res) => {
        if (!cancelled) setAvailable(res.ok)
      })
      .catch(() => {
        if (!cancelled) setAvailable(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go 包文档"
        description="Go 语言标准库文档中文版"
        breadcrumbs={[{ label: "包文档" }]}
      />

      {available === null && (
        <div className="flex h-[calc(100dvh-13rem)] min-h-[600px] items-center justify-center rounded-lg border border-border">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      )}

      {available === false && (
        <div className="rounded-lg border border-border bg-muted/50 p-6 text-sm text-muted-foreground">
          <p className="mb-3">本地标准库文档未部署，可通过以下站点查阅：</p>
          <ul className="space-y-2">
            <li>
              <a
                href="https://studygolang.com/pkgdoc"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 text-primary hover:underline"
              >
                Go 标准库文档中文版（studygolang.com）
                <ExternalLink className="h-3.5 w-3.5" />
              </a>
            </li>
            <li>
              <a
                href="https://pkg.go.dev"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 text-primary hover:underline"
              >
                Go 官方包文档（pkg.go.dev）
                <ExternalLink className="h-3.5 w-3.5" />
              </a>
            </li>
          </ul>
        </div>
      )}

      {available === true && (
        <>
          {/* 撑满剩余视口，接近旧版整页 frameset 的体验 */}
          <div className="flex h-[calc(100dvh-13rem)] min-h-[600px] overflow-hidden rounded-lg border border-border">
            <iframe
              src={INDEX_URL}
              className="h-full w-1/5 min-w-[180px] border-0 border-r border-border"
              title="标准库文档目录"
            />
            <iframe
              src={MAIN_URL}
              name="main"
              className="h-full flex-1 border-0"
              title="标准库文档内容"
            />
          </div>
        </>
      )}
    </PageLayout>
  )
}
