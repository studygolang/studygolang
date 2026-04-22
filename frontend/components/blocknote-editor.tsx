"use client"

import { useEffect } from "react"
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

  // 初始化时加载内容
  useEffect(() => {
    if (!editor || !content) return

    const loadContent = async () => {
      try {
        const blocks = await editor.tryParseMarkdownToBlocks(content)
        editor.replaceBlocks(editor.document, blocks)
      } catch (error) {
        console.error("Failed to parse markdown:", error)
      }
    }

    loadContent()
  }, []) // 只在初始化时加载一次

  // 监听编辑器变化，转换为 Markdown
  useEffect(() => {
    if (!editor) return

    const handleChange = async () => {
      try {
        const markdown = await editor.blocksToMarkdownLossy(editor.document)
        onChange(markdown)
      } catch (error) {
        console.error("Failed to convert to markdown:", error)
      }
    }

    // 订阅编辑器变化
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
