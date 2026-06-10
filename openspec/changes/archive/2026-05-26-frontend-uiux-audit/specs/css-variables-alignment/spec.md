## ADDED Requirements

### Requirement: 背景色阶变量定义
index.css SHALL 在 `:root` / `.dark` 中定义 visual-style-guide.md 规定的全部 5 级背景色阶变量：`--bg-canvas`(#0A0B0E)、`--bg-sidebar`(#111318)、`--bg-card`(#1A1D24)、`--bg-hover`(#22262F)、`--bg-active`(#2C313B)。

#### Scenario: 背景色阶变量完整
- **WHEN** 检查 index.css 中 CSS 自定义属性定义
- **THEN** 5 个背景色阶变量全部存在且色值与 visual-style-guide.md 一致

#### Scenario: 组件通过变量引用背景色
- **WHEN** 检查组件 Tailwind 类名中的背景色引用
- **THEN** 背景色通过 CSS 变量或对应的 Tailwind 语义 token 引用，不存在硬编码色值

### Requirement: 文字色阶变量定义
index.css SHALL 定义全部 3 级文字色阶变量：`--text-primary`(#E8EBF0)、`--text-secondary`(#8B91A0)、`--text-tertiary`(#5A6070)。

#### Scenario: 文字色阶变量完整
- **WHEN** 检查 index.css 中 CSS 自定义属性定义
- **THEN** 3 个文字色阶变量全部存在且色值正确

### Requirement: Agent 标识色变量统一定义
index.css SHALL 在 `:root` 中统一定义 3 个 Agent 标识色变量：`--agent-claude`(#DA7756)、`--agent-opencode`(#10B981)、`--agent-orchestrator`(#EAB308)。组件 SHALL 通过 `var(--agent-xxx)` 引用，不通过内联 style 动态拼接变量名。

#### Scenario: Agent 标识色变量存在
- **WHEN** 检查 index.css 中 CSS 自定义属性定义
- **THEN** `--agent-claude`、`--agent-opencode`、`--agent-orchestrator` 三个变量全部存在且色值正确

#### Scenario: 组件通过变量引用 Agent 标识色
- **WHEN** 检查 MessageBubble 等 Agent 相关组件
- **THEN** Agent 标识色通过 `var(--agent-claude)` 等固定变量名引用，不存在 `var(--agent-${type})` 动态拼接模式

### Requirement: 语义色变量定义
index.css SHALL 定义全部语义色变量：`--color-brand`(#6366F1)、`--color-success`(#22C55E)、`--color-warning`(#F59E0B)、`--color-error`(#EF4444)。

#### Scenario: 语义色变量完整
- **WHEN** 检查 index.css 中 CSS 自定义属性定义
- **THEN** 4 个语义色变量全部存在且色值正确

### Requirement: 变量命名一致性
所有 CSS 自定义属性 SHALL 使用 visual-style-guide.md 中定义的变量名（`--bg-xxx`、`--text-xxx`、`--color-xxx`、`--agent-xxx`），不引入额外的非标准命名。

#### Scenario: 无非标准变量名
- **WHEN** 检查 index.css 中所有 CSS 自定义属性
- **THEN** 所有变量名均来自 visual-style-guide.md 定义或 Tailwind/shadcn 标准变量
