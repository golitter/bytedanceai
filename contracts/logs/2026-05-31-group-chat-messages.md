# 契约变更：agent-request 新增 group_chat_messages 字段

**日期**: 2026-05-31

## 变更原因

群聊跨 Agent 记忆功能需要 Go Backend 在 RunTask 时查询其他 Agent 的窗口消息，并注入到发送给 AgentEnd 的请求体中。需要一个标准化字段来传递这些消息。

## 变更文件

- `contracts/schemas/agent-request.yaml` — 新增 `group_chat_messages` 字段

## 对比结果

```diff
 properties:
   config:
     type: ["object", "null"]
     additionalProperties: true
+  group_chat_messages:
+    type: array
+    default: []
+    description: "群聊窗口消息（Backend 注入，其他 Agent 自上次发言以来的消息）"
+    items:
+      type: object
+      properties:
+        role:
+          type: string
+        agent_name:
+          type: string
+        content:
+          type: string
```

## 跨端影响

| 端 | 影响 |
|----|------|
| Backend (Go) | `generated.AgentRequest` 新增 `GroupChatMessages` 字段，RunTask handler 填充 |
| AgentEnd (Python) | `generated.AgentRequest` 新增 `group_chat_messages` 字段，API 层解析 |
| Frontend (TS) | `generated.AgentRequest` 新增字段，但 Frontend 不直接使用（无影响） |
