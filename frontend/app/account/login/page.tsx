"use client"

import { useState, useEffect, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Github } from "lucide-react"
import { useAuth } from "@/lib/auth-context"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8090"

function LoginForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const rawRedirect = searchParams.get("redirect") || "/"
  const redirect = (rawRedirect.startsWith("/") && !rawRedirect.startsWith("//")) ? rawRedirect : "/"
  const { login } = useAuth()

  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [rememberMe, setRememberMe] = useState(true)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  // 检查是否已登录
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch(`${API_BASE}/api/v1/user/me`, { credentials: "include" })
        if (res.ok) {
          const data = await res.json()
          if (data.code === 0 && data.data?.user) {
            // 已登录,跳转到首页
            router.replace(redirect)
          }
        }
      } catch {
        // 未登录,继续显示登录页面
      }
    }
    checkAuth()
  }, [router, redirect])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    try {
      const res = await fetch(`${API_BASE}/api/v1/user/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          passwd: password,
        }),
        credentials: "include",
      })

      if (res.ok) {
        const data = await res.json()
        if (data.code === 0) {
          if (data.data && data.data.uid) {
            login({
              uid: data.data.uid,
              username: data.data.username,
              name: data.data.username,
              avatar: "",
              is_root: false,
              is_vip: false,
            })
            // 用 router.replace 保持 SPA 状态（auth context 已更新），避免整页刷新丢失 context
            router.replace(redirect)
          } else {
            // data.data 为空却不报错，说明后端响应异常（如账号被冻结时只回 msg）
            setError(data.msg || "登录失败")
          }
        } else {
          setError(data.msg || "登录失败")
        }
      } else {
        setError("登录失败,请检查用户名和密码")
      }
    } catch {
      setError("网络错误,请稍后重试")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-muted/30 p-4">
      <Link href="/" className="mb-6 text-2xl font-bold text-primary hover:opacity-80 transition-opacity">
        Go语言中文网
      </Link>
      <div className="w-full max-w-4xl grid md:grid-cols-3 gap-6">
        <Card className="md:col-span-2">

          <CardHeader>
            <CardTitle>登录</CardTitle>
            <CardDescription>登录 Go语言中文网</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <Label htmlFor="username">用户名或邮箱</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="password">密码</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="请输入密码"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="remember"
                  checked={rememberMe}
                  onCheckedChange={(checked) => setRememberMe(checked as boolean)}
                />
                <Label htmlFor="remember" className="text-sm font-normal cursor-pointer">
                  记住登录状态
                </Label>
              </div>
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "登录中..." : "登录"}
              </Button>
            </form>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">第三方账号登录</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button variant="outline" className="w-full" asChild>
                <a href="/oauth/github/login">
                  <Github className="mr-2 h-4 w-4" />
                  GitHub
                </a>
              </Button>
              <Button variant="outline" className="w-full" asChild>
                <a href="/oauth/gitea/login">
                  <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M4.209 4.603c-.247 0-.525.02-.84.088-.333.07-1.14.317-1.84.87C1.627 6.107 1.17 7.16 1.17 8.44c0 1.254.472 2.462 1.279 3.6.395.555 1.015 1.37 1.778 2.134 1.06 1.06 2.236 2.023 3.496 2.854.59.39 1.17.736 1.74 1.036.296.155.566.287.824.406l.066.03c.257.117.49.217.725.303.14.05.296.1.466.148l.043.012c.17.047.335.083.517.083.275 0 .525-.06.754-.17.23-.11.434-.27.618-.47.425-.46.608-1.12.505-1.74-.044-.263-.132-.52-.26-.754-.128-.235-.296-.448-.496-.63-.32-.293-.72-.49-1.157-.564l-.042-.007c-.107-.017-.212-.046-.312-.086l-.01-.004c-.1-.04-.195-.09-.282-.15-.088-.06-.17-.13-.244-.21l-.012-.013c-.297-.33-.5-.74-.576-1.185-.015-.09-.022-.182-.022-.274 0-.14.018-.28.054-.415.11-.412.347-.78.682-1.043.15-.118.314-.218.49-.297.09-.04.183-.073.28-.1l.033-.01c.27-.074.555-.097.835-.066.376.04.73.185 1.02.42.267.215.475.5.604.828.045.09.08.185.107.283l.01.038c.063.237.083.486.058.732-.04.385-.185.75-.42 1.05-.12.153-.26.288-.42.4-.04.028-.08.053-.122.077l-.018.01c-.142.082-.287.148-.44.198l-.042.013c-.18.058-.37.09-.56.098-.13.004-.258-.003-.385-.02l-.015-.003c-.236-.035-.46-.115-.66-.235l-.006-.004c-.196-.118-.37-.273-.508-.458-.085-.114-.156-.238-.21-.37l-.008-.02c-.098-.237-.14-.496-.12-.754.006-.065.015-.13.027-.193.09-.447.335-.85.688-1.13.132-.104.276-.193.43-.264.083-.038.17-.07.258-.095l.022-.006c.256-.066.524-.076.784-.028.325.06.63.2.882.41.165.138.31.303.424.49l.01.017c.08.13.143.27.187.417l.007.025c.063.21.088.43.073.648-.013.18-.053.36-.12.53l-.01.023c-.08.188-.19.363-.327.515l-.012.013c-.175.19-.39.34-.632.438l-.015.006c-.192.076-.398.117-.606.12h-.01c-.076 0-.152-.004-.228-.013l-.007-.001c-.186-.023-.366-.08-.532-.167l-.008-.005c-.164-.087-.313-.203-.436-.344l-.005-.006c-.1-.114-.185-.243-.25-.383l-.005-.012c-.088-.19-.134-.398-.134-.608v-.008c0-.17.03-.34.087-.5.042-.115.098-.225.168-.325l.007-.01c.11-.154.25-.283.413-.378l.008-.005c.155-.09.328-.15.507-.174l.006-.001c.12-.016.243-.016.363 0l.007.001c.16.023.314.076.455.157l.005.003c.14.08.268.186.374.312l.004.005c.086.103.157.22.208.346l.004.01c.063.156.094.324.094.494v.003c0 .122-.017.243-.05.36-.032.113-.08.222-.143.322l-.004.006c-.088.138-.207.254-.348.337l-.004.003c-.127.074-.268.122-.414.14l-.004.001c-.08.01-.16.01-.24 0l-.003-.001c-.128-.017-.25-.062-.36-.132l-.003-.002c-.1-.064-.19-.148-.262-.247l-.002-.003c-.058-.08-.105-.17-.137-.266l-.002-.006c-.04-.12-.058-.248-.05-.376v-.005c.007-.09.024-.177.05-.262.025-.08.06-.158.104-.23l.002-.004c.065-.105.152-.197.254-.268l.003-.002c.09-.063.19-.11.296-.136l.004-.001c.072-.018.147-.027.222-.027h.004c.06 0 .12.007.178.02l.003.001c.08.02.157.053.227.098l.003.002c.066.043.126.097.175.16l.002.003c.04.05.073.107.097.168l.002.005c.032.08.05.167.05.255v.004c0 .05-.006.1-.018.15-.01.043-.026.085-.047.124l-.002.003c-.03.053-.07.1-.117.14l-.002.002c-.04.033-.086.06-.136.078l-.003.001c-.038.014-.08.022-.12.022h-.002c-.025 0-.05-.003-.074-.01l-.002-.001c-.03-.008-.06-.02-.086-.037l-.002-.001c-.025-.016-.048-.036-.067-.06l-.001-.001c-.016-.02-.03-.043-.04-.068l-.001-.003c-.01-.025-.016-.053-.016-.08v-.002c0-.01.002-.02.005-.03.003-.01.008-.018.014-.026l.001-.001c.008-.01.018-.018.03-.024h.001c.01-.005.022-.008.034-.008h.001"/></svg>
                  Gitea
                </a>
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base">还没有帐号？</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <Link href="/account/register" className="block text-sm text-primary hover:underline">
                注册
              </Link>
              <Link href="/account/forgot-password" className="block text-sm text-primary hover:underline">
                忘记了密码？
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default function LoginPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">加载中...</div>}>
      <LoginForm />
    </Suspense>
  )
}
