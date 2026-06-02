# 2026-06-02 CLI Adapter LangSmith Trace 接入

## 变更原因

Phase 5.2：CLI Adapter（Claude Code / OpenCode / Codex）需要 LangSmith trace 能力，用于调优提示词时查看 Agent 的完整决策过程和工具调用。

## 变更文件

无。未修改 `contracts/schemas/` 中的任何 YAML 文件。

## 对比结果

不适用。

## 跨端影响

无。本次改动为 agentend 内部变更：
- 新建 `agentend/src/adapters/trace.py`（RunTree trace wrapper）
- 修改 `agentend/src/api/v1/agent.py`（`_execute_stream` 接入 trace wrapper）

Wrapper 对事件流完全透明，不改变 SSE 事件格式、请求/响应结构，Frontend 和 Backend 无需任何改动。

## 契约变更

无。
