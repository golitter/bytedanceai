## Why

群聊右侧栏目前只有群成员和历史搜索等静态信息面板。Agent 协作过程中，用户需要实时了解当前工作区的 Git 状态（分支、提交历史、当前分支切换）和终端操作能力。将 Git Graph 可视化面板和交互式终端嵌入右侧栏，可以让用户在群聊上下文中直接观察和操控 Agent 的工作环境，无需切换到外部工具。

## What Changes

- 在群聊右侧栏新增 **Git Graph 面板**，展示提交树（双分支：`main` + 群聊 Agent 分支），支持点击分支标签切换当前分支视图，HEAD 指示点跟随切换
- 在群聊右侧栏新增 **Terminal 面板**，支持常用命令（`git status`、`git log`、`git branch`、`git checkout`、`ls`、`pwd`、`npm run build` 等），终端提示符显示当前分支名
- Git Graph 与 Terminal 联动：切换分支时，Git Graph 的 HEAD 节点、分支标签高亮、终端提示符和 `git status`/`git log`/`git branch` 输出全部同步更新
- 右侧栏区块顺序调整为：群成员 → Git Graph → Terminal

## Capabilities

### New Capabilities
- `git-graph-panel`: 群聊右侧栏 Git Graph 可视化面板，支持双分支提交树渲染、分支标签点击切换、HEAD 节点指示、commit tooltip
- `terminal-panel`: 群聊右侧栏交互式终端面板，支持 `git checkout/switch` 分支切换联动、常用命令模拟、分支感知的输出

### Modified Capabilities
- `right-sidebar`: 右侧栏区块结构变更，新增 Git Graph 和 Terminal 两个可折叠区块，调整区块顺序为 MembersSection → GitGraphPanel → TerminalPanel

## Impact

- **前端组件**：新增 `GitGraphPanel`、`TerminalPanel` 两个 React 组件，修改 `RightSidebar` 容器组件
- **样式**：新增 Git Graph 行渲染、SVG 轨道线、终端窗口等相关 CSS
- **数据流**：需要从后端/Agent 端获取实际 Git 状态数据（当前 Demo 为静态 mock）
- **API**：后续需新增 Git 状态查询 API（本 change 先以 Demo/静态数据实现）
