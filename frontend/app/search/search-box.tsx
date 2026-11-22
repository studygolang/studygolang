"use client"

import { useRouter, useSearchParams } from "next/navigation"
import { useState, useTransition } from "react"
import { Search } from "lucide-react"

interface SearchBoxProps {
  defaultValue?: string
  defaultType?: string
}

// 搜索字段：空=全文，title=标题，author=作者（与后端 DoSearch 的 df 对应）
const fieldOptions = [
  { value: "", label: "全文" },
  { value: "title", label: "标题" },
  { value: "author", label: "作者" },
]

export function SearchBox({ defaultValue = "", defaultType = "" }: SearchBoxProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [query, setQuery] = useState(defaultValue)
  const [field, setField] = useState(defaultType)
  const [isPending, startTransition] = useTransition()

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const q = query.trim()
    if (!q) return
    startTransition(() => {
      const params = new URLSearchParams(searchParams.toString())
      params.set("q", q)
      params.delete("page")
      params.delete("p")
      params.delete("f")
      if (field) {
        params.set("type", field)
      } else {
        params.delete("type")
      }
      router.push(`/search?${params.toString()}`)
    })
  }

  return (
    <form onSubmit={handleSubmit} className="flex w-full max-w-2xl gap-2">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="搜索话题、文章、项目、资源..."
          className="w-full rounded-md border border-border bg-background py-2.5 pl-9 pr-4 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>
      <select
        value={field}
        onChange={(e) => setField(e.target.value)}
        aria-label="搜索字段"
        className="rounded-md border border-border bg-background px-3 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
      >
        {fieldOptions.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      <button
        type="submit"
        disabled={isPending || !query.trim()}
        className="rounded-md bg-primary px-5 py-2.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {isPending ? "搜索中..." : "搜索"}
      </button>
    </form>
  )
}
