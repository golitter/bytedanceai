## Requirements

### Requirement: Chat page four-column layout
ImPage 布局 SHALL 从三栏扩展为四栏：`IconSidebar(56px)` → `ConversationList(280px)` → `ChatArea(flex-1)` → `RightSidebar(280px)`。RightSidebar 仅群聊时显示。

#### Scenario: Group chat four-column layout
- **WHEN** 用户打开群聊会话
- **THEN** 布局为四栏：IconSidebar → ConversationList → ChatArea → RightSidebar

#### Scenario: Single chat three-column layout
- **WHEN** 用户打开单聊会话
- **THEN** 布局保持三栏：IconSidebar → ConversationList → ChatArea，ChatArea 占满剩余宽度

### Requirement: Chat sidebar with session list
The system SHALL display a sidebar containing: app logo/title, a "New Chat" button, and a list of sessions belonging to the current task. Each session item MUST show its title and relative timestamp.

#### Scenario: Session list display
- **WHEN** user opens a task's chat page
- **THEN** sidebar shows all sessions for that task, each with agent type and relative time

#### Scenario: Create new session
- **WHEN** user clicks "New Chat" button
- **THEN** system creates a new session via backend API and selects it as the active session

#### Scenario: Switch session
- **WHEN** user clicks a session in the sidebar
- **THEN** system loads that session's messages and switches the active session

#### Scenario: Active session highlight
- **WHEN** a session is currently selected
- **THEN** that session item shows a 2px left brand-color border and hover background (`#22262F`)

#### Scenario: Empty session list
- **WHEN** no sessions exist for the current task
- **THEN** sidebar shows "New Chat" button and a prompt to start a conversation

### Requirement: Chat area with message display
系统 SHALL 渲染聊天区域，包含：顶部标题栏、可滚动的消息列表（含时间分隔线）、底部消息输入框。ChatArea 是 Smart 组件，通过 `useChatStream` hook 管理流式状态。MessageList 在渲染消息时，根据相邻消息的时间差动态插入 TimeDivider 组件。MessageList 支持上拉到顶部时加载更早的历史消息。

#### Scenario: 空聊天状态
- **WHEN** 当前 session 没有消息
- **THEN** 聊天区域显示居中的欢迎消息"开始一段新对话吧"

#### Scenario: 带时间分隔线的消息列表
- **WHEN** 消息列表中存在满足时间分隔条件的相邻消息
- **THEN** 在对应位置插入 TimeDivider 组件，显示相对时间

#### Scenario: 上拉加载历史消息
- **WHEN** 用户滚动到消息列表顶部（scrollTop === 0）且 has_more 为 true
- **THEN** 系统触发加载更早的历史消息，显示加载指示器

#### Scenario: 历史消息加载完成
- **WHEN** 更早的历史消息加载完成并插入到列表顶部
- **THEN** 系统恢复滚动位置，用户仍然看到加载前正在阅读的消息

#### Scenario: 全部历史消息已加载
- **WHEN** has_more 为 false
- **THEN** 上拉不再触发加载

### Requirement: Virtual list for message scrolling
系统 SHALL 使用 `@tanstack/react-virtual` 实现虚拟滚动，消息数量超过 50 条时启用。时间分隔线作为列表项参与虚拟化，estimateSize 约 40px。MessageList 是 Smart 组件管理虚拟列表状态。

#### Scenario: 短消息列表
- **WHEN** 当前 session 有 50 条或更少的消息（含时间分隔线）
- **THEN** 系统直接渲染所有消息和时间分隔线，不使用虚拟化

#### Scenario: 长消息消息列表虚拟化
- **WHEN** 当前 session 有超过 50 条消息（含时间分隔线）
- **THEN** 系统使用虚拟滚动渲染可见区域的消息和时间分隔线

#### Scenario: 时间分隔线的动态高度
- **WHEN** 虚拟列表中的时间分隔线被渲染
- **THEN** 使用约 40px 的估计高度，实际高度由动态测量确定

### Requirement: User message bubble
User messages SHALL be right-aligned with brand-color background (`rgba(99,102,241,0.08)`), brand-color border (`rgba(99,102,241,0.15)`), and 10px border-radius. MessageBubble is a Dumb component receiving only props.

#### Scenario: User message display
- **WHEN** a user message is rendered
- **THEN** it appears right-aligned with the specified brand-color tinted background and border

### Requirement: Agent message bubble
Agent messages SHALL be left-aligned with card background (`#1A1D24`), a 3px left color bar matching the agent's identity color, and 10px border-radius. Each agent message MUST display an AgentAvatar and agent name.

#### Scenario: Agent message display
- **WHEN** an agent message is rendered
- **THEN** it appears left-aligned with the agent's identity color bar, avatar, and name label

#### Scenario: Streaming agent message
- **WHEN** an agent is currently streaming a response
- **THEN** the message shows a blinking cursor `▌` at the end of the accumulated text

### Requirement: System message display
System messages SHALL be centered, using secondary text color (`#8B91A0`), 12px font size, with no bubble background.

#### Scenario: System message rendering
- **WHEN** a system message is rendered
- **THEN** it appears centered with muted text styling and no background card

### Requirement: Agent avatar with status indicator
The system SHALL render agent avatars as rounded squares (8px radius) at 32px size. When a custom avatar_url is set, the image SHALL be displayed. When no custom avatar is set, the agent's identity color is used as background. A 4px status dot MUST appear in the bottom-right corner. AgentAvatar is a Dumb component.

#### Scenario: Agent avatar rendering with custom image
- **WHEN** an agent message is displayed and the agent has a custom avatar_url
- **THEN** a 32px rounded-square avatar with the custom image is shown

#### Scenario: Agent avatar rendering with DiceBear
- **WHEN** an agent message is displayed and the agent has no custom avatar_url
- **THEN** a 32px rounded-square avatar with DiceBear generated initials image is shown

#### Scenario: Agent status indicators
- **WHEN** agent status changes
- **THEN** status dot updates: green `#22C55E` (ready, pulsing), yellow `#F59E0B` (running, rotating), gray `#5A6070` (offline), red `#EF4444` (error)

### Requirement: Message input with send functionality
The system SHALL provide a textarea input with: auto-expanding height (min 48px, max 200px), Enter to send, Shift+Enter for newline, and a send button. Input MUST be disabled during streaming. MessageInput is a Dumb component with local textarea state only.

#### Scenario: Send message via Enter key
- **WHEN** user types a message and presses Enter (without Shift)
- **THEN** system sends the message and clears the input

#### Scenario: Newline via Shift+Enter
- **WHEN** user presses Shift+Enter in the input
- **THEN** a newline is inserted without sending

#### Scenario: Input disabled during streaming
- **WHEN** agent is streaming a response
- **THEN** textarea and send button are disabled

#### Scenario: Auto-expanding textarea
- **WHEN** user types multi-line content
- **THEN** textarea height grows to fit content up to 200px, then scrolls internally

### Requirement: Route integration
The system SHALL use `/tasks/:taskId` as the chat page route. `/tasks` shows the task list. `/` redirects to `/tasks`.

#### Scenario: Navigate to chat
- **WHEN** user clicks a task in the task list
- **THEN** browser navigates to `/tasks/:taskId` and loads the chat page

#### Scenario: Direct URL access
- **WHEN** user navigates directly to `/tasks/:taskId`
- **THEN** system loads the task data and renders the chat page

### Requirement: 加载历史消息到聊天窗口
前端 SHALL 在进入 task 聊天页面时，调用 `GET /api/tasks/:taskId/messages` 加载历史消息，填充到 ChatArea 的消息列表中。

#### Scenario: 重新打开有历史的 task
- **WHEN** 用户导航到 `/tasks/:taskId`，该 task 有历史消息
- **THEN** 系统从后端加载消息并渲染到 MessageList，自动滚动到底部

#### Scenario: 首次打开无历史的 task
- **WHEN** 用户导航到 `/tasks/:taskId`，该 task 无历史消息
- **THEN** 显示空状态欢迎消息"开始一段新对话吧"

### Requirement: Agent 头像展示与编辑入口
前端 SHALL 在 AgentAvatar 组件中展示自定义头像（有 avatar_url 时）或 DiceBear 生成头像（无 avatar_url 时），并支持点击触发编辑弹窗。

#### Scenario: 展示自定义头像
- **WHEN** agent 配置了 avatar_url
- **THEN** AgentAvatar 渲染该 URL 的图片

#### Scenario: 展示 DiceBear 头像
- **WHEN** agent 未配置 avatar_url，名称为 "claude-code"
- **THEN** AgentAvatar 渲染 DiceBear 生成的 initials 头像

#### Scenario: 点击头像编辑
- **WHEN** 用户点击 agent 头像或名称
- **THEN** 弹出编辑弹窗，可修改名称和上传头像

### Requirement: NewChatDialog 增加 repoPath 输入
NewChatDialog SHALL 增加一个 repoPath 文本输入框，用户可输入本地仓库路径。

#### Scenario: 填写 repoPath 创建 task
- **WHEN** 用户在 NewChatDialog 填写 repoPath 为 "/home/user/project"
- **THEN** 创建 task 请求 body 包含 `repo_path: "/home/user/project"`

#### Scenario: repoPath 输入框可选
- **WHEN** 用户不填写 repoPath
- **THEN** 创建 task 请求 body 不包含 repo_path 或为空

### Requirement: 发送消息前 repoPath 校验反馈
前端 SHALL 在发送消息时先触发 repoPath 校验，校验中显示 loading 状态，校验失败展示具体错误。

#### Scenario: 校验中 loading
- **WHEN** 发送按钮被点击，repoPath 校验请求进行中
- **THEN** 发送按钮显示 loading 状态，输入框临时禁用

#### Scenario: 校验失败错误提示
- **WHEN** repoPath 校验返回 "路径不存在"
- **THEN** 聊天区域顶部展示红色错误横幅"路径不存在: /xxx"，不发送消息

#### Scenario: 校验通过正常发送
- **WHEN** repoPath 校验通过
- **THEN** loading 状态结束，消息正常发送
