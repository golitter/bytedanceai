## Why

当前 Agent 流式输出只在 SSE 流**完全结束后**才写入 MySQL，流式内容仅存在于 React 组件的 `useReducer` 内存中。用户切 Tab、刷新页面或网络断开时，正在输出的内容全部丢失。对于 IM 场景，Agent 一旦激活，其输出必须完整持久化，前端重连后应无缝恢复。

## What Changes

- **BREAKING**: `POST /api/tasks/:taskId/run` 从同步 SSE 代理改为异步——立即返回 `202 Accepted` + `message_id`，后台 goroutine 独立消费 agentend 流
- 新增 `GET /api/tasks/:taskId/stream` SSE 端点：支持前端（重）连接，先发 MySQL 已积累内容，再续接 Redis Stream 实时事件
- Backend 后台 goroutine 将每个 SSE event 发布到 Redis Stream（`agent:{session_id}:{message_id}`），同时批量 flush 内容到 MySQL（每 500 字符或每 2 秒）
- MySQL Message 表新增 `message_id`（UUID）、`status`（streaming/completed/failed）、`last_seq`（Redis Stream ID）字段
- 前端用 Zustand 全局 store 替代 `useReducer`，切 Tab 不丢失状态
- 前端 SSE 改用 GET 请求 + 支持断线自动重连（同一 `message_id`）
- Agentend 无改动

## Capabilities

### New Capabilities
- `stream-incremental-persist`: Backend 端增量持久化——goroutine 消费 agentend 流，实时发布到 Redis Stream，批量 flush 到 MySQL。含 BatchWriter 策略（500 字符 / 2 秒 / 流结束）
- `stream-reconnect`: Backend GET /stream SSE 端点 + 前端断线重连。重连时先发 MySQL 已存内容（伪装成 text events），再从 Redis Stream XREAD BLOCK 续接实时事件，前端无感知

### Modified Capabilities
- `message-persistence`: Agent message 从"流结束后一次性写入"改为"流开始时创建 status=streaming，增量 UPDATE，流结束设 completed"
- `chat-streaming`: 状态管理从 `useReducer`（组件级）迁移到 Zustand 全局 store（页面级生命周期），支持多 Tab 共享状态 + 刷新后自动重连

## Impact

- **Backend**: `handler/task.go`、`handler/stream.go`（新增）、`model/message.go`、新增 `stream/` 包、`go.mod` 新增 `go-redis` 依赖
- **Frontend**: `stores/chat.ts`（新增）、`hooks/use-chat-stream.ts`、`lib/api.ts`、`lib/sse.ts`
- **Infrastructure**: 需要部署 Redis 实例
- **Agentend**: 无改动
- **API**: `POST /run` 响应格式变更（从 SSE 流改为 202 JSON），新增 `GET /stream` 端点
