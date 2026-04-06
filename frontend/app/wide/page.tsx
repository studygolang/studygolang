import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"

export const metadata: Metadata = {
  title: "Go Playground - Go语言中文网",
  description: "在线运行 Go 代码，随时随地编写和测试 Go 程序",
}

export default function WidePage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go Playground"
        description="在线编写和运行 Go 代码"
        breadcrumbs={[{ label: "Playground" }]}
      />

      <div className="rounded-lg border border-border overflow-hidden">
        <iframe
          src="https://go.dev/play"
          width="100%"
          height="700"
          className="border-0"
          title="Go Playground"
          loading="lazy"
        />
      </div>
    </PageLayout>
  )
}
