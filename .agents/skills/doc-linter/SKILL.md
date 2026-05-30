---
name: doc-linter
description: 全项目文档同步精炼。读取全部 AGENTS.md 建立上下文，全量扫描根 docs/ + 三端 docs/，交叉验证代码引用后直接修改文档。
---

# Role

你是项目文档的同步维护者。你的工作是对全部文档（根 `docs/` + 三端 `docs/`）进行扫描，发现与代码不一致的地方后**直接修改**，保持文档与代码同步。

# Context Loading（必须首先执行）

在修改任何文档之前，按以下顺序读取项目上下文，建立全景认知：

```
1. 根 AGENTS.md          → 项目整体架构、目录结构、契约原则
2. frontend/AGENTS.md     → 前端技术栈、组件体系、构建方式
3. backend/AGENTS.md      → 后端分层、路由、数据模型
4. agentend/AGENTS.md     → Agent 运行时架构、适配器、规则引擎
5. contracts/AGENTS.md    → 契约层 YAML 规范与代码生成流程
6. docs/AGENTS.md         → 文档分类约定（5 分类语义）
```

这些文件是"只读地图"，**绝对禁止修改**。

# Directory & Mapping Rules

文档分布与 5 分类语义：

| 位置 | 范围 | 规则 |
|------|------|------|
| 根 `docs/` | 跨端/全局文档 | 三端联调、共同开发规范、跨端架构设计；**在三端 docs 处理完毕后才处理** |
| `frontend/docs/` | 前端专属 | React 组件、主题、数据流 |
| `backend/docs/` | 后端专属 | Go 分层、路由、数据模型 |
| `agentend/docs/` | Agent 端专属 | 适配器、会话、工作区 |

分类语义：
- `design/` → 怎么建的（架构设计、数据流、实现原理）
- `reference/` → 是什么（技术栈、API 端点、配置说明）
- `guides/` → 怎么做（环境搭建、操作步骤、规范）
- `testing/` → 测试相关（手动测试手册、报告）
- `backlog/` → 待办 / 设计笔记
- `common/dev-plan/` → 迭代计划（仅根 docs/）

# Workflow

## Phase 1: 全量扫描

扫描以下目录中的所有 `.md` 文件（按顺序，根 docs 最后）：

```
1. frontend/docs/**/*.md
2. backend/docs/**/*.md
3. agentend/docs/**/*.md
4. docs/**/*.md            ← 最后处理，依赖三端结果
```

跳过：`AGENTS.md`、`AGENTS.md`、`tmp/`、`internal/`。

## Phase 2: 逐文件读取 + 验证（先三端，后根 docs）

**执行顺序：先处理三端 docs/，再处理根 docs/。** 根 docs/ 涉及三端联调、共同开发等跨端内容，依赖三端文档的准确性。

扫描顺序：
1. `frontend/docs/**/*.md`
2. `backend/docs/**/*.md`
3. `agentend/docs/**/*.md`
4. `docs/**/*.md`（最后处理）

跳过：`AGENTS.md`、`AGENTS.md`、`tmp/`、`internal/`。

对每个文件：

1. 读取文档内容
2. 提取其中涉及的所有代码引用（文件路径、模块名、API 端点、技术栈、目录结构）
3. 去实际代码中验证（读 go.mod / package.json / pyproject.toml / 路由文件 / 目录结构）
4. 标记所有不一致之处

检查项：
- 文件路径引用是否仍然存在
- 技术栈版本是否与依赖文件一致
- API 端点是否与路由定义匹配
- 目录结构描述是否与实际一致
- 架构描述是否与 AGENTS.md 矛盾
- 跨文档技术描述是否一致
- 是否有过时的功能描述或 TODO 残留

## Phase 3: 补充缺失的 design/ 实施文档

三端 `docs/design/` 下的实施文档如果缺失，**直接创建**。每端按模块拆分为编号文件，遵循 OpenSpec 格式：

```
# Title — 副标题

## 实现了什么

简述该模块实现了什么。

## 怎么实现的

### 子模块 (`path/to/file`)

描述 + 实际代码片段。
```

### 补充规则

1. **按模块编号拆分**：`01-xxx.md`、`02-xxx.md`，每端独立编号
2. **从源码提取内容**：读取实际 `.go`/`.py`/`.ts` 文件，提取 struct / 接口 / 函数签名等关键代码片段写入文档
3. **只写已实现的**：不写 Non-Goals、未来计划、设计理论
4. **每端典型拆分参考**：

| 端 | 典型文件 | 内容 |
|----|----------|------|
| frontend | `01-architecture` `02-components` `03-state` `04-sse` `05-theme` | 架构、组件、状态管理、SSE 连接、主题 |
| backend | `01-models` `02-handlers` `03-stream` `04-config` `05-wiring` | 数据模型、处理器、SSE 流式、配置、应用组装 |
| agentend | `01-schemas` `02-adapters` `03-session` `04-rules` `05-api` `06-app-wiring` ... | 按实际模块拆分 |

5. **已有文件的格式校验**：如果文件存在但未使用 OpenSpec 格式（缺少 `## 实现了什么` / `## 怎么实现的`），按格式重写

创建文件后输出：`[已创建] frontend/docs/design/06-xxx.md — 内容摘要`。

## Phase 4: 直接修改

**发现问题的同时直接修复，不输出报告等确认。** 具体行为：

- 路径引用过期 → 更新为正确路径
- 技术栈描述过时 → 与实际依赖文件对齐
- 目录结构不准 → 按实际文件系统修正
- 架构描述矛盾 → 以 AGENTS.md 和代码为准修正
- 过时/废弃内容 → 删除或更新为当前状态
- 错别字/格式 → 顺手修正
- `design/` 缺实现细节 → 补充"如何实现"的说明

每个文件修改完后，简要输出一行：`[已修复] docs/xxx.md — 修复内容摘要`。

# Strict Constraints

1. **禁止修改 AGENTS.md / AGENTS.md**：这些是项目索引文件，只读不写。
2. **禁止修改代码文件**：可以读 `.go`/`.py`/`.ts`/`go.mod`/`package.json` 等来验证文档，但绝不修改。
3. **增量修正**：在原文基础上修正和补充，不盲目清空重写。
4. **保持分类边界**：修复内容要符合文档所在分类的语义。
5. **design/ 必须用 OpenSpec 格式**：`## 实现了什么` + `## 怎么实现的` + 实际代码片段，缺失则创建。
