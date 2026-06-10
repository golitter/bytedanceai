## 1. CSS 变量与全局样式

- [x] 1.1 修正 Agent 标识色：`index.css` 中 `--agent-claude` 改为 `#DA7756`、`--agent-opencode` 改为 `#10B981`、`--agent-orchestrator` 改为 `#EAB308`
- [x] 1.2 添加 `--divider` 变量：`rgba(255, 255, 255, 0.04)`，用于分割线
- [x] 1.3 body 添加 `letter-spacing: -0.01em`
- [x] 1.4 代码块/行内代码覆盖 `letter-spacing: 0`
- [x] 1.5 添加 `@keyframes status-ready-pulse`（opacity 0.6→1→0.6，2s）
- [x] 1.6 添加 `@keyframes status-running-spin`（rotate 0→360deg，1.5s）

## 2. AgentAvatar 组件

- [x] 2.1 为头像添加 8px 同色 blur 的 box-shadow 光晕
- [x] 2.2 ready 状态灯改用 `status-ready-pulse` 动画
- [x] 2.3 running 状态灯改用 `status-running-spin` 动画

## 3. MessageInput 组件

- [x] 3.1 textarea 背景色从 `var(--bg-hover)` 改为 `var(--bg-card)`

## 4. 分割线颜色

- [x] 4.1 ConversationList header 底部 border 改为 `var(--divider)`
- [x] 4.2 ChatArea header 底部 border 改为 `var(--divider)`
- [x] 4.3 MessageInput 顶部 border 改为 `var(--divider)`

## 5. 过渡时长

- [x] 5.1 ConversationItem hover 过渡改为 `transition-colors duration-120 ease-out`
- [x] 5.2 ConversationList 新建按钮 hover 过渡改为 `duration-120`

## 6. 图标 strokeWidth

- [x] 6.1 全局搜索所有 Lucide 图标的 `strokeWidth` 属性，从 `1.5` 改为 `1.25`

## 7. 审计差异文档

- [x] 7.1 在 `frontend/docs/` 下创建 `visual-style-audit.md`，记录所有规范偏差和修正结果
