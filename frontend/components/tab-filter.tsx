"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"

const primaryTabs = [
  { id: "recommend", label: "\u63a8\u8350" },
  { id: "hot", label: "\u6700\u70ed" },
  { id: "latest", label: "\u6700\u65b0" },
]

const tagFilters = [
  "\u5168\u90e8",
  "Go\u57fa\u7840",
  "\u5fae\u670d\u52a1",
  "Web\u5f00\u53d1",
  "\u95ee\u4e0e\u7b54",
  "\u9177\u5de5\u4f5c",
  "\u4eba\u5de5\u667a\u80fd",
  "Kubernetes",
  "\u5f00\u6e90\u9879\u76ee",
]

export function TabFilter() {
  const [activeTab, setActiveTab] = useState("recommend")
  const [activeTag, setActiveTag] = useState("\u5168\u90e8")

  return (
    <div className="space-y-3">
      {/* Primary tabs */}
      <div className="flex items-center gap-1 border-b border-border">
        {primaryTabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              "relative px-4 py-2.5 text-sm font-medium transition-colors",
              activeTab === tab.id
                ? "text-primary"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            {tab.label}
            {activeTab === tab.id && (
              <span className="absolute bottom-0 left-1/2 h-0.5 w-8 -translate-x-1/2 rounded-full bg-primary" />
            )}
          </button>
        ))}
      </div>

      {/* Tag filters */}
      <div className="flex flex-wrap items-center gap-2">
        {tagFilters.map((tag) => (
          <button
            key={tag}
            onClick={() => setActiveTag(tag)}
            className={cn(
              "rounded-full px-3 py-1 text-xs font-medium transition-all",
              activeTag === tag
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
          >
            {tag}
          </button>
        ))}
      </div>
    </div>
  )
}
