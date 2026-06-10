## ADDED Requirements

### Requirement: Contact group data model
系统 SHALL 新增 `contact_groups` 和 `contact_group_items` 两张数据库表。`contact_groups` 存储分组定义（group_id, name, sort_order）。`contact_group_items` 存储分组与会话的多对多关联（group_id, task_id, sort_order）。未出现在 `contact_group_items` 中的 task 隐式属于"未分组"。两张表 SHALL 通过 AutoMigrate 自动创建。

#### Scenario: Auto-migrate creates tables on startup
- **WHEN** Backend 启动
- **THEN** `contact_groups` 和 `contact_group_items` 表自动创建，schema 符合设计文档定义

### Requirement: Contact group CRUD API
Backend SHALL 提供以下端点管理通讯录分组：

- `GET /api/contact-groups` — 返回所有分组及其 task 列表 + 未分组 task IDs
- `POST /api/contact-groups` — 创建分组（body: `{name}`）
- `PUT /api/contact-groups/:groupId` — 重命名分组（body: `{name}`）
- `DELETE /api/contact-groups/:groupId` — 删除分组（关联项移至未分组）
- `POST /api/contact-groups/:groupId/items` — 添加 task 到分组（body: `{task_id}`）
- `DELETE /api/contact-groups/:groupId/items/:taskID` — 从分组移除 task

#### Scenario: Create a new group
- **WHEN** 用户通过 `POST /api/contact-groups` 创建分组 `{name: "工作项目"}`
- **THEN** 系统生成唯一 group_id，创建分组记录并返回

#### Scenario: Delete a group preserves tasks
- **WHEN** 用户删除一个包含 3 个会话的分组
- **THEN** 分组被删除，3 个会话的 `contact_group_items` 记录被删除，会话移至"未分组"区域，会话本身不被删除

#### Scenario: Move task to group
- **WHEN** 用户通过 `POST /api/contact-groups/:groupId/items` 将 task 分配到分组
- **THEN** 创建关联记录，该 task 在通讯录中显示在对应分组下

#### Scenario: Remove task from group
- **WHEN** 用户通过 `DELETE /api/contact-groups/:groupId/items/:taskID` 从分组移除 task
- **THEN** 删除关联记录，该 task 回到"未分组"区域

### Requirement: ContactsPage component
前端 SHALL 新建 `ContactsPage` 组件，作为通讯录主页面。布局包含：搜索栏、置顶分组（不可折叠）、自定义分组列表（可折叠，悬停显示编辑/删除按钮）、未分组区域、新建分组按钮。

每个会话卡片（ContactCard）显示头像、名称、描述、最后活跃时间，点击 SHALL 导航到对应聊天会话。悬停显示操作菜单（分配到分组、置顶）。

#### Scenario: View contacts page with groups
- **WHEN** 用户点击 IconSidebar 的"通讯录"按钮
- **THEN** 显示 ContactsPage，展示置顶分组、自定义分组（含成员）、未分组会话

#### Scenario: Create group inline
- **WHEN** 用户点击"新建分组"按钮并输入名称后确认
- **THEN** 系统创建新分组，分组出现在通讯录页面，输入框恢复为按钮

#### Scenario: Delete group moves tasks to ungrouped
- **WHEN** 用户点击分组的删除按钮并确认
- **THEN** 分组从通讯录消失，该分组下的会话移至"未分组"区域

### Requirement: Enable contacts navigation
IconSidebar 中的"通讯录"导航按钮 SHALL 移除 `disabled` 属性，点击后切换到 contacts Tab。ImPage SHALL 将 contacts 的 PlaceholderPage 替换为 `<ContactsPage />`。

#### Scenario: Contacts tab accessible
- **WHEN** 用户点击 IconSidebar 的"通讯录"图标
- **THEN** 主内容区显示 ContactsPage，通讯录图标变为高亮激活状态

### Requirement: LeaveTask cleans up contact group items
Backend `LeaveTask`（退出群聊）的事务中 SHALL 同时删除被删除 task 在 `contact_group_items` 中的关联记录，确保通讯录分组不留悬空引用。

#### Scenario: Leave group removes from contact groups
- **WHEN** 一个被分配到分组的群聊会话被退出删除
- **THEN** 该会话在 `contact_group_items` 中的关联记录被一并删除，通讯录页面不再显示该会话
