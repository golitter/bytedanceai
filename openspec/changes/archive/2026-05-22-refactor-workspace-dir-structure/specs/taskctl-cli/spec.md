## MODIFIED Requirements

### Requirement: Path self-bootstrapping
taskctl exe SHALL determine its task_id and session_id by parsing its own filesystem path. The exe is located at `worktrees/{task_id}/{session_id}/.{agent_type}/skills/taskctl/exe`. It SHALL resolve shared directory as `worktrees/{task_id}/shared/.agent/`. It SHALL determine agent_type by inspecting which config directory (`.claude` or `.opencode`) contains the skills.

#### Scenario: Successful path resolution
- **WHEN** exe is located at `/abs/worktrees/task-123/sess-abc/.claude/skills/taskctl/exe`
- **THEN** exe SHALL resolve task_id="task-123", session_id="sess-abc", agent_type="claude-code", shared_dir=`/abs/worktrees/task-123/shared/.agent/`

#### Scenario: Invalid path structure
- **WHEN** exe cannot find "worktrees" in its path ancestry
- **THEN** exe SHALL print an error message and exit with non-zero code

### Requirement: sub-memory command with session isolation
taskctl SHALL provide a `sub-memory` command that reads files under `shared/.agent/memory/{session_id}/`, where session_id is resolved from the exe's own path. The exe SHALL NOT expose memory belonging to other sessions.

#### Scenario: Read own sub-memory
- **WHEN** exe with session_id="sess-abc" invokes `taskctl sub-memory`
- **THEN** exe SHALL print all files in `shared/.agent/memory/sess-abc/` concatenated, sorted by filename

#### Scenario: No sub-memory for this session
- **WHEN** `taskctl sub-memory` is invoked and the session's memory directory is empty or missing
- **THEN** exe SHALL print "(no sub-memory)"

#### Scenario: Session cannot access other sessions' memory
- **WHEN** exe with session_id="sess-abc" is asked to read memory
- **THEN** exe SHALL only read from `memory/sess-abc/` and SHALL NOT read from `memory/sess-def/` or any other session directory
