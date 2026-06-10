## ADDED Requirements

### Requirement: Leave group with confirmation dialog
系统 SHALL 在群聊 RightSidebar 中提供"退出群聊"按钮。点击后弹出二次确认对话框，明确列出将被彻底删除的内容：群聊消息、Agent 工作区（Git worktree）、任务分支、Agent 分支、代码差异快照。用户确认后执行删除，取消则关闭对话框。

#### Scenario: User confirms leave group
- **WHEN** 用户点击"退出群聊" → 确认对话框 → 点击"确认退出"
- **THEN** 系统调用 `DELETE /api/tasks/:taskId/leave`，会话从列表中移除，当前选中会话被清除，页面显示空状态

#### Scenario: User cancels leave group
- **WHEN** 用户点击"退出群聊" → 确认对话框 → 点击"取消"
- **THEN** 对话框关闭，不执行任何删除操作

### Requirement: Three-tier cascade cleanup on leave
Backend `LeaveTask` handler SHALL 按序执行以下清理，确保不留残留：

1. 查询 task 下所有 sessionID
2. 逐一调用 AgentEnd `DELETE /v1/session/{id}`（终止进程，best-effort）
3. 调用 AgentEnd `DELETE /v1/workspace/task/{task_id}`（清理所有 worktree + Git 分支，best-effort）
4. 数据库事务删除：Task → Session → SessionAgent → Message → Announcement → DiffSnapshot → ContactGroupItem

AgentEnd 调用失败 SHALL NOT 阻断数据库删除。

#### Scenario: Full cleanup succeeds
- **WHEN** LeaveTask 被调用，AgentEnd 正常响应
- **THEN** 所有 Agent 进程被终止，所有 worktree 和 Git 分支被删除，数据库中所有关联记录被级联删除

#### Scenario: AgentEnd unavailable during leave
- **WHEN** LeaveTask 被调用，AgentEnd 调用超时或返回错误
- **THEN** 系统仍 SHALL 完成数据库删除，AgentEnd 侧的 worktree 由后台 `_inactive_cleanup_loop` 兜底清理

### Requirement: AgentEnd task-level cleanup API
AgentEnd SHALL 暴露 `DELETE /v1/workspace/task/{task_id}` 端点，封装已有的 `WorkspaceManager.cleanup_by_task()` 方法，清理该 task 下所有活跃 workspace（agent worktree + task-base worktree + 所有关联 Git 分支）。

#### Scenario: Cleanup all task workspaces
- **WHEN** `DELETE /v1/workspace/task/{task_id}` 被调用
- **THEN** 系统删除该 task 下所有活跃 worktree 目录、删除 `task/{task_id}` 分支、删除每个 `agent/{session_id}/{task_id}` 分支，返回 `{"cleaned": <count>}`

### Requirement: Frontend refresh after leave
退出群聊成功后，前端 SHALL 清除当前会话选中状态（`clearNavigation()`），并刷新会话列表（`invalidateQueries(['conversations'])`），被删除的会话从列表中消失。

#### Scenario: UI state after successful leave
- **WHEN** 退出群聊 API 返回成功
- **THEN** 当前选中会话被清除，聊天区域显示空状态，会话列表中不再显示该会话
