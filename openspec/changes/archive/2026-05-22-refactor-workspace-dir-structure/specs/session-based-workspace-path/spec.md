## ADDED Requirements

### Requirement: Session-based worktree path generation
系统 SHALL 使用 `session_id` 替代 `agent_name` 生成 worktree 路径。路径格式 SHALL 为 `{repo_parent}/worktrees/{task_id}/{session_id}`。

#### Scenario: Generate worktree path with session_id
- **WHEN** 调用 `_generate_worktree_path(repo_path="/repos/project", task_id="task-123", session_id="sess-abc")`
- **THEN** SHALL 返回 `/repos/worktrees/task-123/sess-abc`

#### Scenario: Same task different sessions get different paths
- **WHEN** 同一 task_id="task-123" 下创建 session_id="sess-1" 和 session_id="sess-2" 两个 workspace
- **THEN** 两个 worktree 路径 SHALL 分别为 `{base}/task-123/sess-1` 和 `{base}/task-123/sess-2`

### Requirement: Session-based branch name generation
系统 SHALL 使用 `session_id` 替代 `agent_name` 生成分支名。分支格式 SHALL 为 `agent/{session_id}/{task_id}`。

#### Scenario: Generate branch name with session_id
- **WHEN** 调用 `_generate_branch_name(session_id="sess-abc", task_id="task-123")`
- **THEN** SHALL 返回 `"agent/sess-abc/task-123"`
