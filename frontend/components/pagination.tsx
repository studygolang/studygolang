import Link from "next/link"

interface PaginationProps {
  currentPage: number
  hasMore: boolean
  buildUrl: (page: number) => string
}

export function Pagination({ currentPage, hasMore, buildUrl }: PaginationProps) {
  return (
    <div className="mt-8 flex items-center justify-center gap-2">
      {currentPage > 1 && (
        <Link
          href={buildUrl(currentPage - 1)}
          className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
        >
          上一页
        </Link>
      )}
      <span className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground">
        {currentPage}
      </span>
      {hasMore && (
        <Link
          href={buildUrl(currentPage + 1)}
          className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
        >
          下一页
        </Link>
      )}
    </div>
  )
}
