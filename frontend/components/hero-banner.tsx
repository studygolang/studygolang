"use client"

import Link from "next/link"
import { useEffect, useState } from "react"
import { Sparkles, Code2, Users, BookOpen } from "lucide-react"
import type { Announcement, SiteStats } from "@/lib/types"
import { announcementAPI } from "@/lib/api"
import { formatNum } from "@/lib/utils"

interface HeroBannerProps {
  stats?: SiteStats
}

// 默认公告（当 API 无数据或出错时显示）
const DEFAULT_ANNOUNCEMENTS: Announcement[] = [
  {
    id: 0,
    title: "欢迎来到 Go 语言中文网",
    content: "中国最大的 Go 语言社区，与全国 Gopher 一起学习成长",
    type: 1,
    priority: 0,
    start_time: "",
    end_time: "",
    created_at: "",
  },
]

export function HeroBanner({ stats }: HeroBannerProps) {
  const [announcements, setAnnouncements] = useState<Announcement[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false

    async function loadAnnouncements() {
      try {
        // 获取 type=1（公告）的列表
        const data = await announcementAPI.getList({ type: 1 })
        if (!cancelled) {
          // 按 priority 降序、created_at 降序排序
          const sorted = [...(data?.list ?? [])].sort(
            (a, b) =>
              b.priority - a.priority ||
              new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
          )
          setAnnouncements(sorted.slice(0, 3))
        }
      } catch {
        if (!cancelled) {
          setAnnouncements([])
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    loadAnnouncements()
    return () => {
      cancelled = true
    }
  }, [])

  // 最终展示：API 返回优先，出错或空时用默认值
  const displayAnnouncements = loading
    ? DEFAULT_ANNOUNCEMENTS
    : announcements.length > 0
      ? announcements
      : DEFAULT_ANNOUNCEMENTS

  // 取第一条公告作为主标语
  const primary = displayAnnouncements[0]

  return (
    <div className="border-b border-border bg-card">
      <div className="mx-auto max-w-7xl px-4 py-5 lg:px-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          {/* 公告区域 */}
          <div className="flex items-center gap-3">
            <Sparkles className="h-4 w-4 shrink-0 text-primary" />
            {loading ? (
              <span className="h-5 w-48 animate-pulse rounded bg-muted" />
            ) : primary ? (
              <Link
                href="/announcements"
                className="text-sm font-medium text-foreground transition-colors hover:text-primary"
                title={primary.title !== primary.content ? primary.title : undefined}
              >
                {primary.content || primary.title}
              </Link>
            ) : (
              <span className="text-sm font-medium text-muted-foreground">暂无公告</span>
            )}
          </div>

          {/* 站点统计 */}
          <div className="hidden items-center gap-5 sm:flex">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Users className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.user ?? 0)} 会员</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <BookOpen className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.topic ?? 0)} 主题</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Code2 className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.project ?? 0)} 项目</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
