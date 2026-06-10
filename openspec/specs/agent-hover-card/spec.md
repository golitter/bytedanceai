## ADDED Requirements

### Requirement: Hover 卡片触发
系统 SHALL 在用户鼠标悬停在聊天气泡中的 Agent 头像上时（延迟 300ms），展示 Agent 浮动卡片。仅在 `MessageBubble` 的 `variant: 'agent'` 的 `AgentAvatar` 上触发，会话列表头像不触发。

#### Scenario: 悬停显示卡片
- **WHEN** 用户将鼠标悬停在聊天气泡的 Agent 头像上超过 300ms
- **THEN** 系统在头像附近弹出 Agent Hover Card

#### Scenario: 移出隐藏卡片
- **WHEN** 用户将鼠标移出头像和卡片区域超过 200ms
- **THEN** 卡片自动消失

#### Scenario: 会话列表不触发
- **WHEN** 用户将鼠标悬停在会话列表的头像上
- **THEN** 不展示任何浮动卡片

### Requirement: Hover 卡片内容展示
Agent Hover Card SHALL 展示三层信息：身份（头像 + agentName + agentType + status）、Skills 摘要、元数据。

#### Scenario: 展示 Agent 身份
- **WHEN** 卡片显示
- **THEN** 卡片顶部展示 32px 头像、agentName（粗体）、agentType（副标题）、status 指示灯

#### Scenario: 展示 Skills 摘要
- **WHEN** Agent 有关联的 Skills
- **THEN** 每个 Skill 显示 name（粗体）和 description（截断为一行，超出显示省略号）

#### Scenario: Skills 溢出处理
- **WHEN** Agent 的 Skills 数量超过 N 个（N 默认 3）
- **THEN** 显示前 N 个 Skills + "+M 更多" 文字

#### Scenario: 展示元数据
- **WHEN** 卡片显示
- **THEN** 底部展示 session_id（截断显示）

### Requirement: Hover 卡片跳转详情页
Agent Hover Card 底部 SHALL 提供跳转链接，点击后导航到 Agent 详情页。

#### Scenario: 点击跳转
- **WHEN** 用户点击卡片底部的"查看 Agent 详情"链接
- **THEN** 系统导航到 `/agent/:sessionId` 路由
