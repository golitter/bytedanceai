## MODIFIED Requirements

### Requirement: API 错误响应格式统一
所有后端 HTTP 错误响应 SHALL 使用 `vo.BadRequest()`、`vo.NotFound()`、`vo.InternalError()` 等辅助函数，确保响应体统一为 `{"code": <int>, "msg": "<string>", "data": null}` 格式。SSE 流式错误响应保持 SSE 事件格式不受影响。

#### Scenario: 非流式端点错误返回标准 JSON
- **WHEN** 任何非 SSE 端点处理请求时发生业务错误
- **THEN** 响应体包含 `code`、`msg`、`data` 三个字段，`code` 非 0 表示错误

#### Scenario: SSE 端点初始阶段错误返回标准 JSON
- **WHEN** SSE 流尚未建立（连接验证阶段）时发生错误
- **THEN** 返回标准 JSON 格式错误响应（使用 vo.Response 辅助函数）

#### Scenario: SSE 流已建立后的错误保持 SSE 格式
- **WHEN** SSE 流已经向客户端发送数据后发生错误
- **THEN** 通过 SSE event 通知客户端错误，不切换为 JSON 响应
