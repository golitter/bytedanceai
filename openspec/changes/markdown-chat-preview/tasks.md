## 1. 用户消息 Markdown 渲染

- [x] 1.1 修改 `MessageRenderer.tsx`：用户消息从 `{msg.content}` 改为 `<MarkdownRenderer content={msg.content} />`
- [x] 1.2 修改 `MessageBubble.tsx`：用户气泡 div 添加 `[&_a]:`、`[&_pre]:`、`[&_code]:`、`[&_blockquote]:` 等 Tailwind 子选择器，确保 Markdown 元素在 `bg-primary-soft` 背景上可读
- [ ] 1.3 验证：发送含标题、列表、代码块、表格、引用的用户消息，确认渲染正确且纯文本消息无异常

## 2. MessageInput 受控组件改造

- [x] 2.1 引入 `useState<string>` 管理 textarea 值，替换所有 `textareaRef.current.value` 直接赋值
- [x] 2.2 修改 `handleSend`：从 state 读取值，发送后 `setValue('')`
- [x] 2.3 修改 `insertMention`：通过 state 拼接 mention 文本而非直接操作 DOM
- [ ] 2.4 验证：确认 @提及、Enter 发送、IME 输入、Shift+Enter 换行功能正常

## 3. Markdown 模式切换按钮

- [x] 3.1 在输入栏 `border-t` 与 textarea 之间新增工具栏行，包含 Markdown 切换按钮（`FileText` 图标 + "Markdown" 文案）
- [x] 3.2 添加 `mdMode: boolean` state，按钮点击时切换，控制工具栏样式（active/inactive）
- [ ] 3.3 验证：确认切换按钮 UI 正确，默认关闭不影响现有输入行为

## 4. 双栏实时预览布局

- [x] 4.1 实现双栏条件渲染：关闭时显示单栏（现有布局），开启时显示 `[textarea flex:1] [divider] [preview flex:1] [send-btn]` 三段 flex 布局
- [x] 4.2 预览区复用 `MarkdownRenderer`，通过 `useMemo` + debounce(150ms) 缓冲 textarea 值的 Markdown 解析
- [x] 4.3 空输入时预览区显示 placeholder 文案
- [x] 4.4 切换模式时同步 textarea 内容（单栏 ↔ 双栏）
- [ ] 4.5 验证：输入 Markdown 确认预览区实时渲染，切换单/双栏内容不丢失

## 5. 自动增高 + 同步滚动

- [x] 5.1 双栏输入区自动增高：`min-height: 120px`，根据 textarea `scrollHeight` 动态调整，`max-height: 60vh`
- [x] 5.2 内容超出最大高度时 textarea 和预览区各自 `overflow-y: auto`
- [x] 5.3 实现同步滚动：绑定 textarea 的 `scroll`/`click`/`keyup` 事件，按比例映射预览区 `scrollTop`
- [x] 5.4 预览区重新渲染后恢复滚动位置
- [ ] 5.5 验证：粘贴长 Markdown，确认输入区自动增高到上限后预览跟随光标滚动

## 6. 端到端验证

- [ ] 6.1 启动前端 `make run-frontend`，完整测试所有场景：发送 Markdown 用户消息、双栏预览、@提及、发送、自动增高、同步滚动
- [ ] 6.2 更新设计文档 `frontend/docs/design/markdown-rendering-and-preview.md` 标记为已实现
