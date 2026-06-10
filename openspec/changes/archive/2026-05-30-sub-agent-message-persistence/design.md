## Context

当前架构中，用户消息 → 后端 → agentend(orchestrator) 这条链路的消息能持久化。但 orchestrator 分派子 Agent 时，在 agentend 内部直接调用 `adapter.stream_chat()`，子 Agent 的响应只在 SSE 实时流中可见，不经过后端的 StreamWriter，因此不会写入 MySQL。

关键文件：
- `agentend/src/orchestrator/execution/engine.py` — `_execute_task()` 直接调用 adapter
- `backend/internal/handler/task.go` — `RunTask()` 创建 Message + StreamWriter
- `backend/internal/stream/writer.go` — StreamWriter 负责增量持久化

现有 RunTask 流程已完整：创建 user message → 创建 agent message(streaming) → 启动 goroutine 消费 SSE → StreamWriter 增量写入 MySQL。子 Agent 只需复用这个流程。

## Goals / Non-Goals

**Goals:**
- 子 Agent（claude-code、opencode）的响应消息写入 MySQL，刷新后可恢复
- 复用后端已有的 RunTask + StreamWriter 机制，不引入新的持久化路径
- 子 Agent 消息通过 Redis Stream 推送到前端，实现实时显示
- 前端无改动 —— 现有的 SSE 订阅 + getTaskMessages 已能处理多 session 消息

**Non-Goals:**
- 不修改 contracts 层 —— Message model 和 SSE event 格式不变
- 不修改前端代码
- 不引入消息队列（Kafka/RabbitMQ）
- 不做子 Agent 会话的完整 history 管理（只存 run 产生的消息）

## Decisions

### D1: agentend 回调后端 RunTask API，而非后端主动拉取

**选择**: ExecutionEngine 调用后端 `POST /api/tasks/:taskId/run` 发起子 Agent run。

**备选**:
- (B) agentend 直接写 MySQL — 需要数据库连接，违反架构分层
- (C) orchestrator 将子 Agent 响应合并到自己的消息中 — 丢失独立性，前端无法区分

**理由**: RunTask 已具备完整的消息创建 + StreamWriter + Redis 推送链路。回调 RunTask 是最小改动量，零重复代码。后端天然成为所有消息的单一写入点。

### D2: 通过 agentend 的后端 client 调用，非 agentend 内部 API

**选择**: agentend 新增一个 `backend_client` 模块，内部发起 HTTP 请求到后端 RunTask。

**理由**: agentend 已有 `backend_client` 的概念（虽然目前主要用于接收请求），新增一个出站 HTTP client 符合现有架构。后端 RunTask 本身就是 HTTP endpoint，无需新增内部 RPC。

### D3: 子 Agent 的 SSE 事件通过后端 Redis Stream 推送到前端

**选择**: 子 Agent 的 RunTask 走标准流程 → StreamWriter 写 Redis Stream → 前端 SSE 消费。

**理由**: 前端已通过 `GET /stream?message_id=xxx` 订阅 Redis Stream。子 Agent 的 message_id 不同，前端只需为每个子 Agent 的 message_id 建立独立 SSE 连接。但这需要前端改动...

**修正**: 实际上前端目前只订阅 orchestrator 的 SSE 流。子 Agent 的实时事件仍由 orchestrator 的 SSE 流转发（通过 `RUNTIME_EXECUTING` / `TEXT`(带 agent 字段) / `RUNTIME_COMPLETED` 事件）。持久化走后端 StreamWriter 是独立路径，不依赖前端 SSE 订阅。

因此最终方案：
1. **实时显示**: orchestrator 的 SSE 流继续转发子 Agent 事件（现有行为不变）
2. **持久化**: ExecutionEngine 额外回调后端 RunTask，后端为子 Agent 创建 Message 并 StreamWriter 持久化
3. **刷新恢复**: 前端从 getTaskMessages 加载所有 session 的消息

### D4: ExecutionEngine 改为 fire-and-forget 调用 RunTask

**选择**: ExecutionEngine 调用 RunTask 后，不等待其完成。子 Agent 的执行路径有两条：
1. agentend 内部继续执行 adapter（实时 SSE 推送到前端，通过 orchestrator 流）
2. 后端 RunTask 也执行同一个 adapter（StreamWriter 持久化到 MySQL）

**问题**: 这会导致子 Agent 被执行两次！

**修正方案**: ExecutionEngine 不再直接调用 adapter，而是**只回调后端 RunTask**。后端 RunTask 通过 agentend 的 `/v1/agent/stream` 执行子 Agent。这样：
1. ExecutionEngine 调用 `POST /api/tasks/:taskId/run`（传子 Agent 的 session_id）
2. 后端创建 Message，调用 agentend 的 `/v1/agent/stream`
3. agentend 执行子 Agent adapter，SSE 返回后端
4. 后端 StreamWriter 持久化 + Redis 推送
5. ExecutionEngine 等待后端返回，收集子 Agent 结果

这样子 Agent 只执行一次，且消息被正确持久化。

### D5: ExecutionEngine 如何等待子 Agent 完成

**选择**: ExecutionEngine 调用 RunTask 后，通过 SSE 订阅后端返回的 `message_id`，收集文本内容作为 TaskResult。

**理由**: RunTask 返回 `202 + message_id`，ExecutionEngine 可以：
- 方案 A: 订阅后端 SSE `GET /stream?message_id=xxx` 收集内容
- 方案 B: 调用后端的 messages API 轮询消息状态
- 方案 C: RunTask 改为同步等待完成（阻塞 goroutine，不合适）

方案 A 最干净：复用现有 SSE 机制，ExecutionEngine 作为 SSE 消费者收集子 Agent 的文本输出。

## Risks / Trade-offs

- **[延迟增加]** 子 Agent 执行多一跳（agentend → 后端 → agentend）→ 影响可控，后端和 agentend 通常同机部署，网络延迟 <1ms
- **[orchestrator 流中子 Agent 实时事件减少]** ExecutionEngine 不再 yield 子 Agent 的 TEXT 事件，前端实时显示依赖后端 Redis Stream → 但前端目前只订阅 orchestrator 的 SSE，需要前端为每个子 Agent 建立额外 SSE 连接，或由 orchestrator 流转发关键状态事件
- **[RunTask 并发]** orchestrator 可能同时分派多个子 Agent → RunTask 已在 goroutine 中运行，后端天然支持并发

## Open Questions

- 前端如何发现并订阅子 Agent 的 SSE 流？可能需要 orchestrator 在 dispatch 阶段通过 SSE 推送子 Agent 的 `message_id`，前端据此建立额外 SSE 连接。但这涉及前端改动，是否在本 change 范围内？
