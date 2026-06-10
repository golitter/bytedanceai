## ADDED Requirements

### Requirement: API 响应状态校验
所有 API 函数 SHALL 在解析响应体之前检查 `res.ok`。当响应状态码不在 200-299 范围内时，SHALL 抛出结构化错误（包含 status code 和错误消息）。

#### Scenario: 服务端返回 500
- **WHEN** fetchTasks 调用后服务端返回 500
- **THEN** 抛出包含 status=500 和错误消息的 ApiError，而非 TypeError

#### Scenario: 服务端返回 404
- **WHEN** fetchTask 调用后服务端返回 404
- **THEN** 抛出包含 status=404 的 ApiError，调用方可据此显示"未找到"

### Requirement: 统一响应处理辅助函数
`lib/api.ts` SHALL 导出 `handleResponse<T>(res: Response): Promise<T>` 辅助函数，统一执行 `res.ok` 检查和 JSON 解析。所有 API 函数 SHALL 使用此辅助函数处理响应。

#### Scenario: handleResponse 处理成功响应
- **WHEN** handleResponse 收到 ok=true 的 Response
- **THEN** 解析 JSON 并返回 data 字段

#### Scenario: handleResponse 处理失败响应
- **WHEN** handleResponse 收到 ok=false 的 Response
- **THEN** 解析错误消息并抛出 ApiError（含 status、message）

### Requirement: 分页加载错误不静默吞没
ChatArea 的 `loadMoreMessages` 函数 SHALL 在分页请求失败时向用户显示错误提示（toast），不静默吞没异常。

#### Scenario: 分页加载失败
- **WHEN** loadMoreMessages 的 API 调用失败
- **THEN** 显示 toast 错误提示，不只在 catch 中重置 loading 状态
