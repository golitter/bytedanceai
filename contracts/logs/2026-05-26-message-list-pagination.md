# 2026-05-26: 消息列表 cursor 分页 API

## 变更类型: MODIFY

## 影响范围
- Backend: `GET /api/tasks/:taskId/messages` 新增分页参数和响应包装
- Frontend: `api.ts`、`stores/chat.ts`、`MessageList.tsx` 适配新分页协议
- AgentEnd: 无变更

## 说明
- `ListMessages` API 新增 `limit`（默认 20）和 `before`（消息 uint ID）查询参数
- 响应结构从 `[...messages]` 变为 `{ data: [...messages], has_more: boolean }`
- 不传 `limit` 和 `before` 时行为不变（返回全部消息，`has_more=false`），向后兼容
- 前端初始加载传 `limit=20`，上拉到顶部时传首条消息 `dbId` 作为 `before` 翻页

## 契约变更
- 无 schema 文件修改（Message 数据模型未变，变更为 API 层查询参数和响应包装）
