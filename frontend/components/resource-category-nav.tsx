import Link from "next/link"
import { fetchAPINullable } from "@/lib/api"
import { cn } from "@/lib/utils"

interface ResourceCategory {
  catid: number
  name: string
}

export async function ResourceCategoryNav({ currentCatid = 0 }: { currentCatid?: number }) {
  const data = await fetchAPINullable<{ categories: ResourceCategory[] }>("/resources/categories", {
    cache: "no-store",
  })
  const categories = data?.categories
  if (!categories?.length) return null

  return (
    <nav aria-label="资源分类" className="mb-4 flex flex-wrap gap-2 rounded-lg border border-border bg-card p-4">
      <Link
        href="/resources"
        aria-current={currentCatid === 0 ? "page" : undefined}
        className={cn(
          "rounded-full px-3 py-1.5 text-sm font-medium transition-colors",
          currentCatid === 0
            ? "bg-primary text-primary-foreground"
            : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
        )}
      >
        全部
      </Link>
      {categories.map((category) => (
        <Link
          key={category.catid}
          href={`/resources/cat/${category.catid}`}
          aria-current={currentCatid === category.catid ? "page" : undefined}
          className={cn(
            "rounded-full px-3 py-1.5 text-sm font-medium transition-colors",
            currentCatid === category.catid
              ? "bg-primary text-primary-foreground"
              : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
          )}
        >
          {category.name}
        </Link>
      ))}
    </nav>
  )
}
