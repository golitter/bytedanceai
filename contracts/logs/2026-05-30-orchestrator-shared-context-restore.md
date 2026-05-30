# 2026-05-30 orchestrator shared 规划上下文恢复

## 变更原因

Orchestrator Agent 生命周期重构后，规划结果只保留在内存 `PlanOutput` 中，未继续写入 `shared/.agent/config.yaml` 和 `plans/*.md`。子 Agent 的 `taskctl summary` 依赖 shared 目录中的 `config.yaml`，并使用真实 session_id 筛选当前 Agent 的任务，导致子 Agent 无法读取分配内容。

同时，Orchestrator 入口仍会在带 `repo_path` 时自动创建代码 worktree，并将该 worktree 加入规划工具可读范围，与 Phase 5 中“Orchestrator 只做协调、不拥有代码工作区”的设计不一致。

## 变更文件

- `agentend/src/orchestrator/planning/graph.py` — dispatch 阶段将 `PlanOutput` 写入 `shared/.agent/config.yaml`、`plans/overview.md`、`plans/task-*.md`
- `agentend/src/orchestrator/planning/tools.py` — `read_file` / `list_dir` / `write_file` 的相对路径从 `shared_dir` 解析
- `agentend/src/orchestrator/planning/prompts.py` — 更新工具说明，标注读取范围仅限 shared 目录
- `agentend/src/adapters/orchestrator.py` — Orchestrator 规划工具可读范围收紧为 `shared_dir`
- `agentend/src/api/v1/agent.py` — Orchestrator 不再自动创建代码 worktree，并校验 `shared_dir` 必须指向当前 task 的标准 shared 目录
- `backend/internal/handler/task.go` — Orchestrator config 补齐 `shared_dir`、`repo_path`、子 Agent 真实 `session_id`
- `agentend/tests/test_orchestrator_shared.py` — 新增 shared 规划链路回归测试

## 对比结果

无 schema 变更。`contracts/schemas/agent-request.yaml` 中 `AgentRequest.config` 已定义为可扩展 object（`additionalProperties: true`），本次只规范内部使用的配置键：

- `agents[]`
- `agents[].session_id`
- `repo_path`
- `shared_dir`
- `task_id`

这些字段均通过既有 `config` 扩展对象承载，不改变 AgentRequest、StreamEvent 或其他跨端结构。

## 跨端影响

- **Frontend**: 无影响。SSE 事件类型和消息结构不变。
- **Backend**: Orchestrator 分支构建 `AgentRequest.config` 时新增 `shared_dir`、`repo_path`，并保证 `agents[].session_id` 传递真实子 Agent session id。
- **AgentEnd**: Orchestrator 规划阶段恢复 shared 落盘；子 Agent 可通过 `taskctl summary` 按真实 session id 读取自己的任务。Orchestrator 不再拥有代码 worktree，规划工具只读 shared 上下文。

## 契约变更

无。本次改动不修改 `contracts/schemas/`；仅记录 `AgentRequest.config` 既有扩展字段的跨端使用约定。
