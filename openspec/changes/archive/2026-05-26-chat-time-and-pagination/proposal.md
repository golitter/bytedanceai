## Why

当前聊天系统的消息列表缺少时间维度信息（用户无法感知消息的时间跨度），且历史消息一次性全量加载，在长对话场景下会导致性能问题。这两个问题直接影响用户对聊天节奏的感知和系统的可扩展性。

## What Changes

- **时间分隔线**：在消息列表中根据时间跨度插入分隔线，显示相对时间（"14:30"、"昨天 22:10"、"3天前"）
- **相对时间格式化工具**：新增 `utils/time.ts`，提供 `formatRelativeTime` 和 `shouldShowTimeSeparator` 函数
- **消息分页 API**：后端 `GET /api/tasks/:taskId/messages` 新增 cursor 分页支持（`limit` + `before` 参数），返回 `has_more` 标识
- **前端分页加载**：前端上拉到顶部时加载更早的消息，加载后恢复滚动位置

## Capabilities

### New Capabilities
- `message-time-display`: 消息列表的时间分隔线显示规则和相对时间格式化

### Modified Capabilities
- `message-persistence`: 消息历史加载 API 新增 cursor 分页参数（`limit`、`before`）和 `has_more` 返回字段
- `chat-ui`: MessageList 组件增加时间分隔线渲染和上拉加载历史消息交互

## Impact

- **后端 API**: `GET /api/tasks/:taskId/messages` 请求参数和响应结构变更（向后兼容，不传分页参数时行为不变）
- **前端组件**: `MessageList.tsx`（虚拟滚动 + 时间分隔线 + 上拉加载）、`MessageBubble.tsx`（无变化）
- **前端状态**: `stores/chat.ts` 新增 `hasMore` 状态和 `loadMoreMessages` action
- **前端工具**: 新增 `utils/time.ts` 时间格式化工具
- **前端 API**: `lib/api.ts` 的 `getTaskMessages` 新增分页参数
