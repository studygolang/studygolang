"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { Pencil } from "lucide-react"
import { Button } from "@/components/ui/button"
import { SubjectFollowButton } from "@/components/subject-follow-button"
import { SubjectContributeDialog } from "@/components/subject-contribute-dialog"
import { useAuth } from "@/lib/auth-context"

interface SubjectActionsProps {
  sid: number
  ownerUid: number
  followed: boolean
}

export function SubjectActions({ sid, ownerUid, followed }: SubjectActionsProps) {
  const router = useRouter()
  const { user, isLoggedIn } = useAuth()

  const canEdit = isLoggedIn && (user?.is_root || user?.uid === ownerUid)

  return (
    <div className="flex shrink-0 items-center gap-2">
      {canEdit && (
        <Button asChild variant="ghost" size="sm" className="gap-1.5">
          <Link href={`/subject/modify/${sid}`}>
            <Pencil className="h-4 w-4" />
            编辑
          </Link>
        </Button>
      )}
      {isLoggedIn && (
        <SubjectContributeDialog sid={sid} onSuccess={() => router.refresh()} />
      )}
      <SubjectFollowButton sid={sid} initiallyFollowed={followed} />
    </div>
  )
}