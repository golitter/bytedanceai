## ADDED Requirements

### Requirement: Design documents must be cross-validated against codebase

All design documents in `discussions/` SHALL be validated against actual source code in `frontend/`, `backend/`, `agentend/`, and `contracts/` before implementation begins.

#### Scenario: Event type names match across documents
- **WHEN** two documents reference the same event type (e.g., workspace branch creation)
- **THEN** the event name SHALL be identical (not `workspace.branch_created` in one and `workspace.branch.created` in another)

#### Scenario: Data models match existing code
- **WHEN** a document proposes extending an existing model (e.g., AgentRequest)
- **THEN** the existing fields SHALL match what's currently in `contracts/schemas/` or `src/schemas/`

### Requirement: Design documents must not contradict each other

When multiple documents describe the same feature, they SHALL use consistent terminology and design decisions.

#### Scenario: Go Backend involvement is consistently stated
- **WHEN** one document says "Go Backend 无需改动" and another lists Go Backend changes
- **THEN** the discrepancy SHALL be resolved before implementation

### Requirement: Design documents must follow contract-first principle

Any new protocol fields or event types SHALL reference or plan for `contracts/schemas/` YAML updates.

#### Scenario: New event types require contract updates
- **WHEN** 20+ new RuntimeEvent types are proposed
- **THEN** a plan for updating `contracts/schemas/event-types.yaml` SHALL be included
