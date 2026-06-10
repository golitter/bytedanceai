## MODIFIED Requirements

### Requirement: REASON node performs LLM tool-calling loop with skill discovery

The REASON node SHALL execute a tool-calling loop that integrates progressive skill discovery (L1 scan → L2 semantic selection → L3 content load) as internal function calls before entering the LLM loop. The node SHALL inject pin constraints and evolution history into the system prompt. The REASON node SHALL propagate the LangGraph runnable config to `llm_with_tools.ainvoke()` so that tracing callbacks (e.g., LangSmith) are triggered for every LLM call within the loop.

#### Scenario: Skill discovery integrated into REASON
- **WHEN** REASON node starts
- **THEN** it performs discover_skills → select_skills → load_l2_content as function calls, then enters LLM tool-calling loop with skill content injected

#### Scenario: Pin and evolution context injected
- **WHEN** REASON node builds system prompt
- **THEN** pin constraints from PinMemory.get_context() and evolution history from EvolutionStore.get_recent_experience() are included

#### Scenario: LangGraph config propagated to LLM calls
- **WHEN** REASON node calls `llm_with_tools.ainvoke(messages)`
- **THEN** it SHALL retrieve the current LangGraph runnable config via `get_config()` and pass it as the `config` parameter to `ainvoke()`

#### Scenario: Config retrieval fails gracefully
- **WHEN** `get_config()` raises `RuntimeError` (outside LangGraph context)
- **THEN** REASON node SHALL set config to `None` and proceed with `llm_with_tools.ainvoke(messages, config=None)`
