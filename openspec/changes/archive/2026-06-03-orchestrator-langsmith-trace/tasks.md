## 1. 环境变量配置

- [x] 1.1 在 `agentend/.env` 中添加 `LANGSMITH_TRACING`、`LANGSMITH_API_KEY`、`LANGSMITH_PROJECT` 三个环境变量（API key 留空占位）
- [x] 1.2 确认 `langsmith` 包已在 `agentend` 环境中可用（已安装 `0.8.5`）

## 2. reason_node config 传播

- [x] 2.1 在 `agentend/src/orchestrator/planning/graph.py` 顶部添加 `from langchain_core.runnables import get_config` import
- [x] 2.2 在 `reason_node` 函数内（LLM 调用循环前）添加 `get_config()` 调用，用 try/except RuntimeError 兜底
- [x] 2.3 将 `llm_with_tools.ainvoke(messages)` 改为 `llm_with_tools.ainvoke(messages, config=llm_config)`

## 3. 验证

- [ ] 3.1 配置有效的 `LANGSMITH_API_KEY`，运行一个 orchestrator 任务，在 LangSmith UI 中确认能看到 graph nodes + LLM calls + tool calls
- [ ] 3.2 不设 `LANGSMITH_TRACING`，运行同样的任务，确认行为无变化、无报错
