"use client"

import { useMemo } from "react"
import { marked } from "marked"
import { sanitizeHtml } from "@/lib/sanitize"

marked.setOptions({
  breaks: true,
  gfm: true,
})

interface MarkdownContentProps {
  content: string
  className?: string
}

export function MarkdownContent({ content, className }: MarkdownContentProps) {
  const html = useMemo(() => {
    const raw = marked.parse(content) as string
    return sanitizeHtml(raw)
  }, [content])

  return (
    <div
      className={`prose prose-sm max-w-none text-foreground ${className ?? ""}`}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  )
}
