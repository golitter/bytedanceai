## MODIFIED Requirements

### Requirement: AdapterRegistry supports ORCHESTRATOR and CODEX agent types
The AdapterRegistry SHALL include OrchestratorAdapter registered under `AgentType.ORCHESTRATOR` (value `"orchestrator"`) and CodexAdapter registered under `AgentType.CODEX` (value `"codex"`). Both SHALL be instantiatable via `registry.get()` just like existing adapters.

#### Scenario: Orchestrator adapter resolves from registry
- **WHEN** `registry.get("orchestrator")` is called
- **THEN** it SHALL return the `OrchestratorAdapter` class without raising ValueError

#### Scenario: Codex adapter resolves from registry
- **WHEN** `registry.get("codex")` is called
- **THEN** it SHALL return the `CodexAdapter` class without raising ValueError

#### Scenario: Codex adapter instantiates
- **WHEN** `CodexAdapter()` is constructed
- **THEN** it SHALL create an empty `_processes` dict and be ready for `chat` or `stream_chat` calls
