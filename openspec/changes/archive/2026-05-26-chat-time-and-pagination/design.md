## Context

当前聊天系统的消息列表是纯线性的——只有消息气泡，没有时间维度信息。用户无法直观感知消息的时间跨度（"这是 5 分钟前的？还是昨天的？"）。后端 `GET /api/tasks/:taskId/messages` 一次性返回全部消息（无分页），随着对话增长会有性能问题。

前端已有 TanStack Virtual 虚拟滚动（消息 >50 条时启用），但缺少上拉加载历史消息的能力。

## Goals / Non-Goals

**Goals:**
- 在消息列表中根据时间跨度自动插入时间分隔线
- 支持上拉到顶部时加载更早的历史消息（cursor 分页）
- 加载历史消息后保持用户当前阅读位置不跳动
- 虚拟滚动与时间分隔线、分页加载正确配合

**Non-Goals:**
- 不做"刚刚/X分钟前"这类持续更新的相对时间（避免定时器开销）
- 不做消息时间戳的行内显示（如微信消息气泡内的小时间，后续可扩展）
- 不做 WebSocket 推送新消息（已有 SSE 流式推送）
- 不修改消息模型（不新增字段）

## Decisions

### 1. 时间分隔线：遍历时动态计算，不存入数据

**选择：** 在 `MessageList` 渲染时，遍历消息数组并计算相邻消息的时间差，动态插入 `TimeDivider` 元素。

**原因：** 时间分隔线是纯 UI 展示逻辑，不应污染数据层。对比在 store 中预处理出"消息+分隔线"混合数组，渲染时计算更简单、更符合 React 声明式风格。

**替代方案：** 在 store 中维护一个 `displayItems: (Message | TimeDivider)[]` 混合数组 — 增加了 store 复杂度，每次消息变化都要重算，不采用。

### 2. 分页策略：cursor-based（基于消息 ID）

**选择：** 使用消息的 `id`（uint 主键）作为游标，`before` 参数指定"加载此 ID 之前的消息"。

**原因：** 聊天场景标准做法。相比 offset 分页，cursor 分页不受新消息插入影响（offset 会偏移）。使用 uint 主键而非 `message_id`（UUID）是因为 uint 主键有天然递增顺序，查询更高效。

**API 设计：**
```
GET /api/tasks/:taskId/messages?limit=20&before=42
Response: { data: [...messages], has_more: true }
```
不传 `before` 时返回最新 N 条（向后兼容）。

### 3. 虚拟滚动 + 分页加载：扁平化数组

**选择：** 时间分隔线和消息统一放入一个扁平的 `displayItems` 数组，交给虚拟滚动渲染。

**原因：** TanStack Virtual 需要一个连续的列表。将消息和时间分隔线混合成统一数组，每个元素有 `type: 'message' | 'time-divider'`，虚拟滚动只需关心 count 和 estimateSize。

**替代方案：** 在虚拟列表外部渲染分隔线 — 无法配合虚拟化，长列表中会丢失分隔线。

### 4. 滚动位置恢复：记录 scrollHeight 差值

**选择：** 加载历史消息后，`scrollTop = newScrollHeight - oldScrollHeight`。

**原因：** prepend 消息后内容变长，scrollHeight 增大。用户原来看到的消息需要保持在视口内。用高度差值恢复是最简单可靠的方案。

### 5. 初始加载策略

**选择：** 进入聊天页面时，不传 `before`，只加载最新 `limit` 条消息。如果返回 `has_more: true`，启用上拉加载。

**原因：** 用户进入聊天最关心的是最新消息。如果对话很长，不需要一开始就加载全部历史。

## Risks / Trade-offs

- **[虚拟滚动动态高度]** 时间分隔线增加了列表项类型，需要为分隔线提供合理的 `estimateSize`（约 40px）→ 虚拟滚动已有动态测量机制，风险可控
- **[后端兼容性]** 修改 ListMessages API 但保持不传参数时返回全部消息 → 向后兼容，风险低
- **[快速上拉重复触发]** 用户快速上拉可能触发多次加载 → 用 `isLoading` 锁防止并发
- **[时间分隔线在虚拟列表中]** 虚拟化模式下，分隔线作为列表项参与虚拟化 → 已在 decision 3 中覆盖，统一处理
