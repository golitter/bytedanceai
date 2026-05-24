# Claude Adapter 逐 token 流式输出

## 变更原因

Claude CLI 默认 `stream-json` 输出按消息粒度（等 LLM 生成完整 assistant 消息才输出一行 JSON），导致前端要等数秒才看到内容。启用 `--include-partial-messages` 后，CLI 额外输出 `stream_event` / `content_block_delta` 逐 token 增量文本事件。

## 变更文件

无 schema 文件变更（`contracts/schemas/*.yaml` 未修改）。EventType 枚举不变，只是 adapter 产生更多 TEXT 事件。

## 对比结果

### Agentend (`agentend/src/adapters/claude.py`)

- `_build_command` 末尾追加 `--include-partial-messages` flag
- `_parse_stream_line` 返回类型改为 `StreamEvent | None`，新增 `stream_event` 分支：
  - `content_block_delta` + `text_delta` → yield `StreamEvent(TEXT)`
  - `thinking_delta` 及元事件（`message_start` / `content_block_start` / `content_block_stop` / `message_delta` / `message_stop`）→ 忽略
- `_TYPE_MAP` 移除 `"assistant"` 映射（文本全权由 `content_block_delta` 输出，避免重复）
- `stream_chat` 增加 `if event:` 过滤 `None` 事件

### Frontend

无改动（已支持逐 chunk `STREAM_TEXT` dispatch）。

### Backend

无改动（SSE 透传行为不变）。

## 跨端影响

- **Frontend**: 无影响 — 前端 SSE 解析逻辑不变，只是收到更频繁、更小的 TEXT 事件
- **Backend**: 无影响 — 后端 SSE 透传不变
- **Agentend**: Claude adapter 从消息粒度改为逐 token 粒度输出 TEXT 事件

## 契约变更

- 无 schema 变更（`contracts/schemas/*.yaml` 未修改）
- `EventType` 枚举不变，adapter 行为完全向后兼容
