// 草稿自动保存工具
// 使用 localStorage 存储，JSON 序列化

export interface DraftData {
  title: string
  content: string
  savedAt: number
}

const DRAFT_PREFIX = "draft:"

/** 生成草稿存储 key */
export function draftKey(type: string, id?: string): string {
  return `${DRAFT_PREFIX}${type}:${id || "new"}`
}

/** 保存草稿到 localStorage */
export function saveDraft(
  key: string,
  data: { title: string; content: string }
): void {
  const draft: DraftData = {
    title: data.title,
    content: data.content,
    savedAt: Date.now(),
  }
  try {
    localStorage.setItem(key, JSON.stringify(draft))
  } catch {
    // localStorage 已满或不可用时静默失败
  }
}

/** 从 localStorage 加载草稿 */
export function loadDraft(key: string): DraftData | null {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return null
    const draft: DraftData = JSON.parse(raw)
    if (typeof draft.title === "string" && typeof draft.content === "string") {
      return draft
    }
    return null
  } catch {
    return null
  }
}

/** 清除指定草稿 */
export function clearDraft(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch {
    // 静默失败
  }
}
