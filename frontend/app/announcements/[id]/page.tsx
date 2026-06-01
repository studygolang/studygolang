import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { ArrowLeft, Calendar, Clock } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Badge } from "@/components/ui/badge"
import { announcementAPI } from "@/lib/api"
import type { Announcement } from "@/lib/types"
import { getAnnouncementType } from "@/components/announcement/announcement-type"

function formatDate(timeStr: string): string {
  try {
    return new Date(timeStr).toLocaleString("zh-CN", { hour12: false })
  } catch {
    return timeStr
  }
}

async function getAnnouncement(id: number): Promise<Announcement | null> {
  try {
    const data = await announcementAPI.getDetail(id, { next: { revalidate: 60 } })
    return data.announcement ?? null
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
  const announcementId = parseInt(id, 10)
  if (!announcementId || Number.isNaN(announcementId)) {
    return { title: "公告不存在 - Go语言中文网" }
  }
  const announcement = await getAnnouncement(announcementId)
  if (!announcement) {
    return { title: "公告不存在 - Go语言中文网" }
  }
  const description = announcement.content?.slice(0, 150) || announcement.title
  return {
    title: `${announcement.title} - Go语言中文网`,
    description,
    openGraph: {
      title: announcement.title,
      description,
      type: "article",
    },
  }
}

export default async function AnnouncementDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const announcementId = parseInt(id, 10)
  if (!announcementId || Number.isNaN(announcementId)) {
    notFound()
  }

  const announcement = await getAnnouncement(announcementId)
  if (!announcement) {
    notFound()
  }

  const typeConfig = getAnnouncementType(announcement.type)

  return (
    <PageLayout>
      <PageHeader
        title={announcement.title}
        breadcrumbs={[
          { label: "公告中心", href: "/announcements" },
          { label: announcement.title },
        ]}
      />

      <article className="rounded-lg border border-border bg-card p-6">
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Badge variant={typeConfig.variant}>{typeConfig.label}</Badge>
          </div>

          <h1 className="text-2xl font-semibold leading-tight text-foreground">
            {announcement.title}
          </h1>

          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted-foreground">
            <span className="flex items-center gap-1">
              <Clock className="h-3.5 w-3.5" />
              发布于 {formatDate(announcement.created_at)}
            </span>
            {announcement.start_time && announcement.end_time && (
              <span className="flex items-center gap-1">
                <Calendar className="h-3.5 w-3.5" />
                有效期 {announcement.start_time.slice(0, 10)} ~ {announcement.end_time.slice(0, 10)}
              </span>
            )}
          </div>

          <div className="border-t border-border pt-4">
            <div className="prose prose-sm max-w-none whitespace-pre-wrap text-foreground">
              {announcement.content}
            </div>
          </div>
        </div>
      </article>

      <div className="mt-6">
        <Link
          href="/announcements"
          className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-primary"
        >
          <ArrowLeft className="h-4 w-4" />
          返回公告列表
        </Link>
      </div>
    </PageLayout>
  )
}
