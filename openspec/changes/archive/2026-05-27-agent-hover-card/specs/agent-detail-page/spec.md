## ADDED Requirements

### Requirement: Agent 详情页路由
系统 SHALL 提供 `/agent/:sessionId` 路由，展示该 Session 对应 Agent 的完整信息。

#### Scenario: 访问详情页
- **WHEN** 用户导航到 `/agent/:sessionId`
- **THEN** 系统展示 Agent 详情页，包含头部信息、Skills 列表、元数据三个区域

#### Scenario: 返回对话
- **WHEN** 用户点击详情页左上角的返回按钮
- **THEN** 系统导航回之前的对话页面

### Requirement: 详情页 Agent 头部
详情页 SHALL 展示 Agent 身份信息：64px 头像、agentName、agentType、status。

#### Scenario: 展示头部信息
- **WHEN** 详情页加载完成
- **THEN** 页面顶部展示 64px 头像、agentName、agentType、status 指示灯

### Requirement: 详情页 Skills 列表
详情页 SHALL 展示 Agent 的所有 Skills，每个 Skill 显示完整的 name、description、builtin 标记和 source。

#### Scenario: 展示完整 Skills
- **WHEN** Agent 有关联的 Skills
- **THEN** 每个 Skill 以卡片形式展示，包含 name（粗体）、完整 description（不截断）、builtin 标记（圆点 + "builtin" 文字）、source 信息

#### Scenario: 无 Skills
- **WHEN** Agent 没有关联的 Skills
- **THEN** Skills 区域显示"暂无技能"占位文案

### Requirement: 详情页元数据
详情页 SHALL 展示 Agent 的技术元数据：session_id、task_id、workspace 路径、创建时间、消息数量。

#### Scenario: 展示元数据
- **WHEN** 详情页加载完成
- **THEN** 元数据区域展示 session_id、task_id、workspace（如有）、创建时间、该 Session 下的消息总数
