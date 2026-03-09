"use client"

import { useState } from "react"
import Link from "next/link"
import { ArrowRight, Sparkles, Code2, Users, BookOpen } from "lucide-react"
import { Badge } from "@/components/ui/badge"

const announcements = [
  {
    badge: "\u65b0\u7248\u53d1\u5e03",
    title: "Go 1.26 RC1 \u5df2\u53d1\u5e03\uff0c\u65b0\u7279\u6027\u62a2\u5148\u770b",
    href: "/topics/1",
  },
  {
    badge: "\u793e\u533a\u62db\u52df",
    title: "\u5bfb\u627e\u793e\u533a\u65e5\u5e38\u8fd0\u8425\u3001\u529f\u80fd\u5f00\u53d1\u3001\u7ef4\u62a4\u5fd7\u613f\u8005",
    href: "/topics/volunteer",
  },
]

export function HeroBanner() {
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

          {/* Quick stats */}
          <div className="hidden items-center gap-5 sm:flex">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Users className="h-3.5 w-3.5" />
              <span>{"189,432 \u4f1a\u5458"}</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <BookOpen className="h-3.5 w-3.5" />
              <span>{"52,871 \u4e3b\u9898"}</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Code2 className="h-3.5 w-3.5" />
              <span>{"1,286 \u9879\u76ee"}</span>
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
                className={`h-1 rounded-full transition-all ${
                  i === activeAnnouncement
                    ? "w-4 bg-primary"
                    : "w-1 bg-border hover:bg-muted-foreground"
                }`}
                aria-label={`\u516c\u544a ${i + 1}`}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
