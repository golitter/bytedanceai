## Why

当前 AgentEnd 已支持单 Agent 对话（Claude Code / OpenCode），通过适配器模式调用 CLI 进程。但缺少多 Agent 协作能力 — 用户无法在一次请求中协调多个 Agent 分工完成复杂任务。Orchestrator 是 AgentHub 群聊模式的核心组件，Phase 1 先实现"规划 + 写任务文件"这一步，打通 Orchestrator → shared/.agent/ → taskctl 的链路。

## What Changes

- 新增 `OrchestratorAdapter`，作为 `BaseAgentAdapter` 的实现，通过 LangGraph StateGraph 驱动规划流程
- 新增 LangGraph graph：`plan`（LLM structured output）→ `write_shared`（文件 IO），将用户需求拆解为多个 task markdown 文件写入 `shared/.agent/` 目录
- 新增 `LlmConfig` 配置段（`config.yaml` + `config.py`），用于 Anthropic API 调用
- 扩展 `AgentType` 枚举增加 `ORCHESTRATOR`，扩展 `EventType` 枚举增加 `PLANNING`
- 修改 API 层（`agent.py`）透传 `config.agents`、`config.task_id`、`config.shared_dir` 到 adapter kwargs

## Capabilities

### New Capabilities
- `orchestrator-planning`: Orchestrator 接收用户消息 + agent 列表，通过 LLM 拆解任务，将 overview.md、tasks/*.md、config.yaml 写入 shared/.agent/ 目录

### Modified Capabilities
- `adapter-registry`: 注册 OrchestratorAdapter 为新的 agent type
- `event-types`: 新增 PLANNING 事件类型用于前端展示规划进度

## Impact

- **新增依赖**: `langgraph>=0.4`, `langchain-anthropic>=0.3`, `langchain-core>=0.3`
- **配置变更**: `config.yaml` 新增 `llm` 段（provider, model, api_key）
- **API 变更**: `POST /v1/agent/execute` 和 `POST /v1/agent/stream` 的 `config` 字段新增 `agents`、`task_id`、`shared_dir` 子字段（仅 `agent_type=orchestrator` 时使用）
- **文件系统**: 运行时写入 `shared/.agent/` 目录（config.yaml, overview.md, tasks/*.md）
- **现有功能无破坏性变更**: 现有 Claude Code / OpenCode adapter 不受影响
