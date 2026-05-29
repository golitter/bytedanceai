# 2026-05-29 orchestrator Agent 生命周期重构

## 变更原因

将 Orchestrator 从无状态线性 Pipeline（discover→select→load_l2→plan→write_shared）重构为有记忆的 Agent 生命周期图（skill_prepare→reason→dispatch→execute→review→evolve→save_mem），支持闲聊/编排双模式、跨轮次记忆、拓扑分波执行和失败重规划。

## 变更文件

- `agentend/src/orchestrator/planning/graph.py` — 重写：新 GraphState + skill_prepare/reason/dispatch/review/evolve/save_mem 节点 + MemorySaver + conditional edges
- `agentend/src/orchestrator/planning/prompts.py` — PLAN_PROMPT → REASON_PROMPT 双模式（闲聊直接回复 + 编排调用 plan_and_dispatch）
- `agentend/src/orchestrator/planning/tools.py` — 新增 plan_and_dispatch 工具
- `agentend/src/orchestrator/planning/skill_loader.py` — 提取 select_skills/load_l2_content 为独立函数
- `agentend/src/orchestrator/execution/dispatcher.py` — 新增 topological_sort 按依赖分波
- `agentend/src/orchestrator/execution/wave.py` — 新增 wave executor 子图
- `agentend/src/adapters/orchestrator.py` — 重写 stream_chat：Lifecycle 图 + 实时流式执行 + 错误兜底

## 对比结果

无 schema 变更。EventType 枚举、StreamEvent 结构均未修改。

## 跨端影响

- **Frontend**: 无影响。SSE 事件流格式不变（TEXT/PLANNING/DONE/ERROR/RUNTIME_* 均为已有事件类型）。PLANNING 事件的 node 值从 discover/select/load_l2/plan/write_shared 变为 skill_prepare/reason/dispatch，前端仅透传 node 字段，不依赖具体值。
- **Backend**: 无影响。后端 SSE scanner 透传事件，无需修改。Go HTTP 客户端 60s 超时不变，skill_prepare 节点保证在前几秒内产出首个 SSE 事件保活。
- **AgentEnd**: 内部重构。OrchestratorAdapter 对外接口（stream_chat/chat）签名不变。

## 契约变更

无。本次改动为 agentend 内部实现变更，不涉及 `contracts/schemas/` 中的任何契约定义。

## SSE 事件行为变化

| 场景 | 旧行为 | 新行为 |
|------|--------|--------|
| 闲聊 | 仍走完整 Pipeline，PLANNING(plan) 返回空 tasks | reason 节点直接返回 TEXT 事件，跳过 dispatch/execute |
| 编排 | PLANNING(discover/select/load_l2/plan/write_shared) | PLANNING(skill_prepare/reason/dispatch) + RUNTIME_* + TEXT |
| 执行 | engine.execute 直接 yield 事件 | _stream_wave AsyncIterator 实时流式产出 |
| 错误 | graph 异常时 SSE 断开无 DONE | try/except 兜底，发送 ERROR + DONE |
