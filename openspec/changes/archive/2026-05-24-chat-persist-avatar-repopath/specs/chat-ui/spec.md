## ADDED Requirements

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

## MODIFIED Requirements

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
