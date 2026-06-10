## Context

当前三端（Frontend React/TS、Backend Go/Gin、AgentEnd Python/FastAPI）在同一个 monorepo 中独立开发。AgentEnd 已有完整的 SSE streaming 实现（EventType 枚举、StreamEvent 模型、3 种 Adapter），但 Frontend 和 Backend 仍为空壳。

跨端协议的类型定义散落在 `agentend/src/schemas/` 中，其他两端无法引用。OpenSpec specs/ 定义了协议语义但缺少代码执行链路。当 Agent 在一端修改协议时，其他端无法感知变更，导致架构漂移。

## Goals / Non-Goals

**Goals:**
- 建立 `contracts/` 作为三端共享契约的 Single Source of Truth
- 从 JSON Schema 自动生成 TypeScript / Go / Python 类型文件
- pre-commit hook 引导 Agent 在契约变更时进行对比、更新、写日志
- human secret 门控确保 Agent 无法自主提交
- 迁移 AgentEnd 现有类型定义到 contracts/ 源

**Non-Goals:**
- 不做硬编码检测（由类型系统自行兜底）
- 不做运行时 schema 验证中间件
- 不做跨端集成测试
- 不做 OpenAPI / protobuf 替换（当前用 JSON Schema）
- 不改造 OpenSpec 系统本身

## Decisions

### Decision 1: JSON Schema 作为契约格式
**选择**: JSON Schema (Draft 2020-12)
**替代方案**: TypeScript .ts 源文件、Python Pydantic 源文件、Protobuf、OpenAPI
**理由**: JSON Schema 是语言中立的，Python/TypeScript/Go 都有成熟的 schema-to-code 工具链。不需要引入 protobuf 的编译依赖，也不需要选某一端语言作为"权威源"。

### Decision 2: 生成脚本用 Python
**选择**: `scripts/generate_contracts.py`
**替代方案**: Node.js 脚本、Makefile 内联、Go 程序
**理由**: 项目已有 Python 运行时（AgentEnd），Python 有 `jsonschema2pydantic`、`datamodel-code-generator` 等成熟库。生成脚本本身不需要跨语言。

### Decision 3: 生成文件提交到 git
**选择**: `generated/` 目录下的文件提交到 git，不 gitignore
**替代方案**: gitignore 生成文件，CI 时生成
**理由**: 提交到 git 使得变更在 `git diff` 中可见，Agent 和人类都能直接看到契约变更影响了哪些类型。CI 只需做一致性校验（`make generate && git diff --exit-code`）。

### Decision 4: pre-commit hook 分两阶段
**选择**: Stage 1 检测 + 引导 → Stage 2 secret 验证
**替代方案**: 单阶段全检查、仅依赖 CI
**理由**: 两阶段使得 Agent 可以完成第一阶段（契约对比和更新），但被 secret 阻断后交给 human。如果只有一个阶段，要么 Agent 完全无法工作，要么 human 没有门控点。

### Decision 5: contracts/logs/ 作为审计记录
**选择**: Agent 每次契约变更写入 `contracts/logs/YYYY-MM-DD-<desc>.md`
**替代方案**: git log 中记录、OpenSpec 变更日志、单独的 CHANGELOG.md
**理由**: contracts/logs/ 与 contracts/schemas/ 同目录，Agent 和人类都能快速定位。格式化的 markdown 比 git commit message 更结构化。未来 Agent 可以读 logs/ 理解契约演化历史。

### Decision 6: 契约检测的文件范围
**选择**: 按目录模式匹配（`agentend/src/schemas/**`、`agentend/src/api/**` 等）
**替代方案**: AST 分析、import 图追踪、手动标记
**理由**: 目录模式匹配简单可靠，不需要解析代码。项目已有清晰的目录结构约定。误触发的代价低（只是多跑一次对比），漏触发的代价高。

## Risks / Trade-offs

**[生成脚本维护成本]** → 生成脚本是一次性投入，后续只需维护 JSON Schema。如果工具链需要升级，只改脚本不改契约。Mitigation: 生成脚本保持简单，只做 schema → types 的直接映射。

**[pre-commit hook 降低 commit 速度]** → 契约变更时 Agent 需要执行完整的对比-更新-写日志流程，可能增加 2-3 分钟。Mitigation: 非契约变更直接通过，不受影响。契约变更的额外时间换取的是防止跨端漂移。

**[JSON Schema 的表达能力有限]** → 复杂的业务规则（如 SessionState 的合法转换）不容易用 JSON Schema 表达。Mitigation: 这些规则保留在各端代码中，contracts/ 只管"形状"（字段、类型、枚举值）。

**[Agent 可能不遵守引导信息]** → Agent 看到引导信息后可能跳过某些步骤。Mitigation: 一致性检查是幂等的，不完成就无法通过。Secret 验证确保 human 做最终检查。

**[contracts/ 目录膨胀]** → logs/ 可能随时间积累大量文件。Mitigation: logs/ 文件是 markdown，体积小。如果需要，后续可以按月归档。
