import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { JobList } from "@/components/job-list"
import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

export const metadata = {
  title: "酷工作 - Go语言中文网",
  description: "Go 语言相关职位招聘信息，连接 Gopher 与优质企业",
}

export default function JobsPage() {
  return (
    <PageLayout
      sidebar={false}
    >
      <PageHeader
        title="酷工作"
        description="发现 Go 语言相关的优质工作机会"
        breadcrumbs={[{ label: "酷工作" }]}
        actions={
          <Link href="/jobs/new">
            <Button
              size="sm"
              className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90"
            >
              <Plus className="h-3.5 w-3.5" />
              {"发布职位"}
            </Button>
          </Link>
        }
      />
      <JobList />
    </PageLayout>
  )
}
