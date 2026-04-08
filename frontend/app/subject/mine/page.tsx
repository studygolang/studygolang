"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { fetchAPI } from "@/lib/api"
import type { Subject } from "@/lib/types"
import { Loader2, BookOpen, PlusCircle, Users, FileText } from "lucide-react"
import Link from "next/link"

interface MineSubjectData {
  subjects: Subject[]
}

export default function MySubjectPage() {
  const router = useRouter()
  const [subjects, setSubjects] = useState<Subject[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchAPI<MineSubjectData>("/subject/mine", { credentials: "include" })
      .then((data) => setSubjects(data?.subjects ?? []))
      .catch((err: Error) => {
        if (err.message.includes("未登录") || err.message.includes("401")) {
          router.push("/account/login?redirect=/subject/mine")
        } else {
          setError(err.message)
        }
      })
      .finally(() => setLoading(false))
  }, [router])

  const renderContent = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      )
    }

    if (error) {
      return (
        <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive">
          {error}
        </div>
      )
    }

    if (subjects.length === 0) {
      return (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <BookOpen className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm font-medium text-foreground">还没有专栏</p>
          <p className="mt-1 text-sm text-muted-foreground">创建专栏，分享你的 Go 学习心得</p>
          <Button asChild className="mt-4 gap-2">
            <Link href="/subject/new">
              <PlusCircle className="h-4 w-4" />
              创建专栏
            </Link>
          </Button>
        </div>
      )
    }

    return (
      <div className="space-y-4">
        <div className="flex justify-end">
          <Button asChild size="sm" className="gap-2">
            <Link href="/subject/new">
              <PlusCircle className="h-4 w-4" />
              创建专栏
            </Link>
          </Button>
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          {subjects.map((subject) => (
            <Card key={subject.id} className="hover:shadow-md transition-shadow">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-2">
                  <CardTitle className="text-base leading-snug">
                    <Link
                      href={`/subject/${subject.id}`}
                      className="hover:text-primary transition-colors"
                    >
                      {subject.name}
                    </Link>
                  </CardTitle>
                  <Badge variant="secondary" className="shrink-0 text-xs">
                    专栏
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                {subject.intro && (
                  <p className="text-sm text-muted-foreground line-clamp-2">{subject.intro}</p>
                )}
                <div className="flex items-center gap-4 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <FileText className="h-3.5 w-3.5" />
                    {subject.article_num} 篇文章
                  </span>
                  <span className="flex items-center gap-1">
                    <Users className="h-3.5 w-3.5" />
                    {subject.follower_num} 关注
                  </span>
                </div>
                <div className="flex gap-2 pt-1">
                  <Button asChild size="sm" variant="outline" className="h-7 text-xs">
                    <Link href={`/subject/${subject.id}`}>查看</Link>
                  </Button>
                  <Button asChild size="sm" variant="ghost" className="h-7 text-xs">
                    <Link href={`/subject/modify?sid=${subject.id}`}>编辑</Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="我的专栏"
        description="管理你创建的专栏"
        breadcrumbs={[{ label: "专栏", href: "/subject" }, { label: "我的专栏" }]}
      />
      {renderContent()}
    </PageLayout>
  )
}
