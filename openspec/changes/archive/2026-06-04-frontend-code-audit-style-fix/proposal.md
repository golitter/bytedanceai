## Why

前端新增功能模块（SkillsHub、ContactsPage、PlanReviewCard、AgentProfilePage 等）的视觉实现偏离了 visual-style-guide.md 和 development-strategy.md 定义的规范，具体表现为：Lucide 图标 `strokeWidth` 不统一（部分使用 1.5 而非规范的 1.25）、新模块使用 emoji 替代 Lucide 图标破坏一致性、按钮文字使用 `text-white` 而非 CSS 变量体系、部分组件使用 `text-red-500` 等直接色值而非语义 token。这些偏差在暗色主题下导致视觉不一致和可访问性问题。

## What Changes

- 修正全部 Lucide 图标 `strokeWidth` 为 `1.25`（8 处违规：RightSidebar、PlanReviewCard、ConversationItem、ContactsPage、AgentProfilePage）
- 将 SkillsHubPage 等新模块中的 emoji 图标（🛡️📦⚙️⚠️✅❌）替换为 Lucide React 图标，保持全站图标风格一致
- 将 `text-white` 替换为 `text-primary-foreground`（3 处：AgentProfilePage、SkillsHubPage x2）
- 将 `text-red-500` 等 Tailwind 直接色值替换为语义 token `text-destructive`（SkillsHubPage、DeleteConfirmDialog）
- 将 SessionCleanupPage 中 inline style `var(--text-secondary)` 替换为 `text-text-secondary` Tailwind 类

## Capabilities

### New Capabilities
- `icon-style-normalization`: 统一全站 Lucide 图标 strokeWidth=1.25，消除 emoji 图标，替换为 Lucide 图标
- `color-token-cleanup`: 消除 text-white、text-red-500 等硬编码/直接色值，统一使用语义 token

### Modified Capabilities
<!-- 无既有 spec 的需求变更 -->

## Impact

- 受影响文件：`frontend/src/components/chat/RightSidebar.tsx`、`frontend/src/components/cards/PlanReviewCard.tsx`、`frontend/src/components/im/ConversationItem.tsx`、`frontend/src/components/im/ContactsPage.tsx`、`frontend/src/pages/AgentProfilePage.tsx`、`frontend/src/pages/SkillsHubPage.tsx`、`frontend/src/pages/admin/SessionCleanupPage.tsx`
- 无 API / 后端影响，纯前端样式修正
- 无破坏性变更，所有修改为等价视觉替换
