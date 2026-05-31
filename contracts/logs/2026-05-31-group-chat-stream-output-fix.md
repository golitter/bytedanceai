# 群聊流式输出修复：重复块去重、结构化失败、摘要卡片

## 变更原因

Orchestrator 编排多 Agent 时，SSE 流中存在以下问题：
1. 同一 `question_id` 的 ask-agent 块、同一 `task_id` 的 runtime 块会重复出现，导致前端渲染重复卡片。
2. 子 Agent 超时或出错时仅以 `[Timeout]`/`[Error]` 纯文本标记嵌入 `TaskResult.content`，前端无法解析为结构化 UI。
3. 子 Agent 运行时产生的 `RUNTIME_TEXT` 进度文本在持久化历史中与最终结果文本重复，刷新后出现双份内容。
4. 最终聚合报告直接复制子 Agent 执行日志全文，没有简洁摘要。

本次变更在 Agentend 端区分 runtime 进度与持久化结果、输出结构化失败元数据、生成简洁摘要；在前端去重重复块、解析 legacy 标记为失败卡片、展示最终摘要卡片。

## 变更文件

本次变更 **未修改** `contracts/schemas/*.yaml`。

`contracts/schemas/event-types.yaml` 中 `StreamEvent.content` 已定义为开放对象（`additionalProperties: true`），以下新增字段在现有 schema 下合法：

- `RuntimeCompleted` 事件 content 新增 `failure_type`、`failure_reason` 字段（结构化失败信息）
- Markdown 内容块通过 `:::` 围栏标记嵌入 JSON 元数据（task-failure、final-summary），属于文本渲染约定，不涉及 schema 变更

## 对比结果

### 变更前

子 Agent 超时时，`TaskResult.content` 中仅包含：

```
[Timeout] claude-code 在执行任务 task-001 时超时
```

聚合输出直接拼接所有子 Agent 日志全文。

前端无法识别重复块，runtime 文本在历史刷新后与最终结果共存。

### 变更后

`RuntimeCompleted` 事件在失败时携带：

```json
{
  "type": "runtime_completed",
  "content": {
    "task_id": "task-001",
    "agent": "claude-code",
    "success": false,
    "failure_type": "timeout",
    "failure_reason": "exceeded 300s limit"
  }
}
```

最终聚合输出为简洁摘要块：

```json
{
  "status": "partial",
  "completed": 3,
  "failed": 1,
  "nextAction": "可尝试重试失败任务",
  "details": [{ "task_id": "...", "agent": "...", "status": "completed|failed", "summary": "..." }]
}
```

前端按 `question_id`/`task_id` 去重；runtime 文本在收到持久化结果后自动清除。

## 跨端影响

- **AgentEnd**: `ExecutionEngine` 区分 runtime 进度与持久化结果；`_task_failure_block` 输出结构化失败元数据；`Aggregator` 生成简洁摘要分离完成/失败任务。
- **Backend**: 无变更，透传 SSE 事件。
- **Frontend**: `block-reducer` 新增 `stripRuntimeStreamingText`、`parseFailureMarker`、按 ID 去重逻辑；新增 `FinalSummaryCard`、`TaskFailureCard` 组件；长卡片支持内部滚动与展开。
- **Contracts**: 无 schema 变更；记录 `RuntimeCompleted` 扩展字段与 markdown 内容块渲染约定。

## 契约变更

`RuntimeCompleted` 事件 content 新增可选字段：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `failure_type` | `string` | 否 | 失败类型，取值 `timeout` 或 `error`，仅 `success=false` 时存在。 |
| `failure_reason` | `string` | 否 | 失败原因描述，仅 `success=false` 时存在。 |

Markdown 内容块渲染约定（非 schema 变更）：

| 标记 | 用途 |
|------|------|
| `:::task-failure` | 结构化任务失败卡片，内嵌 JSON 含 `task_id`、`agent`、`failureType`、`reason` |
| `:::final-summary` | 最终摘要卡片，内嵌 JSON 含 `status`、`completed`、`failed`、`nextAction`、`details` |

兼容性：

- 旧客户端忽略未知字段和未知围栏标记，回退为纯文本渲染，不受影响。
- 新客户端在 `failure_type` 缺失时兼容旧 `[Timeout]`/`[Error]` 纯文本标记。
