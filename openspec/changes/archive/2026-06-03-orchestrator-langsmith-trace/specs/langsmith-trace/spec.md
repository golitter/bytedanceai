## ADDED Requirements

### Requirement: LangSmith tracing enabled by environment variables
The system SHALL read `LANGSMITH_TRACING`, `LANGSMITH_API_KEY`, and `LANGSMITH_PROJECT` environment variables to control LangSmith trace reporting. When `LANGSMITH_TRACING` is not set or set to a falsy value, no trace data SHALL be sent and there SHALL be zero runtime overhead.

#### Scenario: Tracing enabled
- **WHEN** `LANGSMITH_TRACING=true` and `LANGSMITH_API_KEY` is set
- **THEN** all Orchestrator LLM calls, graph node transitions, and tool calls SHALL be reported to LangSmith under the configured project

#### Scenario: Tracing disabled
- **WHEN** `LANGSMITH_TRACING` is not set or empty
- **THEN** no trace data SHALL be sent and Orchestrator behavior SHALL be identical to before this change

### Requirement: LLM config propagation in reason_node
The `reason_node` function in `graph.py` SHALL propagate the LangGraph runnable config to `llm_with_tools.ainvoke()` so that LangSmith callbacks are triggered for every LLM call within the tool-calling loop.

#### Scenario: Config propagation succeeds
- **WHEN** `reason_node` is executed within a LangGraph run context
- **THEN** `get_config()` SHALL return the current runnable config and it SHALL be passed to `llm_with_tools.ainvoke(messages, config=llm_config)`

#### Scenario: Config propagation falls back gracefully
- **WHEN** `get_config()` raises `RuntimeError` (not in a LangGraph context)
- **THEN** `llm_config` SHALL be set to `None` and `llm_with_tools.ainvoke(messages, config=None)` SHALL execute normally without trace callbacks

### Requirement: Complete LLM call visibility
Each LLM call in the reason_node tool-calling loop SHALL be visible in LangSmith with: the full input messages array (system prompt + memory + human message), the response (text or tool_calls), token usage, and latency.

#### Scenario: Text response captured
- **WHEN** the LLM returns a text response (no tool calls)
- **THEN** LangSmith SHALL show the full input messages and the text response content

#### Scenario: Tool call response captured
- **WHEN** the LLM returns tool calls
- **THEN** LangSmith SHALL show the full input messages and the tool calls with their names and arguments
