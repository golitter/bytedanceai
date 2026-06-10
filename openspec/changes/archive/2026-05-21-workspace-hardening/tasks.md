## 1. WorkspaceStore 持久化

- [x] 1.1 Create `agentend/src/workspace/store.py` — 定义 `WorkspaceStoreProtocol` 抽象接口（load_all, save, delete, query_by_task, query_by_status）
- [x] 1.2 Implement `JsonFileWorkspaceStore` — 继承 `WorkspaceStoreProtocol`，JSON 文件存储，内部 asyncio.Lock 保护并发写，路径默认 `logs/workspaces.json`

## 2. 模型扩展

- [x] 2.1 Update `agentend/src/workspace/models.py` — `Workspace` 新增 `container_id: str | None = None` 字段
- [x] 2.2 Update `agentend/src/workspace/models.py` — 新增 task branch 常量生成函数 `task_branch_name(task_id) -> "task/{task_id}"`

## 3. GitOps 两级分支支持

- [x] 3.1 Update `agentend/src/workspace/git_ops.py` — 新增 `task_branch_create(repo_path, task_id)` 方法，创建 `task/{task_id}` branch（from main），已存在则跳过
- [x] 3.2 Update `agentend/src/workspace/git_ops.py` — 修改 `worktree_add` 签名，新增 `base_branch: str | None = None` 参数，worktree 创建时从 base_branch 创建 agent branch
- [x] 3.3 Update `agentend/src/workspace/git_ops.py` — 新增 `worktree_list()` 方法，执行 `git worktree list --porcelain` 并解析输出为 `list[tuple[path, branch]]`

## 4. WorkspaceManager 改造

- [x] 4.1 Update `agentend/src/workspace/manager.py` — 构造函数接受 `store: WorkspaceStoreProtocol` 和 `ttl_seconds: int` 参数，初始化 `_locks: dict[str, asyncio.Lock]`
- [x] 4.2 Rewrite `create()` — 先创建 task branch（调用 GitOps.task_branch_create），再从 task branch 创建 agent worktree，整个过程包在 `async with self._get_lock(task_id)` 中，成功后 `store.save()`
- [x] 4.3 Update `merge()` — 默认 target_branch 改为 `task/{task_id}`（agent → task），支持显式传入 `main`（task → main），包在 lock 中，成功后 `store.save()`
- [x] 4.4 Update `cleanup()` — 包在 lock 中，成功后 `store.save()`，清理完毕检查该 task 是否还有 ACTIVE workspace，无则删除 lock
- [x] 4.5 Update `cleanup_by_task()` — 遍历该 task 所有 ACTIVE workspace 串行 cleanup

## 5. 启动恢复

- [x] 5.1 Create `agentend/src/workspace/recovery.py` — 实现 `parse_worktree_list(output: str) -> list[tuple[str, str]]` 解析 `git worktree list --porcelain` 输出
- [x] 5.2 Implement `recover_workspaces(git_ops, store, repo_path)` — 调用 `git_ops.worktree_list()` 获取物理 worktree，与 store 中记录 reconcile：store+worktree 都有 → 恢复；store 有 worktree 无 → CLEANED；worktree 有 store 无 → 清理 orphan

## 6. TTL 清理

- [x] 6.1 Update `agentend/src/workspace/manager.py` — 新增 `_ttl_task: asyncio.Task | None` 和 `start_ttl_cleanup()`/`stop_ttl_cleanup()` 方法，后台循环 sleep + scan `store.query_by_status(ACTIVE)` + cleanup 过期
- [x] 6.2 Update `agentend/src/app/config.py` — 新增 `WORKSPACE_TTL_SECONDS: int = 3600` 和 `WORKSPACE_TTL_CHECK_INTERVAL: int = 300` 配置项

## 7. App Wiring

- [x] 7.1 Update `agentend/src/app/dependencies.py` — `create_workspace_manager` 改为实例化 `JsonFileWorkspaceStore` 并传入 `WorkspaceManager`，传入 TTL 配置
- [x] 7.2 Update `agentend/src/app/main.py` — lifespan startup 调用 `recover_workspaces()`，启动 TTL task；shutdown 停止 TTL task + 清理所有 ACTIVE workspace

## 8. Tests

- [x] 8.1 Update `agentend/tests/test_workspace.py` — 新增 WorkspaceStore 测试（save/load/delete/query），新增 create 时 task branch 创建逻辑测试，新增 merge 默认到 task branch 测试
- [x] 8.2 Create `agentend/tests/test_recovery.py` — 测试 parse_worktree_list 解析，测试 reconcile 三种场景（恢复/cleaned/orphan）
- [x] 8.3 Update `agentend/tests/test_workspace.py` — 测试 Lock 并发安全（同 task 串行、不同 task 并行），测试 TTL cleanup 过期/未过期
