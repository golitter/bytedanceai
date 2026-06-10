## ADDED Requirements

### Requirement: Service health display
服务健康页面 SHALL 展示三端服务（Frontend、Backend、AgentEnd）的运行状态。

#### Scenario: Display service status
- **WHEN** 用户打开服务健康页面
- **THEN** 显示 3 个服务卡片，每个包含：服务名称、运行状态（Running/Down + 绿/红指示灯）、运行时长、版本号、端口号、上次检查时间

#### Scenario: Running service indicator
- **WHEN** 服务状态为 running
- **THEN** 显示绿色脉冲圆点 + "Running" 文字

#### Scenario: Down service indicator
- **WHEN** 服务状态为 down
- **THEN** 显示红色圆点 + "Down" 文字

### Requirement: Manual refresh
用户 SHALL 能手动刷新服务健康状态。

#### Scenario: Click refresh
- **WHEN** 用户点击刷新按钮
- **THEN** 按钮显示加载状态，重新检测所有服务，更新状态和上次检查时间

### Requirement: Service health API
后端 SHALL 提供 `GET /api/admin/services` 接口，返回三端服务状态。

#### Scenario: Fetch service health
- **WHEN** 前端请求 `GET /api/admin/services`
- **THEN** 返回服务列表，每项包含 name、status、uptime、version、port、lastCheck
