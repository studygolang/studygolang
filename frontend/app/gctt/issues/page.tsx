import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPI } from "@/lib/api"

export const metadata: Metadata = {
  title: "GCTT 选题列表 - Go语言中文网",
  description: "GCTT Go 中文翻译组待翻译和已翻译的选题列表",
}

interface GCTTIssue {
  id: number
  title: string
  translator: string
  label: string
  state: number
  translating_at: number
  translated_at: number
  created_at: string
}

async function getIssues(p: number, state?: string) {
  try {
    const params = new URLSearchParams()
    if (p > 1) params.set("p", String(p))
    if (state) params.set("state", state)
    const qs = params.toString()
    return await fetchAPI<{
      issues: GCTTIssue[]
      total: number
      page: number
      has_more: boolean
    }>(`/gctt/issues${qs ? "?" + qs : ""}`, { cache: "no-store" })
  } catch {
    return null
  }
}

export default async function GCTTIssuesPage({
  searchParams,
}: {
  searchParams: Promise<{ p?: string; state?: string }>
}) {
  const sp = await searchParams
  const page = Math.max(1, parseInt(sp.p ?? "1", 10) || 1)
  const data = await getIssues(page, sp.state)
  const issues = data?.issues ?? []
  const hasMore = data?.has_more ?? false

  return (
    <PageLayout>
      <PageHeader
        title="选题列表"
        breadcrumbs={[
          { label: "GCTT", href: "/gctt" },
          { label: "选题列表" },
        ]}
        actions={
          <div className="flex gap-2">
            <Link
              href="/gctt/issues"
              className={`rounded-md px-3 py-1.5 text-sm font-medium ${
                !sp.state ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
            >
              全部
            </Link>
            <Link
              href="/gctt/issues?state=open"
              className={`rounded-md px-3 py-1.5 text-sm font-medium ${
                sp.state === "open" ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
            >
              待翻译
            </Link>
            <Link
              href="/gctt/issues?state=closed"
              className={`rounded-md px-3 py-1.5 text-sm font-medium ${
                sp.state === "closed" ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
            >
              已完成
            </Link>
          </div>
        }
      />

      {issues.length > 0 ? (
        <div className="space-y-3">
          {issues.map((issue) => (
            <Card key={issue.id}>
              <CardContent className="p-4">
                <div className="font-medium">{issue.title}</div>
                <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
                  {issue.label && (
                    <span className="rounded bg-primary/10 px-1.5 py-0.5 text-primary">
                      {issue.label}
                    </span>
                  )}
                  <span className={issue.state === 0 ? "text-orange-500" : "text-green-500"}>
                    {issue.state === 0 ? "待翻译" : "已完成"}
                  </span>
                  {issue.translator && <span>译者: {issue.translator}</span>}
                  <span>{issue.created_at?.slice(0, 10)}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <p className="text-sm text-muted-foreground">暂无选题</p>
        </div>
      )}

      {/* 分页 */}
      <div className="mt-6 flex justify-center gap-2">
        {page > 1 && (
          <Link
            href={page > 2 ? `/gctt/issues?p=${page - 1}${sp.state ? "&state=" + sp.state : ""}` : `/gctt/issues${sp.state ? "?state=" + sp.state : ""}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            上一页
          </Link>
        )}
        {hasMore && (
          <Link
            href={`/gctt/issues?p=${page + 1}${sp.state ? "&state=" + sp.state : ""}`}
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            下一页
          </Link>
        )}
      </div>
    </PageLayout>
  )
}
