"use client"

import React, { useEffect, useState } from "react"
import { useRouter, useParams } from "next/navigation"
import { Loader2, Info } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { bookWriteAPI } from "@/lib/api"
import { toast } from "sonner"

export default function BookModifyPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  // 表单字段
  const [name, setName] = useState("")
  const [ename, setEname] = useState("")
  const [cover, setCover] = useState("")
  const [author, setAuthor] = useState("")
  const [translator, setTranslator] = useState("")
  const [pubDate, setPubDate] = useState("")
  const [price, setPrice] = useState("")
  const [onlineUrl, setOnlineUrl] = useState("")
  const [downloadUrl, setDownloadUrl] = useState("")
  const [buyUrl, setBuyUrl] = useState("")
  const [desc, setDesc] = useState("")
  const [tags, setTags] = useState("")
  const [catalogue, setCatalogue] = useState("")
  const [isFree, setIsFree] = useState(false)

  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/books/modify/${id}`)
      return
    }
    fetchBook()
  }, [id, isLoggedIn, authLoading, router])

  async function fetchBook() {
    try {
      const data = await bookWriteAPI.getEdit(id)
      const book = data.book
      setName(book.name || "")
      setEname(book.ename || "")
      setCover(book.cover || "")
      setAuthor(book.author || "")
      setTranslator(book.translator || "")
      setPubDate(book.pub_date || "")
      setPrice(String(book.price || ""))
      setOnlineUrl(book.online_url || "")
      setDownloadUrl(book.download_url || "")
      setBuyUrl(book.buy_url || "")
      setDesc(book.desc || "")
      setTags(book.tags || "")
      setCatalogue(book.catalogue || "")
      setIsFree(book.is_free || false)
    } catch (err: any) {
      toast.error(err.message || "加载失败")
    } finally {
      setLoading(false)
    }
  }

  async function handleSave() {
    if (!name.trim()) {
      toast.error("请填写书名")
      return
    }
    if (!author.trim()) {
      toast.error("请填写作者")
      return
    }

    setSubmitting(true)
    try {
      await bookWriteAPI.update(id, {
        name: name.trim(),
        ename: ename.trim(),
        cover: cover.trim(),
        author: author.trim(),
        translator: translator.trim(),
        pub_date: pubDate.trim(),
        price: price.trim(),
        online_url: onlineUrl.trim(),
        download_url: downloadUrl.trim(),
        buy_url: buyUrl.trim(),
        desc: desc.trim(),
        tags: tags.trim(),
        catalogue: catalogue.trim(),
        is_free: isFree ? "1" : "0",
      })
      toast.success("保存成功")
      setTimeout(() => router.push(`/books/${id}`), 1200)
    } catch (err: any) {
      toast.error(err.message || "保存失败")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <SiteHeader />
        <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6 flex items-center justify-center">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </main>
        <SiteFooter />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6">
        <div className="mb-6">
          <h1 className="text-xl font-semibold text-foreground">编辑图书</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改图书信息后点击保存</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="name" className="text-sm font-medium">
                  书名 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="name"
                  placeholder="中文书名"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  maxLength={100}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="ename" className="text-sm font-medium">
                  英文书名
                </Label>
                <Input
                  id="ename"
                  placeholder="英文书名"
                  value={ename}
                  onChange={(e) => setEname(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="author" className="text-sm font-medium">
                  作者 <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="author"
                  placeholder="作者姓名"
                  value={author}
                  onChange={(e) => setAuthor(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="translator" className="text-sm font-medium">
                  译者
                </Label>
                <Input
                  id="translator"
                  placeholder="译者姓名"
                  value={translator}
                  onChange={(e) => setTranslator(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="cover" className="text-sm font-medium">
                  封面图片
                </Label>
                <Input
                  id="cover"
                  placeholder="封面图片地址"
                  value={cover}
                  onChange={(e) => setCover(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="pubDate" className="text-sm font-medium">
                  出版日期
                </Label>
                <Input
                  id="pubDate"
                  placeholder="如：2020-01"
                  value={pubDate}
                  onChange={(e) => setPubDate(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="price" className="text-sm font-medium">
                  价格
                </Label>
                <Input
                  id="price"
                  placeholder="如：99.00"
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5 flex items-end">
                <label className="flex items-center gap-2 cursor-pointer pb-2">
                  <input
                    type="checkbox"
                    checked={isFree}
                    onChange={(e) => setIsFree(e.target.checked)}
                    className="h-4 w-4"
                  />
                  <span className="text-sm font-medium">免费资源</span>
                </label>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="onlineUrl" className="text-sm font-medium">
                在线阅读
              </Label>
              <Input
                id="onlineUrl"
                placeholder="在线阅读地址"
                value={onlineUrl}
                onChange={(e) => setOnlineUrl(e.target.value)}
                className="h-10"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="downloadUrl" className="text-sm font-medium">
                  下载地址
                </Label>
                <Input
                  id="downloadUrl"
                  placeholder="下载链接"
                  value={downloadUrl}
                  onChange={(e) => setDownloadUrl(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="buyUrl" className="text-sm font-medium">
                  购买链接
                </Label>
                <Input
                  id="buyUrl"
                  placeholder="购买地址"
                  value={buyUrl}
                  onChange={(e) => setBuyUrl(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="desc" className="text-sm font-medium">
                内容简介
              </Label>
              <Textarea
                id="desc"
                placeholder="请输入图书简介"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
                rows={4}
                className="resize-none"
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="catalogue" className="text-sm font-medium">
                目录
              </Label>
              <Textarea
                id="catalogue"
                placeholder="图书目录"
                value={catalogue}
                onChange={(e) => setCatalogue(e.target.value)}
                rows={6}
                className="resize-none"
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="tags" className="text-sm font-medium">
                标签
                <span className="ml-1 font-normal text-muted-foreground">（逗号分隔）</span>
              </Label>
              <Input
                id="tags"
                placeholder="如：Go,编程,Web开发"
                value={tags}
                onChange={(e) => setTags(e.target.value)}
                className="h-10"
              />
            </div>

            <div className="flex items-start gap-2 rounded-md bg-muted/50 p-3 text-xs text-muted-foreground">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
              <span>
                带有 <strong className="text-foreground">*</strong> 的为必填项，保存后可以继续编辑。
              </span>
            </div>
          </div>

          <div className="flex items-center justify-between border-t border-border px-5 py-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              取消
            </button>
            <Button size="sm" disabled={submitting} onClick={handleSave} className="gap-2">
              {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
              {submitting ? "保存中..." : "保存"}
            </Button>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
