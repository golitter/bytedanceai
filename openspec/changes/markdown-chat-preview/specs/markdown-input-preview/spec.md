## ADDED Requirements

### Requirement: Markdown mode toggle button
MessageInput 组件 SHALL 在输入区上方显示一个 Markdown 模式切换按钮，默认关闭。关闭时输入区保持当前单栏 textarea 行为不变。

#### Scenario: Default state is off
- **WHEN** 用户打开聊天页面
- **THEN** 输入栏显示普通单栏 textarea + 发送按钮，Markdown 按钮未激活

#### Scenario: Toggle on
- **WHEN** 用户点击 Markdown 按钮
- **THEN** 按钮变为激活状态，输入区从单栏切换为双栏布局

#### Scenario: Toggle off
- **WHEN** 用户再次点击已激活的 Markdown 按钮
- **THEN** 按钮恢复未激活状态，输入区恢复为单栏，textarea 内容保留

### Requirement: Dual-pane input layout
开启 Markdown 模式后，输入区 SHALL 变为三段 flex 布局：左栏 textarea 编辑区（flex:1）、中间 1px 分隔线、右栏 MarkdownRenderer 实时预览区（flex:1）、最右侧发送按钮（固定宽度）。

#### Scenario: Dual-pane display
- **WHEN** Markdown 模式开启
- **THEN** 输入区显示为左编辑/右预览双栏，发送按钮在最右侧

#### Scenario: Empty preview
- **WHEN** Markdown 模式开启且 textarea 为空
- **THEN** 预览区显示 placeholder 提示文案

#### Scenario: Live preview update
- **WHEN** 用户在 textarea 中输入 Markdown 内容
- **THEN** 预览区实时渲染对应 Markdown（debounce ≤ 150ms）

#### Scenario: Content sync on toggle
- **WHEN** 用户在 Markdown 模式下输入内容后关闭 Markdown 模式
- **THEN** 单栏 textarea 保留已输入的内容

### Requirement: Auto-growing input pane
双栏输入区 SHALL 根据内容自动增高，最小高度 120px，最大高度为视口高度的 60%（60vh）。

#### Scenario: Content grows pane
- **WHEN** 用户输入的内容超过当前输入区高度
- **THEN** 输入区和预览区同步增高，不超过 60vh

#### Scenario: Content shrinks pane
- **WHEN** 用户删除内容使所需高度减小
- **THEN** 输入区高度相应缩小，不低于 120px

#### Scenario: Max height reached
- **WHEN** 内容超过 60vh 最大高度
- **THEN** textarea 和预览区各自独立滚动

### Requirement: Sync scroll between editor and preview
当双栏内容超出最大高度时，预览区 SHALL 根据 textarea 的滚动比例同步滚动，跟随光标位置。

#### Scenario: Scroll sync on textarea scroll
- **WHEN** 用户在 textarea 中滚动
- **THEN** 预览区按相同滚动比例滚动到对应位置

#### Scenario: Scroll sync on cursor movement
- **WHEN** 用户在 textarea 中点击或用键盘移动光标
- **THEN** 预览区同步滚动到与光标位置对应的位置

#### Scenario: Scroll sync on input
- **WHEN** 用户输入新内容触发预览区重新渲染
- **THEN** 预览区恢复到与 textarea 当前滚动比例对应的位置
