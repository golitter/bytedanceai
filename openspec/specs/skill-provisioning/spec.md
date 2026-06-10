## MODIFIED Requirements

### Requirement: GitOps atomic git operations
系统 SHALL 提供 `GitOps` 内部工具类，封装以下 git 原子操作：`worktree_add`、`worktree_remove`、`branch_create`、`add_and_commit`、`merge_branch`、`get_current_branch`、`write_exclude`。所有操作 MUST 通过 `asyncio.create_subprocess_exec` 异步执行 git 命令。`write_exclude` 为同步文件操作。

#### Scenario: Create worktree successfully
- **WHEN** 调用 `gitops.worktree_add(repo_path="/repos/project", path="/worktrees/task-1/sess-abc", branch="agent/sess-abc/task-1")`
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/sess-abc -b agent/sess-abc/task-1`，成功返回 `True`

#### Scenario: Create worktree on existing branch
- **WHEN** 目标 branch 已存在
- **THEN** SHALL 执行 `git worktree add /worktrees/task-1/sess-abc agent/sess-abc/task-1`（不带 `-b`），成功返回 `True`

#### Scenario: Remove worktree
- **WHEN** 调用 `gitops.worktree_remove("/worktrees/task-1/sess-abc")`
- **THEN** SHALL 执行 `git worktree remove /worktrees/task-1/sess-abc --force`，成功返回 `True`

#### Scenario: Write exclude entries to git info/exclude
- **WHEN** 调用 `gitops.write_exclude(worktree_path="/worktrees/task-1/sess-abc", entries=["/.claude"])`
- **THEN** SHALL 解析 worktree 的 `.git` 文件定位真正的 git dir，将 `/.claude` 追加到 `info/exclude`（不重复追加）

#### Scenario: Write exclude resolves worktree git dir
- **WHEN** worktree 的 `.git` 是一个文件内容为 `gitdir: /repos/project/.git/worktrees/sess-abc`
- **THEN** SHALL 解析出 git dir 为 `/repos/project/.git/worktrees/sess-abc`，在其中的 `info/exclude` 写入

#### Scenario: Write exclude for plain repo (not worktree)
- **WHEN** `.git` 是一个目录
- **THEN** SHALL 直接在 `.git/info/exclude` 写入

#### Scenario: Git operation fails
- **WHEN** git 命令返回非零退出码
- **THEN** SHALL 返回 `False`，并记录 stderr 内容到日志

## REMOVED Requirements

### Requirement: SkillProvisioner._write_git_exclude
**Reason**: Exclude 逻辑已移至 WorkspaceManager，provisioner 不再负责 git exclude 管理。
**Migration**: WorkspaceManager.create() 中直接调用 GitOps.write_exclude()，provisioner 不再调用 `_write_git_exclude`。
