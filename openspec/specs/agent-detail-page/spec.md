## MODIFIED Requirements

### Requirement: Agent detail page Skills section
Agent 详情页 Skills 区域 SHALL 展示当前 Agent 的所有 skills（builtin + external），每个 skill card 显示名称、描述、来源标签。External skill 额外显示「移除」按钮。区域底部显示「导入外部 Skill」按钮（仅 adapter 层 agent 可见）。

#### Scenario: Display mixed skills
- **WHEN** Agent 拥有 2 个 builtin skill 和 1 个 external skill
- **THEN** 渲染 3 个 skill cards，builtin 显示绿色标签无操作按钮，external 显示 indigo 标签和「移除」按钮

#### Scenario: Import button for adapter agent
- **WHEN** Agent 的 agent_type 为 `claude-code` / `opencode` / `codex`
- **THEN** Skills 区域底部显示虚线「导入外部 Skill」按钮

#### Scenario: No import button for orchestrator
- **WHEN** Agent 的 agent_type 为 `orchestrator`
- **THEN** Skills 区域不显示导入按钮

#### Scenario: Remove external skill
- **WHEN** 用户点击 external skill 的「移除」按钮
- **THEN** 发送 `DELETE /api/skills/:name/sessions/:sessionId`，成功后 card 动画移除

#### Scenario: Import new skill
- **WHEN** 用户通过导入对话框成功导入 skill
- **THEN** 新 skill card 以动画出现在列表中，skill 计数更新

#### Scenario: Empty skills
- **WHEN** Agent 无任何 skills
- **THEN** 显示空状态提示和导入按钮（adapter agent）
