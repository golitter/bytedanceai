## 1. Phase 0：上下文加载与基线确认

- [x] 1.1 读取全部 6 个 AGENTS.md（root / frontend / backend / agentend / docs / contracts）+ 3 个 `docs/reference/details.md`，建立文档基线
- [x] 1.2 执行 `wc -l` 统计每个 AGENTS.md / details.md 行数，标记已超线或接近上限的文件
- [x] 1.3 读取 [.agents/skills/agentsmd-linter/SKILL.md](.agents/skills/agentsmd-linter/SKILL.md) 与 [.agents/skills/doc-linter/SKILL.md](.agents/skills/doc-linter/SKILL.md)，确认本次 lint pass 完全沿用其 workflow
- [x] 1.4 执行 `codegraph_status` 确认索引健康；如不可用，回退到 `find` / `Grep` / `Read`

## 2. Phase A：执行 agentsmd-linter（先 AGENTS.md，后 details.md）

- [x] 2.1 验证根 [AGENTS.md](AGENTS.md) 目录树（用 `find` 比对实际 `frontend/` `backend/` `agentend/` `contracts/` `docs/` `docker/` `scripts/` `logs/` 等条目），修复悬空条目 / 补充新增
- [x] 2.2 验证根 AGENTS.md 中列出的 `make` 命令与 [Makefile](Makefile) 实际 target 一一对应
- [x] 2.3 验证 [frontend/AGENTS.md](frontend/AGENTS.md) 目录树、命令、引用路径，按 skill 规则修正
- [x] 2.4 验证 [backend/AGENTS.md](backend/AGENTS.md) 目录树、命令、引用路径，按 skill 规则修正（注意行数已 68，接近 70 上限）
- [x] 2.5 验证 [agentend/AGENTS.md](agentend/AGENTS.md) 目录树、命令、引用路径，按 skill 规则修正
- [x] 2.6 验证 [docs/AGENTS.md](docs/AGENTS.md) 文档分类约定 + 索引完整性（注意行数已 68，接近上限）
- [x] 2.7 验证 [contracts/AGENTS.md](contracts/AGENTS.md) 契约层描述 + 代码生成流程引用
- [x] 2.8 验证 [frontend/docs/reference/details.md](frontend/docs/reference/details.md) 索引覆盖率（扫描 `frontend/docs/**/*.md` 与索引条目比对）
- [x] 2.9 验证 [backend/docs/reference/details.md](backend/docs/reference/details.md) 索引覆盖率
- [x] 2.10 验证 [agentend/docs/reference/details.md](agentend/docs/reference/details.md) 索引覆盖率
- [x] 2.11 对任一 AGENTS.md 中内联展开的文档索引，迁移到对应 details.md，AGENTS.md 只保留一行引用
- [x] 2.12 对超出 70 / 150 行硬上限的 AGENTS.md / details.md，按 Phase 4 四手段压缩
- [ ] 2.13 提交批次 1：`docs(agents): agentsmd-linter 全项目 AGENTS.md + details.md 同步精炼`

## 3. Phase B：执行 doc-linter — frontend docs/

- [x] 3.1 扫描 [frontend/docs/](frontend/docs/) 下所有 `.md`（跳过 `tmp/`、`internal/`）
- [x] 3.2 验证 `frontend/docs/reference/tech-stack.md` 与 [frontend/package.json](frontend/package.json) 实际依赖版本一致
- [x] 3.3 验证 `frontend/docs/design/` 下文件是否使用 OpenSpec 格式（`## 实现了什么` + `## 怎么实现的`），未符合的重写
- [x] 3.4 识别 frontend 已实现但缺 design 文档的模块，按 `NN-<module>.md` 编号补齐
- [x] 3.5 验证 `frontend/docs/guides/` 操作指南中的命令与脚本可执行
- [x] 3.6 修复 frontend 文档中过时的代码引用、路径、目录结构、API 端点描述

## 4. Phase B：执行 doc-linter — backend docs/

- [x] 4.1 扫描 [backend/docs/](backend/docs/) 下所有 `.md`
- [x] 4.2 验证 `backend/docs/reference/tech-stack.md` 与 [backend/go.mod](backend/go.mod) 实际依赖版本一致
- [x] 4.3 验证 `backend/docs/design/` 下文件 OpenSpec 格式合规性，未符合的重写
- [x] 4.4 识别 backend 已实现但缺 design 文档的模块，按编号补齐
- [x] 4.5 验证 backend 文档中的 API 端点与 [backend/](backend/) 路由注册一致（用 `codegraph_search` 优先定位路由注册）
- [x] 4.6 修复 backend 文档中过时的代码引用、路径、目录结构、数据模型描述

## 5. Phase B：执行 doc-linter — agentend docs/

- [x] 5.1 扫描 [agentend/docs/](agentend/docs/) 下所有 `.md`
- [x] 5.2 验证 agentend 文档中 Python 依赖版本与 [agentend/pyproject.toml](agentend/pyproject.toml) 一致
- [x] 5.3 验证 `agentend/docs/design/` 下文件 OpenSpec 格式合规性，未符合的重写
- [x] 5.4 识别 agentend 已实现但缺 design 文档的模块（如 adapter / session / rule-engine），按编号补齐
- [x] 5.5 验证 agentend 文档中的 API 端点与 FastAPI / uvicorn 路由定义一致
- [x] 5.6 修复 agentend 文档中过时的代码引用、路径、目录结构、适配器描述

- [ ] 5.7 提交批次 2：`docs(frontend,backend,agentend): doc-linter 三端 docs 同步精炼`

## 6. Phase B：执行 doc-linter — 根 docs/

- [x] 6.1 扫描 [docs/](docs/) 下所有 `.md`（跳过 `tmp/`、`internal/`、`common/dev-plan/` 历史归档按既有规则不动）
- [x] 6.2 验证 [docs/guides/](docs/guides/) 中提到的命令、脚本路径有效
- [x] 6.3 验证 [docs/design/](docs/design/) 跨端架构文档与三端 design/ 不矛盾
- [x] 6.4 验证 [docs/reference/](docs/reference/) 中跨端引用与三端实际一致
- [x] 6.5 修复根 docs 中跨文档矛盾、过时描述、错误引用
- [ ] 6.6 提交批次 3：`docs(common): doc-linter 根 docs 同步精炼`

## 7. Phase C：终验

- [x] 7.1 执行 `wc -l` 全量检查所有 AGENTS.md ≤ 70 行、所有 details.md ≤ 150 行
- [x] 7.2 执行 `git diff --name-only HEAD~3..HEAD`，确认仅 `.md` 文件被修改，无代码 / 依赖 / 配置文件被改
- [x] 7.3 抽样 5 个被修改的文件，重读验证语义等价（压缩未损及定位信息）
- [x] 7.4 用 `codegraph_status` 复核索引健康；如发现本次修改涉及符号级影响，记录到本 change 的 design.md Open Questions 中（或留待下次 lint pass）
- [x] 7.5 在本 change 的归档阶段（`/opsx:archive`）前，确保所有 tasks 勾选完成
