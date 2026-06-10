## 1. 集中常量

- [x] 1.1 新建 `frontend/src/lib/constants.ts`，定义 `AGENT_NAMES` 和 `AGENT_DESCRIPTIONS`（从 ChatArea、ConversationItem、NewChatDialog 中提取合并）
- [x] 1.2 `ChatArea.tsx` 删除内联 `AGENT_NAMES`，改从 `lib/constants` 导入
- [x] 1.3 `ConversationItem.tsx` 删除内联 `AGENT_NAMES`，改从 `lib/constants` 导入
- [x] 1.4 `NewChatDialog.tsx` 删除内联 `AGENT_DESCRIPTIONS`，改从 `lib/constants` 导入

## 2. useHoverStyle hook

- [x] 2.1 新建 `frontend/src/hooks/use-hover-style.ts`，实现 `useHoverStyle(hoverBg?, normalBg?)` hook
- [x] 2.2 `ConversationItem.tsx` 使用 `useHoverStyle` 替换手写 hover
- [x] 2.3 `ConversationList.tsx` 使用 `useHoverStyle` 替换手写 hover
- [x] 2.4 `NewChatDialog.tsx` 使用 `useHoverStyle` 替换手写 hover

## 3. 类型安全与防御

- [x] 3.1 `api.ts` Session 接口 `agent_name` 改为可选（`agent_name?: string`）
- [x] 3.2 `api.ts` `createConversation` 添加 `detail.sessions[0]` 空数组保护

## 4. 验证

- [x] 4.1 TypeScript 编译无报错
