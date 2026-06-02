# 2026-06-02 Codex Adapter 拼接 system_prompt_append

## 变更原因

Codex Adapter 之前丢弃了 `system_prompt_append` 参数，导致 rule 引擎注入的系统约束（安全规则、群聊上下文、输出技能等）无法传递给 Codex CLI。

## 变更文件

无。未修改 `contracts/schemas/` 中的任何 YAML 文件。

## 对比结果

不适用。

## 跨端影响

无。本次改动为 agentend 内部变更：
- 修改 `agentend/src/adapters/codex.py` 的 `_build_command` 和 `stream_chat`

与 OpenCode Adapter 行为对齐：有 `system_prompt_append` 时拼成 `[系统约束: ...]` 前缀与 message 合并。

## 契约变更

无。
