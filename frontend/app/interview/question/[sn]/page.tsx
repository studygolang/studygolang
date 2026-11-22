import Link from "next/link"
import { BookOpen, ChevronLeft, Eye, MessageSquare } from "lucide-react"
import { Metadata } from "next"
import { cookies } from "next/headers"
import { notFound } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { OBJTYPE_INTERVIEW, type Comment, type InterviewQuestion } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"
import { SafeHtml } from "@/components/safe-html"
import { LikeButton } from "@/components/like-button"
import { CommentSection } from "@/components/comment-section"

export const dynamic = "force-dynamic"

interface InterviewQuestionData {
  question: InterviewQuestion
  replies?: Comment[]
  likeflag?: number
}

const levelLabels = ["初级", "中级", "高级"]
const levelColors = [
  "bg-green-100 text-green-700",
  "bg-yellow-100 text-yellow-700",
  "bg-red-100 text-red-700",
]

interface QuestionPageProps {
  params: Promise<{ sn: string }>
}

export async function generateMetadata({ params }: QuestionPageProps): Promise<Metadata> {
  const { sn } = await params
  return {
    title: `Go 面试题 #${sn} - Go语言中文网`,
    description: `Go 语言面试题详情及答案解析`,
  }
}

export default async function InterviewQuestionPage({ params }: QuestionPageProps) {
  const { sn } = await params
  const requestCookies = await cookies()
  const data = await fetchAPINullable<InterviewQuestionData>(
    `/interviews/question/${sn}`,
    {
      cache: "no-store",
      headers: { Cookie: requestCookies.toString() },
    }
  )

  if (!data?.question) {
    return notFound()
  }

  const { question } = data

  return (
    <PageLayout>
      {/* Back link */}
      <div className="mb-4">
        <Link
          href="/interview"
          className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-primary"
        >
          <ChevronLeft className="h-4 w-4" />
          返回面试题列表
        </Link>
      </div>

      {/* Question card */}
      <Card className="mb-4">
        <CardHeader className="p-5 pb-3">
          <div className="flex flex-wrap items-center gap-2">
            <BookOpen className="h-4 w-4 text-primary" />
            <CardTitle className="text-base font-semibold">题目</CardTitle>
            {question.level >= 0 && question.level <= 2 && (
              <Badge className={`text-[10px] font-semibold ${levelColors[question.level]}`}>
                {levelLabels[question.level]}
              </Badge>
            )}
          </div>
        </CardHeader>
        <CardContent className="px-5 pb-5">
          <SafeHtml
            html={question.question}
            className="prose prose-sm max-w-none text-sm leading-relaxed text-foreground"
          />
        </CardContent>
      </Card>

      {/* Answer card */}
      <Card className="mb-4">
        <CardHeader className="p-5 pb-3">
          <CardTitle className="flex items-center gap-2 text-base font-semibold">
            <BookOpen className="h-4 w-4 text-primary" />
            答案解析
          </CardTitle>
        </CardHeader>
        <CardContent className="px-5 pb-5">
          <SafeHtml
            html={question.answer}
            className="prose prose-sm max-w-none text-sm leading-relaxed text-foreground"
          />
        </CardContent>
      </Card>

      {/* Meta */}
      <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <Eye className="h-3 w-3" />
          {question.viewnum} 次查看
        </span>
        <LikeButton
          objid={question.id}
          objtype={OBJTYPE_INTERVIEW}
          initialLiked={Boolean(data.likeflag)}
          initialCount={question.likenum}
          className="-ml-3 h-auto px-3 py-0 text-xs"
        />
        <span className="flex items-center gap-1">
          <MessageSquare className="h-3 w-3" />
          {question.cmtnum} 评论
        </span>
        {question.source && (
          <span>来源：{question.source}</span>
        )}
      </div>

      <div className="mt-4">
        <CommentSection
          objid={question.id}
          objtype={OBJTYPE_INTERVIEW}
          comments={data.replies ?? []}
        />
      </div>
    </PageLayout>
  )
}
