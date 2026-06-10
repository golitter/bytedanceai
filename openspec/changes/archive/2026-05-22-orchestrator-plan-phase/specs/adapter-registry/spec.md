## MODIFIED Requirements

### Requirement: AdapterRegistry supports ORCHESTRATOR agent type
The AdapterRegistry SHALL include OrchestratorAdapter registered under `AgentType.ORCHESTRATOR` (value `"orchestrator"`). It SHALL be instantiatable via `registry.get("orchestrator")` just like existing adapters.

#### Scenario: Orchestrator adapter resolves from registry
- **WHEN** `registry.get("orchestrator")` is called
- **THEN** it SHALL return the `OrchestratorAdapter` class without raising ValueError

#### Scenario: Orchestrator adapter instantiates
- **WHEN** `OrchestratorAdapter()` is constructed
- **THEN** it SHALL create a compiled LangGraph graph internally and be ready for `chat` or `stream_chat` calls
