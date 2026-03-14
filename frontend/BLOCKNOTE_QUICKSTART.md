# BlockNote 编辑器快速开始

## 🎉 已完成

发布页面的 Markdown 编辑器已升级为 **BlockNote** - 一个现代化的块编辑器。

## ✨ 主要特性

### 1. 所见即所得编辑
- 无需切换预览模式
- 实时看到格式化效果
- 类似 Notion 的编辑体验

### 2. 斜杠命令 (/)
输入 `/` 快速插入:
- 标题 (H1, H2, H3)
- 列表 (无序、有序、待办)
- 代码块
- 引用
- 表格
- 分割线

### 3. 拖拽排序
- 鼠标悬停在块左侧
- 出现拖拽手柄 (⋮⋮)
- 拖拽到新位置

### 4. 格式化工具栏
选中文本后自动显示:
- **粗体** (Ctrl/Cmd + B)
- *斜体* (Ctrl/Cmd + I)
- `代码` (Ctrl/Cmd + E)
- 链接 (Ctrl/Cmd + K)
- 颜色和背景色

### 5. Markdown 支持
自动识别 Markdown 语法:
```
# 标题 1
## 标题 2
- 列表项
1. 有序列表
> 引用
```代码块```
```

## 🚀 快速测试

1. 启动开发服务器:
```bash
cd frontend
pnpm dev
```

2. 访问发布页面:
```
http://localhost:3000/publish
```

3. 尝试以下操作:
   - 输入 `/` 查看命令菜单
   - 选中文本使用工具栏
   - 拖拽块重新排序
   - 插入表格并编辑
   - 使用 Markdown 语法

## 📝 使用技巧

### 快捷键
- `Ctrl/Cmd + B` - 粗体
- `Ctrl/Cmd + I` - 斜体
- `Ctrl/Cmd + E` - 行内代码
- `Ctrl/Cmd + K` - 插入链接
- `Ctrl/Cmd + Enter` - 提交表单
- `Enter` - 新建块
- `Backspace` - 删除空块

### Markdown 快捷输入
- `# ` + 空格 → 一级标题
- `## ` + 空格 → 二级标题
- `- ` + 空格 → 无序列表
- `1. ` + 空格 → 有序列表
- `[ ] ` + 空格 → 待办事项
- `> ` + 空格 → 引用
- ` ``` ` + 空格 → 代码块

### 表格操作
1. 输入 `/` 选择 "表格"
2. 点击单元格编辑内容
3. 右键菜单添加/删除行列
4. 拖拽列边界调整宽度

## 🎨 样式定制

编辑器样式文件: `frontend/app/publish/blocknote.css`

可以自定义:
- 工具栏样式
- 代码块样式
- 表格样式
- 深色模式
- 响应式布局

## 🔧 技术细节

### 组件位置
- 编辑器组件: `frontend/components/blocknote-editor.tsx`
- 发布页面: `frontend/app/publish/page.tsx`
- 样式文件: `frontend/app/publish/blocknote.css`

### 依赖包
```json
{
  "@blocknote/core": "^0.47.1",
  "@blocknote/mantine": "^0.47.1",
  "@blocknote/react": "^0.47.1"
}
```

### 数据流
```
用户输入 → BlockNote → Markdown → content state → 后端 API
```

## 📚 更多资源

- [BlockNote 官方文档](https://www.blocknotejs.org/)
- [BlockNote GitHub](https://github.com/TypeCellOS/BlockNote)
- [示例和演示](https://www.blocknotejs.org/examples)

## ❓ 常见问题

**Q: 如何插入图片?**
A: 输入 `/` 选择 "图片",或直接拖拽图片到编辑器

**Q: 支持哪些代码语言高亮?**
A: 支持常见编程语言,包括 Go, JavaScript, Python, Java 等

**Q: 可以导入现有 Markdown 文件吗?**
A: 可以,编辑器会自动解析 Markdown 内容

**Q: 如何保存草稿?**
A: 点击 "保存草稿" 按钮 (功能待实现)

## 🎯 下一步

- [ ] 测试所有编辑功能
- [ ] 实现图片上传
- [ ] 添加自动保存草稿
- [ ] 优化移动端体验
- [ ] 添加快捷键帮助面板

---

**享受全新的编辑体验!** 🎊
