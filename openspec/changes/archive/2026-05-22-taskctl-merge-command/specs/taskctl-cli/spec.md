## MODIFIED Requirements

### Requirement: parsePath returns only used values
`parsePath` SHALL return `taskID`, `sessionID`, and `sharedDir`. The return signature SHALL be `(taskID, sessionID, sharedDir string, err error)` with exactly four return values.

#### Scenario: parsePath signature
- **WHEN** `parsePath` is called
- **THEN** it SHALL return `(taskID, sessionID, sharedDir string, err error)` with exactly four return values

#### Scenario: parsePath extracts taskID
- **WHEN** exe is located at `/abs/worktrees/task-123/sess-abc/.claude/skills/taskctl/exe`
- **THEN** parsePath SHALL return taskID="task-123", sessionID="sess-abc", sharedDir=`/abs/worktrees/task-123/shared/.agent/`
