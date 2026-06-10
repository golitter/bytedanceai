## Context

AgentHub 前端基于 React 19 + Tailwind CSS 4 + shadcn/ui 构建，设计系统定义在 `visual-style-guide.md`（Dark Utilitarian 风格），开发实践定义在 `development-strategy.md`。经过多轮功能迭代，代码中累积了 92 处风格违规和代码规范偏差。本次变更是一次纯前端的代码审计与修正，不涉及后端或协议变更。

当前 CSS 变量体系已在 `index.css` 中正确建立（`--bg-canvas` 到 `--bg-active` 五级灰阶、Agent 标识色、语义色等），但组件层存在大量绕过变量直接硬编码的情况。

## Goals / Non-Goals

**Goals:**
- 消除全部 30 处视觉规范违规（硬编码颜色、圆角超限、阴影滥用、strokeWidth 不统一）
- 统一 18 处 className 拼接方式为 `cn()` 调用
- 确保 Lucide 图标全局使用 `strokeWidth={1.25}`
- 标记 Zustand chat store 中 Server State 的技术债务

**Non-Goals:**
- 不重构 Zustand → TanStack Query 的状态管理层（涉及 SSE streaming 架构，是独立的重构任务）
- 不新增功能或组件
- 不修改 `index.css` 的 CSS 变量定义（当前已与 visual-style-guide 对齐）
- 不处理 admin 页面的业务逻辑

## Decisions

### D1: 颜色修复策略 — 使用 CSS 变量而非 Tailwind 语义类

**选择**: inline style 中的 `#fff` → `var(--primary-foreground)`；Tailwind 类中的 `text-white` → `text-primary-foreground`

**理由**: inline style 无法使用 Tailwind 类，必须用 CSS 变量。对于 Tailwind 类的场景，使用 `text-primary-foreground` 语义类保持与 shadcn 主题系统一致。

**替代方案**: 全部改用 inline style 的 CSS 变量 — 被否决，因为 Tailwind 类更符合项目约定，且 `text-white` 改为 `text-primary-foreground` 已足够语义化。

### D2: cn() 迁移策略 — 逐文件替换

**选择**: 将所有 `className={\`...${cond}\`}` 替换为 `cn('...', cond && '...')` 形式

**理由**: cn() 来自 `lib/utils.ts`（clsx + tailwind-merge），已是项目标准做法。模板字面量无法处理 Tailwind 类冲突（如条件覆盖 `bg-red-500` 时不会自动清除 `bg-blue-500`）。

**替代方案**: 保留模板字面量仅用于简单条件 — 被否决，统一使用 cn() 降低认知负担。

### D3: 圆角修正映射

**选择**: `rounded-xl`（12px）根据上下文修正：
- Agent 头像区域 → `rounded-lg`（8px），符合头像规范
- 覆盖层/overlay → 跟随父元素圆角

**理由**: visual-style-guide 明确规定最大圆角 12px 仅用于面板，头像用 8px。当前 `rounded-xl` 在 Tailwind 默认配置下为 12px（与 `--radius-xl` 一致），但语义上属于头像/图片区域，应降级。

### D4: 阴影清理策略

**选择**:
- `components/ui/dialog.tsx` 的 `shadow-lg` → 保留（Dialog 属于弹出菜单/模态框，是允许阴影的场景）
- `components/layout/IconSidebar.tsx` 的 `shadow-lg`（tooltip）→ 保留（tooltip 属于弹出菜单类）

**理由**: visual-style-guide 明确"唯一允许阴影的场景——弹出菜单/下拉框"。Dialog 和 Tooltip 都属于此类别。

### D5: strokeWidth 统一策略

**选择**: 全部 Lucide 图标统一为 `strokeWidth={1.25}`

**修正清单**:
- `AskAgentCard.tsx`: 1.5/1.6 → 1.25（5 处）
- `MessageBubble.tsx`: 1.5 → 1.25（1 处）
- `FinalSummaryCard.tsx`: 1.7 → 1.25（3 处）
- `TaskFailureCard.tsx`: 1.7 → 1.25（1 处）

**理由**: visual-style-guide 规定 `strokeWidth={1.25}` 模拟细线视觉风格，与 Phosphor Light 风格一致。

### D6: Zustand Server State 技术债务标记

**选择**: 在 chat store 文件头部添加注释标记技术债务，不在本次变更中重构

**理由**: chat store 中的 messages/streaming 使用 Zustand 是因 SSE streaming 架构需要高频更新（rAF 批量刷新）。迁移到 TanStack Query 需要重新设计 streaming 数据流，是独立的大重构。当前方案虽然违反 development-strategy 的状态分类原则，但功能正确且性能合理。

## Risks / Trade-offs

- **[视觉回归风险]** → 批量修改颜色/圆角/阴影可能引入微小的视觉差异。缓解：逐个检查每处修改的前后效果，确保语义一致。
- **[cn() 性能]** → 18 处新增 cn() 调用增加微小运行时开销。缓解：cn() 内部仅是 clsx + tailwind-merge，开销可忽略。
- **[shadow-lg 保留]** → 保留 Dialog/Tooltip 的 shadow-lg 可能被视为不一致。缓解：这符合设计规范中的例外条款。
- **[Zustand 债务未解决]** → messages 仍在 Zustand 中，不解决根本的状态管理违规。缓解：明确标记为技术债务，后续独立任务处理。
