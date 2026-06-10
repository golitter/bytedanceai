## Context

前端在多 Agent 运行时快速迭代中，视觉规范执行出现偏差。当前代码与 `visual-style-guide.md`（Dark Utilitarian 风格）和 `development-strategy.md`（工程策略）存在以下差距：

- AgentAvatar 使用 `box-shadow` 发光效果，违反"几乎不用阴影"原则
- 多处使用 `text-white` 而非 CSS 变量 `--text-primary`
- PlanCard 中 Agent 标识色被用作 UI 框架色，超出"仅限头像+消息色条"边界
- GroupAvatar 字号 8px/9px 不在定义体系内
- 状态灯动画使用 `ease-in-out` 而非规定的 `ease-out`
- 部分组件内联样式可用 Tailwind 工具类替代

## Goals / Non-Goals

**Goals:**
- 将所有组件样式对齐到 visual-style-guide.md 规范
- 消除纯白文字、越界阴影、错误缓动、超范围字号
- Agent 标识色严格限制在头像 + 消息色条场景
- 内联样式尽可能迁移到 Tailwind 类
- 确保修改不引入功能回归

**Non-Goals:**
- 不重构组件结构或状态管理
- 不修改 API 契约或后端逻辑
- 不新增组件或功能特性
- 不调整 CSS 变量体系（index.css 已与规范对齐）

## Decisions

### D1: text-white → text-primary-foreground

AgentAvatar 和 GroupAvatar 的头像文字当前用 `text-white`（纯白 #FFFFFF）。改为已有的 `text-primary-foreground` CSS 变量（对应 `#E8EBF0`）。

理由：规范明确禁止纯白文字。`text-primary-foreground` 已在 index.css 定义。

### D2: 移除 AgentAvatar box-shadow 发光

AgentAvatar 当前使用 `boxShadow: 0 0 8px ${shadowColor}` 模拟发光光晕。规范中"唯一允许的发光效果，仅头像使用"——但这条与"几乎不用阴影"的核心原则冲突。按更严格的解读，发光效果应通过背景色差表达（用更深的同色背景圈），不用 box-shadow。

替代方案：使用 `background: ${color}20`（同色低透明度背景圈）代替发光阴影。

### D3: PlanCard agent badge 改用语义色

PlanCard 中 task agent badge 直接使用 `AGENT_COLORS[agent]` 作为文字色和背景色。规范要求 Agent 标识色仅用于头像和消息色条。

修改方案：badge 使用 `text-secondary` 文字色 + `bg-card` 背景，不再使用 Agent 标识色。Agent 身份通过左侧消息色条（已有）区分。

### D4: 动画缓动统一 ease-out

AgentAvatar 状态灯 ready 脉冲使用 `ease-in-out`，改为 `ease-out`。

### D5: GroupAvatar 字号对齐

GroupAvatar 中 `text-[9px]` 和 `text-[8px]` 分别改为 `text-[11px]`（微标注）和 `text-[10px]`（icon 内数字，无对应档位，保留但提升到 10px）。

### D6: 内联样式 Tailwind 化

AgentAvatar 中 width/height/borderRadius 等内联样式迁移到 Tailwind 类，通过动态 class 拼接实现。仅在 Tailwind 无法覆盖时保留内联（如动态颜色值）。

## Risks / Trade-offs

- **视觉微调用户可感知** → 改动均为微调（颜色微调、阴影移除），不涉及布局变化，影响极小
- **AgentAvatar 发光移除可能降低头像辨识度** → 通过同色背景圈 + 状态灯补偿视觉区分
- **PlanCard badge 不再显示 Agent 颜色** → 左侧色条已提供 Agent 身份区分，badge 改为中性色不影响可读性
