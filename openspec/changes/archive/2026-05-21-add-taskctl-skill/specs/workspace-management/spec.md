## MODIFIED Requirements

### Requirement: Workspace creation provisions skills
WorkspaceManager.create() SHALL call the skill provisioner after creating the git worktree. The provisioner SHALL receive the workspace path, task_id, and agent_name to determine the correct skill target directory and shared directory location.

#### Scenario: Create workspace with skill provisioning
- **WHEN** WorkspaceManager.create(repo_path, task_id="task-123", agent_name="claude_code") is called
- **THEN** after the git worktree is created, the provisioner SHALL be called with (worktree_path, task_id, agent_name) to distribute builtin skills and initialize shared directories

#### Scenario: Workspace creation failure does not partially provision
- **WHEN** WorkspaceManager.create() fails during worktree creation
- **THEN** the provisioner SHALL NOT be invoked
