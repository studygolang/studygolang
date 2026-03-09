// StudyGolang 前端类型定义

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
  name: string
  avatar: string
  nid: number
  node?: TopicNode
  lastreplyuid: number
  lastreplyname: string
  replynum: number
  likenum: number
  viewnum: number
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
  topics: Topic[]
  tab: string
  tab_list: TopicNode[]
}

export interface TopicDetailData {
  topic: Topic
  replies: TopicReply[]
}

// ======================== 文章相关 ========================
export interface Article {
  id: number
  title: string
  content: string
  summary: string
  author: string
  author_uid: number
  author_avatar: string
  tags: string
  pubdate: string
  ctime: string
  mtime: string
  viewnum: number
  likenum: number
  cmtnum: number
  top: number
  url: string
}

export interface ArticleListData extends Pagination {
  articles: Article[]
}

export interface ArticleDetailData {
  article: Article
  prev?: Article
  next?: Article
  comments?: Comment[]
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
export interface Resource {
  id: number
  title: string
  url: string
  cover: string
  author: string
  author_uid: number
  catid: number
  catname: string
  rtype: number
  desc: string
  viewnum: number
  cmtnum: number
  likenum: number
  ctime: string
}

export interface ResourceListData extends Pagination {
  resources: Resource[]
}

// ======================== 晨读相关 ========================
export interface Reading {
  id: number
  title: string
  url: string
  cover: string
  rtype: number
  lang: string
  desc: string
  clicknum: number
  ctime: string
}

export interface ReadingListData {
  readings: Reading[]
  has_prev: boolean
  has_next: boolean
  prev_id: number
  next_id: number
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
  tagline: string
  website: string
  city: string
  github: string
  company: string
  level: number
  balance: number
  follow_count: number
  fans_count: number
  ctime: string
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
export interface SearchResult {
  id: string
  title: string
  content: string
  url: string
  type: string
  author: string
  ctime: string
}

export interface SearchData {
  results: SearchResult[]
  total: number
  keyword: string
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
