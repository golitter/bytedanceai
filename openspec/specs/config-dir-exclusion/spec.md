## ADDED Requirements

### Requirement: Exclude agent config directory from git tracking
WorkspaceManager 在创建 worktree 后 SHALL 将 agent 配置目录（`.claude` 或 `.opencode`）写入 worktree 的 `.git/info/exclude`，防止 agent 运行时生成的 untracked 文件被 `git add -A` 提交。

#### Scenario: ClaudeCode workspace excludes .claude directory
- **WHEN** 调用 `manager.create(..., agent_type=AgentType.CLAUDE_CODE)` 成功创建 worktree
- **THEN** worktree 的 `.git/info/exclude` 文件中 SHALL 包含 `/.claude`

#### Scenario: OpenCode workspace excludes .opencode directory
- **WHEN** 调用 `manager.create(..., agent_type=AgentType.OPENCODE)` 成功创建 worktree
- **THEN** worktree 的 `.git/info/exclude` 文件中 SHALL 包含 `/.opencode`

#### Scenario: Exclude does not affect tracked files
- **WHEN** 仓库已有 tracked 的 `.claude/CLAUDE.md` 文件，且 worktree 的 exclude 中包含 `/.claude`
- **THEN** 对 `.claude/CLAUDE.md` 的修改 SHALL 仍然被 `git status` 检测到并可正常提交

#### Scenario: Exclude prevents new untracked files from being staged
- **WHEN** agent 在 `.claude/memory/` 下生成了新文件 `notes.md`，且 worktree 的 exclude 中包含 `/.claude`
- **THEN** `git add -A` SHALL 不再暂存 `.claude/memory/notes.md`

### Requirement: Diff 生成层排除 skill 文件
在 `config-dir-exclusion` 已有的 git exclude 策略基础上，diff 生成层 SHALL 额外过滤 provisioner 写入的 skill 文件，确保即使文件未被 git exclude 覆盖也不会出现在 diff 输出中。

#### Scenario: git exclude 未覆盖的 skill 文件仍被 diff 过滤
- **WHEN** workspace 的 `.git/info/exclude` 未包含 agent config 目录（如 exclude 写入失败）
- **AND** `.codex/skills/render/render` 为 untracked 文件
- **THEN** `get_diff()` 返回的 diff 中 SHALL 仍不包含该文件
