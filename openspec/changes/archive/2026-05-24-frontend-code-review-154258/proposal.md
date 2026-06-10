## Why

前端代码与两份核心规范文档（visual-style-guide.md、development-strategy.md）存在系统性偏差。核心问题包括：大量组件绕过 CSS 变量系统直接使用 Tailwind 原始色阶、VSG 与开发策略文档之间存在矛盾（图标库选择）、存在死代码和设计冗余。需要在开发 Phase 2 Chat UI 前清理这些技术债，否则后续开发会基于不一致的基础继续累积问题。

## What Changes

- 修复 SessionList.tsx 等组件中的原始 Tailwind 色阶（`bg-gray-200` 等），迁移到 CSS 变量系统
- 消除 `text-white` / `#fff` 纯白文字使用，替换为 `--text-primary` (`#E8EBF0`)
- 解决文档矛盾：统一图标库选择（Lucide vs Phosphor）
- 消除 AgentAvatar.tsx 中的硬编码 hex 色值，迁移到 CSS 变量
- 修复 ChatArea.tsx 中的浅色主题颜色残留（`#FEF2F2`）
- 清理死代码：空 store、未使用组件、未使用 mutations
- 修复 use-chat-stream.ts 中的魔法字符串，使用契约生成的常量
- 补齐缺失的 CSS 变量（代码块背景 `#0D0F14` 等）
- 修复 AgentAvatar 状态点尺寸（8px → 4px）和冗余常量定义
- 修复 MessageList 中非虚拟化路径缺失的 props 传递

## Capabilities

### New Capabilities
- `css-variable-audit`: 全面审计并修复 CSS 变量系统覆盖缺失，确保所有组件使用设计令牌而非硬编码色值
- `dead-code-cleanup`: 识别并清理前端死代码、冗余常量、未使用的 store/mutations
- `doc-consistency-fix`: 解决 visual-style-guide.md 与 development-strategy.md 之间的矛盾，统一规范

### Modified Capabilities
- `visual-compliance`: 现有视觉合规 spec 需要扩展覆盖范围（SessionList 色阶、AgentAvatar 硬编码、纯白文字禁止等）
- `frontend-constants`: 补充 use-chat-stream 中事件类型的常量定义，消除魔法字符串

## Impact

- 受影响代码：`src/components/` 下几乎所有组件、`src/stores/`、`src/hooks/`、`src/lib/api.ts`
- 受影响文档：`frontend/docs/visual-style-guide.md`（图标部分）
- CSS 变量系统：`src/index.css` 需补充缺失变量
- 无 API 变更，无 breaking change
