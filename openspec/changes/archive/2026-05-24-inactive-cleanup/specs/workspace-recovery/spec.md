## MODIFIED Requirements

### Requirement: Startup workspace recovery
系统 SHALL 在启动时执行 workspace 恢复流程，调用 `git worktree list --porcelain` 扫描所有物理 worktree，与持久化存储中的 workspace 记录 reconcile。

#### Scenario: Store and worktree both exist — recover
- **WHEN** store 中有 workspace（status=ACTIVE）且 worktree 目录存在
- **THEN** SHALL 保持 workspace 为 ACTIVE 状态，加载到内存

#### Scenario: Store exists but worktree missing — mark cleaned
- **WHEN** store 中有 workspace（status=ACTIVE）但 worktree 目录不存在
- **THEN** SHALL 将 workspace status 更新为 CLEANED，更新 store

#### Scenario: Worktree exists but store missing — orphan cleanup
- **WHEN** `git worktree list` 发现 worktree 但 store 中无对应记录
- **THEN** SHALL 执行 `git worktree remove --force` 清理 orphan，记录警告日志

#### Scenario: No worktrees and empty store — clean start
- **WHEN** `git worktree list` 返回空且 store 为空
- **THEN** SHALL 正常启动，不执行任何恢复操作

## ADDED Requirements

### Requirement: No shutdown cleanup
系统 SHALL 在 FastAPI lifespan shutdown 时 NOT 执行任何 workspace 清理操作。只停止后台轮询任务。

#### Scenario: Shutdown preserves workspaces
- **WHEN** FastAPI app shutdown 且存在 ACTIVE workspace
- **THEN** SHALL 不调用 cleanup 或 cleanup_by_task，workspace 和 worktree 保留在磁盘

#### Scenario: Shutdown stops background tasks only
- **WHEN** FastAPI app shutdown
- **THEN** SHALL 仅停止 inactive-cleanup 轮询任务，不清理 session 或 workspace
