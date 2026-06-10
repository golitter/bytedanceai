## MODIFIED Requirements

### Requirement: Session data model
系统 SHALL 定义 `Session` 数据类，包含字段：`id`（str）、`agent_type`（str）、`state`（SessionState 枚举）、`process`（可选 asyncio.subprocess.Process）、`workspace_path`（str）、`created_at`（datetime）、`last_active`（datetime）、`history`（list[dict]）、`metadata`（dict）。`workspace_path` 字段 MUST 为必填 str 类型。

#### Scenario: Create session with workspace
- **WHEN** 调用 `SessionManager.create("claude-code", workspace_path="/workspaces/task-1/frontend")`
- **THEN** SHALL 创建 Session，`workspace_path` 为 `/workspaces/task-1/frontend`

#### Scenario: Create session without workspace fails
- **WHEN** 调用 `SessionManager.create("claude-code")` 不传 workspace_path
- **THEN** SHALL 抛出 `ValueError`，提示 workspace_path 为必填

### Requirement: SessionManager CRUD operations
`SessionManager` SHALL 提供 `create`、`get`、`list`、`update_state`、`destroy` 方法。`create` 方法 MUST 接受 `workspace_path` 参数。

#### Scenario: Create session
- **WHEN** 调用 `create(agent_type, workspace_path="/workspaces/xxx")`
- **THEN** SHALL 创建包含 workspace_path 的 Session 对象
