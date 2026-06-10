## MODIFIED Requirements

### Requirement: Skill provisioning on workspace creation
The provisioner SHALL copy all built-in skills to the agent's skill directory when a workspace is created. For ClaudeCode agents, the target directory SHALL be `{worktree_path}/.claude/skills/`. For OpenCode agents, the target SHALL be `{worktree_path}/.opencode/skills/`. The provisioner SHALL receive `agent_type` (not `agent_name`) to determine the target directory.

#### Scenario: Provisioning for ClaudeCode workspace
- **WHEN** a workspace is created with agent_type=AgentType.CLAUDE_CODE
- **THEN** the provisioner SHALL copy `src/skills/builtin/taskctl/exe` and `src/skills/builtin/taskctl/skill.md` to `{worktree_path}/.claude/skills/taskctl/`

#### Scenario: Provisioning for OpenCode workspace
- **WHEN** a workspace is created with agent_type=AgentType.OPENCODE
- **THEN** the provisioner SHALL copy `src/skills/builtin/taskctl/exe` and `src/skills/builtin/taskctl/skill.md` to `{worktree_path}/.opencode/skills/taskctl/`

### Requirement: Shared directory initialization
The provisioner SHALL create the `shared/.agent/memory/` directory structure at the task level when a workspace is created. This includes `memory/common/` and `memory/{session_id}/`.

#### Scenario: First session in task creates shared structure
- **WHEN** the first workspace for task_id="task-123" is created with session_id="sess-abc"
- **THEN** the provisioner SHALL create `worktrees/task-123/shared/.agent/memory/common/` and `worktrees/task-123/shared/.agent/memory/sess-abc/`

#### Scenario: Second session in task gets its own memory directory
- **WHEN** a second workspace for task_id="task-123" is created with session_id="sess-def"
- **THEN** the provisioner SHALL create `worktrees/task-123/shared/.agent/memory/sess-def/` without recreating common/ or sess-abc/

### Requirement: Workspace creation provisions skills
WorkspaceManager.create() SHALL call the skill provisioner after creating the git worktree. The provisioner SHALL receive the workspace path, task_id, session_id, and agent_type to determine the correct skill target directory and shared directory location.

#### Scenario: Create workspace with skill provisioning
- **WHEN** WorkspaceManager.create(repo_path, task_id="task-123", session_id="sess-abc", agent_type=AgentType.CLAUDE_CODE) is called
- **THEN** after the git worktree is created, the provisioner SHALL be called with (worktree_path, task_id, session_id, agent_type) to distribute builtin skills and initialize shared directories

#### Scenario: Workspace creation failure does not partially provision
- **WHEN** WorkspaceManager.create() fails during worktree creation
- **THEN** the provisioner SHALL NOT be invoked
