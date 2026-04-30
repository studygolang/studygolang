'use client'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <h2 className="text-xl font-semibold mb-4">出错了！</h2>
      <p className="text-muted-foreground mb-6">页面加载失败，请稍后重试</p>
      {error.digest && (
        <p className="text-sm text-muted-foreground/60 mb-4">错误 ID: {error.digest}</p>
      )}
      <button onClick={reset} className="px-4 py-2 bg-primary text-primary-foreground rounded-md">
        重试
      </button>
    </div>
  )
}
