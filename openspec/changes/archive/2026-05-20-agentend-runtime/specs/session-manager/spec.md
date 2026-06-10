## ADDED Requirements

### Requirement: Session data model
系统 SHALL 定义 `Session` 数据类，包含字段：`id`（str）、`agent_type`（str）、`state`（SessionState 枚举）、`process`（可选 asyncio.subprocess.Process）、`workspace_path`（可选 str）、`created_at`（datetime）、`last_active`（datetime）、`history`（list[dict]）、`metadata`（dict）。

#### Scenario: Create new session
- **WHEN** 调用 `SessionManager.create("claude-code", {"workspace": "/tmp"})` 创建会话
- **THEN** SHALL 生成唯一 `session_id`，`state` 为 `IDLE`，`created_at` 和 `last_active` 为当前时间，返回完整 `Session` 对象

### Requirement: SessionState state machine
`SessionState` SHALL 定义枚举值：`IDLE`、`RUNNING`、`COMPLETED`、`INTERRUPTED`、`ERROR`。状态转移 MUST 遵循：IDLE → RUNNING → COMPLETED / INTERRUPTED / ERROR。

#### Scenario: Valid state transition
- **WHEN** session 当前状态为 `IDLE`，执行开始
- **THEN** 状态 SHALL 转移为 `RUNNING`

#### Scenario: Invalid state transition
- **WHEN** session 当前状态为 `COMPLETED`，尝试转移到 `RUNNING`
- **THEN** SHALL 抛出 `ValueError` 指示非法状态转移

### Requirement: SessionManager CRUD operations
`SessionManager` SHALL 提供 `create`、`get`、`list`、`update_state`、`destroy` 方法，全部在内存中管理。

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
