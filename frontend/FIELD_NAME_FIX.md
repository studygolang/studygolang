# 字段名称不匹配修复

## 问题描述

发布话题后,首页和话题列表页都报错:
```
Cannot read properties of undefined (reading 'toString')
```

错误发生在 `formatNumber` 函数中,尝试访问 `topic.viewnum`。

## 根本原因

前端类型定义与后端 API 返回的字段名不匹配:

### 前端类型定义 (修复前)
```typescript
export interface Topic {
  replynum: number  // ❌ 前端使用这个
  likenum: number   // ❌ 前端使用这个
  viewnum: number   // ❌ 前端使用这个
  ...
}
```

### 后端返回数据
```json
{
  "tid": 1,
  "title": "这是第一篇文章",
  "reply": 0,      // ✅ 后端返回这个
  "like": 0,       // ✅ 后端返回这个
  "view": 3,       // ✅ 后端返回这个
  "viewnum": null, // ❌ 这些字段是 null
  "replynum": null,
  "likenum": null
}
```

### 问题
- 后端使用短字段名: `view`, `reply`, `like`
- 前端期望长字段名: `viewnum`, `replynum`, `likenum`
- 访问 `undefined` 字段导致运行时错误

## 解决方案

### 1. 更新类型定义

在 `frontend/lib/types.ts` 中同时支持两种字段名:

```typescript
export interface Topic {
  tid: number
  title: string
  content: string
  uid: number
  name: string
  avatar: string
  nid: number
  node?: TopicNode
  user?: User
  lastreplyuid: number
  lastreplyname: string
  // 后端返回的字段名 (主要)
  reply: number      // 回复数
  like: number       // 点赞数
  view: number       // 浏览数
  // 兼容旧字段名 (可选)
  replynum?: number
  likenum?: number
  viewnum?: number
  top: number
  ctime: string
  mtime: string
  permission: number
}
```

### 2. 更新组件

在所有使用这些字段的组件中使用兼容的访问方式:

#### topic-list.tsx
```typescript
// 修复前
{formatNumber(topic.viewnum)}
{topic.replynum}
{topic.likenum}

// 修复后
{formatNumber(topic.view || topic.viewnum || 0)}
{topic.reply || topic.replynum || 0}
{topic.like || topic.likenum || 0}
```

#### topic-detail.tsx
```typescript
// 修复前
{topic.viewnum}
{liked ? topic.likenum + 1 : topic.likenum}

// 修复后
{topic.view || topic.viewnum || 0}
{liked ? (topic.like || topic.likenum || 0) + 1 : (topic.like || topic.likenum || 0)}
```

## 兼容性说明

修复后的代码:
- **优先使用**: 后端返回的 `view`, `reply`, `like`
- **回退到**: 旧的 `viewnum`, `replynum`, `likenum` (如果存在)
- **默认值**: 如果都不存在,使用 `0`

这样可以同时兼容:
- 当前后端返回的数据结构
- 未来可能的字段名变更
- 避免 `undefined` 导致的运行时错误

## 测试验证

```bash
# 1. 测试后端 API
curl 'http://localhost:8090/api/v1/topics?p=1' | jq '.data.list[0] | {view, reply, like}'
# 输出: {"view": 3, "reply": 0, "like": 0}

# 2. 测试首页
curl http://localhost:3000/ | grep "这是第一篇文章"
# 输出: 这是第一篇文章

# 3. 测试话题列表页
curl http://localhost:3000/topics | grep "这是第一篇文章"
# 输出: 这是第一篇文章

# 4. 浏览器访问
# http://localhost:3000/
# http://localhost:3000/topics
# 应该都正常显示
```

## 相关文件

- `frontend/lib/types.ts` - 类型定义
- `frontend/components/topic-list.tsx` - 话题列表组件
- `frontend/components/topic-detail.tsx` - 话题详情组件

## 后续建议

1. **统一字段名**: 考虑在后端统一使用长字段名或短字段名
2. **类型安全**: 使用 TypeScript 的可选链操作符 `?.` 和空值合并操作符 `??`
3. **数据验证**: 在前端添加数据验证,确保必需字段存在
4. **API 文档**: 更新 API 文档,明确字段名称和类型

## 修复状态

✅ **已修复** - 首页和话题列表页现在可以正常显示
