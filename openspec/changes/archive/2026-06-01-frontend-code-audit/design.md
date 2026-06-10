## Context

AgentHub 前端基于 React 19 + Vite 8 + TypeScript + Tailwind CSS 4 + shadcn/ui + Zustand 5 + TanStack Query 5 技术栈。项目已建立完整的视觉规范（visual-style-guide.md）和开发策略（development-strategy.md），并已有 10+ 份 OpenSpec 规范定义了色彩、圆角、过渡、组件拆分等规则。

**当前状态**：CSS 变量体系已完备（OKLCH 色彩空间），shadcn/ui 基础组件已安装，Geist 字体和 Lucide 图标已集成。但代码审计发现规范与实际实现之间存在显著脱节：

| 类别 | 违规数 | 严重度 |
|------|--------|--------|
| `transition-colors`（应只用 transform/opacity） | ~40 处 | 系统性 |
| 内联 style 代替 Tailwind 类 | 120+ 处（admin 页面为主） | 系统性 |
| Tailwind 直接色类（text-red-500 等） | ~30 处 | 系统性 |
| 硬编码 rgba()/hex 色值 | 6 处 | 中等 |
| 硬编码 Agent 类型字符串 | ~20 处 | 中等 |
| text-[10px]（不在字号体系中） | 6 处 | 轻微 |
| 硬编码用户名 | 4 处 | 轻微 |
| chat.ts 超 974 行（建议 200 行） | 1 处 | 架构级 |
| Admin 页面用 useState+useEffect 代替 TanStack Query | 7 页面 | 架构级 |

**约束**：
- 仅限前端 `frontend/src/` 目录
- 不引入新依赖
- 不改变功能行为，只做规范化和重构
- 需参考已有 skills（design-taste-frontend, minimalist-ui 等）的设计理念

## Goals / Non-Goals

**Goals:**
- 将所有硬编码色值迁移到 CSS 变量 / Tailwind 语义类
- 将 `transition-colors` / `transition-all` 替换为 `transition-[transform,opacity]`
- 将 admin 页面的内联 style 迁移到 Tailwind 工具类
- 将散落的魔法字符串收敛到 `constants.ts`
- 修正 6 处 `text-[10px]` 为允许的字号
- 修正 4 处硬编码用户名
- 将 `chat.ts` 按领域拆分为多个 store
- 将 7 个 admin 页面从 `useState+useEffect` 迁移到 TanStack Query

**Non-Goals:**
- 不添加新的 UI 功能或交互
- 不修改后端 API 或数据模型
- 不引入新的第三方库
- 不做性能优化（虚拟列表、React.lazy 等）— 属于独立优化任务
- 不修改 `generated/` 目录下的契约生成代码
- 不修改测试文件（本次审计不涉及测试覆盖）

## Decisions

### D1. 色值迁移策略：Tailwind 语义类优先

**决策**：对于 Tailwind 直接色类（`text-red-500`、`bg-green-500/10`），优先替换为已有的 CSS 变量 Tailwind 语义类（`text-destructive`、`bg-success/10`、`text-foreground`）。对于不存在精确对应的色值，在 `index.css` 中补充 CSS 变量。

**替代方案**：用 `style={{ color: 'var(--color-error)' }}` — 但违反 Dev Strategy 4.1「优先用工具类」原则。

**理由**：Tailwind 工具类是项目的首选样式方案，内联 style 仅用于动态计算值。

### D2. 过渡替换策略：`transition-[transform,opacity]`

**决策**：将 `transition-colors` 统一替换为 `transition-[transform,opacity]`。由于 hover 态背景色变化不再走 transition，视觉上改为即时切换（120ms 的 `transition-colors` 本身违反 VSG 8 规范）。

**替代方案**：保留 `transition-colors` 但用 `120ms` — 仍违反规范（只允许 transform/opacity）。

**理由**：VSG 第 8 章明确禁止 `transition-colors`，规范级别高于视觉偏好。

### D3. Admin 页面数据获取：渐进迁移到 TanStack Query

**决策**：为每个 admin 页面创建对应的 API 函数和 `useQuery` hook，替换 `useState + useEffect + fetch` 模式。集中定义 admin API 到 `lib/api.ts` 已有的 admin 模块中。

**替代方案**：暂不迁移，标注 TECH DEBT — 但已有规范明确要求 Server State 走 TanStack Query。

**理由**：消除 7 个页面的手写 loading/error 逻辑，统一缓存和刷新行为。

### D4. chat.ts 拆分：按领域域拆为 3-4 个 store

**决策**：将 974 行的 `chat.ts` 拆分为：
- `navigation-store.ts`（导航状态）
- `session-store.ts`（会话列表、切换）
- `message-store.ts`（消息、streaming、runtime blocks）

**替代方案**：保持单文件但加分区注释 — 文件已超 900 行，无法维护。

**理由**：Dev Strategy 5.4 明确「Store 按领域拆分，不合并成一个大 store」。

### D5. 魔法字符串：收敛到 constants.ts + 生成类型

**决策**：在 `lib/constants.ts` 中补充 `AGENT_TYPES`、`MESSAGE_ROLES`、`CHAT_STATUSES` 常量，组件中引用常量代替硬编码字符串。

**理由**：Dev Strategy 第九章「消灭魔法字符串」。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| 色值替换可能导致视觉回归 | 每个 task 完成后截图对比，确认无视觉差异 |
| chat.ts 拆分可能破坏组件间的状态依赖 | 拆分时保持 store 接口不变，组件无需改动 |
| admin 页面迁移到 TanStack Query 可能引入缓存问题 | admin 数据是低频读取，设 `staleTime: 30s` 即可 |
| transition 移除 color 过渡可能让 hover 感觉"生硬" | 这是 VSG 规范要求，保持一致性优先 |
| 大规模修改可能引入回归 | 按 task 粒度提交，每个 task 可独立回滚 |
