import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { AnnouncementList } from "@/components/announcement/announcement-card"
import Link from "next/link"
import { announcementAPI } from "@/lib/api"
import {
  ANNOUNCEMENT_ALL_TYPE_ID,
  ANNOUNCEMENT_TYPE_FILTERS,
} from "@/components/announcement/announcement-type"

export const metadata: Metadata = {
  title: "公告中心 - Go语言中文网",
  description: "Go语言中文社区公告、活动和重要通知",
  openGraph: {
    title: "公告中心 - Go语言中文网",
    description: "Go语言中文社区公告、活动和重要通知",
    type: "website",
  },
}

interface AnnouncementsPageProps {
  searchParams: Promise<{ p?: string; type?: string }>
}

async function getAnnouncementsData(page: number, type: number | undefined) {
  try {
    const data = await announcementAPI.getList(
      { p: page, type },
      { next: { revalidate: 60 } }
    )
    return data
  } catch (error) {
    console.error("Failed to fetch announcements:", error)
    return { list: [], total: 0, has_more: false }
  }
}

const TYPE_FILTERS = ANNOUNCEMENT_TYPE_FILTERS

export default async function AnnouncementsPage({ searchParams }: AnnouncementsPageProps) {
  const { p: pageStr, type: typeStr } = await searchParams
  const page = Math.max(1, parseInt(pageStr ?? "1", 10) || 1)
  const typeParam = typeStr ? parseInt(typeStr, 10) : undefined
  const type = typeParam && [1, 2, 3].includes(typeParam) ? typeParam : undefined

  const announcementsData = await getAnnouncementsData(page, type)
  const announcements = announcementsData.list ?? []
  const total = announcementsData.total ?? 0
  const hasMore = announcementsData.has_more ?? false

  // 计算分页
  const pageSize = 20
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : (hasMore ? page + 1 : page)
  const pageNumbers: number[] = []
  const start = Math.max(1, page - 2)
  const end = Math.min(totalPages, page + 2)
  for (let i = start; i <= end; i++) {
    pageNumbers.push(i)
  }

  return (
    <PageLayout>
      <PageHeader
        title="公告中心"
        description="查看社区公告、活动和重要通知"
        breadcrumbs={[{ label: "公告中心" }]}
      />

      <div className="rounded-lg border border-border bg-card p-4">
        {/* Type Filter */}
        <div className="space-y-3">
          <div className="flex items-center gap-1 border-b border-border">
            {TYPE_FILTERS.map((filter) => (
              <Link
                key={filter.id}
                href={`/announcements${filter.id === ANNOUNCEMENT_ALL_TYPE_ID ? "" : `?type=${filter.id}`}`}
                className={`relative px-4 py-2.5 text-sm font-medium transition-colors ${
                  (type ?? ANNOUNCEMENT_ALL_TYPE_ID) === filter.id
                    ? "text-primary"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {filter.label}
                {(type ?? ANNOUNCEMENT_ALL_TYPE_ID) === filter.id && (
                  <span className="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary" />
                )}
              </Link>
            ))}
          </div>
        </div>

        <div className="mt-4">
          <AnnouncementList announcements={announcements} />
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="mt-6 flex items-center justify-center gap-2">
            {page > 1 ? (
              <Link
                href={`/announcements${type ? `?type=${type}` : ""}${page > 2 ? `&p=${page - 1}` : ""}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                上一页
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                上一页
              </button>
            )}
            {pageNumbers.map((p) => (
              <Link
                key={p}
                href={`/announcements${type ? `?type=${type}` : ""}${p > 1 ? `&p=${p}` : ""}`}
                className={`rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                  p === page
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
                }`}
              >
                {p}
              </Link>
            ))}
            {end < totalPages && (
              <>
                <span className="px-1 text-sm text-muted-foreground">...</span>
                <Link
                  href={`/announcements${type ? `?type=${type}` : ""}&p=${totalPages}`}
                  className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                >
                  {totalPages}
                </Link>
              </>
            )}
            {page < totalPages ? (
              <Link
                href={`/announcements${type ? `?type=${type}` : ""}&p=${page + 1}`}
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
              >
                下一页
              </Link>
            ) : (
              <button
                className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
                disabled
              >
                下一页
              </button>
            )}
          </div>
        )}
      </div>
    </PageLayout>
  )
}