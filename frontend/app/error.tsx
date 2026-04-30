'use client'

import { useEffect } from 'react'
import Link from 'next/link'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      console.error(error)
    }
  }, [error])

  return (
    <div className="flex flex-col items-center justify-center py-24">
      <h1 className="text-6xl font-bold text-muted-foreground mb-4">出错了</h1>
      <p className="text-muted-foreground mb-8">页面加载失败，请稍后重试</p>
      {error.digest && (
        <p className="text-sm text-muted-foreground/60 mb-4">错误 ID: {error.digest}</p>
      )}
      <div className="flex gap-4">
        <button
          onClick={reset}
          className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
        >
          重试
        </button>
        <Link
          href="/"
          className="px-4 py-2 border border-border rounded-md"
        >
          返回首页
        </Link>
      </div>
    </div>
  )
}
