# Agent 路由与 RunTask 响应契约

## 变更原因

Agent 路由实现新增了群聊 mention 直达、群聊无 `@` 交给 Orchestrator、单聊保持当前 Agent、以及 `skip_user_message=true` 内部分派跳过路由的语义。Backend `POST /tasks/:taskId/run` 响应也从仅返回 `message_id/status` 扩展为返回实际执行 `session_id`、`agent_type`、`route_id`、`route_mode` 等字段，Frontend 需要据此订阅实际 SSE session。

同时，群聊历史新增 `mode=group` / `primary_session_id` 查询约定，用于显示 task 下全部 user 消息并避免 Orchestrator 镜像的子 Agent 消息重复。

## 变更文件

- `contracts/schemas/agent-routing.yaml` — 新增 Agent 路由、RunTask 请求/响应、群聊消息查询参数契约。
- `scripts/generate_contracts.py` — 新增 `agent-routing` schema 到三端生成映射。
- `contracts/AGENTS.md` — 补充 schema 列表与生成文件映射。
- `agentend/src/generated/agent_routing.py` — 由 `make generate` 生成。
- `frontend/src/generated/agent-routing.ts` — 由 `make generate` 生成。
- `backend/internal/generated/agent_routing.go` — 由 `make generate` 生成。

## 对比结果

变更前：

- RunTask Backend API 为手写类型，契约层没有描述 `{ message_id, status }` 之外的路由结果。
- Task session 的 `route_id` / `mention_label` / `aliases` 仅存在于实现设计中，没有共享契约。
- 群聊消息查询仅约定 `session_id`，没有 `mode=group` 的可见性过滤契约。

变更后：

- 新增 `RouteMode`：`direct` / `orchestrator` / `unchanged`。
- 新增 `AgentRoute`：包含 `session_id`、`agent_type`、`agent_name`、`route_id`、`mention_label`、`aliases`。
- 新增 `RunTaskRequest`：包含 `message`、`session_id`、可选 `agent_type`、`cwd`、`skip_user_message`。
- 新增 `RunTaskResponse`：包含 `message_id`、`status`、实际执行 `session_id`、`agent_type`、可选 `agent_name`、`route_id`、`route_mode`。
- 新增 `GroupMessagesQuery`：描述 `mode`、`primary_session_id`、`session_id`、`limit`、`before` 查询参数。

## 跨端影响

- **Frontend**：`submitMessage` 可使用生成的 `RunTaskResponse` 类型，并基于响应中的实际 `session_id` 建立 SSE；mention 菜单可消费 `AgentRoute` 等价字段。
- **Backend**：`RunTask` 响应需要持续返回路由字段；`GetTask` session 数据和 Orchestrator config 应复用同一套 `route_id` 生成逻辑。
- **AgentEnd**：内部 `BackendClient.run_task` 仍只依赖 `message_id`，对新增响应字段向后兼容；可按需引用生成类型描述路由响应。

## 契约变更

| 类型 | 字段 / 枚举 | 说明 |
|------|-------------|------|
| `RouteMode` | `direct` | 显式 `@` 唯一非 Orchestrator Agent，直达目标 Agent |
| `RouteMode` | `orchestrator` | 群聊自动编排，交给 Orchestrator |
| `RouteMode` | `unchanged` | 单聊或内部 dispatch，保持请求目标 |
| `AgentRoute` | `route_id` | 后端生成的唯一路由 ID |
| `AgentRoute` | `mention_label` | 前端插入 `@` 时使用的标签 |
| `RunTaskResponse` | `session_id` | 实际执行并产生 SSE 的 session |
| `RunTaskResponse` | `route_mode` | 本次消息的实际路由模式 |
| `GroupMessagesQuery` | `mode=group` | 按群聊可见性返回消息 |
