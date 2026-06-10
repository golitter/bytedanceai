## Context

Current SSE streaming architecture routes every token event through Redis Stream:

```
AgentEnd → StreamWriter → Redis XADD → StreamHandler XRead (200ms block) → SSE → Frontend
```

StreamWriter batches TEXT events (500ms / 2KB) before publishing to Redis. StreamHandler polls Redis with 200ms blocking XRead. Combined worst-case latency is ~700ms. This feels sluggish for chat UI where users expect typewriter-style token delivery.

The root cause is architectural: Redis Stream is a durable event log used as realtime transport. Token streams are transient high-frequency data that need immediate flush, not durable queue semantics. Final messages (tool calls, artifacts, complete responses) are the data that needs durability.

## Goals / Non-Goals

**Goals:**
- Reduce SSE token delivery latency from ~250-700ms to ~1-10ms
- Decouple realtime transport from durable replay without removing Redis
- Maintain all existing guarantees: replay, reconnect, history, agent switching
- Zero frontend changes — SSE event format stays identical

**Non-Goals:**
- Replacing SSE with WebSocket (HTTP is not the bottleneck)
- Removing Redis Stream entirely (still needed for replay/fanout)
- Multi-instance distributed pub/sub (single backend instance for now)
- Changing agentend or frontend code

## Decisions

### Decision 1: In-memory Go channel hub, not Redis Pub/Sub

**Choice**: `RuntimeHub` with `map[string]*RuntimeStream` where each stream has subscriber channels.

**Alternatives considered**:
- Redis Pub/Sub: Adds network hop, no replay, no good for single-instance case
- Ring buffer: More complex, overkill for current scale
- Direct callback: No fanout support, harder lifecycle management

**Rationale**: Go channels give ~1μs publish latency, native fanout via per-subscriber channels, and simple lifecycle. No external dependency.

### Decision 2: Hub owns monotonic sequence IDs

**Choice**: Each `RuntimeStream` has an `atomic.Uint64` counter. Every published event gets a monotonically increasing seq. This seq is separate from Redis stream IDs.

**Rationale**: Enables correct catch-up + live tail handoff:
1. Subscriber calls `Subscribe()`, gets `<-chan Event` + `currentSeq`
2. Replay from Redis covers `(lastKnownSeq, currentSeq]`
3. Live events start from `seq > currentSeq`
4. No gap, no duplication when sequence spaces are separate

### Decision 3: Dual-write in StreamWriter — hub first, Redis async

**Choice**: StreamWriter calls `hub.Publish()` synchronously (immediate SSE delivery), then continues existing Redis XADD + MySQL flush asynchronously.

**Rationale**: Zero risk approach. Redis persistence path is completely untouched. Hub is additive. If hub fails, Redis still works. Gradual migration possible.

### Decision 4: SSE handler three-phase handshake

**Choice**: `serveStreaming()` restructured to:
1. Phase 1: MySQL history (unchanged)
2. Phase 2: Subscribe hub → get `currentSeq` → replay Redis gap `(lastSeq, currentSeq]`
3. Phase 3: Consume hub channel for realtime events

**Alternatives considered**:
- Hub-only (no Redis replay): Loses gap-fill on reconnect, risky
- Subscribe-then-replay (naive): Race condition — events published during replay are missed or duplicated

**Rationale**: Three-phase with `currentSeq` as the handoff point is the classic "catch-up + live tail" pattern from event streaming systems. Correct by construction.

### Decision 5: Hub lifecycle tied to StreamWriter

**Choice**: StreamWriter creates the hub stream on first publish and closes it on `finish()`. Subscribers receive a `Done` event when the stream closes.

**Rationale**: StreamWriter already owns the stream lifecycle. Hub stream is a view into that lifecycle. No separate coordination needed.

## Risks / Trade-offs

**[Memory usage]** Hub holds subscriber channels in memory. For N concurrent streams, N channels exist. → Mitigation: Each stream has bounded channel buffer (256 events). Stale streams are closed by StreamWriter timeout (30 min). Single-instance only — no cross-process coordination.

**[Hub publish failure is silent]** If no subscribers exist, `hub.Publish()` is a no-op. Events are still persisted to Redis/MySQL. → Mitigation: This is actually correct behavior — no frontend connected means no delivery needed. Reconnect picks up from Redis.

**[Sequence gap between hub seq and Redis seq]** Hub seq and Redis stream IDs are independent. Replay logic must map between them. → Mitigation: Replay uses Redis XRead range `(lastRedisID, lastRedisIDBeforeSubscribe]` which is well-defined. Hub seq is only used for the handoff point, not for replay indexing.

**[Migration risk]** Dual-write means events flow through both paths simultaneously. → Mitigation: Feature is additive. Redis path unchanged. Hub can be disabled by simply not calling `hub.Publish()`. Rollback = remove hub calls.
