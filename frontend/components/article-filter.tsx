"use client"

import { useState } from "react"
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
  const [activeSort, setActiveSort] = useState("latest")
  const [activeCategory, setActiveCategory] = useState("全部")

  return (
    <div className="rounded-lg border border-border bg-card p-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        {/* Sort tabs */}
        <div className="flex items-center gap-1">
          {sortTabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveSort(tab.id)}
              className={cn(
                "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
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

      {/* Category tags */}
      <div className="mt-3 flex flex-wrap items-center gap-2">
        {categories.map((cat) => (
          <button
            key={cat}
            onClick={() => setActiveCategory(cat)}
            className={cn(
              "rounded-full px-3 py-1 text-xs font-medium transition-all",
              activeCategory === cat
                ? "bg-primary/10 text-primary"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
          >
            {cat}
          </button>
        ))}
      </div>
    </div>
  )
}
