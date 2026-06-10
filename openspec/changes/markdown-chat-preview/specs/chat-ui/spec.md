## MODIFIED Requirements

### Requirement: Chat area with message display
系统 SHALL 渲染聊天区域，包含：顶部标题栏、可滚动的消息列表（含时间分隔线）、底部消息输入框。ChatArea 是 Smart 组件，通过 `useChatStream` hook 管理流式状态。MessageList 在渲染消息时，根据相邻消息的时间差动态插入 TimeDivider 组件。MessageList 支持上拉到顶部时加载更早的历史消息。

**用户消息 SHALL 使用 MarkdownRenderer 渲染**，与 Agent 消息一致，支持标题、列表、代码块、表格、引用等 Markdown 元素。

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

#### Scenario: 用户消息 Markdown 渲染
- **WHEN** 用户发送包含 Markdown 语法的消息
- **THEN** 消息气泡内以 MarkdownRenderer 渲染，正确显示标题、粗体、代码块、列表、表格等元素

#### Scenario: 用户消息纯文本兼容
- **WHEN** 用户发送不含 Markdown 语法的普通文本
- **THEN** 消息气泡内正常显示文本内容，无明显格式异常
