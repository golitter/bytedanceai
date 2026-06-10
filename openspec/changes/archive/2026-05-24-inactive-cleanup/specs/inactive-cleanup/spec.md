## ADDED Requirements

### Requirement: Inactive cleanup background task
Agentend SHALL 在 lifespan startup 时启动一个后台 asyncio 定时任务，每 2 小时轮询 DB（只读），查询 status == "inactive" 的 session，执行清理。

#### Scenario: Cleanup inactive sessions
- **WHEN** 轮询发现 session status == "inactive" 且对应 workspace 存在
- **THEN** SHALL 销毁该 session 对应的 worktree 资源

#### Scenario: Skip active sessions
- **WHEN** 轮询发现 session status 不为 "inactive"（如 running/completed/error）
- **THEN** SHALL 不执行任何清理操作

#### Scenario: Stop cleanup on shutdown
- **WHEN** FastAPI app shutdown
- **THEN** SHALL 仅取消后台轮询 task，不清理任何 session 或 workspace 资源

### Requirement: Task-level cleanup when all sessions inactive
Agentend SHALL 在每次轮询时检查每个 task 的 session 状态，当 task 下所有 session 都为 inactive 时，执行 cleanup_by_task。

#### Scenario: All sessions inactive
- **WHEN** task_abc 下所有 session status 都 == "inactive"
- **THEN** SHALL 执行 cleanup_by_task，清理该 task 的所有 worktree 和分支

#### Scenario: Mixed session states
- **WHEN** task_abc 下有 session 为 completed，有 session 为 inactive
- **THEN** SHALL 只清理 inactive session 的 worktree，不触发 task 级 cleanup

### Requirement: DB read-only constraint
Agentend 轮询 DB 时 MUST 只执行 SELECT 查询，不得对 DB 执行任何写操作（INSERT/UPDATE/DELETE）。session status 的变更由 Backend API 负责。

#### Scenario: Query sessions
- **WHEN** 轮询任务执行
- **THEN** SHALL 只执行 `SELECT session_id, task_id, status FROM sessions` 类查询

### Requirement: Cleanup interval configuration
轮询间隔 SHALL 可通过 config.yaml 配置，默认值 7200 秒（2 小时）。

#### Scenario: Use default interval
- **WHEN** config.yaml 未配置 cleanup_interval
- **THEN** SHALL 使用 7200 秒作为轮询间隔

#### Scenario: Use custom interval
- **WHEN** config.yaml 配置 cleanup_interval=3600
- **THEN** SHALL 使用 3600 秒作为轮询间隔

### Requirement: Cleanup logging
每次轮询清理 SHALL 记录：扫描了多少 session、清理了多少 inactive session、清理了多少 task。

#### Scenario: Cleanup cycle logged
- **WHEN** 轮询清理周期执行
- **THEN** SHALL 输出日志："Inactive cleanup: scanned X sessions, cleaned Y sessions, cleaned Z tasks"
