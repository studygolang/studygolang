import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

interface User {
  uid: number
  username: string
  name?: string
  avatar?: string
}

interface Article {
  id: number
  title: string
  author?: string
  ctime?: string
  viewnum?: number
}

interface Topic {
  tid: number
  title: string
  author?: string
  ctime?: string
  view?: number
}

interface Resource {
  id: number
  title: string
  author?: string
  ctime?: string
}

interface Project {
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
  articles?: Article[]
  topics?: Topic[]
  resources?: Resource[]
  projects?: Project[]
}

async function fetchFromAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const base = process.env.API_BASE_URL || 'http://localhost:8090'
  const res = await fetch(`${base}/api/v1${path}`, options)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const json = await res.json()
  return json.data
}

async function getFavorites(username: string, objtype: number = 1): Promise<FavoriteData | null> {
  try {
    return await fetchFromAPI<FavoriteData>(
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
  const objtype = parseInt(objtypeStr || '1', 10)

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

      <Tabs defaultValue={objtype.toString()} className="space-y-4">
        <TabsList>
          <TabsTrigger value="1" asChild>
            <Link href={`/favorites/${username}?objtype=1`}>文章</Link>
          </TabsTrigger>
          <TabsTrigger value="2" asChild>
            <Link href={`/favorites/${username}?objtype=2`}>话题</Link>
          </TabsTrigger>
          <TabsTrigger value="3" asChild>
            <Link href={`/favorites/${username}?objtype=3`}>资源</Link>
          </TabsTrigger>
          <TabsTrigger value="4" asChild>
            <Link href={`/favorites/${username}?objtype=4`}>项目</Link>
          </TabsTrigger>
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
                      href={`/projects/${project.id}`}
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
