import { Suspense } from "react"
import Link from "next/link"
import Image from "next/image"
import { Coins, Trophy } from "lucide-react"
import type { Metadata } from "next"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import type { User } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

export const metadata: Metadata = {
  title: "财富排行榜 - Go语言中文网",
  description: "Go语言中文网用户积分财富排行榜",
}

interface TopRichData {
  users: User[]
}

async function RichRankList() {
  const data = await fetchAPINullable<TopRichData>("/top/rich", { cache: "no-store" })

  if (!data || !data.users || data.users.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
        <Coins className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm text-muted-foreground">暂无数据</p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {data.users.map((user, index) => (
        <Card key={user.uid} className="transition-all hover:shadow-sm">
          <CardContent className="flex items-center gap-4 p-4">
            {/* 排名 */}
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-bold">
              {index === 0 && <span className="text-yellow-500">🥇</span>}
              {index === 1 && <span className="text-gray-400">🥈</span>}
              {index === 2 && <span className="text-amber-600">🥉</span>}
              {index > 2 && (
                <span className="text-muted-foreground">{index + 1}</span>
              )}
            </div>

            {/* 头像 */}
            <Link href={`/user/${user.username}`}>
              <div className="relative h-10 w-10 shrink-0 overflow-hidden rounded-full bg-muted">
                {user.avatar ? (
                  <Image
                    src={user.avatar}
                    alt={user.username}
                    fill
                    className="object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center text-sm font-semibold text-muted-foreground">
                    {user.username[0]?.toUpperCase()}
                  </div>
                )}
              </div>
            </Link>

            {/* 用户信息 */}
            <div className="min-w-0 flex-1">
              <Link
                href={`/user/${user.username}`}
                className="font-medium text-foreground hover:text-primary"
              >
                {user.name || user.username}
              </Link>
              {user.introduce && (
                <p className="truncate text-xs text-muted-foreground">
                  {user.introduce}
                </p>
              )}
            </div>

            {/* 积分 */}
            <div className="flex items-center gap-1 text-sm font-semibold text-amber-600">
              <Coins className="h-4 w-4" />
              <span>{user.balance ?? 0}</span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

export default function TopRichPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="财富排行榜"
        breadcrumbs={[{ label: "排行榜" }, { label: "财富榜" }]}
      />
      <div className="mb-4 flex gap-3">
        <Link
          href="/top/dau"
          className="rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
        >
          日活跃榜
        </Link>
        <Link
          href="/top/rich"
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"
        >
          <Trophy className="mr-1.5 inline h-4 w-4" />
          财富榜
        </Link>
      </div>
      <Suspense
        fallback={
          <div className="space-y-3">
            {Array.from({ length: 10 }).map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-lg bg-muted" />
            ))}
          </div>
        }
      >
        <RichRankList />
      </Suspense>
    </PageLayout>
  )
}
