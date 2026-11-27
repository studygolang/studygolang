import type { Metadata } from "next"
import { cookies } from "next/headers"
import { PageLayout } from "@/components/page-layout"
import { toPlainExcerpt } from "@/lib/utils"
import { PageHeader } from "@/components/page-header"
import { TopicDetail } from "@/components/topic-detail"
import { ViewStats } from "@/components/view-stats"
import type { TopicDetailData, TopicAppend } from "@/lib/types"

import { fetchAPI } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"

// 话题有可见性（公开/登录可见/付费可见/仅自己可见），后端按当前登录用户判定。
// 必须动态渲染并透传 Cookie：ISR 会把匿名结果缓存给所有人（受限话题永远读不到），
// 也可能把付费用户的结果缓存给访客（内容泄露）。
export const dynamic = "force-dynamic"

async function authHeaders(): Promise<Record<string, string>> {
  const store = await cookies()
  const cookie = store
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ")
  return cookie ? { Cookie: cookie } : {}
}

async function getTopicAppends(id: string): Promise<TopicAppend[]> {
  try {
    const data = await fetchAPI<{ appends: TopicAppend[] }>(
      `/topics/${id}/appends`,
      { cache: "no-store", headers: await authHeaders() }
    )
    return data?.appends ?? []
  } catch {
    return []
  }
}

async function getTopicDetail(id: string): Promise<TopicDetailData | null> {
  try {
    return await fetchAPI<TopicDetailData>(
      `/topics/${id}`,
      { cache: "no-store", headers: await authHeaders() }
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
  const description = (topic.content ? toPlainExcerpt(topic.content) : "") || "Go语言中文社区话题讨论"

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
  const [data, appends] = await Promise.all([
    getTopicDetail(id),
    getTopicAppends(id),
  ])

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
          dangerouslySetInnerHTML={{ __html: sanitizeHtml(JSON.stringify(jsonLd)) }}
        />
      )}
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "主题讨论", href: "/topics" },
          { label: topic ? topic.title.slice(0, 20) + (topic.title.length > 20 ? "..." : "") : `#${id}` },
        ]}
      />
      <TopicDetail
        id={id}
        topic={topic}
        replies={replies}
        appends={appends}
        canView={data?.can_view ?? true}
        likeflag={data?.likeflag}
        hadcollect={data?.hadcollect}
        prevNext={data?.prev_next}
      />
      <ViewStats viewUserNum={data?.view_user_num} viewSource={data?.view_source} />
    </PageLayout>
  )
}
