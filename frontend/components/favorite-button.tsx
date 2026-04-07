'use client'

import { useState } from 'react'
import { Bookmark } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { favoriteAPI } from '@/lib/api'
import { toast } from 'sonner'
import { useRouter } from 'next/navigation'

interface FavoriteButtonProps {
  objid: number
  objtype: number
  initialFavorited?: boolean
  className?: string
}

export function FavoriteButton({
  objid,
  objtype,
  initialFavorited = false,
  className,
}: FavoriteButtonProps) {
  const [favorited, setFavorited] = useState(initialFavorited)
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  const handleFavorite = async () => {
    if (loading) return

    setLoading(true)
    const newFavorited = !favorited

    try {
      await favoriteAPI.toggle(objid, objtype, newFavorited)
      setFavorited(newFavorited)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : ''
      if (msg.includes('未登录') || msg.includes('600')) {
        router.push('/account/login')
      } else {
        toast.error(msg || '操作失败')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={handleFavorite}
      disabled={loading}
      className={cn(
        'gap-1 transition-all',
        favorited && 'text-yellow-600 hover:text-yellow-700',
        className
      )}
    >
      <Bookmark
        className={cn(
          'h-4 w-4 transition-all',
          favorited && 'fill-current'
        )}
      />
      <span>{favorited ? '已收藏' : '收藏'}</span>
    </Button>
  )
}
