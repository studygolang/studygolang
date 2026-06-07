"use client"

import "./blocknote.css"
import { Suspense, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { SiteHeader } from "@/components/site-header"
import { SiteFooter } from "@/components/site-footer"
import { useAuth } from "@/lib/auth-context"
import { PublishHeader } from "./components/publish-header"
import { TopicFields } from "./components/topic-fields"
import { ProjectFields } from "./components/project-fields"
import { ResourceFields } from "./components/resource-fields"
import { BookFields } from "./components/book-fields"
import { TagsInput } from "./components/tags-input"
import { EditorSection } from "./components/editor-section"
import { FormActions } from "./components/form-actions"
import { usePublish } from "./components/use-publish"
import type { ContentType } from "./components/use-publish"

export default function PublishPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background flex items-center justify-center"><p>加载中...</p></div>}>
      <PublishContent />
    </Suspense>
  )
}

function PublishContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { isLoggedIn, isLoading: authLoading } = useAuth()

  // 从 URL 参数读取初始内容类型
  const initialType = (searchParams.get("type") || "topic") as ContentType
  const validTypes: ContentType[] = ["topic", "article", "project", "resource", "book"]

  const publishState = usePublish({
    initialType: validTypes.includes(initialType) ? initialType : "topic"
  })

  // 未登录重定向
  useEffect(() => {
    if (authLoading) return
    if (!isLoggedIn) {
      router.replace("/account/login?redirect=/publish")
    }
  }, [isLoggedIn, authLoading, router])

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-3xl px-4 py-8 lg:px-6">
        {/* 页面标题和内容类型选择 */}
        <PublishHeader
          contentType={publishState.contentType}
          onContentTypeChange={publishState.setContentType}
        />

        {/* 表单卡片 */}
        <div className="rounded-lg border border-border bg-card">
          <div className="p-5 space-y-5">
            {/* 标题 */}
            <div className="space-y-1.5">
              <Label htmlFor="title" className="text-sm font-medium">
                {publishState.contentType === "project" ? "项目名称" : publishState.contentType === "book" ? "书名" : "标题"} <span className="text-destructive">*</span>
              </Label>
              <Input
                id="title"
                placeholder={
                  publishState.contentType === "topic"
                    ? "请输入主题标题，简明扼要地描述你的问题或话题"
                    : publishState.contentType === "article"
                    ? "请输入文章标题"
                    : publishState.contentType === "resource"
                    ? "请输入资源标题"
                    : publishState.contentType === "book"
                    ? "请输入图书名称"
                    : "请输入项目名称"
                }
                value={publishState.title}
                onChange={(e) => publishState.setTitle(e.target.value)}
                maxLength={120}
                className="h-10"
              />
            </div>

            {/* 项目特定字段 */}
            {publishState.contentType === "project" && (
              <ProjectFields
                src={publishState.src}
                setSrc={publishState.setSrc}
                category={publishState.category}
                setCategory={publishState.setCategory}
                lang={publishState.lang}
                setLang={publishState.setLang}
                home={publishState.home}
                setHome={publishState.setHome}
                doc={publishState.doc}
                setDoc={publishState.setDoc}
                licence={publishState.licence}
                setLicence={publishState.setLicence}
                os={publishState.os}
                setOs={publishState.setOs}
              />
            )}

            {/* 资源特定字段 */}
            {publishState.contentType === "resource" && (
              <ResourceFields
                resourceForm={publishState.resourceForm}
                setResourceForm={publishState.setResourceForm}
                resourceUrl={publishState.resourceUrl}
                setResourceUrl={publishState.setResourceUrl}
              />
            )}

            {/* 图书特定字段 */}
            {publishState.contentType === "book" && (
              <BookFields
                bookAuthor={publishState.bookAuthor}
                setBookAuthor={publishState.setBookAuthor}
                bookTranslator={publishState.bookTranslator}
                setBookTranslator={publishState.setBookTranslator}
                bookCover={publishState.bookCover}
                setBookCover={publishState.setBookCover}
                bookPubDate={publishState.bookPubDate}
                setBookPubDate={publishState.setBookPubDate}
                bookLang={publishState.bookLang}
                setBookLang={publishState.setBookLang}
                bookIsFree={publishState.bookIsFree}
                setBookIsFree={publishState.setBookIsFree}
                bookOnlineUrl={publishState.bookOnlineUrl}
                setBookOnlineUrl={publishState.setBookOnlineUrl}
                bookDownloadUrl={publishState.bookDownloadUrl}
                setBookDownloadUrl={publishState.setBookDownloadUrl}
                bookBuyUrl={publishState.bookBuyUrl}
                setBookBuyUrl={publishState.setBookBuyUrl}
                bookPrice={publishState.bookPrice}
                setBookPrice={publishState.setBookPrice}
              />
            )}

            {/* 节点 + 标签（并排） */}
            <div className="flex gap-4">
              {/* 节点（仅话题类型显示） */}
              {publishState.contentType === "topic" && (
                <TopicFields
                  nid={publishState.nid}
                  setNid={publishState.setNid}
                  nodeGroups={publishState.nodeGroups}
                />
              )}

              {/* 标签 */}
              <TagsInput
                tags={publishState.tags}
                setTags={publishState.setTags}
                tagInput={publishState.tagInput}
                setTagInput={publishState.setTagInput}
                onTagKeyDown={publishState.handleTagKeyDown}
              />
            </div>

            {/* BlockNote 编辑器 + 提示信息 */}
            <EditorSection
              contentType={publishState.contentType}
              resourceForm={publishState.resourceForm}
              content={publishState.content}
              setContent={publishState.setContent}
            />

            {/* 错误提示 */}
            {publishState.error && <p className="text-sm text-destructive">{publishState.error}</p>}
          </div>

          {/* 草稿恢复提示 */}
          {publishState.draftData && (
            <div className="mx-5 mb-4 flex items-center gap-3 rounded-md border border-blue-200 bg-blue-50 px-4 py-3 dark:border-blue-800 dark:bg-blue-950">
              <div className="flex-1">
                <p className="text-sm font-medium text-blue-800 dark:text-blue-200">
                  检测到未完成的草稿
                </p>
                <p className="text-xs text-blue-600 dark:text-blue-400">
                  标题：{publishState.draftData.title || "（无标题）"}
                  {" · "}
                  保存于 {new Date(publishState.draftData.savedAt).toLocaleString()}
                </p>
              </div>
              <button
                type="button"
                onClick={publishState.restoreDraft}
                className="rounded-md bg-blue-600 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-blue-700"
              >
                恢复草稿
              </button>
              <button
                type="button"
                onClick={publishState.discardDraft}
                className="text-xs text-muted-foreground transition-colors hover:text-foreground"
              >
                放弃
              </button>
            </div>
          )}

          {/* 底部操作栏 */}
          <FormActions
            submitting={publishState.submitting}
            onCancel={publishState.handleCancel}
            onPublish={publishState.handlePublish}
            contentType={publishState.contentType}
            title={publishState.title}
            content={publishState.content}
          />
        </div>
      </main>

      <SiteFooter />
    </div>
  )
}
