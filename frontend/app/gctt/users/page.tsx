import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPI } from "@/lib/api"

export const metadata: Metadata = {
  title: "GCTT 译者列表 - Go语言中文网",
  description: "GCTT Go 中文翻译组全部译者列表",
}

interface GCTTUser {
  id: number
  username: string
  avatar: string
  num: number
  words: number
  role_name: string
  avg_time: number
}

async function getUsers() {
  try {
    const data = await fetchAPI<{ users: GCTTUser[] }>("/gctt/users", {
      cache: "no-store",
    })
    return data.users ?? []
  } catch {
    return []
  }
}

export default async function GCTTUsersPage() {
  const users = await getUsers()

  return (
    <PageLayout>
      <PageHeader
        title="GCTT 译者列表"
        breadcrumbs={[
          { label: "GCTT", href: "/gctt" },
          { label: "译者列表" },
        ]}
      />

      {users.length > 0 ? (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {users.map((user) => (
            <Link key={user.id} href={`/gctt/${user.username}`}>
              <Card className="h-full transition-shadow hover:shadow-md">
                <CardContent className="flex items-center gap-3 p-4">
                  <img
                    src={user.avatar}
                    alt={user.username}
                    className="h-10 w-10 rounded-full"
                  />
                  <div className="min-w-0 flex-1">
                    <div className="truncate font-medium">{user.username}</div>
                    <div className="text-xs text-muted-foreground">
                      {user.role_name} · {user.num} 篇 · {user.words} 词
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <p className="text-sm text-muted-foreground">暂无译者</p>
        </div>
      )}
    </PageLayout>
  )
}
