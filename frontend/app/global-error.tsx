'use client'

import { useEffect } from 'react'

export default function GlobalError({
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
    <html lang="zh-CN">
      <body style={{ margin: 0, padding: 0 }}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', fontFamily: 'system-ui, sans-serif' }}>
          <h1 style={{ fontSize: '3.75rem', fontWeight: 'bold', color: '#9ca3af', marginBottom: '1rem' }}>系统错误</h1>
          <p style={{ color: '#6b7280', marginBottom: '2rem' }}>系统遇到问题，请刷新页面重试</p>
          {error.digest && (
            <p style={{ fontSize: '0.875rem', color: '#9ca3af', marginBottom: '1rem' }}>错误 ID: {error.digest}</p>
          )}
          <button
            onClick={reset}
            style={{ padding: '0.5rem 1rem', backgroundColor: '#2563eb', color: 'white', border: 'none', borderRadius: '0.375rem', cursor: 'pointer' }}
          >
            重新加载
          </button>
        </div>
      </body>
    </html>
  )
}
