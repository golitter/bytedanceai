## 1. 验证 CLI 流式能力

- [x] 1.1 实测 OpenCode CLI `--format json` 输出粒度：运行 `opencode run "说三个字" --format json`，观察 text 事件是逐 token 还是完整消息
- [x] 1.2 实测 Claude CLI `--include-partial-messages` 输出格式：确认 `content_block_delta` 事件结构和频率

## 2. ClaudeCodeAdapter 改造

- [x] 2.1 `_build_command` 追加 `--include-partial-messages` flag
- [x] 2.2 `_parse_stream_line` 新增 `type: "stream_event"` 分支：解析 `event.type`，`content_block_delta` 时提取 `event.delta.text` yield TEXT 事件，其他 stream_event 子类型忽略
- [x] 2.3 调整 `_TYPE_MAP`：`"assistant"` 映射改为忽略（或仅在无 stream_event 时映射），避免重复文本
- [x] 2.4 验证：curl 后端 SSE 端点，确认逐 token 事件输出

## 3. OpenCodeAdapter 改造（视实测结果）

- [x] 3.1 若 OpenCode 已逐 token：无需改动，确认现有解析逻辑正确
- [x] 3.2 若 OpenCode 非逐 token：调研替代方案（streaming flag、SDK 直接调用等），更新 adapter

## 4. 端到端验证

- [x] 4.1 启动全三端，前端发送消息，确认 Agent 回复逐字显示
- [x] 4.2 测试流式中断（abort）功能正常
