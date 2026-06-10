## MODIFIED Requirements

### Requirement: POST /v1/agent/stream endpoint
系统 SHALL 提供 `POST /v1/agent/stream` 端点，接收 `AgentRequest`，通过 SSE 流式返回 `StreamEvent` 序列。执行前 MUST 检查 workspace：若 `workspace_path` 存在则直接使用，若 `repo_path` 存在则自动创建 workspace。

#### Scenario: Auto-create workspace from repo_path
- **WHEN** 请求包含 `repo_path="/repos/project"` 和 `task_id="task-123"`，无 `workspace_path`
- **THEN** SHALL 调用 `WorkspaceManager.create` 创建 workspace，将 `workspace_path` 传入 Adapter 执行

#### Scenario: Use existing workspace_path
- **WHEN** 请求包含 `workspace_path="/workspaces/task-1/frontend"`
- **THEN** SHALL 直接绑定到 session，跳过 workspace 创建

## ADDED Requirements

### Requirement: POST /v1/workspace/create endpoint
系统 SHALL 提供 `POST /v1/workspace/create` 端点，手动创建 workspace。

#### Scenario: Create workspace via API
- **WHEN** 发送 `{"repo_path": "/repos/project", "task_id": "task-123", "agent_name": "frontend"}`
- **THEN** SHALL 创建 git worktree，返回 `{"workspace_id": "...", "worktree_path": "...", "branch_name": "agent/frontend/task-123"}`

### Requirement: POST /v1/workspace/{id}/commit endpoint
系统 SHALL 提供 `POST /v1/workspace/{id}/commit` 端点，提交 workspace 的变更。

#### Scenario: Commit workspace changes
- **WHEN** 发送 `POST /v1/workspace/{id}/commit` body `{"message": "feat: login page"}`
- **THEN** SHALL 在 worktree 执行 git commit，返回 `{"success": true}`

### Requirement: POST /v1/workspace/{id}/merge endpoint
系统 SHALL 提供 `POST /v1/workspace/{id}/merge` 端点，将 workspace branch 合并到目标 branch。

#### Scenario: Merge to main
- **WHEN** 发送 `POST /v1/workspace/{id}/merge` body `{"target_branch": "main"}`
- **THEN** SHALL 执行 merge，成功返回 `{"success": true}`，冲突返回 `{"success": false, "error": "merge conflict"}`

### Requirement: DELETE /v1/workspace/{id} endpoint
系统 SHALL 提供 `DELETE /v1/workspace/{id}` 端点，清理 workspace（移除 worktree）。

#### Scenario: Cleanup workspace
- **WHEN** 发送 `DELETE /v1/workspace/{id}`
- **THEN** SHALL 执行 `git worktree remove`，返回 `{"success": true}`

### Requirement: GET /v1/workspace endpoint
系统 SHALL 提供 `GET /v1/workspace` 端点，列出所有 workspace。

#### Scenario: List workspaces
- **WHEN** 调用 `GET /v1/workspace`
- **THEN** SHALL 返回 workspace 列表，每个元素包含 workspace_id、task_id、agent_name、branch_name、status、worktree_path

### Requirement: Config OPENCODE_CLI_PATH
系统配置 SHALL 新增 `OPENCODE_CLI_PATH` 环境变量，默认值为 `"opencode"`。

#### Scenario: Use default OpenCode CLI path
- **WHEN** 未设置 `OPENCODE_CLI_PATH`
- **THEN** SHALL 使用 `"opencode"` 作为 CLI 路径
