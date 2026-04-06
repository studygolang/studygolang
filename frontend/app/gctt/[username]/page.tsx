import type { Metadata } from "next"
import { notFound } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPI } from "@/lib/api"

interface GCTTUser {
  id: number
  username: string
  avatar: string
  num: number
  words: number
  role_name: string
  avg_time: number
  joined_at: string
}

async function getUser(username: string): Promise<GCTTUser | null> {
  try {
    const data = await fetchAPI<{ user: GCTTUser }>(`/gctt/${username}`, {
      cache: "no-store",
    })
    return data.user ?? null
  } catch {
    return null
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ username: string }>
}): Promise<Metadata> {
  const { username } = await params
  const user = await getUser(username)
  return {
    title: user
      ? `${user.username} - GCTT 译者 - Go语言中文网`
      : "译者不存在 - Go语言中文网",
    description: user
      ? `${user.username}，GCTT Go 中文翻译组${user.role_name}，已翻译 ${user.num} 篇文章`
      : undefined,
  }
}

export default async function GCTTUserPage({
  params,
}: {
  params: Promise<{ username: string }>
}) {
  const { username } = await params
  const user = await getUser(username)

  if (!user || user.id === 0) {
    notFound()
  }

  return (
    <PageLayout>
      <PageHeader
        title={user.username}
        breadcrumbs={[
          { label: "GCTT", href: "/gctt" },
          { label: user.username },
        ]}
      />

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center gap-4">
            <img
              src={user.avatar}
              alt={user.username}
              className="h-16 w-16 rounded-full"
            />
            <div>
              <h1 className="text-xl font-bold">{user.username}</h1>
              <p className="text-sm text-muted-foreground">{user.role_name}</p>
            </div>
          </div>

          <div className="mt-6 grid grid-cols-3 gap-4 text-center">
            <div className="rounded-lg border border-border p-3">
              <div className="text-2xl font-bold">{user.num}</div>
              <div className="text-xs text-muted-foreground">翻译篇数</div>
            </div>
            <div className="rounded-lg border border-border p-3">
              <div className="text-2xl font-bold">{user.words}</div>
              <div className="text-xs text-muted-foreground">翻译词数</div>
            </div>
            <div className="rounded-lg border border-border p-3">
              <div className="text-2xl font-bold">
                {user.avg_time > 0 ? `${user.avg_time}h` : "-"}
              </div>
              <div className="text-xs text-muted-foreground">平均耗时</div>
            </div>
          </div>

          {user.joined_at && (
            <div className="mt-4 text-sm text-muted-foreground">
              加入时间：{user.joined_at?.slice(0, 10)}
            </div>
          )}
        </CardContent>
      </Card>
    </PageLayout>
  )
}
