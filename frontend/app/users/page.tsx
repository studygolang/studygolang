import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card } from "@/components/ui/card"
import { Users } from "lucide-react"
import Link from "next/link"
import type { UserListData } from "@/lib/types"

import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "会员列表 - Go语言中文网",
  description: "Go语言中文网活跃会员和新加入会员列表",
  openGraph: {
    title: "会员列表 - Go语言中文网",
    description: "Go语言中文网活跃会员和新加入会员列表",
    type: "website",
  },
}

import { fetchAPI } from "@/lib/api"

async function getUsers() {
  try {
    return await fetchAPI<UserListData>(
      '/users',
      { cache: 'no-store' }
    )
  } catch {
    return { active_users: [], new_users: [], total: 0 }
  }
}

export default async function UsersPage() {
  const data = await getUsers()
  const activeUsers = data.active_users ?? []
  const newUsers = data.new_users ?? []
  const total = data.total ?? 0

  return (
    <PageLayout>
      <PageHeader
        title="会员列表"
        description={`Go语言中文网共有 ${total.toLocaleString()} 位会员`}
        breadcrumbs={[{ label: "会员列表" }]}
      />

      <Tabs defaultValue="active" className="mt-4">
        <TabsList>
          <TabsTrigger value="active">
            <Users className="mr-1.5 h-4 w-4" />
            活跃会员
          </TabsTrigger>
          <TabsTrigger value="new">新加入会员</TabsTrigger>
        </TabsList>

        <TabsContent value="active" className="mt-4">
          <UserGrid users={activeUsers} />
        </TabsContent>

        <TabsContent value="new" className="mt-4">
          <UserGrid users={newUsers} />
        </TabsContent>
      </Tabs>
    </PageLayout>
  )
}

interface UserGridProps {
  users: UserListData['active_users']
}

function UserGrid({ users }: UserGridProps) {
  if (users.length === 0) {
    return (
      <div className="py-12 text-center text-muted-foreground">暂无数据</div>
    )
  }

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
      {users.map((user) => (
        <Link key={user.uid} href={`/user/${user.username}`}>
          <Card className="flex flex-col items-center gap-2 p-4 transition-colors hover:bg-muted/50">
            <Avatar className="h-14 w-14">
              <AvatarImage src={user.avatar} alt={user.name || user.username} />
              <AvatarFallback className="text-sm">
                {(user.name || user.username).slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="w-full text-center">
              <p className="truncate text-sm font-medium">{user.name || user.username}</p>
              {user.balance > 0 && (
                <p className="text-xs text-muted-foreground">{user.balance} 积分</p>
              )}
            </div>
          </Card>
        </Link>
      ))}
    </div>
  )
}