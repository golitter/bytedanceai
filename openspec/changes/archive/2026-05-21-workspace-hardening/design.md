## Context

Workspace 模块 MVP（workspace-opencode change）已实现：`GitOps` 封装 git worktree 原子操作、`WorkspaceManager` 管理内存中的 workspace 生命周期、`Workspace` dataclass 描述状态。当前所有 Agent branch 直接从 main 创建并 merge 回 main（扁平结构），状态存储在 `dict[str, Workspace]` 内存字典中。

现有代码关键接口：
- `GitOps.worktree_add(repo_path, path, branch)` — 创建 worktree，已支持已有 branch 复用
- `WorkspaceManager.create(repo_path, task_id, agent_name)` — 创建 Workspace 并调用 GitOps
- `SessionMappingStore` — 已建立 JSON 文件持久化模式

约束：
- Python 3.10+，uv 管理依赖，不引入新外部依赖
- 延续 `SessionMappingStore` 的 JSON 文件存储模式
- `asyncio.create_subprocess_exec` 调用 git 命令
- 不改 Adapter 层执行逻辑

## Goals / Non-Goals

**Goals:**

- 两级分支：main → task/{task_id} → agent/{agent_name}/{task_id}，Agent merge 只到 task branch
- 持久化：workspace 状态存 JSON 文件，抽象接口可后续迁移 SQLite
- 并发安全：per-workspace asyncio.Lock 保护写操作
- 启动恢复：扫描 git worktree list 与存储 reconcile
- TTL 清理：后台定时清理过期 workspace

**Non-Goals:**

- 不做 Container/Docker 隔离（仅预留 container_id 字段）
- 不引入 SQLite（保持 JSON 文件模式）
- 不改 Adapter 层（ClaudeCodeAdapter、OpenCodeAdapter）
- 不做自动冲突解决
- 不做 workspace 历史版本/审计日志

## Decisions

### D1: 两级分支策略 — task branch 作为 Agent 的公共父级

**选择**: 引入 `task/{task_id}` 集成分支。

```
main
  └── task/task-123 (from main)
        ├── agent/frontend/task-123 (from task/task-123)
        └── agent/backend/task-123  (from task/task-123)

merge 流程：
  agent/frontend/task-123 → task/task-123
  agent/backend/task-123  → task/task-123
  task/task-123           → main (显式 API 调用)
```

**理由**: Agent branch 之间互相隔离，merge 到 task branch 是任务内集成。task → main 只做一次，由上层 Orchestrator 或 API 显式触发。比扁平结构更安全。

**替代方案**: 所有 Agent branch 直接 merge 到 main（当前方式）— 已证明有并发污染风险。

**创建时机**: 收到第一个 Agent 请求时自动创建 task branch + agent branch，对调用方透明。task branch 已存在时直接 checkout。

### D2: 持久化方案 — JSON 文件 + 抽象接口

**选择**: `WorkspaceStore` 使用 JSON 文件存储，和 `SessionMappingStore` 保持一致。定义 `WorkspaceStoreProtocol` 抽象接口。

**理由**: 项目无 DB 依赖，workspace 量级在十级。JSON 文件零依赖、可读、易调试。抽象接口让后续迁移 SQLite 只需实现新的 Store 类。

**存储路径**: `logs/workspaces.json`（与 `logs/session_mappings.json` 对称）

**接口设计**:
```python
class WorkspaceStoreProtocol:
    def load_all() -> dict[str, Workspace]
    def save(workspace: Workspace) -> None
    def delete(workspace_id: str) -> None
    def query_by_task(task_id: str) -> list[Workspace]
    def query_by_status(status: WorkspaceStatus) -> list[Workspace]
```

### D3: 并发锁 — per-task-group asyncio.Lock

**选择**: 按 task_id 分组加锁。同一 task 下的 create/merge/cleanup 操作串行化，不同 task 之间并行。

```python
self._locks: dict[str, asyncio.Lock]  # key = task_id

async def _get_lock(self, task_id: str) -> asyncio.Lock:
    if task_id not in self._locks:
        self._locks[task_id] = asyncio.Lock()
    return self._locks[task_id]
```

**理由**: 同一 task 下可能有多个 agent workspace 互相影响（共享 task branch），需要串行。不同 task 完全独立，不需要互斥。

**替代方案**: per-workspace Lock — 粒度太细，无法保护 task branch 的并发创建。

### D4: 启动恢复策略 — git worktree list + store reconcile

**选择**: 启动时执行 `git worktree list --porcelain`，解析出所有 worktree 路径和 branch，与 Store 中的 workspace 记录对比。

**三种情况**:
1. **Store 有，worktree 有** → 恢复到 ACTIVE 状态
2. **Store 有，worktree 没有** → 标记为 CLEANED
3. **Store 没有，worktree 有** → orphan，自动清理

**理由**: `git worktree list` 是 git 原生命令，输出格式稳定（porcelain 格式可解析）。这比扫描文件系统更可靠。

### D5: TTL 清理 — 可配置超时 + asyncio 后台 task

**选择**: 配置项 `WORKSPACE_TTL_SECONDS`（默认 3600），启动时创建后台 asyncio task 定期扫描。

```python
async def _ttl_cleanup_loop(self):
    while True:
        await asyncio.sleep(WORKSPACE_TTL_CHECK_INTERVAL)
        for ws in self.store.query_by_status(WorkspaceStatus.ACTIVE):
            if (datetime.now() - ws.created_at).total_seconds() > self._ttl:
                await self.cleanup(ws.id)
```

**理由**: asyncio task 零依赖，和 FastAPI lifespan 配合好。用 store.query_by_status 避免全量扫描。

### D6: container_id 预留字段

**选择**: `Workspace` 模型新增 `container_id: str | None = None`，不影响任何逻辑。

**理由**: 为后续 Container 隔离预留关联字段，不引入任何复杂度。

## Risks / Trade-offs

- **[task branch 冲突]** → agent branch merge 到 task branch 时可能冲突 → 和当前 merge 到 main 一样 abort + 返回错误，但影响范围缩小到 task 内
- **[JSON 文件并发写]** → 多协程同时写同一个 JSON 文件 → WorkspaceStore 内部加 asyncio.Lock 保护文件写操作
- **[recovery 误删 orphan]** → 有价值的 worktree 被自动清理 → 配置 orphan 保留策略（默认清理，可配置为标记不删）
- **[TTL 误杀活跃 workspace]** → Agent 执行时间超过 TTL → TTL 基于最后活跃时间而非创建时间，且提供 API 更新活跃时间
- **[per-task Lock 内存泄漏]** → task 完成后 Lock 对象残留 → cleanup 时同步删除对应 Lock
