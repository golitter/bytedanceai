## Context

Phase 1 实现了 Go Backend 的 Session/Task CRUD + AgentEnd SSE 透传，但数据模型与 AgentEnd 的实际概念不匹配。

AgentEnd 的实际概念：
- `task_id` = 群组任务，决定 git branch `task/{task_id}` 和 worktree 根目录 `worktrees/{task_id}/`
- `session_id` = Agent 的会话，由调用方传入，AgentEnd 用 `session_id::task_id` 映射到 `cli_session_id`
- 同一 task 下可有多个 session（多 agent 协作）

当前实现的问题：
- Session 被当作顶层容器（有 repo_path），Task 被当作一次性消息
- session_id 由 Go 后端自动生成，而非调用方传入
- API 路径以 `/api/sessions` 为核心

## Goals / Non-Goals

**Goals:**
- 修正 Task/Session 模型，使其语义与 AgentEnd 一致
- Task 成为顶层实体（群组任务，绑定 repo_path）
- Session 成为 Agent 会话（session_id 由调用方传入，属于某个 Task）
- API 路径以 `/api/tasks` 为核心
- 保持 SSE 透传机制不变

**Non-Goals:**
- 不修改 AgentEnd 代码
- 不碰前端（Phase 2 前端对接时再适配）
- 不改变 AgentEnd Client 的封装
- 不做认证、EventLog 持久化等 Phase 2+ 内容

## Decisions

### D1: Task 为顶层实体，Session 从属
Task 持有 `repo_path`、`title`，是用户的"群组任务"。Session 记录一个 Agent 在某个 Task 里的会话，`session_id` 由调用方传入。

**理由**: 与 AgentEnd 的 `task_id`（workspace 隔离）和 `session_id`（CLI session 映射）语义对齐。

### D2: session_id 由调用方传入
`POST /api/tasks/:taskId/run` 接受 `session_id` 参数。Go 后端不负责生成 session_id，只负责记录和传递。

- 调用方传新 UUID → AgentEnd 创建新 session + workspace
- 调用方传已有 session_id → AgentEnd 自动 resume cli_session_id

**理由**: AgentEnd 用 `session_id::task_id` 做 CLI session 映射，session_id 的生命周期由调用方控制。

### D3: Session 记录按需创建
调用 `/api/tasks/:taskId/run` 时，如果数据库中没有对应 session_id 的记录，自动创建。

**替代方案**: 单独的 `POST /api/tasks/:taskId/sessions` 端点 — 增加复杂度，Phase 1 不需要。

### D4: 删除旧 spec，新建替代 spec
`session-api` 和 `task-run-api` 被删除，用 `task-api` 和 `session-agent-api` 替代。`agent-list-api` 改名为 `agent-types-api`。

## Risks / Trade-offs

- **[BREAKING API 变更]** → Phase 1 尚无前端对接，影响可控
- **[AutoMigrate 不可逆]** → 旧表结构会被修改，开发阶段可接受，需清空旧数据
- **[调用方需管理 session_id]** → 前端需自行生成和追踪 session_id，增加调用方复杂度
