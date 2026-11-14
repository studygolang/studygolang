import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import type { User } from "@/lib/types"
import { encodeProjectUri } from "@/lib/utils"

// 后端统一返回 favorites 数组，元素类型由 objtype 决定，
// 字段是 model.Article / Topic / Resource / OpenProject 的并集子集。
interface FavoriteItem {
  id?: number
  tid?: number
  title?: string
  name?: string
  uri?: string
  desc?: string
  author?: string
  ctime?: string
  viewnum?: number
}

interface FavoriteData {
  user: User
  objtype: number
  total: number
  page: number
  has_more: boolean
  favorites: FavoriteItem[]
}

import { fetchAPI } from "@/lib/api"

// 后端 model/comment.go iota: TypeTopic=0, TypeArticle=1, TypeResource=2, TypeProject=4
// UI 用连续 1-4 显示，但实际传给后端的是 iota 编号
const FAVORITE_TABS = [
  { value: '1', label: '文章', objtype: 1 },
  { value: '2', label: '话题', objtype: 0 },
  { value: '3', label: '资源', objtype: 2 },
  { value: '4', label: '项目', objtype: 4 },
] as const

// 后端 objtype → UI tab value（用于 SSR 时根据 URL objtype 选中默认 tab）
const OBJTYPE_TO_TAB: Record<number, string> = Object.fromEntries(
  FAVORITE_TABS.map((t) => [t.objtype, t.value])
)

function itemHref(objtype: number, item: FavoriteItem): string {
  switch (objtype) {
    case 0:
      return `/topics/${item.tid}`
    case 2:
      return `/resources/${item.id}`
    case 4:
      return `/p/${encodeProjectUri(item.uri || item.id || "")}`
    default:
      return `/articles/${item.id}`
  }
}

async function getFavorites(
  username: string,
  objtype: number,
  page: number
): Promise<FavoriteData | null> {
  try {
    return await fetchAPI<FavoriteData>(
      `/users/${username}/favorites?objtype=${objtype}&p=${page}`,
      { cache: 'no-store' }
    )
  } catch {
    return null
  }
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ username: string }>
}): Promise<Metadata> {
  const { username } = await params
  const data = await getFavorites(username, 1, 1)

  if (!data?.user) {
    return {
      title: "用户不存在 - Go语言中文网",
    }
  }

  return {
    title: `${data.user.username} 的收藏 - Go语言中文网`,
    description: `查看 ${data.user.username} 在 Go语言中文网的收藏内容`,
  }
}

export default async function FavoritesPage({
  params,
  searchParams,
}: {
  params: Promise<{ username: string }>
  searchParams: Promise<{ objtype?: string; p?: string }>
}) {
  const { username } = await params
  const { objtype: objtypeStr, p: pageStr } = await searchParams
  // 默认显示"文章"（TypeArticle=1，与后端 default 一致）
  const objtype = parseInt(objtypeStr || '1', 10)
  const page = Math.max(1, parseInt(pageStr || '1', 10) || 1)
  // URL 里的 objtype 是后端 iota 值（0,1,2,4），UI tab 用连续 1-4 编号
  const activeTab = OBJTYPE_TO_TAB[objtype] ?? '1'

  const data = await getFavorites(username, objtype, page)

  if (!data?.user) {
    return (
      <PageLayout>
        <PageHeader title="用户不存在" />
        <div className="text-center py-12">
          <p className="text-muted-foreground">该用户不存在</p>
        </div>
      </PageLayout>
    )
  }

  const { user, favorites, total } = data
  const emptyLabel = FAVORITE_TABS.find((t) => t.objtype === objtype)?.label ?? '内容'

  return (
    <PageLayout>
      <PageHeader title={`${user.username} 的收藏`} />

      <div className="mb-6 flex items-center gap-4">
        <Avatar className="h-16 w-16">
          <AvatarImage src={user.avatar} alt={user.username} />
          <AvatarFallback>{user.username[0]?.toUpperCase()}</AvatarFallback>
        </Avatar>
        <div>
          <h2 className="text-xl font-semibold">{user.username}</h2>
          <p className="text-sm text-muted-foreground">共收藏 {total} 项内容</p>
        </div>
      </div>

      <Tabs value={activeTab} className="space-y-4">
        <TabsList>
          {FAVORITE_TABS.map((tab) => (
            <TabsTrigger key={tab.value} value={tab.value} asChild>
              <Link href={`/favorites/${username}?objtype=${tab.objtype}`}>{tab.label}</Link>
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="space-y-4">
        {favorites && favorites.length > 0 ? (
          favorites.map((item) => (
            <Card key={item.tid ?? item.id}>
              <CardHeader>
                <CardTitle>
                  <Link
                    href={itemHref(objtype, item)}
                    className="hover:text-primary transition-colors"
                  >
                    {item.title || item.name}
                  </Link>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                  {item.author && <span>作者: {item.author}</span>}
                  {item.ctime && <span>{item.ctime}</span>}
                  {item.viewnum !== undefined && <span>{item.viewnum} 浏览</span>}
                </div>
                {item.desc && (
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">{item.desc}</p>
                )}
              </CardContent>
            </Card>
          ))
        ) : (
          <p className="text-center text-muted-foreground py-8">暂无收藏的{emptyLabel}</p>
        )}
      </div>

      {(page > 1 || data.has_more) && (
        <div className="mt-6 flex items-center justify-center gap-4 text-sm">
          {page > 1 && (
            <Link
              href={`/favorites/${username}?objtype=${objtype}&p=${page - 1}`}
              className="text-primary hover:underline"
            >
              上一页
            </Link>
          )}
          <span className="text-muted-foreground">第 {page} 页</span>
          {data.has_more && (
            <Link
              href={`/favorites/${username}?objtype=${objtype}&p=${page + 1}`}
              className="text-primary hover:underline"
            >
              下一页
            </Link>
          )}
        </div>
      )}
    </PageLayout>
  )
}
