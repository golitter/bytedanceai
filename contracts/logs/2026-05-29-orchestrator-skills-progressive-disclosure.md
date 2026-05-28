# 2026-05-29 orchestrator-skills 渐进式披露

## 变更原因

Orchestrator 规划模块新增 skill 渐进式披露能力（L1→L2→L3）和 tool-calling agent loop，需要记录对跨端的影响评估。

## 变更文件

- `agentend/src/adapters/orchestrator.py` — 规划前调用 SkillProvisioner
- `agentend/src/orchestrator/planning/graph.py` — 新增 discover/select/load_l2 节点 + plan_node agent loop
- `agentend/src/orchestrator/planning/prompts.py` — prompt 注入 L2 指令
- `agentend/src/orchestrator/planning/skill_loader.py` — 新增（L1/L2/L3 加载）
- `agentend/src/orchestrator/planning/tools.py` — 新增（5 个工具定义）
- `agentend/agents.json` — orchestrator 添加 config_dir
- `agentend/src/app/config.py` — SkillsConfig 添加 builtin_dir_resolved

## 对比结果

无 schema 变更。AgentType 枚举、event-types、request/response 结构均未修改。

## 跨端影响

- **Frontend**: 无影响。SSE 事件流格式不变（PLANNING 事件仅新增 discover/select/load_l2 节点名，前端已有 node 字段透传逻辑）。
- **Backend**: 无影响。不涉及后端 API 或数据模型变更。
- **AgentEnd**: 内部重构。图结构从 2 节点扩展为 5 节点，plan_node 从单次 LLM 调用改为 tool-calling agent loop。对外接口（stream_chat）签名和 SSE 事件格式不变。

## 契约变更

无。本次改动为 agentend 内部实现变更，不涉及 `contracts/schemas/` 中的任何契约定义。
