## Context

Orchestrator 使用 LangGraph + LangChain ChatOpenAI 进行任务规划和调度。当前 `reason_node` 内部手动调用 `llm_with_tools.ainvoke(messages)` 时没有传入 LangGraph 的 runnable config，导致 LLM 子调用脱离 graph trace 树，LangSmith 无法捕获。

项目已安装 `langsmith==0.8.5`，LangChain/LangGraph 内置支持自动 trace，只需环境变量开启 + config 传播即可。

## Goals / Non-Goals

**Goals:**
- Orchestrator 的全部 LLM 调用（完整 messages 输入、response、tool_calls）自动上报 LangSmith
- LangGraph 的 graph 节点转换（skill_prepare → reason → dispatch → execute → review → evolve → save_mem）自动可见
- 每次 LLM 调用的 token usage 和耗时自动记录
- 不设环境变量时零开销、零行为变化

**Non-Goals:**
- CLI Adapter（Claude Code / OpenCode / Codex）的 trace 接入（Phase 5.2）
- 自研 trace 系统或自定义 tracer
- LangSmith 数据的持久化到自有存储
- 前端 trace 可视化

## Decisions

### Decision 1：使用 LangSmith 而非自研 trace

**选择**：直接使用 LangSmith，通过环境变量开启。

**备选**：自研 `BaseCallbackHandler` tracer，输出到本地 JSON 文件。

**理由**：LangSmith 已是 LangChain 生态的标准 observability 方案，项目已安装依赖。自研 tracer 需要维护数据模型、序列化、UI 查看，投入不成比例。当前阶段目标是快速看到 prompt 以便调优，不是建设 trace 基础设施。

### Decision 2：`get_config()` 传播而非 `callbacks` 参数注入

**选择**：在 `reason_node` 内用 `langchain_core.runnables.get_config()` 获取 LangGraph 传入的 config，传给 `llm_with_tools.ainvoke(messages, config=llm_config)`。

**备选**：在 `OrchestratorAdapter.stream_chat()` 中创建 callback handler 注入 config。

**理由**：`get_config()` 是 LangGraph 官方推荐的节点内 config 获取方式。graph 的 `astream(initial_state, config=config)` 已经包含 callbacks，只需在节点内取出并传递。这种方式不需要修改 adapter 层代码，改动最小（1 处）。

### Decision 3：环境变量名使用新版 `LANGSMITH_TRACING`

**选择**：`LANGSMITH_TRACING=true`

**备选**：`LANGCHAIN_TRACING_V2=true`（旧写法）

**理由**：官方文档推荐新版变量名，语义更清晰。旧写法仍可用但已不推荐。

## Risks / Trade-offs

**[Risk] `get_config()` 在非 LangGraph 上下文中抛 `RuntimeError`** → 用 try/except 兜底，fallback 到 `config=None`（不传 config 时 LLM 正常调用，只是不触发 callbacks）。

**[Risk] LangSmith 为云端服务，数据隐私** → `LANGSMITH_TRACING` 不设就不上报。生产环境可按需开启。

**[Risk] reason_node 的 tool 调用（`tool_fn.invoke(tc["args"])`）不经过 LangChain callback** → 目前 reason_node 里的工具调用是同步 `invoke`，不走 LangChain 的 tool callback 机制。这些不会出现在 LangSmith trace 中，但 LLM 返回的 tool_call（包含工具名和参数）会被记录，对于 prompt 调优已足够。
