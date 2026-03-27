import { Suspense } from "react"
import Link from "next/link"
import { Plus } from "lucide-react"
import { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Button } from "@/components/ui/button"
import { JobList } from "@/components/job-list"
import type { JobListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "酷工作 - Go语言中文网",
  description: "Go 语言相关职位招聘信息，连接 Gopher 与优质企业",
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

async function JobItems({ page }: { page: number }) {
  const data = await fetchAPI<JobListData>(`/jobs?p=${page}`, { cache: "no-store" })
  const jobs = data?.jobs ?? []
  const hasMore = data?.has_more ?? false
  const total = data?.total ?? 0

  return (
    <div>
      <JobList jobs={jobs} />
      {/* Pagination */}
      {total > 0 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          {page > 1 && (
            <Link
              href={`/jobs?p=${page - 1}`}
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
            >
              上一页
            </Link>
          )}
          <span className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground">
            {page}
          </span>
          {hasMore && (
            <Link
              href={`/jobs?p=${page + 1}`}
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
            >
              下一页
            </Link>
          )}
        </div>
      )}
    </div>
  )
}

interface JobsPageProps {
  searchParams: Promise<{ page?: string }>
}

export default async function JobsPage({ searchParams }: JobsPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.page || "1", 10))

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="酷工作"
        description="发现 Go 语言相关的优质工作机会"
        breadcrumbs={[{ label: "酷工作" }]}
        actions={
          <Button
            size="sm"
            disabled
            className="gap-1.5"
            title="职位发布功能即将开放"
          >
            <Plus className="h-3.5 w-3.5" />
            {"发布职位"}
          </Button>
        }
      />
      <Suspense
        fallback={
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-24 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <JobItems page={page} />
      </Suspense>
    </PageLayout>
  )
}
