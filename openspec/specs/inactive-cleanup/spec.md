## ADDED Requirements

### Requirement: Inactive session polling
系统 SHALL 在 agentend FastAPI lifespan startup 时启动后台 asyncio task，每 2 小时轮询数据库查询 status == `inactive` 的 session 列表。轮询间隔 SHALL 可通过 `config.yaml` 的 `cleanup_interval` 配置，默认 7200 秒。

#### Scenario: Poll and find inactive sessions
- **WHEN** 后台轮询任务查询 DB，发现存在 status == `inactive` 的 session
- **THEN** SHALL 对每个 inactive session 执行 worktree 清理

#### Scenario: Poll and find no inactive sessions
- **WHEN** 后台轮询任务查询 DB，未发现 inactive session
- **THEN** SHALL 不执行任何清理操作，等待下一个轮询周期

#### Scenario: Custom cleanup interval
- **WHEN** `config.yaml` 配置 `cleanup_interval: 3600`
- **THEN** SHALL 使用 3600 秒作为轮询间隔

#### Scenario: Default cleanup interval
- **WHEN** `config.yaml` 未配置 `cleanup_interval`
- **THEN** SHALL 使用 7200 秒（2 小时）作为默认轮询间隔

### Requirement: Session-level cleanup
系统 SHALL 对 status == `inactive` 的 session 执行 worktree 清理。清理操作 SHALL 调用 `git worktree remove` 释放物理目录资源。

#### Scenario: Clean inactive session worktree
- **WHEN** session 的 status == `inactive` 且 worktree 目录存在
- **THEN** SHALL 执行 `git worktree remove` 清理 worktree，记录日志

#### Scenario: Skip session with no worktree
- **WHEN** session 的 status == `inactive` 但 worktree 目录不存在
- **THEN** SHALL 跳过清理，记录警告日志

### Requirement: Task-level cleanup
系统 SHALL 检查每个 task 下所有 session 的状态分布。仅当 task 下所有 session 都 == `inactive` 时，才触发 task 级全量清理。

#### Scenario: All sessions inactive — cleanup task
- **WHEN** task 下所有 session 的 status 均 == `inactive`
- **THEN** SHALL 执行 task 级清理（`cleanup_by_task`），释放该 task 关联的所有资源

#### Scenario: Some sessions not inactive — skip task cleanup
- **WHEN** task 下存在 status != `inactive` 的 session（如 completed、error）
- **THEN** SHALL 不执行 task 级清理，保留该 task 的资源

### Requirement: Shutdown behavior
系统 SHALL 在 FastAPI lifespan shutdown 时只停止轮询任务，不执行资源清理。

#### Scenario: Shutdown stops polling only
- **WHEN** FastAPI app shutdown
- **THEN** SHALL 取消 inactive-cleanup 后台轮询 task，不调用 cleanup 或 destroy

### Requirement: DB read-only constraint
agentend 的清理轮询任务 SHALL 仅通过只读 DB 查询获取清理指令。所有 session 状态变更 MUST 由 backend API 执行写入。

#### Scenario: Agentend reads inactive status from DB
- **WHEN** agentend 轮询 DB 查询 session 状态
- **THEN** SHALL 只执行 SELECT 查询，不修改 session 状态

#### Scenario: Backend writes inactive status
- **WHEN** 用户通过前端停用 session
- **THEN** 前端 SHALL 调用 backend `PATCH /api/sessions/:sessionId` API，由 backend 写入 `inactive` 状态

### Requirement: Cleanup logging
每次清理周期 SHALL 记录清理结果摘要。

#### Scenario: Cleanup cycle logged
- **WHEN** 一个清理周期完成
- **THEN** SHALL 输出日志："Inactive cleanup: X sessions cleaned, Y tasks cleaned, Z skipped"
