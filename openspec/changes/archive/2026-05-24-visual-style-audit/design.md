## Context

当前前端 IM 聊天界面（commit 1bbb088）已实现基本功能，但视觉样式与 `frontend/docs/visual-style-guide.md` 规范存在 8 处偏差。本变更的目标是对齐实现与规范，不涉及功能变更。

当前状态：
- Agent 标识色使用了品牌色变体，而非规范指定的暖色系
- AgentAvatar 缺少规范要求的同色模糊光晕
- 状态灯动效与规范不一致
- 使用了 Lucide 图标而非规范指定的 Phosphor Icons
- 部分组件的间距、背景色、过渡时长与规范不符

## Goals / Non-Goals

**Goals:**
- 将所有组件视觉样式对齐到 visual-style-guide.md
- 修正 CSS 变量值
- 评估图标库迁移的可行性

**Non-Goals:**
- 不改变组件结构或交互逻辑
- 不新增功能
- 不修改后端 API 或数据结构

## Decisions

### D1: Agent 标识色修正

直接替换 `index.css` 中的三个 CSS 变量值：
- `--agent-claude`: `#6366F1` → `#DA7756`
- `--agent-opencode`: `#F59E0B` → `#10B981`
- `--agent-orchestrator`: `#22C55E` → `#EAB308`

所有引用这些变量的组件（AgentAvatar、MessageBubble）自动生效，无需修改组件代码。

### D2: AgentAvatar 光晕效果

在 AgentAvatar 组件中添加 `box-shadow`，使用 Agent 标识色 + 8px blur。由于 CSS 变量不能直接用于 `rgba()` 的 alpha 参数，需要在组件中维护一个颜色映射。

### D3: 状态灯动效修正

- `ready`：添加 opacity 脉冲动画（0.6→1→0.6，2s 循环）
- `running`：改为旋转动画（1.5s 循环），替代当前的 pulse

在 `index.css` 中定义 `@keyframes` 规则，AgentAvatar 通过 className 引用。

### D4: 图标库 — 暂不迁移

Phosphor Icons 与 Lucide 的 API 风格不同，全量迁移涉及 8+ 组件文件，风险较高。建议：
- 本期保持 Lucide，但将 `strokeWidth` 从 `1.5` 调整为 `1.25` 接近 Phosphor Light 的视觉感
- 后续单独做图标库迁移

### D5: 全局 letter-spacing

在 `index.css` 的 `body` 或 `html` 规则中添加 `letter-spacing: -0.01em`。代码块中覆盖为 `0`。

### D6: 细节对齐

| 项 | 当前 | 目标 | 修改位置 |
|------|------|------|------|
| 输入框背景 | `var(--bg-hover)` | `var(--bg-card)` | MessageInput.tsx |
| 列表项 padding | `px-3 py-2.5` | `px-3 py-2.5` | 实际为 12px 10px，符合规范 10px 12px 的倒序（Tailwind px/py 含义），无需改 |
| hover 过渡 | `transition-colors`（默认 150ms） | `transition-colors duration-120` | ConversationItem、ConversationList |
| 分割线颜色 | `rgba(255,255,255,0.06)` | `rgba(255,255,255,0.04)` | 分割线 border 区别于边框 |

## Risks / Trade-offs

- [Agent 色变更视觉冲击] → Claude 从紫色变暖橙、Orchestrator 从绿变金，用户需适应。缓解：颜色变更与规范一致，无个性化空间
- [暂不迁移图标库] → 与规范不完全一致。缓解：调低 strokeWidth 模拟细线风格，降低视觉差距
- [全局 letter-spacing 可能影响布局] → 间距收紧后部分文字可能换行。缓解：`-0.01em` 极小，实际影响可忽略
