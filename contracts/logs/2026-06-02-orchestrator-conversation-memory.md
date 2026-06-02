# 2026-06-02 Orchestrator 对话记忆持久化

**类型**: agentend-only（无跨端契约变更）

## 变更说明

Orchestrator 新增 `ConversationMemoryStore`，将对话链（HumanMessage + AIMessage + ToolMessage）持久化到 `{shared_dir}/memory/conversation_memory.json`，跨轮保留最近 10 轮。

同时将 4 个动态上下文（Pin、Evolution、replan、群聊上下文）从系统提示词拆分到消息列表，按语义使用 SystemMessage / HumanMessage。用户 query 从系统提示词移至独立 HumanMessage。

## 影响范围

- 仅影响 `agentend` 内部逻辑
- 无 `contracts/schemas/` 变更
- 前后端 API 协议不变
