"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Coins, TrendingUp, TrendingDown, Loader2 } from "lucide-react"
import Link from "next/link"
import { fetchAPI } from "@/lib/api"

interface BalanceDetail {
  id: number
  desc: string
  num: number
  ctime: string
}

interface BalanceData {
  details: BalanceDetail[]
  total: number
  page: number
  has_more: boolean
}

async function fetchBalance(page: number): Promise<BalanceData> {
  return fetchAPI<BalanceData>(`/balance?p=${page}`, { credentials: "include" })
}

export default function BalancePage() {
  const router = useRouter()
  const [data, setData] = useState<BalanceData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)

  useEffect(() => {
    setLoading(true)
    setError(null)
    fetchBalance(page)
      .then(setData)
      .catch((err: Error) => {
        if (err.message.includes("未登录") || err.message.includes("401")) {
          router.push("/account/login?redirect=/balance")
        } else {
          setError(err.message)
        }
      })
      .finally(() => setLoading(false))
  }, [page, router])

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
              <span className="text-4xl font-bold text-amber-600">{data.total}</span>
              <span className="text-sm text-muted-foreground">积分</span>
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
                          <p className="text-sm font-medium">{item.desc}</p>
                          <p className="text-xs text-muted-foreground">{item.ctime}</p>
                        </div>
                      </div>
                      <span className={`text-sm font-semibold ${item.num > 0 ? "text-green-600" : "text-red-600"}`}>
                        {item.num > 0 ? `+${item.num}` : item.num}
                      </span>
                    </div>
                  ))}
                </div>
              )}

              {/* 分页 */}
              <div className="mt-4 flex items-center justify-center gap-2">
                {page > 1 && (
                  <button
                    onClick={() => setPage((p) => p - 1)}
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
                    onClick={() => setPage((p) => p + 1)}
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
        <Link href="/top/rich" className="text-primary hover:underline">
          财富排行榜
        </Link>
      </div>
    </PageLayout>
  )
}
