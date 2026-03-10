"use client"

import { useRouter, useSearchParams } from "next/navigation"
import { cn } from "@/lib/utils"

const primaryTabs = [
  { id: "all", label: "全部" },
  { id: "hot", label: "最热" },
  { id: "latest", label: "最新" },
]

// 首页 TabFilter：tab 切换通过 URL 参数 ?tab=xxx 驱动，实际数据由 SSR 父页面根据 tab 从 API 获取
export function TabFilter() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const activeTab = searchParams.get("tab") ?? "all"

  function handleTabChange(tabId: string) {
    const params = new URLSearchParams(searchParams.toString())
    params.set("tab", tabId)
    router.push(`/?${params.toString()}`)
  }

  return (
    <div className="space-y-3">
      {/* Primary tabs */}
      <div className="flex items-center gap-1 border-b border-border">
        {primaryTabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => handleTabChange(tab.id)}
            className={cn(
              "relative cursor-pointer px-4 py-2.5 text-sm font-medium transition-colors",
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
    </div>
  )
}
