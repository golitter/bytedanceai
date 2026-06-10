## ADDED Requirements

### Requirement: Diff 主题色值变量化
系统 SHALL 在 `index.css` 的 `.dark` 块中新增以下 CSS 变量，替换 diff 主题覆盖中的硬编码色值：

- `--diff-insert-color` — 新增行文字色（对应 `#22C55E`）
- `--diff-insert-bg` — 新增行浅背景（对应 `rgba(34, 197, 94, 0.08)`）
- `--diff-delete-color` — 删除行文字色（对应 `#EF4444`）
- `--diff-delete-bg` — 删除行浅背景（对应 `rgba(239, 68, 68, 0.08)`）
- `--diff-insert-bg-strong` — 新增代码区深背景（对应 `rgba(34, 197, 94, 0.1)`）
- `--diff-delete-bg-strong` — 删除代码区深背景（对应 `rgba(239, 68, 68, 0.1)`）

diff 主题覆盖区的 `.diff-gutter-insert`、`.diff-gutter-delete`、`.diff-code-insert`、`.diff-code-delete` MUST 引用这些变量，不再直接写 hex/rgba 值。

#### Scenario: diff 新增行颜色由变量驱动
- **WHEN** 查看 `.diff-gutter-insert` 和 `.diff-code-insert` 的样式
- **THEN** color/background 引用 `var(--diff-insert-color)`/`var(--diff-insert-bg)`/`var(--diff-insert-bg-strong)` 而非硬编码的 `#22C55E` 或 `rgba(34, 197, 94, ...)`

#### Scenario: diff 视觉效果不变
- **WHEN** 在暗色模式下查看 diff 卡片
- **THEN** 新增行和删除行的颜色与变量化之前完全一致（相同的色值，只是通过变量引用）

#### Scenario: diff 变量可在 Tailwind theme 中引用
- **WHEN** 查看 `@theme inline` 块
- **THEN** diff 相关的 CSS 变量也被映射为 `--color-diff-*` 形式，可在 Tailwind 工具类中使用
