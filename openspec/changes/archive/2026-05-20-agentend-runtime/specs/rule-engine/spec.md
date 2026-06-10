## ADDED Requirements

### Requirement: BaseRule abstract interface
系统 SHALL 定义 `BaseRule` 抽象基类，包含字段 `name`（str）、`description`（str）、`phase`（"pre" | "post"）、`priority`（int，默认 0），以及抽象方法 `check` 和 `enforce`。

#### Scenario: Rule implements check and enforce
- **WHEN** 一个 Rule 类继承 `BaseRule`
- **THEN** 该类 MUST 实现 `check(context: dict) -> bool` 和 `enforce(context: dict) -> dict` 方法

### Requirement: RuleEngine evaluate with priority ordering
`RuleEngine` SHALL 按 `priority` 降序执行已注册的 Rule。每个 Rule 先调用 `check`，通过后调用 `enforce` 修改上下文。任何 Rule 的 `check` 返回 `False` SHALL 立即终止评估并返回失败。

#### Scenario: All rules pass
- **WHEN** 3 个 Rule 的 `check` 均返回 `True`，priority 分别为 10、5、0
- **THEN** SHALL 按 priority 10 → 5 → 0 的顺序执行 `enforce`，返回 `(True, merged_context)`

#### Scenario: Rule check fails
- **WHEN** priority=5 的 Rule 的 `check` 返回 `False`
- **THEN** SHALL 停止后续 Rule 评估，返回 `(False, context)`，包含失败 Rule 的名称和原因

### Requirement: RuleEngine enforce output schema
`enforce` 方法的返回值 SHALL 包含以下可选字段：`system_prompt_append`（str）、`allowed_tools`（list[str]）、`max_turns`（int）、`blocked`（bool）、`error`（str）。

#### Scenario: Rule injects system prompt constraint
- **WHEN** SafetyRule 的 `enforce` 返回 `{"system_prompt_append": "禁止执行 rm -rf"}`
- **THEN** RuleEngine SHALL 合并所有 Rule 的 `system_prompt_append` 字段，传递给 Adapter

#### Scenario: Rule restricts allowed tools
- **WHEN** ScopeRule 的 `enforce` 返回 `{"allowed_tools": ["Read", "Write", "Bash"]}`
- **THEN** RuleEngine SHALL 取所有 Rule 的 `allowed_tools` 交集，传递给 Adapter

### Requirement: RuleRegistry dynamic registration
系统 SHALL 提供 `RuleRegistry`，支持动态注册和查找 Rule 实例。启动时 MUST 自动加载内置 Rule。

#### Scenario: Register and list rules
- **WHEN** 注册 `SafetyRule()` 和 `ScopeRule()`
- **THEN** `registry.list()` SHALL 返回两个 Rule 实例的列表

### Requirement: Built-in SafetyRule
系统 MUST 内置 `SafetyRule`（priority=10），`check` 始终返回 `True`，`enforce` 追加安全约束到 system_prompt 并限制危险工具。

#### Scenario: SafetyRule enforce
- **WHEN** SafetyRule 的 `enforce` 被调用
- **THEN** SHALL 返回包含安全提示文本的 `system_prompt_append`，以及排除了 `dangerouslyDisableSandbox` 等危险工具的 `allowed_tools`

### Requirement: Built-in ScopeRule
系统 MUST 内置 `ScopeRule`（priority=5），`check` 验证 workspace_path 是否在允许范围内，`enforce` 注入作用域约束。

#### Scenario: ScopeRule check passes
- **WHEN** 请求的 `workspace_path` 在允许范围内
- **THEN** `check` SHALL 返回 `True`，`enforce` 注入 "只允许修改指定目录下文件" 的约束

#### Scenario: ScopeRule check fails
- **WHEN** 请求的 `workspace_path` 不在允许范围内或为空
- **THEN** `check` SHALL 返回 `False`，`error` 包含 "workspace path not allowed" 信息
