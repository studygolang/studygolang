"use client"

import { useState } from "react"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { saveDraft, draftKey } from "@/lib/draft"

interface FormActionsProps {
  submitting: boolean
  onCancel: () => void
  onPublish: () => void
  contentType: string
  title: string
  content: string
}

export function FormActions({
  submitting,
  onCancel,
  onPublish,
  contentType,
  title,
  content,
}: FormActionsProps) {
  const [savedAt, setSavedAt] = useState<number | null>(null)

  function handleSaveDraft() {
    const key = draftKey(contentType)
    saveDraft(key, { title, content })
    setSavedAt(Date.now())
  }

  return (
    <div className="flex items-center justify-between border-t border-border px-5 py-4">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onCancel}
          className="text-sm text-muted-foreground transition-colors hover:text-foreground"
        >
          取消
        </button>
        {savedAt && (
          <span className="text-xs text-muted-foreground">
            草稿已保存于 {new Date(savedAt).toLocaleTimeString()}
          </span>
        )}
      </div>
      <div className="flex items-center gap-3">
        <Button
          variant="outline"
          size="sm"
          disabled={submitting || (!title.trim() && !content.trim())}
          onClick={handleSaveDraft}
        >
          保存草稿
        </Button>
        <Button size="sm" disabled={submitting} onClick={onPublish} data-publish-btn className="gap-2">
          {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
          {submitting ? "发布中..." : "发布"}
        </Button>
      </div>
    </div>
  )
}
