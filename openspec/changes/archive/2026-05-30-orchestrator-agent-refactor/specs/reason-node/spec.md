## ADDED Requirements

### Requirement: REASON node performs LLM tool-calling loop with skill discovery

The REASON node SHALL execute a tool-calling loop that integrates progressive skill discovery (L1 scan → L2 semantic selection → L3 content load) as internal function calls before entering the LLM loop. The node SHALL inject pin constraints and evolution history into the system prompt.

#### Scenario: Skill discovery integrated into REASON
- **WHEN** REASON node starts
- **THEN** it performs discover_skills → select_skills → load_l2_content as function calls, then enters LLM tool-calling loop with skill content injected

#### Scenario: Pin and evolution context injected
- **WHEN** REASON node builds system prompt
- **THEN** pin constraints from PinMemory.get_context() and evolution history from EvolutionStore.get_recent_experience() are included

### Requirement: REASON node distinguishes text output from plan_and_dispatch tool call

The REASON node SHALL determine output type based on LLM response: if no tool_calls, output is text (chitchat/direct answer); if tool_calls contains `plan_and_dispatch`, output is plan (orchestration).

#### Scenario: Chitchat produces text output
- **WHEN** user sends "你好" and LLM responds with text and no tool_calls
- **THEN** REASON returns `{"output_type": "text", "text": "<response>"}`

#### Scenario: Task triggers plan_and_dispatch
- **WHEN** user sends "分析项目架构" and LLM calls `plan_and_dispatch(overview="...", tasks=[...])`
- **THEN** REASON returns `{"output_type": "plan", "plan": PlanOutput(overview=..., tasks=[...])}`

#### Scenario: LLM uses info-gathering tools before deciding
- **WHEN** user sends a complex request and LLM calls read_file to gather info
- **THEN** REASON executes the tool, appends ToolMessage, and continues the loop

#### Scenario: LLM reaches max iterations
- **WHEN** LLM tool-calling loop reaches MAX_ITERATIONS without plan_or_text output
- **THEN** REASON returns `{"output_type": "text", "text": "规划超时，请重新描述需求"}`

### Requirement: plan_and_dispatch tool signals orchestration intent

The system SHALL define a `plan_and_dispatch` tool with parameters `overview: str` and `tasks: list[dict]`. The tool's return value is intercepted by the REASON loop to construct a PlanOutput — the tool itself does not execute any orchestration.

#### Scenario: plan_and_dispatch tool is called
- **WHEN** LLM generates a tool_call with name "plan_and_dispatch" and args `{overview: "...", tasks: [...]}`
- **THEN** the REASON loop extracts the args, constructs PlanOutput, and returns with output_type="plan"
