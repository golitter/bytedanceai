## Context

群聊右侧栏当前包含 MembersSection、HistorySearch、AnnouncementsSection 等面板。Demo 原型（`docs/common/dev-plan/phase5-2-chat-enhanced/terminal-git-demo.html`）已验证了 Git Graph + Terminal 的交互原型：双分支提交树、点击切换分支、终端命令联动。本 change 将该原型转化为 React 组件集成到现有 RightSidebar 中。

当前右侧栏容器组件 `RightSidebar` 已支持可折叠区块模式（section-header + section-body + chevron），新增面板复用此模式。

## Goals / Non-Goals

**Goals:**
- 在群聊右侧栏集成 Git Graph 面板，展示提交树可视化
- 在群聊右侧栏集成 Terminal 面板，支持常用命令
- 实现 Git Graph ↔ Terminal 双向联动（分支切换同步）
- 复用现有 RightSidebar 折叠区块模式

**Non-Goals:**
- 不实现真实 Git 操作（本阶段使用静态/mock 数据，后续接入后端 API）
- 不实现 Terminal 的完整 Shell 模拟（仅支持预定义命令集）
- 不处理多工作区场景（一个群聊对应一个工作区）
- 不实现 WebSocket 实时推送（终端输出为同步模拟）

## Decisions

### Decision 1: 分支状态管理位置
**选择**：在 `RightSidebar` 组件层级使用 React `useState` 管理 `currentBranch` 状态。

**理由**：Git Graph 和 Terminal 都需要读写当前分支，RightSidebar 是它们的共同父组件，状态提升到此处可避免 props drilling。暂不需要全局 store（如 Zustand），后续接入真实 API 时再考虑。

**备选**：
- Zustand store：过度设计，当前只有两个面板消费此状态
- Context：可用但不必要，两个消费者都在同一父组件下

### Decision 2: Git Graph 渲染方式
**选择**：SVG 绘制轨道线和节点，HTML div 绘制每行 commit 信息。

**理由**：Demo 验证了这种混合方式可行。SVG 负责贝塞尔曲线连接，HTML 负责 commit 信息的布局和交互（hover、tooltip），各自发挥优势。

**备选**：
- 纯 SVG：commit 信息文本排版不如 HTML 灵活
- Canvas：性能好但交互复杂度高，commit 数量（<50）不需要

### Decision 3: Terminal 命令处理方式
**选择**：前端命令映射表（Record<string, () => string>），匹配命令后返回格式化输出。

**理由**：Demo 阶段不需要真实 Shell。命令映射表简洁直观，易于扩展。`git checkout` 命令触发分支状态更新。

**备选**：
- 调用后端 API 执行真实命令：本阶段不需要，后续 Phase 6 实现
- xterm.js：重量级，Demo 阶段过度

### Decision 4: 组件拆分策略
**选择**：三个独立组件文件：
- `GitGraphPanel.tsx`：Git Graph 面板（含 SVG 渲染逻辑）
- `TerminalPanel.tsx`：终端面板（含命令处理逻辑）
- 修改 `RightSidebar.tsx`：集成两个新面板

**理由**：各面板职责清晰、独立。GitGraphPanel 接收 `currentBranch` + `onBranchChange` props；TerminalPanel 接收 `currentBranch` + `onBranchChange` props。

## Risks / Trade-offs

- **[风险] 静态数据与后续 API 接入的差异** → 通过 Props 传入 commits 数据源，后续替换 mock 为 API 调用只需修改数据获取层，组件无需改动
- **[风险] 右侧栏高度有限，三个面板同时展开可能溢出** → Terminal 面板设置 `flex:1` 占满剩余空间，Git Graph 设置 `max-height` 限制高度并支持滚动
- **[权衡] Terminal 功能有限，用户可能期望完整 Shell** → 在终端面板标注 "Demo Mode"，后续接入真实 PTY 后体验对齐
