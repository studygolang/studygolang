import Link from "next/link"
import {
  Eye,
  ExternalLink,
  Globe,
  BookOpen,
  Github,
  MessageSquare,
  ThumbsUp,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { formatNum } from "@/lib/utils"
import type { Project } from "@/lib/types"

interface ProjectListProps {
  projects?: Project[]
}

export function ProjectList({ projects = [] }: ProjectListProps) {
  if (projects.length === 0) {
    return (
      <div className="py-8 text-center text-sm text-muted-foreground">
        暂无项目
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {projects.map((project) => (
        <Card
          key={project.id}
          className="group transition-all hover:border-primary/20 hover:shadow-sm"
        >
          <CardContent className="p-4 sm:p-5">
            <div className="flex gap-4">
              {/* Logo */}
              {project.logo ? (
                <img
                  src={project.logo}
                  alt={project.name}
                  className="h-16 w-16 shrink-0 rounded-lg object-cover ring-1 ring-border"
                />
              ) : (
                <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-lg bg-primary/10 ring-1 ring-border">
                  <Globe className="h-8 w-8 text-primary" />
                </div>
              )}

              {/* Content */}
              <div className="min-w-0 flex-1 space-y-2">
                <div className="flex items-start gap-2">
                  <Link
                    href={`/p/${project.uri || project.id}`}
                    className="text-base font-semibold leading-snug text-foreground transition-colors group-hover:text-primary"
                  >
                    {project.name}
                  </Link>
                  {project.lang && (
                    <Badge variant="secondary" className="text-[10px] font-medium">
                      {project.lang}
                    </Badge>
                  )}
                  {project.category && (
                    <Badge className="bg-primary/10 text-[10px] font-medium text-primary">
                      {project.category}
                    </Badge>
                  )}
                </div>

                {/* Description */}
                {project.desc && (
                  <p className="line-clamp-2 text-sm leading-relaxed text-muted-foreground">
                    {project.desc}
                  </p>
                )}

                {/* Meta stats */}
                <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5">
                  {project.home && (
                    <a
                      href={project.home}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
                    >
                      <Globe className="h-3 w-3" />
                      官网
                    </a>
                  )}
                  {project.doc && (
                    <a
                      href={project.doc}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
                    >
                      <BookOpen className="h-3 w-3" />
                      文档
                    </a>
                  )}
                  {project.src && (
                    <a
                      href={project.src}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1 text-xs text-primary transition-colors hover:text-primary/80"
                    >
                      <Github className="h-3 w-3" />
                      源码
                    </a>
                  )}
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Eye className="h-3 w-3" />
                    {formatNum(project.viewnum)}
                  </span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <MessageSquare className="h-3 w-3" />
                    {formatNum(project.cmtnum)}
                  </span>
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <ThumbsUp className="h-3 w-3" />
                    {formatNum(project.likenum)}
                  </span>
                </div>
              </div>

              {/* External link */}
              {project.uri && (
                <a
                  href={project.uri}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="mt-1 flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-border text-muted-foreground transition-colors hover:border-primary/30 hover:bg-primary/5 hover:text-primary"
                  aria-label={"访问 " + project.name}
                >
                  <ExternalLink className="h-4 w-4" />
                </a>
              )}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
