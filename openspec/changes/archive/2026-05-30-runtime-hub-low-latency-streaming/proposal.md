## Why

SSE token streaming latency is 250-700ms because token events flow through two batching layers — StreamWriter buffers TEXT events 500ms/2KB before Redis XADD, then StreamHandler polls Redis with XRead Block=200ms. This "durable-first" architecture (Redis Stream on the critical path) is designed for high-throughput event logging, not ultra-low-latency token transport. Chat UI needs a "realtime-first" architecture where token arrival triggers immediate flush.

## What Changes

- Introduce an in-memory `RuntimeHub` that publishes token events to SSE handlers via Go channels, bypassing Redis on the realtime path
- StreamWriter becomes a dual-writer: `hub.Publish()` (immediate, ~1ms) + async Redis XADD (background, latency irrelevant)
- StreamHandler SSE endpoint switches from Redis XRead polling to `hub.Subscribe()` for live events, with Redis retained only for replay/reconnect gap-fill
- Add monotonic sequence IDs owned by RuntimeHub (not Redis) to enable correct catch-up + live tail handoff
- Redis Stream and MySQL persistence are untouched — zero data loss risk

## Capabilities

### New Capabilities
- `runtime-hub`: In-memory pub/sub hub for low-latency token streaming. Owns sequence IDs, manages subscriber channels per stream key, supports subscribe/publish/close lifecycle. Decouples realtime transport from durable replay.

### Modified Capabilities
- `chat-streaming`: SSE handler adds hub.Subscribe() path for realtime events alongside existing MySQL history + Redis replay. Realtime events no longer sourced from XRead polling.
- `stream-incremental-persist`: StreamWriter adds hub.Publish() call on the hot path before async Redis XADD. TEXT event batching on the Redis path remains for efficiency but no longer blocks SSE delivery.

## Impact

- **Backend `internal/stream/`**: StreamWriter gains hub dependency; new `hub.go` file for RuntimeHub
- **Backend `internal/handler/stream.go`**: `serveStreaming()` Phase 2 switches from XRead loop to hub.Subscribe
- **Backend `internal/handler/task.go`**: StreamWriter creation passes hub instance
- **Frontend**: No changes — SSE event format is identical, latency improvement is transparent
- **Redis**: Still used for replay/reconnect, just not on realtime critical path
- **MySQL**: Unchanged
