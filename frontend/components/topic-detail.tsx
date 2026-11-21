"use client"

import Link from "next/link"
import { Clock, Eye } from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { TopicAppendForm } from "@/components/topic-append-form"
import { CommentSection } from "@/components/comment-section"
import { ContentActions } from "@/components/content-actions"
import { MarkdownContent } from "@/components/markdown-content"
import { useAuth } from "@/lib/auth-context"
import { formatTime } from "@/lib/utils"
import { OBJTYPE_TOPIC, type Topic, type TopicReply, type TopicAppend } from "@/lib/types"

interface TopicDetailProps {
  id: string
  topic?: Topic
  replies?: TopicReply[]
  appends?: TopicAppend[]
  canView?: boolean
  likeflag?: number
  hadcollect?: number
}

export function TopicDetail({ id, topic, replies = [], appends = [], canView = true, likeflag, hadcollect }: TopicDetailProps) {
  if (!topic) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-muted-foreground">
          话题不存在或已被删除
        </CardContent>
      </Card>
    )
  }

  const tid = parseInt(id)

  return (
    <div className="space-y-4">
      {/* Topic Card */}
      <Card>
        <CardContent className="p-5 sm:p-6">
          {/* Title */}
          <h1 className="text-balance text-xl font-bold leading-snug text-foreground sm:text-2xl">
            {topic.title}
          </h1>

          {/* Meta */}
          <div className="mt-3 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-7 w-7">
                <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                  {topic.user?.username ? topic.user.username.charAt(0).toUpperCase() : (topic.name ? topic.name.charAt(0).toUpperCase() : "?")}
                </AvatarFallback>
              </Avatar>
              <Link
                href={`/user/${topic.user?.username || topic.name}`}
                className="text-sm font-medium text-foreground hover:text-primary"
              >
                {topic.user?.username || topic.name}
              </Link>
            </div>
            {topic.node && (
              <>
                <Separator orientation="vertical" className="h-4" />
                <Badge variant="secondary" className="bg-primary/10 text-xs text-primary">
                  {topic.node.name}
                </Badge>
              </>
            )}
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Clock className="h-3 w-3" />
              {formatTime(topic.ctime)}
            </span>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Eye className="h-3 w-3" />
              {topic.view || topic.viewnum || 0}
            </span>
          </div>

          {/* Content */}
          <MarkdownContent content={topic.content ?? ""} className="mt-6" />

          {/* Actions */}
          <Separator className="my-5" />
          <ContentActions
            objid={tid}
            objtype={OBJTYPE_TOPIC}
            likeCount={topic.like || topic.likenum || 0}
            initialLiked={Boolean(likeflag)}
            initialFavorited={Boolean(hadcollect)}
            ownerUid={topic.uid}
            editHref={`/topics/modify/${tid}`}
            showTop
            isTop={topic.top === 1}
          />
        </CardContent>
      </Card>

      {/* 追加内容展示 */}
      {appends.length > 0 && (
        <Card>
          <CardContent className="p-5 sm:p-6">
            <h2 className="text-base font-semibold text-foreground">附言</h2>
            <div className="mt-4 space-y-4 divide-y divide-border">
              {appends.map((append, index) => (
                <div key={append.id} className={index > 0 ? "pt-4" : ""}>
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <span className="font-medium">附言 {index + 1}</span>
                    <span>{formatTime(append.created_at)}</span>
                  </div>
                  <MarkdownContent content={append.content} className="mt-2" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 追加表单（仅作者可见，客户端判断） */}
      <TopicAppendSection tid={tid} topicUid={topic.uid} />

      {canView && (
        <CommentSection
          objid={tid}
          objtype={OBJTYPE_TOPIC}
          comments={replies}
          ownerUid={topic.uid}
          closed={topic.close_reply}
        />
      )}
    </div>
  )
}

// 追加表单区域：客户端判断当前用户是否为作者
function TopicAppendSection({ tid, topicUid }: { tid: number; topicUid?: number }) {
  const { user } = useAuth()

  if (topicUid === undefined || !user || user.uid !== topicUid) return null
  return <TopicAppendForm tid={tid} />
}
