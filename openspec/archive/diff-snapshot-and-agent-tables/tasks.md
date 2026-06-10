## 1. Backend Models

- [x] 1.1 Create `backend/internal/model/diff_snapshot.go` with DiffSnapshot struct (id, snapshot_id UNIQUE, session_id, diff_content LONGTEXT, status VARCHAR 16 default pending, created_at, updated_at)
- [x] 1.2 Create `backend/internal/model/session_agent.go` with SessionAgent struct (id, session_id UNIQUE, agent_type, agent_name, avatar_url, created_at, updated_at)
- [x] 1.3 Register both models in `backend/cmd/server/main.go` AutoMigrate call

## 2. Backend Diff Snapshot Handler

- [x] 2.1 Create `backend/internal/handler/diff_snapshot.go` with GetDiffSnapshot handler (GET /api/diff-snapshots/:snapshotId)
- [x] 2.2 Add SaveDiffSnapshot handler (PUT /api/diff-snapshots/:snapshotId) with auto-cancel logic: cancel same-session pending snapshots before upsert
- [x] 2.3 Add terminal state guard: reject PUT on committed/reverted/cancelled snapshots with 409
- [x] 2.4 Register diff-snapshot routes in main.go (GET + PUT /api/diff-snapshots/:snapshotId)

## 3. Backend Session Agent Handler

- [x] 3.1 Modify `backend/internal/handler/task.go`: create session_agents row alongside session creation (same transaction)
- [x] 3.2 Modify `backend/internal/handler/task.go`: RunTask also writes/updates session_agents row
- [x] 3.3 Modify `backend/internal/handler/avatar.go`: UpdateSession writes to session_agents table instead of session table
- [x] 3.4 Modify session listing queries to LEFT JOIN session_agents, returning agent_type/agent_name/avatar_url in response

## 4. Agent Render Skill

- [x] 4.1 Modify `agentend/src/skills/builtin/render/card_diff.go`: import uuid, generate UUID v4, output `snapshotId: {uuid}` in diff block
- [x] 4.2 Ensure go.mod has uuid dependency (or use crypto/rand fallback)

## 5. Frontend Block Types & Parsing

- [x] 5.1 Update `frontend/src/lib/block-types.ts`: change diff variant from `{ type: 'diff' }` to `{ type: 'diff'; snapshotId: string }`
- [x] 5.2 Update `frontend/src/lib/block-reducer.ts`: extract snapshotId field when parsing diff block, return null if missing

## 6. Frontend DiffCard Rewrite

- [x] 6.1 Rewrite `frontend/src/components/cards/DiffCard.tsx`: accept snapshotId prop, implement snapshot-first render logic (GET snapshot → fallback to workspace diff → PUT pending)
- [x] 6.2 Implement settled rendering: show frozen diff with badge (已接受/已拒绝/已取消) + hide action buttons + grey out content
- [x] 6.3 Implement commit/revert flow: POST commit/revert → PUT snapshot with committed/reverted status → update UI

## 7. Frontend MessageBubble Integration

- [x] 7.1 Update `frontend/src/components/chat/MessageBubble.tsx`: BlockRenderer passes snapshotId from diff block to DiffCard

## 8. Cleanup & Verification

- [x] 8.1 Remove session-level diff-snapshot code from earlier attempt (handler/session.go GetDiffSnapshot/SaveDiffSnapshot, main.go diff-snapshot routes under /session/)
- [x] 8.2 Build all three projects (backend go build, frontend pnpm build, agent render skill if applicable)
- [ ] 8.3 Manual test: agent outputs diff → DiffCard creates pending snapshot → accept → page refresh → snapshot restored with committed badge
