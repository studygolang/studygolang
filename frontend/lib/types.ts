// StudyGolang 前端类型定义
// 字段名与后端 Go model JSON tag 保持一致

export interface APIResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

export interface Pagination {
  page: number
  total: number
  has_more: boolean
}

// ======================== 动态相关 ========================
// 动态类型常量（与后端 model 保持一致）
export const OBJTYPE_TOPIC = 1
export const OBJTYPE_ARTICLE = 2
export const OBJTYPE_RESOURCE = 3
export const OBJTYPE_PROJECT = 4
export const OBJTYPE_BOOK = 5

// Feed 动态（聚合话题、文章、项目、资源等）
// 注意：后端 model.Feed 没有 JSON tag，字段名首字母大写
// 但 updated_at 有 JSON tag，是小写的
export interface Feed {
  Id: number
  Title: string
  Objid: number
  Objtype: number      // 1=话题 2=文章 3=资源 4=项目 5=书籍
  Uid: number
  Author: string
  Nid: number
  Lastreplyuid: number
  Lastreplytime: string
  Tags: string
  Cmtnum: number
  Likenum: number
  Top: number
  Seq: number
  State: number
  CreatedAt: string
  updated_at: string   // 后端有 JSON tag，是小写的
  // 关联数据
  User?: User
  Lastreplyuser?: User
  Node?: TopicNode | { name: string }
  Uri: string          // 详情页路径（不含域名）
}

export interface FeedListData extends Pagination {
  feeds?: Feed[]
  tab: string
  tab_list: TopicNode[]
}

// ======================== 话题相关 ========================
export interface TopicNode {
  id: number
  name: string
  ename: string
  parent_id: number
  seq: number
  pid: number
  intro: string
  logo: string
  style: string
}

export interface Topic {
  tid: number
  title: string
  content: string
  uid: number
  name: string  // 列表接口返回的用户名
  avatar: string
  nid: number
  node?: TopicNode
  user?: User  // 详情接口返回的完整用户信息
  lastreplyuid: number
  lastreplyname: string
  // 后端返回的字段名
  reply: number      // 回复数
  like: number       // 点赞数
  view: number       // 浏览数
  // 兼容旧字段名
  replynum?: number
  likenum?: number
  viewnum?: number
  top: number
  ctime: string
  mtime: string
  permission: number
}

export interface TopicReply {
  id: number
  tid: number
  uid: number
  name: string
  avatar: string
  content: string
  ctime: string
}

export interface TopicListData extends Pagination {
  // /home 接口返回 "topics" 字段；/topics、/topics/node/:nid 接口返回 "list" 字段
  // 两个字段均为可选，各页面按实际接口取对应字段
  topics?: Topic[]
  list?: Topic[]
  tab: string
  tab_list: TopicNode[]
}

export interface TopicDetailData {
  topic: Topic
  replies: TopicReply[]
}

// ======================== 文章相关 ========================
// 对应后端 model.Article 的 JSON 字段
export interface Article {
  id: number
  title: string
  content: string
  // 后端 Article 无 summary 字段，用 content 前段截取作摘要
  author: string       // 原文作者
  name?: string        // 发布者用户名（来自 User 关联）
  lang: number
  pub_date: string
  url: string
  tags: string
  viewnum: number
  likenum: number
  cmtnum: number
  top: number
  ctime: string
  // mtime 在后端 Article 中没有直接的 JSON 字段（xorm 只读），按实际情况可能返回
  markdown: boolean
  gctt: boolean
}

// 文章列表 API 返回：{ list, total, page, has_more }
export interface ArticleListData {
  list: Article[]      // 后端返回字段名是 list，非 articles
  total: number
  page: number
  has_more: boolean
}

// 文章详情 API 返回：{ article, replies, prev_next }
export interface ArticleComment {
  id: number
  uid: number
  name: string
  avatar: string
  content: string
  ctime: string
  floor: number
}

export interface ArticleDetailData {
  article: Article
  replies?: ArticleComment[]   // 后端字段名是 replies，非 comments
  prev_next?: Article[]        // 后端字段名是 prev_next（数组，[prev, next]）
}

// ======================== 项目相关 ========================
export interface Project {
  id: number
  name: string
  uri: string
  category: string
  lang: string
  logo: string
  desc: string
  homepage: string
  doc_url: string
  src_url: string
  download_url: string
  author: string
  author_uid: number
  author_avatar: string
  star: number
  fork: number
  watch: number
  score: number
  viewnum: number
  cmtnum: number
  likenum: number
  ctime: string
}

export interface ProjectListData extends Pagination {
  projects: Project[]
}

export interface ProjectDetailData {
  project: Project
  comments?: Comment[]
}

// ======================== 资源相关 ========================
// 对应后端 model.ResourceInfo（Resource + ResourceEx 合并）
export interface Resource {
  id: number
  title: string
  url: string
  uid: number
  catid: number
  catname?: string     // 由逻辑层动态注入
  form: string         // 资源形式
  content: string      // 资源内容/描述
  tags: string
  viewnum: number
  cmtnum: number
  likenum: number
  ctime: string
  // 列表中还有 user 子对象
  user?: {
    uid: number
    username: string
    name: string
    avatar: string
  }
  host?: string        // 资源域名
}

// 资源列表 API 返回：{ resources, has_more }（无 total 字段）
export interface ResourceListData {
  resources: Resource[]
  has_more: boolean
}

// ======================== 晨读相关 ========================
// 对应后端 model.MorningReading 的 JSON 字段
export interface Reading {
  id: number
  content: string      // 晨读标题/描述
  rtype: number        // 类型
  inner: number        // 是否站内
  url: string          // 外链
  moreurls: string     // 多个 URL，逗号分隔
  username: string     // 发布者用户名
  clicknum: number
  ctime: string
  rdate?: string       // 晨读日期，格式 2006-01-02
  urls?: string[]      // moreurls 解析后的数组
}

// 晨读列表 API 返回：{ readings, rtype, page: { has_prev, prev_id, has_next, next_id } }
export interface ReadingPage {
  has_prev: boolean
  prev_id: number
  has_next: boolean
  next_id: number
}

export interface ReadingListData {
  readings: Reading[]
  rtype: number
  page: ReadingPage
}

// ======================== 书籍相关 ========================
export interface Book {
  id: number
  name: string
  cover: string
  author: string
  translator: string
  lang: string
  desc: string
  url: string
  buyer_show: number
  status: number
  viewnum: number
  cmtnum: number
  likenum: number
  ctime: string
}

export interface BookListData extends Pagination {
  books: Book[]
}

// ======================== 用户相关 ========================
export interface User {
  uid: number
  username: string
  name: string
  avatar: string
  email: string
  is_root: boolean
  is_vip: boolean
  open_id: string
  open: number       // 是否公开个人信息（1=公开，0=不公开）
  introduce: string  // 个人简介（后端字段名）
  website: string
  city: string
  github: string
  company: string
  level: number
  balance: number
  ctime: string
  weight?: number    // 活跃度权重（DAU 排名用）
}

export interface UserListData {
  active_users: User[]
  new_users: User[]
  total: number
}

export interface Me {
  uid: number
  username: string
  name: string
  avatar: string
  is_root: boolean
  is_vip: boolean
}

export interface LoginData {
  token: string
  uid: number
  username: string
}

// ======================== 评论相关 ========================
export interface Comment {
  id: number
  uid: number
  name: string
  avatar: string
  content: string
  ctime: string
  floor: number
}

// ======================== 侧边栏相关 ========================
export interface SiteStats {
  article: number
  project: number
  topic: number
  resource: number
  book: number
  comment: number
  user: number
}

export interface FriendLink {
  id: number
  name: string
  url: string
  logo: string
}

// ======================== 搜索相关 ========================
// 对应后端 model.Document 的 JSON 字段
export interface SearchResult {
  id: string
  objid: number
  objtype: number   // 1=话题 2=文章 3=资源 4=项目
  title: string
  content: string
  author: string
  uid: number
  pub_time: string  // 后端字段名，非 ctime
  tags: string
  viewnum: number
  cmtnum: number
  likenum: number
}

export interface SearchData {
  results: SearchResult[]
  total: number
  keyword: string
  page: number
  has_more: boolean
}

// ======================== 招聘相关 ========================
export interface Job {
  id: number
  title: string
  company: string
  company_size: string
  city: string
  salary_min: number
  salary_max: number
  experience: string
  tags: string
  uid: number
  status: number
  ctime: string
}

export interface JobListData extends Pagination {
  jobs: Job[]
  total: number
}

// ======================== 面试题相关 ========================
export interface InterviewQuestion {
  id: number
  sn: number
  show_sn: string
  question: string  // 已经过 Markdown 渲染，是 HTML 字符串
  answer: string    // 已经过 Markdown 渲染，是 HTML 字符串
  level: number     // 0=低 1=中 2=高
  viewnum: number
  cmtnum: number
  likenum: number
  source: string
  created_at: string
}

// 面试题列表 API 返回
export interface InterviewListData {
  questions: InterviewQuestion[]
  total: number
  page: number
  total_pages: number
  has_more: boolean
}

// ======================== Wiki 相关 ========================
export interface Wiki {
  id: number
  title: string
  uri: string
  content: string
  uid: number
  author: string
  ctime: string
  mtime: string
}
