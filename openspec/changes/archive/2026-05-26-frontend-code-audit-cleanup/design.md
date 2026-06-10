## Context

前端经多轮 spec 迭代（css-variable-audit、style-normalization、code-dedup、frontend-constants 等），CSS 变量体系和基础常量已建立。但实际代码审计发现以下残留问题：

1. **常量散落**：`API_BASE = '/api'` 在 `api.ts`、`DiffCard.tsx`、`AttachmentCard.tsx`、`ImageCard.tsx` 四处独立声明；Agent 颜色映射在 `AgentAvatar.tsx`（`AGENT_COLORS` + `AGENT_SHADOW_COLORS`）和 `MessageBubble.tsx`（`AGENT_STRIP_COLORS`）中各定义一份，且两者完全相同。
2. **重复工具逻辑**：文件名提取 `path.split('/').pop()` 在 `AttachmentCard.tsx`、`ImageCard.tsx`、`diff-parser.ts`、`block-reducer.ts` 中各出现一次。
3. **魔法数字**：`MessageInput` 中 `maxHeight: 200`、hint 超时 `3000`；`MessageList` 中虚拟化阈值 `50`、滚动判定 `60`；`AgentAvatar` 中动画时长 `'2s'`、`'1.5s'`。
4. **DiffCard 过大**：337 行，混合了快照 CRUD、视图切换、文件 Tab 管理、操作按钮四类职责。
5. **CSS diff 主题硬编码**：`index.css` 中 diff-gutter/diff-code 的 success/error 色值直接写 hex/rgba，未引用 `--color-success`/`--color-danger-bg` 变量。
6. **CHANGE_TYPE_MAP 重复**：DiffCard.tsx 和 DiffFileTabs.tsx 各自维护了几乎相同的变更类型→颜色映射。

当前约束：
- 不引入新依赖
- 纯前端重构，不改后端 API
- 保持所有现有行为不变

## Goals / Non-Goals

**Goals:**
- 消除 `API_BASE` 四处重复，统一到 `lib/constants.ts`
- 收敛 Agent 颜色映射为一份定义，所有组件引用同一来源
- 提取 `getFileName` 到 `lib/utils.ts`，消除 4 处 `path.split('/').pop()` 重复
- 提取魔法数字为命名常量
- 将 DiffCard.tsx 按职责拆分为 <200 行的子组件
- 将 diff 主题 CSS 中的硬编码色值替换为 CSS 变量

**Non-Goals:**
- 不重构状态管理（TanStack Query / Zustand 架构不变）
- 不修改后端接口或契约层
- 不引入新 UI 组件库或新依赖
- 不做视觉风格变更（颜色/间距/圆角等不变）

## Decisions

### D1: API_BASE 统一到 lib/constants.ts

**选择**：在 `lib/constants.ts` 中 export `API_BASE`，所有消费方 import 它。

**理由**：`lib/constants.ts` 已有 `AGENT_NAMES` 等项目级常量，API 基础路径属于同类。`api.ts` 中的 `const API_BASE` 改为 import。

**替代方案**：放在 `api.ts` 然后其他文件 import `api.ts` — 不好，因为 `api.ts` 是 API 函数文件，不应承担常量导出的职责。

### D2: Agent 颜色映射统一到 lib/constants.ts

**选择**：新增 `AGENT_COLORS: Record<AgentType, string>` 到 `lib/constants.ts`，值为 CSS 变量引用（`var(--agent-claude)` 等）。删除 `AgentAvatar.tsx` 中的 `AGENT_COLORS`/`AGENT_SHADOW_COLORS` 和 `MessageBubble.tsx` 中的 `AGENT_STRIP_COLORS`。

**理由**：三处映射完全相同（key = agentType, value = 对应 CSS 变量）。`AGENT_SHADOW_COLORS` 在 `AgentAvatar` 中使用，与 `AGENT_COLORS` 的值完全一致（都是同一 CSS 变量），直接复用 `AGENT_COLORS` 即可。

### D3: getFileName 工具函数提取到 lib/utils.ts

**选择**：在 `lib/utils.ts` 中新增 `getFileName(path: string): string`，逻辑为 `path.split('/').pop() || path`。`diff-parser.ts` 中已有的 `getFileName` 改为从 `utils.ts` import，`AttachmentCard.tsx`、`ImageCard.tsx` 中的内联 `path.split('/').pop()` 替换为函数调用。

**理由**：`lib/utils.ts` 已有 `cn()` 等工具函数，是通用工具的归属地。`diff-parser.ts` 的 `getFileName` 应该放在更通用的位置。

### D4: 魔法数字提取为命名常量

**选择**：在各自的组件文件顶部定义常量（不在 constants.ts 中，因为这些是组件内部实现细节）：

```typescript
// MessageInput.tsx
const MAX_INPUT_HEIGHT = 200
const HINT_DISPLAY_DURATION = 3000
const MIN_INPUT_HEIGHT = 48

// MessageList.tsx
const VIRTUALIZE_THRESHOLD = 50  // 已有，保持
const SCROLL_BOTTOM_THRESHOLD = 60

// AgentAvatar.tsx
const STATUS_READY_DURATION = '2s'
const STATUS_RUNNING_DURATION = '1.5s'
```

**理由**：development-strategy.md 原则"消灭魔法字符串"。这些值是组件内部调优参数，不适合放到全局 constants。

### D5: DiffCard 拆分策略

**选择**：按视觉区块拆分，不按数据流拆分。DiffCard 保留为编排组件（状态 + 数据获取），拆出：

- `DiffHeader` — 顶栏（文件统计 + 视图切换 + 操作按钮 + 状态 badge）
- `DiffFileInfo` — 文件信息条（路径 + 变更类型 badge + 增删统计）

DiffCard 变为 ~150 行，DiffHeader ~80 行，DiffFileInfo ~40 行。

**理由**：数据获取和快照管理逻辑耦合紧密，拆开反而增加 props drilling。按视觉区块拆分是 development-strategy.md 推荐的方式。

**替代方案**：拆出 `useDiffSnapshot` hook 管理快照 CRUD — 不选，因为只有一个消费方，抽 hook 是过度抽象（"抽象是负债"原则）。

### D6: CSS diff 主题变量化

**选择**：在 `.dark` 中新增 diff 专用 CSS 变量，替换硬编码色值：

```css
--diff-insert-color: var(--color-success);          /* #22C55E */
--diff-insert-bg: rgba(34, 197, 94, 0.08);
--diff-delete-color: var(--color-error);             /* #EF4444 */
--diff-delete-bg: rgba(239, 68, 68, 0.08);
--diff-insert-bg-strong: rgba(34, 197, 94, 0.1);
--diff-delete-bg-strong: rgba(239, 68, 68, 0.1);
```

**理由**：visual-style-guide.md 定义了 `--color-success`/`--color-error`，diff 的增删色应与之对应。透明度变体无法通过 CSS 变量计算（CSS 不支持从 hex 派生 rgba），因此保留为独立变量但命名清晰。

### D7: CHANGE_TYPE_MAP 统一

**选择**：将 `DiffFileTabs.tsx` 中的 `CHANGE_TYPE_MAP` 提取为 `DiffCard.tsx` 和 `DiffFileTabs.tsx` 共享的模块级导出，放在 `DiffFileTabs.tsx` 中 export。

**理由**：`DiffFileTabs` 是 DiffFileBadge 的自然归属。DiffCard 中文件信息条的 badge 逻辑应复用 `ChangeTypeBadge` 组件而非重复映射。

## Risks / Trade-offs

- **[重构风险]** DiffCard 拆分涉及大量移动代码 → 逐块移动 + 保持 DiffCard 功能不变，可逐步验证
- **[变量化副作用]** CSS 变量替换可能因 oklch/hex 色彩空间差异产生微小色差 → 使用相同的 hex 值，只是通过变量引用，无色差风险
- **[import 循环]** `constants.ts` 已 import `AgentType`，`api.ts` import `constants.ts` → 无循环依赖风险（单向依赖）
