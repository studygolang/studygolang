"use client"

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { SiteHeader } from '@/components/site-header'
import { SiteFooter } from '@/components/site-footer'
import { userAPI, gcttAPI } from '@/lib/api'

export default function GCTTPublishPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [isTranslator, setIsTranslator] = useState(false)
  const [applying, setApplying] = useState(false)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    async function checkAuth() {
      try {
        // 检查登录状态
        const user = await userAPI.getMe()
        if (!user) {
          toast.error('请先登录')
          router.push('/account/login')
          return
        }

        // 检查译者身份
        const data = await gcttAPI.getMe()
        setIsTranslator(data.is_translator)
      } catch {
        toast.error('获取信息失败')
        router.push('/')
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [router])

  const handleApply = async () => {
    if (applying) return

    setApplying(true)
    try {
      await gcttAPI.apply()
      toast.success('申请成功，已成为译者')
      setIsTranslator(true)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : '申请失败')
    } finally {
      setApplying(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!title.trim() || !content.trim()) {
      toast.error('请填写标题和内容')
      return
    }

    if (submitting) return

    setSubmitting(true)
    try {
      const data = await gcttAPI.publish(title.trim(), content.trim())
      toast.success('发布成功')
      router.push(`/articles/${data.id}`)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : '发布失败')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <>
        <SiteHeader />
        <main className="container mx-auto px-4 py-8">
          <div className="text-center text-muted-foreground">加载中...</div>
        </main>
        <SiteFooter />
      </>
    )
  }

  if (!isTranslator) {
    return (
      <>
        <SiteHeader />
        <main className="container mx-auto px-4 py-8">
          <div className="mx-auto max-w-md space-y-6 text-center">
            <div className="space-y-2">
              <h1 className="text-2xl font-bold">GCTT 翻译计划</h1>
              <p className="text-sm text-muted-foreground">
                您还不是译者，需要先申请成为译者才能发布翻译文章
              </p>
            </div>
            <Button onClick={handleApply} disabled={applying}>
              {applying ? '申请中...' : '申请成为译者'}
            </Button>
          </div>
        </main>
        <SiteFooter />
      </>
    )
  }

  return (
    <>
      <SiteHeader />
      <main className="container mx-auto px-4 py-8">
        <div className="mx-auto max-w-3xl space-y-6">
          <div className="space-y-2">
            <h1 className="text-2xl font-bold">发布翻译文章</h1>
            <p className="text-sm text-muted-foreground">
              发布高质量 Go 相关的翻译文章
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">文章标题</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="请输入翻译后的标题"
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="content">文章内容</Label>
              <Textarea
                id="content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="请输入翻译后的文章内容（Markdown 格式）"
                className="min-h-[400px]"
                required
              />
            </div>

            <div className="flex gap-3">
              <Button type="submit" disabled={submitting}>
                {submitting ? '发布中...' : '发布文章'}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.back()}
              >
                取消
              </Button>
            </div>
          </form>
        </div>
      </main>
      <SiteFooter />
    </>
  )
}
