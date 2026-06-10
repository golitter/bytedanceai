## Context

当前 `src/orchestrator/graph.py` 实现了一个 LangGraph 两节点线性管道（plan → write_shared），产出的 `config.yaml` + `plans/*.md` 文件写在 `shared/.agent/` 目录后流程就结束了。没有任何机制把任务分发给 Agent、收集结果、或汇报。

现有基础设施已经完备：Adapter 统一抽象（ClaudeCodeAdapter / OpenCodeAdapter）、Git Worktree 隔离、Session 管理、Rule Engine、taskctl Go 工具。Orchestrator 需要在这个基础上补全闭环。

## Goals / Non-Goals

**Goals:**
- 让 Orchestrator 从"写文件就结束"升级为"plan → dispatch → collect → aggregate"闭环
- Dispatcher 以 @agent 群聊模式产出调度指令，与前端 IM 交互范式对齐
- Pin Memory 复用现有 `common/` 目录，`_pins.yaml` 做书签，动态 pin/unpin
- Self-Evolution 简单记录成败经验，注入 Planner prompt
- 最小改动，不破坏现有 Adapter / API / Workspace 体系

**Non-Goals:**
- 不做 Event Sourcing / Reducer / Replay
- 不做 Workflow State Machine（pause/resume/cancel）
- 不做 Resource Governor / Token Budget
- 不做 Dynamic DAG / Conflict Semantic Merge
- 不做 Reflection Agent / 自动 Prompt 进化
- 不去 LangGraph 化（现有 graph.py 可继续使用，后续再考虑）

## Decisions

### D1: Dispatcher 以 @ 模式产出调度指令，不直接 HTTP 调 Agent

**选择**：Dispatcher 产出结构化的 `DispatchResult` JSON（含 @agent、任务内容、workspace 路径），由前端/Go Backend 消费后触发实际执行。

**备选**：Dispatcher 直接 HTTP 调 `/v1/agent/execute`。

**理由**：与任务要求的 IM 群聊范式对齐。前端需要渲染 @ 消息卡片，Go Backend 可能需要做额外的调度/审批逻辑。直接调 HTTP 会绕过这些层。@ 模式让 Runtime 只负责决策，执行交给上层。

### D2: Pin Memory 复用 common/ 目录，不新建 pins/ 目录

**选择**：Pin 文件直接存放在 `shared/.agent/memory/common/`，`_pins.yaml` 作为书签层管理 pin 状态和 AI 摘要。Unpin 只删书签不删文件。

**备选**：独立 `pins/` 目录。

**理由**：减少目录层级。common/ 已经是所有 Agent 共享的知识空间（`taskctl common-memory` 已支持）。Pin 只是给某些 common 文件加上"自动注入 prompt"的标记。Unpin 后文件仍是有效的共享知识，不应丢失。

### D3: Self-Evolution 用简单 YAML 文件，不引入数据库

**选择**：`shared/.agent/evolution.yaml` 存最近 20 条编排经验，Planner prompt 注入最近 5 条。

**备选**：SQLite / Redis / Vector DB。

**理由**：MVP 阶段数据量极小（20 条记录），YAML 文件足够。LLM 本身就不稳定，过度设计的存储不会带来额外价值。

### D4: OrchestratorAdapter 内部实现闭环，不拆为独立服务

**选择**：在 `OrchestratorAdapter.stream_chat` 内部串联 plan → dispatch → collect → aggregate。

**备选**：新建独立的 OrchestratorService 或 Background Worker。

**理由**：保持与现有 Adapter 模式一致。`stream_chat` 已经是 async generator，可以在 yield 的过程中串联多个阶段。后续如果需要后台执行，再拆分。

## Risks / Trade-offs

- **[LLM 不稳定性]** Planner/Aggregator 的 LLM 调用可能超时或输出格式错误 → 沿用现有 `_extract_json` + 异常冒泡策略，后续加重试
- **[Token 成本]** Pin 摘要 + Evolution 经验注入 Prompt 会增加 token 消耗 → MVP 阶段量级可控（摘要 1-3 句 + 5 条经验），后续加 Budget 限制
- **[Pin AI 摘要延迟]** 每次 Pin 新内容时调 LLM 生成摘要会增加延迟 → 摘要生成是异步的，不阻塞主流程
- **[@ 调度模式需要前端配合]** Dispatcher 产出的 JSON 需要前端/Go Backend 消费 → 先定义好接口格式，前端侧另行实现
