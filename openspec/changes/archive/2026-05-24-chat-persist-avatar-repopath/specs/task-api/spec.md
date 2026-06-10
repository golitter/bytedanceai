## ADDED Requirements

### Requirement: 查询 Task 消息历史
系统 SHALL 提供 `GET /api/tasks/:taskId/messages` 端点，返回指定 task 下的所有消息，按 created_at 升序排列。

#### Scenario: 返回消息列表
- **WHEN** 发送 `GET /api/tasks/{taskId}/messages`
- **THEN** 返回 HTTP 200，data 为 Message 数组，按 created_at 升序

#### Scenario: 空消息列表
- **WHEN** task 下无消息
- **THEN** 返回 HTTP 200，data 为空数组 `[]`

#### Scenario: Task 不存在
- **WHEN** task_id 不存在
- **THEN** 返回 HTTP 404

### Requirement: repoPath 校验端点
系统 SHALL 提供 `POST /api/validate-repo-path` 端点，接收 `{"repo_path": "..."}`, 转发给 agentend 校验，返回校验结果。

#### Scenario: 转发校验
- **WHEN** 发送 `POST /api/validate-repo-path` 并传入 `{"repo_path": "/valid/path"}`
- **THEN** backend 转发给 agentend，返回 agentend 的校验结果

#### Scenario: agentend 不可达
- **WHEN** agentend 服务不可用
- **THEN** 返回 HTTP 503，提示服务不可用

## MODIFIED Requirements

### Requirement: 创建 Task（群组任务）
系统 SHALL 提供 `POST /api/tasks` 端点，接收 `{"title": "...", "repo_path": "..."}` JSON body，创建 Task 记录并返回 `task_id`。repo_path 为可选字段。

#### Scenario: 成功创建
- **WHEN** 发送 `POST /api/tasks` 并传入 `{"title": "my-project", "repo_path": "/path/to/repo"}`
- **THEN** 返回 HTTP 201，body 为 `{"code": 0, "data": {"task_id": "uuid", "title": "my-project", "repo_path": "/path/to/repo", "status": "active", ...}}`

#### Scenario: 不指定 repo_path
- **WHEN** 发送 `POST /api/tasks` 并传入 `{"title": "my-project"}`
- **THEN** 返回 HTTP 201，repo_path 为空字符串

#### Scenario: 缺少 title
- **WHEN** 发送 `POST /api/tasks` 不传 title 字段
- **THEN** 返回 HTTP 400，提示参数错误
