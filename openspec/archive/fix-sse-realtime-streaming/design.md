## Context

当前 IM 聊天系统的 Agent 回复采用 SSE（Server-Sent Events）流式推送架构：

```
前端 connectSSE (fetch + ReadableStream)
  → Vite dev proxy (/api → localhost:8080)
  → 后端 ServeStream handler (Gin)
  → Redis Stream XRead (实时事件)
  → StreamWriter goroutine (消费 agentend SSE)
```

问题：前端使用 `fetch` + `ReadableStream` 模拟 SSE 连接。Vite dev server 的 `http-proxy` 在处理 `fetch` 的流式响应时存在缓冲行为，导致 SSE 事件无法实时到达前端。刷新页面后通过 `getTaskMessages` 从 MySQL 加载历史消息才能看到完整内容。

后端还存在竞态条件：`serveStreaming` 的阻塞循环在首次检查 `stream.IsActive()` 时，后台 goroutine 可能还未注册 StreamWriter（正在等待 agent 服务响应），导致 handler 误判流失败。

## Goals / Non-Goals

**Goals:**
- Agent 回复通过 SSE 实时流式显示，无需刷新页面
- 用户发送的消息立即显示（乐观更新）
- 后端 serveStreaming 正确处理 goroutine 注册延迟
- 保持断线重连能力
- 保持现有 Zustand store + React 组件架构不变

**Non-Goals:**
- 不引入 WebSocket（SSE 已满足单向推送需求）
- 不修改 Redis Stream 架构
- 不修改 agentend（Python Agent 端）代码
- 不做跨用户实时推送（当前架构为单用户 chat）

## Decisions

### D1: 用 EventSource 替换 fetch + ReadableStream

**选择**: 使用浏览器原生 `EventSource` API

**替代方案**:
- (A) 配置 Vite proxy 跳过 SSE 响应的缓冲 — 不可靠，生产环境 nginx 也需配置
- (B) 前端轮询 `/api/tasks/:taskId/messages` — 延迟高，浪费请求
- (C) 保持 fetch 但加 `ReadableStream` polyfill — 治标不治本

**理由**: `EventSource` 是浏览器为 SSE 专门设计的 API，原生处理流式解析、自动重连、事件分发，能正确穿透各种代理和中间层。`fetch` 本质是为请求-响应模式设计的，流式读取是后加的能力，在某些代理环境下会被缓冲。

### D2: 后端 serveStreaming 竞态修复

**选择**: 当 `IsActive()` 返回 false 且消息状态为 "streaming" 时，不立即报错，而是 fall through 到 XRead 阻塞等待

**替代方案**:
- (A) 提前注册 StreamWriter（在 StreamAgent 调用前）— 需要处理注册后的错误清理
- (B) 在 serveStreaming 入口加 sleep 等待 — 不优雅

**理由**: fall through 方案最简洁，利用已有的 5 秒 XRead 阻塞等待，给 goroutine 足够时间注册。无需额外代码处理清理逻辑。

### D3: 前端乐观更新

**选择**: `store.sendMessage()` 在 `submitMessage()` API 调用之前执行

**理由**: 用户消息是本地已知的，无需等后端确认即可显示。API 失败时通过 `streamError` 回滚状态。

## Risks / Trade-offs

- **[EventSource 不支持自定义 header]** → 当前 SSE 连接不需要自定义 header（认证未在 SSE 层），无影响
- **[EventSource 自动重连可能与 store 状态冲突]** → 重连时 `onEvent` 仍通过 store 更新状态；后端 `serveStreaming` 从 MySQL 历史开始再接 Redis Stream，天然支持重连
- **[serveStreaming 竞态修复可能导致长连接等待]** → goroutine 有 30 分钟超时，且客户端断开时 context 会取消，不会无限等待
