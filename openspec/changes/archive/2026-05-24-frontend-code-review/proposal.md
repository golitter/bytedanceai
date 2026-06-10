## Why

最近一次提交（agent 自定义名称 + 自动 session）引入了大量前端代码，对照 `visual-style-guide.md` 和 `development-strategy.md` 审查后，发现多处冗余代码、硬编码色值、重复常量和潜在缺陷，需要在功能迭代前清理。

## What Changes

- 提取 `AGENT_NAMES` 等重复常量到统一模块
- 将 AgentAvatar 中的硬编码色值替换为 CSS 变量
- 修复 `createConversation` 缺少 sessions 空数组保护
- 统一 `Session.agent_name` 类型为可选
- 提取重复的 hover 交互模式为工具 hook

## Capabilities

### New Capabilities

- `frontend-constants`: 集中管理 agent 名称、描述等常量，消除多处重复定义
- `hover-hook`: 提取重复的 mouseenter/mouseleave 交互为 `useHoverStyle` hook

### Modified Capabilities

## Impact

- `frontend/src/lib/constants.ts` — 新建，集中管理常量
- `frontend/src/hooks/use-hover-style.ts` — 新建，提取 hover 交互
- `frontend/src/components/chat/AgentAvatar.tsx` — 硬编码色值改 CSS 变量
- `frontend/src/components/chat/ChatArea.tsx` — 删除重复常量
- `frontend/src/components/im/ConversationItem.tsx` — 删除重复常量
- `frontend/src/components/im/ConversationList.tsx` — 使用 hover hook
- `frontend/src/components/im/NewChatDialog.tsx` — 使用 hover hook
- `frontend/src/lib/api.ts` — 修复类型安全 + 空数组保护
