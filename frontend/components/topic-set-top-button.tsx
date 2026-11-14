"use client"

import { useState } from 'react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { topicWriteAPI } from '@/lib/api'

interface TopicSetTopButtonProps {
  tid: number | string
  isTop: boolean
  onSuccess: () => void
}

export function TopicSetTopButton({ tid, isTop, onSuccess }: TopicSetTopButtonProps) {
  const [loading, setLoading] = useState(false)

  const handleSetTop = async () => {
    if (loading || isTop) return

    setLoading(true)
    try {
      await topicWriteAPI.setTop(tid)
      toast.success('置顶成功')
      onSuccess()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : '操作失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button
      variant="outline"
      size="sm"
      onClick={handleSetTop}
      disabled={loading || isTop}
    >
      {/* 后端只有置顶接口，取消置顶由定时任务按 top_time 到期自动处理 */}
      {isTop ? '已置顶' : '置顶'}
    </Button>
  )
}
