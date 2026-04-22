"use client"

import { sanitizeHtml } from "@/lib/sanitize"

interface SafeHtmlProps {
  html: string
  className?: string
}

/**
 * SafeHtml: 安全渲染 HTML 内容，自动调用 sanitizeHtml() 防 XSS
 * 替代所有 dangerouslySetInnerHTML 的直接使用
 */
export function SafeHtml({ html, className }: SafeHtmlProps) {
  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: sanitizeHtml(html) }}
    />
  )
}
