## MODIFIED Requirements

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

#### Scenario: 长消息列表虚拟化
- **WHEN** 当前 session 有超过 50 条消息（含时间分隔线）
- **THEN** 系统使用虚拟滚动渲染可见区域的消息和时间分隔线

#### Scenario: 时间分隔线的动态高度
- **WHEN** 虚拟列表中的时间分隔线被渲染
- **THEN** 使用约 40px 的估计高度，实际高度由动态测量确定
