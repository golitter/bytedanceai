## Why

群聊场景中，每个子 Agent（Claude Code / OpenCode / Codex）执行时只能看到自己的 CLI session 历史和 Orchestrator 分发的任务内容，完全不知道其他 Agent 说了什么。这导致 Agent 之间无法协作、可能重复工作、缺乏全局视角。需要让 Agent 在执行时能看到「自上次发言以来，其他 Agent 的消息」。

## What Changes

- Go Backend RunTask handler 新增窗口消息查询：查询当前 task 下其他 session 的消息（`status IN ("completed","streaming")`，从当前 agent 上次发言起算），每条截断 2000 字符，注入 `group_chat_messages` 到请求体
- AgentEnd API 层解析 `group_chat_messages`，放入 rule context
- 新增 `GroupChatRule`：有窗口消息时格式化为 `system_prompt_append`
- 新增群聊 Prompt 模板 `group_chat.py`
- **删除 ExecutionEngine 的 short-circuit 路径**，统一走 HTTP（确保 Rules 对所有 Agent 生效）
- Orchestrator 自身也通过 Backend API 查窗口，注入 REASON prompt
- OrchestratorAdapter 简化：去掉 ContextBuilder 依赖

## Capabilities

### New Capabilities
- `group-chat-window-query`: Go Backend 窗口消息查询——基于 task + session 的跨 Agent 消息窗口，含 status 过滤和截断
- `group-chat-rule`: AgentEnd GroupChatRule——将窗口消息格式化为 system_prompt_append，集成到现有 Rules 引擎

### Modified Capabilities
- `rule-engine`: 新增 GroupChatRule 内置规则（优先级 6，phase=pre）
- `sub-agent-backend-dispatch`: ExecutionEngine 删除 short-circuit，统一走 HTTP 路径
- `orchestrator-planning`: REASON prompt 新增 `{orchestrator_context}` 占位，Orchestrator 自身可查询跨 Agent 窗口

## Impact

- **Backend**: `task.go`（RunTask 注入窗口查询）、`message.go`（新增 WindowMessages API）、`main.go`（新路由）
- **AgentEnd**: `agent.py`（解析 group_chat_messages）、`builtin.py`（GroupChatRule）、新文件 `prompts/group_chat.py`、`engine.py`（删 short-circuit）、`orchestrator.py`（简化 + 窗口查询）、`backend_client.py`（新增方法）、`prompts.py`（REASON prompt 扩展）
- **不改**: claude.py、opencode.py、codex.py、Frontend
- **删除**: `context_builder.py`（如存在）
