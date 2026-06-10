## ADDED Requirements

### Requirement: Backend announcement data model
后端 SHALL 提供 `task_announcements` 表，包含字段：`id`（BIGINT 主键）、`task_id`（BIGINT 索引）、`sender_id`（VARCHAR(64)）、`sender_name`（VARCHAR(64)）、`content`（TEXT）、`pinned`（BOOLEAN）、`created_at`（DATETIME）。

#### Scenario: AutoMigrate creates table
- **WHEN** 后端启动
- **THEN** GORM AutoMigrate 自动创建 `task_announcements` 表，包含所有字段和索引

### Requirement: List announcements API
后端 SHALL 提供 `GET /api/tasks/{taskId}/announcements` 接口，返回该 task 下所有公告，按置顶优先、时间倒序排列。

#### Scenario: Fetch announcements for task
- **WHEN** 前端请求 `GET /api/tasks/123/announcements`
- **THEN** 返回 task 123 下所有公告，置顶的排在前面，其余按创建时间倒序

#### Scenario: Empty announcements list
- **WHEN** task 下没有公告
- **THEN** 返回空数组 `[]`

### Requirement: Create announcement API
后端 SHALL 提供 `POST /api/tasks/{taskId}/announcements` 接口，接收 `sender_id`、`sender_name`、`content`、`pinned` 字段，创建新公告。

#### Scenario: Create announcement successfully
- **WHEN** 前端 POST `{ sender_id: "user1", sender_name: "田乐檬", content: "重构已启动", pinned: true }`
- **THEN** 创建公告并返回 201，包含生成的 `id` 和 `created_at`

#### Scenario: Missing required fields
- **WHEN** POST 请求缺少 `content` 字段
- **THEN** 返回 400 错误 "content is required"

### Requirement: Delete announcement API
后端 SHALL 提供 `DELETE /api/tasks/{taskId}/announcements/{id}` 接口，删除指定公告。

#### Scenario: Delete announcement successfully
- **WHEN** 前端 DELETE `/api/tasks/123/announcements/456`
- **THEN** 删除公告 456，返回 200

#### Scenario: Delete non-existent announcement
- **WHEN** 前端 DELETE 不存在的公告 ID
- **THEN** 返回 404 错误

### Requirement: Frontend announcements display
前端 SHALL 展示公告列表，每条公告显示发布者头像、名称、内容、时间，置顶公告带 📌 badge。

#### Scenario: Pinned announcement shows badge
- **WHEN** 某条公告 `pinned === true`
- **THEN** 公告卡片顶部显示 📌 置顶 badge，排在列表最前

#### Scenario: Announcement card layout
- **WHEN** 公告渲染
- **THEN** 显示发布者头像、名称、公告内容、发布时间

### Requirement: Create announcement button
Owner/Admin 角色 SHALL 看到「+ 发布新公告」按钮，点击后弹出输入框，确认后调用创建 API。

#### Scenario: Owner sees create button
- **WHEN** 当前用户角色为 Owner 或 Admin
- **THEN** 公告列表底部显示「+ 发布新公告」按钮

#### Scenario: Member does not see create button
- **WHEN** 当前用户角色为普通 Member
- **THEN** 不显示「+ 发布新公告」按钮

### Requirement: Task pinned_at field
后端 SHALL 在 `tasks` 表新增 `pinned_at` 字段（DATETIME DEFAULT NULL），通过 `PATCH /api/tasks/{taskId}` 更新。

#### Scenario: Pin a task
- **WHEN** 前端 PATCH `{ pinned_at: "2026-06-01T12:00:00Z" }`
- **THEN** task 的 `pinned_at` 更新为指定时间

#### Scenario: Unpin a task
- **WHEN** 前端 PATCH `{ pinned_at: null }`
- **THEN** task 的 `pinned_at` 设为 NULL
