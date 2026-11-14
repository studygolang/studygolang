"use client"

import Link from "next/link"
import { Clock, Eye, ExternalLink, Timer } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { MarkdownContent } from "@/components/markdown-content"
import { SafeHtml } from "@/components/safe-html"
import { CommentSection } from "@/components/comment-section"
import { ContentActions } from "@/components/content-actions"
import { formatTime } from "@/lib/utils"
import { OBJTYPE_ARTICLE, type Article, type ArticleComment, type ArticleGCTT, type Subject } from "@/lib/types"

interface ArticleDetailProps {
  id: string
  article?: Article
  prev?: Article
  next?: Article
  // 后端文章详情评论字段名为 replies，类型为 ArticleComment
  comments?: ArticleComment[]
  gctt?: ArticleGCTT | null
  subjects?: Subject[] | null
  readTime?: number
  likeflag?: number
  hadcollect?: number
}

export function ArticleDetail({
  id,
  article,
  prev,
  next,
  comments = [],
  gctt,
  subjects,
  readTime,
  likeflag,
  hadcollect,
}: ArticleDetailProps) {
  if (!article) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-muted-foreground">
          文章不存在或已被删除
        </CardContent>
      </Card>
    )
  }

  const aid = parseInt(id)
  const tags = (article.tags ?? "").split(",").map((t) => t.trim()).filter(Boolean)
  // 与旧版一致：超过 100 天的文章提示时效性
  const outdated =
    !!article.ctime && (Date.now() - new Date(article.ctime).getTime()) / 86400000 > 100
  const translator = article.gctt ? gctt?.Translator : undefined

  const relatedArticles = [
    ...(prev ? [{ title: prev.title, href: `/articles/${prev.id}` }] : []),
    ...(next ? [{ title: next.title, href: `/articles/${next.id}` }] : []),
  ]

  return (
    <div className="space-y-4">
      {/* Article Card */}
      <Card>
        <CardContent className="p-5 sm:p-8">
          {/* Title */}
          <h1 className="text-balance text-2xl font-bold leading-tight text-foreground sm:text-[28px]">
            {article.title}
          </h1>

          {/* Meta */}
          <div className="mt-4 flex flex-wrap items-center gap-3">
            <div className="flex items-center gap-2">
              <Avatar className="h-8 w-8">
                {article.gctt && gctt?.Avatar && (
                  <AvatarImage src={gctt.Avatar} alt={gctt.Translator} />
                )}
                <AvatarFallback className="bg-primary/10 text-sm font-semibold text-primary">
                  {(translator || article.author || "?").charAt(0).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div>
                {/* GCTT 译文署名译者；自主发布署站内账号；转载只显示原文作者名 */}
                {translator ? (
                  <Link
                    href={`/gctt/${translator}`}
                    className="text-sm font-medium text-foreground hover:text-primary"
                  >
                    {translator}
                  </Link>
                ) : article.is_self && article.author_txt ? (
                  <Link
                    href={`/user/${article.author_txt}`}
                    className="text-sm font-medium text-foreground hover:text-primary"
                  >
                    {article.author_txt}
                  </Link>
                ) : (
                  <span className="text-sm font-medium text-foreground">
                    {article.author_txt || article.author || "未知作者"}
                  </span>
                )}
                <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Clock className="h-3 w-3" />
                    {formatTime(article.ctime)}
                  </span>
                  <span className="flex items-center gap-1">
                    <Eye className="h-3 w-3" />
                    {article.viewnum} 次点击
                  </span>
                  {readTime !== undefined && (
                    <span className="flex items-center gap-1">
                      <Timer className="h-3 w-3" />
                      {readTime >= 1 ? `预计阅读 ${readTime} 分钟` : "预计阅读不到 1 分钟"}
                    </span>
                  )}
                </div>
              </div>
            </div>
            {tags.length > 0 && (
              <div className="ml-auto flex flex-wrap items-center gap-1.5">
                {tags.map((tag) => (
                  <Link key={tag} href={`/tag/${encodeURIComponent(tag)}`}>
                    <Badge variant="secondary" className="bg-primary/10 text-xs text-primary">
                      {tag}
                    </Badge>
                  </Link>
                ))}
              </div>
            )}
          </div>

          {outdated && (
            <p className="mt-4 rounded-md border border-border bg-muted/40 px-3 py-2 text-xs text-muted-foreground">
              这是一篇创建于 {formatTime(article.ctime)} 的文章，其中的信息可能已经有所发展或是发生改变。
            </p>
          )}

          <Separator className="my-6" />

          {/* Content：markdown 走渲染，富文本直接安全输出 HTML（旧版双轨） */}
          {article.markdown ? (
            <MarkdownContent content={article.content ?? ""} />
          ) : (
            <SafeHtml html={article.content ?? ""} className="prose prose-sm max-w-none text-foreground" />
          )}

          {/* GCTT 译文出处与署名 */}
          {article.gctt && gctt && (
            <div className="mt-6 space-y-2 border-t border-border pt-4 text-xs text-muted-foreground">
              {gctt.URL && (
                <p>
                  via:{" "}
                  <a href={gctt.URL} target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                    {gctt.URL}
                  </a>
                </p>
              )}
              <p className="flex flex-wrap gap-x-3 gap-y-1">
                {gctt.Author && (
                  <span>
                    作者：
                    <a href={gctt.AuthorURL || undefined} target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                      {gctt.Author}
                    </a>
                  </span>
                )}
                {gctt.Translator && (
                  <span>
                    译者：
                    <a href={`https://github.com/${gctt.Translator}`} target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                      {gctt.Translator}
                    </a>
                  </span>
                )}
                {!!gctt.Checkers?.length && (
                  <span>
                    校对：
                    {gctt.Checkers.map((checker) => (
                      <a
                        key={checker}
                        href={`https://github.com/${checker}`}
                        target="_blank"
                        rel="external nofollow noopener"
                        className="mr-1 text-primary hover:underline"
                      >
                        {checker}
                      </a>
                    ))}
                  </span>
                )}
              </p>
              <p>
                本文由{" "}
                <a href="https://github.com/studygolang/GCTT" target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                  GCTT
                </a>{" "}
                原创编译，
                <Link href={`/articles/${aid}`} className="text-primary hover:underline">
                  Go语言中文网
                </Link>{" "}
                荣誉推出。翻译工作遵照 CC-BY-NC-SA 协议规定，转载敬请保留原文/译文链接和作者/译者信息。
              </p>
            </div>
          )}

          {/* 转载来源 */}
          {!article.is_self && (article.domain || article.url) && (
            <div className="mt-6 space-y-1 border-t border-border pt-4 text-xs text-muted-foreground">
              {article.domain && (
                <p>
                  本文来自：
                  <a href={`https://${article.domain}`} target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                    {article.name || article.domain}
                  </a>
                </p>
              )}
              {article.author_txt && <p>感谢作者：{article.author_txt}</p>}
              {article.url && (
                <p>
                  查看原文：
                  <a href={article.url} target="_blank" rel="external nofollow noopener" className="text-primary hover:underline">
                    {article.title}
                  </a>
                </p>
              )}
            </div>
          )}

          {/* 所属专栏 */}
          {!!subjects?.length && (
            <div className="mt-6 flex flex-wrap items-center gap-2 border-t border-border pt-4">
              <span className="text-xs text-muted-foreground">所属专栏：</span>
              {subjects.map((subject) => (
                <Link key={subject.id} href={`/subject/${subject.id}`}>
                  <Badge variant="outline" className="text-xs">
                    {subject.name}
                  </Badge>
                </Link>
              ))}
            </div>
          )}

          {/* Actions */}
          <Separator className="my-6" />
          <ContentActions
            objid={aid}
            objtype={OBJTYPE_ARTICLE}
            likeCount={article.likenum || 0}
            initialLiked={Boolean(likeflag)}
            initialFavorited={Boolean(hadcollect)}
            ownerUsername={article.author}
            editHref={`/articles/modify/${aid}`}
          />
        </CardContent>
      </Card>

      {/* Related Articles */}
      {relatedArticles.length > 0 && (
        <Card>
          <CardContent className="p-5">
            <h3 className="mb-3 text-sm font-semibold text-foreground">{"相关推荐"}</h3>
            <div className="space-y-2">
              {relatedArticles.map((a, i) => (
                <Link
                  key={i}
                  href={a.href}
                  className="flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
                >
                  <ExternalLink className="h-3 w-3 shrink-0" />
                  {a.title}
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      <CommentSection
        objid={aid}
        objtype={OBJTYPE_ARTICLE}
        comments={comments}
        closed={article.close_reply}
      />
    </div>
  )
}
