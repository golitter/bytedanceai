# 初始契约逆向提取

## 变更原因

建立 `contracts/` 契约层，将散落在 `agentend/src/schemas/` 中的跨端协议类型定义逆向提取为 YAML 格式的 JSON Schema 契约，作为三端共享的单一来源。

## 变更文件

- `contracts/schemas/event-types.yaml` — EventType 枚举 + StreamEvent 结构
- `contracts/schemas/agent-request.yaml` — AgentType 枚举 + AgentRequest 结构
- `contracts/schemas/agent-response.yaml` — AgentResponse 结构
- `contracts/schemas/session-state.yaml` — SessionState 枚举 + 合法转换表

## 对比结果

本次为初始提取，从以下源文件逆向生成：

| 契约 Schema | 源文件 |
|-------------|--------|
| event-types.yaml | `agentend/src/schemas/events.py` (EventType, StreamEvent) |
| agent-request.yaml | `agentend/src/schemas/request.py` (AgentType, AgentRequest) |
| agent-response.yaml | `agentend/src/schemas/response.py` (AgentResponse) |
| session-state.yaml | `agentend/src/session/models.py` (SessionState, _VALID_TRANSITIONS) |

## 跨端影响

- **Frontend**：无现有代码，生成类型为首次引入
- **Backend**：无现有代码，生成类型为首次引入
- **AgentEnd**：`src/schemas/` 文件改为从 `src/generated/` re-export，API 不变

## 契约变更

所有枚举值和字段定义与原始 Python 代码完全一致：

- EventType: init, text, tool_call, tool_result, artifact, planning, done, error
- AgentType: claude-code, opencode, orchestrator
- SessionState: idle, running, completed, interrupted, error
