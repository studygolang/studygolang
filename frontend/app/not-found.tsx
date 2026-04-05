import Link from "next/link"

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center py-24">
      <h1 className="text-6xl font-bold text-muted-foreground mb-4">404</h1>
      <h2 className="text-xl font-semibold mb-4">页面不存在</h2>
      <p className="text-muted-foreground mb-8">你访问的页面可能已被移动或删除</p>
      <Link href="/" className="px-4 py-2 bg-primary text-primary-foreground rounded-md">
        返回首页
      </Link>
    </div>
  )
}
