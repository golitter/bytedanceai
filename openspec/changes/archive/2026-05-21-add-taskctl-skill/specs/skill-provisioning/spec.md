## ADDED Requirements

### Requirement: Builtin skills directory structure
The system SHALL maintain a `src/skills/builtin/` directory containing all built-in skills. Each skill SHALL be a subdirectory containing at minimum an executable binary (`exe`) and an instruction file (`skill.md`).

#### Scenario: taskctl builtin structure
- **WHEN** the builtin skills directory is inspected
- **THEN** it SHALL contain `src/skills/builtin/taskctl/main.go`, `src/skills/builtin/taskctl/go.mod`, `src/skills/builtin/taskctl/exe` (pre-compiled binary), and `src/skills/builtin/taskctl/skill.md`

### Requirement: Skill provisioning on workspace creation
The provisioner SHALL copy all built-in skills to the agent's skill directory when a workspace is created. For ClaudeCode agents, the target directory SHALL be `{worktree_path}/.claude/skills/`. For OpenCode agents, the target SHALL be `{worktree_path}/.opencode/skills/`.

#### Scenario: Provisioning for ClaudeCode workspace
- **WHEN** a workspace is created with agent_name="claude_code"
- **THEN** the provisioner SHALL copy `src/skills/builtin/taskctl/exe` and `src/skills/builtin/taskctl/skill.md` to `{worktree_path}/.claude/skills/taskctl/`

#### Scenario: Provisioning for OpenCode workspace
- **WHEN** a workspace is created with agent_name="opencode"
- **THEN** the provisioner SHALL copy `src/skills/builtin/taskctl/exe` and `src/skills/builtin/taskctl/skill.md` to `{worktree_path}/.opencode/skills/taskctl/`

### Requirement: Shared directory initialization
The provisioner SHALL create the `shared/.agent/memory/` directory structure at the task level when a workspace is created. This includes `memory/common/` and `memory/{agent_name}/`.

#### Scenario: First agent in task creates shared structure
- **WHEN** the first workspace for task_id="task-123" is created with agent_name="claude_code"
- **THEN** the provisioner SHALL create `worktrees/task-123/shared/.agent/memory/common/` and `worktrees/task-123/shared/.agent/memory/claude_code/`

#### Scenario: Second agent in task gets its own memory directory
- **WHEN** a second workspace for task_id="task-123" is created with agent_name="opencode"
- **THEN** the provisioner SHALL create `worktrees/task-123/shared/.agent/memory/opencode/` without recreating common/ or claude_code/

### Requirement: skill.md is identical across agents
The `skill.md` file SHALL be identical for all agent types. It SHALL instruct the agent on how to invoke the exe and what commands are available.

#### Scenario: skill.md content
- **WHEN** skill.md is copied to any agent workspace
- **THEN** it SHALL contain instructions for using taskctl help, ls, summary, common-memory, and sub-memory commands, referencing the exe at `.skills/taskctl/exe` relative to the agent's skill directory
