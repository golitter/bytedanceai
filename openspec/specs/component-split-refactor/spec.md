## ADDED Requirements

### Requirement: MessageBubble 按职责拆分
MessageBubble（当前 390 行）SHALL 拆分为三个文件，主文件 ≤160 行：
- `BlockRenderer.tsx` — 消息块类型渲染（约 110 行）
- `AgentMessageContent.tsx` — Agent 消息内容 + 缩放功能（约 125 行）
- `MessageBubble.tsx` — 主组件，编排用户/Agent/系统气泡（约 155 行）

拆分后的子组件 SHALL NOT 被外部模块导入。

#### Scenario: BlockRenderer 独立渲染消息块
- **WHEN** MessageBubble 需要渲染消息块列表
- **THEN** 委托给 BlockRenderer 组件，传入 blocks 数组和 agentSessionLookup

#### Scenario: AgentMessageContent 处理缩放
- **WHEN** 用户点击 Agent 消息的放大按钮
- **THEN** AgentMessageContent 打开 Dialog 展示完整消息内容

#### Scenario: 主组件仅做编排
- **WHEN** MessageBubble 接收消息 props
- **THEN** 根据 role 类型委托给对应子组件，自身不包含渲染逻辑

### Requirement: RightSidebar 按职责拆分
RightSidebar（当前 471 行）SHALL 拆分为六个文件，主文件 ≤230 行：
- `useCollapsible.ts` — localStorage 折叠状态 Hook（约 22 行）
- `AgentInfoSection.tsx` — Agent 信息展示区（约 57 行）
- `SidebarPathSection.tsx` — 仓库/任务路径展示（约 58 行）
- `SidebarActions.tsx` — 导出/置顶/删除操作（约 56 行）
- `RightSidebar.tsx` — 主组件，编排各区域（约 225 行）

#### Scenario: AgentInfoSection 独立渲染 Agent 详情
- **WHEN** RightSidebar 需要展示 Agent 信息
- **THEN** 委托给 AgentInfoSection，传入 agent 数据和折叠状态

#### Scenario: SidebarActions 独立处理操作按钮
- **WHEN** 用户点击导出/置顶/删除按钮
- **THEN** SidebarActions 处理交互逻辑，通过回调通知主组件

### Requirement: GitGraphPanel 按职责拆分
GitGraphPanel（当前 322 行）SHALL 拆分为四个文件，主文件 ≤170 行：
- `GraphRenderer.tsx` — SVG 图形渲染（约 90 行）
- `GraphBranchLabels.tsx` — 分支标签列表（约 42 行）
- `GraphTooltip.tsx` — 悬浮提示框（约 20 行）
- `GitGraphPanel.tsx` — 主组件，编排状态和交互（约 170 行）

#### Scenario: GraphRenderer 独立渲染 SVG
- **WHEN** GitGraphPanel 需要渲染 Git 图
- **THEN** 委托给 GraphRenderer，传入节点数据和交互回调

#### Scenario: GraphBranchLabels 显示分支名
- **WHEN** GitGraphPanel 需要显示分支标签列表
- **THEN** 委托给 GraphBranchLabels，传入分支数据和当前分支标识

### Requirement: 拆分后功能无回归
所有拆分后的组件 SHALL 保持与拆分前完全相同的渲染输出和交互行为。SHALL NOT 修改 props 接口、状态逻辑或事件处理。

#### Scenario: 拆分后 MessageBubble 渲染一致
- **WHEN** 对同一组消息数据渲染
- **THEN** 拆分前后产生完全相同的 DOM 结构

#### Scenario: 拆分后 RightSidebar 交互一致
- **WHEN** 用户在拆分后的 RightSidebar 中执行导出/置顶/删除
- **THEN** 行为与拆分前完全一致
