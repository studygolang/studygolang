import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import { ExternalLink } from "lucide-react"
import { fetchAPI } from "@/lib/api"

export const metadata: Metadata = {
  title: "友情链接 - Go语言中文网",
  description: "Go语言中文网友情链接，推荐优质的 Go 语言和编程相关网站",
}

interface FriendLink {
  id: number
  name: string
  url: string
  logo: string
}

async function getFriendLinks(): Promise<FriendLink[]> {
  try {
    const data = await fetchAPI<{ links: FriendLink[] }>(
      "/sidebar/friend/links",
      { cache: "no-store" }
    )
    return data.links ?? []
  } catch {
    return []
  }
}

export default async function LinksPage() {
  const links = await getFriendLinks()

  return (
    <PageLayout>
      <PageHeader
        title="友情链接"
        description="推荐优质的 Go 语言和编程相关网站"
        breadcrumbs={[{ label: "友情链接" }]}
      />

      {links.length > 0 ? (
        <Card>
          <CardContent className="p-6">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {links.map((link) => (
                <a
                  key={link.id}
                  href={link.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="group flex items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-muted/50"
                >
                  {link.logo ? (
                    <img
                      src={link.logo}
                      alt={link.name}
                      className="h-10 w-10 rounded"
                    />
                  ) : (
                    <div className="flex h-10 w-10 items-center justify-center rounded bg-primary/10 text-primary">
                      <ExternalLink className="h-5 w-5" />
                    </div>
                  )}
                  <span className="font-medium group-hover:text-primary">
                    {link.name}
                  </span>
                </a>
              ))}
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <p className="text-sm text-muted-foreground">暂无友情链接</p>
        </div>
      )}
    </PageLayout>
  )
}
