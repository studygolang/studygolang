import type { Metadata } from "next"
import Link from "next/link"
import { notFound } from "next/navigation"
import { MapPin, Building2, Banknote, Clock, Tag, ArrowLeft } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { Job } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

async function getJob(id: string): Promise<Job | null> {
  const data = await fetchAPINullable<{ job: Job }>(`/jobs/${id}`, {
    next: { revalidate: 300 },
  })
  return data?.job ?? null
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>
}): Promise<Metadata> {
  const { id } = await params
  const job = await getJob(id)

  if (!job) {
    return { title: "职位不存在 - Go语言中文网" }
  }

  const description = `${job.company} 招聘 ${job.title}，薪资 ${job.salary_min}K-${job.salary_max}K，城市：${job.city}`
  return {
    title: `${job.title} - ${job.company} | 酷工作 - Go语言中文网`,
    description,
    openGraph: {
      title: `${job.title} - ${job.company}`,
      description,
      type: "website",
    },
  }
}

export default async function JobDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  const job = await getJob(id)

  if (!job) {
    return notFound()
  }

  const tags = job.tags ? job.tags.split(",").filter(Boolean) : []
  const salaryText =
    job.salary_min > 0 || job.salary_max > 0
      ? `${job.salary_min}K - ${job.salary_max}K`
      : "面议"

  return (
    <PageLayout sidebar={false}>
      <div className="mx-auto max-w-3xl">
        {/* 返回按钮 */}
        <div className="mb-4">
          <Button variant="ghost" size="sm" asChild>
            <Link href="/jobs">
              <ArrowLeft className="mr-1 h-4 w-4" />
              返回职位列表
            </Link>
          </Button>
        </div>

        <Card>
          <CardHeader className="pb-4">
            <div className="flex flex-wrap items-start justify-between gap-2">
              <div>
                <CardTitle className="text-2xl">{job.title}</CardTitle>
                <p className="mt-1 text-lg font-medium text-muted-foreground">
                  {job.company}
                </p>
              </div>
              <span className="text-xl font-bold text-primary">{salaryText}</span>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* 基本信息 */}
            <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
              {job.city && (
                <span className="flex items-center gap-1">
                  <MapPin className="h-4 w-4" />
                  {job.city}
                </span>
              )}
              {job.company_size && (
                <span className="flex items-center gap-1">
                  <Building2 className="h-4 w-4" />
                  {job.company_size}
                </span>
              )}
              {job.experience && (
                <span className="flex items-center gap-1">
                  <Clock className="h-4 w-4" />
                  {job.experience}
                </span>
              )}
              <span className="flex items-center gap-1">
                <Banknote className="h-4 w-4" />
                {salaryText}
              </span>
            </div>

            {/* 标签 */}
            {tags.length > 0 && (
              <div className="flex flex-wrap items-center gap-2">
                <Tag className="h-4 w-4 text-muted-foreground" />
                {tags.map((tag) => (
                  <Badge key={tag} variant="secondary">
                    {tag.trim()}
                  </Badge>
                ))}
              </div>
            )}

            {/* 分隔线 */}
            <div className="border-t" />

            {/* 发布时间 */}
            <p className="text-xs text-muted-foreground">
              发布时间：{new Date(job.ctime).toLocaleDateString("zh-CN")}
            </p>
          </CardContent>
        </Card>
      </div>
    </PageLayout>
  )
}
