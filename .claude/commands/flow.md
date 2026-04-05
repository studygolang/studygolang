---
description: 'StudyGolang 标准化交付工作流：需求分析→方案设计→专家评审→编码实现→实现Review→验收。方案需3+专家≥9分通过，实现需3个不同模型专家全通过。'
allowed-tools: 'Agent, Read, Write, Edit, Bash, Glob, Grep, TaskCreate, TaskUpdate, TaskList, TaskGet, AskUserQuestion'
---

# /flow — StudyGolang 标准化交付工作流

你是 **StudyGolang 交付流程编排者**，严格按以下五阶段执行。每个阶段的门禁条件必须满足才能进入下一阶段。用 TaskCreate/TaskUpdate 跟踪进度。

**需求:** $ARGUMENTS

---

## 阶段 1：需求分析 & 方案设计

派发 **studygolang-refactor** agent（架构师角色），完成：

1. **需求理解** — 澄清模糊点，确认边界和约束
2. **代码调研** — 分析现有代码结构、可复用模块、API 缺口
3. **方案设计** — 输出结构化设计文档，包含：
   - 需求摘要
   - 涉及模块（后端 API / 前端页面 / 数据库变更）
   - API 设计（路径、方法、请求/响应结构）
   - 前端页面设计（路由、组件结构、数据获取方式 SSR/ISR）
   - 文件变更清单（新增/修改的文件列表）
   - 风险点和注意事项

**门禁：** 设计文档完成后，进入阶段 2。

---

## 阶段 2：方案设计评审（循环直到通过）

**至少 3 个资深专家 agent 并行评审**，每个独立打分（0-10）并给出意见：

### 专家角色分配（按需求类型选择 3+ 个）

| 专家 | 评审视角 | 适用场景 |
|------|---------|---------|
| studygolang-go-api | 后端架构合理性、API 设计规范、与现有 app/ 层一致性 | 涉及后端变更 |
| studygolang-nextjs-ssr | 前端架构、SSR/ISR 策略、组件设计、SEO 友好性 | 涉及前端变更 |
| studygolang-security | 认证鉴权、输入验证、XSS/CSRF、数据泄露风险 | 涉及用户数据/认证 |
| studygolang-seo-audit | URL 结构、meta 标签、结构化数据、渲染策略 | 涉及内容页面 |
| studygolang-go-tester | 可测试性、测试覆盖策略、边界条件 | 所有变更 |
| studygolang-refactor | 架构一致性、与现有代码风格匹配、复用度 | 所有变更 |

### 评审打分维度

- **完整性**（0-3）：是否覆盖所有需求场景，有无遗漏
- **架构一致性**（0-3）：是否符合项目现有架构和规范
- **可行性**（0-2）：技术方案是否可行，有无潜在阻塞
- **风险控制**（0-2）：风险点是否识别并给出应对方案

### 通过标准

- **所有专家评分均 ≥ 9 分** → 通过，进入阶段 3
- **任一专家评分 < 9 分** → 汇总反馈意见，打回阶段 1 迭代方案
- 迭代后重新评审，直到全部通过

### 评审执行方式

使用 Agent tool 并行派发 3+ 个专家，每个专家：
1. 读取设计文档
2. 从自身专业角度评审
3. 输出：评分 + 各维度得分 + 具体意见 + 改进建议

---

## 阶段 3：编码实现（方案通过后）

根据设计文档，派发合适的专家 agent 执行编码：

### 派发策略

```
后端 API + 前端页面可并行：
  ┌──────────────────────────┐
  │ studygolang-go-api       │  ← 后端 API 实现
  │ studygolang-nextjs-ssr   │  ← 前端页面实现
  └──────────────────────────┘

实现完成后追加（可并行）：
  ┌──────────────────────────┐
  │ studygolang-go-tester    │  ← 单元测试 + 集成测试
  │ studygolang-doc-writer   │  ← 文档更新
  └──────────────────────────┘
```

### 实现要求

- 严格按设计文档编码，不得偏离方案
- 遵循项目 CLAUDE.md 中的规范
- 每个实现 agent 完成后汇报变更文件清单

---

## 阶段 4：实现 Review（循环直到全部通过）

**3 个使用不同模型的资深 agent 并行 Review**，确保：
- 实现严格按照设计方案执行
- 符合项目架构、代码风格和规范
- 无 bug，无性能问题
- 不影响已有功能，不引入新 bug

### Review 执行

使用 Agent tool 并行派发 3 个 Review agent，**必须使用不同模型**：

```
Agent 1: model=opus   → 深度推理 Review（逻辑正确性、边界条件、竞态风险）
Agent 2: model=sonnet → 架构规范 Review（代码风格、项目一致性、设计模式）
Agent 3: model=haiku  → 快速全面 Review（构建通过、无遗漏、无明显 bug）
```

### 每个 Reviewer 检查项

1. **方案一致性**：实现是否与设计文档完全一致
2. **代码质量**：命名、结构、可读性、无重复代码
3. **项目规范**：API 响应格式、SSR/ISR 策略、SEO 规范
4. **功能正确性**：逻辑正确、边界处理、错误处理
5. **回归风险**：是否修改了公共逻辑、是否影响其他模块
6. **构建验证**：
   ```bash
   cd /Users/polarisxu/project/golang/studygolang && go build ./...
   cd /Users/polarisxu/project/golang/studygolang/frontend && pnpm build
   ```

### 通过标准

- **3 个 Reviewer 全部 PASS** → 进入阶段 5 验收
- **任一 Reviewer FAIL** → 汇总所有 Review 反馈，打回阶段 3 迭代修复
- 迭代后重新 Review，直到全部通过
- 每个 Reviewer 的反馈必须具体到文件和行号

### 迭代规则

- 打回时附带所有 Reviewer 的反馈（不仅仅是 FAIL 的那个）
- 修复时不能引入 Reviewer 已通过部分的新问题
- 如果迭代超过 3 轮，暂停并报告用户决策

---

## 阶段 5：验收交付

所有 Review 通过后：

1. 汇总变更清单（新增/修改文件、API 端点、页面路由）
2. 确认构建通过
3. 输出交付报告
4. 请求用户确认验收

---

## Agent 团队

| Agent | 职责 | 阶段 |
|-------|------|------|
| studygolang-refactor | 需求分析、方案设计、架构评审 | 1, 2 |
| studygolang-go-api | 后端 API 开发 | 3 |
| studygolang-nextjs-ssr | 前端 SSR 开发 | 3 |
| studygolang-go-tester | 测试编写 | 3 |
| studygolang-seo-audit | SEO 审计 | 2, 4 |
| studygolang-security | 安全审计 | 2, 4 |
| studygolang-doc-writer | 文档更新 | 3 |

## 项目关键路径

```
后端 API:  internal/http/controller/api/     (新增)
前端页面:  frontend/app/                      (App Router)
已有 API:  internal/http/controller/app/      (可复用)
业务逻辑:  internal/logic/                    (不改)
```

## 技术约束

- URL 结构保持不变（SEO）
- 列表页 SSR (no-store)，详情页 ISR (revalidate: 60)
- 管理后台不迁移
- 所有页面必须有 generateMetadata
- API 响应: `{ "code": 0, "data": {}, "message": "success" }`

---

**现在开始阶段 1：派发 studygolang-refactor 进行需求分析 & 方案设计。**
