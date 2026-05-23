# 主题与样式系统

## 设计原则

- 暗色优先：主设计稿为 `.dark` 模式，浅色模式保留 CSS 变量但未专门调优
- CSS 变量驱动：所有颜色通过 CSS 自定义属性控制，组件中不硬编码颜色值
- Tailwind + 内联 style 混合：布局用 Tailwind class，颜色用 `style={{ color: 'var(--xxx)' }}`

## 色彩体系

### 基础色（暗色模式）

| 变量 | 值 | 用途 |
|------|------|------|
| `--bg-canvas` | `#0A0B0E` | 主画布背景 |
| `--bg-sidebar` | `#111318` | 侧栏背景 |
| `--bg-hover` | `#22262F` | hover / 搜索框背景 |
| `--bg-active` | `#2C313B` | 激活态背景（预留） |
| `--card` | `#1A1D24` | 卡片 / Agent 消息气泡背景 |
| `--popover` | `#22262F` | 弹出层背景 |

### 文本色

| 变量 | 值 | 用途 |
|------|------|------|
| `--text-primary` | `#E8EBF0` | 主文本 |
| `--text-secondary` | `#8B91A0` | 次要文本（时间、描述等） |
| `--text-tertiary` | `#5A6070` | 占位符、辅助信息 |

### 功能色

| 变量 | 值 | 用途 |
|------|------|------|
| `--color-brand` | `#6366F1` | 品牌 / 主色调（Indigo） |
| `--color-success` | `#22C55E` | 成功 / 就绪 |
| `--color-warning` | `#F59E0B` | 警告 / 运行中 |
| `--color-error` | `#EF4444` | 错误 |
| `--border` | `rgba(255,255,255,0.06)` | 统一边框色 |

### Agent 专属色

| 变量 | 值 | Agent |
|------|------|-------|
| `--agent-claude` | `#6366F1` | Claude Code（紫色） |
| `--agent-opencode` | `#F59E0B` | OpenCode（琥珀色） |
| `--agent-orchestrator` | `#22C55E` | Orchestrator（绿色） |

这些颜色同时用于：
- `AgentAvatar` 背景色
- `MessageBubble` 左侧竖线颜色
- 状态指示灯

## 状态指示灯颜色

| 状态 | 颜色 | 说明 |
|------|------|------|
| `ready` | `#22C55E` 绿色 | Agent 在线空闲 |
| `running` | `#F59E0B` 琥珀色 + 脉冲动画 | Agent 正在处理 |
| `offline` | `#5A6070` 灰色 | Agent 离线 |
| `error` | `#EF4444` 红色 | 出错 |

## 字体

```css
--font-sans: 'Geist Variable', sans-serif;
```

代码字体：`'Geist Mono', monospace`（在 CodeBlock 和行内代码中硬编码）。

通过 `@fontsource-variable/geist` 引入，无需外部 CDN。

## 边框策略

统一使用 `rgba(255,255,255,0.06)` 作为暗色模式边框色，在深色背景上提供微妙的分割线效果，不显突兀。

组件中的写法：
```tsx
style={{ borderColor: 'rgba(255,255,255,0.06)' }}
```

## 交互状态

| 状态 | 实现 |
|------|------|
| Hover | `onMouseEnter/Leave` 切换 `backgroundColor: 'var(--bg-hover)'` |
| Active（选中对话） | `borderLeft: 2px solid var(--color-brand)'` + 持久 hover 背景 |
| Disabled | Tailwind `disabled:opacity-40` 或 `disabled:opacity-50` |
| Streaming 光标 | CSS `animate-pulse` + `▌` 字符 |

## 代码块样式

- 背景：`#0D0F14`（比 canvas 更深一层）
- 主题：Shiki `tokyo-night`
- 字体：`Geist Mono`
- 字号：`13px`
- 行高：`1.65`
- 圆角：`rounded-lg`
