import { NextRequest, NextResponse } from "next/server"

// 旧版 /wr 外链包装的兼容实现：老内容中大量 /wr?u=xxx 链接依赖此路由
// shortcut: 不再做 iframe 包装与异步 X-Frame-Options 检测，统一 utm 跳转；需要 iframe 包装时再升级
export async function GET(req: NextRequest) {
  const raw = req.nextUrl.searchParams.get("u")
  let target = raw || new URL("/", req.nextUrl).toString()
  if (raw && !raw.includes("studygolang.com")) {
    target += (target.includes("?") ? "&" : "?") + "utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com"
  }
  return NextResponse.redirect(target, 303)
}
