## Context

多 Agent 聊天系统中，orchestrator 通过一次 RunTask 调用执行完整规划，期间在 SSE 流中交替输出 orchestrator 和子 Agent 的内容。当前后端为每次 RunTask 创建**一条** `model.Message`，所有 Agent 的文本拼接在一起，`agent_type` 固定为初始值（orchestrator）。页面刷新或历史加载后，前端无法区分消息归属。

同时 `ExecutionEngine` 为子 Agent 单独调用 `backend_client.run_task()`，在子 Agent 的 session 下产生重复消息。`ListMessages` 按 `task_id` 查询不分 session，导致重复内容混入。

前端 `streamAgentUpdate` 已能正确处理实时流中的 Agent 切换（通过 Redis 原始 SSE 事件中的 `agent_type` 字段），但该能力仅限实时路径，MySQL 持久化层缺失。

## Goals / Non-Goals

**Goals:**
- MySQL 中每个 Agent 的输出为独立的 Message 行，各有正确的 `agent_type` 和 `agent_name`
- 页面刷新/历史加载时，前端正确显示按 Agent 拆分的消息
- SSE 重放（Phase 1 + serveCompleted）携带 agent 元数据
- 消除 ExecutionEngine 双写导致的重复消息
- 实时流路径（Redis → streamAgentUpdate）不受影响

**Non-Goals:**
- 不修改 Agentend 代码（SSE 事件格式已正确）
- 不修改 Redis stream key 结构
- 不修改前端 store 的 streamAgentUpdate 逻辑
- 不引入新的 Message 数据库字段（schema 不变）

## Decisions

### D1: StreamWriter 内检测 agent_type 变化并创建新 Message

**选择**: 在 `writer.go` 的 `Run()` 方法中，解析 TEXT 事件的 `agent_type`/`agent` 字段，变化时调用 `switchAgent()` 在同 session 下创建新 Message。

**替代方案**:
- A) 在 Message 中加 `agent_segments` JSON 列，API 层拆分 → 增加字段、拆分逻辑复杂
- B) 前端 loadHistory 时按标记拆分 → 需要 agent 边界标记，污染内容

**理由**: 直接在 StreamWriter 层拆分是最自然的，每条 Message 独立完整，API/前端无需特殊处理。

### D2: 原始 Message 保持 streaming 状态直到整轮结束

**选择**: 原始 Message（RunTask 创建的第一条）不随 Agent 切换而 completed，仅在 `finish()` 时标记 completed。中间子消息在 Agent 切走时立即 completed。

**理由**: `ServeStream` 通过 `IsActive(originalMessageID)` 判断流是否活跃。若原始 Message 提前 completed，重连时 `serveCompleted` 会立即发送 done 导致前端丢失后续事件。

### D3: Redis stream key 不变

**选择**: 所有事件继续写入原始 Message 的 Redis stream（`agent:{sessionID}:{originalMessageID}`）。新建的子消息不创建独立 Redis stream。

**理由**: 前端连接到原始 messageID 的 SSE 端点，实时事件全部来自这一个 stream。`streamAgentUpdate` 已经能正确处理 agent 切换。拆分 Redis stream 会导致前端需要管理多个连接，复杂度大增。

### D4: ListMessages 加 session_id 过滤

**选择**: `ListMessages` 新增可选 `session_id` query param。前端传当前 session 的 session_id，只加载当前 session 的消息。

**理由**: ExecutionEngine 为子 Agent 调用 `run_task()` 创建的 Message 属于子 Agent 的 session，与 orchestrator session 的消息是重复内容。按 session 过滤后，orchestrator 视图只显示其自身 session 下的已拆分消息。

### D5: SSE 重放携带 agent 元数据

**选择**: `serveCompleted` 和 `serveStreaming` Phase 1 使用消息自身的 `AgentType`/`AgentName` 构造 SSE 事件，替代原来的 `FormatSSE`（不含 agent 信息）。

**理由**: 重连或加载已完成消息时，前端需要 agent 元数据来正确归属消息。修改后的 SSE 事件格式：`{type:"text", content:{text, agent, agent_type}}`。

## Risks / Trade-offs

- **[Agent 快速切换时可能产生短命 Message]** → 若 Agent 切换极快（毫秒级），可能产生 content 为空的 Message。Mitigation: `switchAgent()` 中检查 buffer 是否为空，空则跳过 finalize 直接更新 agent 信息。
- **[原始 Message streaming 期间重连]** → Phase 1 发送原始 Message 的 MySQL content（仅第一个 Agent 的内容），Phase 2 从 Redis 续读。若 Redis 事件已过期（TTL 600s），部分内容丢失。Mitigation: 现有 TTL 已足够长，且子消息已有独立 MySQL 记录。
- **[ListMessages session_id 过滤为可选]** → 不传 session_id 时行为不变（向后兼容），但可能仍有重复。Mitigation: 前端统一传 session_id。
