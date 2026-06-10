## ADDED Requirements

### Requirement: Import skill to agent
Agent 详情页 Skills 区域 SHALL 展示「导入外部 Skill」按钮（仅 adapter 层 agent 可见），点击弹出导入选择对话框。

#### Scenario: Show import button for adapter agent
- **WHEN** 当前 Agent 的 agent_type 为 `claude-code` / `opencode` / `codex`
- **THEN** Skills 区域底部显示「导入外部 Skill」按钮

#### Scenario: Hide import button for orchestrator
- **WHEN** 当前 Agent 的 agent_type 为 `orchestrator`
- **THEN** Skills 区域不显示「导入外部 Skill」按钮

### Requirement: Import selection dialog
导入选择对话框 SHALL 列出 SkillsHub 中所有 external skills（不含已导入的），用户可勾选后确认导入。

#### Scenario: Show available skills
- **WHEN** 打开导入选择对话框
- **THEN** 列出所有 hub 中 `builtin=false` 的 skills，已导入的显示为灰色不可选

#### Scenario: Select and confirm import
- **WHEN** 用户勾选 1+ 个 skill 并点击「确认导入」
- **THEN** 发送 import 请求，成功后关闭对话框，新 skill 以动画出现在 Skills 列表中

#### Scenario: No skills available
- **WHEN** 所有 external skills 都已导入
- **THEN** 对话框显示空状态提示

### Requirement: Remove imported skill
已导入的 external skill SHALL 显示「移除」按钮，点击后发送删除请求并从列表中动画移除。Builtin skills 不显示移除按钮。

#### Scenario: Remove external skill
- **WHEN** 用户点击 external skill 的「移除」按钮
- **THEN** 发送 `DELETE /api/skills/:name/sessions/:sessionId`，成功后该 skill card 以动画移除

#### Scenario: Builtin skill no remove button
- **WHEN** 渲染 builtin skill
- **THEN** 不显示「移除」按钮

### Requirement: Import validation feedback
前端 SHALL 对导入操作进行前置校验，提供友好提示。

#### Scenario: Duplicate import warning
- **WHEN** 用户在导入对话框中勾选了已导入的 skill
- **THEN** 该项显示为灰色不可选，标注「已导入」

#### Scenario: Import success toast
- **WHEN** 导入成功
- **THEN** 显示「已导入 N 个 Skill」成功 Toast

#### Scenario: Remove success toast
- **WHEN** 移除成功
- **THEN** 显示「已从 Agent 移除 Skill「xxx」」信息 Toast
