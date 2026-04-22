import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { fetchAPI } from "@/lib/api"
import { SafeHtml } from "@/components/safe-html"

export const metadata: Metadata = {
  title: "GCTT - Go中文翻译组 - Go语言中文网",
  description: "GCTT（Go Chinese Translation Team）是 Go 语言中文网的翻译组，致力于将优秀的 Go 语言英文内容翻译成中文",
}

interface GCTTTimeLine {
  id: number
  content: string
  created_at: string
}

interface GCTTUser {
  id: number
  username: string
  avatar: string
  num: number
  words: number
  role_name: string
}

interface GCTTIssue {
  id: number
  title: string
  translator: string
  label: string
  state: number
  created_at: string
}

async function getGCTTIndex() {
  try {
    return await fetchAPI<{
      time_lines: GCTTTimeLine[]
      core_users: GCTTUser[]
      untranslated_issues: GCTTIssue[]
    }>("/gctt", { cache: "no-store" })
  } catch {
    return null
  }
}

export default async function GCTTPage() {
  const data = await getGCTTIndex()

  return (
    <PageLayout>
      <PageHeader
        title="GCTT"
        description="Go 中文翻译组（Go Chinese Translation Team）"
        breadcrumbs={[
          { label: "GCTT" },
        ]}
        actions={
          <Link
            href="/gctt/issues"
            className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
          >
            查看选题列表
          </Link>
        }
      />

      {/* 核心成员 */}
      {data?.core_users && data.core_users.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-base">核心成员</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-4">
              {data.core_users.map((user) => (
                <Link
                  key={user.id}
                  href={`/gctt/${user.username}`}
                  className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 hover:bg-muted/50"
                >
                  <img
                    src={user.avatar}
                    alt={user.username}
                    className="h-8 w-8 rounded-full"
                  />
                  <div>
                    <div className="text-sm font-medium">{user.username}</div>
                    <div className="text-xs text-muted-foreground">
                      {user.role_name} · {user.num} 篇 · {user.words} 词
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 待认领选题 */}
      {data?.untranslated_issues && data.untranslated_issues.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-base">待认领选题</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {data.untranslated_issues.map((issue) => (
                <div key={issue.id} className="rounded-lg border border-border p-3">
                  <div className="font-medium">{issue.title}</div>
                  <div className="mt-1 text-xs text-muted-foreground">
                    {issue.label && (
                      <span className="mr-2 rounded bg-primary/10 px-1.5 py-0.5 text-primary">
                        {issue.label}
                      </span>
                    )}
                    创建于 {issue.created_at?.slice(0, 10)}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 时间线 */}
      {data?.time_lines && data.time_lines.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">动态</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {data.time_lines.map((tl) => (
                <div key={tl.id} className="flex gap-3">
                  <div className="mt-1 h-2 w-2 shrink-0 rounded-full bg-primary" />
                  <div>
                    <SafeHtml html={tl.content} className="text-sm" />
                    <div className="text-xs text-muted-foreground">
                      {tl.created_at?.slice(0, 10)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      <div className="mt-6 flex gap-3">
        <Link
          href="/gctt/users"
          className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
        >
          全部译者
        </Link>
      </div>
    </PageLayout>
  )
}
