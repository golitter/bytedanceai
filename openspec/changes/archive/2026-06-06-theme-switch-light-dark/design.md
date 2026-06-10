## Context

当前前端是纯 Dark 主题。所有自定义语义 CSS 变量（`--bg-canvas`、`--text-primary`、`--agent-*` 等）仅在 `.dark` 块中定义，`:root` 只有 shadcn 默认 token（均为无彩色 oklch）。`IconSidebar` 已有 Settings 导航项但标记为 `disabled`。

Light 主题配色方案已在 `frontend/docs/reference/color-palette.md` L-八 中设计完毕，`frontend/payloads/light-theme-demo.html` 提供了完整视觉验证。

## Goals / Non-Goals

**Goals:**

- 用户可在 Settings 面板中切换 Dark / Light 主题，切换即时生效
- 主题偏好持久化到 `localStorage`，刷新页面后保持
- 首次访问时自动跟随系统 `prefers-color-scheme`
- Light 模式下所有组件颜色正确，与 Dark 模式功能对等
- Settings 入口复用 IconSidebar 已有的 Settings 图标（移除 `disabled`）

**Non-Goals:**

- 不实现更多主题色方案（仅 Dark / Light 两档）
- 不实现 per-session 或 per-agent 主题
- 不涉及后端 API 变更（纯前端功能）
- 不引入新的 UI 框架或 CSS-in-JS 方案（保持 Tailwind + CSS 变量体系）

## Decisions

### D1: 主题切换机制 — `<html>` class 切换

**选择**：在 `<html>` 元素上添加/移除 `dark` class，通过 CSS 变量级联生效。

**理由**：项目已使用 `@custom-variant dark (&:is(.dark *))` 做暗色变体匹配，这正是 shadcn/ui 的标准做法。无需引入额外库或 Context Provider 即可让所有 Tailwind dark: 变体和 CSS 变量自动响应。

**替代方案**：
- CSS `prefers-color-scheme` 媒体查询 → 无法响应用户手动切换
- CSS 变量全部 runtime 切换 → 需要同时维护两套值，复杂度高
- Context Provider + CSS 变量 → 额外层级，与 Tailwind dark: 变体不兼容

### D2: `useTheme` hook — 轻量级实现

**选择**：独立 hook `src/hooks/use-theme.ts`，不引入 Zustand store。

**理由**：主题状态是全局单一值，不需要 middleware、devtools 或复杂更新逻辑。hook 内部使用 `useState` + `localStorage` + `matchMedia` 即可。如果组件需要响应主题变化，直接调用 hook 读取 `resolvedTheme`。

**接口设计**：

```ts
type Theme = 'dark' | 'light' | 'system'
interface UseThemeReturn {
  theme: Theme          // 用户设置的偏好（dark/light/system）
  resolvedTheme: 'dark' | 'light'  // 实际生效的主题
  setTheme: (theme: Theme) => void
}
```

### D3: Settings 面板 — Popover 弹出

**选择**：复用项目已有的 shadcn Popover 组件，点击 Settings 图标弹出面板。

**理由**：Popover 已在 `components/ui/popover.tsx` 中就绪（基于 Radix），无需新建组件。面板内容为简单的主题切换，无需 Dialog 或独立页面。

**布局**：Popover 从 Settings 图标向右弹出，包含：
- 标题 "外观"（12px, text-secondary）
- 三个选项按钮：深色 / 浅色 / 跟随系统（类似 Linear 的三段式切换）
- 当前选中项高亮（bg-primary-soft + text-primary）

### D4: Light 主题 CSS 变量 — `:root` 补全

**选择**：在 `:root` 块中补全所有 `.dark` 块中存在的自定义变量，色值参考 `color-palette.md` L-八。

**关键决策**：
- `--primary` 覆盖为 `#6366F1`（与 Dark 一致的 Indigo，替代 shadcn 默认的 oklch 无彩黑）
- `--ring` 覆盖为 `#6366F1`（品牌色焦点环）
- Agent 标识色 `--agent-*` 两套主题保持一致（这些是功能色，不是装饰色）
- 图表色 `--chart-*` 保留 shadcn 默认灰阶（Light 下图表需求待定）

### D5: 初始化时序

**选择**：在 `main.tsx` 中 `createRoot` 之前执行同步初始化脚本（inline script 或模块顶层调用），避免 FOUC（Flash of Unstyled Content）。

**理由**：如果等 React 渲染后再添加 `dark` class，用户会看到短暂的白屏闪烁。在 `<script>` 标签中同步读取 `localStorage` 并设置 class 可以完全避免 FOUC。

**实现**：在 `index.html` 的 `<head>` 中添加内联 script，读取 localStorage 并立即设置 class。`useTheme` hook 则负责后续的响应式切换。

### D6: 内联暗色样式适配

少量组件使用内联 `rgba(0,0,0,...)` 做 boxShadow（如 `IconSidebar` 的 hover 卡片 `boxShadow: '0 4px 24px rgba(0,0,0,0.4)'`）。Light 模式下需改为动态值。

**选择**：将这些内联值改为 CSS 变量引用（在 `:root` 和 `.dark` 中定义 `--shadow-popup`），或使用 Tailwind 的 `shadow-lg` 等语义类。

## Risks / Trade-offs

- **[FOUC 风险]** → 通过 `index.html` 内联 script 同步初始化 class 规避
- **[Light 模式视觉差异]** → 某些低透明度背景在白底上可能几乎不可见（如 `rgba(99,102,241,0.06)`），需在实际页面中微调透明度
- **[内联样式遗漏]** → 部分组件可能有遗漏的硬编码暗色值，需人工走查一遍
- **[测试覆盖]** → 纯 CSS/视觉变更难以自动化测试，依赖人工视觉验证 + Light demo HTML 对照
