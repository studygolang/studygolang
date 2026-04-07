"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Plus, Check } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { subjectAPI } from "@/lib/api"
import { toast } from "sonner"

interface SubjectFollowButtonProps {
  sid: number
  initiallyFollowed: boolean
  onToggle?: (followed: boolean) => void
  className?: string
}

export function SubjectFollowButton({
  sid,
  initiallyFollowed,
  onToggle,
  className,
}: SubjectFollowButtonProps) {
  const router = useRouter()
  const [followed, setFollowed] = useState(initiallyFollowed)
  const [loading, setLoading] = useState(false)

  const handleFollow = async () => {
    if (loading) return

    setLoading(true)
    const newFollowed = !followed

    try {
      const result = await subjectAPI.follow(sid)
      setFollowed(result.followed)
      onToggle?.(result.followed)
      toast.success(result.followed ? "关注成功" : "已取消关注")
    } catch (error) {
      const message = error instanceof Error ? error.message : "操作失败"
      if (message.includes("登录") || message.includes("600")) {
        router.push("/account/login")
        return
      }
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button
      variant={followed ? "outline" : "default"}
      size="sm"
      onClick={handleFollow}
      disabled={loading}
      className={cn("gap-1.5", className)}
    >
      {followed ? (
        <>
          <Check className="h-4 w-4" />
          <span>已关注</span>
        </>
      ) : (
        <>
          <Plus className="h-4 w-4" />
          <span>关注</span>
        </>
      )}
    </Button>
  )
}
