## Architecture Overview

多 Agent 协作的核心链路：`用户消息 → Backend → AgentEnd Planner → 真实 Agent 执行 → 事件流回前端`。改动集中在三个层：

```
┌─────────────────────────────────────────────────────────────┐
│  Frontend                                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │PlanCard  │  │ToolCard  │  │StatusBdg │  │CoordChan │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
│       └──────────────┴─────────────┴─────────────┘          │
│                         │ SSE events                         │
├─────────────────────────┼───────────────────────────────────┤
│  Backend (Go)           │                                   │
│  stream.go ─────────────┤ 透传所有事件类型，不额外过滤       │
│                         │                                   │
├─────────────────────────┼───────────────────────────────────┤
│  AgentEnd (Python)      │                                   │
│  ┌─────────────────────┐│                                   │
│  │ OrchestratorAdapter ││                                   │
│  │  Phase 1: Planning  ││  ✅ 已有 (LangGraph)              │
│  │  Phase 2: Dispatch  ││  🔧 核心改动：mock → 真实调用      │
│  │  Phase 3: Collect   ││  🔧 新增：执行引擎 + 事件归一化    │
│  │  Phase 4: Aggregate ││  ✅ 已有 (LLM 汇总)               │
│  │  Phase 5: Record    ││  ✅ 已有 (EvolutionStore)          │
│  └─────────────────────┘│                                   │
│       │                  │                                   │
│  ┌────┴────┐  ┌─────────┐│                                   │
│  │ClaudeCode│  │OpenCode ││  真实 Agent Adapter              │
│  │ Adapter  │  │ Adapter ││                                   │
│  └─────────┘  └─────────┘│                                   │
└─────────────────────────────────────────────────────────────┘
```

## Data Model Changes

### 1. 新增 EventType 枚举值

在 `contracts/schemas/event-types.yaml` 中扩展：

```
现有（保留）：
  init, text, tool_call, tool_result, artifact, planning, done, error

新增：
  runtime_executing    # Agent 开始执行任务
  runtime_completed    # Agent 完成任务
  coordination_start   # 协调通道开启
  coordination_message # 协调消息（项目经理 ↔ Agent）
  coordination_done    # 协调通道关闭
```

### 2. OrchestratorAdapter 执行引擎

核心数据流（Phase 2-3 重构）：

```
PlanOutput.tasks
    │
    ▼ (串行遍历)
┌─────────────────────────────────────────────┐
│ for task in plan.tasks:                     │
│                                             │
│   1. yield runtime_executing event          │
│      {task_id, agent, title, status}        │
│                                             │
│   2. adapter = registry.get(task.agent)     │
│      instance = adapter()                   │
│      await instance.create_session(sid)     │
│                                             │
│   3. async for event in                     │
│        instance.stream_chat(sid, content):  │
│      yield event  (透传 Agent 的流式事件)    │
│                                             │
│   4. yield runtime_completed event          │
│      {task_id, agent, success, duration}    │
│                                             │
└─────────────────────────────────────────────┘
```

### 3. 协调通道数据流

```
Planner 生成初版计划后 → 协调通道介入

┌──────────────────┐
│ CoordinationChan │
│                  │
│ for agent in     │
│   relevant_agents│
│                  │
│   yield coord_   │  ← 协调开始事件
│   start          │
│                  │
│   question =     │  ← LLM 生成针对该 agent 的问题
│   llm.generate() │
│                  │
│   yield coord_   │  ← 转发问题事件
│   message(Q→A)   │
│                  │
│   answer =       │  ← adapter.chat() 非流式
│   adapter.chat() │
│                  │
│   yield coord_   │  ← 转发回答事件
│   message(A→O)   │
│                  │
│   yield coord_   │  ← 协调结束
│   done           │
└──────────────────┘
```

## Component Changes

### AgentEnd — OrchestratorAdapter

**文件**: `agentend/src/adapters/orchestrator.py`

改动范围：Phase 2 (Dispatch) 和 Phase 3 (Collect) 的 lines 62-90

- OrchestratorAdapter 构造函数接收 `AdapterRegistry` 实例（或通过全局获取）
- Dispatch 阶段保持不变（Dispatcher 只做 plan → DispatchResult 映射）
- Collect 阶段替换为 ExecutionEngine：
  - 串行遍历 DispatchResult 列表
  - 对每个 task：从 registry 获取 adapter → create_session → stream_chat → 收集事件
  - 透传 agent 的流式事件（text, tool_call, tool_result 等）
  - 在 task 开始/结束时 yield runtime_executing / runtime_completed
  - 收集 TaskResult 供 Phase 4 Aggregate 使用

### AgentEnd — 新增 ExecutionEngine

**新文件**: `agentend/src/orchestrator/execution.py`

职责：
- 接收 DispatchResult 列表 + AdapterRegistry
- 串行执行每个 task
- yield 所有流式事件（agent 事件 + runtime 控制事件）
- 返回 TaskResult 列表

```python
class ExecutionEngine:
    def __init__(self, registry: AdapterRegistry, workspace_root: str): ...

    async def execute(
        self, dispatches: list[DispatchResult]
    ) -> AsyncIterator[tuple[StreamEvent, TaskResult | None]]:
        """Yield (event, optional_result) for each step.
        TaskResult is non-None only when a task completes."""
```

### AgentEnd — 新增 CoordinationChannel

**新文件**: `agentend/src/orchestrator/coordination.py`

职责：
- Planner 生成计划后，向相关 agent 发起单轮 Q&A
- 使用 LLM 生成针对性问题
- 使用 `adapter.chat()` 获取回答
- yield coordination 系列事件
- 返回协调结论文本，供 Planner 参考调整

### AgentEnd — OrchestratorAdapter 新流程

```python
async def stream_chat(self, session_id, message, **kwargs):
    # Phase 1: Planning (已有)
    ...

    # Phase 1.5: Coordination (新增)
    channel = CoordinationChannel(registry, settings.llm)
    coord_context = ""
    async for event in channel.coordinate(plan, agents):
        yield event
    coord_context = channel.summary()

    # Phase 2+3: Execute (重构)
    engine = ExecutionEngine(registry, workspace_root)
    task_results = []
    async for event, result in engine.execute(dispatch_results):
        yield event
        if result:
            task_results.append(result)

    # Phase 4: Aggregate (已有)
    ...

    # Phase 5: Record (已有)
    ...
```

### Contracts — event-types.yaml

**文件**: `contracts/schemas/event-types.yaml`

新增 5 种事件类型定义，每种事件明确其 content 字段结构：

| 事件类型 | content 字段 |
|----------|-------------|
| `runtime_executing` | `{task_id, agent, title, status: "running"}` |
| `runtime_completed` | `{task_id, agent, success, duration, status: "completed"/"failed"}` |
| `coordination_start` | `{round: 1, agents: [...]}` |
| `coordination_message` | `{from, to, text, round}` |
| `coordination_done` | `{rounds, decisions: [...]}` |

### Backend — SSE 透传

**文件**: `backend/internal/stream/writer.go`

当前 StreamWriter 检测特定事件类型（text, error, done）做特殊处理，其他类型直接 XAdd 到 Redis Stream。新增的 `runtime_*` 和 `coordination_*` 事件不需要特殊处理，会自动透传。

需确认：`backend/internal/generated/events.go` 中 EventType 常量列表是否需要同步新增（运行 `make generate` 后自动更新）。如果 Backend 代码有 switch-case 只处理已知类型，需要加 default 分支透传。

### Frontend — 事件渲染扩展

**文件变更范围**:

1. **`frontend/src/hooks/use-chat-stream.ts`** — 新增事件处理 case
2. **`frontend/src/stores/chat.ts`** — 新增 streaming 状态字段（coordination blocks、runtime status）
3. **`frontend/src/lib/block-reducer.ts`** — 新增 plan、coordination、runtime block 类型
4. **新组件**:
   - `frontend/src/components/cards/PlanCard.tsx` — 任务分配卡片
   - `frontend/src/components/cards/ToolCard.tsx` — 工具调用卡片（已有类似，可能复用）
   - `frontend/src/components/cards/CoordChannel.tsx` — 可折叠协调通道
   - `frontend/src/components/cards/RuntimeStatus.tsx` — 状态徽章

### Frontend — 事件到 Block 的映射

```
event.type                    → block.type         → component
─────────────────────────────────────────────────────────────
runtime_executing             → runtime_status     → RuntimeStatus
runtime_completed             → runtime_status     → RuntimeStatus
coordination_start            → coordination       → CoordChannel
coordination_message          → coordination       → CoordChannel (append)
coordination_done             → coordination       → CoordChannel (close)
planning (node=dispatch)      → plan               → PlanCard
tool_call                     → tool_call          → ToolCard
tool_result                   → tool_result        → ToolCard
text                          → text               → (现有渲染)
```

## Key Design Decisions

### D1: 串行调度，不并行

首次实现使用串行执行。理由：
- 并行需要工作区隔离（每个 agent 独立 worktree），增加合并复杂度
- 串行更易调试，事件流天然有序
- 后续迭代可升级为并行

### D2: adapter 实例生命周期

每个 task 执行时创建新的 adapter 实例：
- `adapter_cls = registry.get(agent_type)` → `adapter = adapter_cls()` → `await adapter.create_session(new_session_id)` → `stream_chat` → `destroy_session`
- 不复用跨 task 的 adapter 实例，避免状态污染

### D3: 协调通道单轮 Q&A

先不做多轮协商：
- 项目经理根据 plan 内容，针对每个相关 agent 生成 1 个问题
- 使用 `adapter.chat()` 非流式获取回答
- 收集所有回答后汇总为协调结论
- 协调结论作为补充 context 传给各 agent 的执行 prompt

### D4: 事件命名 dot notation

统一使用 `runtime.executing` 风格（dot notation），但 YAML key 使用 underscore（`runtime_executing`），因为代码生成后 EventType 枚举值使用 underscore 更符合各语言惯例。

### D5: Backend 最小改动

Backend 不解析新事件类型的 content，只做 SSE 透传：
- StreamWriter 的 XAdd 逻辑已有，新事件类型直接走现有路径
- 运行 `make generate` 自动更新 `generated/events.go` 的 EventType 常量
- 如果 Backend 有类型枚举校验，改为 string pass-through

## Risks

| 风险 | 缓解 |
|------|------|
| Agent 执行超时导致整个 stream 卡住 | ExecutionEngine 设置单 task 超时（默认 300s），超时 yield error + continue |
| Adapter 实例化失败（agent 进程不可用） | try-catch 单个 task，标记 failed，继续执行后续 task |
| 协调通道 LLM 调用增加延迟 | 单轮 Q&A，每个 agent 一次 chat() 调用，可接受 |
| 新事件类型前端未识别 | 未识别的事件静默忽略（block-reducer 加 default 分支） |
| 并发 session_id 冲突 | 使用 `f"orch-{task_id}-{dispatch.task_id}"` 格式确保唯一 |
