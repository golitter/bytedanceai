## 1. 时间格式化工具

- [x] 1.1 新建 `frontend/src/utils/time.ts`，实现 `formatRelativeTime(timestamp: number): string` 函数（今天 "HH:mm"、昨天 "昨天 HH:mm"、2-7天 "N天前"、今年 "M月D日 HH:mm"、跨年 "YYYY年M月D日"）
- [x] 1.2 在同文件中实现 `shouldShowTimeSeparator(prevTimestamp?: number, currentTimestamp: number): boolean` 函数（首条消息、间隔 >5分钟、跨天触发）

## 2. TimeDivider 组件

- [x] 2.1 新建 `frontend/src/components/chat/TimeDivider.tsx`，接收 `timestamp` prop，渲染居中时间分隔线（muted 颜色、12px 字号、两侧横线装饰、上下 8px 内边距）

## 3. 后端消息分页 API

- [x] 3.1 修改 `backend/internal/handler/message.go` 的 `ListMessages`，新增 `limit`（默认值 0 表示不限）和 `before`（消息 uint ID）查询参数解析
- [x] 3.2 修改 GORM 查询：有 `before` 时添加 `WHERE id < before` 条件，有 `limit` 时 `LIMIT limit+1`（多取一条判断 has_more）
- [x] 3.3 修改响应结构：返回 `{ data: [...messages], has_more: boolean }`，不传 limit/keep 时保持原行为（返回全部消息、has_more=false）

## 4. 前端 API 层

- [x] 4.1 修改 `frontend/src/lib/api.ts` 的 `getTaskMessages` 函数，新增 `limit?: number` 和 `before?: number` 可选参数
- [x] 4.2 修改响应类型，解析 `has_more` 字段

## 5. 前端 Store 分页状态

- [x] 5.1 修改 `frontend/src/stores/chat.ts`，在 session 状态中新增 `hasMore: boolean` 字段（默认 true）
- [x] 5.2 新增 `loadMoreMessages` action：调用 API 获取更早消息（传入首条消息 ID 作为 before），prepend 到 messages 数组，更新 hasMore
- [x] 5.3 新增 `isLoadingMore: boolean` 状态，防止并发加载

## 6. MessageList 集成

- [x] 6.1 修改 `frontend/src/components/chat/MessageList.tsx`，构建 `displayItems` 扁平数组：遍历消息时调用 `shouldShowTimeSeparator` 判断，在满足条件处插入 `{ type: 'time-divider', timestamp }` 元素
- [x] 6.2 修改虚拟滚动和非虚拟渲染逻辑，根据 `displayItem.type` 分别渲染 `TimeDivider` 或 `MessageRenderer`
- [x] 6.3 新增上拉加载逻辑：监听 `scrollTop === 0` 且 `hasMore && !isLoadingMore` 时触发 `loadMoreMessages`
- [x] 6.4 实现滚动位置恢复：加载前记录 `scrollHeight`，加载后设置 `scrollTop = newScrollHeight - oldScrollHeight`
- [x] 6.5 在加载中状态显示顶部加载指示器（小型 spinner 或 "加载中..." 文本）
