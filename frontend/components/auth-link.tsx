"use client"

import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"

interface AuthLinkProps {
  /** 已登录时跳转目标 */
  href: string
  children: React.ReactNode
  className?: string
  size?: "sm" | "default" | "lg" | "icon"
}

/**
 * AuthLink：已登录时跳转到 href，未登录时跳转到登录页（携带 redirect 参数）
 */
export function AuthLink({ href, children, className, size = "sm" }: AuthLinkProps) {
  const router = useRouter()

  function handleClick() {
    const uid = typeof window !== "undefined" ? localStorage.getItem("uid") : null
    if (uid) {
      router.push(href)
    } else {
      router.push(`/account/login?redirect=${encodeURIComponent(href)}`)
    }
  }

  return (
    <Button
      size={size}
      className={className}
      onClick={handleClick}
    >
      {children}
    </Button>
  )
}
