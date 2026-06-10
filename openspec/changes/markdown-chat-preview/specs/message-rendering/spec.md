## ADDED Requirements

### Requirement: User bubble Markdown style adaptation
用户消息气泡内的 MarkdownRenderer 渲染内容 SHALL 在 `bg-primary-soft` 背景上保持良好可读性，包括链接颜色、代码块背景、引用块等元素。

#### Scenario: Link readability in user bubble
- **WHEN** 用户消息包含 Markdown 链接
- **THEN** 链接在 `bg-primary-soft` 背景上有足够对比度，颜色与气泡背景区分明显

#### Scenario: Code block readability in user bubble
- **WHEN** 用户消息包含行内代码或代码块
- **THEN** 代码背景色与 `bg-primary-soft` 有层次区分，代码文本清晰可读

#### Scenario: Blockquote readability in user bubble
- **WHEN** 用户消息包含引用块
- **THEN** 引用块在 `bg-primary-soft` 内有可见的左边框和背景区分
