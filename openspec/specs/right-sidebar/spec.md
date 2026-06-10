## MODIFIED Requirements

### Requirement: Right sidebar section structure
RightSidebar SHALL 包含以下区块，从上到下依次为：群成员、Git Graph、Terminal。历史消息搜索和群公告区块移除。

#### Scenario: Sections render in correct order
- **WHEN** RightSidebar 渲染
- **THEN** 区块按以下顺序显示：MembersSection → GitGraphPanel → TerminalPanel

### Requirement: Right sidebar styling
RightSidebar 宽度 SHALL 从 280px 调整为 300px，以容纳 Git Graph 和 Terminal 面板的信息密度。

#### Scenario: Visual consistency
- **WHEN** RightSidebar 渲染
- **THEN** 宽度为 300px，背景色为 `var(--sidebar)`，左边框 `1px solid var(--sidebar-border)`，与 ConversationList 视觉对齐
