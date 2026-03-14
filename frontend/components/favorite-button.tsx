'use client'

import { useState } from 'react'
import { Bookmark } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

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

  const handleFavorite = async () => {
    if (loading) return

    setLoading(true)
    const newFavorited = !favorited

    try {
      const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8090'
      const res = await fetch(
        `${base}/api/v1/favorites/${objid}?objtype=${objtype}&collect=${newFavorited ? 1 : 0}`,
        {
          method: 'POST',
          credentials: 'include',
        }
      )

      const json = await res.json()

      if (json.code === 0) {
        setFavorited(newFavorited)
      } else if (json.code === 600) {
        window.location.href = '/login'
      } else {
        alert(json.msg || '操作失败')
      }
    } catch (error) {
      console.error('收藏失败:', error)
      alert('操作失败，请稍后重试')
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
