## Requirements

### Requirement: Aggregator summarizes multiple agent results via LLM
Aggregator SHALL accept a list of task results and plan overview, call LLM to generate a human-readable summary, and return the aggregated report string.

#### Scenario: Two-agent result aggregation
- **WHEN** two agents complete their tasks and results are collected
- **THEN** Aggregator calls LLM with both results and overview, returns a unified report string

#### Scenario: Empty results
- **WHEN** no agent results are provided
- **THEN** Aggregator SHALL return an empty string without calling LLM

### Requirement: Aggregator report is included in DONE event
The aggregated report SHALL be included in the final DONE StreamEvent content.

#### Scenario: Stream ends with aggregated report
- **WHEN** OrchestratorAdapter completes the full plan→dispatch→aggregate cycle
- **THEN** the DONE event payload includes the aggregated report text
