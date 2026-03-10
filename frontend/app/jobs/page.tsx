import { Suspense } from "react"
import Link from "next/link"
import { Plus, Briefcase } from "lucide-react"
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

  return <JobList jobs={jobs} />
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
          // TODO: /jobs/new 发布职位页面尚未实现
          <a href="#">
            <Button
              size="sm"
              className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90"
            >
              <Plus className="h-3.5 w-3.5" />
              {"发布职位"}
            </Button>
          </a>
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
