## Context

前端项目在 Phase 2-3 快速迭代中，积累了约 17 处与 visual-style-guide.md / development-strategy.md 不一致的代码。主要集中在三类问题：

1. **色彩类名误用**（`text-secondary` 暗色不可见、`#FFFFFF` 纯白文字）—— 直接影响用户体验
2. **组件/Hook 超限**（4 个文件 >300 行）—— 影响可维护性
3. **魔法字符串散落**（30+ 处中文 UI 字符串）—— 影响一致性和可维护性

当前前端技术栈：React 19 + Tailwind CSS 4 + shadcn/ui (radix-nova) + Zustand 5 + TanStack Query 5。

## Goals / Non-Goals

**Goals:**
- 修正所有导致暗色模式下文字不可见的色彩类名误用
- 拆分 4 个超限文件，每个拆分后主文件 ≤200 行
- 将散落的中文 UI 字符串集中到 `lib/ui-text.ts`
- 所有修改为纯重构，零功能变更

**Non-Goals:**
- 不引入新依赖或新设计模式
- 不改动组件 API（props 接口）—— 拆分后的子组件仅作为内部实现
- 不处理 `any` 类型（当前代码中未发现严重 `any` 问题）
- 不做国际化（i18n），仅集中管理现有中文字符串
- 不重构状态管理架构

## Decisions

### D1: 色彩修正策略 — 就地替换

**决策**: 直接将 `text-secondary` → `text-text-secondary`，`bg-secondary` → `bg-bg-card` 或 `bg-bg-hover`，`#FFFFFF` → `#E8EBF0`。

**理由**: Tailwind CSS 4 的色彩 token 体系已在 index.css 中完整定义（`--color-text-secondary` → `text-text-secondary`），只需对齐类名。不需要引入抽象层。

**替代方案**: 创建语义别名（如 `text-muted`）→ 过度抽象，违背 development-strategy.md "第三次重复再抽象"原则。

### D2: 组件拆分 — 内部子组件文件

**决策**: 将大组件按职责拆为同目录下的独立文件，主组件仅做编排。拆分后的子组件不导出给外部使用。

| 原文件 | 行数 | 拆分方案 |
|--------|------|----------|
| MessageBubble.tsx | 390 | → BlockRenderer.tsx (110行) + AgentMessageContent.tsx (125行) + MessageBubble.tsx (155行) |
| RightSidebar.tsx | 471 | → AgentInfoSection.tsx (57行) + SidebarPathSection.tsx (58行) + SidebarActions.tsx (56行) + useCollapsible.ts (22行) + RightSidebar.tsx (225行) |
| GitGraphPanel.tsx | 322 | → GraphRenderer.tsx (90行) + GraphBranchLabels.tsx (42行) + GraphTooltip.tsx (20行) + GitGraphPanel.tsx (170行) |
| use-chat-stream.ts | 350 | 保持整体（SSE 事件处理是紧密耦合的状态机，强行拆分会增加复杂度） |

**理由**: MessageBubble、RightSidebar、GitGraphPanel 有清晰的职责边界，拆分收益明确。use-chat-stream 是 SSE 状态机，事件处理与状态流转紧密耦合，拆分反而增加理解成本。

**替代方案**: 对 use-chat-stream 抽出 useSSEConnection → 虽可减少行数，但 SSE 事件处理需要直接 dispatch 到 store，抽离后回调链更复杂，得不偿失。

### D3: UI 字符串集中 — 按语义分组常量对象

**决策**: 在 `src/lib/ui-text.ts` 中按语义分组（UI_ACTIONS、UI_STATUS、UI_MESSAGES、UI_LABELS），使用 `as const` 确保类型安全。组件中直接引用常量。

**理由**: 按语义分组比单层平铺更易维护，`as const` 提供字面量类型。不需要 i18n 库的复杂度。

**替代方案**: 使用 i18next 等 i18n 库 → 过度工程，当前只需集中管理。

### D4: shake 动画保留

**决策**: 保留现有 shake 动画（已使用 `transform: translateX`），仅确认合规。

**理由**: 审计发现 shake 动画实际上只操作 `transform` 属性，符合 style guide "只动画 transform 和 opacity" 的约束。无需修改。

## Risks / Trade-offs

- **[拆分引入的导入链]** → 子组件仅在 chat 目录内部引用，不影响外部 API。通过 TypeScript 编译验证无遗漏。
- **[常量引用遗漏]** → 逐文件 grep 替换，确保所有散落字符串都已集中。编译时 TypeScript 会捕获拼写错误。
- **[use-chat-stream 不拆分]** → 350 行超过 200 行红线，但作为状态机其内部逻辑不可分割。接受此 trade-off，通过注释分区提升可读性。
