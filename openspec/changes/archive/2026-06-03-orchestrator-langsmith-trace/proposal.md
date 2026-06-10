## Why

调优 Orchestrator 提示词时，需要看到每次 LLM 调用的完整输入/输出（system prompt、message history、tool calls、token 用量）。目前没有任何 trace 机制，调试全靠日志猜测，没有 trace 基本就是盲调。项目已安装 `langsmith==0.8.5`，接入成本极低。

## What Changes

- 在 `reason_node` 中通过 `get_config()` 将 LangGraph 的 runnable config 传播到内部 `ChatOpenAI.ainvoke()` 调用，使 LangSmith 回调能穿透到 LLM 调用层
- 在 `.env` 中配置 LangSmith 环境变量（`LANGSMITH_TRACING`、`LANGSMITH_API_KEY`、`LANGSMITH_PROJECT`），启用自动 trace
- Orchestrator 的 graph 节点转换、LLM 调用（完整 messages/response）、tool calls、token usage 将自动上报 LangSmith

## Capabilities

### New Capabilities
- `langsmith-trace`: Orchestrator LLM 调用链的 LangSmith 自动 trace 接入，覆盖 graph nodes、ChatOpenAI 输入/输出、tool calls

### Modified Capabilities
- `reason-node`: 修改 LLM 调用方式，从无 config 调用改为传播 LangGraph config

## Impact

- **代码**：`agentend/src/orchestrator/planning/graph.py`（reason_node 函数，1 处改动）
- **配置**：`agentend/.env`（新增 3 行环境变量）
- **依赖**：无新增依赖（`langsmith` 已安装）
- **风险**：`get_config()` 在非 LangGraph 上下文中会抛 `RuntimeError`，已用 try/except 兜底；不设环境变量时零开销
