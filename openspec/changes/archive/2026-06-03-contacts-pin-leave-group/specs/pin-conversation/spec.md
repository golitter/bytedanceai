## ADDED Requirements

### Requirement: Pin toggle in RightSidebar
系统 SHALL 在群聊 RightSidebar 中提供"置顶会话"按钮。点击后调用 `updateTaskPin(taskId, new Date().toISOString())` 置顶，再次点击调用 `updateTaskPin(taskId, null)` 取消置顶。按钮状态 SHALL 反映当前置顶状态（已置顶显示高亮 + "取消置顶"文案，未置顶显示普通样式 + "置顶会话"文案）。

#### Scenario: Pin a conversation
- **WHEN** 用户在群聊 RightSidebar 中点击"置顶会话"按钮
- **THEN** 系统调用 `PATCH /api/tasks/:taskId` 设置 `pinned_at` 为当前时间，会话列表刷新后该会话排到最前面，按钮变为"取消置顶"高亮状态

#### Scenario: Unpin a conversation
- **WHEN** 用户在已置顶会话的 RightSidebar 中点击"取消置顶"按钮
- **THEN** 系统调用 `PATCH /api/tasks/:taskId` 设置 `pinned_at` 为 null，会话列表刷新后该会话恢复到非置顶区域，按钮变为"置顶会话"普通样式

### Requirement: Pinned conversations sorted first
会话列表 SHALL 按置顶优先排序：置顶会话按 `pinned_at` 降序排列在最前，非置顶会话按 `created_at` 降序排在后面。Backend `ListTasks` SHALL 使用 `ORDER BY pinned_at DESC NULLS LAST, created_at DESC`。

#### Scenario: Pinned conversations appear first in list
- **WHEN** 会话列表加载时存在已置顶和未置顶的会话
- **THEN** 置顶会话 SHALL 全部排在非置顶会话之前，置顶会话之间按置顶时间降序排列

### Requirement: Pin icon on conversation item
ConversationItem 组件 SHALL 在会话名称旁显示 Pin 图标（仅当 `pinnedAt` 存在时）。图标使用 Lucide Pin 组件，旋转 45 度，颜色为 primary。

#### Scenario: Pin icon visible on pinned conversation
- **WHEN** 一个已置顶的会话渲染在 ConversationItem 中
- **THEN** 会话名称右侧 SHALL 显示一个旋转 45 度的 Pin 图标，颜色为主题 primary 色

#### Scenario: No pin icon on unpinned conversation
- **WHEN** 一个未置顶的会话渲染在 ConversationItem 中
- **THEN** 不 SHALL 显示 Pin 图标
