## Context

当前 agentend 的清理策略是三层防线：TTL 自动过期（1h）、shutdown 全量清理、手动 destroy。这导致 agent 的工作产物（worktree、session 状态）无法跨服务重启保留。

三方数据流：
- **Frontend** → 用户操作入口
- **Backend** (Go/MySQL) → 数据库读写，维护 session status
- **Agentend** (Python) → 资源管理（worktree、进程），只读 DB

Session 状态已有自动同步机制：agent 进程完成后 agentend 会将 status（completed/error/interrupted）同步给 backend。

## Goals / Non-Goals

**Goals:**
- 用户通过前端手动停用 session 时，才触发后续清理
- Agentend 服务重启后，保留 workspace 和 session 状态
- Agentend 通过定时轮询 DB（只读）获取清理指令
- 清理粒度：session 级（清 worktree）+ task 级（全量清理，当 task 下所有 session 都 inactive）

**Non-Goals:**
- 不实现实时清理通知（WebSocket/回调），2h 轮询延迟可接受
- 不修改 agent 进程自身生命周期管理
- 不改变 workspace recovery 的启动恢复逻辑

## Decisions

### 1. 清理触发源：DB 轮询而非事件驱动

**选择**：Agentend 每 2h 读 DB 查询 inactive session
**备选**：Backend 通过 API/WebSocket 实时通知 agentend
**理由**：用户确认 2h 延迟可接受；轮询实现简单，不引入新的跨端通信依赖；agentend 保持只读 DB 的约束

### 2. Task 级清理条件：全 inactive 而非 "无 running"

**选择**：task 下所有 session 都 == inactive 才触发 cleanup_by_task
**备选**：task 下无 running session（completed/error 也算可清理）
**理由**：用户明确要求只清用户手动停用的；completed/error 的 session 代表有价值的执行结果，应由用户决定何时释放

### 3. Shutdown 时不清理资源

**选择**：lifespan shutdown 只停轮询任务，不调用 cleanup/destroy
**备选**：保留 shutdown 清理作为安全网
**理由**：workspace 有 JSON 持久化 + worktree 是 git 管理的物理目录，重启后可恢复；session 状态已同步给 backend

### 4. Session status 扩展：新增 inactive 值

**选择**：在现有状态枚举（idle/running/completed/interrupted/error）基础上新增 inactive
**理由**：inactive 语义明确（用户手动停用），与 agent 自然结束的状态区分清晰；string 类型字段无需改 DB schema

## Risks / Trade-offs

- **[磁盘占用]** workspace 不自动清理，长期运行可能积累大量 worktree → 用户通过前端定期停用不再需要的 session/task 来释放
- **[2h 窗口]** 停用后最长 2h 才真正释放资源 → 用户已确认可接受
- **[DB 轮询开销]** 每次 2h 全量扫描 session 表 → 数据量小时可忽略；数据量大时可加 WHERE status='inactive' 条件索引优化
- **[重启后 session 内存状态丢失]** session 运行时状态（history、process）在内存中，重启后无法恢复 → workspace 可恢复，session 的 process 需要重新建立；这是已有的限制，不在本次范围内解决
