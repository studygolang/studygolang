"use client"

import { useEffect, useState, useRef } from "react"
import { useRouter } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Calendar, Gift, CheckCircle2, Loader2 } from "lucide-react"
import { fetchAPI } from "@/lib/api"
import { Button } from "@/components/ui/button"
import Link from "next/link"

interface LoginMission {
  uid: number
  date: number
  had_redeem: boolean
}

interface MissionData {
  login_mission: LoginMission | null
  had_redeem: boolean
}

export default function MissionPage() {
  const router = useRouter()
  const [data, setData] = useState<MissionData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [redeeming, setRedeeming] = useState(false)
  const redeemingRef = useRef(false)

  useEffect(() => {
    setLoading(true)
    setError(null)
    fetchAPI<MissionData>("/mission/daily", { credentials: "include" })
      .then(setData)
      .catch((err: Error) => {
        if (err.message.includes("未登录") || err.message.includes("401")) {
          // 防止重定向循环：已在登录页则不跳转
          if (typeof window !== "undefined" && !window.location.pathname.startsWith("/account/login")) {
            router.push("/account/login?redirect=/mission/daily")
          }
        } else {
          setError(err.message)
        }
      })
      .finally(() => setLoading(false))
  }, [router])

  const handleRedeem = async () => {
    if (redeemingRef.current) return
    redeemingRef.current = true
    setRedeeming(true)
    try {
      await fetchAPI<null>("/mission/daily/redeem", {
        method: "POST",
        credentials: "include",
      })
      setData((prev) =>
        prev
          ? {
              ...prev,
              had_redeem: true,
              login_mission: prev.login_mission
                ? { ...prev.login_mission, had_redeem: true }
                : prev.login_mission,
            }
          : null
      )
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "领取失败"
      setError(message)
    } finally {
      setRedeeming(false)
      redeemingRef.current = false
    }
  }

  const content = (() => {
    if (loading) {
      return (
        <div className="flex items-center justify-center py-12">
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

    if (!data) return null

    return (
      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Calendar className="h-5 w-5 text-primary" />
              每日登录奖励
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-100">
                  <Gift className="h-5 w-5 text-amber-600" />
                </div>
                <div>
                  <p className="font-medium">每日登录</p>
                  <p className="text-sm text-muted-foreground">
                    奖励 <span className="font-semibold text-amber-600">+10</span> 积分
                  </p>
                </div>
              </div>
              {data.had_redeem ? (
                <div className="flex items-center gap-2 text-sm text-green-600">
                  <CheckCircle2 className="h-4 w-4" />
                  已领取
                </div>
              ) : (
                <Button size="sm" onClick={handleRedeem} disabled={redeeming}>
                  {redeeming ? "领取中..." : "领取奖励"}
                </Button>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">任务说明</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            <ul className="space-y-1">
              <li>每日登录可获得 10 积分奖励</li>
              <li>发布话题、文章、资源可获得额外积分</li>
              <li>评论被采纳可获得积分奖励</li>
            </ul>
          </CardContent>
        </Card>
      </div>
    )
  })()

  return (
    <PageLayout sidebar={false}>
      <PageHeader title="每日任务" breadcrumbs={[{ label: "任务中心" }]} />
      {content}
      <div className="mt-4 text-center text-sm text-muted-foreground">
        查看{" "}
        <Link href="/balance" className="text-primary hover:underline">
          积分余额
        </Link>
        {" · "}
        <Link href="/gift" className="text-primary hover:underline">
          礼品兑换
        </Link>
      </div>
    </PageLayout>
  )
}
