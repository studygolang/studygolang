"use client"

import Link from "next/link"
import { Heart, Send, Settings } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/lib/auth-context"

export function UserActions({ username }: { username: string }) {
  const { user, isLoading } = useAuth()

  return (
    <div className="mt-4 flex flex-wrap gap-2">
      <Button asChild size="sm" variant="outline" className="gap-1.5">
        <Link href={`/favorites/${username}`}>
          <Heart className="h-3.5 w-3.5" />
          {user?.username === username ? "我的收藏" : "TA的收藏"}
        </Link>
      </Button>
      {!isLoading && user && user.username !== username && (
        <Button asChild size="sm" variant="outline" className="gap-1.5">
          <Link href={`/messages/send?username=${encodeURIComponent(username)}`}>
            <Send className="h-3.5 w-3.5" />
            发送消息
          </Link>
        </Button>
      )}
      {!isLoading && user?.username === username && (
        <Button asChild size="sm" variant="outline" className="gap-1.5">
          <Link href={`/user/${username}/settings`}>
            <Settings className="h-3.5 w-3.5" />
            编辑信息
          </Link>
        </Button>
      )}
    </div>
  )
}
