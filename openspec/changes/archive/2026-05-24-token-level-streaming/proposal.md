## Why

Agent 回复不是逐 token 流式推送到前端的。Claude CLI 默认的 `stream-json` 输出按消息粒度（等 LLM 生成完整 assistant 消息才输出一行 JSON），导致用户要等数秒甚至十几秒才看到内容，体验上等于非流式。

## What Changes

- Claude CLI 命令新增 `--include-partial-messages` flag，启用 `stream_event` → `content_block_delta` 逐 token 输出
- `_parse_stream_line` 新增 `stream_event` 类型解析，从 `event.delta.text` 提取增量文本
- `_TYPE_MAP` 新增 `"stream_event"` 映射（作为 TEXT 事件 yield）
- OpenCode CLI 实测 `--format json` 的输出粒度，若非逐 token 则寻找替代方案（如 `--stream` flag 或 NDJSON streaming 模式）
- 前端无需改动（已支持逐 chunk `STREAM_TEXT` dispatch）

## Capabilities

### New Capabilities

- `token-streaming`: 逐 token 流式输出能力——adapter 从 CLI stdout 逐行解析增量文本事件，每收到一个 content_block_delta 立即 yield StreamEvent(TEXT)

### Modified Capabilities

- `stream-protocol`: 新增 Claude CLI `stream_event` / `content_block_delta` 到 EventType.TEXT 的映射规则
- `opencode-adapter`: stream_chat 从"模拟流式"改为真正的逐 token 流式（取决于 OpenCode CLI 能力）

## Impact

- `agentend/src/adapters/claude.py` — 命令构建 + 解析逻辑
- `agentend/src/adapters/opencode.py` — 命令构建 + 解析逻辑（视实测结果）
- `contracts/schemas/` — 无变更（EventType 枚举不变，只是产生更多 TEXT 事件）
- 前端 — 无变更
- 后端 — 无变更（透传 SSE 行为不变）
