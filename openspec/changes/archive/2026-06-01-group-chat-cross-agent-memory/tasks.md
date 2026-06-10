## 1. Go Backend — 窗口查询核心

- [x] 1.1 在 `backend/internal/handler/task.go` 的 RunTask handler 中，新增 `fetchGroupChatWindow(taskID, sessionID)` 方法：查询其他 session 的消息（`status IN ("completed","streaming")`，从当前 session 最后一条 agent 消息时间起算），每条截断 2000 rune 字符
- [x] 1.2 在 RunTask handler 中，调用 `fetchGroupChatWindow` 并将结果作为 `group_chat_messages` 注入请求体发送给 AgentEnd
- [x] 1.3 窗口查询失败时降级为空列表 `[]`，记录 warning 日志，不阻塞 RunTask

## 2. Go Backend — 窗口查询 API 路由

- [x] 2.1 在 `backend/internal/handler/message.go` 新增 `WindowMessages` handler，复用 `fetchGroupChatWindow` 的查询逻辑
- [x] 2.2 在 `backend/cmd/server/main.go` 注册路由 `GET /api/tasks/:taskId/messages/window?session_id=xxx`

## 3. AgentEnd — 群聊 Prompt 模板

- [x] 3.1 新建 `agentend/src/orchestrator/prompts/group_chat.py`，实现 `GROUP_CHAT_CONTEXT` 模板和 `build_group_chat_context()` 函数，支持 user/agent 消息格式化，空输入返回空字符串

## 4. AgentEnd — GroupChatRule

- [x] 4.1 在 `agentend/src/rules/builtin.py` 新增 `GroupChatRule` 类（name="group_chat", phase="pre", priority=6），`check` 检查 `group_chat_messages` 非空，`enforce` 调用 `build_group_chat_context` 返回 `system_prompt_append`
- [x] 4.2 确认 GroupChatRule 在系统启动时被 RuleRegistry 自动加载

## 5. AgentEnd — API 层解析

- [x] 5.1 修改 `agentend/src/api/v1/agent.py` 的 stream endpoint，从请求体提取 `group_chat_messages` 放入 rule context（缺失时默认为 `[]`）

## 6. AgentEnd — ExecutionEngine 去 short-circuit

- [x] 6.1 删除 `agentend/src/orchestrator/execution/engine.py` 中的 `_get_adapter()`、`_iter_adapter_with_timeout()` 方法和 short-circuit 分支逻辑
- [x] 6.2 删除 ExecutionEngine 构造函数的 `adapter_registry` 参数，统一走 `BackendClient.run_task()` HTTP 路径
- [x] 6.3 删除 `agentend/src/orchestrator/execution/context_builder.py`（如存在）
- [x] 6.4 更新 OrchestratorAdapter `_handle_execute()` 中 ExecutionEngine 的构造调用，移除 `adapter_registry` 参数

## 7. AgentEnd — Orchestrator 自身上下文

- [x] 7.1 在 `agentend/src/clients/backend_client.py` 新增 `get_agent_window_messages(task_id, session_id)` 方法，调用 `GET /api/tasks/:taskId/messages/window`
- [x] 7.2 修改 `agentend/src/adapters/orchestrator.py` 的 `stream_chat()`，在构建 initial_state 前查询自身窗口消息，格式化后注入 `orchestrator_context`
- [x] 7.3 修改 `agentend/src/orchestrator/planning/prompts.py` 的 REASON_PROMPT，新增 `{orchestrator_context}` 占位

## 8. 验证

- [x] 8.1 验证单聊无影响：单聊走 Backend → AgentEnd → adapter，窗口查询为空，GroupChatRule 不触发
- [ ] 8.2 验证窗口过滤逻辑：无历史消息 → 全量查询；有历史消息 → 窗口过滤；status 过滤生效
- [x] 8.3 验证消息截断：超 2000 字符的消息被正确截断并添加 `[截断]` 后缀
- [ ] 8.4 验证端到端群聊：2 Agent 群聊 → 第 2 轮 agent 的 system_prompt_append 包含第 1 轮另一个 Agent 的输出
- [x] 8.5 验证 Rules 共存：子 Agent 的 system_prompt_append 同时包含 Safety/Soul/Skill/GroupChat 内容
- [ ] 8.6 验证 Wave 隔离：并行 wave 内 Agent 互相看不到输出，串行 wave 间能看到
- [ ] 8.7 验证 Orchestrator 自身：@Agent 后再 @Orchestrator → REASON prompt 包含之前其他 Agent 的消息
