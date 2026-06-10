## 1. 契约层：新增 runtime_text 事件类型

- [x] 1.1 在 `contracts/schemas/` 的 StreamEvent YAML 中添加 `runtime_text` 事件类型定义，content 包含 task_id、agent、text 字段
- [x] 1.2 运行 `make generate` 生成三端类型（Python generated/events.py、Go generated/、TypeScript types）

## 2. Agentend：ExecutionEngine 流式透传改造

- [x] 2.1 在 `agentend/src/schemas/events.py` 的 EventType 中添加 `RUNTIME_TEXT` 枚举值
- [x] 2.2 重构 `engine.py` 的 `_execute_task` 方法：将 `_collect_stream()` 从全量缓冲改为边收边 yield RUNTIME_TEXT + 保留 collected 列表
- [x] 2.3 实现进程内短路调用路径：根据 agent_type 查找 adapter 实例，直接调用 adapter.stream_chat()，跳过 BackendClient HTTP 回环
- [x] 2.4 短路路径中仍调用 BackendClient.run_task() 创建 message 记录，用 message_id 关联
- [x] 2.5 对不支持短路的 agent_type 保留 BackendClient.run_task() + stream_result() HTTP 路径，同样流式透传 RUNTIME_TEXT

## 3. Backend：SSE Handler 批量读取

- [x] 3.1 修改 `stream.go` serveStreaming 的实时阶段：XREAD 从 `Count:1, Block:5s` 改为 `Count:100, Block:200ms`
- [x] 3.2 批量读取后遍历写出所有 data 行，最后统一一次 Flush（替代逐条 Flush）
- [x] 3.3 StreamWriter 的 EventType 处理逻辑中新增 runtime_text 类型的处理（透传即可，不特殊解析）

## 4. Backend：StreamWriter 事件合并

- [x] 4.1 在 `writer.go` 的 Run scanFunc 回调中实现 TEXT 事件合并：累积不超过 500ms 或 2KB 后执行 XADD
- [x] 4.2 非 TEXT 事件（done、error、tool_call、runtime_text）先 flush 已累积文本再立即 XADD
- [x] 4.3 更新 flushThreshold 和 flushInterval 常量适配新的合并策略

## 5. Frontend：识别 RUNTIME_TEXT 事件

- [x] 5.1 在前端 chat store 的 SSE event handler 中添加 `runtime_text` 类型处理逻辑
- [x] 5.2 根据 task_id + agent 归属 RUNTIME_TEXT 文本到对应子 Agent 消息块，消息块不存在时自动创建
- [x] 5.3 MessageList 组件中渲染子 Agent 消息块的实时文本（复用现有 streamingContent 渲染逻辑）

## 6. 集成验证

- [x] 6.1 单 Agent 模式端到端测试：验证 token 流式延迟改善
- [x] 6.2 Orchestrator 多 Agent 模式端到端测试：验证子 Agent token 实时透传
- [x] 6.3 断线重连测试：验证 RUNTIME_TEXT 事件在 Redis Stream 中的回放正确性
