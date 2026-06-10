## MODIFIED Requirements

### Requirement: 契约生成器处理所有已定义 schema
`scripts/generate_contracts.py` 的 `OUTPUT_MAP` SHALL 包含 `contracts/schemas/` 目录下所有 YAML schema 文件的映射，确保无遗漏。当前缺失的 `message` 和 `validate-repo-path` SHALL 被添加。

#### Scenario: 所有 schema 文件都被生成器处理
- **WHEN** `contracts/schemas/` 中新增 YAML schema 文件
- **THEN** 该 schema SHALL 在 `OUTPUT_MAP` 中有对应条目，`make generate` 会为其生成三端类型

#### Scenario: 新增 schema 后 make generate 成功
- **WHEN** 在 OUTPUT_MAP 中添加 message 和 validate-repo-path 后执行 `make generate`
- **THEN** 三端 generated/ 目录中出现对应的类型文件，无报错
