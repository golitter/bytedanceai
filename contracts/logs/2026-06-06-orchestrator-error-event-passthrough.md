# Orchestrator ERROR 事件端到端透传

## 变更原因

Orchestrator reason 节点异常时，后端仅生成 `"Reasoning failed"` 硬编码文本，前端无法区分正常输出与错误状态，用户看不到具体错误原因。需要将已有的 `error` 事件类型从 Agent 端完整透传至前端，实现端到端错误可见性。

## 变更文件

**无 schema 变更** — `contracts/schemas/event-types.yaml` 中的 `error` 事件类型及其 `error` content 字段已满足需求，`make generate` 生成文件无差异。

本次变更仅涉及三端消费侧代码对已有契约的充分利用。

## 跨端影响

### AgentEnd
- `agentend/src/orchestrator/planning/graph.py` — reason 节点 catch 块携带异常类名和消息
- `agentend/src/adapters/orchestrator.py` — `_handle_reason` 处理 `output_type=error` 分支，生成 `EventType.ERROR` StreamEvent
- `agentend/tests/test_orchestrator_presentation.py` — 新增 `test_orchestrator_reason_error_becomes_stream_error_event`

### Backend
- `backend/internal/stream/writer.go` — 新增 `persistRuntimeBlockEvent` 处理 planning/runtime_status/coordination 事件；error 事件标记消息 `failed`；转发场景下正确创建子消息
- `backend/internal/stream/writer_test.go` — 新增 error 事件持久化测试 + mock DAO 实现
- `backend/internal/dao/gorm/message_dao.go` — 重构 group 消息可见性查询，用 EXISTS 子查询替代 IN 子查询
- `backend/internal/dao/gorm/message_dao_test.go` — 新增可见性 SQL 断言测试

### Frontend
- `frontend/src/stores/message-store.ts` — `streamError` 时若无 streaming content 则生成包含错误文本的失败消息
- `frontend/src/stores/__tests__/chat.test.ts` — 新增 streamError 空 content 场景测试
- `frontend/src/lib/ui-text.ts` — 新增 `FAILED` 状态标签
- `frontend/src/components/chat/ChatArea.tsx` — `UI_LABELS` 分组规范化
- `frontend/src/components/chat/SidebarActions.tsx` — `UI_LABELS` 分组规范化
- `frontend/src/pages/AgentProfilePage.tsx` — `UI_MESSAGES.IMPORTED` 迁移
- `frontend/src/components/chat/git-graph-types.ts` — 新增 `GitBranchInfo` 类型导出
