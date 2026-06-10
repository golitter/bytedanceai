## 1. Theme & CSS Foundation

- [x] 1.1 Rewrite `index.css` `.dark` block: map shadcn CSS variables to visual-style-guide color values (background → `#0A0B0E`, card → `#1A1D24`, foreground → `#E8EBF0`, muted-foreground → `#8B91A0`, primary → `#6366F1`, border → `rgba(255,255,255,0.06)`, input → `rgba(255,255,255,0.06)`, ring → `#6366F1`)
- [x] 1.2 Add custom CSS variables in `.dark`: `--bg-canvas`, `--bg-sidebar`, `--bg-hover`, `--bg-active`, `--text-primary`, `--text-secondary`, `--text-tertiary`, `--color-brand`, `--color-success`, `--color-warning`, `--color-error`, `--agent-claude`, `--agent-opencode`, `--agent-orchestrator`
- [x] 1.3 Add `.dark` class to `<html>` in `index.html` as default, keep `:root` light values as placeholder

## 2. Dependencies & Config

- [x] 2.1 Install new dependencies: `react-markdown`, `remark-gfm`, `shiki`, `@tanstack/react-virtual`
- [x] 2.2 Add Vite proxy for `/api` → `http://localhost:8080` in `vite.config.ts`
- [x] 2.3 Clean up demo code: remove counter from `stores/app.ts`, remove demo imports

## 3. API Layer

- [x] 3.1 Rewrite `lib/api.ts`: define types for Task, Session, AgentType, StreamEvent; implement `fetchTasks()`, `fetchTask(id)`, `createTask(title)`, `deleteTask(id)`, `fetchAgentTypes()`
- [x] 3.2 Create `lib/sse.ts`: SSE connection with fetch + ReadableStream + TextDecoder, cross-chunk `\n\n` boundary buffering, parse `data:` JSON lines, return AbortController
- [x] 3.3 Create `hooks/use-chat-stream.ts`: useReducer with discriminated union ChatState (idle/loading/streaming/tool_running/done/error), wire SSE events to dispatch, manage message accumulation

## 4. Stores

- [x] 4.1 Create `stores/chat.ts`: Zustand store — `currentSessionId`, `currentTaskId`, navigation actions (`setCurrentSession`, `setCurrentTask`). Only navigation state, NOT messages or streaming state.
- [x] 4.2 Create `hooks/use-sessions.ts`: TanStack Query hooks — `useSessions(taskId)`, `useCreateSession()`, `useDeleteSession()` with `invalidateQueries` on mutation success

## 5. Chat Components — Dumb Layer

- [x] 5.1 Create `components/chat/AgentAvatar.tsx`: rounded-square avatar (8px radius, 32px), agent identity color background, 4px status dot (green/yellow/gray/red with appropriate animations)
- [x] 5.2 Create `components/chat/MessageBubble.tsx`: three variants — user (right-aligned, brand tint bg + border), agent (left-aligned, card bg, 3px identity color bar), system (centered, muted text, no bubble). Content rendered as children.
- [x] 5.3 Create `components/chat/MessageInput.tsx`: auto-expanding textarea (min 48px, max 200px), Enter send / Shift+Enter newline, send button, disabled prop for streaming state

## 6. Markdown & Code Rendering

- [x] 6.1 Create `components/markdown/MarkdownRenderer.tsx`: react-markdown + remark-gfm, override `code`/`pre` for code blocks, dark styling for tables/inline code, Geist Mono for inline code with `#0D0F14` bg
- [x] 6.2 Create `components/markdown/CodeBlock.tsx`: Shiki syntax highlighting (Tokyo Night), dynamic import with plain-text fallback, line numbers (`#5A6070`), Geist Mono 13px, `#0D0F14` bg, 8px border-radius

## 7. Chat Components — Smart Layer

- [x] 7.1 Create `components/chat/MessageList.tsx`: render MessageBubble list with MarkdownRenderer content, @tanstack/react-virtual for >50 messages, auto-scroll with "scroll to bottom" button, respect manual scroll position
- [x] 7.2 Create `components/chat/ChatArea.tsx`: compose header (task title) + MessageList + MessageInput, wire `useChatStream` hook for streaming state, empty state welcome message, streaming cursor on active agent message
- [x] 7.3 Create `components/chat/ChatSidebar.tsx`: app logo, "New Chat" button, session list with active highlight (2px brand border + `#22262F` bg), empty state prompt, use `useSessions` hook

## 8. Page & Route Integration

- [x] 8.1 Create `pages/ChatPage.tsx`: three-column flex layout (ChatSidebar 260px + ChatArea flex-1 + hidden reserved slot), load task data via `useQuery` with taskId from URL params
- [x] 8.2 Rewrite `pages/TaskList.tsx`: dark-themed task list with "New Task" button (create task via mutation), each task card links to `/tasks/:taskId`
- [x] 8.3 Update `main.tsx` routes: `/` redirect to `/tasks`, `/tasks` → TaskList, `/tasks/:taskId` → ChatPage
- [x] 8.4 Update `App.tsx`: remove demo content, render `<Routes>` only

## 9. Polish & Verification

- [x] 9.1 Verify visual adherence: background hierarchy (5-level gray), border opacity `rgba(255,255,255,0.06)`, text color levels, agent identity colors, 4px spacing base
- [x] 9.2 Verify SSE streaming end-to-end: send message → text events stream → cursor blinks → done → message finalized → input re-enabled
- [x] 9.3 Verify code block highlighting: trigger code response → Tokyo Night syntax colors, line numbers, Geist Mono font, `#0D0F14` background
- [x] 9.4 Verify virtual list: load >50 messages → smooth scrolling, dynamic height measurement, no performance degradation
