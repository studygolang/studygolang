"use client"

import Link from "next/link"
import { Sparkles, Code2, Users, BookOpen } from "lucide-react"
import type { SiteStats } from "@/lib/types"

// TODO: 后端实现 GET /api/v1/announcements 接口后，接入真实公告数据

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
  return (
    <div className="border-b border-border bg-card">
      <div className="mx-auto max-w-7xl px-4 py-5 lg:px-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          {/* Brand slogan */}
          <div className="flex items-center gap-3">
            <Sparkles className="h-4 w-4 shrink-0 text-primary" />
            <Link
              href="/topics"
              className="text-sm font-medium text-foreground transition-colors hover:text-primary"
            >
              {"中国最大的 Go 语言社区，与全国 Gopher 一起学习成长"}
            </Link>
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
      </div>
    </div>
  )
}
