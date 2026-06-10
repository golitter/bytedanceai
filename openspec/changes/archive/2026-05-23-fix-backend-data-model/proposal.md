## Why

Phase 1 实现的数据模型与 AgentEnd 的实际概念不匹配。当前 Session 被当作顶层容器、Task 被当作一次性消息，但实际上 **Task 才是群组任务**（绑定 repo_path），**Session 是 Agent 的会话**（session_id 由调用方传入）。这导致 API 设计错误，前端无法正确对接。

## What Changes

- **BREAKING** 重定义 Session/Task 数据模型：Task = 群组任务（顶层），Session = Agent 会话（属于 Task）
- **BREAKING** API 路径重构：从 `/api/sessions` 为核心改为 `/api/tasks` 为核心
- **BREAKING** `session_id` 不再由 Go 后端生成，改由调用方传入（AgentEnd 用 `session_id::task_id` 映射 cli_session_id）
- 修改 Handler 层适配新的模型关系和 API 路径
- 修改路由注册

## Capabilities

### New Capabilities
- `task-api`: Task（群组任务）的 CRUD REST API — `/api/tasks`
- `session-agent-api`: Session（Agent 会话）的运行 API — `POST /api/tasks/:taskId/run`

### Modified Capabilities
- `agent-list-api`: 重命名为 `agent-types-api`，路径从 `/api/agents` 改为 `/api/agent-types`
- `agentend-client`: 无需修改（已正确封装）
- `session-api`: 删除（被 `task-api` + `session-agent-api` 替代）
- `task-run-api`: 删除（被 `session-agent-api` 替代）

## Impact

- `backend/internal/model/` — Session 和 Task 结构体字段互换
- `backend/internal/handler/` — 重写所有 handler（session → task + session-agent）
- `backend/cmd/server/main.go` — 路由注册更新
- API 端点：`/api/sessions` → `/api/tasks`，`/api/sessions/:sid/tasks/run` → `/api/tasks/:taskId/run`
- 数据库：`sessions` 和 `tasks` 表结构变化（AutoMigrate 自动处理）
- 不涉及 AgentEnd 和前端改动
