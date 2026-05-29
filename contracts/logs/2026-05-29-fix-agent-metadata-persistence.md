# 契约变更：修复 agent 元数据在 SSE 合并时丢失 + 无 Sub-Agent 时 overview 重复输出

## 变更原因

TEXT 事件的 `agent` + `agent_type` 元数据在 Backend `StreamWriter.flushTextBuffer()` 合并发布到 Redis 时被丢弃（`FormatSSE` 只保留 `text`），导致前端无法检测 Agent 切换，所有内容合并为一条 Orchestrator 消息。此外，`currentAgentName` 在 `agent_type` 未变时不同步更新，以及无 Sub-Agent 时 Phase 4 聚合回退到 overview 导致重复输出。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。修复均为实现层面。

## 修复内容

### 1. `flushTextBuffer()` 保留 agent 元数据

已在前序提交中将 `FormatSSE()` 替换为 `FormatSSEWithMeta()`，此处不再赘述。

### 2. `currentAgentName` 同步更新（writer.go）

TEXT 事件处理新增：当 `agent_type` 未变但 TEXT 事件携带 `agent` 名称时，同步更新 `currentAgentName`，确保 Redis SSE 携带完整元数据。

### 3. 无 Sub-Agent 时跳过 Phase 4 TEXT（orchestrator.py）

`Aggregator.aggregate([], overview)` 返回空字符串时，不再 `yield TEXT(text=aggregated or overview)`，避免 overview 被输出两次。

## 跨端影响

- **Backend**: `writer.go` TEXT 事件处理逻辑变更
- **AgentEnd**: `orchestrator.py` Phase 4 条件输出
- **Frontend**: 无变更（行为修正后前端现有逻辑即可正确工作）
- **Contracts**: 无 schema 变更
