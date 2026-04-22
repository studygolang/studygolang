import {
  MessageSquare,
  FileText,
  FolderGit2,
  Link2,
  BookOpen,
} from "lucide-react"
import { cn } from "@/lib/utils"

type ContentType = "topic" | "article" | "project" | "resource" | "book"

const CONTENT_TYPES: { type: ContentType; icon: React.ReactNode; label: string; desc: string }[] = [
  { type: "topic", icon: <MessageSquare className="h-5 w-5" />, label: "主题", desc: "发起一个讨论话题" },
  { type: "article", icon: <FileText className="h-5 w-5" />, label: "文章", desc: "发表原创或翻译文章" },
  { type: "project", icon: <FolderGit2 className="h-5 w-5" />, label: "项目", desc: "分享开源项目" },
  { type: "resource", icon: <Link2 className="h-5 w-5" />, label: "资源", desc: "分享学习资源" },
  { type: "book", icon: <BookOpen className="h-5 w-5" />, label: "图书", desc: "推荐 Go 相关图书" },
]

interface PublishHeaderProps {
  contentType: ContentType
  onContentTypeChange: (type: ContentType) => void
}

export function PublishHeader({ contentType, onContentTypeChange }: PublishHeaderProps) {
  return (
    <div className="mb-6 flex items-center justify-between">
      <div>
        <h1 className="text-xl font-semibold text-foreground">发布内容</h1>
        <p className="mt-1 text-sm text-muted-foreground">与社区分享你的知识和发现</p>
      </div>

      {/* 内容类型选择 - 紧凑的标签页样式 */}
      <div className="flex items-center gap-1 rounded-lg border border-border bg-muted/50 p-1">
        {CONTENT_TYPES.map(({ type, icon, label }) => (
          <button
            key={type}
            type="button"
            onClick={() => onContentTypeChange(type)}
            className={cn(
              "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-all",
              contentType === type
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <span className="h-4 w-4">{icon}</span>
            {label}
          </button>
        ))}
      </div>
    </div>
  )
}
