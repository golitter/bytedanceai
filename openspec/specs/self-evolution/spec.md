## Requirements

### Requirement: Evolution store records orchestration outcomes
After each orchestration cycle, the system SHALL record an entry to `shared/.agent/evolution.yaml` containing: timestamp, user message summary, plan summary, results summary, success boolean, and per-agent performance (agent_id, success, duration).

#### Scenario: Successful orchestration recorded
- **WHEN** a plan→dispatch→aggregate cycle completes successfully
- **THEN** an entry is appended to `evolution.yaml` with `success: true` and agent durations

#### Scenario: Failed orchestration recorded
- **WHEN** an agent fails during execution
- **THEN** an entry is appended with `success: false` and failure details in results_summary

### Requirement: History is capped at 20 entries
The evolution store SHALL only keep the most recent 20 entries. Older entries are discarded on each new write.

#### Scenario: 25th orchestration
- **WHEN** 25 entries exist and a new one is recorded
- **THEN** the oldest 5 entries are removed, keeping only the latest 20

### Requirement: Recent experience is injected into Planner prompt
`EvolutionStore.get_recent_experience(n=5)` SHALL return a formatted string of the most recent N entries for Planner prompt injection.

#### Scenario: Planner receives history
- **WHEN** Planner builds prompt and calls `evolution.get_recent_experience(5)`
- **THEN** returned string includes the 5 most recent entries with success/failure indicators under "## 最近编排经验" heading
