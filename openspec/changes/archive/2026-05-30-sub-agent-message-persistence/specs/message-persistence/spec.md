## MODIFIED Requirements

### Requirement: 保存 agent message（增量持久化）
系统 SHALL 在 Agent 流开始时创建 Message 记录（status=streaming, content=""），在流式输出过程中通过 BatchWriter 增量 UPDATE content，流结束时设置 status=completed/failed。此要求 SHALL 同等适用于 orchestrator 和子 Agent（claude-code、opencode）。

#### Scenario: orchestrator agent message 创建
- **WHEN** 前端发送 `POST /api/tasks/:taskId/run` 到 orchestrator session
- **THEN** 后端创建 Message 记录（role=agent, status=streaming, content="", agent_type="orchestrator"）

#### Scenario: 子 Agent agent message 创建
- **WHEN** agentend 通过 `POST /api/tasks/:taskId/run` 请求执行子 Agent
- **THEN** 后端为子 Agent 创建独立的 Message 记录（role=agent, status=streaming, content="", agent_type=子Agent类型, session_id=子Agent的session_id）

#### Scenario: agent message 增量更新
- **WHEN** StreamWriter 触发 flush
- **THEN** UPDATE messages SET content=累积内容, last_seq=Redis Stream ID WHERE message_id=?

#### Scenario: agent message 完成
- **WHEN** goroutine 消费完 agentend 流（收到 Done event）
- **THEN** final flush 全部内容，UPDATE status=completed

#### Scenario: agent message 失败
- **WHEN** goroutine 消费过程中遇到错误或超时
- **THEN** flush 已有内容，UPDATE status=failed

#### Scenario: SSE 流中断（保留旧行为）
- **WHEN** SSE 流中途断开未收到 done 事件
- **THEN** goroutine 继续运行至完成或超时，已流式返回的部分内容通过 StreamWriter 持久化
