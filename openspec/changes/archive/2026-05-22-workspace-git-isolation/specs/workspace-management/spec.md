## MODIFIED Requirements

### Requirement: WorkspaceManager merge branch
`WorkspaceManager` SHALL 提供 `merge(workspace_id, target_branch)` 方法。默认将 agent branch 合并到对应的 task branch（`task/{task_id}`）。MUST 支持 `target_branch` 参数指定合并到 main。

#### Scenario: Merge agent branch to task branch (default)
- **WHEN** 调用 `merge(workspace_id)` 不指定 target_branch
- **THEN** SHALL 将 `agent/{session_id}/{task_id}` 合并到 `task/{task_id}`，成功后 status 不变（仍为 ACTIVE）

#### Scenario: Merge agent branch to task branch via API
- **WHEN** 通过 API `POST /v1/workspace/{workspace_id}/merge` 不传 target_branch 字段
- **THEN** SHALL 将 agent branch 合并到 `task/{task_id}`（不传参时 target_branch 为 None，由 manager fallback 到 task branch）

#### Scenario: Merge task branch to main (explicit)
- **WHEN** 调用 `merge(workspace_id, target_branch="main")`
- **THEN** SHALL 在主仓库执行 `git checkout main && git merge task/{task_id}`，成功后 status 为 MERGED

#### Scenario: Merge with conflicts
- **WHEN** merge 产生冲突
- **THEN** SHALL 执行 `git merge --abort` 回滚，返回 `False`

## ADDED Requirements

### Requirement: WorkspaceManager merge task branch to main
`WorkspaceManager` SHALL 提供 `merge_task_to_main(repo_path, task_id)` 方法，将 `task/{task_id}` 分支合并到 `main`。此方法为平台级操作，不关联特定 workspace。

#### Scenario: Merge task branch to main successfully
- **WHEN** 调用 `merge_task_to_main(repo_path="/repos/project", task_id="task-123")`
- **THEN** SHALL 在主仓库执行 `git checkout main && git merge task/task-123`，成功返回 `True`

#### Scenario: Merge task branch to main with conflicts
- **WHEN** `task/{task_id}` 合并到 main 产生冲突
- **THEN** SHALL 执行 `git merge --abort` 回滚，返回 `False`

### Requirement: API endpoint for task-to-main merge
API SHALL 提供 `POST /v1/workspace/task/{task_id}/merge-to-main` 端点，接收 `repo_path` 参数，调用 `WorkspaceManager.merge_task_to_main()`。

#### Scenario: Call task-to-main merge endpoint
- **WHEN** 调用 `POST /v1/workspace/task/task-123/merge-to-main` 且 body 包含 `{"repo_path": "/repos/project"}`
- **THEN** SHALL 调用 `manager.merge_task_to_main("/repos/project", "task-123")`，返回 `{"success": true/false}`

#### Scenario: Task-to-main merge with conflict
- **WHEN** merge 产生冲突
- **THEN** SHALL 返回 `{"success": false, "error": "merge conflict"}`

### Requirement: MergeRequest default target branch
`MergeRequest` 的 `target_branch` 字段 SHALL 默认为 `None`（而非 `"main"`），使 manager 的 fallback 逻辑（默认合到 `task/{task_id}`）生效。

#### Scenario: API merge without target_branch
- **WHEN** 调用 `POST /v1/workspace/{id}/merge` body 为 `{}`（不传 target_branch）
- **THEN** target_branch SHALL 为 `None`，manager SHALL 合并到 `task/{task_id}`
