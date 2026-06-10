## Why

前端代码在快速迭代中积累了多处与 `visual-style-guide.md` 和 `development-strategy.md` 不一致的实现：硬编码颜色、纯白文字、Agent 标识色越界使用、头像阴影违规、动画缓动函数错误、字号超出定义体系、内联样式可 Tailwind 化但未做。需要一次系统性审计将代码对齐到设计规范，确保 UI 一致性。

## What Changes

- 修正 `text-white` → `text-primary-foreground`（AgentAvatar、GroupAvatar 等组件）
- 移除 AgentAvatar 上的 `box-shadow` 发光效果（违反"几乎不用阴影"原则）
- Agent 标识色仅限头像 + 消息色条场景（修正 PlanCard 等越界使用）
- `rounded-full` 仅保留给 Badge/胶囊标签，移除卡片/按钮上的误用
- 动画缓动统一为 `ease-out`（修正 AgentAvatar 状态灯 `ease-in-out`）
- 字号对齐到 11/12/13/14/20px 体系（修正 GroupAvatar 8px/9px）
- 内联样式迁移至 Tailwind 工具类（AgentAvatar、GroupAvatar、MessageInput）
- PlanCard agent badge 改用语义色而非 Agent 标识色

## Capabilities

### New Capabilities
- `style-enforcement`: 对现有组件执行视觉规范修正，确保与 visual-style-guide.md 对齐

### Modified Capabilities

## Impact

- `frontend/src/components/chat/` — AgentAvatar、GroupAvatar、MessageBubble、AgentHoverCard
- `frontend/src/components/cards/` — PlanCard、CoordChannel、RuntimeStatus
- `frontend/src/components/im/` — ConversationItem
- `frontend/src/index.css` — 可能需要补充 `--text-primary-foreground` 等缺失变量
- 纯样式修改，不影响功能逻辑和 API 契约
