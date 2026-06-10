## ADDED Requirements

### Requirement: Workspace list with disk usage
工作区管理页面 SHALL 展示所有工作区列表，包含磁盘占用信息。

#### Scenario: Display workspaces
- **WHEN** 用户打开工作区管理页面
- **THEN** 显示工作区表格，包含列：工作区 ID、任务、Agent、分支名、磁盘占用（MB + 进度条）、状态、操作

#### Scenario: Display statistics bar
- **WHEN** 工作区管理页面加载
- **THEN** 顶部显示统计条：总数、活跃数、已清理数、总磁盘占用

### Requirement: Clean workspace
用户 SHALL 能清理活跃或已合并的工作区。

#### Scenario: Clean single workspace
- **WHEN** 用户点击某活跃工作区的清理按钮
- **THEN** 系统删除该工作区的 git worktree 和分支，更新状态为 cleaned，磁盘占用清零，显示 toast 通知

#### Scenario: Batch clean merged workspaces
- **WHEN** 用户点击"清理已合并工作区"按钮
- **THEN** 系统批量清理所有已合并状态的工作区，显示 toast 通知清理数量

### Requirement: Workspace API
后端 SHALL 提供工作区管理相关 API。

#### Scenario: List workspaces
- **WHEN** 前端请求 `GET /api/admin/workspaces`
- **THEN** 返回工作区列表，每项包含 id、task、agent、branch、diskMB、status

#### Scenario: Delete workspace
- **WHEN** 前端发送 `DELETE /api/admin/workspaces/:id`
- **THEN** 后端删除对应 git worktree 和分支，返回 `{ success: true }`
