"use client"

import Link from "next/link"
import { Pencil } from "lucide-react"
import { useAuth } from "@/lib/auth-context"

interface WikiEditButtonProps {
  uri: string
  ownerUid: number
}

// 与后端 logic.CanEdit(*model.Wiki) 对齐：root/admin 或所有者可编辑
// （编辑时间窗由后端校验，超时后保存会被拒绝）
export function WikiEditButton({ uri, ownerUid }: WikiEditButtonProps) {
  const { user } = useAuth()
  const canEdit = !!user && (user.is_root || user.is_admin || user.uid === ownerUid)
  if (!canEdit) return null

  return (
    <Link
      href={`/wiki/${uri}/edit`}
      className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-primary"
    >
      <Pencil className="h-3.5 w-3.5" />
      编辑
    </Link>
  )
}