## Why

前端已有 10+ 份 style/spec 规范（color-token-enforcement、style-enforcement、transition-compliance 等），但代码审计发现实际实现中仍存在与规范脱节的硬编码色值、内联样式、className 拼接不规范等问题。同时 `dev-practice-audit` 规范（组件分层、状态管理、Hook 规范）尚未落地。需要一次全面的代码审计与修正，将现有 skills 定义的设计规范真正贯彻到每一行前端代码中。

## What Changes

- **全面扫描 src/ 下所有 .tsx/.css 文件**，对照 visual-style-guide.md 和 development-strategy.md 逐项检查违规
- **消灭硬编码色值**：将所有 `#fff`、`rgba(0,0,0,0.75)`、`rgba(34,197,94,0.1)` 等替换为 CSS 变量或 Tailwind 语义类
- **统一 border-radius**：消除 `rounded-[6px]` 等自定义值，使用 CSS radius token（rounded-sm/md/lg/xl）
- **内联样式迁移**：将 `style={{ color: '#fff' }}` 等 inline style 迁移到 Tailwind 工具类
- **transition 合规**：确保所有过渡只涉及 transform/opacity，时长 120-300ms ease-out
- **组件分层审查**：检查 Page/Smart/Dumb 三层模型是否遵循，拆分违反单一职责的巨型组件
- **状态管理审查**：Server State 走 TanStack Query、Global State 走 Zustand、Local State 走 useState，消灭跨界使用
- **Hook 规范审查**：纯计算函数不写 Hook、Hook 体量红线 150-200 行
- **TypeScript 严格性**：消灭 `any` 类型、魔法字符串，利用 discriminated unions

## Capabilities

### New Capabilities
- `code-audit-enforcement`: 前端代码全面审计与修正——涵盖色彩变量合规、组件架构分层、状态管理分类、Hook 规范、TypeScript 严格性，确保 visual-style-guide.md 和 development-strategy.md 的所有规则在代码中 100% 落地

### Modified Capabilities
<!-- 无需修改已有 spec 的行为要求，本次是实现层面落地 -->

## Impact

- **代码文件**：`frontend/src/` 下所有 `.tsx`、`.css` 文件，预计涉及 20-40 个文件
- **依赖**：无新依赖引入
- **视觉影响**：修正后视觉一致性和代码可维护性显著提升，但无功能行为变更
- **风险**：色值和类名替换属于低风险修改，但需逐文件验证无回归
