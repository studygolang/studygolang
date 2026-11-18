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
import { useAuth } from "@/lib/auth-context"
import { gcttAPI, captchaAPI } from '@/lib/api'
import { CaptchaField } from '@/components/captcha-field'

export default function GCTTPublishPage() {
  const router = useRouter()
  const { isLoggedIn, isLoading: authLoading } = useAuth()
  const [loading, setLoading] = useState(true)
  const [isTranslator, setIsTranslator] = useState(false)
  const [applying, setApplying] = useState(false)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [cover, setCover] = useState('')
  const [author, setAuthor] = useState('')
  const [authorUrl, setAuthorUrl] = useState('')
  const [translator, setTranslator] = useState('')
  const [checker, setChecker] = useState('')
  const [sourceUrl, setSourceUrl] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [captchaRequired, setCaptchaRequired] = useState(false)
  const [captchaId, setCaptchaId] = useState('')
  const [captchaSolution, setCaptchaSolution] = useState('')
  const [captchaKey, setCaptchaKey] = useState(0)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      toast.error('请先登录')
      router.push('/account/login')
      return
    }

    async function checkTranslator() {
      try {
        const data = await gcttAPI.getMe()
        setIsTranslator(data.is_translator)
        captchaAPI.required()
          .then((captcha) => setCaptchaRequired(captcha.required))
          .catch(() => {})
      } catch {
        toast.error('获取信息失败')
        router.push('/')
      } finally {
        setLoading(false)
      }
    }

    checkTranslator()
  }, [authLoading, isLoggedIn, router])

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
    if (!author.trim() || !authorUrl.trim() || !translator.trim() || !checker.trim() || !sourceUrl.trim()) {
      toast.error('请填写原作者、作者链接、译者、校对和原文链接')
      return
    }
    if (captchaRequired && !captchaSolution.trim()) {
      toast.error('请填写验证码')
      return
    }

    if (submitting) return

    setSubmitting(true)
    try {
      const captcha = captchaRequired
        ? { id: captchaId, solution: captchaSolution.trim() }
        : undefined
      const data = await gcttAPI.publish(title.trim(), content.trim(), {
        cover: cover.trim(),
        author: author.trim(),
        author_url: authorUrl.trim(),
        translator: translator.trim(),
        checker: checker.trim(),
        url: sourceUrl.trim(),
      }, captcha)
      toast.success('发布成功')
      router.push(`/articles/${data.id}`)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : '发布失败')
      if (captchaRequired) {
        setCaptchaSolution('')
        setCaptchaKey((key) => key + 1)
      }
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

            <div className="space-y-2">
              <Label htmlFor="cover">封面</Label>
              <Input
                id="cover"
                value={cover}
                onChange={(e) => setCover(e.target.value)}
                placeholder="封面图片地址（可选）"
              />
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="author">原作者</Label>
                <Input
                  id="author"
                  value={author}
                  onChange={(e) => setAuthor(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="author-url">作者链接</Label>
                <Input
                  id="author-url"
                  type="url"
                  value={authorUrl}
                  onChange={(e) => setAuthorUrl(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="translator">译者 GitHub ID</Label>
                <Input
                  id="translator"
                  value={translator}
                  onChange={(e) => setTranslator(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="checker">校对者 GitHub ID</Label>
                <Input
                  id="checker"
                  value={checker}
                  onChange={(e) => setChecker(e.target.value)}
                  placeholder="多个以逗号分隔"
                  required
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="source-url">原文链接</Label>
              <Input
                id="source-url"
                type="url"
                value={sourceUrl}
                onChange={(e) => setSourceUrl(e.target.value)}
                required
              />
            </div>

            {captchaRequired && (
              <CaptchaField
                key={captchaKey}
                value={captchaSolution}
                onChange={setCaptchaSolution}
                captchaId={captchaId}
                onCaptchaIdChange={setCaptchaId}
              />
            )}

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
