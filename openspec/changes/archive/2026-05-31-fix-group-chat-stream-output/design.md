## Context

Group-chat execution currently mixes three different concepts in the same visible message stream:

- delegation state (`ask_card_start` / `ask_card_done`)
- transient execution progress (`runtime_executing` / `runtime_text` / `runtime_completed`)
- durable conversation answers (`text` events persisted as agent messages)

Browser inspection of the `tttt` chat showed the result of this mixing: empty orchestrator bubbles, child-agent answers repeated in both runtime cards and standalone messages, timeout markers appended to normal prose, and final summaries that read like internal audit logs. The backend and agentend already persist enough metadata to improve presentation, but the current frontend reducer and orchestrator formatting do not maintain a clean boundary between transient progress, durable output, and summary.

## Goals / Non-Goals

**Goals:**

- Make live streaming and refreshed history converge to the same visible structure.
- Keep runtime progress useful during execution without freezing it into final answer cards when durable child-agent text exists.
- Separate task failures/timeouts from normal answer prose.
- Make group-chat final summaries concise by default while retaining detailed evidence behind expansion.
- Preserve agent identity using user-defined names where available.
- Add tests that reproduce the observed duplicate/timeout/empty-message cases.

**Non-Goals:**

- Redesign the entire chat visual system.
- Change how agents execute tasks or how task scheduling chooses agents.
- Remove detailed child-agent logs entirely; they remain available behind expansion/details.
- Introduce a new external dependency.
- Migrate or rewrite historical database rows. The implementation may normalize rendering of existing rows at read/render time.

## Decisions

### Decision: Treat `runtime_text` as transient progress

`runtime_text` SHALL be displayed only as live task progress while no durable child-agent `text` output is available for that same task/message. Once official agent text begins, the frontend reducer clears the runtime transcript and lets the durable agent message own the answer.

Alternative considered: stop emitting `runtime_text` from agentend. This would remove useful live progress and make long child-agent runs feel frozen, so the change keeps `runtime_text` but scopes it to progress UI.

### Decision: Normalize at the presentation boundary first

The frontend SHALL hide empty messages, coalesce duplicate structured blocks, and normalize agent names. This makes existing persisted rows render better without requiring a database migration. Agentend/backend can still be improved to stop creating noisy rows in future tasks.

Alternative considered: only fix backend persistence. That would not repair historical rows and would leave reload/live divergence until all old data is cleaned.

### Decision: Keep ask-agent cards as status summaries

Ask-agent cards SHALL show source, target, question summary, and status. They SHALL NOT duplicate the full child-agent answer when a separate child-agent message exists. The card can keep a concise answer summary for collapsed state and linking/context.

Alternative considered: inline the full child-agent answer inside the ask-card. This works for simple Q&A, but in group-chat coding tasks it makes the orchestrator card too large and duplicates child-agent messages.

### Decision: Represent task failures structurally

Timeouts and errors SHALL be parsed or emitted as task failure blocks/cards instead of appended to Markdown prose. Final orchestrator summaries can mention failure, but raw `[Timeout]` or `[Error]` markers should not appear as accidental sentence suffixes.

Alternative considered: keep raw markers and style them with regex in Markdown. That is brittle and still leaves the content model polluted.

### Decision: Summary-first final reports

Orchestrator final messages SHALL default to a concise outcome card: overall status, completed tasks, failed tasks, and next action. Detailed per-task logs or long tables SHALL be available through expansion.

Alternative considered: continue rendering full Markdown reports inside a fixed-height card. This is technically simple, but the default chat view remains hard to scan.

## Risks / Trade-offs

- [Risk] Hiding empty messages could hide intentional blank separator messages. → Mitigation: hide only agent/system messages with empty content and no structured blocks.
- [Risk] Clearing runtime transcripts too aggressively could remove useful progress if no durable text arrives. → Mitigation: clear only when non-empty durable `text` arrives; retain runtime failure/completion status.
- [Risk] Historical rows with raw timeout text may still render imperfectly. → Mitigation: add legacy parsing in block reduction for `[Timeout]` / `[Error]` markers and render them as structured failure blocks.
- [Risk] Final summary extraction may misclassify arbitrary Markdown. → Mitigation: prefer structured agentend emission for new runs and use conservative frontend fallback for old rows.
- [Risk] More structured rendering increases reducer complexity. → Mitigation: keep reducer helpers small and cover each observed case with store/block tests.

## Migration Plan

1. Implement frontend normalization and rendering changes behind existing data structures.
2. Add parsing of legacy timeout/error markers so old rows render as failure blocks.
3. Adjust agentend orchestrator aggregation to produce concise summary content and structured detail metadata where feasible.
4. Verify live run and reload show equivalent structure for a group-chat coding task.
5. Rollback by reverting frontend rendering/reducer changes; backend and agentend changes should remain compatible with existing text events.
