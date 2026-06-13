"use client"

import { useCallback, useEffect, useState } from "react"
import { ArrowLeft, X } from "lucide-react"
import { fetchAPI } from "@/lib/api"
import { cn } from "@/lib/utils"

const DISMISS_KEY = "sg_new_version_banner_dismissed"

type BannerState = "loading" | "visible" | "hidden"

export function VersionBanner() {
  const [state, setState] = useState<BannerState>("loading")
  const [switching, setSwitching] = useState(false)

  useEffect(() => {
    try {
      const dismissed = window.localStorage.getItem(DISMISS_KEY) === "1"
      setState(dismissed ? "hidden" : "visible")
    } catch {
      setState("visible")
    }
  }, [])

  const switchToOld = useCallback(async () => {
    setSwitching(true)
    try {
      await fetchAPI("/version/switch", {
        method: "POST",
        body: JSON.stringify({ version: "old" }),
      })
      window.location.reload()
    } catch {
      setSwitching(false)
    }
  }, [])

  const dismiss = useCallback(() => {
    try {
      window.localStorage.setItem(DISMISS_KEY, "1")
    } catch {
      // localStorage 不可用时仅隐藏 UI
    }
    setState("hidden")
  }, [])

  if (state !== "visible") {
    return null
  }

  return (
    <div
      role="region"
      aria-label="版本切换提示"
      className={cn(
        "flex items-center justify-center gap-3 bg-gradient-to-r from-teal-600 to-cyan-600 px-4 py-2 text-sm text-white",
        "relative"
      )}
    >
      <span className="font-medium">您正在体验全新界面</span>
      <span className="hidden text-white/80 sm:inline">如遇问题可随时切回旧版</span>
      <button
        type="button"
        onClick={switchToOld}
        disabled={switching}
        className={cn(
          "inline-flex items-center gap-1 rounded-md bg-white/15 px-3 py-1 text-xs font-medium",
          "transition-colors hover:bg-white/25 disabled:opacity-60"
        )}
      >
        <ArrowLeft className="h-3 w-3" />
        {switching ? "切换中…" : "返回旧版"}
      </button>
      <button
        type="button"
        onClick={dismiss}
        aria-label="关闭提示"
        className="absolute right-2 top-1/2 -translate-y-1/2 rounded p-1 text-white/70 transition-colors hover:bg-white/15 hover:text-white"
      >
        <X className="h-3.5 w-3.5" />
      </button>
    </div>
  )
}
