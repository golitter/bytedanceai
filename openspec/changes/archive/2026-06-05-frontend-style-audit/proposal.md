## Why

前端代码在快速迭代中积累了与 [visual-style-guide.md](../../frontend/docs/reference/visual-style-guide.md) 及 [development-strategy.md](../../frontend/docs/guides/development-strategy.md) 不一致的实现：`text-secondary` 误用导致暗色背景下文字不可见、部分组件超过 300 行未拆分、shake 动画违反"仅 transform/opacity"约束、`#FFFFFF` 纯白文字违反色彩禁令。这些问题影响可读性和可维护性，需要系统性修正。

## What Changes

- **修正 `text-secondary` 误用**：MarkdownRenderer 中 2 处 `text-secondary` 改为 `text-text-secondary`，避免暗色背景不可见
- **修正纯白文字**：`--prose-heading-h1: #FFFFFF` 改为 `#E8EBF0`
- **修正 `bg-secondary` 误用**：ToolCard、GitGraphPanel 中 2 处 `bg-secondary` 评估并替换为语义正确的 `bg-bg-card` 或 `bg-bg-hover`
- **修正动画约束**：shake 动画改为纯 `transform: translateX`（已是，仅确认合规）
- **拆分超限组件**：MessageBubble（389 行）、RightSidebar（470 行）、GitGraphPanel（321 行）按职责拆分为子组件
- **拆分超限 Hook**：use-chat-stream（349 行）按职责拆分为多个小 Hook
- **集中 UI 字符串**：散落在组件中的中文 UI 字符串提取到 constants

## Capabilities

### New Capabilities
- `color-usage-normalization`: 修正 Tailwind CSS 色彩类误用（text-secondary → text-text-secondary、bg-secondary 替换、纯白文字修正），确保暗色模式下所有文字可见
- `component-split-refactor`: 拆分 MessageBubble、RightSidebar、GitGraphPanel、use-chat-stream 四个超限文件，按职责和变化频率拆为独立子模块
- `ui-strings-centralization`: 将散落在组件中的中文 UI 字符串集中到 constants，消灭魔法字符串

### Modified Capabilities
（无既有 spec 需要修改）

## Impact

- **影响范围**：`frontend/src/components/`（chat、cards、im、layout）、`frontend/src/hooks/`、`frontend/src/index.css`、`frontend/src/lib/constants.ts`
- **风险**：纯重构，无 API 变更、无功能行为变更，风险低
- **依赖**：无新增依赖
- **测试**：通过现有页面手动验证（暗色主题下文字可见性、组件功能无回归）
