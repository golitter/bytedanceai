# 2026-06-01-orchestrator-task-base-worktree

## 变更原因

Orchestrator 之前只能看到 `shared/.agent/` 元数据目录，无法访问代码库，导致任务拆解时无法根据实际代码结构做精准规划。现在为每个任务创建 `task-base` worktree（基于 `task/{task_id}` 分支），让 orchestrator 能读取完整代码。

## 变更文件

无 schema 文件修改。`task_base_path` 是 agentend 内部从 `repo_path` + `task_id` 推导的路径，不涉及跨端请求协议变更。

## 变更内容

### 新增目录结构

```
worktrees/{task_id}/
├── task-base/              ← 新增：task/{task_id} 分支的 worktree（只读）
├── shared/.agent/          ← 不变：元数据 + .orchestrator/skills/
└── {session_X}/            ← 不变：Sub-Agent worktree
```

### 新增方法

| 文件 | 方法 | 说明 |
|------|------|------|
| `workspace/git_ops.py` | `task_base_worktree_create(repo_path, task_id)` | 创建 task-base worktree |
| `workspace/git_ops.py` | `task_base_worktree_remove(repo_path, task_id)` | 删除 task-base worktree |
| `workspace/manager.py` | `create_task_base(repo_path, task_id)` | 带锁的 task-base 创建入口 |

### 修改方法

| 文件 | 方法 | 说明 |
|------|------|------|
| `workspace/manager.py` | `cleanup_by_task()` | 删 task 分支前先移除 task-base worktree |
| `api/v1/agent.py` | `_resolve_workspace()` | orchestrator 请求时创建 task-base |
| `api/v1/agent.py` | `_orchestrator_kwargs()` | 计算 `task_base_path` 并传入 kwargs |
| `adapters/orchestrator.py` | `stream_chat()` | `allowed_read_dirs` 加入 task-base |
| `planning/graph.py` | `GraphState` | 新增 `task_base_path` 字段 |
| `planning/graph.py` | `skill_prepare_node()` | 传递 `task_base_path` 给 prompt |
| `planning/prompts.py` | `build_reason_prompt()` | 提示词描述两个目录 + 子 Agent 用相对路径 |

## 跨端影响

- **Frontend**: 无影响
- **Backend**: 无影响（worktree 管理纯 agentend 侧）
- **AgentEnd**: orchestrator 现在能读取代码库，提示词新增工作区目录描述
- **Contracts**: 无 schema 变更
