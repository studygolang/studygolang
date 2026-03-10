"use client"

import { useState } from "react"
import Link from "next/link"
import { ArrowRight, Sparkles, Code2, Users, BookOpen } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import type { SiteStats } from "@/lib/types"

// TODO: 公告内容应从后端 API 获取，以下为临时占位内容
const announcements = [
  {
    badge: "新版发布",
    title: "Go 1.26 RC1 已发布，新特性抢先看",
    href: "/topics", // 临时指向话题列表，待后端提供公告 API 后替换为真实话题链接
  },
  {
    badge: "社区招募",
    title: "寻找社区日常运营、功能开发、维护志愿者",
    href: "/topics", // 临时指向话题列表，待后端提供公告 API 后替换为真实话题链接
  },
]

interface HeroBannerProps {
  stats?: SiteStats
}

function formatNum(n: number | undefined): string {
  if (n == null) return "—"
  if (n >= 10000) return (n / 10000).toFixed(1) + "万"
  if (n >= 1000) return (n / 1000).toFixed(1) + "k"
  return String(n)
}

export function HeroBanner({ stats }: HeroBannerProps) {
  const [activeAnnouncement, setActiveAnnouncement] = useState(0)

  return (
    <div className="border-b border-border bg-card">
      <div className="mx-auto max-w-7xl px-4 py-5 lg:px-6">
        {/* Announcement bar */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-3 overflow-hidden">
            <Sparkles className="h-4 w-4 shrink-0 text-primary" />
            <div className="flex items-center gap-2 overflow-hidden">
              <Badge variant="secondary" className="shrink-0 bg-primary/10 text-xs font-semibold text-primary">
                {announcements[activeAnnouncement].badge}
              </Badge>
              <Link
                href={announcements[activeAnnouncement].href}
                className="truncate text-sm font-medium text-foreground transition-colors hover:text-primary"
              >
                {announcements[activeAnnouncement].title}
              </Link>
              <ArrowRight className="h-3.5 w-3.5 shrink-0 text-primary" />
            </div>
          </div>

          {/* Quick stats from real API */}
          <div className="hidden items-center gap-5 sm:flex">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Users className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.user)} 会员</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <BookOpen className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.topic)} 主题</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Code2 className="h-3.5 w-3.5" />
              <span>{formatNum(stats?.project)} 项目</span>
            </div>
          </div>
        </div>

        {/* Announcement dots */}
        {announcements.length > 1 && (
          <div className="mt-2 flex items-center gap-1.5">
            {announcements.map((_, i) => (
              <button
                key={i}
                onClick={() => setActiveAnnouncement(i)}
                className={`h-1 cursor-pointer rounded-full transition-all ${
                  i === activeAnnouncement ? "w-4 bg-primary" : "w-1 bg-border hover:bg-muted-foreground"
                }`}
                aria-label={`公告 ${i + 1}`}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
