"use client"

import { useState, useEffect } from "react"
import { FileText, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { subjectAPI } from "@/lib/api"
import { toast } from "sonner"

interface Article {
  id: number
  title: string
}

interface SubjectContributeDialogProps {
  sid: number
  onSuccess: () => void
}

export function SubjectContributeDialog({
  sid,
  onSuccess,
}: SubjectContributeDialogProps) {
  const [open, setOpen] = useState(false)
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState<number | null>(null)

  useEffect(() => {
    if (open) {
      loadArticles()
    }
  }, [open])

  async function loadArticles() {
    setLoading(true)
    try {
      const result = await subjectAPI.myArticles()
      setArticles(result.articles || [])
    } catch (error) {
      const message = error instanceof Error ? error.message : "加载失败"
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  async function handleContribute(aid: number) {
    setSubmitting(aid)
    try {
      await subjectAPI.contribute(sid, aid)
      toast.success("投稿成功")
      setOpen(false)
      onSuccess()
    } catch (error) {
      const message = error instanceof Error ? error.message : "投稿失败"
      toast.error(message)
    } finally {
      setSubmitting(null)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="gap-1.5">
          <FileText className="h-4 w-4" />
          <span>投稿文章</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>投稿文章到专栏</DialogTitle>
          <DialogDescription>选择一篇文章投稿到专栏</DialogDescription>
        </DialogHeader>

        {loading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : articles.length === 0 ? (
          <div className="py-8 text-center text-sm text-muted-foreground">
            暂无可投稿的文章
          </div>
        ) : (
          <div className="space-y-2">
            {articles.map((article) => (
              <div
                key={article.id}
                className="flex items-center justify-between rounded-lg border border-border p-3"
              >
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium">{article.title}</p>
                </div>
                <Button
                  size="sm"
                  onClick={() => handleContribute(article.id)}
                  disabled={submitting === article.id}
                  className="ml-3 shrink-0"
                >
                  {submitting === article.id ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    "投稿"
                  )}
                </Button>
              </div>
            ))}
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
