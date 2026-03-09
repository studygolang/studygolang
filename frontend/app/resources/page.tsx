import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ResourceList } from "@/components/resource-list"
import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"

export const metadata = {
  title: "资源索引 - Go语言中文网",
  description: "Go语言学习资源大全，教程、工具、视频、文档一站式索引",
}

export default function ResourcesPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="资源索引"
        description="精选 Go 语言学习资源，教程、工具、视频、文档一网打尽"
        breadcrumbs={[{ label: "资源索引" }]}
        actions={
          <Link href="/resources/new">
            <Button
              size="sm"
              className="gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90"
            >
              <Plus className="h-3.5 w-3.5" />
              {"分享资源"}
            </Button>
          </Link>
        }
      />
      <ResourceList />
    </PageLayout>
  )
}
