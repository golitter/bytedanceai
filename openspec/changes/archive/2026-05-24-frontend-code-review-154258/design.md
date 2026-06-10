## Context

前端代码库在快速迭代中积累了与两份核心规范的系统性偏差：

1. **visual-style-guide.md** 定义了完整的 Dark Utilitarian 设计系统（CSS 变量、色阶、动效），但部分组件绕过了变量系统直接使用 Tailwind 原始色阶或硬编码 hex 值。
2. **development-strategy.md** 定义了状态管理、组件拆分、代码健壮性等工程规范，但存在死代码、魔法字符串、文档间矛盾。

两份文档自身也存在矛盾：VSG 指定 Phosphor Icons，开发策略指定 Lucide React，实际代码使用 Lucide。需要统一。

当前前端目录结构：`src/components/`（20+ 组件）、`src/hooks/`、`src/stores/`、`src/lib/`、`src/generated/`。

## Goals / Non-Goals

**Goals:**
- 所有组件统一使用 CSS 变量系统，消除硬编码色值
- 解决 VSG 与开发策略之间的文档矛盾
- 清理死代码和冗余定义
- 补齐缺失的 CSS 变量和设计令牌
- 确保后续 Phase 2 开发基于一致的基础

**Non-Goals:**
- 不重构组件架构或调整组件拆分方式
- 不引入新的第三方依赖
- 不修改后端 API 或 contracts
- 不做性能优化（虚拟列表等属于 Phase 2 范畴）

## Decisions

### D1: 图标库统一为 Lucide React

**选择**：保留 Lucide React，更新 VSG 文档。

**理由**：代码库全面使用 Lucide，开发策略明确指定 Lucide，pnpm 已安装。Lucide 生态更成熟，与 shadcn/ui 深度集成。VSG 中的 Phosphor 建议仅是风格偏好，Lucide 通过 `strokeWidth={1.25}` 可模拟细线视觉。

**替代方案**：迁移到 Phosphor — 成本高，无功能收益，且 shadcn/ui 内置 Lucide。

### D2: CSS 变量补齐策略

**选择**：在 `index.css` 的 `:root` / `.dark` 中补充缺失变量（`--code-bg`、`--color-danger-bg` 等），组件逐步迁移。

**理由**：保持设计令牌单一来源。硬编码值散落在各组件中难以维护和主题切换。

**替代方案**：直接使用 Tailwind 色阶 — 简单但破坏主题系统一致性。

### D3: 死代码清理范围

**选择**：删除明确无引用的死代码（空 store、未导入组件、未使用 mutations），保留可能有未来用途的类型定义。

**理由**：死代码增加认知负担，误导开发者。类型定义成本低且有助于 IDE 提示，暂保留。

### D4: 魔法字符串修复方式

**选择**：使用 `generated/` 中的契约类型，补充 `lib/constants.ts` 中缺失的事件类型常量。

**理由**：契约优先原则——`contracts/schemas/` 为单一来源，`make generate` 生成 TypeScript 类型。手写类型与契约类型并存会导致漂移。

## Risks / Trade-offs

- **[CSS 变量迁移可能引入视觉回归]** → 逐组件修改后视觉验证，使用 dev server 对比修改前后效果
- **[删除死代码可能影响未发现的引用]** → 使用 TypeScript 编译器 + grep 确认无引用后再删除
- **[文档更新可能改变团队预期]** → 统一文档后通知团队，VSG 图标部分从 Phosphor 改为 Lucide 是唯一实质性变更
- **[AgentAvatar 常量合并可能影响导入方]** → 合并后更新所有导入路径，TypeScript 编译会捕获遗漏
