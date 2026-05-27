"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Star, GitFork, Eye, ExternalLink, Github } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { fetchAPINullable } from "@/lib/api"
import { formatNum } from "@/lib/utils"
import { Pagination } from "@/components/pagination"
import { ProjectSortSelect } from "@/components/project-sort-select"
import type { ProjectListData, Project } from "@/lib/types"

interface ProjectsClientProps {
  page: number
  sort: string
  total: number
  projects?: Project[]
  initialHasMore?: boolean
}

export function ProjectsClient({ page, sort, total: initialTotal, projects: initialProjects = [], initialHasMore: initialHasMore = false }: ProjectsClientProps) {
  const [projects, setProjects] = useState<Project[]>(initialProjects)
  const [hasMore, setHasMore] = useState(initialHasMore)
  const [loading, setLoading] = useState(initialProjects.length === 0)
  const [count, setCount] = useState(initialTotal)

  useEffect(() => {
    // 如果没有初始数据（SSR 未提供），则客户端请求
    if (initialProjects.length === 0) {
      setLoading(true)
      fetchAPINullable<ProjectListData>(
        `/projects?p=${page}&sort=${sort}`,
        { cache: "no-store" }
      ).then((data) => {
        if (data && data.projects) {
          setProjects(data.projects)
          setHasMore(data.has_more)
          if (data.total) {
            setCount(data.total)
          }
        }
        setLoading(false)
      })
    }
  }, [page, sort, initialProjects.length])

  return (
    <div className="rounded-lg border border-border bg-card p-4">
      {/* Sort Controls */}
      <div className="mb-4 flex items-center justify-between">
        <span className="text-sm text-muted-foreground">共 {count} 个项目</span>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">排序:</span>
          <ProjectSortSelect currentSort={sort} />
        </div>
      </div>

      {/* Loading State */}
      {loading && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="h-48 animate-pulse rounded-lg bg-muted" />
          ))}
        </div>
      )}

      {/* Empty State */}
      {!loading && projects.length === 0 && (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <Github className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">暂无项目数据</p>
        </div>
      )}

      {/* Project Grid */}
      {!loading && projects.length > 0 && (
        <div>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {projects.map((project) => (
              <Card
                key={project.id}
                className="group transition-all hover:border-primary/20 hover:shadow-sm"
              >
                <CardContent className="flex h-full flex-col p-5">
                  {/* Header */}
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                      {project.logo ? (
                        <img
                          src={project.logo}
                          alt={project.name}
                          className="h-9 w-9 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10">
                          <Github className="h-4.5 w-4.5 text-primary" />
                        </div>
                      )}
                      <div>
                        <Link href={`/p/${project.uri}`}>
                          <h3 className="text-sm font-bold text-foreground group-hover:text-primary">
                            {project.name}
                          </h3>
                        </Link>
                        <p className="text-xs text-muted-foreground">{project.author}</p>
                      </div>
                    </div>
                    {project.src_url && (
                      <a
                        href={project.src_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-secondary hover:text-primary"
                        aria-label={"查看 " + project.name + " 源码"}
                      >
                        <ExternalLink className="h-4 w-4" />
                      </a>
                    )}
                  </div>

                  {/* Description */}
                  <p className="mt-3 flex-1 line-clamp-3 text-sm leading-relaxed text-muted-foreground">
                    {project.desc || "暂无描述"}
                  </p>

                  {/* Category */}
                  {project.category && (
                    <div className="mt-3 flex flex-wrap gap-1.5">
                      <Badge variant="secondary" className="text-[10px] font-medium">
                        {project.category}
                      </Badge>
                      {project.lang && (
                        <Badge variant="outline" className="text-[10px] font-medium">
                          {project.lang}
                        </Badge>
                      )}
                    </div>
                  )}

                  {/* Stats */}
                  <div className="mt-3 flex items-center gap-4 border-t border-border pt-3">
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Star className="h-3.5 w-3.5" />
                      {formatNum(project.star)}
                    </div>
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <GitFork className="h-3.5 w-3.5" />
                      {formatNum(project.fork)}
                    </div>
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Eye className="h-3.5 w-3.5" />
                      {formatNum(project.watch)}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Pagination */}
          <Pagination
            currentPage={page}
            hasMore={hasMore}
            buildUrl={(p) => `/projects?p=${p}&sort=${sort}`}
          />
        </div>
      )}
    </div>
  )
}
