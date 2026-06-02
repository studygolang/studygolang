"use client"

import dynamic from "next/dynamic"
import { Label } from "@/components/ui/label"
import { Info } from "lucide-react"

// 动态导入 BlockNote 编辑器,避免 SSR 问题
const BlockNoteEditor = dynamic(
  () => import("@/components/blocknote-editor").then((mod) => ({ default: mod.BlockNoteEditor })),
  {
    ssr: false,
    loading: () => <div className="p-4 text-sm text-muted-foreground">加载编辑器...</div>
  }
)

interface EditorSectionProps {
  contentType: "topic" | "article" | "project" | "resource" | "book"
  resourceForm?: string
  content: string
  setContent: (content: string) => void
}

export function EditorSection({ contentType, resourceForm, content, setContent }: EditorSectionProps) {
  // 资源类型只在"包括内容"时显示编辑器
  if (contentType === "resource" && resourceForm !== "content") {
    return null
  }

  const getLabel = () => {
    if (contentType === "project") return "项目描述"
    if (contentType === "resource") return "资源内容"
    if (contentType === "book") return "图书简介"
    return "内容"
  }

  const getPlaceholder = () => {
    if (contentType === "project") return "请详细描述项目的功能特点、使用场景……"
    if (contentType === "resource") return "请详细描述资源内容……"
    if (contentType === "book") return "请简要介绍这本图书的主要内容、适合的读者群体……"
    return "请详细描述你的问题或话题……"
  }

  return (
    <>
      <div className="space-y-1.5">
        <Label className="text-sm font-medium">
          {getLabel()} <span className="text-destructive">*</span>
        </Label>
        <div className="rounded-md border border-border overflow-hidden">
          <BlockNoteEditor
            content={content}
            onChange={setContent}
            placeholder={getPlaceholder()}
          />
        </div>
      </div>

      {/* 提示信息 */}
      <div className="flex items-start gap-2 rounded-md bg-muted/50 px-3 py-2.5 text-xs text-muted-foreground">
        <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-primary" />
        <span>
          内容支持 <strong className="text-foreground">Markdown</strong> 语法，发布前请确认标题准确、内容完整、节点正确，发布后可以继续编辑。
        </span>
      </div>
    </>
  )
}
