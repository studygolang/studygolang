"use client"

import React, { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { useAuth } from "@/lib/auth-context"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Trash2, Mail, Send, Bell, Loader2 } from "lucide-react"
import Link from "next/link"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"

const API_BASE = ""

type MessageType = "system" | "inbox" | "outbox"

interface Message {
  id: number
  content: string
  hasread: string
  ctime: string
  user?: {
    uid: number
    username: string
    avatar: string
  }
  title?: string
  objtitle?: string
  objurl?: string
}

interface MessageListData {
  messages: Message[]
  page: number
  total: number
  has_more: boolean
}

async function fetchMessages(
  type: MessageType,
  page: number
): Promise<MessageListData> {
  try {
    const res = await fetch(
      `${API_BASE}/api/v1/messages?type=${type}&p=${page}`,
      { credentials: "include" }
    )
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const json = await res.json()
    if (json.code !== 0) throw new Error(json.msg)
    return json.data
  } catch {
    return { messages: [], page, total: 0, has_more: false }
  }
}

async function deleteMessage(id: number, type: MessageType): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/messages/${id}?type=${type}`, {
      method: "DELETE",
      credentials: "include",
    })
    const json = await res.json()
    return json.code === 0
  } catch {
    return false
  }
}

const TAB_LABELS: Record<MessageType, { label: string; icon: React.ReactNode }> = {
  system: { label: "系统消息", icon: <Bell className="h-4 w-4" /> },
  inbox: { label: "收件箱", icon: <Mail className="h-4 w-4" /> },
  outbox: { label: "发件箱", icon: <Send className="h-4 w-4" /> },
}

function MessageItem({
  msg,
  msgtype,
  onDelete,
}: {
  msg: Message
  msgtype: MessageType
  onDelete: (id: number) => void
}) {
  const isUnread = msg.hasread === "未读"

  return (
    <div
      className={`flex items-start justify-between gap-4 rounded-lg border p-4 transition-colors ${
        isUnread ? "border-primary/30 bg-primary/5" : "border-border bg-card"
      }`}
    >
      <div className="min-w-0 flex-1">
        {/* 发送者 */}
        {msg.user && (
          <div className="mb-1 flex items-center gap-2">
            <Link
              href={`/user/${msg.user.username}`}
              className="text-sm font-medium text-foreground hover:text-primary"
            >
              {msg.user.username}
            </Link>
            {isUnread && (
              <span className="rounded-full bg-primary px-1.5 py-0.5 text-[10px] text-primary-foreground">
                新
              </span>
            )}
          </div>
        )}

        {/* 消息标题 + 对象 */}
        {msg.title && (
          <p className="text-sm text-muted-foreground">
            {msg.title}
            {msg.objurl && msg.objtitle && (
              <Link
                href={msg.objurl}
                className="ml-1 text-primary hover:underline"
              >
                {msg.objtitle}
              </Link>
            )}
          </p>
        )}

        {/* 消息内容 */}
        {msg.content && (
          <p className="mt-1 text-sm text-foreground line-clamp-2">
            {msg.content}
          </p>
        )}

        {/* 时间 */}
        <p className="mt-2 text-xs text-muted-foreground">{msg.ctime}</p>
      </div>

      {/* 删除按钮 */}
      <button
        onClick={() => onDelete(msg.id)}
        className="shrink-0 rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
        title="删除"
      >
        <Trash2 className="h-4 w-4" />
      </button>
    </div>
  )
}

export default function MessagePage() {
  const params = useParams()
  const router = useRouter()
  const msgtype = (params.msgtype as MessageType) || "system"
  const { isLoggedIn } = useAuth()

  const [messages, setMessages] = useState<Message[]>([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [hasMore, setHasMore] = useState(false)
  const [loading, setLoading] = useState(true)
  const [deleteId, setDeleteId] = useState<number | null>(null)
  const [deleting, setDeleting] = useState(false)

  // 验证 msgtype
  const validTypes: MessageType[] = ["system", "inbox", "outbox"]
  const currentType: MessageType = validTypes.includes(msgtype)
    ? msgtype
    : "system"

  useEffect(() => {
    // 未登录重定向
    if (!isLoggedIn) {
      router.replace(`/account/login?redirect=/message/${currentType}`)
      return
    }

    setLoading(true)
    fetchMessages(currentType, page).then((data) => {
      setMessages(data.messages ?? [])
      setTotal(data.total ?? 0)
      setHasMore(data.has_more ?? false)
      setLoading(false)
    })
  }, [currentType, page, isLoggedIn, router])

  function handleTabChange(type: string) {
    setPage(1)
    router.push(`/message/${type}`)
  }

  async function handleDelete() {
    if (deleteId === null) return
    setDeleting(true)
    const ok = await deleteMessage(deleteId, currentType)
    setDeleting(false)
    setDeleteId(null)
    if (ok) {
      setMessages((prev) => prev.filter((m) => m.id !== deleteId))
      setTotal((t) => Math.max(0, t - 1))
    }
  }

  const pageSize = 20
  const totalPages = total > 0 ? Math.ceil(total / pageSize) : hasMore ? page + 1 : page

  return (
    <PageLayout>
      <PageHeader
        title="消息中心"
        description="查看系统通知和私信"
        breadcrumbs={[{ label: "消息中心" }]}
        actions={
          <Link href="/message/send">
            <Button size="sm" className="gap-1.5">
              <Send className="h-3.5 w-3.5" />
              发私信
            </Button>
          </Link>
        }
      />

      {/* Tab 切换 */}
      <Tabs value={currentType} onValueChange={handleTabChange} className="mt-4">
        <TabsList>
          {(Object.entries(TAB_LABELS) as [MessageType, { label: string; icon: React.ReactNode }][]).map(
            ([type, { label, icon }]) => (
              <TabsTrigger key={type} value={type} className="gap-1.5">
                {icon}
                {label}
              </TabsTrigger>
            )
          )}
        </TabsList>
      </Tabs>

      {/* 消息列表 */}
      <div className="mt-4 space-y-3">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : messages.length === 0 ? (
          <div className="rounded-lg border border-dashed border-border py-12 text-center text-sm text-muted-foreground">
            暂无消息
          </div>
        ) : (
          messages.map((msg) => (
            <MessageItem
              key={msg.id}
              msg={msg}
              msgtype={currentType}
              onDelete={(id) => setDeleteId(id)}
            />
          ))
        )}
      </div>

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="mt-6 flex items-center justify-center gap-2">
          {page > 1 ? (
            <button
              onClick={() => setPage(page - 1)}
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
            >
              上一页
            </button>
          ) : (
            <button
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
              disabled
            >
              上一页
            </button>
          )}
          <span className="text-sm text-muted-foreground">
            第 {page} / {totalPages} 页
          </span>
          {hasMore || page < totalPages ? (
            <button
              onClick={() => setPage(page + 1)}
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-secondary-foreground hover:bg-secondary/80"
            >
              下一页
            </button>
          ) : (
            <button
              className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
              disabled
            >
              下一页
            </button>
          )}
        </div>
      )}

      {/* 删除确认对话框 */}
      <AlertDialog open={deleteId !== null} onOpenChange={(open) => !open && setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认删除</AlertDialogTitle>
            <AlertDialogDescription>
              确定要删除这条消息吗？删除后无法恢复。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleting}>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleting}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
              删除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </PageLayout>
  )
}
