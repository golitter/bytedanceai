## ADDED Requirements

### Requirement: 创建 Session
系统 SHALL 提供 `POST /api/sessions` 端点，接收 `{"title": "..."}` JSON body，创建 Session 记录并返回 `session_id`。

#### Scenario: 成功创建
- **WHEN** 发送 `POST /api/sessions` 并传入 `{"title": "test"}`
- **THEN** 返回 HTTP 201，body 为 `{"code": 0, "data": {"session_id": "uuid", "title": "test", "status": "active", ...}}`

#### Scenario: 缺少 title
- **WHEN** 发送 `POST /api/sessions` 不传 title 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: 查询 Session 列表
系统 SHALL 提供 `GET /api/sessions` 端点，返回所有 Session 记录。

#### Scenario: 返回列表
- **WHEN** 发送 `GET /api/sessions`
- **THEN** 返回 HTTP 200，data 为 Session 数组

### Requirement: 查询 Session 详情
系统 SHALL 提供 `GET /api/sessions/:id` 端点，返回指定 Session。

#### Scenario: Session 存在
- **WHEN** 发送 `GET /api/sessions/{session_id}`
- **THEN** 返回 HTTP 200，data 为该 Session 详情

#### Scenario: Session 不存在
- **WHEN** 发送 `GET /api/sessions/{不存在的_id}`
- **THEN** 返回 HTTP 404

### Requirement: 删除 Session
系统 SHALL 提供 `DELETE /api/sessions/:id` 端点，软删除指定 Session。

#### Scenario: 成功删除
- **WHEN** 发送 `DELETE /api/sessions/{session_id}`
- **THEN** 返回 HTTP 200，该 Session 从数据库删除
