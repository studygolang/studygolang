"use client"

import Link from "next/link"
import {
  Github,
  TrendingUp,
  Users,
  MessageSquare,
  BookOpen,
  Trophy,
  ArrowRight,
  Flame,
  Calendar,
  BarChart3,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import type { Reading, SiteStats, Topic, Comment, User, FriendLink } from "@/lib/types"

/* ---------- Login Card ---------- */
export function LoginCard() {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="mb-3 text-center">
          <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <Users className="h-6 w-6 text-primary" />
          </div>
          <p className="text-sm font-medium text-foreground">{"加入 Go 中文社区"}</p>
          <p className="mt-1 text-xs text-muted-foreground">{"与全国 Gopher 一起学习交流"}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" className="flex-1 gap-1.5 text-xs">
            <Github className="h-3.5 w-3.5" />
            GitHub
          </Button>
          <Button size="sm" className="flex-1 bg-primary text-xs text-primary-foreground hover:bg-primary/90">
            {"注册账号"}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Daily Interview Question ---------- */
export function DailyQuestion() {
  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Flame className="h-4 w-4 text-primary" />
          {"Go 今日面试题"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <Link
          href="/interview/today"
          className="group block rounded-md bg-secondary/50 p-3 transition-colors hover:bg-secondary"
        >
          <p className="text-sm font-medium leading-relaxed text-foreground group-hover:text-primary">
            {"Go 中 slice 和 array 的区别是什么？如何避免 slice 的内存泄漏？"}
          </p>
          <span className="mt-2 inline-flex items-center gap-1 text-xs font-medium text-primary">
            {"查看解答"}
            <ArrowRight className="h-3 w-3" />
          </span>
        </Link>
      </CardContent>
    </Card>
  )
}

/* ---------- Morning Reading ---------- */
interface MorningReadingProps {
  readings?: Reading[]
}

export function MorningReading({ readings }: MorningReadingProps) {
  if (!readings || readings.length === 0) return null

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Calendar className="h-4 w-4 text-primary" />
          {"今日晨读"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-2.5">
          {readings.slice(0, 7).map((r, i) => (
            <Link
              key={r.id}
              href={r.url || `/readings/${r.id}`}
              className="flex items-start gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
            >
              <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded bg-secondary text-xs font-semibold text-secondary-foreground">
                {i + 1}
              </span>
              <span className="line-clamp-2 leading-snug">{r.title}</span>
            </Link>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Trending Topics ---------- */
interface TrendingTopicsProps {
  topics?: Topic[]
}

export function TrendingTopics({ topics }: TrendingTopicsProps) {
  if (!topics || topics.length === 0) return null

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-sm font-semibold">
            <TrendingUp className="h-4 w-4 text-primary" />
            {"热门话题"}
          </CardTitle>
        </div>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-2.5">
          {topics.slice(0, 5).map((t, i) => {
            const views = t.viewnum >= 1000 ? (t.viewnum / 1000).toFixed(1) + "k" : String(t.viewnum)
            return (
              <Link
                key={t.tid}
                href={`/topics/${t.tid}`}
                className="flex items-start gap-2 text-sm transition-colors hover:text-primary"
              >
                <span
                  className={`mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded text-xs font-bold ${
                    i < 2 ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground"
                  }`}
                >
                  {i + 1}
                </span>
                <span className="line-clamp-2 flex-1 leading-snug text-foreground">{t.title}</span>
                <span className="shrink-0 text-xs text-muted-foreground">{views}</span>
              </Link>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Latest Comments ---------- */
interface LatestCommentsProps {
  comments?: Comment[]
}

function formatTime(ctime: string): string {
  try {
    const date = new Date(ctime)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    if (minutes < 60) return `${minutes} 分钟前`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours} 小时前`
    return `${Math.floor(hours / 24)} 天前`
  } catch {
    return ctime
  }
}

export function LatestComments({ comments }: LatestCommentsProps) {
  if (!comments || comments.length === 0) return null

  const displayComments = comments.slice(0, 3)

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <MessageSquare className="h-4 w-4 text-primary" />
          {"最新评论"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-3">
          {displayComments.map((comment, i) => (
            <div key={comment.id} className="space-y-1">
              <div className="flex items-center gap-2">
                <Avatar className="h-5 w-5">
                  <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                    {comment.name ? comment.name.charAt(0).toUpperCase() : "?"}
                  </AvatarFallback>
                </Avatar>
                <span className="text-xs font-medium text-foreground">{comment.name}</span>
                <span className="text-[10px] text-muted-foreground">{formatTime(comment.ctime)}</span>
              </div>
              <p className="line-clamp-1 pl-7 text-xs text-muted-foreground">{comment.content}</p>
              {i < displayComments.length - 1 && <Separator className="mt-2" />}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Active Members ---------- */
interface ActiveMembersProps {
  users?: User[]
}

export function ActiveMembers({ users }: ActiveMembersProps) {
  if (!users || users.length === 0) return null

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Trophy className="h-4 w-4 text-primary" />
          {"活跃会员"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="flex flex-wrap gap-2">
          {users.map((user) => {
            const displayName = user.name || user.username
            return (
              <Link
                key={user.uid}
                href={`/user/${user.username}`}
                className="group flex items-center gap-1.5 rounded-full bg-secondary py-1 pl-1 pr-2.5 transition-colors hover:bg-primary/10"
              >
                <Avatar className="h-5 w-5">
                  <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                    {displayName ? displayName.charAt(0).toUpperCase() : "?"}
                  </AvatarFallback>
                </Avatar>
                <span className="text-xs text-secondary-foreground group-hover:text-primary">{displayName}</span>
              </Link>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Stats Card ---------- */
interface StatsCardProps {
  stats?: SiteStats
}

function formatStatNum(num: number): string {
  if (num >= 10000) return (num / 10000).toFixed(1) + "万"
  if (num >= 1000) return (num / 1000).toFixed(1) + "k"
  return String(num)
}

export function StatsCard({ stats }: StatsCardProps) {
  // 注意：使用 ?? 而非 ? 避免 0 被当成 falsy 误显示 placeholder
  const displayStats = [
    { label: "会员数", value: stats?.user != null ? formatStatNum(stats.user) : "—", icon: Users },
    { label: "主题数", value: stats?.topic != null ? formatStatNum(stats.topic) : "—", icon: BookOpen },
    { label: "评论数", value: stats?.comment != null ? formatStatNum(stats.comment) : "—", icon: MessageSquare },
    { label: "文章数", value: stats?.article != null ? formatStatNum(stats.article) : "—", icon: BarChart3 },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="text-sm font-semibold">{"社区统计"}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="grid grid-cols-2 gap-3">
          {displayStats.map((stat) => (
            <div key={stat.label} className="rounded-md bg-secondary/50 p-2.5 text-center">
              <stat.icon className="mx-auto mb-1 h-4 w-4 text-primary" />
              <p className="text-sm font-bold text-foreground">{stat.value}</p>
              <p className="text-[10px] text-muted-foreground">{stat.label}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Friend Links ---------- */
interface FriendLinksProps {
  links?: FriendLink[]
}

export function FriendLinks({ links }: FriendLinksProps) {
  if (!links || links.length === 0) return null

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="text-sm font-semibold">{"友情链接"}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="flex flex-wrap gap-2">
          {links.map((link, i) => (
            <a
              key={link.id ?? i}
              href={link.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-muted-foreground transition-colors hover:text-primary"
            >
              {link.name}
            </a>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
