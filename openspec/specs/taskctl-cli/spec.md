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

When invoked without additional arguments, it SHALL print all files concatenated, sorted by filename. When invoked with a filename argument, it SHALL print only the content of that single file.

#### Scenario: Read all sub-memory files
- **WHEN** exe with session_id="sess-abc" invokes `taskctl sub-memory`
- **THEN** exe SHALL print all files in `shared/.agent/memory/sess-abc/` concatenated, sorted by filename

#### Scenario: Read single sub-memory file
- **WHEN** exe with session_id="sess-abc" invokes `taskctl sub-memory notes.md`
- **THEN** exe SHALL print only the content of `shared/.agent/memory/sess-abc/notes.md`

#### Scenario: Specified file not found
- **WHEN** exe invokes `taskctl sub-memory nonexistent.md` and the file does not exist
- **THEN** exe SHALL print an error message to stderr and exit with non-zero code

#### Scenario: No sub-memory for this session
- **WHEN** `taskctl sub-memory` is invoked and the session's memory directory is empty or missing
- **THEN** exe SHALL print "(no sub-memory)"

#### Scenario: Session cannot access other sessions' memory
- **WHEN** exe with session_id="sess-abc" is asked to read memory
- **THEN** exe SHALL only read from `memory/sess-abc/` and SHALL NOT read from `memory/sess-def/` or any other session directory

### Requirement: common-memory command
taskctl SHALL provide a `common-memory` command that reads files under `shared/.agent/memory/common/`.

When invoked without additional arguments, it SHALL print all files concatenated, sorted by filename. When invoked with a filename argument, it SHALL print only the content of that single file.

#### Scenario: Read all common memory files
- **WHEN** `taskctl common-memory` is invoked
- **THEN** exe SHALL print all files in `shared/.agent/memory/common/` concatenated, sorted by filename

#### Scenario: Read single common memory file
- **WHEN** `taskctl common-memory config.md` is invoked
- **THEN** exe SHALL print only the content of `shared/.agent/memory/common/config.md`

#### Scenario: Specified file not found in common memory
- **WHEN** `taskctl common-memory nonexistent.md` is invoked and the file does not exist
- **THEN** exe SHALL print an error message to stderr and exit with non-zero code

#### Scenario: No common memory
- **WHEN** `taskctl common-memory` is invoked and the common memory directory is empty or missing
- **THEN** exe SHALL print "(no common memory)"

### Requirement: write-sub-memory command
taskctl SHALL provide a `write-sub-memory` command that writes content to a file under `shared/.agent/memory/{session_id}/`. The content SHALL be read from stdin if stdin has data (non-character-device), otherwise from remaining positional arguments joined by spaces. The write SHALL be atomic (write to temp file then rename).

#### Scenario: Write from positional arguments
- **WHEN** exe invokes `taskctl write-sub-memory notes.md hello world`
- **THEN** exe SHALL write "hello world" to `shared/.agent/memory/{session_id}/notes.md`

#### Scenario: Write from stdin
- **WHEN** exe pipes content: `echo "hello world" | taskctl write-sub-memory notes.md`
- **THEN** exe SHALL write the stdin content to `shared/.agent/memory/{session_id}/notes.md`

#### Scenario: No content provided
- **WHEN** exe invokes `taskctl write-sub-memory notes.md` with no stdin data and no content arguments
- **THEN** exe SHALL print an error message to stderr and exit with non-zero code

#### Scenario: Atomic write on crash
- **WHEN** write operation is interrupted mid-write
- **THEN** the target file SHALL either contain the previous content or the complete new content, never a partial write

## ADDED Requirements

### Requirement: Consistent error exit codes
All taskctl commands SHALL exit with code 0 on success and non-zero on failure. Read failures (file not found, directory read error) SHALL exit with code 1, same as write failures.

#### Scenario: Read failure exits non-zero
- **WHEN** any read command encounters an error (missing file, permission denied)
- **THEN** exe SHALL exit with code 1

#### Scenario: Successful operation exits zero
- **WHEN** any command completes successfully
- **THEN** exe SHALL exit with code 0

### Requirement: parsePath returns only used values
`parsePath` SHALL return `taskID`, `sessionID`, and `sharedDir`. The return signature SHALL be `(taskID, sessionID, sharedDir string, err error)` with exactly four return values.

#### Scenario: parsePath signature
- **WHEN** `parsePath` is called
- **THEN** it SHALL return `(taskID, sessionID, sharedDir string, err error)` with exactly four return values

#### Scenario: parsePath extracts taskID
- **WHEN** exe is located at `/abs/worktrees/task-123/sess-abc/.claude/skills/taskctl/exe`
- **THEN** parsePath SHALL return taskID="task-123", sessionID="sess-abc", sharedDir=`/abs/worktrees/task-123/shared/.agent/`
