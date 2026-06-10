## Why

当前工作区目录结构使用 `agent_name`（如 `claude-code`）作为路径层级：`worktrees/{task_id}/{agent_name}`。这导致同一 agent 类型的多个 session 无法并发运行——它们的 worktree 路径会冲突。需要将目录维度从 agent 类型切换为 session，以支持同 task 下多 session 并行，同时 `agent_type` 仅作为元数据标识，不再参与路径构建。

## What Changes

- **BREAKING**: 工作区目录路径从 `worktrees/{task_id}/{agent_name}` 变更为 `worktrees/{task_id}/{session_id}`
- **BREAKING**: Git 分支名从 `agent/{agent_name}/{task_id}` 变更为 `agent/{session_id}/{task_id}`
- Workspace model 新增 `agent_type` 字段，用于替代原 `agent_name` 在 skills 和 config 映射中的角色
- Workspace 创建时必须传入 `session_id`，且 `session_id` 参与路径生成
- Skills provisioner 使用 `agent_type` 而非 `agent_name` 来查找配置目录
- Shared 目录结构保持 `shared/.agent/memory/` 不变，子目录改为按 session 隔离

## Capabilities

### New Capabilities

- `session-based-workspace-path`: 用 session_id 替代 agent_name 构建工作区目录和分支命名

### Modified Capabilities

- `workspace-isolation`: worktree 路径生成逻辑变更，隔离维度从 agent_type 切换为 session_id
- `workspace-management`: Workspace model 新增 agent_type 字段，workspace 创建和查找逻辑适配新路径
- `skill-provisioning`: provisioner 使用 agent_type 而非 agent_name 进行配置目录映射
- `taskctl-cli`: taskctl 解析自身路径获取 session_id 而非 agent_name

## Impact

- **核心文件**: `workspace/models.py`, `workspace/manager.py`, `workspace/store.py`, `api/v1/agent.py`, `skills/provisioner.py`
- **API 层**: request schema 需确保 session_id 必传，agent_type 独立于路径
- **存储**: 已持久化的 workspace 记录（workspaces.json）路径格式不兼容，需迁移或重建
- **分支命名**: 现有 git 分支 `agent/{agent_name}/{task_id}` 不兼容，需清理后重建
