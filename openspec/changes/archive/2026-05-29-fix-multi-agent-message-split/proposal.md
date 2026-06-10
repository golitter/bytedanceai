## Why

Orchestrator 运行时，所有 Agent 输出合并为一条 MySQL Message，agent_type 固定为初始值。页面刷新/历史加载后，所有内容混在一个气泡下，子 Agent 消息无法区分。同时 ExecutionEngine 为子 Agent 单独调用 RunTask 产生重复消息，导致前端显示 `orchestrator: abc nih ok` + `claudecode: nih` 的错乱。

## What Changes

- Backend StreamWriter 解析 SSE TEXT 事件中的 `agent_type`/`agent` 字段，检测 Agent 切换时 finalize 当前消息并在同 session 下新建 Message
- Backend `ServeStream` 重放（Phase 1 + serveCompleted）时携带消息自身的 AgentType/AgentName 元数据
- Backend `ListMessages` 新增可选 `session_id` 过滤参数，消除 ExecutionEngine 双写导致的跨 session 重复
- Frontend `getTaskMessages` 传入 `session_id`，只加载当前 session 的消息

## Capabilities

### New Capabilities
- `agent-switch-message-split`: StreamWriter 检测 agent_type 变化时自动创建新 Message，实现多 Agent 输出按 Agent 拆分为独立数据库记录

### Modified Capabilities
- `message-persistence`: Message 持久化从"单条合并"变为"按 Agent 拆分"，单次 RunTask 可产生多条 Message
- `chat-streaming`: SSE 重放时携带 agent 元数据（agent_type + agent_name）

## Impact

- **Backend**: `stream/writer.go`（核心改动）、`handler/stream.go`（重放逻辑）、`handler/message.go`（过滤参数）
- **Frontend**: `use-chat-stream.ts`（传 session_id）
- **数据库**: 无 schema 变更，同 session 下新增 Message 行
- **API**: `ListMessages` 新增 `session_id` query param（向后兼容）
- **Agentend**: 无改动
