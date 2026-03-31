import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * 格式化数字：>=10000 显示 x.xw，>=1000 显示 x.xk，否则原样
 */
export function formatNum(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1) + "w"
  if (n >= 1000) return (n / 1000).toFixed(1) + "k"
  return String(n)
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
