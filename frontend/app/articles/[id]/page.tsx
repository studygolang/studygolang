import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ArticleDetail } from "@/components/article-detail"
import type { ArticleDetailData } from "@/lib/types"

import { fetchAPI } from "@/lib/api"
import { sanitizeHtml } from "@/lib/sanitize"

async function getArticleDetail(id: string): Promise<ArticleDetailData | null> {
  try {
    return await fetchAPI<ArticleDetailData>(
      `/articles/${id}`,
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
  const data = await getArticleDetail(id)

  if (!data?.article) {
    return {
      title: "文章不存在 - Go语言中文网",
    }
  }

  const { article } = data
  // 后端 Article 无 summary 字段，截取 content 前 150 字作描述
  const description = article.content?.slice(0, 150) || "Go语言技术文章"

  return {
    title: `${article.title} - Go语言中文网`,
    description,
    openGraph: {
      title: article.title,
      description,
      type: "article",
      publishedTime: article.ctime,
      authors: [article.author],
      tags: article.tags ? article.tags.split(",") : [],
    },
  }
}

export default async function ArticleDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const data = await getArticleDetail(id)

  const article = data?.article
  // 后端返回 prev_next 数组（[prev, next]），replies 为评论列表
  const prevNext = data?.prev_next ?? []
  const prev = prevNext[0]
  const next = prevNext[1]
  const comments = data?.replies ?? []

  // JSON-LD 结构化数据
  const jsonLd = article
    ? {
        "@context": "https://schema.org",
        "@type": "Article",
        headline: article.title,
        // 后端 Article 无 summary 字段，截取 content 作描述
        description: article.content?.slice(0, 150),
        author: {
          "@type": "Person",
          name: article.author,
        },
        datePublished: article.ctime,
        dateModified: article.ctime,  // 后端 Article mtime 不序列化，使用 ctime
        publisher: {
          "@type": "Organization",
          name: "Go语言中文网",
          url: "https://studygolang.com",
        },
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
          { label: "技术文章", href: "/articles" },
          { label: article?.title ? article.title.slice(0, 20) + (article.title.length > 20 ? "..." : "") : `#${id}` },
        ]}
      />
      <ArticleDetail
        id={id}
        article={article}
        prev={prev}
        next={next}
        comments={comments}
      />
    </PageLayout>
  )
}
