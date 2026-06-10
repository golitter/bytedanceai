## ADDED Requirements

### Requirement: rounded-full 仅限头像和胶囊 Badge
所有使用 `rounded-full` 的元素 SHALL 仅限以下场景：用户头像（`border-radius: 50%`）、胶囊形 Badge/标签。状态指示器圆点 SHALL 使用固定尺寸（如 `w-1.5 h-1.5`）配合 `rounded-full`（4px 圆点本身是圆形，属合理用法）。非胶囊的按钮、卡片、面板 SHALL 使用规范圆角值（6/8/10/12px）。

#### Scenario: 非 Badge 元素不使用 rounded-full
- **WHEN** 检查所有组件文件的 `rounded-full` 用法
- **THEN** 仅在用户头像、胶囊 Badge/标签、4px 状态圆点中出现 `rounded-full`

#### Scenario: 进度条和状态指示器使用规范圆角
- **WHEN** 组件需要圆角元素（非头像非 Badge）
- **THEN** 使用 `rounded-sm`（6px）、`rounded`（8px）、`rounded-lg`（10px）等规范值

### Requirement: diff 代码区域使用等宽字体
`.diff-code` CSS 类 SHALL 使用 Geist Mono 等宽字体，与 `font-sans` UI 字体区分。

#### Scenario: diff 代码字体为 Geist Mono
- **WHEN** 渲染 diff 文件的代码区域
- **THEN** 字体为 `'Geist Mono', 'Geist Variable', monospace`，字号 13px

### Requirement: Lucide 图标统一 strokeWidth
所有 Lucide React 图标组件 SHALL 设置 `strokeWidth={1.25}` 以匹配细线风格。

#### Scenario: 全量 Lucide 图标具有 strokeWidth 属性
- **WHEN** 检查所有使用 lucide-react 图标的组件
- **THEN** 每个图标组件均包含 `strokeWidth={1.25}` 属性

### Requirement: Agent 头像同色模糊光晕
Agent 头像 SHALL 显示 8px 的同色模糊光晕（`box-shadow: 0 0 8px <agent-color>`），这是唯一允许的发光效果。

#### Scenario: Agent 头像显示光晕
- **WHEN** 渲染 Agent 头像组件
- **THEN** 头像外围有 8px 的同色 box-shadow 光晕

### Requirement: 大写 Badge 字间距
大写标签和 Badge 文字 SHALL 设置 `letter-spacing: +0.05em`，提升全大写文本可读性。

#### Scenario: Badge 组件包含字间距
- **WHEN** 渲染全大写或大写主导的 Badge/标签
- **THEN** 文字具有 `letter-spacing: 0.05em` 样式
