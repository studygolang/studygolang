import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { marked } from 'marked'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function markdownToPlainText(markdown: string): string {
  const html = marked.parse(markdown, { async: false })
  const body = new DOMParser().parseFromString(html, "text/html").body
  return (body.textContent || "")
    .replace(/\u00a0/g, " ")
    .replace(/[ \t]+/g, " ")
    .replace(/\n\s*\n+/g, "\n")
    .trim()
}

/**
 * 项目 uri 转路径段：逐段编码，保留 / 作为 catch-all 路由的分隔符。
 * 历史 uri 由项目名派生，可能含 / # ? 空格等字符。
 */
export function encodeProjectUri(uri: string | number): string {
  return String(uri)
    .split("/")
    .map(encodeURIComponent)
    .join("/")
}

/**
 * 格式化数字：>=10000 显示 x.xw，>=1000 显示 x.xk，否则原样
 */
export function formatNum(n: number): string {
  // 防御性兜底：旧响应/缓存里计数字段可能缺失
  const num = n ?? 0
  if (num >= 10000) return (num / 10000).toFixed(1) + "w"
  if (num >= 1000) return (num / 1000).toFixed(1) + "k"
  return String(num)
}

// 服务端安全的纯文本摘要（meta description 用）：strip 标签 + 压空白，不依赖 DOMParser
export function toPlainExcerpt(html: string, max = 150): string {
  return html
    .replace(/<[^>]+>/g, " ")
    .replace(/&nbsp;/g, " ")
    .replace(/\s+/g, " ")
    .trim()
    .slice(0, max)
}

/**
 * 格式化日期字符串为中文日期格式
 */
export function formatDate(dateStr: string): string {
  if (!dateStr) return ""
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleDateString("zh-CN", { year: "numeric", month: "2-digit", day: "2-digit" })
}

/**
 * 格式化时间为相对时间（"x 分钟前" 等）
 */
export function formatTime(ctime: string): string {
  try {
    const date = new Date(ctime)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    if (minutes < 60) return `${minutes} 分钟前`
    const hours = Math.floor(minutes / 60)
    if (hours < 24) return `${hours} 小时前`
    const days = Math.floor(hours / 24)
    if (days < 30) return `${days} 天前`
    return ctime.slice(0, 10)
  } catch {
    return ctime
  }
}

// 提取内容中被 @ 提及的用户名，返回逗号分隔、去重的字符串，
// 随表单发送 usernames 字段，后端 SendSysMsgAtUsernames 据此发系统消息。
// 仅匹配行首或空白后的 @，避免把邮箱 foo@bar.com 误判为提及。
export function extractMentionedUsernames(content: string): string {
  if (!content) return ''
  const re = /(?:^|\s)@([a-zA-Z0-9_]{1,20})/g
  const names = new Set<string>()
  let m: RegExpExecArray | null
  while ((m = re.exec(content)) !== null) {
    names.add(m[1])
  }
  return Array.from(names).join(',')
}
