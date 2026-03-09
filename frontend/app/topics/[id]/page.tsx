import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { TopicDetail } from "@/components/topic-detail"
import type { TopicDetailData } from "@/lib/types"

async function fetchFromAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const base = process.env.API_BASE_URL || 'http://localhost:8088'
  const res = await fetch(`${base}/api/v1${path}`, options)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  return json.data
}

async function getTopicDetail(id: string): Promise<TopicDetailData | null> {
  try {
    return await fetchFromAPI<TopicDetailData>(
      `/topics/${id}`,
      { next: { revalidate: 60 } }
    )
  } catch {
    return null
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>
}): Promise<Metadata> {
  const { id } = await params
  const data = await getTopicDetail(id)

  if (!data?.topic) {
    return {
      title: "话题不存在 - Go语言中文网",
    }
  }

  const { topic } = data
  const description = topic.content?.slice(0, 150) || "Go语言中文社区话题讨论"

  return {
    title: `${topic.title} - Go语言中文网`,
    description,
    openGraph: {
      title: topic.title,
      description,
      type: "article",
      publishedTime: topic.ctime,
      authors: [topic.name],
    },
  }
}

export default async function TopicDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const data = await getTopicDetail(id)

  const topic = data?.topic
  const replies = data?.replies ?? []

  // JSON-LD 结构化数据
  const jsonLd = topic
    ? {
        "@context": "https://schema.org",
        "@type": "Article",
        headline: topic.title,
        description: topic.content?.slice(0, 150),
        author: {
          "@type": "Person",
          name: topic.name,
        },
        datePublished: topic.ctime,
        dateModified: topic.mtime || topic.ctime,
        publisher: {
          "@type": "Organization",
          name: "Go语言中文网",
          url: "https://studygolang.com",
        },
        commentCount: topic.replynum,
      }
    : null

  return (
    <PageLayout>
      {jsonLd && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      )}
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "主题讨论", href: "/topics" },
          { label: topic ? topic.title.slice(0, 20) + (topic.title.length > 20 ? "..." : "") : `#${id}` },
        ]}
      />
      <TopicDetail id={id} topic={topic} replies={replies} />
    </PageLayout>
  )
}
