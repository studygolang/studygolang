import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicDetail } from "@/components/topic-detail"

export const metadata = {
  title: "主题详情 - Go语言中文网",
}

export default async function TopicDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  return (
    <PageLayout>
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "主题讨论", href: "/topics" },
          { label: `#${id}` },
        ]}
      />
      <TopicDetail id={id} />
    </PageLayout>
  )
}
