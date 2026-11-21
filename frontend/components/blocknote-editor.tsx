"use client"

import { useEffect, useRef } from "react"
import { useCreateBlockNote } from "@blocknote/react"
import { BlockNoteView } from "@blocknote/mantine"
import type { BlockNoteEditor as BlockNoteEditorType } from "@blocknote/core"
import "@blocknote/mantine/style.css"
import "@blocknote/core/fonts/inter.css"
import { quickUploadImage } from "@/lib/upload"

interface BlockNoteEditorProps {
  content: string
  onChange: (markdown: string) => void
  placeholder?: string
}

/**
 * 自定义图片上传函数，将文件上传到后端并返回 URL
 * BlockNote 会在用户粘贴/拖拽图片时调用此函数
 */
async function uploadFile(file: File): Promise<string> {
  const result = await quickUploadImage(file)
  return result.url
}

export function BlockNoteEditor({ content, onChange, placeholder }: BlockNoteEditorProps) {
  // 创建编辑器实例，启用图片上传
  const editor = useCreateBlockNote({
    uploadFile,
  })

  const lastMarkdownRef = useRef<string | null>(null)

  useEffect(() => {
    if (!content || content === lastMarkdownRef.current) return

    let cancelled = false
    const loadContent = async () => {
      try {
        const blocks = await editor.tryParseMarkdownToBlocks(content)
        if (cancelled) return
        lastMarkdownRef.current = content
        editor.replaceBlocks(editor.document, blocks)
      } catch (error) {
        console.error("Failed to parse markdown:", error)
      }
    }

    loadContent()
    return () => {
      cancelled = true
    }
  }, [editor, content])

  useEffect(() => {
    const handleChange = async () => {
      try {
        const markdown = await editor.blocksToMarkdownLossy(editor.document)
        lastMarkdownRef.current = markdown
        onChange(markdown)
      } catch (error) {
        console.error("Failed to convert to markdown:", error)
      }
    }

    return editor.onChange(handleChange)
  }, [editor, onChange])

  if (!editor) {
    return <div className="p-4 text-sm text-muted-foreground">加载编辑器...</div>
  }

  return (
    <div className="blocknote-wrapper min-h-[400px]">
      <BlockNoteView
        editor={editor}
        theme="light"
        data-theming-css-variables-demo
        slashMenu={true}
        formattingToolbar={true}
        linkToolbar={true}
        sideMenu={true}
        filePanel={true}
        tableHandles={true}
      />
    </div>
  )
}
