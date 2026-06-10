## ADDED Requirements

### Requirement: 次要文字色 SHALL 使用 text-text-secondary 类
所有需要显示次要文字色（对应 `--text-secondary: #8B91A0`）的组件 SHALL 使用 Tailwind 类 `text-text-secondary`，不得使用以下替代方式：
- `text-secondary`（映射到 shadcn secondary-foreground，暗色背景下不可见）
- `style={{ color: 'var(--text-secondary)' }}`（inline style 方式不统一）

#### Scenario: 修正 text-secondary 为 text-text-secondary
- **WHEN** 组件当前使用 Tailwind 类 `text-secondary` 渲染次要文字
- **THEN** 替换为 `text-text-secondary`，确保暗色背景下文字可见（色值 `#8B91A0`）

#### Scenario: 修正 inline style 为 Tailwind 类
- **WHEN** 组件当前使用 `style={{ color: 'var(--text-secondary)' }}` 渲染次要文字
- **THEN** 移除 inline style，改用 Tailwind 类 `text-text-secondary`

#### Scenario: 新组件使用正确的文字色类名
- **WHEN** 开发者编写新组件需要次要文字色
- **THEN** 使用 `text-text-secondary`，而非 `text-secondary` 或 inline style

### Requirement: 硬编码语义色值 SHALL 替换为 CSS 变量或 Tailwind 语义类
组件中不得出现硬编码的十六进制色值用于表达语义状态（如成功/警告/错误）。SHALL 使用以下映射：
- `#EF4444`（错误红）→ `bg-destructive` 或 `bg-[var(--destructive)]`
- `#F59E0B`（警告黄）→ `bg-[var(--color-warning)]`
- `#22C55E`（成功绿）→ `bg-[var(--color-success)]`

#### Scenario: TerminalPanel 状态灯色值替换
- **WHEN** TerminalPanel 渲染 Agent 进程状态指示器
- **THEN** 使用 CSS 变量引用（`var(--color-success)` 等）或 Tailwind 语义类，不硬编码 `#EF4444`/`#F59E0B`/`#22C55E`

#### Scenario: 全局无语义色硬编码
- **WHEN** 扫描 `frontend/src/` 下所有 TSX 文件
- **THEN** 不存在用于 UI 元素的硬编码 `#EF4444`、`#F59E0B`、`#22C55E`、`#6366F1` 色值
