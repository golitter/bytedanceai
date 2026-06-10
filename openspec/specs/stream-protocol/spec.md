## MODIFIED Requirements

### Requirement: 流式调用使用异步 LLM 接口
Agent 端所有在 async 上下文中执行的 LLM 调用 SHALL 使用 `await llm.ainvoke()` 而非阻塞的 `llm.invoke()`，避免阻塞事件循环。受影响模块：`orchestrator/pin_memory.py`（`_generate_summary`）、`orchestrator/aggregator.py`（`aggregate`）。

#### Scenario: async 上下文中不使用同步 LLM 调用
- **WHEN** 在 async def 函数中需要调用 LLM
- **THEN** 使用 `await llm.ainvoke()` 而非 `llm.invoke()`

#### Scenario: 改动后调用链完整可用
- **WHEN** `_generate_summary` 和 `aggregate` 改为 async
- **THEN** 所有上游调用方使用 `await` 调用，无遗漏
