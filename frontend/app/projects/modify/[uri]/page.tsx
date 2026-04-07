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
import { projectWriteAPI, userAPI } from "@/lib/api"
import { toast } from "sonner"

export default function ProjectModifyPage() {
  const router = useRouter()
  const params = useParams()
  const uri = params.uri as string

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  // 表单字段
  const [name, setName] = useState("")
  const [category, setCategory] = useState("")
  const [logo, setLogo] = useState("")
  const [home, setHome] = useState("")
  const [doc, setDoc] = useState("")
  const [src, setSrc] = useState("")
  const [download, setDownload] = useState("")
  const [desc, setDesc] = useState("")
  const [tags, setTags] = useState("")
  const [author, setAuthor] = useState("")
  const [licence, setLicence] = useState("")
  const [lang, setLang] = useState("")
  const [os, setOS] = useState("")

  useEffect(() => {
    async function checkAuth() {
      const me = await userAPI.getMe()
      if (!me) {
        router.replace(`/account/login?redirect=/projects/modify/${uri}`)
        return
      }
      fetchProject()
    }
    checkAuth()
  }, [uri])

  async function fetchProject() {
    try {
      const data = await projectWriteAPI.getEdit(uri)
      const project = data.project
      setName(project.name || "")
      setCategory(project.category || "")
      setLogo(project.logo || "")
      setHome(project.home || "")
      setDoc(project.doc || "")
      setSrc(project.src || "")
      setDownload(project.download || "")
      setDesc(project.desc || "")
      setTags(project.tags || "")
      setAuthor(project.author || "")
      setLicence(project.licence || "")
      setLang(project.lang || "")
      setOS(project.os || "")
    } catch (err: any) {
      toast.error(err.message || "加载失败")
    } finally {
      setLoading(false)
    }
  }

  async function handleSave() {
    if (!name.trim()) {
      toast.error("请填写项目名称")
      return
    }
    if (!src.trim()) {
      toast.error("请填写源码地址")
      return
    }

    setSubmitting(true)
    try {
      await projectWriteAPI.update(uri, {
        name: name.trim(),
        category: category.trim(),
        logo: logo.trim(),
        home: home.trim(),
        doc: doc.trim(),
        src: src.trim(),
        download: download.trim(),
        desc: desc.trim(),
        tags: tags.trim(),
        author: author.trim(),
        licence: licence.trim(),
        lang: lang.trim(),
        os: os.trim(),
      })
      toast.success("保存成功")
      setTimeout(() => router.push(`/projects/${uri}`), 1200)
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
          <h1 className="text-xl font-semibold text-foreground">编辑项目</h1>
          <p className="mt-1 text-sm text-muted-foreground">修改项目信息后点击保存</p>
        </div>

        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            <div className="space-y-1.5">
              <Label htmlFor="name" className="text-sm font-medium">
                项目名称 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="name"
                placeholder="请输入项目名称"
                value={name}
                onChange={(e) => setName(e.target.value)}
                maxLength={100}
                className="h-10"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="category" className="text-sm font-medium">
                  分类
                </Label>
                <Input
                  id="category"
                  placeholder="如：Web框架"
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="logo" className="text-sm font-medium">
                  Logo URL
                </Label>
                <Input
                  id="logo"
                  placeholder="项目 Logo 图片地址"
                  value={logo}
                  onChange={(e) => setLogo(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="src" className="text-sm font-medium">
                源码地址 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="src"
                placeholder="如：https://github.com/xxx/xxx"
                value={src}
                onChange={(e) => setSrc(e.target.value)}
                className="h-10"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="home" className="text-sm font-medium">
                  主页
                </Label>
                <Input
                  id="home"
                  placeholder="项目主页"
                  value={home}
                  onChange={(e) => setHome(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="doc" className="text-sm font-medium">
                  文档地址
                </Label>
                <Input
                  id="doc"
                  placeholder="项目文档"
                  value={doc}
                  onChange={(e) => setDoc(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="download" className="text-sm font-medium">
                下载地址
              </Label>
              <Input
                id="download"
                placeholder="项目下载地址"
                value={download}
                onChange={(e) => setDownload(e.target.value)}
                className="h-10"
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="desc" className="text-sm font-medium">
                项目描述
              </Label>
              <Textarea
                id="desc"
                placeholder="请输入项目描述"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
                rows={4}
                className="resize-none"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="author" className="text-sm font-medium">
                  作者
                </Label>
                <Input
                  id="author"
                  placeholder="项目作者"
                  value={author}
                  onChange={(e) => setAuthor(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="licence" className="text-sm font-medium">
                  许可证
                </Label>
                <Input
                  id="licence"
                  placeholder="如：MIT"
                  value={licence}
                  onChange={(e) => setLicence(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <Label htmlFor="lang" className="text-sm font-medium">
                  开发语言
                </Label>
                <Input
                  id="lang"
                  placeholder="如：Go"
                  value={lang}
                  onChange={(e) => setLang(e.target.value)}
                  className="h-10"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="os" className="text-sm font-medium">
                  操作系统
                </Label>
                <Input
                  id="os"
                  placeholder="如：跨平台"
                  value={os}
                  onChange={(e) => setOS(e.target.value)}
                  className="h-10"
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="tags" className="text-sm font-medium">
                标签
                <span className="ml-1 font-normal text-muted-foreground">（逗号分隔）</span>
              </Label>
              <Input
                id="tags"
                placeholder="如：web,framework,go"
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
