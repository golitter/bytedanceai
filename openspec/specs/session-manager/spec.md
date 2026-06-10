## ADDED Requirements

### Requirement: Session data model
系统 SHALL 定义 `Session` 数据类，包含字段：`id`（str）、`agent_type`（str）、`state`（SessionState 枚举）、`process`（可选 asyncio.subprocess.Process）、`workspace_path`（str）、`created_at`（datetime）、`last_active`（datetime）、`history`（list[dict]）、`metadata`（dict）。`workspace_path` 字段 MUST 为必填 str 类型。

#### Scenario: Create session with workspace
- **WHEN** 调用 `SessionManager.create("claude-code", workspace_path="/workspaces/task-1/frontend")`
- **THEN** SHALL 创建 Session，`workspace_path` 为 `/workspaces/task-1/frontend`

#### Scenario: Create session without workspace fails
- **WHEN** 调用 `SessionManager.create("claude-code")` 不传 workspace_path
- **THEN** SHALL 抛出 `ValueError`，提示 workspace_path 为必填

### Requirement: SessionState state machine
`SessionState` SHALL 定义枚举值：`IDLE`、`RUNNING`、`COMPLETED`、`INTERRUPTED`、`ERROR`、`INACTIVE`。状态转移 MUST 遵循：IDLE → RUNNING → COMPLETED / INTERRUPTED / ERROR；任意非 INACTIVE 状态 → INACTIVE（仅由 backend API 写入）。

#### Scenario: Valid state transition
- **WHEN** session 当前状态为 `IDLE`，执行开始
- **THEN** 状态 SHALL 转移为 `RUNNING`

#### Scenario: Invalid state transition
- **WHEN** session 当前状态为 `COMPLETED`，尝试转移到 `RUNNING`
- **THEN** SHALL 抛出 `ValueError` 指示非法状态转移

### Requirement: INACTIVE status — backend API only
`INACTIVE` 状态 SHALL 仅由 backend `PATCH /api/sessions/:sessionId` API 写入，表示用户手动停用该 session。Agentend 不得直接修改 session 状态为 `INACTIVE`，只通过只读 DB 查询感知该状态。

#### Scenario: User deactivates session via backend
- **WHEN** 用户通过前端点击"停用"按钮
- **THEN** 前端 SHALL 调用 `PATCH /api/sessions/:sessionId` 且 body 为 `{ "status": "inactive" }`，由 backend 写入 DB

#### Scenario: Backend validates inactive transition
- **WHEN** backend 收到 `PATCH /api/sessions/:sessionId` 且 `status` 值为 `"inactive"`
- **THEN** SHALL 校验 session 存在且当前状态非 `INACTIVE`，通过后更新 DB 中 session 的 status 为 `inactive`

#### Scenario: Session not found
- **WHEN** backend 收到 `PATCH /api/sessions/:sessionId` 但 session 不存在
- **THEN** SHALL 返回 HTTP 404

#### Scenario: Invalid status value
- **WHEN** backend 收到 `PATCH /api/sessions/:sessionId` 且 `status` 值不为 `"inactive"`
- **THEN** SHALL 返回 HTTP 400 参数校验错误

### Requirement: SessionManager CRUD operations
`SessionManager` SHALL 提供 `create`、`get`、`list`、`update_state`、`destroy` 方法。`create` 方法 MUST 接受 `workspace_path` 参数。

#### Scenario: Create session
- **WHEN** 调用 `create(agent_type, workspace_path="/workspaces/xxx")`
- **THEN** SHALL 创建包含 workspace_path 的 Session 对象

#### Scenario: Get existing session
- **WHEN** 调用 `get(session_id)` 且 session 存在
- **THEN** SHALL 返回对应的 `Session` 对象

#### Scenario: Get non-existent session
- **WHEN** 调用 `get(session_id)` 且 session 不存在
- **THEN** SHALL 返回 `None`

#### Scenario: Destroy session
- **WHEN** 调用 `destroy(session_id)` 且 session 存在且有运行中进程
- **THEN** SHALL 先终止进程，再从内存中移除 session，返回 `True`

#### Scenario: Destroy non-existent session
- **WHEN** 调用 `destroy(session_id)` 且 session 不存在
- **THEN** SHALL 返回 `False`

### Requirement: Session history tracking
`Session` 的 `history` 列表 SHALL 自动记录每次 chat/stream_chat 的请求消息和响应摘要。`last_active` 时间戳 MUST 在每次交互时更新。

#### Scenario: History updated after chat
- **WHEN** 对 session 执行一次 `chat` 调用
- **THEN** session 的 `history` SHALL 新增请求消息记录，`last_active` SHALL 更新为当前时间
