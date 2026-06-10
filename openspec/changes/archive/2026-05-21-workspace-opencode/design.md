## Context

AgentEnd Runtime 已完成 MVP：Claude Code Adapter + Session Manager + Rule Engine + SSE 流式 API。当前所有 Agent 共享同一工作目录，多 Agent 并行操作同一仓库会导致文件冲突。

现有代码关键接口：
- `BaseAgentAdapter.stream_chat(session_id, message, **kwargs)` — `**kwargs` 已预留扩展空间
- `Session(workspace_path: str | None)` — 字段已存在但未实际使用
- `AgentRequest(workspace_path: str | None)` — 同上

OpenCode 是一个 Go 编写的开源终端 AI Coding Agent，支持多模型提供商，非交互模式通过 `opencode -p "prompt" -f json` 调用。

关键约束：
- Python 类型提示，模块导出风格 `from src.module import xxx`
- uv 管理依赖
- Git worktree 通过 `asyncio.create_subprocess_exec` 调用 git 命令实现
- OpenCode CLI **不支持** `stream-json`，只支持 `text` 和 `json` 输出格式
- OpenCode CLI **不支持** `--session` / `--system-prompt` / `--allowedTools` 参数

## Goals / Non-Goals

**Goals:**

- WorkspaceManager 基于 git worktree 为每个任务创建独立工作空间（独立 branch + 独立目录）
- OpenCodeAdapter 接入 OpenCode CLI 非交互模式
- Adapter 执行时自动绑定到对应 workspace 的 cwd
- Workspace 全生命周期管理：create → bind → commit → merge → cleanup
- 不破坏现有 ClaudeCodeAdapter 功能

**Non-Goals:**

- 不做 Docker 容器隔离（后续阶段）
- 不做冲突自动解决（merge conflict 由用户手动处理）
- 不做 OpenCode TUI 交互模式，只做非交互 `-p` 模式
- 不做 OpenCode session 持久化复用（OpenCode 内部通过 SQLite 管理，CLI 不暴露）
- 不做文件权限隔离（依赖 OS 级权限）

## Decisions

### D1: Workspace 隔离方案 — Git Worktree

**选择**: Git worktree（`git worktree add <path> -b <branch>`）

**理由**: 一个 repo 共享 `.git` 对象但拥有多个独立工作目录和 branch。比 `git clone` 多份省磁盘和带宽，比 Docker 轻量。

**目录结构**:
```
/repos/project-a/              # 主仓库（main branch）
/worktrees/
    task-123/frontend/          # agent/frontend/task-123
    task-123/backend/           # agent/backend/task-123
    task-456/reviewer/          # agent/reviewer/task-456
```

**Branch 命名**: `agent/{agent_name}/{task_id}`，如 `agent/frontend/task-123`

### D2: WorkspaceManager 职责边界

**选择**: WorkspaceManager 负责 worktree 创建/删除、branch 创建/合并；不负责 Agent 执行。

```
WorkspaceManager
├── create(repo_path, task_id, agent_name) → Workspace
├── get(workspace_id) → Workspace | None
├── list() → list[Workspace]
├── commit(workspace_id, message) → bool
├── merge(workspace_id, target_branch) → bool
├── cleanup(workspace_id) → bool
└── cleanup_by_task(task_id) → int    # 清理该任务所有 workspace
```

GitOps 封装底层 git 命令：
```
GitOps (内部工具类)
├── worktree_add(repo, path, branch) → bool
├── worktree_remove(path) → bool
├── branch_create(repo, name, base) → bool
├── add_and_commit(path, message) → bool
├── merge_branch(repo, branch, target) → bool
└── get_current_branch(path) → str
```

### D3: Adapter cwd 绑定方式

**选择**: 在 `asyncio.create_subprocess_exec` 时通过 `cwd=workspace_path` 参数指定工作目录。

**ClaudeCodeAdapter 改动**: `_build_command` 不变，`stream_chat` 从 kwargs 取 `cwd` 传入 subprocess。
**OpenCodeAdapter**: 同理，使用 `cwd` + OpenCode 的 `-c` 参数双保险。

### D4: OpenCode 流式输出策略

**选择**: 非流式 `opencode -p "..." -f json -q`，拿到完整 JSON 结果后拆分为多个 StreamEvent 一次性 yield。

**理由**: OpenCode CLI 只支持 `text` 和 `json` 输出格式，无 `stream-json`。完整 JSON 结果中包含消息和 tool_use 信息，可以解析后按类型拆分为 text / tool_call / tool_result / done 事件序列。

**替代方案**: 使用 `text` 格式 + 逐行 stdout 读取模拟流式 — 但无法区分 text 和 tool_call 事件类型，JSON 格式更结构化。

**后备**: 如果 OpenCode 非交互模式 JSON 输出格式不包含 tool 信息，降级为纯 text 输出，所有内容包装为 `type="text"` 的 StreamEvent。

### D5: OpenCode Rule 约束注入方式

**选择**: 将约束文本拼接到 prompt 前缀中。

```python
prompt = f"[系统约束: {rules_text}]\n\n{original_message}"
```

**理由**: OpenCode CLI 不支持 `--system-prompt` / `--append-system-prompt` 参数，只能通过 prompt 文本注入约束。不如 Claude Code 的 CLI 参数注入干净，但功能等价。

### D6: OpenCode Session 管理

**选择**: 不复用 OpenCode 内部 session，每次 `stream_chat` 调用都是独立的 `opencode -p` 进程。上下文连续性通过上层 Session Manager 的 history 机制维护（将历史消息拼接到 prompt 中）。

**理由**: OpenCode CLI 非交互模式不暴露 session 管理。内部 SQLite session 无法跨进程复用。

### D7: Workspace 与 Session 的绑定关系

**选择**: Session 创建时关联 workspace_path。一个 Session 对应一个 Workspace。

```
Session
  ├── workspace_path: str    # 指向 worktree 目录路径
  └── agent_type: str        # 决定用哪个 Adapter

Workspace
  ├── session_id: str        # 反向引用
  ├── repo_path: str         # 主仓库路径
  ├── worktree_path: str     # worktree 目录路径
  └── branch_name: str       # agent/{name}/{task_id}
```

## Risks / Trade-offs

- **[Git worktree 冲突]** → 多 Agent 修改同一文件时 merge conflict → MVP 不自动解决，merge 失败返回错误，由用户或 Orchestrator 处理
- **[Worktree 残留]** → Agent 异常退出未 cleanup → WorkspaceManager 提供 `cleanup_by_task` 批量清理，app shutdown 时也尝试清理
- **[OpenCode CLI 输出格式不确定]** → JSON 格式可能不包含 tool_use 详情 → 降级为 text 输出，后续根据实际格式调整解析逻辑
- **[OpenCode 无流式]** → 用户等待时间长于 Claude Code → 在 StreamEvent 中标注来源 agent_type，前端可据此展示不同等待状态
- **[Prompt 注入约束不如 CLI 参数可靠]** → Agent 可能忽略 prompt 中的约束 → 后期 OpenCode 支持 `--system-prompt` 参数时迁移
- **[OpenCode 项目已归档，迁移至 Crush]** → CLI 接口可能变化 → Adapter 隔离了具体 CLI 差异，只需修改 `_build_command` 和 `_parse_output`
