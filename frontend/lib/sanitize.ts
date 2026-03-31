// HTML 安全渲染工具 —— 防止 XSS 攻击
// 使用 isomorphic-dompurify 同时支持 SSR 和客户端环境

import DOMPurify from 'isomorphic-dompurify'

// 允许的标签：保留富文本格式，但过滤 script/iframe 等危险标签
const ALLOWED_TAGS = [
  'a', 'abbr', 'b', 'blockquote', 'br', 'code', 'dd', 'del', 'dl', 'dt',
  'em', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr', 'i', 'img', 'kbd',
  'li', 'mark', 'ol', 'p', 'pre', 's', 'strong', 'sub', 'sup', 'table',
  'tbody', 'td', 'th', 'thead', 'tr', 'ul',
]

// 允许的属性：仅保留安全的展示属性
const ALLOWED_ATTR = ['href', 'src', 'alt', 'title', 'class', 'target', 'rel']

/**
 * 清理 HTML 内容，移除潜在的 XSS 攻击向量
 * 用于所有 dangerouslySetInnerHTML 渲染场景
 */
export function sanitizeHtml(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS,
    ALLOWED_ATTR,
  })
}
