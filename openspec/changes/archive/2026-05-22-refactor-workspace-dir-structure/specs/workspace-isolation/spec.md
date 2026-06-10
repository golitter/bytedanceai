## MODIFIED Requirements

### Requirement: GitOps atomic git operations
系统 SHALL 提供 `GitOps` 内部工具类，封装以下 git 原子操作：`worktree_add`、`worktree_remove`、`branch_create`、`add_and_commit`、`merge_branch`、`get_current_branch`。所有操作 MUST 通过 `asyncio.create_subprocess_exec` 异步执行 git 命令。

#### Scenario: Create worktree successfully
- **WHEN** 调用 `gitops.worktree_add(repo_path="/repos/project", path="/worktrees/task-1/sess-abc", branch="agent/sess-abc/task-1")`
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/sess-abc -b agent/sess-abc/task-1`，成功返回 `True`

#### Scenario: Create worktree on existing branch
- **WHEN** 目标 branch 已存在
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/sess-abc agent/sess-abc/task-1`（不带 `-b`），成功返回 `True`

#### Scenario: Remove worktree
- **WHEN** 调用 `gitops.worktree_remove("/worktrees/task-1/sess-abc")`
- **THEN** SHALL 执行 `git worktree remove /worktrees/task-1/sess-abc --force`，成功返回 `True`

#### Scenario: Git operation fails
- **WHEN** git 命令返回非零退出码
- **THEN** SHALL 返回 `False`，并记录 stderr 内容到日志
