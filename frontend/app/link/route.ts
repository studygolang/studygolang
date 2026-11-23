import { NextRequest, NextResponse } from "next/server"

// 旧版 /link 广告外链跳转的兼容实现
export async function GET(req: NextRequest) {
  const raw = req.nextUrl.searchParams.get("url")
  let target = raw || new URL("/", req.nextUrl).toString()
  if (raw) {
    target += (target.includes("?") ? "&" : "?") + "utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com"
  }
  return NextResponse.redirect(target, 303)
}
