# 数据流与状态管理

## 状态分层

```
┌─────────────────────────────────────────────────┐
│  Zustand Store (chat.ts)                        │
│  ─ 纯客户端导航状态                              │
│  ─ currentSessionId / currentTaskId             │
├─────────────────────────────────────────────────┤
│  React Query (use-conversations / use-sessions) │
│  ─ 服务端数据缓存 + 自动刷新                     │
│  ─ conversations / tasks / sessions             │
├─────────────────────────────────────────────────┤
│  useReducer (use-chat-stream)                   │
│  ─ 聊天流式状态机（组件级，不持久化）             │
│  ─ messages / streamingContent / status         │
└─────────────────────────────────────────────────┘
```

## 聊天状态机

`use-chat-stream.ts` 使用 `useReducer` 管理聊天流的全生命周期：

```
         sendMessage()
  idle ──────────────► loading
                          │
                    SSE init event
                          │
                          ▼
                      streaming ◄──── STREAM_TEXT
                       │  ▲
              STREAM_TOOL_CALL  STREAM_TOOL_RESULT
                       │  ▲
                       ▼  │
                    tool_running
                       │
            ┌──────────┼──────────┐
            │          │          │
        STREAM_DONE  error    abort()
            │          │          │
            ▼          ▼          ▼
          done       error      (idle)
```

### Action 说明

| Action | 触发条件 | 状态变更 |
|--------|----------|----------|
| `SEND_MESSAGE` | 用户发送消息 | → `loading`，追加 user 消息 |
| `STREAM_START` | SSE `init` 事件 | → `streaming`，记录 agentType |
| `STREAM_TEXT` | SSE `text` 事件 | 追加 `streamingContent` |
| `STREAM_TOOL_CALL` | SSE `tool_call` 事件 | → `tool_running`，记录 toolName |
| `STREAM_TOOL_RESULT` | SSE `tool_result` 事件 | → `streaming` |
| `STREAM_DONE` | SSE `done` 事件 | → `done`，将 streamingContent 转为 agent 消息 |
| `STREAM_ERROR` | SSE `error` 或网络错误 | → `error` |
| `RESET` | 手动重置 | → `idle` |

### State 结构

```typescript
interface ChatState {
  status: 'idle' | 'loading' | 'streaming' | 'tool_running' | 'done' | 'error'
  messages: ChatMessage[]
  streamingContent: string
  streamingAgentType?: AgentType
  error: Error | null
  toolName?: string
}
```

## SSE 连接

`lib/sse.ts` 封装了基于 `fetch` + `ReadableStream` 的 SSE 客户端：

```
客户端                        后端
  │                            │
  │  POST /api/tasks/:id/run   │
  │  { message, session_id,    │
  │    agent_type }            │
  │ ──────────────────────────►│
  │                            │
  │  SSE: data: {"type":"init"}
  │ ◄──────────────────────────│
  │  SSE: data: {"type":"text","content":{"text":"..."}}
  │ ◄──────────────────────────│
  │  SSE: data: {"type":"tool_call","content":{"name":"..."}}
  │ ◄──────────────────────────│
  │  SSE: data: {"type":"done"}
  │ ◄──────────────────────────│
```

实现要点：
- 使用 `ReadableStream` 逐块读取，不依赖 EventSource（支持 POST）
- 以 `\n\n` 分割 SSE 帧，解析 `data: ` 前缀的 JSON
- 返回 `AbortController`，支持中断流
- 丢弃不完整的帧（保留在 buffer 中等待下次拼接）

## API 层

`lib/api.ts` 封装了所有 REST API 调用：

### Task 相关

| 函数 | 方法 | 路径 | 说明 |
|------|------|------|------|
| `fetchTasks` | GET | `/api/tasks` | 获取任务列表 |
| `fetchTask` | GET | `/api/tasks/:id` | 获取任务详情（含 sessions） |
| `createTask` | POST | `/api/tasks` | 创建任务 |
| `deleteTask` | DELETE | `/api/tasks/:id` | 删除任务 |

### Session 相关

| 函数 | 方法 | 路径 | 说明 |
|------|------|------|------|
| `patchSession` | PATCH | `/api/sessions/:id` | 更新 session 状态 |

### Agent 相关

| 函数 | 方法 | 路径 | 说明 |
|------|------|------|------|
| `fetchAgentTypes` | GET | `/api/agent-types` | 获取可用 Agent 类型列表 |

### Conversation 视图

`Conversation` 是 Task + Session 的扁平化视图：

```typescript
interface Conversation {
  taskId: string
  sessionId: string
  agentType: AgentType
  title: string
  lastActiveAt: string
  taskTitle: string
  status: string
}
```

`fetchConversations()` 聚合所有 Task → Session，按 `lastActiveAt` 降序排列。
`createConversation()` 创建 Task → 取首个 Session → 返回 Conversation。

## React Query Hooks

### useConversations / useCreateConversation

```typescript
// 查询对话列表
useQuery({ queryKey: ['conversations'], queryFn: fetchConversations })

// 新建对话（成功后自动 invalidate）
useMutation({ mutationFn: createConversation, onSuccess: invalidate ['conversations'] })
```

### useSessions / useCreateTask / useDeactivateSession

```typescript
// 查询 Task 详情，提取 sessions
useQuery({ queryKey: ['task', taskId], queryFn: fetchTask })

// 创建/删除任务，成功后 invalidate ['tasks']
// 停用 session，成功后 invalidate ['task', taskId]
```

## Zustand Store

`stores/chat.ts` — 聊天导航状态（纯客户端）：

```typescript
{
  currentSessionId: string | null   // 当前选中的对话
  currentTaskId: string | null      // 当前关联的任务
  setCurrentSession: (id) => void
  setCurrentTask: (id) => void
  clearNavigation: () => void
}
```

`stores/app.ts` — 全局 store 占位，当前未使用。
