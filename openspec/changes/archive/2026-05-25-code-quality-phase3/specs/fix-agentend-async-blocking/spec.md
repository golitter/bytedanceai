## ADDED Requirements

### Requirement: pin_memory 同步 LLM 调用改为异步
`src/orchestrator/pin_memory.py` 的 `_generate_summary` 函数 SHALL 使用 `await llm.ainvoke()` 替代 `llm.invoke()`，并将其改为 `async` 函数。调用方 `pin()` 和 `pin_existing()` SHALL 使用 `await` 调用。

#### Scenario: pin 不阻塞事件循环
- **WHEN** `pin()` 方法执行 LLM 摘要生成
- **THEN** 使用 `await llm.ainvoke()` 而非阻塞的 `llm.invoke()`，不阻塞 asyncio 事件循环

### Requirement: aggregator 同步 LLM 调用改为异步
`src/orchestrator/aggregator.py` 的 `aggregate` 函数 SHALL 使用 `await llm.ainvoke()` 替代 `llm.invoke()`，并将其改为 `async` 函数。调用方 SHALL 使用 `await` 调用。

#### Scenario: aggregate 不阻塞事件循环
- **WHEN** `aggregate()` 执行 LLM 结果聚合
- **THEN** 使用 `await llm.ainvoke()` 而非阻塞调用

### Requirement: _extract_json 添加错误处理
`src/orchestrator/graph.py` 的 `_extract_json` 函数 SHALL 捕获 `json.JSONDecodeError`，在解析失败时返回 None 而非抛出异常。

#### Scenario: LLM 返回非法 JSON
- **WHEN** LLM 返回的内容无法解析为 JSON
- **THEN** `_extract_json` 返回 None，不抛出异常

### Requirement: plan_node 添加 LLM 错误处理
`src/orchestrator/graph.py` 的 `plan_node` 函数 SHALL 捕获 LLM 调用异常和 Pydantic 验证异常，在失败时返回 fallback 状态。

#### Scenario: LLM 调用网络失败
- **WHEN** LLM API 调用因网络错误失败
- **THEN** `plan_node` 返回 `{"plan": None}` 而非崩溃

#### Scenario: LLM 返回不符合 PlanOutput 模型的 JSON
- **WHEN** LLM 返回的 JSON 无法通过 `PlanOutput.model_validate` 验证
- **THEN** `plan_node` 返回 `{"plan": None}` 而非崩溃
