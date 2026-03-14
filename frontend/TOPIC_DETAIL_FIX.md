# 话题详情页修复

## 问题描述

发布话题成功后,跳转到 `/topics/1` 显示"话题不存在或已被删除",但数据库中记录存在。

## 根本原因

前端类型定义与后端 API 返回的数据结构不匹配:

### 后端返回结构
```json
{
  "code": 0,
  "data": {
    "topic": {
      "tid": 1,
      "title": "这是第一篇文章",
      "uid": 6,
      "user": {
        "uid": 6,
        "username": "polarisxu",
        "email": "polaris@studygolang.com",
        ...
      },
      ...
    },
    "replies": []
  }
}
```

### 前端类型定义 (修复前)
```typescript
export interface Topic {
  tid: number
  title: string
  uid: number
  name: string  // ❌ 只有这个字段
  ...
}
```

### 问题
- 列表接口返回的是 `name` 字段 (用户名字符串)
- 详情接口返回的是 `user` 对象 (完整用户信息)
- 前端组件使用 `topic.name` 访问用户名,但详情接口中没有这个字段

## 解决方案

### 1. 更新类型定义

在 `frontend/lib/types.ts` 中添加 `user` 字段:

```typescript
export interface Topic {
  tid: number
  title: string
  content: string
  uid: number
  name: string      // 列表接口返回的用户名
  avatar: string
  nid: number
  node?: TopicNode
  user?: User       // ✅ 详情接口返回的完整用户信息
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
```

### 2. 更新组件

在 `frontend/components/topic-detail.tsx` 中使用兼容的字段访问:

```typescript
// 修复前
<AvatarFallback>
  {topic.name ? topic.name.charAt(0).toUpperCase() : "?"}
</AvatarFallback>
<Link href={`/user/${topic.name}`}>
  {topic.name}
</Link>

// 修复后
<AvatarFallback>
  {topic.user?.username 
    ? topic.user.username.charAt(0).toUpperCase() 
    : (topic.name ? topic.name.charAt(0).toUpperCase() : "?")}
</AvatarFallback>
<Link href={`/user/${topic.user?.username || topic.name}`}>
  {topic.user?.username || topic.name}
</Link>
```

## 兼容性说明

修复后的代码同时兼容:
- **列表接口**: 使用 `topic.name` (字符串)
- **详情接口**: 优先使用 `topic.user.username`,回退到 `topic.name`

## 测试验证

```bash
# 1. 测试后端 API
curl http://localhost:8090/api/v1/topics/1 | jq '.data.topic.user.username'
# 输出: "polarisxu"

# 2. 测试前端页面
curl http://localhost:3000/topics/1 | grep "这是第一篇文章"
# 输出: 这是第一篇文章

# 3. 浏览器访问
# http://localhost:3000/topics/1
# 应该正常显示话题内容
```

## 相关文件

- `frontend/lib/types.ts` - 类型定义
- `frontend/components/topic-detail.tsx` - 话题详情组件
- `internal/logic/topic.go` - 后端逻辑 (FindByTid 方法)
- `internal/http/controller/api/topic.go` - 后端 API 控制器

## 后续建议

1. **统一数据结构**: 考虑让列表接口也返回 `user` 对象,保持一致性
2. **类型安全**: 使用 TypeScript 的类型守卫确保字段存在
3. **错误处理**: 添加更详细的错误信息,帮助调试

## 修复状态

✅ **已修复** - 话题详情页现在可以正常显示
