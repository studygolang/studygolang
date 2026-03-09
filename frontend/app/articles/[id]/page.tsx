import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { ArticleDetail } from "@/components/article-detail"
import type { ArticleDetailData } from "@/lib/types"

async function fetchFromAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const base = process.env.API_BASE_URL || 'http://localhost:8088'
  const res = await fetch(`${base}/api/v1${path}`, options)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  return json.data
}

async function getArticleDetail(id: string): Promise<ArticleDetailData | null> {
  try {
    return await fetchFromAPI<ArticleDetailData>(
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
  const description = article.summary || article.content?.slice(0, 150) || "Go语言技术文章"

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
  const prev = data?.prev
  const next = data?.next
  const comments = data?.comments ?? []

  // JSON-LD 结构化数据
  const jsonLd = article
    ? {
        "@context": "https://schema.org",
        "@type": "Article",
        headline: article.title,
        description: article.summary || article.content?.slice(0, 150),
        author: {
          "@type": "Person",
          name: article.author,
        },
        datePublished: article.ctime,
        dateModified: article.mtime || article.ctime,
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
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      )}
      <PageHeader
        title=""
        breadcrumbs={[
          { label: "技术文章", href: "/articles" },
          { label: article ? article.title.slice(0, 20) + (article.title.length > 20 ? "..." : "") : `#${id}` },
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
