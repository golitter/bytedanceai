## MODIFIED Requirements

### Requirement: SessionState state machine
`SessionState` SHALL 定义枚举值：`IDLE`、`RUNNING`、`COMPLETED`、`INTERRUPTED`、`ERROR`、`INACTIVE`。状态转移 MUST 遵循：IDLE → RUNNING → COMPLETED / INTERRUPTED / ERROR。`INACTIVE` 状态不通过 agent 自动状态转移产生，MUST 仅由外部系统（Backend API）设置。

#### Scenario: Valid state transition
- **WHEN** session 当前状态为 `IDLE`，执行开始
- **THEN** 状态 SHALL 转移为 `RUNNING`

#### Scenario: Invalid state transition
- **WHEN** session 当前状态为 `COMPLETED`，尝试转移到 `RUNNING`
- **THEN** SHALL 抛出 `ValueError` 指示非法状态转移

#### Scenario: Inactive set by external API
- **WHEN** Backend 接收到 PATCH /api/sessions/:sessionId { status: "inactive" } 请求
- **THEN** SHALL 将 session status 更新为 `inactive`，无论之前是何种状态（running/completed/error/interrupted）

## ADDED Requirements

### Requirement: Backend PATCH session status API
Backend SHALL 提供 `PATCH /api/sessions/:sessionId` API，接受 `{ "status": "inactive" }` 请求体，将指定 session 的 status 更新为 `inactive`。

#### Scenario: Deactivate session
- **WHEN** 前端调用 `PATCH /api/sessions/s-123` body `{ "status": "inactive" }`
- **THEN** Backend SHALL 将 session s-123 的 status 更新为 `inactive`，返回 200

#### Scenario: Session not found
- **WHEN** 前端调用 `PATCH /api/sessions/non-existent`
- **THEN** SHALL 返回 404

#### Scenario: Invalid status value
- **WHEN** 前端调用 `PATCH /api/sessions/s-123` body `{ "status": "invalid_status" }`
- **THEN** SHALL 返回 400，提示 status 值无效

### Requirement: Frontend session deactivation UI
Frontend session 列表 SHALL 为每个 session 提供"停用"操作入口，点击后调用 Backend PATCH API 将 session status 设为 inactive。

#### Scenario: User deactivates session
- **WHEN** 用户在 session 列表点击"停用"按钮
- **THEN** Frontend SHALL 调用 `PATCH /api/sessions/:sessionId` 并刷新列表状态显示
