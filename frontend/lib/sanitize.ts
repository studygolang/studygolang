// HTML 安全渲染工具 —— 防止 XSS 攻击
// 使用 isomorphic-dompurify 同时支持 SSR 和客户端环境

import DOMPurify from 'isomorphic-dompurify'

// 允许的标签：保留富文本格式，但过滤 script 等危险标签。
// iframe 仅放行 Wide Playground（见下方 hook），其余一律移除。
const ALLOWED_TAGS = [
  'a', 'abbr', 'b', 'blockquote', 'br', 'code', 'dd', 'del', 'dl', 'dt',
  'em', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr', 'i', 'iframe', 'img', 'kbd',
  'li', 'mark', 'ol', 'p', 'pre', 's', 'strong', 'sub', 'sup', 'table',
  'tbody', 'td', 'th', 'thead', 'tr', 'ul',
]

// 允许的属性：仅保留安全的展示属性
const ALLOWED_ATTR = ['href', 'src', 'alt', 'title', 'class', 'target', 'rel', 'width', 'height']

// 旧版 EmbedWide 会把主题里的 Wide 代码片段转成 iframe，只允许这一个来源
const WIDE_IFRAME_SRC = /^https:\/\/wide\.b3log\.org\/playground([/?#]|$)/

DOMPurify.addHook('uponSanitizeElement', (node, data) => {
  if (data.tagName !== 'iframe') return
  const src = (node as Element).getAttribute?.('src') ?? ''
  if (!WIDE_IFRAME_SRC.test(src)) {
    node.parentNode?.removeChild(node)
  }
})

/**
 * 清理 HTML 内容，移除潜在的 XSS 攻击向量
 * 用于所有 dangerouslySetInnerHTML 渲染场景
 */
export function sanitizeHtml(html: string): string {
  const clean = DOMPurify.sanitize(html, {
    ALLOWED_TAGS,
    ALLOWED_ATTR,
    ADD_ATTR: ['target'],
  })

  // 为 target="_blank" 链接强制添加 rel="noopener noreferrer" 防止 tabnabbing
  return clean.replace(
    /<a\s+([^>]*target=["']_blank["'][^>]*)>/gi,
    (_, attrs) => {
      if (/rel=["'][^"']*noopener/.test(attrs)) return `<a ${attrs}>`
      return `<a ${attrs} rel="noopener noreferrer">`
    }
  )
}
