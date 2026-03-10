/**
 * 酷工作页面
 *
 * 注意：后端目前没有招聘（jobs）相关 API。
 * 页面内容由 JobList 组件提供，使用 Mock 数据展示。
 *
 * TODO: 后端实现 GET /api/v1/jobs 接口后，将此页面改为 SSR，从 API 获取真实数据。
 * TODO: "发布职位" 按钮指向 /jobs/new，该页面尚未实现，接入后端 API 后同步创建。
 */
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
