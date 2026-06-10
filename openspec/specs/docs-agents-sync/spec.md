# docs-agents-sync

项目文档（6 个 AGENTS.md + 根 docs/ + 三端 docs/）与实际代码、依赖、构建配置保持一致的同步约束。规则来源：仓库内 `.agents/skills/agentsmd-linter` 与 `.agents/skills/doc-linter` 两个 lint skill。

## Purpose

文档与代码漂移的典型来源：新增能力后未同步更新 AGENTS.md 目录树与 details.md 索引；重构后文件路径变化；依赖版本升级但 reference/tech-stack.md 仍写旧版本；`docs/design/` 下历史文档未按 OpenSpec 格式书写；AGENTS.md 内联了文档索引而非一行引用 details.md。本 capability 收口这些约束，让任何 lint pass 都可基于统一规则验证与修复。

## Requirements

### Requirement: AGENTS.md 目录树 SHALL 与实际文件系统一致

每个 AGENTS.md（root / frontend / backend / agentend / docs / contracts）中列出的目录树 MUST 与 `ls` / `find` 实际结果一致：列出的目录和文件均存在；未列出的重要新增（如 `src/` 下新子目录、新增能力对应的目录）MUST 补充并附准确注释。

#### Scenario: AGENTS.md 目录树无悬空条目
- **WHEN** 执行 `find` 验证某 AGENTS.md 中列出的任一目录或文件
- **THEN** 该目录或文件实际存在于文件系统

#### Scenario: 新增重要目录被记录
- **WHEN** 某 `src/` 子目录或跨端可见的功能目录在文件系统中存在且与项目能力相关
- **THEN** 对应 AGENTS.md 的目录树包含该条目并附描述其用途的注释

### Requirement: AGENTS.md 中的 Makefile 命令 SHALL 与实际 Makefile 对齐

AGENTS.md 中描述的 `make` 命令 MUST 与根 Makefile 中的实际 target 一致；新增 target MUST 补充，已废弃 target MUST 移除。

#### Scenario: AGENTS.md 命令可执行
- **WHEN** 在仓库根目录执行 AGENTS.md 中列出的任一 `make` 命令
- **THEN** 该 target 存在于 Makefile 且不会因 "No rule to make target" 失败

### Requirement: AGENTS.md 文件引用 SHALL 全部可达

AGENTS.md 中引用的任何相对路径（如 `docs/AGENTS.md`、`docs/guides/git-conventions.md`）MUST 指向实际存在的文件。

#### Scenario: 文件引用无 404
- **WHEN** 点击 AGENTS.md 中的任一 markdown 链接或路径
- **THEN** 目标文件存在

### Requirement: AGENTS.md SHALL 仅一行引用 details.md 不内联展开文档索引

AGENTS.md MUST NOT 内联展开 design / reference / guides / testing 文档的完整索引；统一以一行 `详见 [docs/reference/details.md](docs/reference/details.md)` 形式指向 details.md。

#### Scenario: AGENTS.md 不展开文档索引
- **WHEN** 读取任一子项目 AGENTS.md
- **THEN** 文档索引仅出现在对应 `docs/reference/details.md`，AGENTS.md 只保留一行引用

### Requirement: details.md 索引 SHALL 覆盖 docs/ 下所有文档

三端 `docs/reference/details.md` MUST 索引本端 `docs/` 下所有 `.md` 文件（`tmp/`、`internal/` 除外），分类符合 5 分类语义（design / reference / guides / testing / backlog）。

#### Scenario: details.md 无遗漏
- **WHEN** 扫描某端 `docs/**/*.md`（排除 `tmp/`、`internal/`）
- **THEN** 每个文件都在对应 `docs/reference/details.md` 中有索引条目，且分类正确

### Requirement: AGENTS.md 与 details.md SHALL 不超行数上限

每个 AGENTS.md MUST ≤ 70 行；每个 `details.md` MUST ≤ 150 行。超出时 MUST 通过压缩目录树注释 / 合并短段落 / 精简命令块 / 删减冗余链接等手段压回上限内，且不损失定位信息。

#### Scenario: AGENTS.md 行数符合上限
- **WHEN** 执行 `wc -l <AGENTS.md>`
- **THEN** 输出行数 ≤ 70

#### Scenario: details.md 行数符合上限
- **WHEN** 执行 `wc -l <details.md>`
- **THEN** 输出行数 ≤ 150

### Requirement: 文档中代码引用 SHALL 与当前代码一致

根 `docs/` 与三端 `docs/` 下所有 `.md` 文件中的代码引用（文件路径、模块名、API 端点、技术栈版本、目录结构描述）MUST 与当前代码、`go.mod` / `package.json` / `pyproject.toml` / 路由定义一致。

#### Scenario: 文件路径可达
- **WHEN** 文档中引用了某个代码文件路径（如 `src/components/Foo.tsx`）
- **THEN** 该路径在文件系统中存在

#### Scenario: 技术栈版本对齐
- **WHEN** 文档中描述某依赖的版本（如 React、Go、Python）
- **THEN** 版本号与对应 `package.json` / `go.mod` / `pyproject.toml` 中的实际版本一致

#### Scenario: API 端点匹配路由定义
- **WHEN** 文档中描述后端 API 端点（如 `GET /api/sessions`）
- **THEN** 该端点在 backend 路由注册中存在且方法一致

### Requirement: 三端 docs/design/ SHALL 按 OpenSpec 格式书写

三端 `docs/design/` 下的实施文档 MUST 包含 `## 实现了什么` + `## 怎么实现的` 两个章节，后者 MUST 包含实际代码片段（struct / 接口 / 函数签名等）。已实现但缺文档的模块 MUST 按编号（`NN-xxx.md`）补齐。

#### Scenario: 已有 design 文档符合格式
- **WHEN** 读取三端 `docs/design/` 下任一 `.md`
- **THEN** 同时包含 `## 实现了什么` 与 `## 怎么实现的` 章节

#### Scenario: 缺失 design 文档被补齐
- **WHEN** 三端某已实现模块在 `docs/design/` 下无对应文档
- **THEN** 后续 lint pass 在归档前补齐对应 `NN-<module>.md`，且符合上述格式

### Requirement: 文档间技术描述 SHALL 不互相矛盾

跨文档（包括跨端、跨分类）对同一技术对象的描述 MUST 一致；如与 AGENTS.md 冲突，以 AGENTS.md + 实际代码为准修正文档。

#### Scenario: 跨文档技术描述一致
- **WHEN** 同一架构对象（如会话管理、SSE 协议、契约生成流程）在多份文档中被描述
- **THEN** 各文档对其职责边界、模块构成、对外接口的描述一致

### Requirement: 文档 SHALL 无过时 TODO / 废弃残留

文档中残留的已实现 TODO、已废弃功能描述、过时的"未来计划"MUST 清理或更新为当前状态。

#### Scenario: TODO 残留清理
- **WHEN** 扫描文档发现 TODO / FIXME 标记
- **THEN** 标记关联的事项要么已完成（移除 TODO），要么仍有意义（保留并等待后续 change 处理）

### Requirement: 文档 lint pass SHALL 严格不动代码与配置

任意 lint pass（基于 `agentsmd-linter` / `doc-linter` skill 的回合）MUST NOT 修改 `.go` / `.py` / `.ts` / `go.mod` / `package.json` / `pyproject.toml` / `Makefile` / `.env*` 等代码、依赖、构建、配置文件。lint skill 自身规则文件（`.agents/skills/*/SKILL.md`）也 MUST NOT 被 lint pass 自身修改。

#### Scenario: 仅文档被修改
- **WHEN** `git diff --name-only HEAD~N..HEAD` 查看某次 lint pass 的提交
- **THEN** 所有变更文件后缀为 `.md`（必要时含 `openspec/changes/.../*.md` 元数据），不出现代码 / 依赖 / 配置文件
