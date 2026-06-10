## MODIFIED Requirements

### Requirement: Workspace data model
系统 SHALL 定义 `Workspace` 数据类，包含字段：`id`（str）、`task_id`（str）、`agent_name`（str）、`repo_path`（str）、`worktree_path`（str）、`branch_name`（str）、`session_id`（可选 str）、`status`（WorkspaceStatus 枚举）、`created_at`（datetime）、`container_id`（可选 str，默认 None）。

#### Scenario: Create workspace instance
- **WHEN** 创建 `Workspace(task_id="task-123", agent_name="frontend", repo_path="/repos/project")`
- **THEN** `branch_name` SHALL 自动生成为 `"agent/frontend/task-123"`，`status` 为 `ACTIVE`，`container_id` 为 `None`

#### Scenario: Create workspace with container_id
- **WHEN** 创建 `Workspace(..., container_id="docker-abc")`
- **THEN** `container_id` SHALL 为 `"docker-abc"`

### Requirement: WorkspaceManager create workspace
`WorkspaceManager` SHALL 提供 `create(repo_path, task_id, agent_name)` 方法，自动创建两级分支和 git worktree。MUST 先创建 `task/{task_id}` 集成分支（从 main），再创建 `agent/{agent_name}/{task_id}`（从 task branch）。

#### Scenario: Create workspace with new task branch
- **WHEN** 调用 `manager.create(repo_path="/repos/project", task_id="task-123", agent_name="frontend")` 且 task branch 不存在
- **THEN** SHALL 先创建 branch `task/task-123`（from main），再创建 branch `agent/frontend/task-123`（from task/task-123），创建 worktree 目录

#### Scenario: Create workspace with existing task branch
- **WHEN** 调用 `manager.create(...)` 且 `task/task-123` 已存在
- **THEN** SHALL 直接从 `task/task-123` 创建 `agent/{agent_name}/{task_id}`，不重复创建 task branch

#### Scenario: Create multiple workspaces for same task
- **WHEN** 同一 task_id 创建 frontend 和 backend 两个 workspace
- **THEN** 两个 workspace SHALL 从同一个 task branch 创建，拥有独立 agent branch 和目录

### Requirement: WorkspaceManager merge branch
`WorkspaceManager` SHALL 提供 `merge(workspace_id, target_branch)` 方法。默认将 agent branch 合并到对应的 task branch（`task/{task_id}`）。MUST 支持 `target_branch` 参数指定合并到 main。

#### Scenario: Merge agent branch to task branch (default)
- **WHEN** 调用 `merge(workspace_id)` 不指定 target_branch
- **THEN** SHALL 将 `agent/{agent_name}/{task_id}` 合并到 `task/{task_id}`，成功后 status 不变（仍为 ACTIVE）

#### Scenario: Merge task branch to main (explicit)
- **WHEN** 调用 `merge(workspace_id, target_branch="main")` 且 workspace 是 task branch 的代理
- **THEN** SHALL 在主仓库执行 `git checkout main && git merge task/{task_id}`，成功后 status 为 MERGED

#### Scenario: Merge with conflicts
- **WHEN** merge 产生冲突
- **THEN** SHALL 执行 `git merge --abort` 回滚，返回 `False`

### Requirement: WorkspaceManager concurrent operation safety
`WorkspaceManager` MUST 使用 per-task-group `asyncio.Lock` 保护 create、merge、cleanup 操作。同一 task_id 下的写操作串行化，不同 task_id 之间并行。

#### Scenario: Concurrent create for same task
- **WHEN** 两个请求同时为同一 task_id 创建不同 agent workspace
- **THEN** SHALL 串行执行，第一个完成后再执行第二个

#### Scenario: Concurrent create for different tasks
- **WHEN** 两个请求分别为不同 task_id 创建 workspace
- **THEN** SHALL 并行执行，互不阻塞

#### Scenario: Lock cleanup on workspace cleanup
- **WHEN** 一个 task_id 下所有 workspace 都已 cleanup
- **THEN** SHALL 删除对应的 Lock 对象，避免内存泄漏

### Requirement: WorkspaceManager store integration
`WorkspaceManager` MUST 通过 `WorkspaceStoreProtocol` 接口持久化所有 workspace 状态变更。create 时 save，cleanup/merge 时 update。

#### Scenario: Persist on create
- **WHEN** `create()` 成功创建 workspace
- **THEN** SHALL 调用 `store.save(workspace)` 持久化

#### Scenario: Persist on status change
- **WHEN** cleanup 或 merge 改变 workspace status
- **THEN** SHALL 调用 `store.save(workspace)` 更新持久化记录

### Requirement: WorkspaceManager cleanup on shutdown
App shutdown 时 SHALL 调用 `WorkspaceManager` 清理所有 ACTIVE workspace，并停止 TTL 后台 task。

#### Scenario: App shutdown cleanup
- **WHEN** FastAPI app lifespan shutdown
- **THEN** SHALL 取消 TTL task，遍历所有 ACTIVE workspace 执行 cleanup
