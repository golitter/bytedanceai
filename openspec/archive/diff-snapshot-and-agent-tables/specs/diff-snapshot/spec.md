## ADDED Requirements

### Requirement: DiffSnapshot model
系统 SHALL 提供 `diff_snapshots` 数据表，包含以下字段：`id`（自增主键）、`snapshot_id`（VARCHAR 36，UNIQUE）、`session_id`（VARCHAR 128）、`diff_content`（LONGTEXT）、`status`（VARCHAR 16，默认 `'pending'`）、`created_at`、`updated_at`。`snapshot_id` 为 agent 生成的 UUID，是查询主键。

#### Scenario: AutoMigrate creates table
- **WHEN** 后端启动执行 AutoMigrate
- **THEN** SHALL 自动创建 `diff_snapshots` 表，包含上述字段及 `idx_session` 索引

### Requirement: GET diff snapshot by snapshotId
系统 SHALL 提供 `GET /api/diff-snapshots/:snapshotId` 端点，根据 snapshotId 返回对应记录。返回 JSON 包含 `snapshot_id`、`session_id`、`diff_content`、`status`、`created_at`、`updated_at`。

#### Scenario: Snapshot exists
- **WHEN** 调用 `GET /api/diff-snapshots/abc-123`
- **THEN** SHALL 返回 200 及该 snapshot 完整数据

#### Scenario: Snapshot not found
- **WHEN** 调用 `GET /api/diff-snapshots/nonexistent`
- **THEN** SHALL 返回 404

### Requirement: PUT diff snapshot with auto-cancel
系统 SHALL 提供 `PUT /api/diff-snapshots/:snapshotId` 端点，接收 `{ session_id, diff, status }` JSON body。当 `status` 为 `pending` 时，MUST 先将同 `session_id` 下所有 `status='pending'` 且 `snapshot_id != 当前值` 的记录更新为 `cancelled`，再 upsert 当前记录。

#### Scenario: Create new pending snapshot
- **WHEN** 调用 `PUT /api/diff-snapshots/aaa` body `{ session_id: "s1", diff: "...", status: "pending" }` 且 s1 无其他 pending 记录
- **THEN** SHALL 创建 snapshot 记录，status 为 pending

#### Scenario: New pending cancels previous pending
- **WHEN** 调用 `PUT /api/diff-snapshots/bbb` body `{ session_id: "s1", diff: "...", status: "pending" }` 且 s1 已有 pending 记录 aaa
- **THEN** SHALL 将 aaa 的 status 更新为 `cancelled`，bbb 创建为 `pending`

#### Scenario: Commit an existing snapshot
- **WHEN** 调用 `PUT /api/diff-snapshots/aaa` body `{ session_id: "s1", diff: "...", status: "committed" }`
- **THEN** SHALL 更新 aaa 的 status 为 `committed`，diff_content 为请求中的值

#### Scenario: Revert an existing snapshot
- **WHEN** 调用 `PUT /api/diff-snapshots/aaa` body `{ session_id: "s1", diff: "...", status: "reverted" }`
- **THEN** SHALL 更新 aaa 的 status 为 `reverted`，diff_content 为请求中的值

### Requirement: DiffSnapshot status transitions
状态流转 SHALL 为：`pending → committed | reverted | cancelled`。终态（committed、reverted、cancelled）MUST 不可变更。PUT 请求对终态记录 SHALL 返回 409 Conflict。

#### Scenario: Update terminal state
- **WHEN** 调用 `PUT /api/diff-snapshots/aaa` 且 aaa 当前 status 为 `committed`
- **THEN** SHALL 返回 409 Conflict，不更新记录

### Requirement: DiffCard snapshot-driven rendering
DiffCard 组件 SHALL 接收 `snapshotId` 和 `sessionId` props。首次渲染时调 `GET /api/diff-snapshots/:snapshotId`，若 404 则从 workspace diff 获取内容并创建 pending snapshot。若 200 且 status 为终态，则直接用 snapshot 的 diff_content 渲染并显示对应徽章。

#### Scenario: First render creates pending
- **WHEN** DiffCard(snapshotId="aaa") 首次渲染且 GET 返回 404
- **THEN** SHALL 调 `GET /api/session/{sid}/diff` 获取 workspace diff，再调 `PUT /api/diff-snapshots/aaa` 创建 pending，显示 diff + 操作按钮

#### Scenario: Page reload shows settled state
- **WHEN** DiffCard(snapshotId="aaa") 渲染且 GET 返回 `{ status: "committed", diff_content: "..." }`
- **THEN** SHALL 用 snapshot 的 diff_content 渲染 diff，显示"已接受"徽章，隐藏操作按钮

#### Scenario: Cancelled snapshot shows greyed out
- **WHEN** DiffCard(snapshotId="aaa") 渲染且 GET 返回 `{ status: "cancelled", diff_content: "..." }`
- **THEN** SHALL 用 snapshot 的 diff_content 渲染 diff（灰化），显示"已取消"徽章，无操作按钮
