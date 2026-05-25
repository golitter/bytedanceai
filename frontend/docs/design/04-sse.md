# SSE — 连接与数据流

## 实现了什么

基于 EventSource 的 SSE 客户端，配合两步式 API 调用（POST 提交消息 + GET 连接 SSE 流）实现实时流式通信。所有 REST API 调用集中在 `lib/api.ts`。

## 怎么实现的

### SSE 客户端 (`src/lib/sse.ts`)

封装 `EventSource` 连接，接收 SSE 事件并解析为 `StreamEvent` 类型。支持自动重连和手动中断：

```typescript
interface SSEOptions {
  url: string
  params?: Record<string, string>
  onEvent: (event: StreamEvent) => void
  onError?: (error: Error) => void
  reconnect?: boolean
}

export function connectSSE({
  url, params, onEvent, onError, reconnect = false,
}: SSEOptions): AbortController {
  const controller = new AbortController()

  const qs = params ? '?' + new URLSearchParams(params).toString() : ''
  // 开发环境绕过 Vite 代理（Vite 会缓冲 SSE 响应）
  const baseUrl = import.meta.env.DEV ? 'http://localhost:8080' : ''
  const fullUrl = `${baseUrl}${url}${qs}`

  const es = new EventSource(fullUrl)

  es.onmessage = (e: MessageEvent) => {
    const data = typeof e.data === 'string' ? e.data : ''
    if (!data.trim()) return
    try {
      const event: StreamEvent = JSON.parse(data)
      onEvent(event)
    } catch {
      console.warn('Failed to parse SSE event:', data)
    }
  }

  es.onerror = () => {
    if (es.readyState === EventSource.CLOSED) {
      if (!controller.signal.aborted) {
        onError?.(new Error('SSE connection closed'))
      }
      return
    }
    if (!reconnect) {
      es.close()
      if (!controller.signal.aborted) {
        onError?.(new Error('SSE connection error'))
      }
    }
    // reconnect=true 时 EventSource 自动重连
  }

  controller.signal.addEventListener('abort', () => {
    es.close()
  })

  return controller
}
```

关键设计点：
- 开发环境直接连接 `http://localhost:8080` 绕过 Vite 代理（Vite 会缓冲 SSE 响应）
- `AbortController` 支持手动中断流，abort 时关闭 EventSource
- `reconnect` 参数控制是否让 EventSource 自动重连

### 两步式 SSE 流程

SSE 通信分为两步：先 POST 提交消息获取 `message_id`，再用该 ID 连接 SSE 流：

```
客户端                              后端
  │                                  │
  │  POST /api/tasks/:id/run         │  提交消息
  │  { message, session_id,          │
  │    agent_type }                  │
  │ ────────────────────────────────►│
  │  ◄── { message_id }             │
  │                                  │
  │  GET /api/tasks/:id/stream       │  连接 SSE 流
  │  ?session_id=&message_id=        │
  │ ────────────────────────────────►│
  │                                  │
  │  SSE: {"type":"init"}            │
  │ ◄────────────────────────────────│
  │  SSE: {"type":"text","content":{"text":"..."}} │
  │ ◄────────────────────────────────│
  │  SSE: {"type":"tool_call","content":{"name":"..."}} │
  │ ◄────────────────────────────────│
  │  SSE: {"type":"done"}            │
  │ ◄────────────────────────────────│
```

### API 层 (`src/lib/api.ts`)

所有 REST API 调用集中在 `lib/api.ts`，使用原生 `fetch` 封装。

**消息提交** — POST `/api/tasks/:id/run`：

```typescript
export async function submitMessage(
  taskId: string,
  body: { message: string; session_id: string; agent_type?: string },
): Promise<{ message_id: string; status: string }> {
  const res = await fetch(`${API_BASE}/tasks/${taskId}/run`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const json = await res.json().catch(() => ({}))
    throw new Error(json.msg || `HTTP ${res.status}`)
  }
  const json = await res.json()
  return json.data
}
```

**历史消息获取** — GET `/api/tasks/:id/messages`：

```typescript
export async function getTaskMessages(taskId: string): Promise<TaskMessage[]> {
  const res = await fetch(`${API_BASE}/tasks/${taskId}/messages`)
  const json = await res.json()
  return json.data
}
```

### 对话聚合

`fetchConversations()` 将 Task + Session 扁平化为 `Conversation` 视图，按 `lastActiveAt` 降序排列：

```typescript
export async function fetchConversations(): Promise<Conversation[]> {
  const tasks = await fetchTasks()
  const details = await Promise.all(tasks.map((t) => fetchTask(t.task_id)))
  const convos: Conversation[] = []
  for (const detail of details) {
    for (const s of detail.sessions) {
      convos.push({
        taskId: s.task_id,
        sessionId: s.session_id,
        agentType: s.agent_type,
        agentName: s.agent_name ?? '',
        title: s.agent_name || s.agent_type,
        lastActiveAt: s.updated_at,
        taskTitle: detail.task.title,
        status: s.status,
      })
    }
  }
  convos.sort((a, b) => new Date(b.lastActiveAt).getTime() - new Date(a.lastActiveAt).getTime())
  return convos
}
```

`createConversation()` 创建 Task -> 取首个 Session -> 返回 Conversation。

### API 接口总览

| 函数 | 方法 | 路径 | 说明 |
|------|------|------|------|
| `fetchTasks` | GET | `/api/tasks` | 获取任务列表 |
| `fetchTask` | GET | `/api/tasks/:id` | 获取任务详情（含 sessions） |
| `createTask` | POST | `/api/tasks` | 创建任务 |
| `submitMessage` | POST | `/api/tasks/:id/run` | 提交消息，返回 message_id |
| `getTaskMessages` | GET | `/api/tasks/:id/messages` | 获取任务消息列表 |
| `updateSession` | PUT | `/api/sessions/:id` | 更新 session（agent_name / avatar_url） |
| `fetchAgentTypes` | GET | `/api/agent-types` | 获取可用 Agent 类型列表 |
| `uploadAvatar` | POST | `/api/agents/avatar` | 上传头像 |
| `validateRepoPath` | POST | `/api/validate-repo-path` | 校验仓库路径 |
| `fetchConversations` | GET | 多接口聚合 | Task+Session 扁平化对话列表 |
| `createConversation` | POST+GET | 多接口组合 | 创建 Task -> 取 Session -> 返回 Conversation |
