## Why

三端（Frontend/Backend/AgentEnd）由不同 Agent 独立开发时，跨端协议（事件类型、API schema、SSE 格式）的变更无法被其他端感知，导致架构漂移（drift）。OpenSpec specs/ 定义了协议语义但缺少代码层的强制执行链路。需要建立独立的契约层（contracts/），通过代码生成 + pre-commit 引导 + human secret 门控，将"希望别人知道"变成"不更新就无法工作"。

## What Changes

- 新建 `contracts/` 目录，作为三端共享契约的单一来源（Single Source of Truth）
- 将现有的跨端协议定义（EventType、AgentRequest、StreamEvent、SessionState、错误码等）从各端代码逆向提取到 `contracts/schemas/`
- 建立 `make generate` 管线，从 contracts/schemas/ 自动生成 TypeScript / Go / Python 类型文件到各端的 `generated/` 目录
- 新增 `.husky/pre-commit` hook，检测契约相关文件变更时输出引导信息，要求 Agent 对比 contracts/ 并在 `contracts/logs/` 写下变更记录
- commit 需要 human secret 验证，Agent 无法自主完成提交
- 三端代码改为 import 生成的类型，替代手写定义

## Capabilities

### New Capabilities
- `contract-schema`: 契约 schema 定义与代码生成管线——定义 contracts/schemas/ 的结构、格式，以及 make generate 如何从 schema 生成三端类型
- `contract-workflow`: 契约变更工作流——pre-commit hook 的变更检测逻辑、引导信息格式、contracts/logs/ 审计记录规范、human secret 门控机制

### Modified Capabilities
- `event-types`: EventType / StreamEvent 定义源从 agentend 代码迁移到 contracts/schemas/，agentend 改为 import 生成代码
- `stream-protocol`: SSE 事件格式定义源迁移到 contracts/schemas/，三端统一从生成代码获取

## Impact

- **新增目录**: `contracts/`（schemas/、logs/）、各端 `generated/` 目录
- **修改文件**: agentend/src/schemas/events.py、agentend/src/schemas/request.py、agentend/src/schemas/response.py（改为 import 生成代码）
- **新增文件**: .husky/pre-commit、scripts/generate.py（或类似生成脚本）
- **构建依赖**: 新增 make generate 步骤，各端需配置 generated/ 目录的导入路径
- **开发流程**: Agent 工作流增加契约对比环节，commit 需要 human 介入
