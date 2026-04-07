"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { accountAPI } from "@/lib/api"
import { toast } from "sonner"
import { Trash2 } from "lucide-react"

interface SocialUnbindButtonProps {
  bindId: number
  platform: string
  onSuccess: () => void
}

export function SocialUnbindButton({ bindId, platform, onSuccess }: SocialUnbindButtonProps) {
  const [submitting, setSubmitting] = useState(false)

  const handleUnbind = async () => {
    setSubmitting(true)
    try {
      await accountAPI.socialUnbind(bindId, platform)
      toast.success("解绑成功")
      onSuccess()
    } catch (err: any) {
      toast.error(err.message || "解绑失败")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="destructive" size="sm" disabled={submitting}>
          <Trash2 className="h-4 w-4 mr-1.5" />
          解绑
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认解绑</AlertDialogTitle>
          <AlertDialogDescription>
            确定要解绑 {platform} 账号吗？解绑后你将无法使用该账号登录。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction
            onClick={handleUnbind}
            disabled={submitting}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {submitting ? "解绑中..." : "确认解绑"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
