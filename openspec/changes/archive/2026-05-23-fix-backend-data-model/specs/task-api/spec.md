## ADDED Requirements

### Requirement: 创建 Task（群组任务）
系统 SHALL 提供 `POST /api/tasks` 端点，接收 `{"title": "...", "repo_path": "..."}` JSON body，创建 Task 记录并返回 `task_id`。

#### Scenario: 成功创建
- **WHEN** 发送 `POST /api/tasks` 并传入 `{"title": "my-project", "repo_path": "/path/to/repo"}`
- **THEN** 返回 HTTP 201，body 为 `{"code": 0, "data": {"task_id": "uuid", "title": "my-project", "repo_path": "/path/to/repo", "status": "active", ...}}`

#### Scenario: 缺少 title
- **WHEN** 发送 `POST /api/tasks` 不传 title 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: 查询 Task 列表
系统 SHALL 提供 `GET /api/tasks` 端点，返回所有 Task 记录。

#### Scenario: 返回列表
- **WHEN** 发送 `GET /api/tasks`
- **THEN** 返回 HTTP 200，data 为 Task 数组

### Requirement: 查询 Task 详情
系统 SHALL 提供 `GET /api/tasks/:taskId` 端点，返回指定 Task 及其关联的 Sessions。

#### Scenario: Task 存在
- **WHEN** 发送 `GET /api/tasks/{task_id}`
- **THEN** 返回 HTTP 200，data 为该 Task 详情，包含关联的 Session 列表

#### Scenario: Task 不存在
- **WHEN** 发送 `GET /api/tasks/{不存在的_id}`
- **THEN** 返回 HTTP 404

### Requirement: 删除 Task
系统 SHALL 提供 `DELETE /api/tasks/:taskId` 端点，删除指定 Task。

#### Scenario: 成功删除
- **WHEN** 发送 `DELETE /api/tasks/{task_id}`
- **THEN** 返回 HTTP 200，该 Task 从数据库删除

#### Scenario: Task 不存在
- **WHEN** 发送 `DELETE /api/tasks/{不存在的_id}`
- **THEN** 返回 HTTP 404
