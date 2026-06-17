import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import type { TopicNode } from "@/lib/types"
import { loadDraft, clearDraft, draftKey, type DraftData } from "@/lib/draft"

// 客户端组件使用相对路径，通过 next.config.mjs 中的 rewrites 代理到后端
const API_BASE = ""

type ContentType = "topic" | "article" | "project" | "resource" | "book"

// 节点分组类型
type NodeGroup = {
  category: string
  nodes: TopicNode[]
}

async function fetchNodes(): Promise<NodeGroup[]> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/nodes`)
    if (!res.ok) return []
    const json = await res.json()

    if (json.code !== 0 || !Array.isArray(json.data)) return []

    // 后端返回扁平数组：TopicNode[]
    // 按 parent 字段聚合：parent=0 的是分类节点（category），其他归到对应 parent 下
    const allNodes: TopicNode[] = json.data.map((node: any) => ({
      nid: node.nid,
      name: node.name,
      ename: node.ename,
      parent: node.pid ?? node.parent ?? 0,
      seq: node.seq || 0,
      intro: node.intro || '',
      logo: node.logo || '',
      show_index: node.show_index ?? false,
    }))

    // 父节点（parent=0）作为分组标题，其子节点作为组内节点
    const rootNodes = allNodes.filter(n => !n.parent)
    const groups: NodeGroup[] = rootNodes.map(root => ({
      category: root.name,
      nodes: allNodes.filter(n => n.parent === root.nid),
    }))

    // 没有父节点的散落节点（理论上不应该有），合并到"其他"
    const grouped = new Set<number>([0, ...rootNodes.map(n => n.nid)])
    const orphan = allNodes.filter(n => !grouped.has(n.parent))
    if (orphan.length > 0) {
      groups.push({ category: '其他', nodes: orphan })
    }

    return groups
  } catch {
    return []
  }
}

interface UsePublishProps {
  initialType: ContentType
}

export function usePublish({ initialType }: UsePublishProps) {
  const router = useRouter()

  const [contentType, setContentType] = useState<ContentType>(initialType)
  const [nodeGroups, setNodeGroups] = useState<NodeGroup[]>([])
  const [title, setTitle] = useState("")
  const [content, setContent] = useState("")
  const [nid, setNid] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [draftData, setDraftData] = useState<DraftData | null>(null)

  // 项目特定字段
  const [src, setSrc] = useState("")
  const [category, setCategory] = useState("")
  const [home, setHome] = useState("")
  const [doc, setDoc] = useState("")
  const [download, setDownload] = useState("")
  const [licence, setLicence] = useState("")
  const [lang, setLang] = useState("Go")
  const [os, setOs] = useState("")

  // 资源特定字段
  const [resourceForm, setResourceForm] = useState("link")
  const [resourceUrl, setResourceUrl] = useState("")
  const [resourceCatid, setResourceCatid] = useState("")

  // 图书特定字段
  const [bookAuthor, setBookAuthor] = useState("")
  const [bookTranslator, setBookTranslator] = useState("")
  const [bookCover, setBookCover] = useState("")
  const [bookPubDate, setBookPubDate] = useState("")
  const [bookLang, setBookLang] = useState("中文")
  const [bookIsFree, setBookIsFree] = useState(false)
  const [bookOnlineUrl, setBookOnlineUrl] = useState("")
  const [bookDownloadUrl, setBookDownloadUrl] = useState("")
  const [bookBuyUrl, setBookBuyUrl] = useState("")
  const [bookPrice, setBookPrice] = useState("")

  // 加载节点
  useEffect(() => {
    fetchNodes().then(setNodeGroups)
  }, [])

  // 加载草稿
  useEffect(() => {
    const key = draftKey(initialType)
    const draft = loadDraft(key)
    if (draft) {
      setDraftData(draft)
    }
  }, [initialType])

  // 恢复草稿内容
  function restoreDraft() {
    if (!draftData) return
    setTitle(draftData.title)
    setContent(draftData.content)
    setDraftData(null)
  }

  // 放弃草稿
  function discardDraft() {
    const key = draftKey(contentType)
    clearDraft(key)
    setDraftData(null)
  }

  // Ctrl/Cmd+Enter 提交
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        document.querySelector<HTMLButtonElement>('[data-publish-btn]')?.click()
      }
    }
    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  // 标签输入处理
  function handleTagKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if ((e.key === "Enter" || e.key === ",") && tagInput.trim()) {
      e.preventDefault()
      const tag = tagInput.trim()
      if (tags.length < 5 && !tags.includes(tag)) {
        setTags([...tags, tag])
      }
      setTagInput("")
    } else if (e.key === "Backspace" && !tagInput && tags.length > 0) {
      setTags(tags.slice(0, -1))
    }
  }

  // 发布处理
  async function handlePublish() {
    setError("")
    if (!title.trim()) { setError("请填写标题"); return }
    if (contentType === "topic" && !nid) { setError("请选择节点"); return }
    if (contentType === "project" && !src.trim()) { setError("请填写源码地址"); return }
    if (contentType === "resource" && resourceForm === "link" && !resourceUrl.trim()) { setError("请填写资源链接"); return }
    if (contentType === "resource" && resourceForm === "content" && !content.trim()) { setError("请填写资源内容"); return }
    if (contentType === "book" && !bookAuthor.trim()) { setError("请填写作者"); return }

    setSubmitting(true)
    // 发布时清除草稿
    clearDraft(draftKey(contentType))
    try {
      if (contentType === "topic") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("content", content.trim())
        form.set("nid", nid)
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/topics`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.tid ? `/topics/${json.data.tid}` : "/topics")
      } else if (contentType === "article") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("content", content.trim())
        form.set("txt", content.trim())
        form.set("cover", "")
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/articles`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/articles/${json.data.id}` : "/articles")
      } else if (contentType === "project") {
        const form = new URLSearchParams()
        form.set("name", title.trim())
        form.set("src", src.trim())
        form.set("desc", content.trim())
        if (category.trim()) form.set("category", category.trim())
        if (home.trim()) form.set("home", home.trim())
        if (doc.trim()) form.set("doc", doc.trim())
        if (download.trim()) form.set("download", download.trim())
        if (licence.trim()) form.set("licence", licence.trim())
        if (lang.trim()) form.set("lang", lang.trim())
        if (os.trim()) form.set("os", os.trim())
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/projects`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        // 项目详情路由是 /p/[uri]（见 frontend/app/p/[uri]/page.tsx）
        router.push(json.data?.uri ? `/p/${json.data.uri}` : "/projects")
      } else if (contentType === "resource") {
        const form = new URLSearchParams()
        form.set("title", title.trim())
        form.set("form", resourceForm === "link" ? "只是链接" : "包括内容")
        if (resourceForm === "link") {
          form.set("url", resourceUrl.trim())
        } else {
          form.set("content", content.trim())
        }
        if (resourceCatid) form.set("catid", resourceCatid)
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/resources`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/resources/${json.data.id}` : "/resources")
      } else if (contentType === "book") {
        const form = new URLSearchParams()
        form.set("name", title.trim())
        form.set("author", bookAuthor.trim())
        form.set("desc", content.trim())
        if (bookTranslator.trim()) form.set("translator", bookTranslator.trim())
        if (bookCover.trim()) form.set("cover", bookCover.trim())
        if (bookPubDate.trim()) form.set("pub_date", bookPubDate.trim())
        if (bookLang.trim()) form.set("lang", bookLang.trim())
        form.set("is_free", bookIsFree ? "1" : "0")
        if (bookOnlineUrl.trim()) form.set("online_url", bookOnlineUrl.trim())
        if (bookDownloadUrl.trim()) form.set("download_url", bookDownloadUrl.trim())
        if (bookBuyUrl.trim()) form.set("buy_url", bookBuyUrl.trim())
        if (bookPrice.trim()) form.set("price", bookPrice.trim())
        if (tags.length > 0) form.set("tags", tags.join(","))

        const res = await fetch(`${API_BASE}/api/v1/books`, {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          credentials: "include",
          body: form.toString(),
        })
        const json = await res.json()
        if (json.code !== 0) { setError(json.msg || "发布失败"); return }
        router.push(json.data?.id ? `/books/${json.data.id}` : "/books")
      }
    } catch {
      setError("网络错误，请稍后重试")
    } finally {
      setSubmitting(false)
    }
  }

  function handleCancel() {
    router.back()
  }

  return {
    contentType,
    setContentType,
    nodeGroups,
    title,
    setTitle,
    content,
    setContent,
    nid,
    setNid,
    tags,
    setTags,
    tagInput,
    setTagInput,
    handleTagKeyDown,
    submitting,
    error,
    handlePublish,
    handleCancel,
    // 草稿
    draftData,
    restoreDraft,
    discardDraft,
    // 项目字段
    src,
    setSrc,
    category,
    setCategory,
    home,
    setHome,
    doc,
    setDoc,
    download,
    setDownload,
    licence,
    setLicence,
    lang,
    setLang,
    os,
    setOs,
    // 资源字段
    resourceForm,
    setResourceForm,
    resourceUrl,
    setResourceUrl,
    resourceCatid,
    setResourceCatid,
    // 图书字段
    bookAuthor,
    setBookAuthor,
    bookTranslator,
    setBookTranslator,
    bookCover,
    setBookCover,
    bookPubDate,
    setBookPubDate,
    bookLang,
    setBookLang,
    bookIsFree,
    setBookIsFree,
    bookOnlineUrl,
    setBookOnlineUrl,
    bookDownloadUrl,
    setBookDownloadUrl,
    bookBuyUrl,
    setBookBuyUrl,
    bookPrice,
    setBookPrice,
  }
}

export type { ContentType, NodeGroup }
