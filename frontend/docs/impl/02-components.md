# 组件目录

## 页面层

### ImPage (`pages/ImPage.tsx`)

主页面，双栏布局入口。左侧 `ConversationList`，右侧 `ChatArea` 或空状态占位。

- 从 `useConversations()` 获取对话列表
- 从 `useChatNav()` 读取当前选中的 `currentSessionId`
- 选中时渲染 `<ChatArea>`，未选中时显示 "选择一个对话开始聊天" 空状态

---

## IM 侧栏 (`components/im/`)

### ConversationList

侧栏容器组件，固定宽度 280px，包含四部分：

1. **Header**：品牌标识 "AgentHub" + 新建对话按钮（`+`）
2. **搜索栏**：过滤 `agentType` 和 `taskTitle`
3. **对话列表**：渲染 `ConversationItem`，支持加载态和空态
4. **新建弹窗**：`NewChatDialog` 受控组件

数据源：`useConversations()` (React Query) + `useChatNav()` (Zustand)

### ConversationItem

单条对话项，展示：

- 左侧：`AgentAvatar`（带状态指示灯）
- 中间：Agent 名称 + 任务标题（单行截断）
- 右侧：相对时间（"刚刚"、"5分钟前"、"2小时前"等）

交互：
- 点击 → `setCurrentSession(sessionId)`
- Active 态：左侧品牌色竖线 + hover 背景色

### NewChatDialog

新建对话弹窗（shadcn Dialog），展示可用 Agent 列表：

- 通过 `useQuery(['agent-types'])` 从 API 拉取可用 Agent，失败时 fallback 到内置列表
- 每个 Agent 显示：`AgentAvatar` + 名称 + 描述
- 点击 → `createConversation` mutation → 成功后自动选中并关闭弹窗

---

## 聊天区 (`components/chat/`)

### ChatArea

聊天主容器，纵向三段式布局：

1. **Header**（h-12）：`AgentAvatar` + Agent 显示名 + "正在回复..." 状态
2. **消息区**：空态显示 Agent 大头像 + "发送消息开始对话"；有消息时渲染 `MessageList`
3. **输入区**：`MessageInput`

核心 hook：`useChatStream(taskId)` 返回 `{ state, sendMessage, abort }`

状态判断：`['loading', 'streaming', 'tool_running'].includes(state.status)` 为活跃态

### MessageList

消息列表组件，支持两种模式：

- **普通模式**（消息 ≤ 50 条）：直接渲染
- **虚拟滚动**（消息 > 50 条）：`@tanstack/react-virtual` 优化

特性：
- 流式消息以 `id: 'streaming'` 临时追加到列表末尾
- 自动滚动到底部（`autoScroll` 状态跟踪）
- 手动上滑时隐藏自动滚动，显示 "回到底部" 按钮
- 内部 `MessageRenderer` 根据 `msg.role` 选择气泡变体

### MessageBubble

消息气泡组件，三种变体（discriminated union）：

| 变体 | 布局 | 样式 |
|------|------|------|
| `user` | 右对齐 | 半透明品牌色背景 (`rgba(99,102,241,0.08)`) |
| `agent` | 左对齐 + AgentAvatar | 卡片背景 + 左侧 3px Agent 色竖线 |
| `system` | 居中 | 小字灰色文本 |

Agent 气泡流式输出时显示闪烁光标 `▌`。

### MessageInput

输入框组件：

- `textarea` 自动高度（最小 48px，最大 200px）
- `Enter` 发送，`Shift+Enter` 换行
- 发送后清空并重置高度
- 右侧品牌色发送按钮（Lucide `Send` 图标）
- `disabled` 状态：流式输出时禁用

### AgentAvatar

Agent 头像组件：

- 圆角方块，显示 Agent 首字母
- 颜色映射：`claude-code` → 紫色 `#6366F1`，`opencode` → 琥珀色 `#F59E0B`，`orchestrator` → 绿色 `#22C55E`
- 状态指示灯（右下角小圆点）：`ready` 绿、`running` 琥珀闪烁、`offline` 灰、`error` 红
- 可配置 `size`（默认 32px）

---

## Markdown 渲染 (`components/markdown/`)

### MarkdownRenderer

基于 `react-markdown` + `remark-gfm` 的 Markdown 渲染器，自定义组件：

- **代码块**（`code` + `className` 含 language-）：委托给 `CodeBlock`
- **行内代码**：深色背景 + Geist Mono 字体
- **表格**：带边框和表头背景色
- 输出暗色主题 prose 样式

### CodeBlock

代码高亮组件：

- 使用 Shiki（`tokyo-night` 主题）异步高亮
- 高亮完成前 fallback 到纯文本 + 行号显示
- 代码块可横向滚动
- 字体：`Geist Mono`，字号 13px
