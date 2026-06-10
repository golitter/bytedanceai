## MODIFIED Requirements

### Requirement: Planner generates plan via tool-calling agent loop

The planner SHALL use an LLM tool-calling loop (max 10 iterations) where the LLM can call info-gathering tools (read_file, write_file, list_dir, run_skill, load_resource) and the plan_and_dispatch tool. The loop terminates when the LLM responds without tool_calls (text output) or calls plan_and_dispatch (plan output).

#### Scenario: LLM gathers info then outputs plan
- **WHEN** LLM calls read_file to inspect code, then calls plan_and_dispatch with tasks
- **THEN** the loop collects tool results, then returns PlanOutput from plan_and_dispatch args

#### Scenario: LLM responds with text directly
- **WHEN** LLM responds to a chitchat message with text and no tool_calls
- **THEN** the loop terminates immediately, returning the text as output_type="text"

## ADDED Requirements

### Requirement: System prompt supports chitchat and orchestration modes

The system prompt SHALL instruct the LLM that it is a conversational agent capable of both answering questions directly and orchestrating multi-agent tasks. When the user's request requires multi-agent collaboration, the LLM SHALL call the plan_and_dispatch tool. When the request can be answered directly, the LLM SHALL respond with text.

#### Scenario: Prompt enables dual-mode behavior
- **WHEN** REASON node constructs the system prompt
- **THEN** the prompt includes instructions for both direct reply and orchestration, available tools list, agents description, and skills content
