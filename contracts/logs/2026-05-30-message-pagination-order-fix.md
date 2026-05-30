# 2026-05-30: 消息分页排序修复 + 前端运行态块渲染

## 变更类型: FIX

## 影响范围
- Backend: `GET /api/tasks/:taskId/messages` 分页查询排序逻辑
- Frontend: `block-reducer.ts`、`MessageList.tsx`、`stores/chat.ts`
- AgentEnd: 无变更

## 变更原因
1. 后端分页按 `created_at ASC LIMIT 20` 返回最早 20 条，会话超过 20 条后刷新首屏拿不到最新消息
2. 前端 `MessageList` 未渲染 `runtimeBlocks`，子 Agent 运行态输出不可见
3. `runtime_completed` 事件覆盖已有块时丢弃 `streamingText`，导致输出完成后瞬间消失
4. 历史消息中裸 `type: runtime_status / json: ...` 片段被 Markdown 渲染成大量代码块

## 说明
- 后端分页改为 `ORDER BY id DESC LIMIT` + 内存反转，返回最新一页（API 响应格式不变）
- 前端 `block-reducer` 新增 legacy 裸行和 fenced 代码块中运行态片段的解析，合并碎片化 `streamingText`
- `MessageList` 将 `runtimeBlocks` 纳入 streaming 伪消息渲染
- Store `streamRuntimeEvent` 更新时保留已有 `streamingText`；`streamDone` 清空临时 `runtimeBlocks`

## 契约变更
- 无 schema 文件修改（Message 模型和 SSE 事件类型未变，变更为查询实现和前端渲染逻辑）
