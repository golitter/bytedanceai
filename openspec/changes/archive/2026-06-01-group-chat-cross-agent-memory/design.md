## Context

当前群聊场景中，每个子 Agent 执行时只能看到自己的 CLI session 历史和 Orchestrator 分发的任务内容。Agent 之间完全隔离，无法协作。系统已有三层架构（Frontend → Go Backend → AgentEnd），Rules 引擎已有 SafetyRule、ScopeRule、SoulRule 等内置规则。ExecutionEngine 存在 short-circuit 路径（直接调用本地 adapter），绕过了 Backend 和 Rules 引擎。

## Goals / Non-Goals

**Goals:**
- 子 Agent 执行时能看到「自上次发言以来其他 Agent 的消息」（窗口语义）
- Orchestrator 自身也能看到其他 Agent 的消息
- 复用现有 Rules 引擎，以 GroupChatRule 形式注入
- 统一执行路径，确保 Rules 对所有 Agent 生效
- 单聊场景零影响

**Non-Goals:**
- 不做消息摘要（先用截断策略）
- 不改 Frontend
- 不改子 Agent adapter 代码（claude.py / opencode.py / codex.py）
- 不改 Orchestrator 的 LangGraph 图结构

## Decisions

### D1: 窗口查询放在 Go Backend RunTask 中，而非 AgentEnd

**选择**: Go Backend 查 MySQL → 注入请求体 → AgentEnd 消费

**理由**: Backend 拥有 MySQL 连接，是数据的 natural owner。AgentEnd 无直接数据库访问能力。查询放在 Backend 避免了 AgentEnd 反向查数据库的额外依赖。

**替代方案**: AgentEnd 通过 API 查 Backend → 窗口查询 → 再注入。多一次网络往返，且 AgentEnd 需要知道何时该查窗口。

### D2: 不需要 is_group_chat flag，系统天然自适应

**选择**: Backend 每次都查窗口，单聊自然为空

**理由**: 单聊不经过 ExecutionEngine（直接 Frontend → Backend → AgentEnd → adapter）。Backend RunTask 查窗口时，单聊 task 只有 1 个 session，结果自然为空 `[]`，GroupChatRule 收到空列表不触发。无需任何条件分支或 flag。

**替代方案**: Orchestrator 传 `is_group_chat` flag 给 ExecutionEngine。但 ExecutionEngine 只在群聊时被调用，flag 永远为 True，等于没用。

### D3: 删除 ExecutionEngine 的 short-circuit 路径

**选择**: 统一走 HTTP（BackendClient.run_task → Backend → AgentEnd）

**理由**: Short-circuit 绕过 Backend 和 Rules 引擎，导致 Rules 不一致。ExecutionEngine 只在群聊（Orchestrator）场景下被调用，统一走 HTTP 无性能损失。单聊走另一条路径（不经过 ExecutionEngine），不受影响。

### D4: Wave 级别的窗口隔离

**选择**: 依赖 ExecutionEngine 串行处理 wave 的天然时序保证

**理由**: `asyncio.gather` 等 wave 全部完成后才进下一个 wave。Wave 2 的 RunTask 天然在 Wave 1 全部完成之后，MySQL 中 Wave 1 的消息已全部 persisted。无需额外锁或协调机制。

### D5: 消息截断策略

**选择**: 每条消息截断 2000 字符（rune-based）

**理由**: 保守值。CLI Agent 单次输出通常几百到几千字符，2000 字符足够传达"做了什么"的核心信息，同时不会挤占 Agent 的工作 context window。后续可根据实际体验调整。

### D6: 状态过滤包含 streaming 作为兜底

**选择**: `status IN ("completed", "streaming")`

**理由**: Wave 串行处理保证正常流程下只看到 completed 消息。但万一有 timing 极端情况，包含 streaming 作为兜底比完全丢失信息更安全。

## Risks / Trade-offs

- **[内容截断丢失信息]** → 2000 字符可能截断关键代码片段。后续可升级为智能摘要（让 LLM 摘要后再注入），但当前先用截断作为 MVP。
- **[ExecutionEngine 统一 HTTP 增加延迟]** → 多一跳 HTTP 往返（~10-50ms）。在群聊场景中，Agent 之间本身就有等待（wave 串行），额外延迟可忽略。单聊不受影响。
- **[窗口查询失败]** → Backend 查窗口失败时，RunTask 仍应正常执行（降级为无跨 Agent 上下文），不应阻塞任务执行。
- **[Orchestrator 自身查询的网络开销]** → 每次 `stream_chat` 都调 Backend API。暂不优化，后续可加本地缓存或直接查 Redis stream。
