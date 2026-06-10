## ADDED Requirements

### Requirement: System resource display
总览仪表盘 SHALL 展示磁盘、内存、Redis 三项系统资源的使用率和进度条。

#### Scenario: Display resource usage
- **WHEN** 用户打开总览仪表盘
- **THEN** 显示磁盘（已用/总量 GB + 百分比 + 进度条）、内存（已用/总量 GB + 百分比 + 进度条）、Redis（已用/总量 MB + 百分比 + 进度条）

#### Scenario: Resources load on click only
- **WHEN** 用户点击总览仪表盘菜单项
- **THEN** 系统请求数据并展示，不做自动刷新或轮询

#### Scenario: Manual refresh
- **WHEN** 用户点击刷新按钮
- **THEN** 重新请求资源数据并更新展示

### Requirement: Key metrics display
总览仪表盘 SHALL 展示 4 个关键指标：总会话数、活跃 Agent 数、工作区数、今日消息数。

#### Scenario: Display metrics
- **WHEN** 用户打开总览仪表盘
- **THEN** 显示 4 个指标卡片，每个包含数字和标签

### Requirement: Resource API
后端 SHALL 提供 `GET /api/admin/resources` 接口，返回磁盘、内存、Redis 使用信息。

#### Scenario: Fetch resources
- **WHEN** 前端请求 `GET /api/admin/resources`
- **THEN** 返回 JSON 包含 `{ disk: { used, total, unit }, memory: { used, total, unit }, redis: { used, total, unit } }`
