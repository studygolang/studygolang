"use client"

// ResourceList 组件：仅用于展示，所有数据通过 props 从父组件（SSR 页面）传入
// 不使用本地 mock 数据，不做 API 调用

import Link from "next/link"
import {
  ExternalLink,
  ThumbsUp,
  MessageSquare,
  Clock,
  Tag,
  Download,
  Globe,
  BookOpen,
  Video,
  Wrench,
  GraduationCap,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn, formatNum } from "@/lib/utils"
import { useState } from "react"
import type { Resource } from "@/lib/types"

// 资源分类侧边栏（纯 UI，分类过滤通过 URL catid 参数与后端交互）
const CATEGORY_NAV = [
  { label: "全部", icon: Tag, catid: 0 },
  { label: "教程指南", icon: GraduationCap, catid: 1 },
  { label: "开发工具", icon: Wrench, catid: 2 },
  { label: "视频课程", icon: Video, catid: 3 },
  { label: "网站推荐", icon: Globe, catid: 4 },
  { label: "文档翻译", icon: BookOpen, catid: 5 },
  { label: "下载资源", icon: Download, catid: 6 },
]

interface ResourceListProps {
  resources: Resource[]
  /** 当前选中的分类 catid，0 表示全部 */
  activeCatid?: number
}

export function ResourceList({ resources, activeCatid = 0 }: ResourceListProps) {
  // 排序状态：最新/最热/推荐（纯 UI 展示，暂未接 API sort 参数）
  const [activeSort, setActiveSort] = useState("最新")

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      {/* Sidebar categories：点击通过 URL 参数切换分类，交由 SSR 页面处理 */}
      <aside className="w-full shrink-0 lg:w-56">
        <div className="rounded-lg border border-border bg-card p-3">
          <h3 className="mb-2 px-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            {"资源分类"}
          </h3>
          <nav className="space-y-0.5">
            {CATEGORY_NAV.map((cat) => (
              <Link
                key={cat.catid}
                href={cat.catid === 0 ? "/resources" : `/resources?catid=${cat.catid}`}
                className={cn(
                  "flex w-full cursor-pointer items-center gap-2.5 rounded-md px-2.5 py-2 text-left text-sm font-medium transition-colors",
                  activeCatid === cat.catid
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                )}
              >
                <cat.icon className="h-4 w-4 shrink-0" />
                <span className="flex-1">{cat.label}</span>
              </Link>
            ))}
          </nav>
        </div>
      </aside>

      {/* Resource list */}
      <div className="min-w-0 flex-1 space-y-3">
        {/* Result header */}
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            {"共 "}<span className="font-semibold text-foreground">{resources.length}</span>{" 个资源（本页）"}
          </p>
          {/* 排序按钮：TODO 接入 API sort 参数 */}
          <div className="flex items-center gap-2">
            {["最新", "最热", "推荐"].map((sort) => (
              <button
                key={sort}
                onClick={() => setActiveSort(sort)}
                className={cn(
                  "rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
                  activeSort === sort
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-secondary"
                )}
              >
                {sort}
              </button>
            ))}
          </div>
        </div>

        {resources.length === 0 ? (
          <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
            <ExternalLink className="mb-3 h-10 w-10 text-muted-foreground/40" />
            <p className="text-sm text-muted-foreground">暂无资源数据</p>
          </div>
        ) : (
          resources.map((resource) => (
            <Card
              key={resource.id}
              className="group transition-all hover:border-primary/20 hover:shadow-sm"
            >
              <CardContent className="p-4 sm:p-5">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    {/* Title */}
                    <div className="flex items-start gap-2">
                      {resource.catname && (
                        <Badge className="mt-0.5 shrink-0 bg-primary/10 text-[10px] font-semibold text-primary">
                          {resource.catname}
                        </Badge>
                      )}
                      <Link
                        href={`/resources/${resource.id}`}
                        className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                      >
                        {resource.title}
                      </Link>
                    </div>

                    {/* Description：后端 Resource 用 content 字段存描述 */}
                    {resource.content && (
                      <p className="mt-1.5 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                        {resource.content}
                      </p>
                    )}

                    {/* Tags + Meta */}
                    <div className="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                      {resource.tags && (
                        <div className="flex items-center gap-1.5">
                          {resource.tags.split(",").slice(0, 3).map((tag) => (
                            <Badge
                              key={tag}
                              variant="secondary"
                              className="text-[10px] font-medium"
                            >
                              {tag.trim()}
                            </Badge>
                          ))}
                        </div>
                      )}
                      {/* 作者信息：用 user.name 优先，降级用 uid */}
                      <span className="text-xs text-muted-foreground">
                        {resource.user?.name || resource.user?.username || `uid:${resource.uid}`}
                      </span>
                      <span className="flex items-center gap-1 text-xs text-muted-foreground">
                        <Clock className="h-3 w-3" />
                        {resource.ctime}
                      </span>
                      <span className="flex items-center gap-1 text-xs text-muted-foreground">
                        <ThumbsUp className="h-3 w-3" />
                        {formatNum(resource.likenum)}
                      </span>
                      <span className="flex items-center gap-1 text-xs text-muted-foreground">
                        <MessageSquare className="h-3 w-3" />
                        {resource.cmtnum}
                      </span>
                    </div>
                  </div>

                  {/* External link button */}
                  {resource.url && (
                    <a
                      href={resource.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
                      aria-label={"访问 " + resource.title}
                    >
                      <ExternalLink className="h-4 w-4" />
                    </a>
                  )}
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
