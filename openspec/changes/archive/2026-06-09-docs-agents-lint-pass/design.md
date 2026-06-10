## Context

项目文档分布在 6 个 AGENTS.md（root / frontend / backend / agentend / docs / contracts）+ 3 个 `docs/reference/details.md` + 根 `docs/` 与三端 `docs/` 下数十个 `.md` 文件中。文档与代码漂移的典型来源：

- 新增能力（如 `agent-memory`、`skill-import-export`、`theme-switching`、`git-graph-panel`、`contacts-page`）后未同步更新 AGENTS.md 目录树与 details.md 索引
- 重构后文件路径变化，文档引用失效
- 依赖版本升级但 `docs/reference/tech-stack.md` 仍写旧版本
- `docs/design/` 下部分历史文档未按 OpenSpec 格式（`## 实现了什么` + `## 怎么实现的`）书写
- AGENTS.md 内联了文档索引而非一行引用 details.md

两个 lint skill 已在仓库内：
- [.agents/skills/agentsmd-linter/SKILL.md](.agents/skills/agentsmd-linter/SKILL.md) — 6 AGENTS.md + 3 details.md 同步精炼
- [.agents/skills/doc-linter/SKILL.md](.agents/skills/doc-linter/SKILL.md) — 根 docs/ + 三端 docs/ 全量扫描

两个 skill 都采用 Phase 1 全量读取 → Phase 2 交叉验证 → Phase 3 直接修改 → Phase 4 行数硬约束的工作流。

## Goals / Non-Goals

**Goals:**

- 让两个 lint skill 的 workflow 在仓库上完整跑一遍，并把发现的差异修复掉
- 让所有 AGENTS.md 与 details.md 满足行数硬上限（AGENTS.md ≤ 70 行 / details.md ≤ 150 行）
- 让三端 `docs/design/` 缺失的实施文档按 OpenSpec 格式补齐
- 让根 `docs/` 与三端 `docs/` 中的代码引用（路径、模块名、API 端点、技术栈、目录结构）与当前代码一致
- 在 OpenSpec 中留下一份可追溯的 change 记录（即本 change 本身），便于将来按相同方式触发新一轮 lint

**Non-Goals:**

- 不修改 `.go` / `.py` / `.ts` / 依赖清单 / 配置文件 / Makefile 等代码或构建文件
- 不重写 `docs/internal/`、`docs/tmp/`（两个 skill 均明确跳过）
- 不调整文档 5 分类法（`design/` / `reference/` / `guides/` / `testing/` / `backlog/`）本身
- 不写新功能或修改运行时行为
- 不重新设计 lint skill 自身的规则（skill 文件视为不可变输入）
- 不删除历史归档 change / 历史 design 文档

## Decisions

### Decision 1: 严格按 skill 的工作流顺序执行

**选择**：先 `agentsmd-linter`（Phase 1–4），再 `doc-linter`（Phase 1–4）。

**理由**：`doc-linter` 在 Context Loading 阶段要求"先读全部 AGENTS.md 建立全景认知"。先跑 `agentsmd-linter` 可让 AGENTS.md 处于最新状态，避免 `doc-linter` 拿到的"地图"已经过时。

**替代方案**：并行 / 反向顺序 — 都会让 `doc-linter` 在过时的 AGENTS.md 基础上工作，二次返工。

### Decision 2: 三端 docs 顺序固定为 frontend → backend → agentend → 根 docs

**选择**：与 `doc-linter` Phase 2 中明确规定的顺序一致。

**理由**：skill 已经硬编码了"先三端后根"的依赖关系（根 docs 涉及跨端联调，依赖三端文档的准确性）。复用既定顺序，避免流程外偏差。

### Decision 3: 行数硬上限视为强制约束，超出立即压缩

**选择**：AGENTS.md > 70 行 / details.md > 150 行时按 skill Phase 4 列出的 4 个手段（压缩目录树注释 / 合并短段落 / 精简命令块 / 删减冗余链接）压缩到上限内。

**理由**：当前根 AGENTS.md 60 行、frontend AGENTS.md 53 行、backend AGENTS.md 68 行（接近上限）、agentend AGENTS.md 62 行、docs AGENTS.md 68 行（接近上限）、contracts AGENTS.md 51 行 — 多个文件已经接近上限，本轮修复如果新增内容可能直接超线。

**替代方案**：放宽行数上限 — 会破坏 skill 的硬约束，且 skill 文件不在修改范围。

### Decision 4: 三端 docs/design/ 缺失实施文档按 OpenSpec 格式补齐

**选择**：发现缺失时按 `## 实现了什么` + `## 怎么实现的` + 真实代码片段的格式创建，编号遵循每端既有最大编号 +1。

**理由**：`doc-linter` Phase 3 明确要求此格式，且三端 `docs/design/` 已有的多数文档遵循此格式。统一格式便于后续 lint pass 自动校验。

### Decision 5: 文档内的代码引用统一使用相对路径 markdown 链接

**选择**：`[filename.ts](src/filename.ts)` 或 `[file.ts:42](src/file.ts#L42)` 形式。

**理由**：VSCode native extension 上下文要求，且 IDE 内可点击跳转。`codegraph` 工具不依赖此格式，但人工审阅和 IDE 体验受益。

### Decision 6: 修复粒度 — 增量修正，不重写

**选择**：在原文基础上修正和补充，绝不盲目清空重写。

**理由**：两个 skill 在 Strict Constraints 中均明确此项。重写会丢失原有作者的细微表达（如对某模块边界的强调），且 git diff 噪音大。

## Risks / Trade-offs

- **[Risk] 大量文档修改导致 PR diff 难以审阅** → 按 skill 顺序分阶段提交：先 AGENTS.md + details.md 一组提交，再三端 docs 一组提交，再根 docs 一组提交。每组提交信息以 `docs(xxx):` 前缀遵循 [docs/guides/git-conventions.md](docs/guides/git-conventions.md)。
- **[Risk] 修复后又发现新不一致** → 在 [tasks.md](tasks.md) 中安排最后一轮 sweep：先 `git diff` 全量扫一遍本次修改，再用 `codegraph_status` 检查索引是否健康，再做一次 Phase 2 抽样验证。
- **[Risk] 行数压缩误伤关键信息** → 压缩时只动冗余表达（重复的技术栈描述、过时的注释、多余的空行），保留定位信息（目录用途、命令含义、关键交叉引用）。压缩后必须保证语义等价。
- **[Trade-off] 增量修正 vs. 风格统一**：增量优先意味着不会把所有文档强行写成同一种段落风格 — 风格多样性保留，但格式（标题 / 代码块 / 链接）仍按 skill 约束统一。
- **[Trade-off] 不动 skill 文件**：skill 文件本身的规则如果未来需要更新，是另一个 change 的事。本次只把 skill 当作"既定的 lint 规则"执行。
- **[Risk] `doc-linter` 要求读取代码（go.mod / package.json / pyproject.toml / 路由文件）** → 验证阶段只读不改；用 `codegraph_*` 工具优先（毫秒级响应），grep / Read 作为补充。
- **[Risk] 三端 docs/design/ 补齐可能规模爆炸** → 只补"已实现但缺文档"的模块；不写未实现的功能、未来计划或设计理论（与 skill 严格约束 #3 一致）。如果新增文档数量超过 10 个，分批处理并在 tasks.md 中标注。

## Migration Plan

不涉及代码迁移 / 数据迁移 / 部署步骤。

文档修改的 git 提交建议分三批：

1. `docs(agents): agentsmd-linter 全项目 AGENTS.md + details.md 同步精炼` — 6 AGENTS.md + 3 details.md 的修改
2. `docs(frontend,backend,agentend): doc-linter 三端 docs 同步精炼` — 三端 `docs/` 下的修改 + 新增 design/ 实施文档
3. `docs(common): doc-linter 根 docs 同步精炼` — 根 `docs/` 下的修改

回滚策略：`git revert` 三个 commit 即可。文档变更对运行时无影响，无回滚风险。

## Open Questions

- **Q1**：根 `docs/internal/` 目录在两个 skill 中都被列为"跳过"，但内部资料是否需要独立 lint pass？— 默认不需要，由用户在后续 change 中决定。
- **Q2**：本次 lint pass 完成后，是否要把"定期跑一次 lint pass"加入项目维护流程？— 不在本 change 范围；本 change 仅完成一次回合，记录归档后可作为下次的参考模板。
