## ADDED Requirements

### Requirement: Right sidebar container component
系统 SHALL 提供 `RightSidebar` 组件作为右侧栏容器，宽度 280px，仅当 `isGroupChat === true` 时显示。

#### Scenario: Group chat shows right sidebar
- **WHEN** 用户打开一个群聊会话（`isGroupChat === true`）
- **THEN** ChatArea 右侧显示 RightSidebar（280px），布局从三栏变为四栏

#### Scenario: Single chat hides right sidebar
- **WHEN** 用户打开一个单聊会话（`isGroupChat === false` 或 undefined）
- **THEN** 不显示 RightSidebar，布局保持三栏不变

### Requirement: Right sidebar section structure
RightSidebar SHALL 包含四个区块，从上到下依次为：历史消息搜索、群公告、群成员、更多操作。

#### Scenario: Sections render in correct order
- **WHEN** RightSidebar 渲染
- **THEN** 区块按以下顺序显示：HistorySearch → AnnouncementsSection → MembersSection → 更多操作

### Requirement: Collapsible sections
群公告和群成员区块 SHALL 支持点击标题栏折叠/展开，折叠状态 SHALL 持久化到 localStorage。

#### Scenario: Toggle section collapse state
- **WHEN** 用户点击群公告或群成员的标题栏
- **THEN** 对应区块折叠/展开，chevron 图标旋转，状态保存到 localStorage

#### Scenario: Restore collapse state on mount
- **WHEN** RightSidebar 首次挂载
- **THEN** 从 localStorage 读取折叠状态，恢复上次的折叠/展开状态

### Requirement: Right sidebar styling
RightSidebar SHALL 使用 `--sidebar` 背景、`--sidebar-border` 左边框，与 ConversationList 视觉风格一致。

#### Scenario: Visual consistency
- **WHEN** RightSidebar 渲染
- **THEN** 背景色为 `var(--sidebar)`，左边框 `1px solid var(--sidebar-border)`，与 ConversationList 视觉对齐
