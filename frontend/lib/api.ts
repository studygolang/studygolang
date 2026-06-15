// StudyGolang API 客户端
// 服务端（SSR）使用完整 URL，客户端使用相对路径（经 Next.js rewrite 代理）

import type {
  APIResponse,
  ArticleDetailData,
  ArticleListData,
  BindUser,
  BookListData,
  Book,
  Comment,
  CommentDetailData,
  LoginData,
  Me,
  ProjectDetailData,
  ProjectListData,
  ReadingListData,
  Reading,
  ResourceListData,
  Resource,
  SearchData,
  SiteStats,
  Topic,
  TopicDetailData,
  TopicListData,
  TopicNode,
  User,
  UserComment,
  Wiki,
  Announcement,
  AnnouncementListData,
} from './types'

// 服务端使用后端地址，客户端使用相对路径（经 rewrites 代理）
export function getAPIBase(): string {
  if (typeof window === 'undefined') {
    // 服务端 SSR 调用
    return process.env.API_BASE_URL || 'http://localhost:8090'
  }
  // 客户端调用，直接访问后端（避免 rewrite 的 cookie 问题）
  return process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8090'
}

// SSR 超时 10 秒，客户端超时 30 秒
const SSR_TIMEOUT = 10_000
const CLIENT_TIMEOUT = 30_000

function getTimeout(): number {
  return typeof window === 'undefined' ? SSR_TIMEOUT : CLIENT_TIMEOUT
}

/**
 * 创建带超时的 AbortController
 * 如果调用方已提供 signal，则两者任一触发即中止请求
 */
function createTimeoutSignal(existingSignal?: AbortSignal): {
  signal: AbortSignal
  cleanup: () => void
} {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), getTimeout())

  // 如果外部已有 signal，将其 abort 事件转发到 controller
  if (existingSignal) {
    if (existingSignal.aborted) {
      clearTimeout(timeoutId)
      controller.abort()
    } else {
      const onExternalAbort = () => {
        clearTimeout(timeoutId)
        controller.abort()
      }
      existingSignal.addEventListener('abort', onExternalAbort, { once: true })
    }
  }

  const cleanup = () => clearTimeout(timeoutId)
  return { signal: controller.signal, cleanup }
}

/**
 * 统一的 API 请求函数（抛异常模式）
 * 错误时抛出异常，成功时返回 json.data
 */
export async function fetchAPI<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const base = getAPIBase()
  const url = `${base}/api/v1${path}`

  const { signal, cleanup } = createTimeoutSignal(options.signal ?? undefined)

  try {
    const res = await fetch(url, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        ...options.headers,
      },
      ...options,
      signal,
    })

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`)
    }

    const json: APIResponse<T> = await res.json()
    if (json.code !== 0) {
      // 后端错误字段为 msg（非 message）
      throw new Error(json.msg || '请求失败')
    }

    return json.data
  } finally {
    cleanup()
  }
}

/**
 * 可空版本的 API 请求函数
 * 错误时返回 null，成功时返回 json.data
 */
export async function fetchAPINullable<T>(
  path: string,
  options: RequestInit = {}
): Promise<T | null> {
  try {
    return await fetchAPI<T>(path, options)
  } catch {
    return null
  }
}

// ======================== 话题 ========================
export const topicAPI = {
  getList(params: { p?: number; tab?: string } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    if (params.tab) q.set('tab', params.tab)
    const qs = q.toString()
    return fetchAPI<TopicListData>(`/topics${qs ? `?${qs}` : ''}`, fetchOptions)
  },

  getNoReply(params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<TopicListData>(`/topics/no_reply?${q}`, fetchOptions)
  },

  getLast(params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<TopicListData>(`/topics/last?${q}`, fetchOptions)
  },

  getDetail(tid: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<TopicDetailData>(`/topics/${tid}`, fetchOptions)
  },

  getNodeTopics(nid: number | string, params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<TopicListData>(`/topics/node/${nid}?${q}`, fetchOptions)
  },

  getNodes(fetchOptions?: RequestInit) {
    return fetchAPI<TopicNode[]>('/nodes', fetchOptions)
  },
}

// ======================== 文章 ========================
export const articleAPI = {
  getList(params: { p?: number; sort?: string; tag?: string } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    if (params.sort) q.set('sort', params.sort)
    if (params.tag) q.set('tag', params.tag)
    return fetchAPI<ArticleListData>(`/articles?${q}`, fetchOptions)
  },

  getDetail(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<ArticleDetailData>(`/articles/${id}`, fetchOptions)
  },
}

// ======================== 项目 ========================
export const projectAPI = {
  getList(params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<ProjectListData>(`/projects?${q}`, fetchOptions)
  },

  getDetail(uri: string, fetchOptions?: RequestInit) {
    return fetchAPI<ProjectDetailData>(`/projects/${uri}`, fetchOptions)
  },
}

// ======================== 资源 ========================
export const resourceAPI = {
  getList(params: { p?: number; catid?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    if (params.catid) q.set('catid', String(params.catid))
    return fetchAPI<ResourceListData>(`/resources?${q}`, fetchOptions)
  },

  getDetail(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<{ resource: Resource; comments: Comment[] }>(`/resources/${id}`, fetchOptions)
  },

  getCategories(fetchOptions?: RequestInit) {
    return fetchAPI<{ categories: Array<{ id: number; name: string }> }>('/resources/categories', fetchOptions)
  },
}

// ======================== 晨读 ========================
export const readingAPI = {
  getList(params: { lastid?: number; rtype?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.lastid) q.set('lastid', String(params.lastid))
    if (params.rtype !== undefined) q.set('rtype', String(params.rtype))
    return fetchAPI<ReadingListData>(`/readings?${q}`, fetchOptions)
  },

  getDetail(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<Reading>(`/readings/${id}`, fetchOptions)
  },
}

// ======================== 书籍 ========================
export const bookAPI = {
  getList(params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<BookListData>(`/books?${q}`, fetchOptions)
  },

  getDetail(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<{ book: Book; comments: Comment[] }>(`/books/${id}`, fetchOptions)
  },
}

// ======================== 用户 ========================
export const userAPI = {
  login(username: string, passwd: string) {
    return fetchAPI<LoginData>('/user/login', {
      method: 'POST',
      body: JSON.stringify({ username, passwd }),
    })
  },

  register(data: { username: string; email: string; passwd: string; captcha: string }) {
    return fetchAPI<{ uid: number }>('/user/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async getMe(token?: string, fetchOptions?: RequestInit): Promise<Me | null> {
    const headers: Record<string, string> = {}
    if (token) headers['X-Token'] = token
    // token 已迁移到 HttpOnly Cookie，通过 credentials 携带
    const options = { ...fetchOptions, headers, credentials: "include" as RequestCredentials }
    try {
      // 后端返回 { user: {...} }，提取 user 字段
      const data = await fetchAPI<{ user: Me }>('/user/me', options)
      return data?.user ?? null
    } catch {
      return null
    }
  },

  getProfile(username: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ user: User; topics: Topic[] }>(`/user/${username}`, fetchOptions)
  },

  // 获取当前登录用户的个人资料（需要登录）
  getMyProfile(fetchOptions?: RequestInit) {
    return fetchAPI<{ user: User; has_passwd: boolean }>('/user/profile', {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 更新个人资料
  updateProfile(data: {
    name?: string
    email?: string
    city?: string
    company?: string
    github?: string
    website?: string
    introduce?: string
    open?: string
  }) {
    return fetchAPI<null>('/user/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
      credentials: 'include',
    })
  },

  // 更新头像
  updateAvatar(avatar: string) {
    return fetchAPI<null>('/user/avatar', {
      method: 'PUT',
      body: JSON.stringify({ avatar }),
      credentials: 'include',
    })
  },

  // 修改密码
  changePassword(curPasswd: string, newPasswd: string) {
    return fetchAPI<null>('/user/password', {
      method: 'PUT',
      body: JSON.stringify({ cur_passwd: curPasswd, new_passwd: newPasswd }),
      credentials: 'include',
    })
  },

  // 获取用户评论列表
  getUserComments(username: string, params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<{ comments: UserComment[]; total: number; page: number; has_more: boolean }>(
      `/users/${username}/comments?${q}`,
      fetchOptions,
    )
  },
}

// ======================== 评论 ========================
export const commentAPI = {
  getList(objid: number, objtype: number, fetchOptions?: RequestInit) {
    return fetchAPI<Comment[]>(`/comments?objid=${objid}&objtype=${objtype}`, fetchOptions)
  },

  getDetail(cid: number | string, objid: number | string, objtype: number, fetchOptions?: RequestInit) {
    return fetchAPI<CommentDetailData>(`/comments/${cid}/detail?objid=${objid}&objtype=${objtype}`, fetchOptions)
  },

  // 发表评论（使用 form-urlencoded + Cookie 认证）
  create(objid: number, objtype: number, content: string) {
    const form = new URLSearchParams()
    form.set('objtype', String(objtype))
    form.set('content', content)
    return fetchAPI<Comment>(`/comments/${objid}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 修改评论（使用 form-urlencoded + Cookie 认证）
  update(cid: number | string, content: string) {
    const form = new URLSearchParams()
    form.set('content', content)
    return fetchAPI<Comment>(`/comments/${cid}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // @用户自动补全（后端返回裸数组，非标准 APIResponse 格式）
  async getAtUsers(term: string): Promise<Array<{ username: string; avatar: string }>> {
    const base = getAPIBase()
    try {
      const res = await fetch(`${base}/api/v1/at/users?term=${encodeURIComponent(term)}`, {
        credentials: 'include',
      })
      if (!res.ok) return []
      return res.json()
    } catch {
      return []
    }
  },
}

// ======================== 搜索 ========================
export const searchAPI = {
  search(keyword: string, params: { p?: number; type?: string } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams({ q: keyword })
    if (params.p) q.set('p', String(params.p))
    if (params.type) q.set('type', params.type)
    return fetchAPI<SearchData>(`/search?${q}`, fetchOptions)
  },
}

// ======================== 侧边栏 ========================
export const sidebarAPI = {
  getStats(fetchOptions?: RequestInit) {
    return fetchAPI<SiteStats>('/stat/site', fetchOptions)
  },

  // 后端返回 { topics: Topic[] }
  getRecentTopics(limit = 10, fetchOptions?: RequestInit) {
    return fetchAPI<{ topics: Topic[] }>(`/sidebar/topics/recent?limit=${limit}`, fetchOptions)
  },

  getRecentReadings(limit = 7, fetchOptions?: RequestInit) {
    return fetchAPI<{ readings: Reading[] }>(`/sidebar/readings/recent?limit=${limit}`, fetchOptions)
  },

  // 后端返回 { nodes: TopicNode[] }
  getHotNodes(fetchOptions?: RequestInit) {
    return fetchAPI<{ nodes: TopicNode[] }>('/sidebar/nodes/hot', fetchOptions)
  },

  // 后端返回 { users: User[] }
  getActiveUsers(fetchOptions?: RequestInit) {
    return fetchAPI<{ users: User[] }>('/sidebar/users/active', fetchOptions)
  },

  // 后端返回 { links: FriendLink[] }
  getFriendLinks(fetchOptions?: RequestInit) {
    return fetchAPI<{ links: import('./types').FriendLink[] }>('/sidebar/friend/links', fetchOptions)
  },

  // 后端返回 { comments: Comment[], [uid: string]: User }
  getRecentComments(fetchOptions?: RequestInit) {
    return fetchAPI<{ comments: Comment[] }>('/sidebar/comments/recent', fetchOptions)
  },
}

// ======================== 公告 ========================
// 重新导出 AnnouncementListData，使 import type 的使用被 IDE 正确识别
export type { AnnouncementListData } from './types'

export const announcementAPI = {
  getList(params: { p?: number; type?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    if (params.type !== undefined) q.set('type', String(params.type))
    const qs = q.toString()
    return fetchAPI<AnnouncementListData>(`/announcements${qs ? `?${qs}` : ''}`, fetchOptions)
  },
  getDetail(id: number, fetchOptions?: RequestInit) {
    return fetchAPI<{ announcement: Announcement }>(`/announcements/${id}`, fetchOptions)
  },
}

// ======================== Wiki ========================
export const wikiAPI = {
  // 后端返回 { wikis: Wiki[], page: { has_prev, prev_id, has_next, next_id } }
  getList(fetchOptions?: RequestInit) {
    return fetchAPI<{ wikis: Wiki[]; page: { has_prev: boolean; prev_id: number; has_next: boolean; next_id: number } }>('/wiki', fetchOptions)
  },

  getDetail(uri: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ wiki: Wiki }>(`/wiki/${uri}`, fetchOptions)
  },

  // 获取 Wiki 编辑数据（需要登录）
  getEdit(uri: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ wiki: { id: number; title: string; content: string; uri: string } }>(`/wiki/${uri}/edit`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 创建 Wiki
  create(data: { title: string; content: string; uri: string }) {
    const form = new URLSearchParams()
    form.set('title', data.title)
    form.set('content', data.content)
    form.set('uri', data.uri)
    return fetchAPI<{ message: string }>('/wiki', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 更新 Wiki
  update(id: number, data: { title: string; content: string }) {
    const form = new URLSearchParams()
    form.set('id', String(id))
    form.set('title', data.title)
    form.set('content', data.content)
    return fetchAPI<{ message: string }>(`/wiki/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },
}

// ======================== 私信 ========================
export const messageAPI = {
  list(type: string = 'inbox', params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    q.set('type', type)
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<{ messages: any[]; total: number; page: number; has_more: boolean }>(`/messages?${q}`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 发送私信
  send(toUid: number, content: string) {
    return fetchAPI<{ message: string }>('/messages', {
      method: 'POST',
      body: JSON.stringify({ to: toUid, content }),
      credentials: 'include',
    })
  },

  // 获取未读消息计数
  getUnreadCount() {
    return fetchAPI<{ system: number; inbox: number }>('/messages/unread-count', {
      credentials: 'include',
    })
  },

  // 删除私信
  delete(id: number | string) {
    return fetchAPI<null>(`/messages/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
  },
}

// ======================== 项目（补充写操作）========================
export const projectWriteAPI = {
  // 获取项目编辑数据（需要登录）
  getEdit(uri: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ project: any }>(`/projects/${uri}/edit`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 更新项目
  update(uri: string, data: Record<string, string>) {
    const form = new URLSearchParams()
    for (const [key, value] of Object.entries(data)) {
      form.set(key, value)
    }
    return fetchAPI<{ message: string }>(`/projects/${uri}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 检查 URI 唯一性
  checkUri(uri: string) {
    return fetchAPI<{ exists: boolean }>(`/projects/check_uri?uri=${encodeURIComponent(uri)}`)
  },
}

// ======================== 资源（补充写操作）========================
export const resourceWriteAPI = {
  // 获取资源编辑数据
  getEdit(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<{ resource: any }>(`/resources/${id}/edit`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 更新资源
  update(id: number | string, data: Record<string, string>) {
    const form = new URLSearchParams()
    for (const [key, value] of Object.entries(data)) {
      form.set(key, value)
    }
    return fetchAPI<{ message: string }>(`/resources/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },
}

// ======================== 图书（补充写操作）========================
export const bookWriteAPI = {
  // 获取图书编辑数据
  getEdit(id: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<{ book: any }>(`/books/${id}/edit`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 更新图书
  update(id: number | string, data: Record<string, string>) {
    const form = new URLSearchParams()
    for (const [key, value] of Object.entries(data)) {
      form.set(key, value)
    }
    return fetchAPI<{ message: string }>(`/books/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },
}

// ======================== 专栏 ========================
export const subjectAPI = {
  // 关注/取消关注专栏
  follow(sid: number) {
    const form = new URLSearchParams()
    form.set('sid', String(sid))
    return fetchAPI<{ followed: boolean }>('/subject/follow', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 投稿文章到专栏
  contribute(sid: number, aid: number) {
    const form = new URLSearchParams()
    form.set('sid', String(sid))
    form.set('article_id', String(aid))
    return fetchAPI<{ message: string }>('/subject/contribute', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 从专栏移除文章
  removeContribute(sid: number, aid: number) {
    const form = new URLSearchParams()
    form.set('sid', String(sid))
    form.set('article_id', String(aid))
    return fetchAPI<{ message: string }>('/subject/remove_contribute', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 修改专栏
  modify(data: Record<string, string>) {
    const form = new URLSearchParams()
    for (const [key, value] of Object.entries(data)) {
      form.set(key, value)
    }
    return fetchAPI<{ message: string }>('/subject/modify', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 获取我的文章列表（用于投稿选择）
  myArticles(fetchOptions?: RequestInit) {
    return fetchAPI<{ articles: any[] }>('/subject/my_articles', {
      ...fetchOptions,
      credentials: 'include',
    })
  },
}

// ======================== 账户管理 ========================
export const accountAPI = {
  // 激活账户（GET，因邮箱激活链接为浏览器直接打开）
  activate(token: string) {
    return fetchAPI<{ message: string }>('/account/activate?token=' + encodeURIComponent(token))
  },

  // 发送激活邮件
  sendActivateEmail() {
    return fetchAPI<{ message: string }>('/account/send-activate-email', {
      method: 'POST',
      credentials: 'include',
    })
  },

  // 邮件退订确认
  unsubscribePage(token: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ email: string }>(`/account/email/unsubscribe?token=${encodeURIComponent(token)}`, fetchOptions)
  },

  // 执行邮件退订
  unsubscribe(token: string) {
    const form = new URLSearchParams()
    form.set('token', token)
    return fetchAPI<{ message: string }>('/account/email/unsubscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 第三方账号解绑
  socialUnbind(bindId: number, platform: string) {
    const form = new URLSearchParams()
    form.set('bind_id', String(bindId))
    form.set('platform', platform)
    return fetchAPI<{ message: string }>('/account/social/unbind', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 获取已绑定的社交账号列表
  getBindUsers(): Promise<{ bind_users: BindUser[] }> {
    return fetchAPI<{ bind_users: BindUser[] }>('/account/bind_users', {
      credentials: 'include',
    })
  },
}

// ======================== GCTT ========================
export const gcttAPI = {
  // 获取当前用户 GCTT 信息
  getMe(fetchOptions?: RequestInit) {
    return fetchAPI<{ gctt_user: any; is_translator: boolean }>('/gctt/me', {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 申请成为译者
  apply() {
    return fetchAPI<{ gctt_user: any }>('/gctt/apply', {
      method: 'POST',
      credentials: 'include',
    })
  },

  // 发布翻译文章
  publish(title: string, content: string) {
    const form = new URLSearchParams()
    form.set('title', title)
    form.set('content', content)
    return fetchAPI<{ id: number }>('/gctt/articles', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },
}

// ======================== 话题（补充写操作）========================
export const topicWriteAPI = {
  // 置顶/取消置顶话题（管理员）
  setTop(tid: number | string) {
    return fetchAPI<{ message: string }>(`/topics/${tid}/set_top`, {
      method: 'POST',
      credentials: 'include',
    })
  },

  // 话题追加内容
  append(tid: number | string, content: string) {
    const form = new URLSearchParams()
    form.set('content', content)
    return fetchAPI<{ message: string }>(`/topics/${tid}/append`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 获取话题追加内容
  getAppends(tid: number | string, fetchOptions?: RequestInit) {
    return fetchAPI<{ appends: any[] }>(`/topics/${tid}/appends`, fetchOptions)
  },
}

// ======================== 点赞 ========================
export const likeAPI = {
  // 点赞/取消点赞（toggle）
  toggle(objid: number, objtype: number, flag: boolean) {
    const form = new URLSearchParams()
    form.set('objtype', String(objtype))
    form.set('flag', flag ? '1' : '0')
    return fetchAPI<{ liked: boolean }>(`/likes/${objid}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 查询点赞状态
  getStatus(objid: number, objtype: number, fetchOptions?: RequestInit) {
    return fetchAPI<{ has_like: boolean }>(`/likes/${objid}/status?objtype=${objtype}`, fetchOptions)
  },
}

// ======================== 收藏 ========================
export const favoriteAPI = {
  // 收藏/取消收藏（toggle）
  toggle(objid: number, objtype: number, flag: boolean) {
    const form = new URLSearchParams()
    form.set('objtype', String(objtype))
    form.set('collect', flag ? '1' : '0')
    return fetchAPI<{ collected: boolean }>(`/favorites/${objid}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: form.toString(),
      credentials: 'include',
    })
  },

  // 查询收藏状态
  getStatus(objid: number, objtype: number, fetchOptions?: RequestInit) {
    return fetchAPI<{ has_favorite: boolean }>(`/favorites/${objid}/status?objtype=${objtype}`, fetchOptions)
  },

  // 获取用户收藏列表
  listByUsername(username: string, params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
    return fetchAPI<{ favorites: any[]; total: number; page: number; has_more: boolean }>(
      `/users/${username}/favorites?${q}`,
      fetchOptions,
    )
  },
}

// ======================== 任务 ========================
export const missionAPI = {
  // 获取每日任务状态（需登录）
  getDaily(fetchOptions?: RequestInit) {
    return fetchAPI<{ login_mission: { uid: number; date: number } | null; had_redeem: boolean }>(`/mission/daily`, {
      ...fetchOptions,
      credentials: 'include',
    })
  },

  // 领取每日登录奖励
  dailyRedeem() {
    return fetchAPI<{ message: string }>('/mission/daily/redeem', {
      method: 'POST',
      credentials: 'include',
    })
  },

  // 完成任务
  complete(id: number | string) {
    return fetchAPI<{ message: string }>(`/mission/complete/${id}`, {
      method: 'POST',
      credentials: 'include',
    })
  },
}

// ======================== 余额 ========================
export const balanceAPI = {
  // 获取余额详情（需登录）
  getMyBalance(fetchOptions?: RequestInit) {
    return fetchAPI<{ balance: number; incomes: any[]; expenses: any[] }>('/balance', {
      ...fetchOptions,
      credentials: 'include',
    })
  },
}

// ======================== 首页聚合 ========================
export const homeAPI = {
  // 首页数据聚合
  getData(fetchOptions?: RequestInit) {
    return fetchAPI<{
      topics: any[]
      articles: any[]
      projects: any[]
      resources: any[]
      readings: any[]
      nodes: any[]
      stats: any
    }>('/home', fetchOptions)
  },
}

// ======================== 用户登出 ========================
export const authAPI = {
  // 登出
  logout() {
    return fetchAPI<{ message: string }>('/user/logout', {
      method: 'POST',
      credentials: 'include',
    })
  },
}
