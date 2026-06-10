## Context

AgentHub 前端遵循 visual-style-guide.md 定义的 "Dark Utilitarian" 视觉规范，使用 Lucide React 图标库（strokeWidth=1.25 细线风格）和 CSS 变量色彩体系。近期新增的 SkillsHub、ContactsPage、PlanReviewCard、AgentProfilePage 等模块在实现时未严格遵守规范，导致 8 处 strokeWidth 偏差、3 处 text-white 使用、多处 emoji 替代 Lucide 图标、以及直接色值（text-red-500）绕过语义 token。本次变更是一次纯前端样式层面的合规修正，不涉及功能变更。

## Goals / Non-Goals

**Goals:**
- 全站 Lucide 图标 strokeWidth 统一为 1.25，消除 1.5 和 "3" 的偏差
- 新模块 emoji 图标全部替换为 Lucide React 图标，保持视觉一致性
- 消除 `text-white`，改用 `text-primary-foreground` 语义类
- 消除 `text-red-500` 等直接色值，改用 `text-destructive` 语义类
- 消除 inline style 中的 `var(--text-secondary)`，改用 `text-text-secondary` Tailwind 类

**Non-Goals:**
- 不重构大文件结构（message-store.ts、api.ts 等超 300 行文件的拆分不属于本次范围）
- 不调整组件架构或状态管理
- 不新增 ESLint 自定义规则（后续可独立规划）
- 不修改 visual-style-guide.md 本身

## Decisions

### D1: strokeWidth 修正策略 — 全局搜索替换
所有 Lucide 图标的 `strokeWidth={1.5}` 和 `strokeWidth="3"` 直接替换为 `strokeWidth={1.25}`。无需引入全局默认值机制（如 wrapper 组件），因为违规点仅 8 处，逐一修改更可控。

**替代方案（不采用）**：创建一个 LucideIconWrapper 统一设置 strokeWidth。缺点：增加抽象层，违反 development-strategy.md 的"第三次重复再抽象"原则。

### D2: Emoji → Lucide 映射
SkillsHubPage 等处的 emoji 需要映射到语义等价的 Lucide 图标：
- 🛡️ (内置标记) → `Shield` 图标
- 📦 (外部/上传标记) → `Package` 图标
- ⚙️ (内置 Skill 头像) → `Settings` 或 `Wrench` 图标
- ⚠️ (警告) → `AlertTriangle` 图标
- ✅ (成功) → `CheckCircle2` 图标
- ❌ (失败) → `XCircle` 图标

映射遵循"语义优先"原则，选最贴近 emoji 含义的 Lucide 图标。

### D3: text-white 替换目标
`text-white` 在 `bg-primary` 按钮上应替换为 `text-primary-foreground`，这是 shadcn 语义 token，在暗色主题下等价于白色但在亮色主题下会自动适配。保持按钮的文字可读性。

### D4: 直接色值 → 语义 token
- `text-red-500` → `text-destructive`（shadcn 语义 token，暗色主题下同为红色）
- `border-red-500/20` → `border-destructive/20`
- `bg-red-500/10` → `bg-destructive/10`

## Risks / Trade-offs

- **[视觉回归风险]** 替换 emoji 为 Lucide 图标可能改变用户对新功能的视觉认知 → 缓解：emoji 仅用于 SkillsHub 页面的装饰性位置（SectionLabel、HubSkillCard 图标），替换为 Lucide 图标后视觉更统一，不会影响功能认知
- **[遗漏风险]** 审计可能未覆盖所有违规点 → 缓解：在 tasks 中包含全量验证步骤，修改后用 grep 确认零残留
- **[极低风险]** `text-primary-foreground` 在亮色主题下可能不是白色 → 缓解：项目当前仅使用暗色主题，且这正是语义 token 的正确用法
