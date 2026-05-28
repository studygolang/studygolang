import type { Metadata } from "next"
import Link from "next/link"
import { Plus } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Button } from "@/components/ui/button"
import { fetchAPINullable } from "@/lib/api"
import { ProjectsClient } from "./projects-client"
import type { ProjectListData } from "@/lib/types"

export const metadata: Metadata = {
  title: "Go 开源项目 - Go语言中文网",
  description: "发现优秀的 Go 语言开源项目，分享你的作品",
  openGraph: {
    title: "Go 开源项目 - Go语言中文网",
    description: "发现优秀的 Go 语言开源项目，分享你的作品",
    type: "website",
  },
}

interface ProjectsPageProps {
  searchParams: Promise<{ p?: string; sort?: string }>
}

export default async function ProjectsPage({ searchParams }: ProjectsPageProps) {
  const params = await searchParams
  const page = Math.max(1, parseInt(params.p || "1", 10))
  const sort = params.sort || "latest"

  // Fetch initial data for SSR + pass to client to avoid duplicate request
  const initialData = await fetchAPINullable<ProjectListData>(
    `/projects?p=${page}&sort=${sort}`,
    { cache: "no-store" }
  )
  const projects = initialData?.projects || []
  const total = initialData?.total || 0
  const hasMore = initialData?.has_more || false

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="开源项目"
        description="发现和分享优秀的 Go 语言开源项目"
        breadcrumbs={[{ label: "开源项目" }]}
        actions={
          <Link href="/publish">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <Plus className="h-3.5 w-3.5" />
              提交项目
            </Button>
          </Link>
        }
      />
      <ProjectsClient page={page} sort={sort} total={total} projects={projects} initialHasMore={hasMore} />
    </PageLayout>
  )
}
