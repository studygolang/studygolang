'use client'

import { useState } from 'react'
import { ThumbsUp } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { likeAPI, needRelogin } from '@/lib/api'
import { toast } from 'sonner'
import { useRouter } from 'next/navigation'

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
  const router = useRouter()

  const handleLike = async () => {
    if (loading) return

    setLoading(true)
    const newLiked = !liked

    try {
      await likeAPI.toggle(objid, objtype, newLiked)
      setLiked(newLiked)
      setCount(prev => newLiked ? prev + 1 : Math.max(0, prev - 1))
    } catch (err: unknown) {
      if (needRelogin(err)) {
        router.push('/account/login')
      } else {
        toast.error(err instanceof Error ? err.message : '操作失败')
      }
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
        liked && 'text-primary hover:text-primary/80',
        className
      )}
    >
      <ThumbsUp
        className={cn(
          'h-4 w-4 transition-all',
          liked && 'fill-current'
        )}
      />
      <span>{count > 0 ? count : '赞'}</span>
    </Button>
  )
}
