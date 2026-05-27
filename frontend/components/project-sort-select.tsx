"use client"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useRouter } from "next/navigation"

const SORT_OPTIONS = [
  { value: "hot", label: "最热" },
  { value: "latest", label: "最新" },
  { value: "noreply", label: "零回复" },
]

interface ProjectSortSelectProps {
  currentSort: string
}

export function ProjectSortSelect({ currentSort }: ProjectSortSelectProps) {
  const router = useRouter()

  function handleChange(value: string) {
    const params = new URLSearchParams()
    params.set("sort", value)
    params.set("p", "1")
    router.push(`/projects?${params.toString()}`)
  }

  return (
    <Select value={currentSort} onValueChange={handleChange}>
      <SelectTrigger size="sm" className="w-28">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {SORT_OPTIONS.map((opt) => (
          <SelectItem key={opt.value} value={opt.value}>
            {opt.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
