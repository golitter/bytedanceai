# State — 状态管理

## 实现了什么

三层状态架构：Zustand 管理聊天导航与各会话独立流式状态，TanStack React Query 管理服务端数据缓存，`useChatStream` hook 编排 SSE 连接与 store actions 的协作。Agent 消息通过 `reduceEventToBlocks` 解析为 `MessageBlock[]` 结构化块（text / html-render / image / attachment / diff / preview）。

## 怎么实现的

### Zustand Store (`src/stores/chat.ts`)

单一 store 管理导航状态 + 各会话独立聊天状态。导航状态存储 `currentSessionId`，会话状态以 `sessions[sessionId]` 的 Map 形式隔离：

```typescript
export interface ChatMessage {
  id: string
  dbId?: number
  role: 'user' | 'agent' | 'system'
  content: string
  blocks?: MessageBlock[]
  agentType?: AgentType
  timestamp: number
  messageId?: string
  status?: string
}

export type ChatStatus = 'idle' | 'loading' | 'streaming' | 'tool_running' | 'done' | 'error'

export interface ActiveStream {
  messageId: string
  sessionId: string
}

interface SessionChatState {
  messages: ChatMessage[]
  streamingContent: string
  streamingAgentType?: AgentType
  status: ChatStatus
  error: Error | null
  toolName?: string
  activeStream: ActiveStream | null
  hasMore: boolean
  isLoadingMore: boolean
}

interface ChatStoreState {
  nav: ChatNavState
  sessions: Record<string, SessionChatState>
  // Nav actions
  setCurrentSession: (sessionId: string) => void
  clearNavigation: () => void
  // Session actions
  getSession: (sessionId: string) => SessionChatState
  loadHistory: (sessionId: string, messages: ChatMessage[], hasMore?: boolean) => void
  sendMessage: (sessionId: string, message: ChatMessage, activeStream: ActiveStream) => void
  streamStart: (sessionId: string, agentType: AgentType) => void
  streamText: (sessionId: string, text: string) => void
  streamToolCall: (sessionId: string, toolName: string) => void
  streamToolResult: (sessionId: string) => void
  streamDone: (sessionId: string) => void
  streamError: (sessionId: string, error: Error) => void
  resetSession: (sessionId: string) => void
  // Pagination actions
  prependMessages: (sessionId: string, messages: ChatMessage[], hasMore: boolean) => void
  setLoadingMore: (sessionId: string, loading: boolean) => void
}
```

`loadHistory` 对历史 agent 消息同样执行 `reduceEventToBlocks` 解析：

```typescript
loadHistory: (sessionId, messages, hasMore) =>
  set((s) => ({
    sessions: {
      ...s.sessions,
      [sessionId]: {
        ...ensureSession(s, sessionId),
        status: 'done',
        messages: messages.map((msg) =>
          msg.role === 'agent' && msg.content
            ? { ...msg, blocks: reduceEventToBlocks(msg.content) }
            : msg,
        ),
        hasMore: hasMore ?? false,
      },
    },
  })),
```

`prependMessages` 用于向上翻页加载更多历史消息，将新消息插入到列表头部：

```typescript
prependMessages: (sessionId, messages, hasMore) =>
  set((s) => {
    const session = ensureSession(s, sessionId)
    const mapped = messages.map((msg) =>
      msg.role === 'agent' && msg.content
        ? { ...msg, blocks: reduceEventToBlocks(msg.content) }
        : msg,
    )
    return {
      sessions: {
        ...s.sessions,
        [sessionId]: {
          ...session,
          messages: [...mapped, ...session.messages],
          hasMore,
          isLoadingMore: false,
        },
      },
    }
  }),
```

`ensureSession` 工具函数确保访问不存在的 sessionId 时返回初始状态，而非 undefined：

```typescript
function ensureSession(state: ChatStoreState, sessionId: string): SessionChatState {
  return state.sessions[sessionId] ?? { ...initialSessionState }
}
```

### 聊天状态机流转

Store actions 驱动聊天流的全生命周期：

```
         sendMessage()
  idle ──────────────► loading
                          │
                    SSE init event → streamStart()
                          │
                          ▼
                      streaming ◄──── streamText()
                       │  ▲
          streamToolCall()  streamToolResult()
                       │  ▲
                       ▼  │
                    tool_running
                       │
            ┌──────────┼──────────┐
    streamDone()  streamError()  abort()
            │          │          │
            ▼          ▼          ▼
          done       error      (idle)
```

### Store Action 实现

`sendMessage` 追加用户消息，设置状态为 `loading`，记录 `activeStream`：

```typescript
sendMessage: (sessionId, message, activeStream) =>
  set((s) => ({
    sessions: {
      ...s.sessions,
      [sessionId]: {
        ...ensureSession(s, sessionId),
        status: 'loading',
        messages: [...(s.sessions[sessionId]?.messages ?? []), message],
        streamingContent: '',
        error: null,
        activeStream,
      },
    },
  })),
```

`streamText` 追加流式文本，若当前处于 `tool_running` 状态则回退到 `streaming`：

```typescript
streamText: (sessionId, text) =>
  set((s) => {
    const session = ensureSession(s, sessionId)
    return {
      sessions: {
        ...s.sessions,
        [sessionId]: {
          ...session,
          status: session.status === 'tool_running' ? 'streaming' : session.status,
          streamingContent: session.streamingContent + text,
        },
      },
    }
  }),
```

`streamDone` 将 `streamingContent` 通过 `reduceEventToBlocks` 解析为 `MessageBlock[]`，转为 agent 消息追加到 `messages`，清空流式状态：

```typescript
streamDone: (sessionId) =>
  set((s) => {
    const session = ensureSession(s, sessionId)
    const blocks = reduceEventToBlocks(session.streamingContent)
    const agentMessage: ChatMessage = {
      id: `agent-${Date.now()}`,
      role: 'agent',
      content: session.streamingContent,
      blocks,
      agentType: session.streamingAgentType,
      timestamp: Date.now(),
    }
    return {
      sessions: {
        ...s.sessions,
        [sessionId]: {
          ...session,
          status: 'done',
          messages: [...session.messages, agentMessage],
          streamingContent: '',
          streamingAgentType: undefined,
          activeStream: null,
        },
      },
    }
  }),
```

### useChatNav 选择器

暴露导航状态的选择器 hook，组件通过它订阅 `currentSessionId`，避免订阅整个 store 导致不必要的 re-render：

```typescript
export function useChatNav() {
  const currentSessionId = useChatStore((s) => s.nav.currentSessionId)
  const setCurrentSession = useChatStore((s) => s.setCurrentSession)
  const clearNavigation = useChatStore((s) => s.clearNavigation)
  return { currentSessionId, setCurrentSession, clearNavigation }
}
```

### React Query Hooks (`src/hooks/use-conversations.ts`)

对话列表使用 `useQuery` 查询，新建对话使用 `useMutation`，成功后自动 invalidate 缓存：

```typescript
export function useConversations() {
  return useQuery({
    queryKey: ['conversations'],
    queryFn: fetchConversations,
  })
}

export function useCreateConversation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ agentType, agentName, title, repoPath }: { ... }) =>
      createConversation(agentType, agentName, title, repoPath),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
    },
  })
}
```

### useChatStream Hook (`src/hooks/use-chat-stream.ts`)

核心编排 hook，连接 Zustand store 与 SSE 客户端。挂载时加载最近 20 条历史消息（cursor 分页），若发现 `status === 'streaming'` 的 agent 消息则自动重连 SSE。返回 `{ state, sendMessage, abort }`：

```typescript
export function useChatStream(taskId: string, sessionId: string) {
  const store = useChatStore()
  const abortRef = useRef<AbortController | null>(null)

  const connectToStream = useCallback((messageId: string) => {
    abortRef.current?.abort()
    store.streamStart(sessionId, 'claude-code')

    const controller = connectSSE({
      url: `/api/tasks/${taskId}/stream`,
      params: { session_id: sessionId, message_id: messageId },
      reconnect: true,
      onEvent: (event: StreamEvent) => {
        switch (event.type) {
          case EventTypeValues.Init:
            break
          case EventTypeValues.Text:
            store.streamText(sessionId, (event.content?.text as string) ?? '')
            break
          case EventTypeValues.ToolCall:
            store.streamToolCall(sessionId, (event.content?.name as string) ?? 'unknown')
            break
          case EventTypeValues.ToolResult:
            store.streamToolResult(sessionId)
            break
          case EventTypeValues.Done:
            store.streamDone(sessionId)
            break
          case EventTypeValues.Error:
            store.streamError(sessionId, new Error(...))
            break
        }
      },
      onError: (error) => store.streamError(sessionId, error),
    })
    abortRef.current = controller
  }, [taskId, sessionId])
```

发送消息时先追加 user 消息到 store，再调用 `submitMessage` API 获取 `message_id`，然后用该 ID 连接 SSE 流：

```typescript
const sendMessage = useCallback(async (message: string, agentType: AgentType = 'claude-code') => {
  const userMessage: ChatMessage = {
    id: `user-${Date.now()}`, role: 'user', content: message, timestamp: Date.now(),
  }
  store.sendMessage(sessionId, userMessage, { messageId: '', sessionId })

  const result = await submitMessage(taskId, {
    message, session_id: sessionId, agent_type: agentType,
  })
  connectToStream(result.message_id)
}, [taskId, sessionId, connectToStream])
```

历史消息加载（含 cursor 分页）和自动重连逻辑：

```typescript
useEffect(() => {
  let cancelled = false

  getTaskMessages(taskId, { limit: 20 })
    .then((res) => {
      if (cancelled || res.data.length === 0) return
      const chatMessages = res.data.map((m) => ({ ... }))
      store.loadHistory(sessionId, chatMessages, res.has_more)

      // 自动重连：如果发现 streaming 中的消息
      const streaming = res.data.find((m) => m.role === 'agent' && m.status === 'streaming')
      if (streaming && streaming.message_id) {
        connectToStream(streaming.message_id)
      }
    })
    .catch(() => { /* silently ignore */ })

  return () => { cancelled = true }
}, [taskId, sessionId, connectToStream])
```
