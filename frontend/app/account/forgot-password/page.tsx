import type { Metadata } from "next"
import Link from "next/link"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Mail } from "lucide-react"

export const metadata: Metadata = {
  title: "找回密码 - Go语言中文网",
  description: "通过邮件找回你的 Go语言中文网 账号密码",
}

export default function ForgotPasswordPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <Link href="/" className="text-2xl font-bold text-primary">
            Go语言中文网
          </Link>
          <p className="mt-2 text-sm text-muted-foreground">找回你的账号</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-xl">找回密码</CardTitle>
            <CardDescription>通过注册邮箱找回密码</CardDescription>
          </CardHeader>

          <CardContent className="space-y-4">
            <div className="flex flex-col items-center gap-4 rounded-lg bg-muted/50 p-6 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Mail className="h-6 w-6 text-primary" />
              </div>
              <div className="space-y-1">
                <p className="text-sm font-medium text-foreground">请发送邮件至管理员</p>
                <p className="text-sm text-muted-foreground">
                  联系{" "}
                  <a
                    href="mailto:polaris@studygolang.com"
                    className="font-medium text-primary hover:text-primary/80"
                  >
                    polaris@studygolang.com
                  </a>
                  {" "}，说明你的用户名，我们将帮助你重置密码。
                </p>
              </div>
            </div>
          </CardContent>

          <CardFooter>
            <Button asChild variant="outline" className="w-full">
              <Link href="/account/login">返回登录</Link>
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}
