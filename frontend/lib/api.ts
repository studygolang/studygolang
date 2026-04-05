// StudyGolang API 客户端
// 服务端（SSR）使用完整 URL，客户端使用相对路径（经 Next.js rewrite 代理）

import type {
  APIResponse,
  ArticleDetailData,
  ArticleListData,
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
} from './types'

// 服务端使用后端地址，客户端使用相对路径（经 rewrites 代理）
function getAPIBase(): string {
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
      headers: {
        'Content-Type': 'application/json',
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
  getList(params: { p?: number } = {}, fetchOptions?: RequestInit) {
    const q = new URLSearchParams()
    if (params.p) q.set('p', String(params.p))
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

  create(objid: number, content: string, token: string) {
    return fetchAPI<Comment>(`/comments/${objid}`, {
      method: 'POST',
      headers: { 'X-Token': token },
      body: JSON.stringify({ content }),
    })
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

// ======================== Wiki ========================
export const wikiAPI = {
  // 后端返回 { wikis: Wiki[], page: { has_prev, prev_id, has_next, next_id } }
  getList(fetchOptions?: RequestInit) {
    return fetchAPI<{ wikis: Wiki[]; page: { has_prev: boolean; prev_id: number; has_next: boolean; next_id: number } }>('/wiki', fetchOptions)
  },

  getDetail(uri: string, fetchOptions?: RequestInit) {
    return fetchAPI<{ wiki: Wiki }>(`/wiki/${uri}`, fetchOptions)
  },
}
