## ADDED Requirements

### Requirement: Member list display
系统 SHALL 从 `groupAgentTypes`、`groupAgentNames`、`groupSessions` props 提取成员数据，展示所有群成员。

#### Scenario: All members displayed
- **WHEN** 群聊有 5 个成员（1 user + 4 agents）
- **THEN** 成员列表显示 5 行，每行包含头像、名称、角色

### Requirement: Member role display
每行成员 SHALL 显示角色标签：Owner（创建者）、Admin（Orchestrator）、Member（Subagent）。

#### Scenario: Owner member
- **WHEN** 成员为创建者（用户自己）
- **THEN** 显示 Owner badge，使用 `--primary-soft` 背景和 `--primary` 文字色

#### Scenario: Admin member
- **WHEN** 成员为 Orchestrator
- **THEN** 显示 Admin badge，使用 `--warning-soft` 背景和 `--color-warning` 文字色

#### Scenario: Regular member
- **WHEN** 成员为 Subagent
- **THEN** 不显示角色 badge，仅显示 "Subagent · 职责" 文字

### Requirement: Online status indicator
每行成员 SHALL 显示在线状态圆点：在线为绿色（`--color-success`），离线为灰色（`--text-tertiary`）。

#### Scenario: Member with active SSE stream
- **WHEN** 某成员 session 有活跃的 SSE stream
- **THEN** 显示绿色实心圆点

#### Scenario: Member without active stream
- **WHEN** 某成员 session 没有活跃的 SSE stream
- **THEN** 显示灰色实心圆点

### Requirement: Click member to navigate
点击成员行 SHALL 切换到该成员的独立 1v1 会话（更新 `activeSessionId`）。

#### Scenario: Click agent member
- **WHEN** 用户点击某 Agent 成员行
- **THEN** 切换 `activeSessionId` 到该成员的 session，聊天区显示该成员的 1v1 会话

#### Scenario: Click owner (self)
- **WHEN** 用户点击自己的成员行
- **THEN** 不执行跳转（自己没有独立 1v1 会话）

### Requirement: Member avatar colors
成员头像颜色 SHALL 与 `MessageBubble.tsx` 中的 `AGENT_COLORS` 映射一致。

#### Scenario: Agent avatar color
- **WHEN** 成员为 Claude Code（agentType = 'claude_code'）
- **THEN** 头像使用 `--agent-claude` 色（#DA7756）

#### Scenario: Orchestrator avatar color
- **WHEN** 成员为 Orchestrator（agentType = 'orchestrator'）
- **THEN** 头像使用 `--agent-orchestrator` 色（#EAB308）
