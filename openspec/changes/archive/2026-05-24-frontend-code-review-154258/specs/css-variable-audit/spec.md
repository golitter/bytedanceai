## ADDED Requirements

### Requirement: 所有组件 SHALL 使用 CSS 变量而非硬编码色值
组件中的颜色值 MUST 通过 CSS 变量（如 `var(--bg-card)`、`var(--text-primary)`）引用，不得使用 Tailwind 原始色阶（如 `bg-gray-200`、`text-green-700`）或硬编码 hex 值。

#### Scenario: SessionList 使用 CSS 变量
- **WHEN** 渲染 SessionList 组件
- **THEN** 所有背景色、文字色均通过 CSS 变量引用，无 `bg-gray-*`、`text-green-*` 等 Tailwind 色阶

#### Scenario: AgentAvatar 使用 CSS 变量
- **WHEN** 渲染 AgentAvatar 组件
- **THEN** Agent 标识色、阴影色通过 CSS 变量引用（如 `var(--agent-claude)`），无硬编码 hex 值

#### Scenario: ChatArea 错误横幅使用暗色主题色
- **WHEN** 渲染 ChatArea 验证错误横幅
- **THEN** 背景色为暗色主题下的语义色（如 `rgba(239, 68, 68, 0.1)`），非浅色 `#FEF2F2`

### Requirement: 纯白文字 SHALL 被禁止
所有组件 MUST NOT 使用 `text-white`、`#fff`、`#FFFFFF`、`rgb(255,255,255)` 作为文字色，正文 SHALL 使用 `var(--text-primary)` (`#E8EBF0`)。

#### Scenario: 对话框无纯白文字
- **WHEN** 渲染 AgentEditDialog、NewChatDialog
- **THEN** 所有文字色使用 CSS 变量，无 `text-white` 类名

#### Scenario: MessageInput 无纯白文字
- **WHEN** 渲染 MessageInput
- **THEN** 占位符和按钮文字使用 CSS 变量，无 `text-white`

#### Scenario: Button 组件无默认纯白
- **WHEN** 渲染 shadcn Button 默认变体
- **THEN** 文字色为 `var(--text-primary)` 而非 `#FFFFFF`

### Requirement: 缺失的 CSS 变量 SHALL 补齐
`index.css` 中 MUST 定义所有 VSG 规范中的设计令牌，包括代码块背景 `--code-bg`、边框透明度统一值等。

#### Scenario: 代码块背景有变量
- **WHEN** 查看 `index.css` 的 `:root` / `.dark` 定义
- **THEN** 包含 `--code-bg: #0D0F14` 变量

#### Scenario: 组件使用代码块变量
- **WHEN** 渲染 CodeBlock 或 MarkdownRenderer
- **THEN** 代码块背景色引用 `var(--code-bg)`，非硬编码 `#0D0F14`

### Requirement: 边框透明度 SHALL 统一为 0.06
所有组件边框 MUST 使用 `rgba(255, 255, 255, 0.06)`（CSS 变量 `--border`），不得使用 `rgba(255, 255, 255, 0.1)` 等非规范值。

#### Scenario: 各组件边框一致
- **WHEN** 检查所有组件的 border-color 值
- **THEN** 统一使用 `var(--border)` 或 `rgba(255, 255, 255, 0.06)`，无 0.1 等偏差值
