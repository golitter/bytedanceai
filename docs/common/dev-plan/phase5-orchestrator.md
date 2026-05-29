# Phase 5: Orchestrator 群聊协作 — Agent 模式重构

> 目标: Orchestrator 重构为有记忆的 Agent，自主判断闲聊直接回复 vs 任务编排派发，通过 LangGraph 状态图管理完整生命周期。
> 预估: 5-6 天
> 前置: Phase 4 完成 (产物卡片渲染可用)

## 核心设计变更

### 旧设计（Pipeline 模式）

```
discover → select → load_l2 → plan → write_shared
(线性管道，无记忆，所有请求都走完整规划流程)
```

### 新设计（Agent 模式）

Orchestrator 是一个**有记忆的对话 Agent**：
- 闲聊 → LLM 直接回复文本
- 任务 → 调用 `plan_and_dispatch` 工具触发编排
- 每轮对话的历史、工具调用、结果都记入 Memory

```
┌──────────────────────────────────────────────────────┐
│             Orchestrator Agent Graph                  │
│                                                      │
│              ┌──────────┐                            │
│              │  REASON  │ ← LLM + Memory + Skills    │
│              └────┬─────┘                            │
│                   │                                  │
│              输出类型?                                │
│          ┌────────┴────────┐                         │
│          │                 │                         │
│       文本回复          plan_and_dispatch             │
│       (闲聊/直接回答)    (任务编排工具调用)            │
│          │                 │                         │
│          ▼                 ▼                         │
│     ┌─────────┐     ┌──────────┐                    │
│     │  SAVE   │     │ DISPATCH │                    │
│     │  MEM    │     └────┬─────┘                    │
│     └────┬────┘          │                          │
│          │               ▼                          │
│          │        ┌──────────┐                       │
│          │        │ EXECUTE  │ (子图: wave executor) │
│          │        └────┬─────┘                       │
│          │             │                            │
│          │             ▼                            │
│          │       ┌──────────┐                        │
│          │       │  REVIEW  │                        │
│          │       └────┬─────┘                        │
│          │            │                              │
│          │       ┌────┴────┐                         │
│          │    ok │    replan│                        │
│          │       ▼        ▼                          │
│          │ ┌─────────┐ ┌────────┐                    │
│          │ │ EVOLVE  │ │ REASON │ ← 带失败上下文     │
│          │ └────┬────┘ └────────┘                    │
│          │      │                                    │
│          │      ▼                                    │
│          │  ┌─────────┐                              │
│          └─▶│ SAVE MEM│                              │
│             └────┬────┘                              │
│                  │                                   │
│                  ▼                                   │
│                 END                                  │
└──────────────────────────────────────────────────────┘
```

### 关键设计决策

1. **REASON 节点是 Agent 的核心** — LLM 自主判断输出文本还是调用 `plan_and_dispatch` 工具
2. **Memory 跨轮次持久化** — LangGraph MemorySaver 管理，累积对话历史 + 工具调用 + 结果
3. **`plan_and_dispatch` 是工具调用信号** — 不是输出格式约束，LLM 通过 tool_calls 表达编排意图
4. **Execute 是子图** — wave executor 按依赖分波次执行任务
5. **REVIEW 决定重规划** — 失败时带回上下文重入 REASON
6. **闲聊也是 Orchestrator 回复** — 不回落到其他 Agent，Orchestrator 自己回答

## 交付标准

### AgentEnd 验证

```bash
# 启动服务
cd agentend && uv run python -m src.app.main

# 测试 1: 闲聊 (REASON 输出文本)
curl -X POST http://localhost:8001/v1/agent/execute \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "test-001", "session_id": "orch-1",
    "message": "你好",
    "agent_type": "orchestrator",
    "config": {
      "agents": [
        {"id": "claude-code", "name": "Claude Code", "capabilities": "代码生成"},
        {"id": "opencode", "name": "OpenCode", "capabilities": "代码审查"}
      ],
      "task_id": "test-001",
      "shared_dir": "/tmp/orch-test/.agent"
    }
  }'
# 预期: Orchestrator 直接文本回复，不触发编排

# 测试 2: 任务编排 (REASON 调用 plan_and_dispatch)
curl -X POST http://localhost:8001/v1/agent/execute \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "test-002", "session_id": "orch-1",
    "message": "用 Claude Code 写一个 React 登录页，然后用 OpenCode 审查代码质量",
    "agent_type": "orchestrator",
    "config": {
      "agents": [
        {"id": "claude-code", "name": "Claude Code", "capabilities": "代码生成"},
        {"id": "opencode", "name": "OpenCode", "capabilities": "代码审查"}
      ],
      "task_id": "test-002",
      "shared_dir": "/tmp/orch-test/.agent"
    }
  }'
# 预期: 规划 → 派发 → 执行 → 汇总

# 测试 3: 引用上下文 (Memory 跨轮次)
curl -X POST http://localhost:8001/v1/agent/execute \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "test-003", "session_id": "orch-1",
    "message": "刚才的分析结果总结一下",
    "agent_type": "orchestrator",
    "config": {
      "agents": [...],
      "task_id": "test-003",
      "shared_dir": "/tmp/orch-test/.agent"
    }
  }'
# 预期: 从 Memory 中引用 test-002 的结果，直接文本回复
```

### 前端群聊 UI 验证

1. 选择 Orchestrator Agent
2. 发送 "你好" → Orchestrator 直接回复（无 PLANNING 事件）
3. 发送 "写登录页并审查" → 看到规划 → 多 Agent 执行 → 汇总
4. 发送 "总结刚才的结果" → Orchestrator 引用上下文回复

## 实现步骤

### Step 1: 新增 `plan_and_dispatch` 工具

**修改**: `agentend/src/orchestrator/planning/tools.py`

```python
@tool
def plan_and_dispatch(
    overview: str = Field(description="任务编排概述"),
    tasks: list[dict] = Field(description="任务列表，每个任务含 session_id/title/content"),
) -> str:
    """当你判断需要多 Agent 协作完成任务时，调用此工具生成编排计划。
    如果你认为可以直接回答用户，请不要调用此工具，直接输出文本即可。"""
    return "plan_generated"
```

工具定义本身不做任何事——它的返回值会被 REASON 循环拦截，提取参数构造 `PlanOutput`。

### Step 2: 重构 Graph — REASON 节点

**修改**: `agentend/src/orchestrator/planning/graph.py`

将现有的 `discover → select → load_l2 → plan → write_shared` 线性管道重构为 REASON 节点：

- 渐进式 skill 加载（保留现有逻辑，合并为函数内部步骤）
- LLM tool-calling 循环（保留现有 plan_node 的循环结构）
- 判断输出类型：tool_calls 中有 `plan_and_dispatch` → 编排，无 tool_calls → 文本回复

```python
def reason_node(state: GraphState) -> dict:
    # 1. 渐进式 skill 加载 (保留 discover → select → load_l2)
    l1 = discover_skills(skills_dir)
    selected = select_skills(l1, state["message"])
    l2 = load_l2_content(selected, skills_dir)

    # 2. 构建 prompt (注入 memory + pin + evolution)
    system_prompt = build_planner_prompt(
        agents_desc=..., message=state["message"],
        shared_dir=state["shared_dir"], l2_content=l2,
    )

    # 3. LLM tool-calling 循环
    tools = [read_file, write_file, list_dir, run_skill, plan_and_dispatch, ...]
    messages = [SystemMessage(system_prompt)] + state["memory_messages"] + [HumanMessage(state["message"])]

    for _ in range(MAX_ITERATIONS):
        response = llm_with_tools.invoke(messages)

        if not response.tool_calls:
            # 纯文本回复 (闲聊 / 直接回答)
            return {"output_type": "text", "text": response.content}

        for tc in response.tool_calls:
            if tc["name"] == "plan_and_dispatch":
                # 编排信号
                return {"output_type": "plan", "plan": PlanOutput(**tc["args"])}

            # 普通工具 → 执行，继续循环
            result = execute_tool(tc)
            messages.append(ToolMessage(content=result, tool_call_id=tc["id"]))

        messages.append(response)
```

### Step 3: 生命周期节点

**修改**: `agentend/src/orchestrator/planning/graph.py`

新增 Dispatch、Execute、Review、Evolve、SaveMem 节点：

```python
def dispatch_node(state: GraphState) -> dict:
    plan = state["plan"]
    dispatcher = Dispatcher(state["agents"])
    dispatch_results = dispatcher.dispatch(plan)
    # 拓扑分波
    execution_waves = topological_sort(dispatch_results)
    return {"dispatch_results": dispatch_results, "execution_waves": execution_waves}

def execute_node(state: GraphState) -> dict:
    # Wave executor 子图入口
    ...

def review_node(state: GraphState) -> dict:
    task_results = state.get("task_results", [])
    failed = [r for r in task_results if not r.success]
    if failed and state["iteration"] < state["max_iterations"]:
        return {"needs_replan": True, "replan_reason": format_failed_context(failed)}
    return {"needs_replan": False}

def evolve_node(state: GraphState) -> dict:
    evolution = EvolutionStore(state["shared_dir"])
    evolution.record(...)
    return {}

def save_mem_node(state: GraphState) -> dict:
    # 将本轮交互追加到 memory
    return {}
```

### Step 4: 构建新 Graph

**修改**: `agentend/src/orchestrator/planning/graph.py` → `build_graph()`

```python
def build_graph() -> CompiledGraph:
    graph = StateGraph(GraphState)

    graph.add_node("reason", reason_node)
    graph.add_node("dispatch", dispatch_node)
    graph.add_node("execute", execute_node)       # 子图
    graph.add_node("review", review_node)
    graph.add_node("evolve", evolve_node)
    graph.add_node("save_mem", save_mem_node)

    graph.set_entry_point("reason")

    # REASON 输出路由
    graph.add_conditional_edges("reason", route_by_output_type, {
        "text": "save_mem",
        "plan": "dispatch",
        "error": END,
    })

    graph.add_edge("dispatch", "execute")
    graph.add_edge("execute", "review")

    # REVIEW 路由
    graph.add_conditional_edges("review", route_by_review, {
        "ok": "evolve",
        "replan": "reason",
    })

    graph.add_edge("evolve", "save_mem")
    graph.add_edge("save_mem", END)

    return graph.compile(checkpointer=MemorySaver())
```

### Step 5: Memory 集成

**修改**: `agentend/src/orchestrator/planning/graph.py`

使用 LangGraph 的 MemorySaver 管理跨轮次状态：

```python
from langgraph.checkpoint.memory import MemorySaver

# GraphState 新增字段
class GraphState(TypedDict):
    # ... 现有字段 ...
    output_type: str                    # "text" | "plan"
    text: str                           # 文本回复内容
    memory_messages: Annotated[list, add]  # 累积的对话 + 工具消息
```

每轮对话结束后，messages 自动持久化到 checkpointer。下一轮 REASON 自动加载历史。

### Step 6: Execute 子图 (Wave Executor)

**新增**: `agentend/src/orchestrator/execution/wave.py`

```python
def build_execute_subgraph():
    """Execute 子图: 按 wave 分批执行任务"""
    sub = StateGraph(ExecuteState)
    sub.add_node("wave_execute", wave_execute_node)
    return sub.compile()

async def wave_execute_node(state):
    current_wave = state["execution_waves"][state["current_wave_idx"]]
    # 并行执行当前 wave 的所有任务
    results = await asyncio.gather(*[
        execute_single_task(task, state)
        for task in current_wave
    ])
    return {"task_results": results, "current_wave_idx": state["current_wave_idx"] + 1}
```

### Step 7: 重构 OrchestratorAdapter

**修改**: `agentend/src/adapters/orchestrator.py`

```python
class OrchestratorAdapter(BaseAgentAdapter):
    def __init__(self):
        self._graph = build_graph()

    async def stream_chat(self, session_id, message, **kwargs):
        config = {"configurable": {"thread_id": session_id}}  # Memory key

        async for chunk in self._graph.astream(
            {
                "message": message,
                "agents": kwargs["agents"],
                "task_id": kwargs["task_id"],
                "shared_dir": kwargs["shared_dir"],
            },
            config=config,
            stream_mode="updates",
        ):
            node_name, node_output = next(iter(chunk.items()))

            if node_name == "reason":
                if node_output.get("output_type") == "text":
                    yield StreamEvent.create(EventType.TEXT, text=node_output["text"],
                                             agent="Orchestrator", agent_type="orchestrator")
                else:
                    yield StreamEvent.create(EventType.PLANNING, status="completed")

            elif node_name == "dispatch":
                for dr in node_output.get("dispatch_results", []):
                    yield StreamEvent.create(EventType.PLANNING, node="dispatch", dispatch=dr.model_dump())

            elif node_name == "execute":
                yield ...  # 透传子 Agent 执行事件

            elif node_name == "save_mem":
                yield StreamEvent.create(EventType.DONE, text=node_output.get("summary", ""))
```

### Step 8: Go Backend 适配

**修改**: `backend/internal/handler/stream.go`

- 当 agent_type=orchestrator 时，SSE 透传新事件类型
- Orchestrator 侧完整闭环，Go 不参与调度逻辑

### Step 9: 前端群聊 UI

**修改**: `frontend/src/stores/chat.ts`

- Orchestrator TEXT 事件无 PLANNING 前缀 → 直接渲染为普通消息
- PLANNING 事件 → 规划进度卡片
- 子 Agent TEXT 事件 → 带 Agent 标签的消息块

**修改**: `frontend/src/components/chat/MessageBubble.tsx`

- 不同 Agent 用不同颜色标签（Claude Code: 蓝 / OpenCode: 绿 / Orchestrator: 紫）

**新增**: `frontend/src/components/chat/PlanningCard.tsx`

- 规划进度卡片

## State Schema

```python
class GraphState(TypedDict):
    # 输入层（不可变）
    message: str
    agents: list[dict]
    task_id: str
    shared_dir: str
    allowed_read_dirs: list[str]

    # REASON 产出
    output_type: str                            # "text" | "plan"
    text: str                                   # 文本回复
    plan: PlanOutput | None                     # 编排计划

    # DISPATCH 产出
    dispatch_results: list[DispatchResult]
    execution_waves: list[list[DispatchResult]]

    # EXECUTE 产出（累积）
    task_results: Annotated[list[TaskResult], add]
    task_status: dict[str, str]                 # task_id → RUNNING|DONE|FAILED

    # REVIEW 决策
    needs_replan: bool
    replan_reason: str

    # 聚合
    summary: str

    # 元信息
    iteration: Annotated[int, add_one]          # 防无限重规划
    max_iterations: int                         # 默认 3

    # Memory
    memory_messages: Annotated[list, add]       # 累积对话历史
```

## Memory 设计

使用 LangGraph MemorySaver（简单方案）：

```
Orchestrator Memory (per session, LangGraph checkpointer 管理)
│
├── messages: 累积的对话 + 工具调用 + 结果
│   turn 1: [Human("你好"), AI("你好！我是编排器...")]
│   turn 2: [Human("分析项目"), AI(plan_tool_call), ToolMsg(results...)]
│   turn 3: [Human("基于上面的结果细化"), AI(new_plan), ToolMsg(results...)]
│
├── evolution: 规划经验（已有，evolution.yaml）
└── pins: 共享约束（已有，_pins.yaml）
```

REASON 节点构建 prompt 时：
- 从 checkpointer 加载历史 messages
- 注入 pin constraints
- 注入 evolution context
- 渐进式 skill 加载

## 文件清单

```
AgentEnd:
├── src/
│   ├── orchestrator/
│   │   ├── planning/
│   │   │   ├── graph.py              # ✏️ 重构: REASON 节点 + 生命周期图
│   │   │   ├── tools.py              # ✏️ 新增 plan_and_dispatch 工具
│   │   │   ├── prompts.py            # ✏️ 更新 prompt (支持闲聊 + 编排)
│   │   │   └── skill_loader.py       # ✅ 保持
│   │   ├── execution/
│   │   │   ├── engine.py             # ✏️ 适配子图调用
│   │   │   ├── wave.py               # 🆕 Wave executor 子图
│   │   │   ├── dispatcher.py         # ✏️ 新增拓扑分波
│   │   │   ├── state.py              # ✏️ 适配新 GraphState
│   │   │   └── coordination.py       # ✅ 保持
│   │   ├── memory/
│   │   │   ├── evolution.py          # ✅ 保持
│   │   │   └── pin_memory.py         # ✅ 保持
│   │   ├── reporting/
│   │   │   └── aggregator.py         # ✅ 保持
│   │   └── models.py                 # ✏️ 新增 output_type 字段
│   ├── adapters/
│   │   └── orchestrator.py           # ✏️ 重构: 适配新图结构
│   └── schemas/
│       └── events.py                 # ✅ 保持

Go Backend:
├── internal/
│   └── handler/
│       └── stream.go                 # ✏️ 透传新事件类型

Frontend:
├── src/
│   ├── stores/chat.ts                # ✏️ 识别 Orchestrator 直接回复 vs 编排
│   └── components/chat/
│       ├── MessageBubble.tsx         # ✏️ Agent 标签 + 颜色
│       └── PlanningCard.tsx          # 🆕 规划进度卡片
```

## 验证流程

```bash
# 1. 启动三端
make all

# 2. AgentEnd 单元验证 (Step 1-7 完成后)
# 闲聊测试
curl -X POST http://localhost:8001/v1/agent/execute \
  -H "Content-Type: application/json" \
  -d '{ "message": "你好", "agent_type": "orchestrator", ... }'
# 预期: 直接文本回复

# 编排测试
curl -X POST http://localhost:8001/v1/agent/execute \
  -H "Content-Type: application/json" \
  -d '{ "message": "写登录页并审查", "agent_type": "orchestrator", ... }'
# 预期: PLANNING → DISPATCH → EXECUTE → REVIEW → AGGREGATE → DONE

# 3. 全链路验证 (Step 8-9 完成后)
# 浏览器打开 → 选 Orchestrator → 测试闲聊 + 任务编排 + 上下文引用
```
