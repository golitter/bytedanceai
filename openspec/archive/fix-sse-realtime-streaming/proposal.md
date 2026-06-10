## Why

Agent 回复不能实时显示在聊天界面。用户发送消息后，需要手动刷新页面才能看到 Agent 的完整回复。当前的 SSE 流式链路（前端 `fetch` + `ReadableStream` → Vite proxy → 后端 Gin `text/event-stream` → Redis Stream）在开发环境下无法将事件实时推送到前端，刷新后通过 MySQL 历史加载才能看到内容。

## What Changes

- **前端 SSE 连接层**：用浏览器原生 `EventSource` API 替换当前 `fetch` + `ReadableStream` 模拟方案，解决 Vite dev proxy 缓冲流式响应的问题
- **后端 serveStreaming 竞态修复**：当 `IsActive()` 返回 false 但消息状态仍为 "streaming" 时，不再误判为失败，改为等待 goroutine 注册完成
- **前端乐观更新**：用户消息先加入 Zustand store 立即显示，再异步发送 API 请求
- **SSE 连接健壮性增强**：完善错误日志、连接状态跟踪，便于排查链路问题

## Capabilities

### New Capabilities
- `sse-native-streaming`: 用原生 EventSource 替换 fetch 模拟的 SSE 客户端，确保流式事件实时到达前端

### Modified Capabilities

## Impact

- **前端文件**：`frontend/src/lib/sse.ts`（重写）、`frontend/src/hooks/use-chat-stream.ts`（乐观更新）
- **后端文件**：`backend/internal/handler/stream.go`（竞态修复）
- **依赖**：无新增外部依赖（EventSource 是浏览器内置 API）
- **兼容性**：EventSource 在所有现代浏览器中支持，无需 polyfill
