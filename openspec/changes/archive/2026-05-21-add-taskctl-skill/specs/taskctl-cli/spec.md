## ADDED Requirements

### Requirement: Path self-bootstrapping
taskctl exe SHALL determine its task_id and agent_name by parsing its own filesystem path. The exe is located at `worktrees/{task_id}/{agent_name}/.{agent_type}/skills/taskctl/exe`. It SHALL resolve shared directory as `worktrees/{task_id}/shared/.agent/`.

#### Scenario: Successful path resolution
- **WHEN** exe is located at `/abs/worktrees/task-123/claude_code/.claude/skills/taskctl/exe`
- **THEN** exe SHALL resolve task_id="task-123", agent_name="claude_code", shared_dir=`/abs/worktrees/task-123/shared/.agent/`

#### Scenario: Invalid path structure
- **WHEN** exe cannot find "worktrees" in its path ancestry
- **THEN** exe SHALL print an error message and exit with non-zero code

### Requirement: help command
taskctl SHALL provide a `help` command that prints usage information for all available commands.

#### Scenario: User runs help
- **WHEN** `taskctl help` is invoked
- **THEN** exe SHALL print all command names with brief descriptions: help, ls, summary, common-memory, sub-memory

### Requirement: ls command
taskctl SHALL provide an `ls` command that lists the file structure under `shared/.agent/`.

#### Scenario: List shared directory contents
- **WHEN** `taskctl ls` is invoked
- **THEN** exe SHALL recursively list all files and directories under `shared/.agent/` with relative paths

#### Scenario: Shared directory does not exist
- **WHEN** `taskctl ls` is invoked and `shared/.agent/` does not exist
- **THEN** exe SHALL print an error indicating the directory is not found

### Requirement: summary command
taskctl SHALL provide a `summary` command that outputs a task overview by reading `shared/.agent/config.yaml` and all files under `shared/.agent/plans/`.

#### Scenario: Summary with config and plans
- **WHEN** `taskctl summary` is invoked and config.yaml and plans/ exist
- **THEN** exe SHALL print config.yaml content followed by all plan files concatenated

#### Scenario: Summary with missing files
- **WHEN** `taskctl summary` is invoked and some files do not exist
- **THEN** exe SHALL print available files and skip missing ones silently

### Requirement: common-memory command
taskctl SHALL provide a `common-memory` command that reads all files under `shared/.agent/memory/common/`.

#### Scenario: Read common memory
- **WHEN** `taskctl common-memory` is invoked
- **THEN** exe SHALL print all files in `shared/.agent/memory/common/` concatenated, sorted by filename

#### Scenario: No common memory files
- **WHEN** `taskctl common-memory` is invoked and the directory is empty or missing
- **THEN** exe SHALL print "(no common memory)"

### Requirement: sub-memory command with agent isolation
taskctl SHALL provide a `sub-memory` command that reads files under `shared/.agent/memory/{agent_name}/`, where agent_name is resolved from the exe's own path. The exe SHALL NOT expose memory belonging to other agents.

#### Scenario: Read own sub-memory
- **WHEN** exe with agent_name="claude_code" invokes `taskctl sub-memory`
- **THEN** exe SHALL print all files in `shared/.agent/memory/claude_code/` concatenated, sorted by filename

#### Scenario: No sub-memory for this agent
- **WHEN** `taskctl sub-memory` is invoked and the agent's memory directory is empty or missing
- **THEN** exe SHALL print "(no sub-memory)"

#### Scenario: Agent cannot access other agents' memory
- **WHEN** exe with agent_name="claude_code" is asked to read memory
- **THEN** exe SHALL only read from `memory/claude_code/` and SHALL NOT read from `memory/opencode/` or any other agent directory
