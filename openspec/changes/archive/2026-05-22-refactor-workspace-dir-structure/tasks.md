## 1. Workspace Model 调整

- [x] 1.1 修改 `_generate_worktree_path` 签名：`agent_name` 参数改为 `session_id`，生成路径 `{repo_parent}/worktrees/{task_id}/{session_id}`
- [x] 1.2 修改 `_generate_branch_name` 签名：`agent_name` 参数改为 `session_id`，生成 `agent/{session_id}/{task_id}`
- [x] 1.3 Workspace dataclass 新增 `agent_type` 字段（AgentType 枚举类型），`session_id` 从 `str | None` 改为 `str`（必传）
- [x] 1.4 修改 `__post_init__`：使用 `session_id` 而非 `agent_name` 生成 `worktree_path` 和 `branch_name`

## 2. WorkspaceManager 适配

- [x] 2.1 修改 `WorkspaceManager.create()` 签名：新增 `session_id` 和 `agent_type` 参数，传给 Workspace 构造
- [x] 2.2 修改 `create()` 内部逻辑：git worktree 和 branch 操作使用新的路径和分支名格式
- [x] 2.3 修改 `create()` 中调用 provisioner 的参数：传入 `session_id` 和 `agent_type` 替代 `agent_name`

## 3. API 层适配

- [x] 3.1 修改 `api/v1/agent.py` 中 `_resolve_workspace()`：传 `session_id` 和 `agent_type` 给 `workspace_mgr.create()`
- [x] 3.2 修改 `api/v1/workspace.py` 中直接创建 workspace 的接口：确保 `session_id` 必传
- [x] 3.3 验证 request schema 中 `session_id` 为必填字段

## 4. Skills Provisioner 适配

- [x] 4.1 修改 `provision()` 方法：使用 `agent_type` 替代 `agent_name` 查找 skills 目标目录
- [x] 4.2 修改 `init_shared_dirs()` 方法：shared memory 子目录使用 `session_id` 替代 `agent_name`

## 5. taskctl 适配

- [x] 5.1 修改 taskctl 路径解析逻辑：从 `worktrees/{task_id}/{session_id}/.{agent_type}/skills/taskctl/exe` 解析 `task_id` 和 `session_id`
- [x] 5.2 修改 `sub-memory` 命令：读取 `memory/{session_id}/` 替代 `memory/{agent_name}/`
- [x] 5.3 重新编译 taskctl 二进制

## 6. 测试与验证

- [x] 6.1 更新 `workspace/models.py` 相关单元测试：覆盖新的路径和分支名生成逻辑
- [x] 6.2 更新 `workspace/manager.py` 相关测试：验证 create 使用 session_id 的完整流程
- [x] 6.3 更新 provisioner 测试：验证 agent_type 映射和 session_id 隔离
- [x] 6.4 验证 API 端到端流程：创建 workspace → skills 部署 → shared 目录初始化
