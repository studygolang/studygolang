import Link from "next/link"
import { Clock, Calendar } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import type { Announcement } from "@/lib/types"
import { getAnnouncementType } from "@/components/announcement/announcement-type"

interface AnnouncementCardProps {
  announcement: Announcement
}

function formatTime(timeStr: string): string {
  try {
    const date = new Date(timeStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(Math.abs(diff) / 60000)
    if (diff < 0) {
      if (minutes < 60) return `${minutes} 分钟后`
      const hours = Math.floor(minutes / 60)
      if (hours < 24) return `${hours} 小时后`
      return `${Math.floor(hours / 24)} 天后`
    }
    if (minutes === 0) return "刚刚"
    if (minutes < 60) return `${minutes} 分钟前`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours} 小时前`
    const days = Math.floor(hours / 24)
    if (days < 30) return `${days} 天前`
    return timeStr.slice(0, 10)
  } catch {
    return timeStr
  }
}

export function AnnouncementCard({ announcement }: AnnouncementCardProps) {
  const typeConfig = getAnnouncementType(announcement.type)

  return (
    <article className="group relative rounded-lg border border-border bg-card p-4 transition-colors hover:bg-secondary/30">
      <div className="space-y-3">
        {/* Header: Title + Badge */}
        <div className="flex items-start gap-2">
          <Badge variant={typeConfig.variant} className="shrink-0">
            {typeConfig.label}
          </Badge>
          <Link
            href={`/announcements/${announcement.id}`}
            className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
          >
            {announcement.title}
          </Link>
        </div>

        {/* Content Summary */}
        <p className="text-sm text-muted-foreground line-clamp-2">
          {announcement.content}
        </p>

        {/* Footer: Time + Duration */}
        <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
          <span className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {formatTime(announcement.created_at)}
          </span>
          {announcement.start_time && announcement.end_time && (
            <span className="flex items-center gap-1">
              <Calendar className="h-3 w-3" />
              {announcement.start_time.slice(0, 10)} ~ {announcement.end_time.slice(0, 10)}
            </span>
          )}
        </div>
      </div>
    </article>
  )
}

interface AnnouncementListProps {
  announcements?: Announcement[]
}

export function AnnouncementList({ announcements = [] }: AnnouncementListProps) {
  if (announcements.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无公告
      </div>
    )
  }

  return (
    <div className="grid gap-4">
      {announcements.map((announcement) => (
        <AnnouncementCard key={announcement.id} announcement={announcement} />
      ))}
    </div>
  )
}