## Context

当前 agentend 通过子进程调用 CLI 工具（Claude CLI / OpenCode CLI）获取 LLM 回复，使用 `--output-format stream-json` / `--format json` 输出。这些格式按**消息粒度**输出——等 LLM 生成完整 assistant 消息后才输出一行 JSON。前端虽然实现了 SSE 逐 chunk 渲染，但收到的 TEXT 事件里已经包含完整文本，等于非流式。

实测发现 Claude CLI 有 `--include-partial-messages` flag，启用后在 `stream-json` 输出中额外输出 `stream_event` 类型事件，包含 `content_block_delta` 逐 token 增量文本。

OpenCode CLI 当前 `--format json` 已经逐行输出 `type: "text"` 事件，但需要实测确认粒度。

## Goals / Non-Goals

**Goals:**
- Claude adapter 实现逐 token 流式输出
- OpenCode adapter 验证并实现逐 token 流式输出
- 改动仅限 agentend adapter 层，前端和后端无需修改

**Non-Goals:**
- 不改 EventType 枚举或 contracts/schemas
- 不改前端 SSE 解析逻辑
- 不改后端 SSE 透传逻辑
- 不引入新的 SDK 依赖（继续走 CLI 子进程方案）

## Decisions

### D1: Claude CLI 使用 `--include-partial-messages` 启用逐 token 输出

CLI 输出两种事件混合：
1. `type: "stream_event"` — 逐 token 增量，`event.type = "content_block_delta"` 时 `event.delta.text` 为文本片段
2. `type: "assistant"` — 完整消息（仍会输出，作为冗余确认）

**Adapter 改动**：
- `_build_command` 末尾追加 `--include-partial-messages`
- `_parse_stream_line` 新增分支：`type == "stream_event"` → 检查 `event.type`
  - `content_block_delta` → yield `StreamEvent(TEXT, text=event.delta.text)`
  - `message_start` / `content_block_start` / `content_block_stop` / `message_stop` → 忽略（元事件）
- `type: "assistant"` 仍然映射为 TEXT，但因为 stream_event 已经输出了全部 token，此时 content_text 会重复。方案：当检测到已有 stream_event 输出时，跳过 `assistant` 事件（或让 `_TYPE_MAP` 不映射 `assistant`，改由 `stream_event` 全权负责文本输出）

**方案选择**：保留 `assistant` 映射不变，但将 `stream_event` 也映射为 TEXT。由于前端是追加式（`streamingContent += text`），重复输出会导致文本重复。因此需要：
- 将 `stream_event.content_block_delta` 作为主要文本来源
- 将 `type: "assistant"` 映射为空事件或忽略（因为完整文本已通过 delta 输出）

### D2: OpenCode CLI 输出粒度验证

实测确认 OpenCode `--format json` 是否已经逐 token 输出 `type: "text"` 事件。如果是，则无需修改 adapter 逻辑；如果不是，调研 `opencode run` 是否有类似 `--stream` flag。

### D3: 不引入去重机制

不在 adapter 层做事件去重。通过正确的类型映射（`content_block_delta` → TEXT, `assistant` → 忽略/空）避免重复文本。

## Risks / Trade-offs

- **[stream_event 格式不稳定]** Claude CLI 是内部工具，`stream_event` 的事件结构可能随版本变化 → 只依赖 `event.type` 和 `event.delta.text` 两个字段，其余忽略
- **[OpenCode 不支持逐 token]** 如果 OpenCode CLI 无法逐 token 输出 → 保持现有行为（消息粒度），仅 Claude adapter 获得改善
- **[CLI 版本兼容]** `--include-partial-messages` 需要较新版本的 Claude CLI → 在 CLI help 中已确认该 flag 存在

## Open Questions

- OpenCode CLI 实测结果：是否已经逐 token 输出？
