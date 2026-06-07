import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import type { User } from "@/lib/types"

// favorites API 返回的精简数据类型（字段为后端实际返回的子集）
interface FavoriteArticle {
  id: number
  title: string
  author?: string
  ctime?: string
  viewnum?: number
}

interface FavoriteTopic {
  tid: number
  title: string
  author?: string
  ctime?: string
  view?: number
}

interface FavoriteResource {
  id: number
  title: string
  author?: string
  ctime?: string
}

interface FavoriteProject {
  id: number
  name: string
  uri?: string
  description?: string
}

interface FavoriteData {
  user: User
  objtype: number
  total: number
  page: number
  has_more: boolean
  articles?: FavoriteArticle[]
  topics?: FavoriteTopic[]
  resources?: FavoriteResource[]
  projects?: FavoriteProject[]
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

async function getFavorites(username: string, objtype: number = 1): Promise<FavoriteData | null> {
  try {
    return await fetchAPI<FavoriteData>(
      `/users/${username}/favorites?objtype=${objtype}`,
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
  const data = await getFavorites(username)

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
  searchParams: Promise<{ objtype?: string }>
}) {
  const { username } = await params
  const { objtype: objtypeStr } = await searchParams
  // 默认显示"文章"（TypeArticle=1，与后端 default 一致）
  const objtype = parseInt(objtypeStr || '1', 10)
  // URL 里的 objtype 是后端 iota 值（0,1,2,4），UI tab 用连续 1-4 编号
  const activeTab = OBJTYPE_TO_TAB[objtype] ?? '1'

  const data = await getFavorites(username, objtype)

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

  const { user, articles, topics, resources, projects, total } = data

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

      <Tabs defaultValue={activeTab} className="space-y-4">
        <TabsList>
          {FAVORITE_TABS.map((tab) => (
            <TabsTrigger key={tab.value} value={tab.value} asChild>
              <Link href={`/favorites/${username}?objtype=${tab.objtype}`}>{tab.label}</Link>
            </TabsTrigger>
          ))}
        </TabsList>

        <TabsContent value="1" className="space-y-4">
          {articles && articles.length > 0 ? (
            articles.map((article) => (
              <Card key={article.id}>
                <CardHeader>
                  <CardTitle>
                    <Link
                      href={`/articles/${article.id}`}
                      className="hover:text-primary transition-colors"
                    >
                      {article.title}
                    </Link>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    {article.author && <span>作者: {article.author}</span>}
                    {article.ctime && <span>{article.ctime}</span>}
                    {article.viewnum !== undefined && <span>{article.viewnum} 浏览</span>}
                  </div>
                </CardContent>
              </Card>
            ))
          ) : (
            <p className="text-center text-muted-foreground py-8">暂无收藏的文章</p>
          )}
        </TabsContent>

        <TabsContent value="2" className="space-y-4">
          {topics && topics.length > 0 ? (
            topics.map((topic) => (
              <Card key={topic.tid}>
                <CardHeader>
                  <CardTitle>
                    <Link
                      href={`/topics/${topic.tid}`}
                      className="hover:text-primary transition-colors"
                    >
                      {topic.title}
                    </Link>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    {topic.author && <span>作者: {topic.author}</span>}
                    {topic.ctime && <span>{topic.ctime}</span>}
                    {topic.view !== undefined && <span>{topic.view} 浏览</span>}
                  </div>
                </CardContent>
              </Card>
            ))
          ) : (
            <p className="text-center text-muted-foreground py-8">暂无收藏的话题</p>
          )}
        </TabsContent>

        <TabsContent value="3" className="space-y-4">
          {resources && resources.length > 0 ? (
            resources.map((resource) => (
              <Card key={resource.id}>
                <CardHeader>
                  <CardTitle>
                    <Link
                      href={`/resources/${resource.id}`}
                      className="hover:text-primary transition-colors"
                    >
                      {resource.title}
                    </Link>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    {resource.author && <span>作者: {resource.author}</span>}
                    {resource.ctime && <span>{resource.ctime}</span>}
                  </div>
                </CardContent>
              </Card>
            ))
          ) : (
            <p className="text-center text-muted-foreground py-8">暂无收藏的资源</p>
          )}
        </TabsContent>

        <TabsContent value="4" className="space-y-4">
          {projects && projects.length > 0 ? (
            projects.map((project) => (
              <Card key={project.id}>
                <CardHeader>
                  <CardTitle>
                    <Link
                      href={`/p/${project.uri || project.id}`}
                      className="hover:text-primary transition-colors"
                    >
                      {project.name}
                    </Link>
                  </CardTitle>
                </CardHeader>
                {project.description && (
                  <CardContent>
                    <p className="text-sm text-muted-foreground">{project.description}</p>
                  </CardContent>
                )}
              </Card>
            ))
          ) : (
            <p className="text-center text-muted-foreground py-8">暂无收藏的项目</p>
          )}
        </TabsContent>
      </Tabs>
    </PageLayout>
  )
}
