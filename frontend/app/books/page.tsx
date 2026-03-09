import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { BookGrid } from "@/components/book-grid"

export const metadata = {
  title: "Go 图书 - Go语言中文网",
  description: "Go语言经典图书推荐，从入门到精通",
}

export default function BooksPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go 图书"
        description="精选 Go 语言经典图书，助你系统学习"
        breadcrumbs={[{ label: "Go 图书" }]}
      />
      <BookGrid />
    </PageLayout>
  )
}
