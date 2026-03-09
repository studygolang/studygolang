import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ArticleDetail } from "@/components/article-detail"

export const metadata = {
  title: "文章详情 - Go语言中文网",
}

export default async function ArticleDetailPage({
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
          { label: "技术文章", href: "/articles" },
          { label: `#${id}` },
        ]}
      />
      <ArticleDetail id={id} />
    </PageLayout>
  )
}
