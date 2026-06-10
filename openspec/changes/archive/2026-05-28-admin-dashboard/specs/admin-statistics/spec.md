## ADDED Requirements

### Requirement: Session trend chart
数据统计页面 SHALL 展示会话数量趋势图，支持按天/按周切换。

#### Scenario: Display daily trend
- **WHEN** 用户打开数据统计页面，默认显示按天视图
- **THEN** 展示最近 7 天的柱状图，每天一根柱子，悬停显示具体数值

#### Scenario: Switch to weekly view
- **WHEN** 用户点击"按周"切换按钮
- **THEN** 图表切换为最近 4 周的柱状图

### Requirement: Message total with agent breakdown
数据统计页面 SHALL 展示消息总量及各 Agent 的消息占比。

#### Scenario: Display message stats
- **WHEN** 数据统计页面加载
- **THEN** 显示消息总量数字，下方显示各 Agent 消息占比的堆叠进度条和图例

### Requirement: Storage trend chart
数据统计页面 SHALL 展示存储占用的近期趋势。

#### Scenario: Display storage trend
- **WHEN** 数据统计页面加载
- **THEN** 显示最近 7 天的存储占用趋势图（面积样式），每列标注 GB 数值

### Requirement: Statistics API
后端 SHALL 提供 `GET /api/admin/statistics` 接口，返回统计数据。

#### Scenario: Fetch statistics
- **WHEN** 前端请求 `GET /api/admin/statistics`
- **THEN** 返回 JSON 包含 dailySessions、weeklySessions、labels、totalMessages、messagesByAgent、storageDays、storageLabels
