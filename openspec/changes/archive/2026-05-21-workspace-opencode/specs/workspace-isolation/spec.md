## ADDED Requirements

### Requirement: GitOps atomic git operations
系统 SHALL 提供 `GitOps` 内部工具类，封装以下 git 原子操作：`worktree_add`、`worktree_remove`、`branch_create`、`add_and_commit`、`merge_branch`、`get_current_branch`。所有操作 MUST 通过 `asyncio.create_subprocess_exec` 异步执行 git 命令。

#### Scenario: Create worktree successfully
- **WHEN** 调用 `gitops.worktree_add(repo_path="/repos/project", path="/worktrees/task-1/frontend", branch="agent/frontend/task-1")`
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/frontend -b agent/frontend/task-1`，成功返回 `True`

#### Scenario: Create worktree on existing branch
- **WHEN** 目标 branch 已存在
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/frontend agent/frontend/task-1`（不带 `-b`），成功返回 `True`

#### Scenario: Remove worktree
- **WHEN** 调用 `gitops.worktree_remove("/worktrees/task-1/frontend")`
- **THEN** SHALL 执行 `git worktree remove /worktrees/task-1/frontend --force`，成功返回 `True`

#### Scenario: Git operation fails
- **WHEN** git 命令返回非零退出码
- **THEN** SHALL 返回 `False`，并记录 stderr 内容到日志

### Requirement: Workspace data model
系统 SHALL 定义 `Workspace` 数据类，包含字段：`id`（str）、`task_id`（str）、`agent_name`（str）、`repo_path`（str）、`worktree_path`（str）、`branch_name`（str）、`session_id`（可选 str）、`status`（WorkspaceStatus 枚举）、`created_at`（datetime）。

#### Scenario: Create workspace instance
- **WHEN** 创建 `Workspace(task_id="task-123", agent_name="frontend", repo_path="/repos/project")`
- **THEN** `branch_name` SHALL 自动生成为 `"agent/frontend/task-123"`，`status` 为 `ACTIVE`

### Requirement: WorkspaceStatus enum
`WorkspaceStatus` SHALL 定义枚举值：`ACTIVE`、`MERGED`、`CLEANED`。

#### Scenario: Initial workspace status
- **WHEN** workspace 创建
- **THEN** `status` SHALL 为 `ACTIVE`

### Requirement: WorkspaceManager create workspace
`WorkspaceManager` SHALL 提供 `create(repo_path, task_id, agent_name)` 方法，自动创建 git worktree。Worktree 目录 MUST 位于 `{repo_path}/../worktrees/{task_id}/{agent_name}/`。

#### Scenario: Create workspace for agent task
- **WHEN** 调用 `manager.create(repo_path="/repos/project", task_id="task-123", agent_name="frontend")`
- **THEN** SHALL 创建 branch `agent/frontend/task-123`，创建 worktree 目录，返回 `Workspace` 对象

#### Scenario: Create multiple workspaces for same task
- **WHEN** 同一 task_id 创建 frontend 和 backend 两个 workspace
- **THEN** 两个 workspace SHALL 拥有独立 branch 和目录：`agent/frontend/task-123` 和 `agent/backend/task-123`

### Requirement: WorkspaceManager CRUD operations
`WorkspaceManager` SHALL 提供 `get`、`list`、`cleanup`、`cleanup_by_task` 方法。

#### Scenario: Get workspace by id
- **WHEN** 调用 `get(workspace_id)` 且存在
- **THEN** SHALL 返回对应 `Workspace` 对象

#### Scenario: List all workspaces
- **WHEN** 调用 `list()`
- **THEN** SHALL 返回所有 `Workspace` 对象列表

#### Scenario: Cleanup single workspace
- **WHEN** 调用 `cleanup(workspace_id)`
- **THEN** SHALL 执行 `git worktree remove`，将 status 设为 `CLEANED`，返回 `True`

#### Scenario: Cleanup all workspaces for a task
- **WHEN** 调用 `cleanup_by_task(task_id)`
- **THEN** SHALL 清理该 task 下所有 ACTIVE workspace，返回清理数量

### Requirement: WorkspaceManager commit changes
`WorkspaceManager` SHALL 提供 `commit(workspace_id, message)` 方法，在 worktree 目录执行 `git add . && git commit`。

#### Scenario: Commit workspace changes
- **WHEN** 调用 `commit(workspace_id, "feat: login page")`
- **THEN** SHALL 在 worktree 目录执行 `git add -A && git commit -m "feat: login page"`，成功返回 `True`

#### Scenario: Commit with no changes
- **WHEN** worktree 目录无文件变更
- **THEN** SHALL 返回 `False`，不执行 commit

### Requirement: WorkspaceManager merge branch
`WorkspaceManager` SHALL 提供 `merge(workspace_id, target_branch)` 方法，将 workspace 的 branch 合并到目标 branch（默认 `main`）。

#### Scenario: Merge workspace branch to main
- **WHEN** 调用 `merge(workspace_id, "main")`
- **THEN** SHALL 在主仓库执行 `git checkout main && git merge agent/{name}/{task_id}`，成功返回 `True`

#### Scenario: Merge with conflicts
- **WHEN** merge 产生冲突
- **THEN** SHALL 执行 `git merge --abort` 回滚，返回 `False`，错误信息包含 "merge conflict"

### Requirement: WorkspaceManager cleanup on shutdown
App shutdown 时 SHALL 调用 `WorkspaceManager.cleanup_by_task` 清理所有 ACTIVE workspace。

#### Scenario: App shutdown cleanup
- **WHEN** FastAPI app lifespan shutdown
- **THEN** SHALL 遍历所有 ACTIVE workspace 执行 cleanup
