## ADDED Requirements

### Requirement: Dispatcher produces execution waves via topological sort

The Dispatcher SHALL topologically sort DispatchResults based on their depends_on field, producing execution_waves — a list of lists where each inner list contains tasks that can execute in parallel.

#### Scenario: No dependencies — single wave
- **WHEN** all tasks have empty depends_on
- **THEN** execution_waves = [[task-001, task-002, task-003]]

#### Scenario: Chain dependencies — sequential waves
- **WHEN** task-002 depends on task-001, task-003 depends on task-002
- **THEN** execution_waves = [[task-001], [task-002], [task-003]]

#### Scenario: Mixed dependencies
- **WHEN** task-002 and task-003 depend on task-001, no other dependencies
- **THEN** execution_waves = [[task-001], [task-002, task-003]]
