## Context

AgentEnd 已有完整的适配器体系（ClaudeCodeAdapter / OpenCodeAdapter），会话管理，工作区隔离（Git Worktree），技能分发（taskctl）。但缺少"多 Agent 协作"能力 — 无法在一次请求中协调多个 Agent 分工。

现有 taskctl 已支持 Agent 读写 `shared/.agent/` 目录（config.yaml、plans/、memory/）。Orchestrator 的 Phase 1 定位为"AI 项目经理"：理解需求 → 拆任务 → 写文件到 shared/ → 结束。后续 Phase 再接调度和聚合。

## Goals / Non-Goals

**Goals:**
- 实现 OrchestratorAdapter，作为 BaseAgentAdapter 的子类，通过 `agent_type: "orchestrator"` 触发
- 用 LangGraph StateGraph 驱动 plan → write_shared 线性流程
- 用 LLM structured output（`with_structured_output`）生成任务拆解，代替手写 JSON 解析
- 将 overview.md + tasks/*.md + config.yaml 写入 shared/.agent/ 目录
- config.yaml 作为声明式任务索引（类似 docker-compose.yaml），不存 runtime state
- 前端可通过现有 SSE/execute 接口接收规划过程事件

**Non-Goals:**
- 不实现 Agent 调度（Phase 2，由 Go Backend 负责）
- 不实现结果聚合（Phase 4，aggregate node）
- 不实现 DAG 依赖（Phase 5，depends_on）
- 不实现 LangGraph checkpoint / interrupt / human-in-the-loop
- 不修改 taskctl（Phase 2 再增强）
- 不实现自建 Agent（CustomAgentAdapter，独立 change）

## Decisions

### Decision 1: Orchestrator 只负责 Planning，不调度 Agent

**选择**: Orchestrator 写完 shared/ 就结束，调度由外部（Go Backend）负责。

**替代方案**: Orchestrator 内部调用其他 adapter（`registry.get("claude-code")()`）。

**理由**: 避免进程内耦合 — Orchestrator 不需要知道 adapter 的生命周期；天然隔离（一个 Agent 崩溃不影响其他）；调度策略（并行/串行/重试）独立演进。这也符合 filesystem-as-coordination-layer 的架构核心。

### Decision 2: 保留 LangGraph，即使 Phase 1 是线性流程

**选择**: 用 StateGraph 实现 plan → write_shared。

**替代方案**: 纯 async 函数 `plan() + write_shared()`。

**理由**: Phase 2-5 需要 aggregate node、interrupt（等 Agent 完成）、checkpoint（恢复状态）、conditional edges（按依赖路由）。现在用 LangGraph 打好基础，后续加 node/edge 即可，不需要重写。代价只是多引入 langgraph 依赖。

### Decision 3: 一次 LLM 调用完成规划（合并 analyze + plan）

**选择**: 单个 `plan` node，一次 LLM structured output 调用同时生成 overview + tasks。

**替代方案**: analyze node + plan node 两次调用。

**理由**: 两次调用增加延迟和 cost，且 Phase 1 不需要"先分析再规划"的分离。如果后续需要更复杂的推理，可以在 graph 中加 node。

### Decision 4: 用 `with_structured_output()` 而不是手写 JSON 解析

**选择**: `ChatAnthropic.with_structured_output(PlanOutput)` 让 LLM 直接输出 Pydantic model。

**替代方案**: prompt 要求输出 JSON + `_parse_json_response()` 手动解析。

**理由**: structured output 自动处理 markdown 代码块包裹、trailing comma、字段缺失等问题。Anthropic API 原生支持 tool use 模式的 structured output，稳定性远高手写解析。

### Decision 5: 文件名后端生成，不信任 LLM

**选择**: `write_shared_node` 中用 `task-{idx:03d}.md` 生成文件名。

**替代方案**: 让 LLM 在 TaskDef 中输出 `file` 字段。

**理由**: LLM 可能输出非法路径（`../../etc/passwd`）、重名、或非标准格式。后端生成保证一致性和安全性。

### Decision 6: config.yaml 是声明式任务索引

**选择**: config.yaml 只存 `task_id` + `tasks` 列表（id, agent, file），不存 status/session_id/depends_on。

**理由**: config.yaml 是所有组件（Orchestrator、Scheduler、taskctl、Aggregator）之间的 ABI。如果它变成 runtime state（存 status、session），会和调度器耦合，无法恢复、重试、重新调度。保持声明式（像 k8s manifest）是正确边界。

## Risks / Trade-offs

- **[LLM 拆解质量不稳定]** → Phase 1 的 prompt 中加入约束规则（最多 5 个任务、唯一 agent、具体可执行），后续可通过 few-shot 示例或 fine-tuning 改善
- **[structured output 解析失败]** → Anthropic 对 tool use 模式支持好，但极端情况可能失败；需要在 adapter 层加 try/catch，返回 ERROR event
- **[shared/ 目录路径需调用方提供]** → Phase 1 由 API 调用方（Go Backend 或 curl）提供 `shared_dir` 绝对路径，Orchestrator 不自己计算路径
- **[taskctl 读 plans/ vs Orchestrator 写 tasks/ 目录名不一致]** → Phase 1 先用 `tasks/`，跑通后统一调整 taskctl 或 Orchestrator 的目录名
