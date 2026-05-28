# 契约变更：子 Agent 消息持久化 — ExecutionEngine 回调后端 RunTask

## 变更原因

Orchestrator 分派子 Agent 时，原流程在 agentend 内部直接调用 adapter，子 Agent 响应不经过后端 StreamWriter，导致消息不持久化到 MySQL。改为 ExecutionEngine 回调后端 RunTask API，复用 StreamWriter 持久化链路。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。所有改动复用现有契约字段。

## 涉及的契约字段复用

| 契约字段 | 复用方式 |
|----------|----------|
| `AgentRequest.workspace_path` | 后端 RunTask 将 orchestrator 创建的 worktree 路径通过此字段传递给 agentend |
| `AgentRequest.config` | orchestrator 的 agents 列表、task_id 等仍通过此字段传递，无变化 |
| `Message.role` | 子 Agent 消息的 user/agent role 保持不变 |
| `Message.status` | streaming/completed/failed 状态流转保持不变 |
| `Message.agent_type` / `Message.agent_name` | 子 Agent 的 agent_type 和 agent_name 由后端从 session 表读取写入 |

## 后端内部 API 变更（非契约）

`RunTaskReq` 新增两个内部字段，不影响跨端契约：

| 字段 | 类型 | 说明 |
|------|------|------|
| `cwd` | `string` | 子 Agent 的 worktree 路径，传递给 agentend 的 `AgentRequest.workspace_path` |
| `skip_user_message` | `bool` | orchestrator 内部分派时为 true，后端跳过创建 user message |

## 跨端影响

- **AgentEnd**: 新增 `BackendClient` 模块（HTTP 调用后端 RunTask + SSE 订阅）；`ExecutionEngine` 从直接调用 adapter 改为回调后端 API
- **Backend**: `RunTaskReq` 新增 `cwd` 和 `skip_user_message` 字段；RunTask 为子 Agent 创建 Message 并通过 StreamWriter 持久化
- **Frontend**: 无变更 — `getTaskMessages` 已返回所有 session 消息，SSE 订阅机制不变
- **Contracts**: 无变更 — 枚举值和字段定义保持不变
