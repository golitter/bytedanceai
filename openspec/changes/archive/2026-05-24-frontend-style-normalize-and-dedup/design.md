## Context

前端当前存在两套并行的 CSS 变量体系：shadcn/ui 的 27 个语义 token（`--background`、`--foreground`、`--primary` 等）和 16 个自定义 app token（`--bg-canvas`、`--text-primary`、`--color-brand` 等）。两者指向相同的底色值，但应用层组件全部通过 `style={{ color: 'var(--text-primary)' }}` 内联样式引用自定义 token，而 shadcn 原语组件（`ui/*`）通过 Tailwind 类引用语义 token。此外有 11 处硬编码的 `rgba(255,255,255,0.06)` 边框色、6 处硬编码状态色、3 个未使用的 CSS 变量、以及 `App.tsx` 和 `use-sessions.ts` 两处死代码。

核心约束：所有样式变更必须符合 `visual-style-guide.md` 的 Dark Utilitarian 规范（禁止阴影、禁止渐变、禁止纯白/纯黑文字、动画限于 120-300ms transform/opacity）。

## Goals / Non-Goals

**Goals:**

- 消除 CSS 变量双轨制，统一到 shadcn 语义 token + 少量自定义扩展
- 将所有内联 `style={{}}` 迁移到 Tailwind 工具类，使 `@theme inline` 桥接层真正发挥作用
- 消除所有硬编码颜色值，统一通过 CSS 变量引用
- 移除死代码和未使用的 CSS 变量
- 确保迁移后视觉表现不变

**Non-Goals:**

- 不做组件结构重构（不改组件拆分方式或状态管理模式）
- 不做 Light mode 支持（当前只有 dark mode 在用，自定义 token 无 `:root` 定义）
- 不更新 `docs/impl/` 下的过期文档（那是独立的文档维护任务）
- 不引入新的 UI 库或设计系统

## Decisions

### D1: 统一到 shadcn 语义 token，移除重复的自定义 app token

**决策：** 将自定义 app token 逐一映射到 shadcn 语义 token，移除重复定义。

| 自定义 token | → shadcn token | Tailwind 类 |
|---|---|---|
| `--bg-canvas` | `--background` | `bg-background` |
| `--bg-sidebar` | `--sidebar` | `bg-sidebar` |
| `--bg-hover` | `--accent` | `bg-accent` |
| `--text-primary` | `--foreground` | `text-foreground` |
| `--text-secondary` | `--muted-foreground` | `text-muted-foreground` |
| `--color-brand` | `--primary` | `text-primary` / `bg-primary` |
| `--color-error` | `--destructive` | `text-destructive` / `bg-destructive` |
| `--divider` | `--border` | `border-border` |

**保留的自定义 token**（无 shadcn 对应项）：

| Token | 用途 | Tailwind 桥接 |
|---|---|---|
| `--text-tertiary` (#5A6070) | 三级文字 | `text-tertiary` (需加 `@theme` 入口) |
| `--bg-active` (#2C313B) | 激活态背景 | `bg-active` (需加 `@theme` 入口) |
| `--code-bg` (#0D0F14) | 代码块背景 | `bg-code` (需加 `@theme` 入口) |
| `--color-danger-bg` (rgba) | 错误背景色 | `bg-danger` (需加 `@theme` 入口) |
| `--agent-claude` / `--agent-opencode` / `--agent-orchestrator` | Agent 标识色 | `text-agent-claude` 等 |

**理由：** 两套系统指向相同颜色是维护负担的根源。shadcn 语义 token 已在 `@theme inline` 中桥接到 Tailwind，而自定义 token 没有这个桥接，导致组件被迫用内联样式。统一后可全面使用 Tailwind 类。

**备选方案：** 保留双轨但给自定义 token 也加 `@theme` 映射 — 但这意味着两套命名共存，每个开发者需要记住两套同义词。

### D2: 内联样式迁移到 Tailwind 工具类

**决策：** 将所有 `style={{ color: 'var(--xxx)' }}` 替换为对应的 Tailwind 类。

**规则：**
- `style={{ color: 'var(--text-primary)' }}` → `text-foreground`
- `style={{ backgroundColor: 'var(--bg-canvas)' }}` → `bg-background`
- `style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}` → `border-b border-border`
- `style={{ borderColor: 'var(--color-brand)' }}` → `border-primary`

**例外：** 动态计算的样式（如基于 state 的条件颜色）保留内联样式，但用 CSS 变量替代硬编码值。

**理由：** 内联样式无法被 Tailwind 的 purge 优化，不利于主题切换，且比 Tailwind 类冗长。迁移到 Tailwind 类可减少代码量约 30-40%。

### D3: 硬编码颜色替换为 CSS 变量

**决策：** 将 11 处 `rgba(255,255,255,0.06)` 替换为 `var(--border)`，将 6 处硬编码状态色替换为对应 CSS 变量。

### D4: 死代码移除清单

**决策：** 移除以下确认无引用的文件和变量：

| 目标 | 类型 |
|---|---|
| `App.tsx` (returns null, never rendered) | 文件 |
| `hooks/use-sessions.ts` (never imported) | 文件 |
| `--bg-active` (未使用) | CSS 变量 |
| `--color-success` (未使用) | CSS 变量 |
| `--color-warning` (未使用) | CSS 变量 |
| `--divider` (与 `--border` 重复) | CSS 变量 |
| 所有被 D1 映射掉的自定义 app token | CSS 变量 |

**关于 `ui/button.tsx`、`ui/card.tsx`、`ui/input.tsx`：** 虽然当前未被应用组件直接引用，但它们是 shadcn/ui 的基础组件，保留作为未来使用的基础设施。不移除。

### D5: `text-brand` 修复

**决策：** `MessageBubble.tsx` 中使用的 `text-brand` Tailwind 类当前无 `@theme` 映射，静默失效。迁移后统一使用 `text-primary`（对应 `--primary = #6366F1`），不再需要 `text-brand`。

## Risks / Trade-offs

- **[视觉回归风险]** 大量组件的样式属性同时变更，任何映射错误都会导致视觉差异 → 逐组件迁移，每次迁移一个文件后肉眼对比确认；所有颜色值不变，只是引用方式改变
- **[Tailwind 类名冲突]** `text-primary` 在 shadcn 体系中指 `--primary`（品牌色），不是"主要文字色"→ 文档中明确 shadcn 命名约定：`text-foreground` = 主要文字，`text-primary` = 品牌色
- **[内联样式中复杂表达式]** 部分组件有条件拼接的内联样式，不能简单替换 → 保留动态部分的内联样式，仅替换静态部分
