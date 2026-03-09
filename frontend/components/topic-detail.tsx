"use client"

import { useState } from "react"
import Link from "next/link"
import {
  ThumbsUp,
  MessageSquare,
  Bookmark,
  Share2,
  Eye,
  Clock,
  Flag,
  ChevronUp,
  MoreHorizontal,
} from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"

const topicData = {
  title: "Go 1.26 蓄势待发，第一个 RC 版本已发布",
  author: "polaris",
  authorInitial: "P",
  authorLevel: "管理员",
  tag: "Go动态",
  tagColor: "bg-primary/10 text-primary",
  createdAt: "2026-02-25 14:30",
  views: 1850,
  likes: 128,
  bookmarks: 45,
  content: `Go 1.26 将引入多项重要更新，包括泛型增强、性能优化和新的标准库功能。本文详细解析了 RC1 版本的主要变化。

## 主要更新

### 1. 泛型增强

Go 1.26 在泛型方面带来了几个重要的改进：

- **类型约束推断优化**：编译器现在能更好地推断泛型函数的类型参数
- **联合类型约束**：支持更灵活的类型约束组合
- **方法集扩展**：泛型类型的方法集现在可以包含更多操作

\`\`\`go
func Map[S ~[]E, E any, R any](s S, fn func(E) R) []R {
    result := make([]R, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}
\`\`\`

### 2. 性能优化

- 垃圾回收器延迟降低 15%
- 编译速度提升 20%
- 运行时调度器优化，更好地利用多核处理器

### 3. 新标准库功能

- \`maps\` 包新增 \`Collect\` 和 \`Keys\` 函数
- \`slices\` 包新增 \`Chunk\` 和 \`Repeat\` 函数
- \`log/slog\` 包性能优化和新的 Handler 接口

## 如何体验

你可以通过以下命令安装 RC 版本：

\`\`\`bash
go install golang.org/dl/go1.26rc1@latest
go1.26rc1 download
\`\`\`

欢迎大家试用并反馈问题！`,
}

const comments = [
  {
    id: 1,
    author: "alice_go",
    initial: "A",
    content:
      "泛型增强太棒了！之前写通用函数总是要用 interface{}，现在终于可以更优雅了。期待正式版发布！",
    time: "2 小时前",
    likes: 23,
    floor: 1,
  },
  {
    id: 2,
    author: "bob_dev",
    initial: "B",
    content:
      "GC 延迟降低 15% 对我们线上服务帮助很大，之前 P99 延迟偶尔有毛刺就是 GC 导致的。会尽快升级测试。",
    time: "1 小时前",
    likes: 15,
    floor: 2,
  },
  {
    id: 3,
    author: "charlie",
    initial: "C",
    content:
      '请问 slices.Chunk 函数的具体用法是什么？能给个例子吗？\n\n我看文档里写的是 `func Chunk[S ~[]E, E any](s S, n int) iter.Seq[S]`，返回的是 iterator，感觉比之前方便很多。',
    time: "45 分钟前",
    likes: 8,
    floor: 3,
  },
  {
    id: 4,
    author: "polaris",
    initial: "P",
    content:
      "回复 @charlie：\n\n用法很简单：\n```go\nfor chunk := range slices.Chunk(data, 100) {\n    process(chunk)\n}\n```\n\n每次迭代会得到一个最大长度为 100 的切片。非常适合批量处理场景。",
    time: "30 分钟前",
    likes: 19,
    floor: 4,
    isAuthor: true,
  },
]

export function TopicDetail({ id }: { id: string }) {
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)

  return (
    <div className="space-y-4">
      {/* Topic Card */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          {/* Title */}
          <h1 className="text-balance text-xl font-bold leading-snug text-foreground sm:text-2xl">
            {topicData.title}
          </h1>

          {/* Meta */}
          <div className="mt-3 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-7 w-7">
                <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                  {topicData.authorInitial}
                </AvatarFallback>
              </Avatar>
              <Link
                href={`/user/${topicData.author}`}
                className="text-sm font-medium text-foreground hover:text-primary"
              >
                {topicData.author}
              </Link>
              <Badge
                variant="secondary"
                className="bg-primary/10 text-[10px] font-semibold text-primary"
              >
                {topicData.authorLevel}
              </Badge>
            </div>
            <Separator orientation="vertical" className="h-4" />
            <Badge variant="secondary" className={`text-xs ${topicData.tagColor}`}>
              {topicData.tag}
            </Badge>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Clock className="h-3 w-3" />
              {topicData.createdAt}
            </span>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Eye className="h-3 w-3" />
              {topicData.views}
            </span>
          </div>

          {/* Content */}
          <div className="prose prose-sm mt-6 max-w-none text-foreground">
            {topicData.content.split("\n").map((line, i) => {
              if (line.startsWith("## ")) {
                return (
                  <h2 key={i} className="mb-3 mt-6 text-lg font-bold text-foreground">
                    {line.replace("## ", "")}
                  </h2>
                )
              }
              if (line.startsWith("### ")) {
                return (
                  <h3 key={i} className="mb-2 mt-4 text-base font-semibold text-foreground">
                    {line.replace("### ", "")}
                  </h3>
                )
              }
              if (line.startsWith("```")) {
                return null
              }
              if (line.startsWith("- ")) {
                return (
                  <div key={i} className="flex gap-2 py-0.5 text-sm leading-relaxed text-foreground">
                    <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                    <span>{line.replace("- ", "")}</span>
                  </div>
                )
              }
              if (line.trim() === "") return <div key={i} className="h-2" />
              if (
                line.includes("func ") ||
                line.includes("go install") ||
                line.includes("result") ||
                line.includes("for ") ||
                line.includes("go1.26") ||
                line.includes("process(")
              ) {
                return (
                  <pre
                    key={i}
                    className="my-1 overflow-x-auto rounded bg-secondary/60 px-3 py-1 font-mono text-xs leading-relaxed text-foreground"
                  >
                    <code>{line}</code>
                  </pre>
                )
              }
              return (
                <p key={i} className="text-sm leading-relaxed text-foreground">
                  {line}
                </p>
              )
            })}
          </div>

          {/* Actions */}
          <Separator className="my-5" />
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant={liked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                liked
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground"
              }`}
              onClick={() => setLiked(!liked)}
            >
              <ThumbsUp className="h-3.5 w-3.5" />
              {liked ? topicData.likes + 1 : topicData.likes}
            </Button>
            <Button
              variant={bookmarked ? "default" : "outline"}
              size="sm"
              className={`gap-1.5 ${
                bookmarked
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground"
              }`}
              onClick={() => setBookmarked(!bookmarked)}
            >
              <Bookmark className="h-3.5 w-3.5" />
              {bookmarked ? "已收藏" : "收藏"}
            </Button>
            <Button variant="outline" size="sm" className="gap-1.5 text-muted-foreground">
              <Share2 className="h-3.5 w-3.5" />
              {"分享"}
            </Button>
            <Button variant="ghost" size="sm" className="ml-auto gap-1.5 text-muted-foreground">
              <Flag className="h-3.5 w-3.5" />
              {"举报"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Comments Section */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          <h2 className="flex items-center gap-2 text-base font-semibold text-foreground">
            <MessageSquare className="h-4 w-4 text-primary" />
            {"评论"} ({comments.length})
          </h2>

          {/* Comment Input */}
          <div className="mt-4 rounded-lg border border-border p-3">
            <textarea
              placeholder="写下你的评论..."
              className="min-h-[80px] w-full resize-none bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
            />
            <div className="mt-2 flex items-center justify-between">
              <p className="text-xs text-muted-foreground">
                {"支持 Markdown 语法"}
              </p>
              <Button size="sm" className="bg-primary text-primary-foreground hover:bg-primary/90">
                {"发表评论"}
              </Button>
            </div>
          </div>

          {/* Comments List */}
          <div className="mt-6 space-y-0 divide-y divide-border">
            {comments.map((comment) => (
              <div key={comment.id} className="py-4 first:pt-0">
                <div className="flex gap-3">
                  <Avatar className="mt-0.5 h-8 w-8 shrink-0">
                    <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                      {comment.initial}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <Link
                        href={`/user/${comment.author}`}
                        className="text-sm font-semibold text-foreground hover:text-primary"
                      >
                        {comment.author}
                      </Link>
                      {comment.isAuthor && (
                        <Badge variant="secondary" className="bg-primary/10 text-[10px] text-primary">
                          {"楼主"}
                        </Badge>
                      )}
                      <span className="text-xs text-muted-foreground">
                        {comment.time}
                      </span>
                      <span className="ml-auto text-xs text-muted-foreground">
                        #{comment.floor}
                      </span>
                    </div>
                    <div className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                      {comment.content}
                    </div>
                    <div className="mt-2 flex items-center gap-3">
                      <button className="flex items-center gap-1 text-xs text-muted-foreground transition-colors hover:text-primary">
                        <ChevronUp className="h-3.5 w-3.5" />
                        {comment.likes}
                      </button>
                      <button className="text-xs text-muted-foreground transition-colors hover:text-primary">
                        {"回复"}
                      </button>
                      <button className="ml-auto text-muted-foreground transition-colors hover:text-primary">
                        <MoreHorizontal className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
