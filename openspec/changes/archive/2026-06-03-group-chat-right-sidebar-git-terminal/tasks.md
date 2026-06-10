## 1. 准备工作

- [x] 1.1 创建 `GitGraphPanel.tsx` 组件文件骨架（`frontend/src/components/chat/`）
- [x] 1.2 创建 `TerminalPanel.tsx` 组件文件骨架（`frontend/src/components/chat/`）
- [x] 1.3 定义 Git Graph 和 Terminal 相关的 TypeScript 类型（commit 数据结构、分支配置、命令映射类型）

## 2. GitGraphPanel 实现

- [x] 2.1 实现 SVG 轨道线渲染（双 lane 垂直虚线，main 灰色、feature 紫色）
- [x] 2.2 实现 commit 节点和连接线渲染（同 lane 直线、跨 lane 贝塞尔曲线）
- [x] 2.3 实现 commit 行信息渲染（hash、message、author 颜色点、time）
- [x] 2.4 实现分支标签切换功能（点击标签触发 `onBranchChange`，标签样式切换）
- [x] 2.5 实现 HEAD 节点绿色指示（绿色节点 + 半透明光晕 + 行高亮背景）
- [x] 2.6 实现 commit tooltip（hover 显示、leave 隐藏、opacity 过渡动画）
- [x] 2.7 实现区块折叠/展开（复用 section-header 模式，localStorage 持久化）

## 3. TerminalPanel 实现

- [x] 3.1 实现终端窗口布局（标题栏 + 三色圆点 + 路径 + Connected 状态）
- [x] 3.2 实现输出区域（滚动、自动滚到底部、欢迎消息）
- [x] 3.3 实现输入行（分支感知提示符 `(branch) $`、输入框、光标闪烁动画）
- [x] 3.4 实现命令映射表（help、clear、pwd、ls、whoami、npm run build、npm test、echo、cat）
- [x] 3.5 实现 `git status` 命令（分支感知输出，feature 分支显示 modified files，main 分支显示 clean）
- [x] 3.6 实现 `git branch` 命令（列出所有分支，当前分支标 `*`）
- [x] 3.7 实现 `git log` 命令（只显示当前分支可达的 commits，按时间倒序）
- [x] 3.8 实现 `git checkout`/`git switch` 命令（切换分支、触发 `onBranchChange`、更新提示符、处理已存在/不存在分支）
- [x] 3.9 实现区块折叠/展开（Terminal 占满剩余高度 flex:1，localStorage 持久化）

## 4. RightSidebar 集成

- [x] 4.1 在 RightSidebar 中添加 `currentBranch` 状态管理（`useState`，初始值为 feat/group-chat-agent）
- [x] 4.2 集成 GitGraphPanel（传入 commits 数据、currentBranch、onBranchChange）
- [x] 4.3 集成 TerminalPanel（传入 currentBranch、onBranchChange）
- [x] 4.4 调整右侧栏区块顺序为 MembersSection → GitGraphPanel → TerminalPanel
- [x] 4.5 调整右侧栏宽度从 280px 到 300px

## 5. 样式与收尾

- [x] 5.1 添加 GitGraphPanel 相关 CSS（行样式、lane area、tooltip、分支标签）
- [x] 5.2 添加 TerminalPanel 相关 CSS（终端窗口、输出区、输入行、光标动画）
- [x] 5.3 验证 Git Graph ↔ Terminal 双向联动（切换分支时两边同步更新）
- [x] 5.4 验证折叠/展开功能和 localStorage 持久化
