# BlockNote 编辑器集成说明

## 已完成的工作

### 1. 安装依赖
已成功安装以下依赖包：
- @blocknote/core@0.47.1
- @blocknote/mantine@0.47.1
- @blocknote/react@0.47.1

### 2. 创建编辑器组件
文件：`/Users/polarisxu/project/golang/studygolang/frontend/components/blocknote-editor.tsx`

功能：
- 支持 Markdown 输入和输出
- 自动将编辑器内容转换为 Markdown
- 初始化时可加载现有 Markdown 内容
- 实时同步内容变化到父组件

### 3. 更新发布页面
文件：`/Users/polarisxu/project/golang/studygolang/frontend/app/publish/page.tsx`

变更：
- ✅ 移除了旧的 textarea 编辑器
- ✅ 移除了 Markdown 工具栏（BlockNote 自带完整工具栏）
- ✅ 移除了预览功能（BlockNote 是所见即所得编辑器）
- ✅ 移除了相关的 state（preview、textareaRef）
- ✅ 移除了 insertMarkdown 函数
- ✅ 集成 BlockNote 编辑器组件
- ✅ 保留了所有其他功能（标题、节点、标签、提交等）

### 4. 添加自定义样式
文件：`/Users/polarisxu/project/golang/studygolang/frontend/app/publish/blocknote.css`

样式调整：
- 最小高度 300px
- 与页面主题协调
- 支持深色模式
- 移除默认边框，使用页面统一样式

## BlockNote 特性

BlockNote 是一个现代化的块编辑器，提供以下特性：

- ✨ 所见即所得编辑体验
- 🎨 内置完整工具栏（标题、粗体、斜体、代码、列表、引用等）
- 🔄 支持拖拽排序块
- ⚡ 支持斜杠命令（/）快速插入块
- 📝 自动保存为 Markdown 格式
- 🎯 支持代码块语法高亮
- 🔗 支持链接和图片
- 📋 支持表格
- ✅ 支持任务列表

## 使用方式

用户在 `/publish` 页面的使用流程：

1. 选择内容类型（主题/文章/项目）
2. 填写标题
3. 选择节点（如果是主题）
4. 添加标签（可选）
5. 在 BlockNote 编辑器中编写内容：
   - 使用顶部工具栏格式化文本
   - 使用 `/` 命令快速插入块
   - 拖拽块进行重新排序
6. 点击发布按钮

编辑器会自动将内容转换为 Markdown 格式提交到后端 API。

## 技术细节

### 数据流
```
用户输入 → BlockNote 编辑器 → blocksToMarkdownLossy() → Markdown 字符串 → onChange 回调 → content state → 提交到后端
```

### 组件通信
```typescript
<BlockNoteEditor
  content={content}           // 初始 Markdown 内容
  onChange={setContent}       // 内容变化回调
  placeholder="..."           // 占位符文本
/>
```

### Markdown 转换
- 输入：使用 `tryParseMarkdownToBlocks()` 将 Markdown 转换为 BlockNote 块
- 输出：使用 `blocksToMarkdownLossy()` 将块转换为 Markdown

## 测试建议

### 功能测试
1. ✅ 访问 `/publish` 页面
2. ✅ 测试各种 Markdown 格式：
   - 标题（H1-H6）
   - 粗体、斜体
   - 代码块（带语法高亮）
   - 引用
   - 有序/无序列表
   - 链接和图片
3. ✅ 测试提交功能
4. ✅ 验证提交的内容是否正确转换为 Markdown

### 用户体验测试
1. 测试工具栏按钮是否正常工作
2. 测试斜杠命令（/）
3. 测试拖拽排序
4. 测试快捷键（Ctrl+B 粗体等）
5. 测试响应式布局

## 已知问题

### 构建问题
当前 `pnpm build` 失败是由于网络问题导致 Google Fonts 无法加载，与 BlockNote 集成无关。这是一个独立的问题，需要：
- 配置字体本地化
- 或使用 CDN 镜像
- 或禁用 Google Fonts

### 解决方案
可以在 `next.config.mjs` 中配置字体优化：
```javascript
experimental: {
  optimizeFonts: false
}
```

或者使用本地字体文件。

## 下一步

1. 解决 Google Fonts 加载问题
2. 测试编辑器在生产环境的表现
3. 根据用户反馈调整编辑器样式
4. 考虑添加图片上传功能
5. 考虑添加自动保存草稿功能

