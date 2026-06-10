## Why

Group-chat orchestrator answers currently expose internal execution noise instead of a clean collaborative answer. In the `tttt` chat, real browser inspection showed empty agent bubbles, duplicated child-agent output, timeout strings appended to normal prose, and long final reports that obscure the actual outcome.

This change makes group-chat streaming and persisted history present the same coherent structure: delegation cards show assignment/status, child-agent messages show final output, errors are separated from prose, and orchestrator summaries are concise by default with details available on demand.

## What Changes

- Suppress empty persisted agent messages in chat history and live rendering.
- Treat `runtime_text` as transient execution progress, not as durable answer content once official child-agent `text` output exists.
- Render task timeout/error signals as structured failure state instead of appending `[Timeout]` / `[Error]` strings into message prose.
- Keep ask-agent cards as delegation/status summaries and prevent duplicate card/text marker rendering.
- Normalize agent display labels so user-defined names are shown consistently instead of leaking fallback types like `Orchestrator`.
- Present orchestrator final reports as concise outcome cards with expandable details, especially when tasks partially fail.
- Preserve existing SSE contracts unless a spec requires additional metadata for structured task errors.

## Capabilities

### New Capabilities
- `group-chat-result-presentation`: Covers user-facing rendering rules for group-chat delegation cards, child-agent outputs, transient execution progress, task failures, and final orchestrator summaries.

### Modified Capabilities
- `message-rendering`: Long, structured group-chat messages must support concise default presentation, scrolling, expansion, and deduplication of structured blocks.
- `orchestrator-streaming-forward`: Runtime forwarding must distinguish transient progress from durable answer content and expose task failures without polluting normal text.

## Impact

- **Frontend**: chat store block reduction, message hydration, message list filtering, ask-agent/runtime card rendering, long-message summary/expand UI.
- **Agentend**: orchestrator execution and aggregation paths may need to emit/format task failure metadata separately from result prose.
- **Backend**: message listing may filter empty messages or preserve them while frontend hides them; stream persistence may need minor handling for structured runtime/error markers.
- **Contracts**: likely no breaking API change; optional event metadata may be added if needed for task failure blocks.
- **Tests**: add frontend store/rendering tests for duplicate suppression and runtime text lifecycle; add agentend/backend tests if structured error emission changes.
