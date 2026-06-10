## Why

流式输出体验严重滞后：Orchestrator 模式下子 Agent 的 token 经过四层网络跳转（agentend → backend StreamWriter → Redis → backend SSE → agentend BackendClient）后才被全量缓冲，执行完毕才向前端转发一条 `RUNTIME_COMPLETED` 事件，用户在子 Agent 执行期间看不到任何输出。单 Agent 模式也存在 Redis XREAD 5s 阻塞 + Count:1 逐条读取、每个 token 独占一次 Flush 等开销。

## What Changes

- Orchestrator ExecutionEngine 从"全量缓冲"改为"流式透传"：子 Agent 的每个 token 事件边收边 yield，不再等执行完成
- 新增 agentend 内部直接调用 adapter 的"短路路径"：Orchestrator 调度子 Agent 时跳过 BackendClient 的 HTTP 回环，直接在进程内调用 adapter.stream_chat()
- Backend SSE Handler XREAD 改为批量读取（Count: 100）+ 短阻塞（200ms），减少 Redis 往返
- 后端 StreamWriter Run 路径中合并高频小 event，减少 Redis XADD 和 HTTP Flush 调用频次
- **BREAKING**: Orchestrator 新增 `RUNTIME_TEXT` 事件类型用于透传子 Agent token；前端需识别并渲染该事件

## Capabilities

### New Capabilities

- `orchestrator-streaming-forward`: Orchestrator 边收边转发子 Agent token 级事件，包含 ExecutionEngine 流式透传改造和 adapter 短路调用路径
- `streaming-batch-read`: Backend SSE Handler 批量读取 Redis Stream 事件，减少逐条阻塞读取的开销

### Modified Capabilities

- `token-streaming`: 新增 `RUNTIME_TEXT` 事件类型，Orchestrator 通过该类型透传子 Agent 增量文本
- `chat-streaming`: 前端需识别并渲染 `RUNTIME_TEXT` 事件，按 task_id + agent 归属到对应子 Agent 消息块

## Impact

- **agentend**: `orchestrator/execution/engine.py`（核心改造）、`clients/backend_client.py`（短路路径下可能不再需要 stream_result）、`api/v1/agent.py`（内部调用 adapter）
- **backend**: `internal/handler/stream.go`（XREAD 批量读取）、`internal/stream/writer.go`（事件合并）
- **frontend**: chat store + MessageList 组件（识别 RUNTIME_TEXT 事件）
- **contracts**: StreamEvent schema 新增 `runtime_text` 事件类型
- **性能**: 预计端到端流式延迟从"子 Agent 完成后"降低到"接近实时"（< 200ms）
