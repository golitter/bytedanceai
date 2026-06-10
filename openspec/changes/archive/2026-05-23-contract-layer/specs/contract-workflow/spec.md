## ADDED Requirements

### Requirement: pre-commit 契约变更检测
pre-commit hook SHALL 检测 staged 文件中是否包含契约相关变更。契约相关文件包括：`agentend/src/schemas/`、`agentend/src/api/`、`agentend/src/adapters/`、`backend/internal/model/`、`backend/internal/handler/`、`backend/internal/types/`、`frontend/src/types/`、`frontend/src/api/`、`contracts/`。

#### Scenario: 检测到契约相关变更
- **WHEN** staged 文件包含 `agentend/src/schemas/events.py`
- **THEN** hook SHALL 以 exit 1 阻断提交，并输出引导信息

#### Scenario: 无契约相关变更
- **WHEN** staged 文件仅包含 `agentend/src/adapters/claude.py` 的内部实现修改（不涉及 schema 或 API 变更）
- **THEN** hook SHALL 通过检测，进入 secret 验证阶段

### Requirement: 契约变更引导信息
当检测到契约相关变更时，hook SHALL 输出结构化的引导信息，包含变更文件列表、要求 Agent 执行的步骤（对比 contracts/schemas/、更新契约、检查跨端影响、写变更日志、重新 commit）。

#### Scenario: 引导信息格式
- **WHEN** hook 检测到契约变更
- **THEN** 输出 SHALL 包含以下段落：变更文件列表、对比步骤、contracts/logs/ 写入要求、重新 commit 提示

### Requirement: 契约一致性验证
pre-commit hook SHALL 验证各端代码中的类型定义与 `contracts/schemas/` 一致。当检测到不一致时 SHALL 阻断提交。

#### Scenario: 契约一致通过
- **WHEN** Agent 已更新 `contracts/schemas/event-types.json` 且代码中使用的是生成类型
- **THEN** hook SHALL 通过一致性检查

#### Scenario: 契约不一致阻断
- **WHEN** 代码中定义了新的 EventType 值但 `contracts/schemas/event-types.json` 未同步更新
- **THEN** hook SHALL 以 exit 1 阻断，提示 Agent 更新 contracts/

### Requirement: 变更审计日志
Agent 在处理契约变更时 SHALL 在 `contracts/logs/` 目录写入变更记录文件。文件名格式 SHALL 为 `YYYY-MM-DD-<kebab-case-description>.md`。内容 SHALL 包含变更原因、变更文件、对比结果、跨端影响、契约变更详情。

#### Scenario: 写入变更日志
- **WHEN** Agent 新增了 EventType.THINKING 事件类型并更新了 contracts/
- **THEN** Agent SHALL 创建 `contracts/logs/2026-05-23-add-thinking-event.md`，包含变更原因、修改的文件列表、与原 contracts/ 的对比结果、对 Frontend/Backend 的影响说明

#### Scenario: 日志文件内容完整性
- **WHEN** Agent 写入变更日志
- **THEN** 日志 SHALL 包含以下段落：变更原因、变更文件、对比结果、跨端影响、契约变更

### Requirement: Human secret 门控
git commit SHALL 要求 human 输入预置的 secret 才能完成提交。Agent SHALL NOT 通过任何方式获取或推断该 secret。Secret 验证失败 SHALL 以 exit 1 阻断提交。

#### Scenario: Secret 验证成功
- **WHEN** human 在终端执行 git commit 并输入正确的 secret
- **THEN** 提交 SHALL 成功完成

#### Scenario: Secret 验证失败
- **WHEN** 输入的 secret 与预置值不匹配
- **THEN** SHALL 以 exit 1 阻断提交，提示 secret 不匹配

#### Scenario: Agent 无法绕过 secret
- **WHEN** Agent 执行 git commit
- **THEN** SHALL 触发 secret 验证，Agent 无法提供正确的 secret，提交 SHALL 被阻断

### Requirement: 幂等性检查
pre-commit hook 的契约一致性检查 SHALL 是幂等的——多次执行 SHALL 产生相同结果。只要契约和代码一致就放行，无论执行次数。

#### Scenario: 多次执行结果一致
- **WHEN** 契约和代码已保持一致，连续执行 pre-commit hook 两次
- **THEN** 两次执行 SHALL 都通过一致性检查
