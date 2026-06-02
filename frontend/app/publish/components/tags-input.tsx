"use client"

import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"

interface TagsInputProps {
  tags: string[]
  setTags: (tags: string[]) => void
  tagInput: string
  setTagInput: (input: string) => void
  onTagKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void
}

export function TagsInput({ tags, setTags, tagInput, setTagInput, onTagKeyDown }: TagsInputProps) {
  return (
    <div className="min-w-0 flex-1 space-y-1.5">
      <Label htmlFor="tags" className="text-sm font-medium">
        标签
        <span className="ml-1 font-normal text-muted-foreground">（最多5个，回车添加）</span>
      </Label>
      <div className="flex min-h-10 flex-wrap items-center gap-1.5 rounded-md border border-input bg-background px-3 py-1.5 text-sm focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-0">
        {tags.map((tag) => (
          <Badge
            key={tag}
            variant="secondary"
            className="gap-1 py-0.5 text-xs cursor-pointer"
            onClick={() => setTags(tags.filter((t) => t !== tag))}
          >
            {tag} ×
          </Badge>
        ))}
        {tags.length < 5 && (
          <input
            id="tags"
            className="min-w-16 flex-1 bg-transparent outline-none placeholder:text-muted-foreground"
            placeholder={tags.length === 0 ? "输入标签" : ""}
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={onTagKeyDown}
          />
        )}
      </div>
    </div>
  )
}
