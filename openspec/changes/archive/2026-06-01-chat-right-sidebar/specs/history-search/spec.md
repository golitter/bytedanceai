## ADDED Requirements

### Requirement: Search input with debounce
系统 SHALL 提供搜索输入框，输入 ≥1 字符后展示下拉搜索结果，使用 300ms debounce。

#### Scenario: Type keyword shows results
- **WHEN** 用户在搜索框输入 ≥1 字符
- **THEN** 300ms 无新输入后，展示下拉搜索结果列表

#### Scenario: Clear input hides results
- **WHEN** 用户清空搜索框
- **THEN** 下拉结果列表立即关闭

### Requirement: Frontend message filtering
搜索 SHALL 在前端过滤已加载消息，匹配 `message.content` 和 `block.content` 中的关键词，结果按时间倒序排列。

#### Scenario: Match in message content
- **WHEN** 用户搜索 "重构" 且某条消息 `content` 包含 "重构"
- **THEN** 该消息出现在搜索结果中

#### Scenario: Match in block content
- **WHEN** 用户搜索 "jwt" 且某条消息的 `block.content` 包含 "jwt"
- **THEN** 该消息出现在搜索结果中

#### Scenario: Results ordered by time descending
- **WHEN** 搜索匹配多条消息
- **THEN** 结果按时间倒序排列（最新在前）

### Requirement: Search result display
每条搜索结果 SHALL 显示：Agent 头像 + 名称 + 时间 + 匹配文本摘要，匹配关键词使用 `<mark>` 高亮。

#### Scenario: Result shows agent info
- **WHEN** 搜索结果包含 Agent 消息
- **THEN** 显示该 Agent 的头像、名称、消息时间

#### Scenario: Keyword highlighted
- **WHEN** 搜索关键词为 "session"
- **THEN** 结果文本中 "session" 被包裹在 `<mark>` 标签中，使用 `--warning-soft` 背景和 `--color-warning` 文字色

### Requirement: Click result to navigate
点击搜索结果 SHALL 滚动到对应消息位置并添加高亮闪烁动画（800ms 后消失）。

#### Scenario: Click result scrolls to message
- **WHEN** 用户点击某条搜索结果
- **THEN** 聊天区域滚动到对应消息，消息添加高亮闪烁动画，800ms 后动画消失

### Requirement: Dropdown dismissal
搜索下拉 SHALL 在点击外部区域或清空输入时关闭。

#### Scenario: Click outside closes dropdown
- **WHEN** 搜索下拉打开时，用户点击搜索区域外的地方
- **THEN** 下拉关闭

### Requirement: Result count limit
搜索结果 SHALL 最多显示 50 条，超出时显示提示文字。

#### Scenario: More than 50 results
- **WHEN** 搜索匹配超过 50 条消息
- **THEN** 只显示前 50 条，底部显示 "还有 N 条匹配结果"
