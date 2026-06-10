## Why

当前 agentend 在服务 shutdown 时会无条件销毁所有 session 和 workspace，同时通过 TTL 定时清理过期 workspace。这导致 agent 运行产物无法在重启后保留，用户无法在稍后手动决定何时释放资源。需要将清理策略从"自动全清"改为"用户手动停用后才清理"。

## What Changes

- Backend 新增 `PATCH /api/sessions/:sessionId` API，支持将 session status 设为 `inactive`
- Agentend 去掉 shutdown 时的全量 cleanup 逻辑
- Agentend 去掉 TTL 定时清理，替换为每 2h 轮询 DB（只读）的懒清理机制
- Agentend 懒清理规则：只清理 status == `inactive` 的 session 的 worktree；只清理 task 下所有 session 都 == `inactive` 的 task
- Frontend session 列表新增"停用"操作入口

## Capabilities

### New Capabilities
- `inactive-cleanup`: 基于 inactive 状态的懒清理机制——agentend 定时轮询 DB，按规则清理已被用户标记为 inactive 的 session 和 task

### Modified Capabilities
- `session-manager`: session 新增 `inactive` 状态，由前端用户手动触发，backend API 负责写入
- `workspace-ttl`: 去掉 TTL 自动过期，改为由 inactive-cleanup 驱动清理
- `workspace-recovery`: 启动时恢复 workspace 的逻辑不变，但 shutdown 时不再清理

## Impact

- **Backend**: 新增 PATCH session status API，session model 需支持 `inactive` 状态值
- **Agentend**: `lifespan` shutdown 逻辑简化（只停轮询任务，不清理资源）；新增定时 DB 轮询清理任务；去掉 TTL 相关配置和代码
- **Frontend**: session 列表 UI 加"停用"按钮
- **Database**: Session 表 status 字段新增 `inactive` 值（无需改 schema，string 类型已支持）
- **Contracts**: `session-state.yaml` 新增 `inactive` 状态及转换规则
