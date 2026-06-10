## ADDED Requirements

### Requirement: merge command
taskctl SHALL provide a `merge` command that merges the current agent branch into the task branch. The target SHALL always be `task/{taskID}`, derived from parsePath. The command SHALL NOT accept any flags or arguments that change the merge target.

#### Scenario: Successful merge
- **WHEN** exe with taskID="task-123" and sessionID="sess-abc" invokes `taskctl merge`
- **THEN** exe SHALL execute: commit any uncommitted changes, checkout `task/task-123`, merge `agent/sess-abc/task-123`, checkout back to `agent/sess-abc/task-123`, print "merged to task/task-123", exit 0

#### Scenario: Auto-commit before merge
- **WHEN** exe invokes `taskctl merge` and working tree has uncommitted changes
- **THEN** exe SHALL first run `git add -A && git commit -m "auto: merge前自动提交"` before proceeding with merge

#### Scenario: No uncommitted changes
- **WHEN** exe invokes `taskctl merge` and working tree is clean
- **THEN** exe SHALL skip the auto-commit step and proceed directly to merge

#### Scenario: Merge conflict
- **WHEN** exe invokes `taskctl merge` and git merge reports conflicts
- **THEN** exe SHALL run `git merge --abort`, checkout back to the agent branch, print error to stderr, exit 1

#### Scenario: Branch name derivation
- **WHEN** exe with taskID="task-123" and sessionID="sess-abc" runs merge
- **THEN** exe SHALL derive agent branch as `agent/sess-abc/task-123` and task branch as `task/task-123` from path components, without running git commands to determine branch names

#### Scenario: Always returns to agent branch
- **WHEN** exe invokes `taskctl merge` regardless of outcome (success or conflict)
- **THEN** exe SHALL ensure the worktree is checked out on the agent branch `agent/{sessionID}/{taskID}` before exiting
