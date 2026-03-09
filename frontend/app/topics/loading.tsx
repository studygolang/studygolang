import { Skeleton } from "@/components/ui/skeleton"
import { PageLayout } from "@/components/page-layout"

export default function TopicsLoading() {
  return (
    <PageLayout>
      <div className="space-y-4">
        {/* Header skeleton */}
        <div className="space-y-2">
          <Skeleton className="h-8 w-32" />
          <Skeleton className="h-4 w-64" />
        </div>

        {/* Tab filter skeleton */}
        <div className="flex gap-4 border-b border-border pb-2">
          <Skeleton className="h-6 w-12" />
          <Skeleton className="h-6 w-16" />
          <Skeleton className="h-6 w-16" />
          <Skeleton className="h-6 w-12" />
        </div>

        {/* Topic list skeleton */}
        <div className="rounded-lg border border-border bg-card p-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="flex gap-4 py-4 border-b border-border last:border-0">
              <Skeleton className="hidden h-9 w-9 rounded-full sm:block shrink-0" />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-4 w-3/4" />
                <div className="flex gap-3">
                  <Skeleton className="h-4 w-16" />
                  <Skeleton className="h-4 w-12" />
                  <Skeleton className="h-4 w-16" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </PageLayout>
  )
}
