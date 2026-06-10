## MODIFIED Requirements

### Requirement: RuleRegistry dynamic registration
系统 SHALL 提供 `RuleRegistry`，支持动态注册和查找 Rule 实例。启动时 MUST 自动加载内置 Rule。内置 Rule 列表 SHALL 包含：SafetyRule（priority=10）、ScopeRule（priority=5）、SoulRule、GroupChatRule（priority=6）。

#### Scenario: Register and list rules
- **WHEN** 注册 `SafetyRule()` 和 `ScopeRule()` 和 `GroupChatRule()`
- **THEN** `registry.list()` SHALL 返回三个 Rule 实例的列表

#### Scenario: GroupChatRule included in built-in loading
- **WHEN** 系统启动并自动加载内置 Rule
- **THEN** RuleRegistry SHALL 包含 GroupChatRule 实例（priority=6，介于 SafetyRule 和 ScopeRule 之间）
