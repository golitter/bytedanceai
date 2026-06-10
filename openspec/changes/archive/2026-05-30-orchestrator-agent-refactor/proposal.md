## Why

当前 Orchestrator 是一个无状态的线性 Pipeline（discover → select → load_l2 → plan → write_shared），所有请求都走完整规划流程，无法区分闲聊和任务，无跨轮次记忆。需要重构为有记忆的 Agent 模式：LLM 自主判断直接回复还是触发多 Agent 编排，并记住每轮对话历史。

## What Changes

- **REASON 节点替代 Pipeline**：将 5 节点线性管道合并为单个 REASON 节点（含渐进式 skill 加载 + LLM tool-calling 循环），LLM 通过 `plan_and_dispatch` 工具调用信号表达编排意图，通过纯文本输出表达闲聊回复
- **LangGraph 生命周期状态图**：新建 REASON → DISPATCH → EXECUTE → REVIEW → EVOLVE → SAVE_MEM 的完整状态图，替代手写 if/for 流程
- **跨轮次 Memory**：使用 LangGraph MemorySaver 管理对话历史 + 工具调用 + 结果，支持上下文引用
- **Execute 子图（Wave Executor）**：按依赖关系拓扑分波次执行任务，预留 DAG 并行能力
- **REVIEW 重规划**：任务失败时带失败上下文重入 REASON 节点

## Capabilities

### New Capabilities

- `reason-node`: REASON 节点——LLM tool-calling 循环，自主判断输出文本（闲聊）vs 调用 plan_and_dispatch 工具（编排），含渐进式 skill 加载
- `agent-memory`: Orchestrator 跨轮次记忆——LangGraph MemorySaver 持久化对话历史、工具调用记录、执行结果
- `wave-executor`: Execute 子图——按依赖拓扑分波次执行任务，波次内并行，波次间串行
- `lifecycle-graph`: LangGraph 生命周期状态图——REASON → DISPATCH → EXECUTE → REVIEW → EVOLVE → SAVE_MEM，含条件路由和重规划循环

### Modified Capabilities

- `orchestrator-planning`: 从线性 Pipeline（5 节点）重构为 REASON 节点（合并内部步骤），prompt 支持闲聊 + 编排双模式
- `task-dispatcher`: 新增拓扑分波（topological sort），将 DispatchResult 列表按 depends_on 分为 execution_waves

## Impact

- **AgentEnd**：`orchestrator/planning/graph.py`（核心重构）、`orchestrator/planning/tools.py`（新增 plan_and_dispatch）、`orchestrator/planning/prompts.py`（双模式 prompt）、`orchestrator/execution/wave.py`（新增）、`adapters/orchestrator.py`（重构适配新图）
- **依赖**：新增 `langgraph.checkpoint`（MemorySaver），已有 `langgraph` 依赖
- **向后兼容**：OrchestratorAdapter 对外接口（stream_chat/chat）不变，内部实现重构
- **前端**：需识别 Orchestrator 直接回复（无 PLANNING 前缀）vs 编排流程
- **Go Backend**：少量改动，SSE 透传新事件类型
