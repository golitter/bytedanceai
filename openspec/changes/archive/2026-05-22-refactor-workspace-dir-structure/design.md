## Context

当前 workspace 目录结构为 `worktrees/{task_id}/{agent_name}`，其中 `agent_name` 取自 `AgentType.value`（如 `"claude-code"`）。这意味着同一 task 下，同类型 agent 只能有一个 workspace，因为路径会冲突。

实际需求是：同一 task 下同一个 agent 类型可以有多个并发 session（如多个 Claude Code 会话同时工作），每个 session 需要独立的 worktree。

现有代码中 Workspace model 已有 `session_id` 字段（可选），但未参与路径生成。`agent_type` 的信息仅通过 `agent_name` 传递，没有独立字段。

## Goals / Non-Goals

**Goals:**
- 用 `session_id` 替代 `agent_name` 作为目录层级，支持同 task 多 session 并发
- Workspace model 显式区分 `agent_type`（元数据）和 `session_id`（路径维度）
- Skills provisioner 和 taskctl 适配新路径结构

**Non-Goals:**
- 不改变 shared 目录的基础结构（仍为 `shared/.agent/`）
- 不改变 workspace store 的持久化格式（JSON）
- 不做向后兼容迁移——清理旧 workspace 后使用新结构

## Decisions

### Decision 1: 路径结构从 `{task_id}/{agent_name}` 变为 `{task_id}/{session_id}`

**选择**: `worktrees/{task_id}/{session_id}`

**替代方案**:
- `worktrees/{task_id}/{agent_type}/{session_id}`: 三级结构，保留 agent_type 分组。被否决因为增加了路径深度但无实际收益——同一 task 下按 agent_type 分组的场景可通过元数据查询实现。

**理由**: session_id 是唯一的，无需额外层级即可避免冲突。

### Decision 2: 分支命名从 `agent/{agent_name}/{task_id}` 变为 `agent/{session_id}/{task_id}`

**理由**: 分支名需要与 worktree 一一对应，session_id 可直接保证唯一性。

### Decision 3: Workspace model 字段调整

保留 `agent_name` 字段但语义变更——不再用于路径生成，仅存储 agent 类型标识。同时 `agent_type` 作为独立字段存储枚举值。`session_id` 从可选变为必传。

实际方案：将 `agent_name` 保留用于 skills 映射等用途（值仍为 `"claude-code"` 等），新增 `agent_type` 字段存储 `AgentType` 枚举值用于类型判断。路径生成统一使用 `session_id`。

### Decision 4: taskctl 路径解析适配

taskctl exe 路径变为 `worktrees/{task_id}/{session_id}/.{agent_type}/skills/taskctl/exe`。解析逻辑需从路径中提取 `session_id` 而非 `agent_name`，`agent_type` 通过 skills 目录名推断（`.claude` → claude-code，`.opencode` → opencode）。

## Risks / Trade-offs

- **[路径格式不兼容]** → 旧 workspace 的 worktree_path 和 branch_name 格式与新代码不兼容。Mitigation: 升级前清理所有现有 workspace，不提供迁移工具。
- **[session_id 必传]** → API 请求必须包含 session_id，缺失将导致创建失败。Mitigation: API 层校验 session_id 非空。
- **[taskctl 路径解析变更]** → 已部署的 taskctl 二进制无法解析新路径。Mitigation: 重新编译部署。
