import type { Metadata } from "next"
import Link from "next/link"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { SafeHtml } from "@/components/safe-html"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPI, needRelogin } from "@/lib/api"
import { Coins, Smartphone } from "lucide-react"

export const dynamic = "force-dynamic"

export const metadata: Metadata = {
  title: "充值 - Go语言中文网",
  description: "通过支付宝或微信充值获得积分",
}

const priceTable = [
  { amount: 5, coins: 2000 },
  { amount: 10, coins: 5000 },
  { amount: 20, coins: 12000 },
  { amount: 30, coins: 16000 },
  { amount: 50, coins: 30000 },
]

interface BalanceDetail {
  id: number
  desc: string
  type_show: string
  num: number
  balance: number
  created_at: string
}

interface RechargeData {
  balance: number
  details: BalanceDetail[]
  recharge_amount: number
}

async function getRechargeData(page: number, cookieHeader: string): Promise<RechargeData | null> {
  try {
    return await fetchAPI<RechargeData>(`/balance/add?p=${page}`, {
      cache: "no-store",
      headers: { Cookie: cookieHeader },
    })
  } catch (err: unknown) {
    if (needRelogin(err)) {
      const target = page > 1 ? `/balance/add?p=${page}` : "/balance/add"
      redirect(`/account/login?redirect=${encodeURIComponent(target)}`)
    }
    return null
  }
}

export default async function BalanceAddPage({
  searchParams,
}: {
  searchParams: Promise<{ p?: string }>
}) {
  const [{ p }, requestCookies] = await Promise.all([searchParams, cookies()])
  const page = Math.max(1, parseInt(p ?? "1", 10) || 1)
  const data = await getRechargeData(page, requestCookies.toString())
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="充值"
        breadcrumbs={[
          { label: "首页", href: "/" },
          { label: "账户余额", href: "/balance" },
          { label: "充值" },
        ]}
      />

      <div className="space-y-6">
        {/* 充值说明 */}
        <Card>
          <CardContent className="p-6 space-y-4">
            <p className="text-sm text-muted-foreground">
              你可以通过支付宝或微信转账方式向我们充值。目前的实现方式是手工的，
              我们在收到你的充值之后，就会尽快向你的账户发放铜币及开通充值会员的额外功能。
            </p>
            <p className="text-sm font-medium text-foreground">
              请在支付宝或微信的付款说明中填入你的 Go语言中文网 用户名。
            </p>

            {/* 支付方式 */}
            <div className="grid gap-6 sm:grid-cols-2">
              <div>
                <h3 className="mb-2 flex items-center gap-2 text-base font-semibold">
                  <Smartphone className="h-4 w-4 text-blue-500" />
                  支付宝
                </h3>
                <img
                  src="https://static.golangjob.cn/170605/9d11b43988fbbb42a5da9f970a0f6818.png"
                  alt="支付宝收款码"
                  className="w-[220px] rounded-lg border"
                />
              </div>
              <div>
                <h3 className="mb-2 flex items-center gap-2 text-base font-semibold">
                  <Smartphone className="h-4 w-4 text-green-500" />
                  微信
                </h3>
                <img
                  src="https://static.golangjob.cn/img/wxpay.png"
                  alt="微信收款码"
                  className="w-[220px] rounded-lg border"
                />
              </div>
            </div>

            <p className="text-sm text-muted-foreground">
              如果你在扫码支付的过程中忘记了填入用户名，可以在支付结束后联系我们，
              并附上交易尾号的末 4 位。
            </p>
          </CardContent>
        </Card>

        {/* 充值额度 */}
        <Card>
          <CardContent className="p-6">
            <h3 className="mb-4 flex items-center gap-2 text-base font-semibold">
              <Coins className="h-5 w-5 text-amber-500" />
              充值额度
            </h3>
            <div className="space-y-2">
              {priceTable.map((item) => (
                <div
                  key={item.amount}
                  className="flex items-center justify-between rounded-lg border px-4 py-3"
                >
                  <span className="text-sm font-medium">
                    充值 {item.amount} 元
                  </span>
                  <span className="text-sm font-semibold text-amber-600">
                    获得 {item.coins.toLocaleString()} 铜币
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            {data ? (
              <>
                <div className="mb-6 grid gap-4 rounded-lg bg-secondary/50 p-4 sm:grid-cols-2">
                  <div>
                    <p className="text-xs text-muted-foreground">历史充值金额总计</p>
                    <p className="mt-1 text-2xl font-semibold text-foreground">￥{data.recharge_amount}</p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">当前账户余额</p>
                    <p className="mt-1 text-2xl font-semibold text-amber-600">{data.balance} 铜币</p>
                  </div>
                </div>
                <h3 className="mb-2 text-base font-semibold">充值记录</h3>
                {data.details.length > 0 ? (
                  <div className="divide-y divide-border">
                    {data.details.map((item) => (
                      <div key={item.id} className="flex flex-col gap-2 py-3 sm:flex-row sm:items-start sm:justify-between">
                        <div>
                          <p className="text-sm font-medium">{item.type_show}</p>
                          <p className="text-xs text-muted-foreground">{item.created_at}</p>
                          {item.desc && (
                            <SafeHtml
                              html={item.desc}
                              className="mt-1 text-xs text-muted-foreground [&_a]:text-primary [&_a]:hover:underline"
                            />
                          )}
                        </div>
                        <div className="shrink-0 sm:text-right">
                          <p className={`text-sm font-semibold ${item.num >= 0 ? "text-green-600" : "text-red-600"}`}>
                            {item.num >= 0 ? `+${item.num}` : item.num}
                          </p>
                          <p className="text-xs text-muted-foreground">余额 {item.balance}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="py-8 text-center text-sm text-muted-foreground">暂无充值记录</p>
                )}
              </>
            ) : (
              <p className="py-8 text-center text-sm text-muted-foreground">充值记录加载失败，请稍后重试</p>
            )}
          </CardContent>
        </Card>

        {/* VIP 特权 */}
        <Card>
          <CardContent className="p-6 space-y-3">
            <h3 className="text-base font-semibold">充值会员额外功能</h3>
            <ul className="list-inside list-disc space-y-1 text-sm text-muted-foreground">
              <li>置顶自己的主题或文章 1 天 / <span className="text-amber-600">每次消耗 30,000 铜币</span></li>
            </ul>
          </CardContent>
        </Card>

        <div className="flex justify-center gap-4 text-sm text-muted-foreground">
          <Link href="/balance" className="text-primary hover:underline">
            查看余额
          </Link>
          <Link href="/top/rich" className="text-primary hover:underline">
            财富排行榜
          </Link>
        </div>
      </div>
    </PageLayout>
  )
}
