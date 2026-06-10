## Why

Workspace 模块 MVP 已实现基础 git worktree 隔离，但存在四个生产级缺陷：(1) 所有 Agent branch 直接 merge 到 main，多 Agent 并发会互相污染；(2) workspace 状态纯内存存储，进程崩溃或重启后全部丢失，worktree 变成 orphan；(3) create/merge/cleanup 无并发保护，高并发下竞态条件会导致 worktree 状态不一致；(4) 无自动清理机制，worktree 会无限积累占用磁盘。

## What Changes

- **两级分支结构**：引入 `task/{task_id}` 集成分支作为 Agent branch 的父级。Agent worktree 的 branch 从 task branch 创建，merge 时先合入 task branch，最后由 task branch 合入 main。避免多 Agent 直接操作 main。
- **Workspace 持久化**：新增 `WorkspaceStore`，将 workspace 状态持久化到 JSON 文件（与现有 `SessionMappingStore` 模式一致），抽象存储接口以便后续迁移 SQLite。
- **并发锁**：`WorkspaceManager` 新增 per-workspace `asyncio.Lock`，保护 create/merge/cleanup 三个写操作的并发安全。
- **启动恢复**：服务启动时通过 `git worktree list --porcelain` 扫描物理 worktree，与持久化存储 reconcile，自动恢复或清理 orphan worktree。
- **TTL 清理**：后台 asyncio task 定期检查 ACTIVE workspace，超时自动 cleanup，防止 worktree 无限增长。
- **模型扩展**：`Workspace` 新增 `container_id: str | None` 字段预留容器隔离接口。

## Capabilities

### New Capabilities
- `workspace-store`: Workspace 持久化存储，基于 JSON 文件的 CRUD，抽象存储接口（WorkspaceStoreProtocol），支持 load/save/query/delete
- `workspace-recovery`: 启动时通过 git worktree list 扫描物理 worktree，与持久化状态 reconcile，恢复或清理 orphan
- `workspace-ttl`: 后台定时清理过期 ACTIVE workspace 的 TTL 机制，可配置超时时长

### Modified Capabilities
- `workspace-isolation`: 分支结构从扁平（agent/* → main）改为两级（main → task/* → agent/*），Workspace 模型新增 container_id 字段，WorkspaceManager 增加并发锁保护

## Impact

- **新增文件**: `src/workspace/store.py`（持久化）、`src/workspace/recovery.py`（启动恢复）
- **修改文件**: `src/workspace/models.py`（container_id 字段）、`src/workspace/git_ops.py`（task branch 创建/合并逻辑）、`src/workspace/manager.py`（锁、store 集成、TTL）、`src/app/config.py`（TTL 配置项）、`src/app/main.py`（recovery 调用 + TTL task 启动）、`src/app/dependencies.py`（store 实例化）
- **修改文件**: `src/api/v1/workspace.py`（task branch 创建端点）、`src/api/v1/agent.py`（workspace 创建时关联 task branch）
- **新增依赖**: 无
- **数据迁移**: 首次启动 store 文件为空，无需迁移
