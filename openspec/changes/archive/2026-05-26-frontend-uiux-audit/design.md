## Context

AgentHub 前端采用 Dark Utilitarian 视觉风格（参考 Linear / Cursor），已有完整的 visual-style-guide.md（色彩、字体、圆角、边框、阴影、动效、Agent 身份系统）和 development-strategy.md（三层组件模型、状态管理分类、性能策略、Hook 规范）。初始审计显示代码库整体合规度较高，但存在以下已知偏差：

- `MessageBubble.tsx` 中 BlockRenderer 使用 index 作为 key
- CSS 自定义属性中缺少 Agent 标识色的统一定义（`--agent-claude` 等通过内联 style 动态拼接）
- 部分组件的 Tailwind 类名中存在硬编码色值或非标准 token

本设计文档描述如何系统性地审计、修复和建立持续合规机制。

## Goals / Non-Goals

**Goals:**

- 对 `frontend/src/` 全部组件执行逐项视觉风格合规审查，输出完整违规清单
- 对 `frontend/src/` 全部组件执行逐项开发实践合规审查，输出完整违规清单
- 将 CSS 自定义属性与 visual-style-guide.md 完全对齐，补齐缺失 token
- 修复所有已确认的违规项
- 为后续开发提供可参考的合规检查清单

**Non-Goals:**

- 不做视觉重设计或 UI 改版
- 不引入新的 CSS-in-JS 方案或样式系统
- 不修改后端或 Agent 端代码
- 不做自动化 CI 检查流水线（留作后续 change）
- 不引入新的第三方依赖

## Decisions

### D1: 审计方式——静态代码扫描 + 人工审查

**选择**：通过 grep/AST 扫描 Tailwind 类名和 CSS 属性，结合逐文件人工审查。

**理由**：项目规模尚小（~20 组件），全量人工审查可控。自动化 lint（如 stylelint 自定义规则）投入产出比不高，等项目规模增长后再引入。

**替代方案**：Stylelint 自定义规则 → 项目太小，规则开发成本 > 人工审查成本。

### D2: CSS 变量对齐——在 index.css 的 :root / .dark 中统一定义

**选择**：将 visual-style-guide.md 中所有色值、字体、圆角档位定义为 CSS 自定义属性，组件通过 `var(--xxx)` 引用。

**理由**：与 Tailwind CSS 4 的 CSS 变量双主题机制一致，与 development-strategy.md 中"用 CSS 变量实现主题"原则一致。当前 Agent 标识色在 MessageBubble 中通过内联 style 动态拼接变量名（`var(--agent-${type})`），应在 index.css 中统一定义这些变量。

**替代方案**：Tailwind theme.extend 直接定义 → 与现有 index.css 模式不一致。

### D3: 修复策略——最小变更原则

**选择**：只修复合规违规项，不做连带重构。每个修复独立、可回滚。

**理由**：审计的目的是对齐规范，不是借机重构。混合修改会增加 review 难度和回归风险。

### D4: 违规分级——Critical / Warning / Info

**选择**：
- **Critical**：直接违反色彩禁令（渐变、backdrop-blur、纯白纯黑文字）、性能反模式（index 作 key、server state 进 Zustand）
- **Warning**：非标准圆角值、缺失 CSS 变量、可优化但无实际影响的模式
- **Info**：风格建议、更优但非必要的改进

**理由**：分级让修复有优先级，Critical 必修，Warning 建议修，Info 可选。

## Risks / Trade-offs

- **CSS 变量补齐可能影响现有组件渲染** → 逐变量添加，每添加一个变量后验证视觉无变化（新变量只补定义，不改引用方式）
- **修复 index key 可能暴露 block 数据缺少唯一 ID 的问题** → 先检查 MessageBlock 类型是否有 id 字段，若无则在 block-reducer 生成时补充
- **人工审查可能遗漏** → 对照 visual-style-guide.md 逐节检查，使用 checklist 确保覆盖
