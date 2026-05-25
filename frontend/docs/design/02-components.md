# Components — 组件体系

## 实现了什么

基于三层组件模型（Page / Smart / Dumb）构建的 IM 聊天 UI，分为 IM 侧栏、聊天区、Markdown 渲染三大模块。所有组件使用 Tailwind CSS + CSS 变量驱动样式，无硬编码颜色值。

## 怎么实现的

## IM 侧栏 (`components/im/`)

### ConversationList (`src/components/im/ConversationList.tsx`)

侧栏容器组件，固定宽度 280px，包含 Header、搜索栏、对话列表和新建弹窗四部分。数据源为 `useConversations()` (React Query) + `useChatNav()` (Zustand)：

```tsx
export function ConversationList() {
  const [search, setSearch] = useState('')
  const [showNewChat, setShowNewChat] = useState(false)
  const { data: conversations, isLoading } = useConversations()
  const { currentSessionId, setCurrentSession } = useChatNav()

  const filtered = conversations?.filter((c) => {
    if (!search) return true
    const q = search.toLowerCase()
    return c.agentType.includes(q) || c.taskTitle.toLowerCase().includes(q)
  })
  // ... 渲染 Header / Search / ConversationItem 列表 / NewChatDialog
}
```

搜索支持按 `agentType` 和 `taskTitle` 过滤。空态和加载态分别显示对应占位 UI。

### ConversationItem (`src/components/im/ConversationItem.tsx`)

单条对话项，接收 `Conversation` 数据并渲染 Agent 头像、名称、任务标题和相对时间。通过 `useHoverStyle` hook 实现悬停效果：

```tsx
export function ConversationItem({ conversation, isActive, onClick }: ConversationItemProps) {
  const name =
    conversation.agentName || AGENT_NAMES[conversation.agentType] || conversation.agentType
  const hoverStyle = useHoverStyle()

  return (
    <button
      className="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left transition-colors duration-120 ease-out"
      style={{
        backgroundColor: isActive ? 'var(--accent)' : 'transparent',
        borderLeft: isActive ? '2px solid var(--primary)' : '2px solid transparent',
      }}
      onClick={onClick}
      {...(!isActive && hoverStyle)}
    >
```

Active 态通过 `borderLeft: 2px solid var(--primary)` 品牌色竖线标识，非 Active 态使用 `useHoverStyle` 提供的 `onMouseEnter/Leave` 切换背景色。

### NewChatDialog (`src/components/im/NewChatDialog.tsx`)

新建对话弹窗（shadcn Dialog），流程为：输入仓库路径 -> 校验 -> 选择 Agent 类型 -> 创建对话。仓库路径通过 `validateRepoPath()` API 校验后才能选择 Agent：

```tsx
const handleValidate = async () => {
  const path = repoPathRef.current?.value?.trim()
  if (!path) {
    setRepoPathError('请输入仓库路径')
    setRepoPathValidated(false)
    return
  }
  setValidating(true)
  try {
    const result = await validateRepoPath(path)
    if (result.valid) {
      setRepoPathValidated(true)
    } else {
      setRepoPathError(result.errors.join('; '))
    }
  } finally {
    setValidating(false)
  }
}
```

可用 Agent 列表通过 `useQuery({ queryKey: ['agent-types'], queryFn: fetchAgentTypes })` 拉取，失败时 fallback 到内置列表 `['claude-code', 'opencode', 'orchestrator']`。选中后调用 `createConversation` mutation，成功后自动选中并关闭弹窗。

---

## 聊天区 (`components/chat/`)

### ChatArea (`src/components/chat/ChatArea.tsx`)

聊天主容器（Smart 组件），纵向三段式布局：Header（h-12）、消息区、输入区。核心 hook `useChatStream` 返回 `{ state, sendMessage }`：

```tsx
export function ChatArea({ taskId, sessionId, agentType = 'claude-code', agentName, avatarUrl, repoPath }: ChatAreaProps) {
  const { state, sendMessage } = useChatStream(taskId, sessionId)
  const isStreaming = ['loading', 'streaming', 'tool_running'].includes(state.status)
```

发送消息前会先验证 `repoPath`（如果存在），通过 `validateRepoPath()` API 校验路径有效性。Header 区域显示 `AgentAvatar` + Agent 显示名 + "正在回复..." 状态。点击 Header 中的头像或名称可打开 `AgentEditDialog`。

### MessageList (`src/components/chat/MessageList.tsx`)

消息列表组件，支持两种渲染模式，阈值为 50 条消息：

```tsx
const VIRTUALIZE_THRESHOLD = 50

const allMessages =
  isStreaming && streamingContent
    ? [
        ...messages,
        {
          id: 'streaming',
          role: 'agent' as const,
          content: streamingContent,
          agentType: streamingAgentType as AgentType | undefined,
          timestamp: Date.now(),
        },
      ]
    : messages

const useVirtual = allMessages.length > VIRTUALIZE_THRESHOLD
```

流式消息以 `id: 'streaming'` 临时追加到列表末尾。虚拟滚动使用 `@tanstack/react-virtual`，通过 `estimateSize` 根据内容长度估算行高。内置自动滚动逻辑：监听 `scrollHeight - scrollTop - clientHeight < 60` 判断是否在底部，手动上滑时隐藏自动滚动并显示 "回到底部" 按钮。

### MessageBubble (`src/components/chat/MessageBubble.tsx`)

消息气泡组件，使用 TypeScript discriminated union 定义三种变体：

```tsx
interface UserBubbleProps extends BaseProps {
  variant: 'user'
}
interface AgentBubbleProps extends BaseProps {
  variant: 'agent'
  agentType: AgentType
  avatarUrl?: string
  agentName?: string
  status?: 'ready' | 'running' | 'offline' | 'error'
  isStreaming?: boolean
}
interface SystemBubbleProps extends BaseProps {
  variant: 'system'
}
type MessageBubbleProps = UserBubbleProps | AgentBubbleProps | SystemBubbleProps
```

- **user**：右对齐，`bg-primary-soft` 背景 + `border-primary-border` 边框
- **agent**：左对齐 + AgentAvatar，`bg-card` 背景 + 左侧 3px Agent 色竖线，流式输出时显示闪烁光标 `▌`
- **system**：居中，小字 `text-muted-foreground`

### MessageInput (`src/components/chat/MessageInput.tsx`)

输入框组件，textarea 自动高度（最小 48px，最大 200px），`Enter` 发送，`Shift+Enter` 换行：

```tsx
const adjustHeight = useCallback(() => {
  const el = textareaRef.current
  if (!el) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 200)}px`
}, [])
```

发送后清空输入框并重置高度为 48px。流式输出时 `disabled`，按钮和输入框同时变为不可用。

### AgentAvatar (`src/components/chat/AgentAvatar.tsx`)

Agent 头像组件，圆角方块显示首字母或上传的头像图片。颜色映射通过 CSS 变量：

```tsx
const AGENT_COLORS: Record<AgentType, string> = {
  'claude-code': 'var(--agent-claude)',
  opencode: 'var(--agent-opencode)',
  orchestrator: 'var(--agent-orchestrator)',
}
```

状态指示灯（右下角小圆点）使用 `STATUS_COLORS` 映射，`ready` 脉冲动画 `status-ready-pulse`，`running` 旋转动画 `status-running-spin`。支持自定义头像 URL，无自定义头像时若有 `agentName` 则通过 DiceBear API 生成 initials 头像。

### AgentEditDialog (`src/components/chat/AgentEditDialog.tsx`)

Agent 编辑弹窗，支持修改名称和上传头像。上传头像调用 `uploadAvatar` API，保存时调用 `updateSession` API 并 invalidate `conversations` 缓存：

```tsx
const handleSave = async () => {
  const data: { agent_name?: string; avatar_url?: string } = {}
  if (name !== initialName) data.agent_name = name
  if (avatarUrl !== initialAvatarUrl) data.avatar_url = avatarUrl
  if (Object.keys(data).length > 0) {
    await updateSession(sessionId, data)
    await queryClient.invalidateQueries({ queryKey: ['conversations'] })
  }
}
```

---

## Markdown 渲染 (`components/markdown/`)

### MarkdownRenderer (`src/components/markdown/MarkdownRenderer.tsx`)

基于 `react-markdown` + `remark-gfm` 的渲染器，通过自定义 `components` 对象覆盖默认渲染：

```tsx
const components: Components = {
  pre({ children }) {
    return <>{children}</>
  },
  code({ className, children, ...props }) {
    const match = /language-(\w+)/.exec(className ?? '')
    const code = String(children).replace(/\n$/, '')
    if (match) {
      return <CodeBlock code={code} language={match[1]} />
    }
    return (
      <code className="rounded bg-code px-1.5 py-0.5 text-[13px]"
        style={{ fontFamily: "'Geist Mono', monospace", letterSpacing: 0 }}
        {...props}>
        {children}
      </code>
    )
  },
  table({ children }) {
    return (
      <div className="overflow-x-auto">
        <table className="w-full border-collapse text-sm border-border">{children}</table>
      </div>
    )
  },
}
```

代码块（带 `language-` 前缀）委托给 `CodeBlock`，行内代码使用深色背景 + Geist Mono 字体，表格带边框和表头背景色。

### CodeBlock (`src/components/markdown/CodeBlock.tsx`)

代码高亮组件，使用 Shiki（`tokyo-night` 主题）异步高亮。高亮完成前 fallback 到纯文本 + 行号显示：

```tsx
useEffect(() => {
  async function highlight() {
    const shiki = await import('shiki')
    const highlighter = await shiki.createHighlighter({
      themes: ['tokyo-night'],
      langs: language ? [language] : [],
    })
    const result = highlighter.codeToHtml(code, {
      lang: language ?? 'text',
      theme: 'tokyo-night',
    })
    setHtml(result)
    highlighter.dispose()
  }
  if (language) highlight()
}, [code, language])
```

Shiki 通过动态 `import()` 按需加载，`highlighter` 使用后立即 `dispose()` 释放资源。代码块可横向滚动，字体 `Geist Mono`，字号 13px，行高 1.65。
