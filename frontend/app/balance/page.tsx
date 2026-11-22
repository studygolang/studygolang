"use client"

import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Coins, TrendingUp, TrendingDown, Loader2, Sparkles } from "lucide-react"
import Link from "next/link"
import { fetchAPI, needRelogin } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { SafeHtml } from "@/components/safe-html"

interface BalanceDetail {
  id: number
  desc: string
  type_show: string
  num: number
  balance: number
  created_at: string
}

interface InitialMission {
  id: number
  name: string
  fixed: number
}

interface BalanceData {
  balance: number
  details: BalanceDetail[]
  total: number
  page: number
  has_more: boolean
  initial_mission: InitialMission | null
}

async function fetchBalance(page: number): Promise<BalanceData> {
  return fetchAPI<BalanceData>(`/balance?p=${page}`, { credentials: "include" })
}

function BalanceContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const page = Math.max(1, Number(searchParams.get("p")) || 1)
  const [data, setData] = useState<BalanceData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [reload, setReload] = useState(0)
  const [claiming, setClaiming] = useState(false)

  useEffect(() => {
    setLoading(true)
    setError(null)
    fetchBalance(page)
      .then(setData)
      .catch((err: Error) => {
        if (needRelogin(err)) {
          const target = page > 1 ? `/balance?p=${page}` : "/balance"
          router.push(`/account/login?redirect=${encodeURIComponent(target)}`)
        } else {
          setError(err.message)
        }
      })
      .finally(() => setLoading(false))
  }, [page, reload, router])

  const goToPage = (targetPage: number) => {
    const params = new URLSearchParams(searchParams.toString())
    if (targetPage > 1) {
      params.set("p", String(targetPage))
    } else {
      params.delete("p")
    }
    const query = params.toString()
    router.push(query ? `/balance?${query}` : "/balance")
  }

  const handleClaimInitial = async () => {
    const mission = data?.initial_mission
    if (!mission || claiming) return
    setClaiming(true)
    try {
      await fetchAPI<null>(`/mission/complete/${mission.id}`, {
        method: "POST",
        credentials: "include",
      })
      toast.success(`已领取初始资本 +${mission.fixed} 铜币`)
      if (page > 1) router.push("/balance")
      setReload((r) => r + 1)
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "领取失败")
    } finally {
      setClaiming(false)
    }
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="我的积分"
        breadcrumbs={[{ label: "积分中心" }]}
      />

      {/* 总积分卡片 */}
      {data && (
        <Card className="mb-6">
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium text-muted-foreground">
              当前积分余额
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <Coins className="h-8 w-8 text-amber-500" />
              <span className="text-4xl font-bold text-amber-600">{data.balance}</span>
              <span className="text-sm text-muted-foreground">积分</span>
            </div>
          </CardContent>
        </Card>
      )}

      {/* 初始资本：没有任何积分记录时才出现，领取后后端不再下发 */}
      {data?.initial_mission && (
        <Card className="mb-6 border-amber-200 bg-amber-50/50 dark:border-amber-900 dark:bg-amber-950/30">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <Sparkles className="h-5 w-5 text-amber-600" />
              获得初始资本
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <p className="text-sm text-muted-foreground">
              欢迎来到社区。领取初始资本后你才能发布主题、文章等内容；发布和回复会消耗铜币，
              完成有意义的事情（被回复、被采纳等）则会收获铜币。请优先把内容发布在
              <span className="font-medium text-foreground">主题</span>中。
            </p>
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold text-amber-600">
                任务奖励 +{data.initial_mission.fixed} 铜币
              </span>
              <Button size="sm" onClick={handleClaimInitial} disabled={claiming}>
                {claiming ? "领取中..." : "完成任务"}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* 积分明细 */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">积分明细</CardTitle>
        </CardHeader>
        <CardContent>
          {loading && (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          )}

          {error && (
            <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive">
              {error}
            </div>
          )}

          {!loading && !error && data && (
            <>
              {data.details.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <Coins className="mb-3 h-10 w-10 text-muted-foreground/40" />
                  <p className="text-sm text-muted-foreground">暂无积分记录</p>
                </div>
              ) : (
                <div className="divide-y">
                  {data.details.map((item) => (
                    <div
                      key={item.id}
                      className="flex items-center justify-between py-3"
                    >
                      <div className="flex items-center gap-3">
                        <div className={`rounded-full p-1.5 ${item.num > 0 ? "bg-green-100" : "bg-red-100"}`}>
                          {item.num > 0 ? (
                            <TrendingUp className="h-4 w-4 text-green-600" />
                          ) : (
                            <TrendingDown className="h-4 w-4 text-red-600" />
                          )}
                        </div>
                        <div>
                          <p className="text-sm font-medium">{item.type_show}</p>
                          <p className="text-xs text-muted-foreground">{item.created_at}</p>
                          {item.desc && (
                            <SafeHtml
                              html={item.desc}
                              className="text-xs text-muted-foreground [&_a]:text-primary [&_a]:hover:underline"
                            />
                          )}
                        </div>
                      </div>
                      <div className="text-right">
                        <p className={`text-sm font-semibold ${item.num > 0 ? "text-green-600" : "text-red-600"}`}>
                          {item.num > 0 ? `+${item.num}` : item.num}
                        </p>
                        <p className="text-xs text-muted-foreground">余额 {item.balance}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* 分页 */}
              <div className="mt-4 flex items-center justify-center gap-2">
                {page > 1 && (
                  <button
                    onClick={() => goToPage(page - 1)}
                    className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                  >
                    上一页
                  </button>
                )}
                <span className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground">
                  {page}
                </span>
                {data.has_more && (
                  <button
                    onClick={() => goToPage(page + 1)}
                    className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
                  >
                    下一页
                  </button>
                )}
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {/* 排行榜链接 */}
      <div className="mt-4 text-center text-sm text-muted-foreground">
        查看{" "}
        <Link href="/balance/add" className="text-primary hover:underline">
          充值
        </Link>
        {" · "}
        <Link href="/top/rich" className="text-primary hover:underline">
          财富排行榜
        </Link>
      </div>
    </PageLayout>
  )
}

export default function BalancePage() {
  return (
    <Suspense fallback={<div className="flex min-h-40 items-center justify-center"><Loader2 className="h-6 w-6 animate-spin text-muted-foreground" /></div>}>
      <BalanceContent />
    </Suspense>
  )
}
