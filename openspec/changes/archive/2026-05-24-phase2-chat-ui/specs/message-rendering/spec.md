## ADDED Requirements

### Requirement: Markdown rendering in messages
The system SHALL render message content as Markdown using `react-markdown` with `remark-gfm` for GitHub Flavored Markdown support (tables, strikethrough, task lists, autolinks). MarkdownRenderer is a Dumb component receiving `content` as a string prop.

#### Scenario: Standard Markdown rendering
- **WHEN** a message contains Markdown syntax (headings, bold, italic, lists, links)
- **THEN** system renders it as formatted HTML within the message bubble

#### Scenario: GFM table rendering
- **WHEN** a message contains a Markdown table
- **THEN** system renders it as an HTML table with dark background and muted borders

#### Scenario: Inline code rendering
- **WHEN** a message contains inline `code` spans
- **THEN** system renders them with Geist Mono font, code background `#0D0F14`, and 6px border-radius

### Requirement: Code block syntax highlighting
The system SHALL render fenced code blocks with syntax highlighting using Shiki (Tokyo Night theme). Code blocks MUST display with Geist Mono 13px font, line numbers in tertiary text color (`#5A6070`), and 8px border-radius on `#0D0F14` background.

#### Scenario: Highlighted code block with language
- **WHEN** a message contains a fenced code block with a language identifier (e.g., ```python)
- **THEN** system renders it with Shiki syntax highlighting in Tokyo Night colors, line numbers, and code font

#### Scenario: Code block without language
- **WHEN** a message contains a fenced code block without a language identifier
- **THEN** system renders it with monospace font and dark background but no syntax coloring

#### Scenario: Async highlighting fallback
- **WHEN** Shiki is still loading or highlighting is in progress
- **THEN** system renders the code block with plain monospace text first, then replaces with highlighted version when ready

### Requirement: Virtual list for message scrolling
The system SHALL use `@tanstack/react-virtual` to implement virtual scrolling for the message list when message count exceeds 50 items. Below 50 items, direct rendering is used. MessageList is a Smart component managing virtual list state.

#### Scenario: Short message list
- **WHEN** current session has 50 or fewer messages
- **THEN** system renders all messages directly without virtualization

#### Scenario: Long message list virtualization
- **WHEN** current session has more than 50 messages
- **THEN** system uses virtual scrolling to render only visible messages plus a small overscan buffer

#### Scenario: Dynamic height measurement
- **WHEN** messages have varying heights (short text vs long code blocks)
- **THEN** virtual list correctly measures and accounts for dynamic heights, recalculating when streaming content changes

### Requirement: Auto-scroll behavior
The system SHALL automatically scroll to the bottom when a new message arrives or streaming content updates, but MUST NOT force-scroll if the user has manually scrolled up to read history.

#### Scenario: Auto-scroll on new message
- **WHEN** a new message is added while user is at the bottom of the list
- **THEN** list scrolls to show the new message

#### Scenario: Respect manual scroll position
- **WHEN** user has scrolled up to read message history
- **THEN** system does NOT auto-scroll on new streaming content, but shows a "new messages" indicator

#### Scenario: Scroll to bottom button
- **WHEN** user has scrolled up and new messages exist below
- **THEN** a "scroll to bottom" button appears; clicking it scrolls to the latest message
