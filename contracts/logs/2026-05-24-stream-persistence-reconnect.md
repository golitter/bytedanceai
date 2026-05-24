# 流式持久化 + 断线重连

## 变更原因

Agent 流式输出只在 SSE 流完全结束后才写入 MySQL，前端切 Tab、刷新或断网时正在输出的内容全部丢失。需要增量持久化 + 断线重连机制。

## 变更文件

- `schemas/message.yaml` — Message 模型新增字段

## 对比结果

### message.yaml 新增

| 字段 | 类型 | 说明 |
|------|------|------|
| `message_id` | string \| null | 对外唯一标识（UUID），用于流式追踪和重连 |
| `status` | string (streaming/completed/failed) | 消息流式状态，默认 completed |
| `last_seq` | string | Redis Stream 最后已 flush 的 event ID，默认空 |

新增枚举 `MessageStatus`：streaming / completed / failed

## 跨端影响

- **Backend**：Message 模型新增 message_id（UNIQUE）、status、last_seq 字段；POST /run 改为异步返回 202 + message_id；新增 GET /stream SSE 端点
- **Frontend**：Zustand 全局 store 替代 useReducer；SSE 改为 GET 请求 + 断线重连；加载历史时检测 streaming 消息自动重连
- **Agentend**：无改动

## API 变更

| 端点 | 变更 |
|------|------|
| `POST /api/tasks/:taskId/run` | 响应从 SSE 流改为 202 Accepted `{ message_id, status }` |
| `GET /api/tasks/:taskId/stream?session_id=&message_id=` | 新增 SSE 端点，支持历史回放 + 实时续接 |
