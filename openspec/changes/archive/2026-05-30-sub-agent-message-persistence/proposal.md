## Why

Orchestrator 分派子 Agent（claude-code、opencode）时，子 Agent 的执行在 agentend 内部通过 `adapter.stream_chat()` 直接调用完成，绕过了后端的 `StreamWriter`，导致子 Agent 的响应消息不会写入 MySQL。前端刷新后只能从数据库加载到 orchestrator 的消息，子 Agent 的消息丢失。

## What Changes

- Orchestrator 分派子 Agent 时，不再直接在 agentend 内部调用 adapter，而是回调后端 API `POST /api/tasks/:taskId/run`，由后端启动 StreamWriter 为每个子 Agent 创建独立的 Message 记录并持久化流式内容。
- 后端新增内部 API 或复用现有 RunTask 端点，支持 orchestrator（通过 agentend）发起的子 Agent run 请求。
- Agentend 的 ExecutionEngine 改为调用后端 HTTP API 而非直接调用 adapter。
- 前端无需改动 —— 刷新后 `getTaskMessages` 已返回 task 下所有 session 的消息。

## Capabilities

### New Capabilities
- `sub-agent-backend-dispatch`: 子 Agent 通过后端 API 执行，消息走 StreamWriter 持久化

### Modified Capabilities
- `message-persistence`: 新增场景 —— 子 Agent 的 agent message 也需通过 StreamWriter 持久化
- `orchestrator-planning`: OrchestratorAdapter 的执行阶段改为调用后端 API 而非内部 adapter
- `task-dispatcher`: DispatchResult 的消费方从 agentend 内部引擎变为后端 RunTask handler

## Impact

- **agentend**: `orchestrator/execution/engine.py` 重写执行逻辑，从直接调用 adapter 改为 HTTP 回调后端；新增后端 HTTP client 调用
- **backend**: `handler/task.go` 的 RunTask 需支持内部调用（区分前端请求 vs agentend 回调），可能需要内部认证机制
- **contracts**: 无变更 —— Message model 和 SSE event 格式不变
- **frontend**: 无变更
