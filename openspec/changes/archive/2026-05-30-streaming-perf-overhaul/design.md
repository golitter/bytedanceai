## Context

当前流式输出架构中，数据从 Agent CLI 到前端用户需经过多层中转。单 Agent 模式链路为：CLI → agentend adapter SSE → Backend StreamWriter → Redis XADD → Redis XREAD(5s block, Count:1) → Backend SSE Handler → Frontend。Orchestrator 多 Agent 模式更严重：子 Agent 的 token 先走完上述全链路被 BackendClient.stream_result() 读取，全量缓冲到 collected[] 后才发出一条 RUNTIME_COMPLETED 事件，然后再走一遍 StreamWriter → Redis → SSE Handler 才到前端。

核心瓶颈：
1. **Orchestrator 全量缓冲**：engine.py `_collect_stream()` 收集所有文本后才 yield
2. **四层网络回环**：Orchestrator → BackendClient HTTP → Backend → agentend → Backend → Redis → Backend → SSE → agentend
3. **Redis 逐条阻塞读取**：XREAD Count:1 + Block 5s
4. **逐 token Flush**：每条 SSE 消息单独系统调用

## Goals / Non-Goals

**Goals:**
- 子 Agent token 实时透传到前端，延迟 < 200ms
- 单 Agent 模式 Redis 往返开销降低 10x+
- 保持 Redis Stream 架构不变（断线重连、历史回放仍依赖它）

**Non-Goals:**
- 重写前端 Markdown 渲染引擎
- 移除 Redis Stream 层（改为 WebSocket 等）
- 修改 Agent CLI 的输出格式

## Decisions

### Decision 1: Orchestrator 进程内短路调用 adapter

Orchestrator 调度子 Agent 时，不再通过 BackendClient 发 HTTP 请求回 Backend 再到 agentend，而是直接在 agentend 进程内调用 adapter.stream_chat()。

**路径**：`Orchestrator → adapter.stream_chat() → SSE events 直接 yield`

**收益**：消除 Orchestrator ↔ Backend 的 HTTP 回环（省掉 BackendClient.run_task + stream_result 两步），token 直接从 adapter 流入 Orchestrator 的 SSE 响应。

**替代方案**：保留 HTTP 调用但改为流式透传 — 更复杂，收益不如进程内调用大。

**注意**：短路路径下 Orchestrator 仍然需要调用 BackendClient.run_task() 来创建 agent message 记录（用于前端定位消息），但不再需要 stream_result() 读取回环。

### Decision 2: ExecutionEngine 流式透传 token

将 `_collect_stream()` 从"全量收集"改为"边收边 yield"：每个 TEXT 事件包装为 `RUNTIME_TEXT` StreamEvent 立即 yield，同时保留 `collected` 列表用于最终 TaskResult。

**新事件类型**：`RUNTIME_TEXT`，content 包含 `{"task_id": "...", "agent": "...", "text": "增量文本"}`

### Decision 3: Backend SSE Handler 批量 XREAD

将 `stream.go` 中 `serveStreaming` 的实时阶段从 `Count:1, Block:5s` 改为 `Count:100, Block:200ms`。批量读取后一次 Flush，减少系统调用。

**替代方案**：使用 XRANGE + 定时轮询 — 实现更复杂，XREAD 批量已足够。

### Decision 4: StreamWriter 事件合并

在 StreamWriter.Run 的 scanFunc 回调中，将高频小 TEXT 事件合并为较大的批次再 XADD，阈值约 500ms 或 2KB。

**替代方案**：前端节流 — 治标不治本，后端减少 XADD 频次更有效。

## Risks / Trade-offs

- [Risk] 进程内短路调用绕过 Backend StreamWriter，子 Agent token 不经过 Redis → **Mitigation**: 短路路径产生的 RUNTIME_TEXT 事件仍由 Orchestrator 的外层 SSE 流经 Backend StreamWriter → Redis，前端通过同一个 SSE 连接消费
- [Risk] 批量 XREAD 可能引入额外 200ms 延迟 → **Mitigation**: Redis 有消息时 XREAD 立即返回，200ms block 仅在无消息时生效；相比当前 5s block 已大幅改善
- [Risk] RUNTIME_TEXT 新事件类型导致旧前端不识别 → **Mitigation**: 前端忽略未知 type 的 event（当前实现已支持），新前端再添加渲染逻辑
