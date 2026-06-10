## Why

当前 IM 聊天界面的实现与 `frontend/docs/visual-style-guide.md` 规范存在多处偏差（Agent 标识色完全不同、缺少头像光晕、图标库选型不一致等）。需要在实施后续功能前对齐视觉规范，避免偏差累积。

## What Changes

- 修正 Agent 标识色：Claude `#DA7756`、OpenCode `#10B981`、Orchestrator `#EAB308`
- 为 AgentAvatar 添加 8px 同色模糊光晕
- 修正状态灯动效：ready 脉冲（2s）、running 旋转（1.5s）
- 补充全局 letter-spacing: `-0.01em`
- 输入框背景改为 `var(--bg-card)`
- 侧栏列表项 padding 对齐为 `10px 12px`
- hover 过渡时长调整为 120ms ease-out
- 区分边框与分割线颜色（0.06 vs 0.04）
- 评估图标库迁移（Lucide → Phosphor Icons）

## Capabilities

### New Capabilities

- `visual-compliance`: 将现有组件的视觉样式对齐到 visual-style-guide.md 规范

### Modified Capabilities

## Impact

- `frontend/src/index.css` — CSS 变量值、全局字体间距
- `frontend/src/components/chat/AgentAvatar.tsx` — 头像光晕、状态动效
- `frontend/src/components/chat/MessageBubble.tsx` — 边框颜色
- `frontend/src/components/chat/MessageInput.tsx` — 输入框背景色
- `frontend/src/components/im/ConversationItem.tsx` — padding、过渡时长
- `frontend/src/components/im/ConversationList.tsx` — 分割线颜色
- `frontend/package.json` — 可能引入 Phosphor Icons
- 所有使用 Lucide 图标的组件（如迁移图标库）
