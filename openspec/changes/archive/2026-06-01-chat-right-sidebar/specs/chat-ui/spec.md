## MODIFIED Requirements

### Requirement: Chat page four-column layout
ImPage 布局 SHALL 从三栏扩展为四栏：`IconSidebar(56px)` → `ConversationList(280px)` → `ChatArea(flex-1)` → `RightSidebar(280px)`。RightSidebar 仅群聊时显示。

#### Scenario: Group chat four-column layout
- **WHEN** 用户打开群聊会话
- **THEN** 布局为四栏：IconSidebar → ConversationList → ChatArea → RightSidebar

#### Scenario: Single chat three-column layout
- **WHEN** 用户打开单聊会话
- **THEN** 布局保持三栏：IconSidebar → ConversationList → ChatArea，ChatArea 占满剩余宽度
