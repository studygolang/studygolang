import type { Metadata } from "next"
import Link from "next/link"
import { BookOpen, Clock } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import type { Wiki } from "@/lib/types"

interface WikiListData {
  wikis: Wiki[]
  page: {
    has_prev: boolean
    prev_id: number
    has_next: boolean
    next_id: number
  }
}

export const metadata: Metadata = {
  title: "Wiki - Go语言中文网",
  description: "Go语言中文网 Wiki，汇集 Go 语言相关知识文档",
}

async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const base = process.env.API_BASE_URL || "http://localhost:8090"
    const res = await fetch(`${base}/api/v1${path}`, options)
    if (!res.ok) return null
    const json = await res.json()
    return json.code === 0 ? json.data : null
  } catch {
    return null
  }
}

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString("zh-CN", {
      year: "numeric",
      month: "long",
      day: "numeric",
    })
  } catch {
    return dateStr
  }
}

export default async function WikiListPage() {
  // 后端返回 { wikis: [...], page: {...} }，需要取 wikis 字段
  const wikisData = await fetchAPI<WikiListData>("/wiki", { cache: "no-store" })

  const wikiList: Wiki[] = wikisData?.wikis ?? []

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Wiki"
        description="Go 语言相关知识文档，持续更新中"
        breadcrumbs={[{ label: "Wiki" }]}
      />

      {wikiList.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">暂无 Wiki 文档</p>
        </div>
      ) : (
        <div className="rounded-lg border border-border bg-card">
          <ul className="divide-y divide-border">
            {wikiList.map((wiki) => (
              <li key={wiki.id}>
                <Link
                  href={`/wiki/${wiki.uri}`}
                  className="group flex items-start gap-4 px-5 py-4 transition-colors hover:bg-secondary/30"
                >
                  <div className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-primary/10">
                    <BookOpen className="h-4 w-4 text-primary" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <h2 className="text-sm font-semibold leading-snug text-foreground transition-colors group-hover:text-primary">
                      {wiki.title}
                    </h2>
                    <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                      {wiki.author && (
                        <span>{wiki.author}</span>
                      )}
                      <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {formatDate(wiki.mtime || wiki.ctime)}
                      </span>
                    </div>
                  </div>
                  <span className="mt-1 shrink-0 text-xs text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100">
                    查看 →
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}
    </PageLayout>
  )
}
