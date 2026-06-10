## ADDED Requirements

### Requirement: Three-column chat layout
The system SHALL render a three-column layout: ChatSidebar (left, ~260px fixed width), ChatArea (center, flex-1), and a reserved empty slot (right, collapsible, hidden in Phase 2). The layout MUST fill full viewport height (`h-screen`).

#### Scenario: Default layout rendering
- **WHEN** user navigates to `/tasks/:taskId`
- **THEN** system renders ChatSidebar, ChatArea, and hidden reserved area

#### Scenario: Responsive behavior
- **WHEN** viewport width is less than 768px
- **THEN** sidebar collapses to overlay mode, chat area fills full width

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
The system SHALL render a chat area consisting of: a header (task title), a scrollable message list, and a message input at the bottom. ChatArea is a Smart component managing streaming state via `useChatStream` hook.

#### Scenario: Empty chat state
- **WHEN** the current session has no messages
- **THEN** chat area displays a centered welcome message "开始一段新对话吧"

#### Scenario: Message list rendering
- **WHEN** messages exist in the current session
- **THEN** system renders MessageBubble components for each message, with auto-scroll to bottom on new messages

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
The system SHALL render agent avatars as rounded squares (8px radius) at 32px size, using the agent's identity color as background. A 4px status dot MUST appear in the bottom-right corner. AgentAvatar is a Dumb component.

#### Scenario: Agent avatar rendering
- **WHEN** an agent message is displayed
- **THEN** a 32px rounded-square avatar with agent identity color is shown

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
