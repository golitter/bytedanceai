## Context

前一轮审计（2026-05-30）修正了 rounded-full 滥用、diff 字体、Lucide strokeWidth 等表层问题。本轮审计覆盖剩余系统性偏差，分为两个维度：

**视觉规范偏差**：Tailwind 默认 radius scale（`--radius: 0.625rem`）导致 `rounded-md`=8px、`rounded-lg`=10px、`rounded-xl`=14px，与 visual-style-guide.md 规定的精确值（按钮 6px、输入框 8px、面板 12px）不匹配。大量组件使用 opacity 调色（`text-foreground/50`）绕过定义色板。`transition-colors` 广泛用于 hover 动画，违反了"仅动画 transform/opacity"规则。

**代码质量偏差**：`stores/chat.ts` 660 行严重超限，将 server state（消息历史）直接存入 Zustand 而非 TanStack Query。API 层 5 个核心函数缺少 `res.ok` 校验。全应用无 Error Boundary，单点故障导致白屏。魔法字符串（`'claude-code'`、`'streaming'`）散布在 4+ 个文件中。

## Goals / Non-Goals

**Goals:**
- 将全部组件 border-radius 对齐到 visual-style-guide.md 精确值
- 消除 opacity 调色，统一使用 CSS 变量语义 token
- 品牌色仅用于合规场景（选中态、发送按钮、焦点环）
- 过渡动画仅限 transform/opacity 属性
- 为关键模块补充 Error Boundary
- API 层补充 res.ok 校验

**Non-Goals:**
- 不重构 chat store 的状态管理架构（server state → TanStack Query 迁移是独立的架构变更，单独规划）
- 不修改组件拆分策略（前轮已完成 NewChatDialog、MessageList 拆分）
- 不引入新的 UI 组件或功能
- 不修改 Tailwind 配置中的 `--radius` 基础值（使用精确值 `rounded-[Npx]` 覆盖）

## Decisions

### D1: 使用精确像素值覆盖 Tailwind radius scale

**选择**：在组件中使用 `rounded-[6px]`、`rounded-[8px]`、`rounded-[12px]`、`rounded-[9999px]` 而非 `rounded-md`/`rounded-lg`/`rounded-xl`。

**理由**：Tailwind 的 radius scale 基于 `--radius` 基础变量（0.625rem = 10px），生成的 `rounded-md`=8px、`rounded-lg`=10px、`rounded-xl`=14px 与规范不匹配。修改基础变量会影响 shadcn/ui 组件的全局行为。精确像素值既满足规范又不影响 shadcn 生态。

**替代方案**：修改 `--radius` 为 0.5rem 并调整 scale — 风险高，影响所有 shadcn 组件默认样式。

### D2: 定义语义色 token 映射表替换 opacity 调色

**选择**：在 `index.css` 中补充 `--color-muted-text` 等中间色 token，将 `text-foreground/50` → `text-secondary`（`#8B91A0`）、`text-foreground/75` → `text-secondary`、`text-foreground/55` → `text-tertiary`（`#5A6070`）。

**理由**：opacity 调色在不同背景上表现不一致（半透明白色在深色和浅色背景上效果不同），且产生的颜色不在定义色板内。语义 token 保证一致性。

**替代方案**：直接用 `style` 属性写精确 hex — 不如 Tailwind token 一致。

### D3: 用 `transition-[transform,opacity]` 替换 `transition-colors`

**选择**：全局替换 `transition-colors` 为 `transition-[transform,opacity]` 或直接移除（如果元素无 transform/opacity 变化）。

**理由**：visual-style-guide.md 明确"只动画 transform 和 opacity"。`transition-colors` 动画 `color`、`background-color`、`border-color` 均不在允许范围。120-300ms ease-out 时长保持不变。

**替代方案**：保留 `transition-colors` 但限制 duration — 违反规范。

### D4: Error Boundary 采用模块级粒度

**选择**：在 `ImPage.tsx` 中为 ChatArea、AdminContent、ConversationList 分别包裹独立 Error Boundary。每个 Boundary 显示轻量降级 UI（错误提示 + 重试按钮）。

**理由**：模块级粒度保证一个模块崩溃不影响其他模块正常使用。错误信息提供重试入口而非空白页面。

**替代方案**：单个顶层 Error Boundary — 粒度太粗，一个模块出错全页面降级。

### D5: API 层统一错误处理模式

**选择**：在 `lib/api.ts` 中抽取 `handleResponse<T>(res: Response): Promise<T>` 辅助函数，统一检查 `res.ok` 并抛出结构化错误。所有 API 函数改用此辅助函数。

**理由**：5 个核心 API 函数缺少 res.ok 检查，非 2xx 响应导致晦涩的 TypeError。统一辅助函数消除重复代码，错误消息可追踪。

**替代方案**：每个函数独立检查 — 重复代码，容易遗漏。

## Risks / Trade-offs

- **[精确像素值冗长]** → `rounded-[6px]` 不如 `rounded-md` 简洁，但在可读性和合规性之间权衡，合规性优先
- **[移除 transition-colors 影响 hover 体验]** → hover 态背景色变化不再是平滑动画，而是瞬时切换。可通过在 hover 时同时使用 `opacity` 过渡来部分补偿
- **[Error Boundary 增加组件嵌套]** → 每多一个 Boundary 多一层 wrapper，但模块级粒度下影响有限（3 个 Boundary）
- **[API handleResponse 改变错误类型]** → 从 TypeError 变为自定义 ApiError，调用方 catch 逻辑可能需要调整。需同步更新调用方
