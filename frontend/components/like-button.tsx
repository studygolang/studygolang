'use client'

import { useState } from 'react'
import { Heart } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface LikeButtonProps {
  objid: number
  objtype: number
  initialLiked?: boolean
  initialCount?: number
  className?: string
}

export function LikeButton({
  objid,
  objtype,
  initialLiked = false,
  initialCount = 0,
  className,
}: LikeButtonProps) {
  const [liked, setLiked] = useState(initialLiked)
  const [count, setCount] = useState(initialCount)
  const [loading, setLoading] = useState(false)

  const handleLike = async () => {
    if (loading) return

    setLoading(true)
    const newLiked = !liked

    try {
      const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8090'
      const res = await fetch(
        `${base}/api/v1/likes/${objid}?objtype=${objtype}&flag=${newLiked ? 1 : 0}`,
        {
          method: 'POST',
          credentials: 'include',
        }
      )

      const json = await res.json()

      if (json.code === 0) {
        setLiked(newLiked)
        setCount(prev => newLiked ? prev + 1 : Math.max(0, prev - 1))
      } else if (json.code === 600) {
        window.location.href = '/login'
      } else {
        alert(json.msg || '操作失败')
      }
    } catch (error) {
      console.error('点赞失败:', error)
      alert('操作失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={handleLike}
      disabled={loading}
      className={cn(
        'gap-1 transition-all',
        liked && 'text-red-500 hover:text-red-600',
        className
      )}
    >
      <Heart
        className={cn(
          'h-4 w-4 transition-all',
          liked && 'fill-current'
        )}
      />
      <span>{count > 0 ? count : ''}</span>
    </Button>
  )
}
