## 1. AgentEnd — BackendClient

- [x] 1.1 新建 `agentend/src/clients/backend_client.py`，实现 `BackendClient` 类，封装对后端的 HTTP 调用（`POST /api/tasks/:taskId/run`、`GET /api/tasks/:taskId/stream`）
- [x] 1.2 在 `agentend/src/app/config.py` 的 Settings 中新增 `backend_url` 配置项（默认 `http://localhost:8080`）
- [x] 1.3 在 FastAPI app 生命周期中初始化 `BackendClient` 单例，注入到 OrchestratorAdapter

## 2. AgentEnd — ExecutionEngine 重写

- [x] 2.1 重写 `agentend/src/orchestrator/execution/engine.py` 的 `_execute_task()` 方法：不再直接调用 adapter，改为调用 `BackendClient.run_task()` 发起子 Agent run
- [x] 2.2 实现 SSE 结果收集：`_execute_task()` 调用 RunTask 获取 `message_id` 后，通过 `BackendClient.stream_result()` 订阅 SSE，收集 TEXT 事件构建 TaskResult
- [x] 2.3 保留 worktree 创建逻辑（`_ensure_worktree`），将 worktree 路径通过 RunTask 请求传递给后端→agentend
- [x] 2.4 支持并发分派：多个子 Agent 的 `_execute_task()` 使用 `asyncio.gather` 并发执行

## 3. Backend — RunTask 适配

- [x] 3.1 检查 `backend/internal/handler/task.go` 的 `RunTask` 是否能正确处理子 Agent 的 session_id（session 在 CreateTask 时已创建，RunTask 应能找到现有 session）
- [x] 3.2 确认 RunTask 为子 Agent 创建的 Message 记录包含正确的 `agent_type` 和 `agent_name`（从 session 表读取）
- [x] 3.3 确保 RunTask 传入的 `cwd`（worktree 路径）能传递到 agentend 的 adapter 请求中

## 4. Backend — SSE Stream 端点确认

- [x] 4.1 确认 `GET /api/tasks/:taskId/stream?message_id=xxx` 能为子 Agent 的 message_id 正确推送 Redis Stream 事件
- [x] 4.2 确认 agentend 的 ExecutionEngine 能通过 SSE 端点订阅子 Agent 的流式输出

## 5. 集成验证

- [ ] 5.1 启动三端服务，创建包含 orchestrator + claude-code + opencode 的 task
- [ ] 5.2 发送消息触发 orchestrator 分派，验证子 Agent 消息出现在数据库 messages 表
- [ ] 5.3 刷新前端页面，验证子 Agent 的历史消息正确加载和显示
- [ ] 5.4 验证实时流式显示仍正常（orchestrator 流推送 RUNTIME_EXECUTING/COMPLETED 事件）
