## MODIFIED Requirements

### Requirement: AdapterRegistry supports ORCHESTRATOR and CODEX agent types
The AdapterRegistry SHALL include OrchestratorAdapter registered under `AgentType.ORCHESTRATOR` and CodexAdapter registered under `AgentType.CODEX`. It SHALL also include SkillRule in the rule engine so that output skill prompts are injected for all agent types.

#### Scenario: SkillRule injected for all agent types
- **WHEN** any agent type processes a request through the rule engine
- **THEN** SkillRule SHALL append the output skill prompt to the system_prompt_append
