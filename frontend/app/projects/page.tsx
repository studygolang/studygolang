import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ProjectGrid } from "@/components/project-grid"
import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

export const metadata = {
  title: "开源项目 - Go语言中文网",
  description: "发现优秀的 Go 语言开源项目，分享你的作品",
}

export default function ProjectsPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="开源项目"
        description="发现和分享优秀的 Go 语言开源项目"
        breadcrumbs={[{ label: "开源项目" }]}
        actions={
          <Link href="/projects/new">
            <Button size="sm" className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
              <Plus className="h-3.5 w-3.5" />
              {"提交项目"}
            </Button>
          </Link>
        }
      />
      <ProjectGrid />
    </PageLayout>
  )
}
