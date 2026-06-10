## ADDED Requirements

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

### Requirement: Git worktree list parsing
系统 SHALL 实现 `parse_worktree_list(output: str)` 函数，解析 `git worktree list --porcelain` 输出为结构化列表。

#### Scenario: Parse multiple worktrees
- **WHEN** 输出包含多个 worktree 条目
- **THEN** SHALL 解析出每个 worktree 的 path 和 branch，返回 `list[tuple[str, str]]`

#### Scenario: Parse empty output
- **WHEN** `git worktree list` 只返回主仓库（无额外 worktree）
- **THEN** SHALL 返回空列表

### Requirement: Recovery logging
恢复流程 SHALL 记录每次恢复操作的结果：恢复了多少 workspace、标记了多少 CLEANED、清理了多少 orphan。

#### Scenario: Recovery summary logged
- **WHEN** 恢复流程完成
- **THEN** SHALL 输出日志："Workspace recovery: X recovered, Y cleaned, Z orphans removed"

### Requirement: Shutdown preservation
系统 SHALL 在 FastAPI lifespan shutdown 时保留所有 workspace 和 session 状态，不执行资源清理。Workspace 清理职责由 `inactive-cleanup` 能力承担。

#### Scenario: Shutdown preserves workspaces
- **WHEN** FastAPI app shutdown
- **THEN** SHALL 不调用 cleanup 或 destroy，保留所有 workspace 目录和 session 状态

#### Scenario: Workspace survives restart
- **WHEN** agentend 重启后执行 workspace recovery
- **THEN** SHALL 通过 reconciliation 流程恢复之前保留的 workspace
