## 1. Frontend Rendering Normalization

- [x] 1.1 Hide agent/system messages that have no text and no structured blocks.
- [x] 1.2 Coalesce duplicate ask-agent blocks by `question_id` before rendering live or persisted messages.
- [x] 1.3 Coalesce runtime status blocks by `task_id`, preserving only transient progress that should still be visible.
- [x] 1.4 Render ask-agent cards as assignment/status summaries and avoid duplicating full child-agent answers when separate child messages exist.
- [x] 1.5 Resolve agent display labels from session metadata first, with agent type labels only as fallback.
- [x] 1.6 Suppress redundant outer `@agent` labels for self-identifying structured cards.

## 2. Streaming And History Convergence

- [x] 2.1 Treat `runtime_text` as live progress and clear it when durable child-agent `text` arrives for the same logical task/message.
- [x] 2.2 Prevent transient runtime transcripts from being frozen into persisted visible message blocks during agent switches.
- [x] 2.3 Ensure live stream output and refreshed history render the same visible structure for completed group-chat tasks.
- [x] 2.4 Keep long collapsed message cards internally scrollable and clickable/controllable for expanded reading.

## 3. Failure And Summary Presentation

- [x] 3.1 Parse legacy `[Timeout]` markers into structured task failure blocks where safe.
- [x] 3.2 Parse legacy `[Error]` markers into structured error/failure blocks where safe.
- [x] 3.3 Render unparseable marker-like text as plain text so no information is lost.
- [x] 3.4 Add a summary-first final orchestrator report UI for overall status, completed count, failed count, and next action.
- [x] 3.5 Provide expandable details for long Markdown reports, per-task evidence, and logs.

## 4. Orchestrator / Agentend Output

- [x] 4.1 Distinguish forwarded child-agent runtime progress from durable child-agent answer content in orchestrator stream handling.
- [x] 4.2 Emit task timeout metadata including task ID, agent, and reason instead of relying only on marker strings in prose.
- [x] 4.3 Emit task error metadata including task ID, agent, and message instead of relying only on marker strings in prose.
- [x] 4.4 Generate final aggregation output as concise summary plus structured or expandable details.
- [x] 4.5 Avoid copying full child-agent execution logs wholesale into the default final summary.

## 5. Verification

- [x] 5.1 Add frontend store/reducer tests for duplicate ask-agent cards, runtime text clearing, empty-message hiding, and legacy timeout/error parsing.
- [x] 5.2 Add agentend/orchestrator tests for runtime progress vs durable result separation and structured failure metadata.
- [x] 5.3 Run frontend typecheck, lint, and targeted tests.
- [x] 5.4 Run agentend targeted tests for orchestrator streaming/aggregation.
- [x] 5.5 Use the browser to replay a group-chat prompt such as “帮我的前端项目的风格改成铁血风” and verify live output, collapsed card scrolling, expanded reading, labels, failures, and reload consistency.
