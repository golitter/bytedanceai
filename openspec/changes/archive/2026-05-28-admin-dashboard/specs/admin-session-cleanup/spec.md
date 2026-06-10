## ADDED Requirements

### Requirement: Session list with filters
会话清理页面 SHALL 展示会话列表，支持按 Agent 类型筛选。

#### Scenario: Display sessions
- **WHEN** 用户打开会话清理页面
- **THEN** 显示会话表格，包含列：会话 ID、Agent 类型、任务标题、最后活跃时间、状态、操作

#### Scenario: Filter by agent type
- **WHEN** 用户选择特定 Agent 类型筛选
- **THEN** 表格只显示该 Agent 类型的会话

### Requirement: Batch session cleanup
用户 SHALL 能通过复选框选中多个会话并批量清理。

#### Scenario: Select and delete sessions
- **WHEN** 用户选中多个会话并点击清理按钮
- **THEN** 系统删除选中会话及其关联数据（消息、session_agent 记录），显示 toast 通知已清理数量

#### Scenario: Delete single session
- **WHEN** 用户点击某条会话的删除按钮
- **THEN** 系统删除该会话及其关联数据，显示 toast 通知

### Requirement: Session cleanup API
后端 SHALL 提供 `DELETE /api/admin/sessions` 接口，支持批量删除。

#### Scenario: Batch delete API
- **WHEN** 前端发送 `DELETE /api/admin/sessions` 且 body 包含 `{ session_ids: ["id1", "id2"] }`
- **THEN** 后端删除指定会话及其关联消息和 session_agent 记录，返回 `{ deleted: 2 }`
