"use client"

import { useEffect, useState, useCallback } from "react"
import { useRouter } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Loader2 } from "lucide-react"
import { fetchAPI } from "@/lib/api"
import { useAuth } from "@/lib/auth-context"
import Link from "next/link"

interface Gift {
  id: number
  name: string
  description: string
  price: number
  remain_num: number
  buy_limit: number
  typ_show: string
  can_exchange: boolean
}

interface ExchangeRecord {
  id: number
  gift_id: number
  remark: string
  expire_time: string
  created_at: string
}

// 兑换码/履约说明存在 remark 里，是用户兑换后真正要拿到的东西
function exchangeLabel(gift: Gift): string {
  if (gift.remain_num <= 0) return "已兑完"
  if (gift.buy_limit <= 0) return "已达上限"
  if (!gift.can_exchange) return "积分不足"
  return "兑换"
}

function isZeroTime(value: string): boolean {
  return !value || value.startsWith("0001-")
}

export default function GiftPage() {
  const router = useRouter()
  const { isLoggedIn, isLoading: authLoading } = useAuth()
  const [gifts, setGifts] = useState<Gift[]>([])
  const [records, setRecords] = useState<ExchangeRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [tab, setTab] = useState<"list" | "mine">("list")
  const [exchanging, setExchanging] = useState<number | null>(null)
  const [confirmGift, setConfirmGift] = useState<Gift | null>(null)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace("/account/login?redirect=/gift")
    }
  }, [isLoggedIn, authLoading, router])

  const loadData = useCallback(() => {
    setLoading(true)
    setError(null)
    const endpoint = tab === "list" ? "/gift" : "/gift/mine"
    fetchAPI<{ gifts: Gift[] } | { records: ExchangeRecord[] }>(endpoint, { credentials: "include" })
      .then((data) => {
        if (tab === "list") {
          setGifts((data as { gifts: Gift[] }).gifts || [])
        } else {
          setRecords((data as { records: ExchangeRecord[] }).records || [])
        }
      })
      .catch((err: Error) => {
        setError(err.message)
      })
      .finally(() => setLoading(false))
  }, [tab])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleExchange = async () => {
    if (!confirmGift) return
    const giftId = confirmGift.id
    setExchanging(giftId)
    setConfirmGift(null)
    try {
      const form = new URLSearchParams()
      form.set("gift_id", String(giftId))
      await fetchAPI<null>("/gift/exchange", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: form.toString(),
        credentials: "include",
      })
      // 兑换码在「兑换记录」里，切过去让用户直接拿到（tab 变化会触发重新拉取）
      setTab("mine")
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "兑换失败"
      setError(message)
    } finally {
      setExchanging(null)
    }
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader title="礼品兑换" breadcrumbs={[{ label: "积分商城" }]} />
      <div className="mb-4 flex gap-2">
        <button
          onClick={() => setTab("list")}
          className={`rounded-md px-4 py-2 text-sm font-medium ${
            tab === "list" ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground"
          }`}
        >
          礼品列表
        </button>
        <button
          onClick={() => setTab("mine")}
          className={`rounded-md px-4 py-2 text-sm font-medium ${
            tab === "mine" ? "bg-primary text-primary-foreground" : "bg-secondary text-secondary-foreground"
          }`}
        >
          兑换记录
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-sm text-destructive">
          {error}
        </div>
      ) : tab === "list" ? (
        <div className="space-y-4">
          {gifts.length === 0 ? (
            <div className="rounded-lg border border-dashed border-border py-16 text-center">
              <p className="text-muted-foreground">暂无可兑换礼品</p>
            </div>
          ) : (
            gifts.map((gift) => (
              <Card key={gift.id} className="transition-all hover:shadow-sm">
                <CardContent className="p-4">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="font-medium">{gift.name}</h3>
                      <p className="mt-1 text-sm text-muted-foreground">{gift.description}</p>
                      <div className="mt-2 flex items-center gap-3">
                        <span className="text-sm font-semibold text-amber-600">{gift.price} 铜币</span>
                        <span className="text-sm text-muted-foreground">· 剩余 {gift.remain_num}</span>
                        {gift.typ_show && (
                          <span className="text-sm text-muted-foreground">· {gift.typ_show}</span>
                        )}
                      </div>
                    </div>
                    <Button
                      size="sm"
                      onClick={() => setConfirmGift(gift)}
                      disabled={!gift.can_exchange || exchanging === gift.id}
                    >
                      {exchanging === gift.id ? "兑换中..." : exchangeLabel(gift)}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      ) : (
        <div className="space-y-4">
          {records.length === 0 ? (
            <div className="rounded-lg border border-dashed border-border py-16 text-center">
              <p className="text-muted-foreground">暂无兑换记录</p>
            </div>
          ) : (
            records.map((record) => (
              <Card key={record.id}>
                <CardContent className="p-4">
                  <div className="flex items-start justify-between gap-4">
                    <div className="min-w-0 flex-1">
                      <p className="break-all font-mono text-sm">{record.remark || "处理中"}</p>
                      <p className="mt-1 text-sm text-muted-foreground">
                        兑换于 {record.created_at}
                      </p>
                    </div>
                    {!isZeroTime(record.expire_time) && (
                      <span className="shrink-0 text-sm text-muted-foreground">
                        有效期至 {record.expire_time.slice(0, 10)}
                      </span>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {/* 兑换确认对话框 */}
      <AlertDialog open={!!confirmGift} onOpenChange={(open) => !open && setConfirmGift(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认兑换</AlertDialogTitle>
            <AlertDialogDescription>
              确定要使用 <span className="font-semibold text-amber-600">{confirmGift?.price} 铜币</span> 兑换{" "}
              <span className="font-medium">{confirmGift?.name}</span> 吗？兑换后铜币将无法退回。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction onClick={handleExchange}>确认兑换</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <div className="mt-4 text-center text-sm text-muted-foreground">
        查看{" "}
        <Link href="/balance" className="text-primary hover:underline">
          积分余额
        </Link>
        {" · "}
        <Link href="/mission" className="text-primary hover:underline">
          每日任务
        </Link>
      </div>
    </PageLayout>
  )
}
