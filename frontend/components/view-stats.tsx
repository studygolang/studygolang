import type { ViewSource } from "@/lib/types"

interface ViewStatsProps {
  viewUserNum?: number
  viewSource?: ViewSource | null
}

export function ViewStats({ viewUserNum, viewSource }: ViewStatsProps) {
  if (!viewUserNum) return null

  const sources = viewSource?.Id
    ? [
        { label: "Google", value: viewSource.Google },
        { label: "Baidu", value: viewSource.Baidu },
        { label: "Bing", value: viewSource.Bing },
        { label: "Sogou", value: viewSource.Sogou },
        { label: "360搜索", value: viewSource.So },
        { label: "外链", value: viewSource.Other },
      ].filter((source) => source.value > 0)
    : []

  return (
    <section className="mt-4 rounded-lg border border-border bg-card p-4 text-sm">
      <h2 className="font-semibold text-foreground">统计信息</h2>
      <p className="mt-2 text-muted-foreground">
        目前已有 <strong className="text-foreground">{viewUserNum}</strong> 位注册会员查看
      </p>
      {sources.length > 0 && (
        <ul className="mt-2 space-y-1 text-muted-foreground">
          {sources.map((source) => (
            <li key={source.label}>
              {source.value} 次来自 {source.label} 点击
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
