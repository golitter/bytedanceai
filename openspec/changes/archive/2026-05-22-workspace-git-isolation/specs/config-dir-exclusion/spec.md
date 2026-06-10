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
