## ADDED Requirements

### Requirement: Git Graph rendering
系统 SHALL 提供 `GitGraphPanel` 组件，使用 SVG 绘制双分支提交树，每行显示一个 commit（hash、message、author 颜色点、time）。

#### Scenario: Render commit tree with two branches
- **WHEN** GitGraphPanel 接收 commits 数据（包含 main 和 feat/group-chat-agent 两个 lane）
- **THEN** 渲染 SVG 轨道线（两条垂直虚线，main 灰色、feature 紫色），同 lane commit 用直线连接，跨 lane commit 用贝塞尔曲线连接

#### Scenario: Render commit row info
- **WHEN** 每个 commit 渲染
- **THEN** 显示 hash（等宽字体、灰色）、commit message（白色）、author 颜色圆点、time（灰色），message 超长时省略号截断

### Requirement: Branch label switching
GitGraphPanel 顶部 SHALL 显示分支标签（当前分支紫色高亮 `--primary-soft`，其他分支灰色），点击标签 SHALL 切换当前分支。

#### Scenario: Click to switch branch
- **WHEN** 用户点击非当前分支的标签（如当前在 feat/group-chat-agent，点击 main）
- **THEN** main 标签变为 current 样式（紫色），feat/group-chat-agent 变为 default 样式（灰色），HEAD 节点移动到 main 分支最新 commit，触发 `onBranchChange` 回调

#### Scenario: Click current branch label
- **WHEN** 用户点击当前已激活的分支标签
- **THEN** 不触发任何变更

### Requirement: HEAD node indicator
当前分支的最新 commit 节点 SHALL 以绿色（`--color-success`）显示，外围带半透明光晕。

#### Scenario: HEAD node on feature branch
- **WHEN** 当前分支为 feat/group-chat-agent
- **THEN** 该分支最后一个 commit 节点为绿色，半径 4.5px，外围 8px 半透明绿色光晕，对应行背景带 `rgba(99,102,241,0.06)` 高亮

#### Scenario: HEAD node moves on switch
- **WHEN** 用户切换分支到 main
- **THEN** HEAD 绿色指示点移动到 main 分支最新 commit，原 HEAD 恢复为分支颜色

### Requirement: Commit tooltip
鼠标 hover 到 commit 节点或行时 SHALL 显示 tooltip，包含 hash、message、author、lane、time。

#### Scenario: Hover shows tooltip
- **WHEN** 鼠标进入 commit 节点或对应行
- **THEN** 在行上方显示浮动 tooltip，内容为 hash（紫色等宽）、message（白色）、author · lane · time（灰色）

#### Scenario: Mouse leave hides tooltip
- **WHEN** 鼠标离开 commit 节点或行
- **THEN** tooltip 淡出隐藏（opacity 过渡 120ms）

### Requirement: Git Graph section collapsible
Git Graph 区块 SHALL 支持折叠/展开，复用 RightSidebar 的 section-header + chevron 模式。

#### Scenario: Toggle Git Graph section
- **WHEN** 用户点击 Git Graph 标题栏
- **THEN** 区块折叠/展开，chevron 图标旋转，折叠状态持久化到 localStorage
