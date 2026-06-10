## MODIFIED Requirements

### Requirement: Workspace creation provisions skills
WorkspaceManager.create() SHALL call the skill provisioner after creating the git worktree. The provisioner SHALL receive the workspace path, task_id, session_id, and agent_type to determine the correct skill target directory and shared directory location. WorkspaceManager SHALL 在 provisioner 调用之后写入 agent 配置目录的 git exclude 规则。

#### Scenario: Create workspace with skill provisioning
- **WHEN** WorkspaceManager.create(repo_path, task_id="task-123", session_id="sess-abc", agent_type=AgentType.CLAUDE_CODE) is called
- **THEN** after the git worktree is created, the provisioner SHALL be called with (worktree_path, task_id, session_id, agent_type) to distribute builtin skills and initialize shared directories, 然后 WorkspaceManager SHALL 调用 GitOps.write_exclude 将 `/.claude` 写入 worktree 的 `.git/info/exclude`

#### Scenario: Workspace creation failure does not partially provision
- **WHEN** WorkspaceManager.create() fails during worktree creation
- **THEN** the provisioner SHALL NOT be invoked, exclude 规则也 SHALL NOT 被写入
