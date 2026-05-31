# Orchestrator 转发导致子 Agent 消息重复存储

**状态**: 已知问题，待修复
**发现日期**: 2026-05-31
**影响范围**: 群聊场景，跨 Agent 记忆窗口查询会注入重复内容

## 问题描述

群聊中，子 Agent（如 claude-code）的回复消息在 MySQL `messages` 表中被存储两次：

1. **第一次**：子 Agent 自己的 session（由 Backend RunTask 直接创建）
2. **第二次**：Orchestrator 的 session（Orchestrator 将子 Agent 回复透传到自己的 SSE 流后，由流式刷写写入）

### 示例数据

```
id=48  session=75bcaaea (aa)  agent=aa  "我是 aa，一个 Claude Code Agent..."  ← 第1次
id=49  session=aa390b42 (orch) agent=aa  "我是 aa，一个 Claude Code Agent..."  ← 第2次（重复）
```

## 根因

消息流经过两条路径，两条路径都会持久化：

```
用户消息 → Backend RunTask(orchestrator) → AgentEnd OrchestratorAdapter
  │
  ├── Orchestrator 内部通过 BackendClient.run_task() 调子 Agent
  │   → Backend RunTask(aa) → 创建消息到 aa 的 session ← 第1次
  │
  └── 子 Agent 结果通过 Redis stream 回到 Orchestrator
      → Orchestrator 透传到自己的 SSE 流
      → Backend 流式刷写 → 创建消息到 orchestrator 的 session ← 第2次
```

Orchestrator 作为消息中转站，需要把子 Agent 回复透传给前端（前端订阅 orchestrator 的 SSE stream）。副作用是同一个内容被写入两个 session。

## 影响

1. **数据库冗余**：每个子 Agent 的回复占用双倍存储
2. **窗口查询重复注入**：`fetchGroupChatWindow` 查询其他 session 消息时，同一个子 Agent 回复出现两次，导致 `system_prompt_append` 包含重复内容
3. **前端展示可能重复**：如果前端按 session 分别展示，可能出现重复消息

## 可能的修复方案

### 方案 A：窗口查询排除 Orchestrator session（简单但不彻底）

在 `fetchGroupChatWindow` 中额外排除 `agent_type = "orchestrator"` 的 session。

- ✅ 修复窗口注入重复
- ❌ 数据库仍然双写

### 方案 B：Orchestrator 转发时不写 MySQL（彻底）

让 Orchestrator 透传子 Agent 结果时，只走 Redis stream（前端实时可见），不触发 MySQL 持久化。需要区分"原始消息"和"转发消息"。

- ✅ 从根源消除双写
- ❌ 改动较大，需要修改流式刷写逻辑

### 方案 C：前端按 message_id 去重

前端展示时按 `message_id` 或 `content + agent_name + timestamp` 去重。

- ✅ 改动最小
- ❌ 只是隐藏问题，不解决数据冗余

## 相关文件

- `backend/internal/handler/task.go` — RunTask handler，创建消息
- `backend/internal/stream/` — SSE 流式刷写，将 Redis stream 事件持久化到 MySQL
- `agentend/src/adapters/orchestrator.py` — OrchestratorAdapter，透传子 Agent 结果
- `agentend/src/orchestrator/execution/engine.py` — ExecutionEngine，通过 BackendClient 调子 Agent
