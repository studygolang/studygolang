"use client"

// ArticleFilter：文章列表过滤组件
// 排序通过 URL searchParams 传给 SSR 页面（可被搜索引擎感知）
// 分类标签：后端文章列表 API 暂不支持按分类过滤，纯前端 UI 展示，不发起 API 请求

import { useRouter, useSearchParams } from "next/navigation"
import { cn } from "@/lib/utils"

const sortTabs = [
  { id: "latest", label: "最新" },
  { id: "hot", label: "热门" },
  { id: "recommend", label: "推荐" },
]

const categories = [
  "全部",
  "Go基础",
  "Go标准库",
  "Go Web开发",
  "微服务",
  "数据库",
  "Kubernetes",
  "人工智能",
  "Go实战",
  "源码分析",
]

export function ArticleFilter() {
  const router = useRouter()
  const searchParams = useSearchParams()

  const activeSort = searchParams.get("sort") || "latest"
  // 分类过滤暂无后端支持，保留 UI 状态但不触发请求

  function handleSortClick(sortId: string) {
    const params = new URLSearchParams(searchParams.toString())
    params.set("sort", sortId)
    // 切换排序时重置到第一页
    params.delete("page")
    router.push(`/articles?${params.toString()}`)
  }

  return (
    <div className="rounded-lg border border-border bg-card p-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        {/* Sort tabs：通过 URL sort 参数与 SSR 页面联动 */}
        <div className="flex items-center gap-1">
          {sortTabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => handleSortClick(tab.id)}
              className={cn(
                "cursor-pointer rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
                activeSort === tab.id
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground hover:bg-secondary hover:text-foreground"
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Category tags：纯 UI 展示，后端支持 tag 过滤后可改为 URL 参数 */}
      <div className="mt-3 flex flex-wrap items-center gap-2">
        {categories.map((cat) => (
          <span
            key={cat}
            title="分类过滤功能待后端支持"
            className={cn(
              "cursor-default rounded-full px-3 py-1 text-xs font-medium",
              cat === "全部"
                ? "bg-primary/10 text-primary"
                : "bg-secondary text-secondary-foreground"
            )}
          >
            {cat}
          </span>
        ))}
      </div>
    </div>
  )
}
