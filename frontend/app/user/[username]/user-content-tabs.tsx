"use client"

import { useState } from "react"
import { MessageSquare, FileText, MessageCircle, Library, FolderGit2 } from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import type { Topic, Article, Resource, Project, UserComment } from "@/lib/types"
import { userAPI } from "@/lib/api"
import { encodeProjectUri } from "@/lib/utils"

interface UserContentTabsProps {
  username: string
  initialTopics: Topic[]
  initialArticles: Article[]
  initialResources: Resource[]
  initialProjects: Project[]
}

const TABS = [
  { key: "topics", label: "话题", icon: MessageSquare },
  { key: "articles", label: "文章", icon: FileText },
  { key: "resources", label: "资源", icon: Library },
  { key: "projects", label: "项目", icon: FolderGit2 },
  { key: "comments", label: "评论", icon: MessageCircle },
] as const

type TabKey = (typeof TABS)[number]["key"]

export function UserContentTabs({
  username,
  initialTopics,
  initialArticles,
  initialResources,
  initialProjects,
}: UserContentTabsProps) {
  const [activeTab, setActiveTab] = useState<TabKey>("topics")
  const [comments, setComments] = useState<UserComment[]>([])
  const [loading, setLoading] = useState(false)

  const handleTabChange = async (tab: TabKey) => {
    setActiveTab(tab)
    if (tab === "comments" && comments.length === 0) {
      setLoading(true)
      try {
        const res = await userAPI.getUserComments(username, { p: 1 })
        setComments(res?.comments ?? [])
      } catch {
        // 忽略错误，显示空列表
      } finally {
        setLoading(false)
      }
    }
  }

  return (
    <div>
      {/* Tab header */}
      <div className="mb-4 flex flex-wrap gap-1 rounded-lg border border-border bg-muted/50 p-1">
        {TABS.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            onClick={() => handleTabChange(key)}
            className={`flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
              activeTab === key
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Icon className="h-3.5 w-3.5" />
            {label}
          </button>
        ))}
      </div>

      {/* Tab content */}
      {activeTab === "topics" && (
        <div className="space-y-3">
          {initialTopics.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
              <MessageSquare className="mb-2 h-8 w-8 text-muted-foreground/30" />
              <p className="text-sm text-muted-foreground">暂无话题</p>
            </div>
          ) : (
            initialTopics.map((topic) => (
              <Card key={topic.tid} className="transition-colors hover:border-primary/20">
                <CardContent className="p-4">
                  <a
                    href={`/topics/${topic.tid}`}
                    className="text-sm font-medium text-foreground transition-colors hover:text-primary"
                  >
                    {topic.title}
                  </a>
                  <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                    {topic.reply != null && <span>{topic.reply} 回复</span>}
                    {topic.view != null && <span>{topic.view} 浏览</span>}
                    <span>{topic.ctime}</span>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {activeTab === "articles" && (
        <div className="space-y-3">
          {initialArticles.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
              <FileText className="mb-2 h-8 w-8 text-muted-foreground/30" />
              <p className="text-sm text-muted-foreground">暂无文章</p>
            </div>
          ) : (
            initialArticles.map((article) => (
              <Card key={article.id} className="transition-colors hover:border-primary/20">
                <CardContent className="p-4">
                  <a
                    href={`/articles/${article.id}`}
                    className="text-sm font-medium text-foreground transition-colors hover:text-primary"
                  >
                    {article.title}
                  </a>
                  <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                    {article.viewnum != null && <span>{article.viewnum} 浏览</span>}
                    <span>{article.ctime}</span>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {activeTab === "resources" && (
        <div className="space-y-3">
          {initialResources.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
              <Library className="mb-2 h-8 w-8 text-muted-foreground/30" />
              <p className="text-sm text-muted-foreground">暂无资源</p>
            </div>
          ) : (
            initialResources.map((resource) => (
              <Card key={resource.id} className="transition-colors hover:border-primary/20">
                <CardContent className="p-4">
                  <a
                    href={`/resources/${resource.id}`}
                    className="text-sm font-medium text-foreground transition-colors hover:text-primary"
                  >
                    {resource.title}
                  </a>
                  <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                    {resource.catname && <span>{resource.catname}</span>}
                    <span>{resource.viewnum} 浏览</span>
                    <span>{resource.ctime}</span>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {activeTab === "projects" && (
        <div className="space-y-3">
          {initialProjects.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
              <FolderGit2 className="mb-2 h-8 w-8 text-muted-foreground/30" />
              <p className="text-sm text-muted-foreground">暂无项目</p>
            </div>
          ) : (
            initialProjects.map((project) => (
              <Card key={project.id} className="transition-colors hover:border-primary/20">
                <CardContent className="p-4">
                  <a
                    href={`/p/${encodeProjectUri(project.uri || project.id)}`}
                    className="text-sm font-medium text-foreground transition-colors hover:text-primary"
                  >
                    {project.name}
                  </a>
                  <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                    {project.category && <span>{project.category}</span>}
                    <span>{project.viewnum} 浏览</span>
                    <span>{project.ctime}</span>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}

      {activeTab === "comments" && (
        <div className="space-y-3">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            </div>
          ) : comments.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
              <MessageCircle className="mb-2 h-8 w-8 text-muted-foreground/30" />
              <p className="text-sm text-muted-foreground">暂无评论</p>
            </div>
          ) : (
            comments.map((comment) => (
              <Card key={comment.cid} className="transition-colors hover:border-primary/20">
                <CardContent className="p-4">
                  <p className="text-sm text-foreground line-clamp-2">{comment.content}</p>
                  <div className="mt-1 flex items-center gap-3 text-xs text-muted-foreground">
                    <span>{comment.ctime}</span>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      )}
    </div>
  )
}
