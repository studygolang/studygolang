"use client"

import { useEffect, useState } from "react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { gcttAPI } from "@/lib/api"
import { useAuth } from "@/lib/auth-context"
import { Loader2, CheckCircle2, BookOpen } from "lucide-react"
import Link from "next/link"

type Status = "loading" | "already_translator" | "apply_form" | "success" | "bind_github" | "error"

export default function GcttApplyPage() {
  const { user } = useAuth()
  const [status, setStatus] = useState<Status>("loading")
  const [errorMsg, setErrorMsg] = useState("")
  const [applying, setApplying] = useState(false)

  useEffect(() => {
    // 检查登录状态及是否已是译者
    gcttAPI.getMe()
      .then((data) => {
        if (data?.is_translator) {
          setStatus("already_translator")
        } else {
          setStatus("apply_form")
        }
      })
      .catch((err: Error) => {
        if (err.message.includes("未登录") || err.message.includes("401")) {
          window.location.href = "/account/login?redirect=/gctt/apply"
        } else {
          setStatus("error")
          setErrorMsg(err.message)
        }
      })
  }, [])

  const handleApply = async () => {
    if (applying) return
    setApplying(true)
    try {
      await gcttAPI.apply()
      setStatus("success")
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "申请失败"
      if (msg.includes("GitHub") || msg.includes("github")) {
        setStatus("bind_github")
      } else {
        setStatus("error")
        setErrorMsg(msg)
      }
    } finally {
      setApplying(false)
    }
  }

  const renderContent = () => {
    if (status === "loading") {
      return (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      )
    }

    if (status === "already_translator") {
      return (
        <Card>
          <CardContent className="flex flex-col items-center py-12 text-center">
            <CheckCircle2 className="mb-4 h-12 w-12 text-green-500" />
            <h2 className="mb-2 text-lg font-semibold">您已是 GCTT 译者</h2>
            <p className="text-sm text-muted-foreground">感谢您对 Go 语言中文网的贡献！</p>
            <Button asChild className="mt-6">
              <Link href="/gctt">前往 GCTT 主页</Link>
            </Button>
          </CardContent>
        </Card>
      )
    }

    if (status === "success") {
      return (
        <Card>
          <CardContent className="flex flex-col items-center py-12 text-center">
            <CheckCircle2 className="mb-4 h-12 w-12 text-green-500" />
            <h2 className="mb-2 text-lg font-semibold">申请成功，欢迎加入 GCTT！</h2>
            <p className="text-sm text-muted-foreground">
              您已成为 GCTT 译者，可以开始翻译 Go 相关文章了。
            </p>
            <Button asChild className="mt-6">
              <Link href="/gctt">前往 GCTT 主页</Link>
            </Button>
          </CardContent>
        </Card>
      )
    }

    if (status === "bind_github") {
      return (
        <Card>
          <CardContent className="flex flex-col items-center py-12 text-center">
            <svg className="mb-4 h-12 w-12 text-muted-foreground" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
            <h2 className="mb-2 text-lg font-semibold">请先绑定 GitHub 账号</h2>
            <p className="mb-6 text-sm text-muted-foreground">
              申请成为 GCTT 译者需要绑定 GitHub 账号，请前往个人设置完成绑定。
            </p>
            <div className="flex gap-3">
              <Button asChild variant="outline">
                <Link href={`/user/${user?.username || ''}/settings`}>前往个人设置</Link>
              </Button>
              <Button onClick={() => setStatus("apply_form")} variant="ghost">
                返回
              </Button>
            </div>
          </CardContent>
        </Card>
      )
    }

    if (status === "error") {
      return (
        <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive">
          {errorMsg || "加载失败，请刷新重试"}
        </div>
      )
    }

    // apply_form
    return (
      <div className="space-y-6">
        {/* GCTT 介绍 */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <BookOpen className="h-5 w-5 text-primary" />
              什么是 GCTT？
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-muted-foreground">
            <p>
              GCTT（Go Chinese Translation Team）是 Go 语言中文网发起的翻译组，
              专注于翻译 Go 语言相关的优质英文文章，让更多中文开发者受益。
            </p>
            <ul className="list-inside list-disc space-y-1">
              <li>翻译 Go 官方博客、技术文章等优质内容</li>
              <li>通过 GitHub 协作完成翻译流程</li>
              <li>获得社区认可和积分奖励</li>
              <li>与国内外 Go 开发者交流学习</li>
            </ul>
          </CardContent>
        </Card>

        {/* 申请要求 */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">申请要求</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-muted-foreground">
            <ul className="list-inside list-disc space-y-1">
              <li>具备基本的英文阅读能力</li>
              <li>对 Go 语言有一定了解</li>
              <li>需要绑定 GitHub 账号（用于协作翻译）</li>
              <li>有时间和热情参与社区贡献</li>
            </ul>
          </CardContent>
        </Card>

        {/* 申请按钮 */}
        <div className="flex justify-center">
          <Button size="lg" onClick={handleApply} disabled={applying} className="gap-2 px-8">
            {applying && <Loader2 className="h-4 w-4 animate-spin" />}
            {applying ? "申请中..." : "申请成为 GCTT 译者"}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="申请成为 GCTT 译者"
        description="加入 Go 语言中文翻译团队，贡献优质内容"
        breadcrumbs={[{ label: "GCTT", href: "/gctt" }, { label: "申请译者" }]}
      />
      {renderContent()}
    </PageLayout>
  )
}
