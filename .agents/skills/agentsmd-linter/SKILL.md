---
name: agentsmd-linter
description: 全项目 AGENTS.md 同步精炼。扫描全部 5 个 AGENTS.md + 三端 details.md，交叉验证目录结构、文件引用、Makefile 命令后直接修改。
---

# Role

你是项目 AGENTS.md 文件的同步维护者。你的工作是扫描全部 5 个 AGENTS.md 文件及三端 `docs/reference/details.md`，发现与实际代码不一致的地方后**直接修改**，保持文档与代码库同步。

# Target Files（允许修改）

```
AGENTS.md                    # 根目录
frontend/AGENTS.md
backend/AGENTS.md
agentend/AGENTS.md
docs/AGENTS.md
frontend/docs/reference/details.md   # 前端详细文档索引
backend/docs/reference/details.md    # 后端详细文档索引
agentend/docs/reference/details.md   # Agent 端详细文档索引
```

禁止修改：`AGENTS.md`、`*.go`/`*.py`/`*.ts` 等代码文件、`package.json`/`go.mod`/`pyproject.toml` 等依赖文件。

# 文档分工

- **AGENTS.md** — 精简入口：项目简介 + 目录结构 + 命令 + 一行 `详见 [docs/reference/details.md]`
- **docs/reference/details.md** — 详细索引：design / reference / guides / testing 的完整文档列表 + API 端点 + 核心架构等技术参考内容

AGENTS.md 不内联展开文档索引，统一指向 details.md。

# Workflow

## Phase 1: 全量读取

读取全部 5 个 AGENTS.md + 3 个 details.md，建立当前文档状态的全景认知。

## Phase 2: 交叉验证

逐个文件进行以下检查：

### 2.1 目录结构验证

对每个 AGENTS.md 中列出的目录树：

1. 用 `ls` 或 `find` 验证每个列出的目录和文件是否实际存在
2. 检查是否有重要的新目录/文件未被记录（重点关注 `src/` 下的子目录）
3. 验证注释描述是否准确（如 `# （预留）` 标记的目录是否仍为空）

### 2.2 Makefile 命令验证

读取根 `Makefile`，验证每个 AGENTS.md 中列出的 `make` 命令是否实际存在，以及是否有新命令未被记录。

### 2.3 文件引用验证

验证 AGENTS.md 中引用的所有路径是否存在。

### 2.4 details.md 索引验证

对三端 `docs/reference/details.md`：

1. 验证每个索引条目的文件是否实际存在
2. 扫描实际 `docs/` 目录，发现未被索引的新文档，补充到 details.md
3. 验证 AGENTS.md 是否正确引用了 `docs/reference/details.md`（一行链接，不内联展开）

### 2.5 docs/AGENTS.md 索引验证

对 docs/AGENTS.md 中的文档索引：

1. 验证每个索引条目的文件是否实际存在
2. 扫描实际 `docs/` 目录，发现未被索引的新文档
3. 验证分类是否正确

## Phase 3: 直接修改

**发现问题的同时直接修复，不输出报告等确认。** 具体行为：

- 目录结构过时 → 按实际文件系统更新目录树
- 新增目录/文件未记录 → 补充到目录树，注释准确描述用途
- Makefile 命令缺失/多余 → 与 Makefile 实际内容对齐
- 文件引用 404 → 更新为正确路径或移除
- details.md 缺条目 → 补充新文档到对应分类
- docs/AGENTS.md 缺条目 → 补充新文档到对应分类
- AGENTS.md 内联了文档索引 → 移至 details.md，AGENTS.md 只保留一行引用
- 注释描述不准 → 修正为当前状态

每个文件修改完后，简要输出一行：`[已修复] AGENTS.md — 修复内容摘要` 或 `[已修复] docs/reference/details.md — 修复内容摘要`。

## Phase 4: 精简到行数限制

对每个修改过的文件检查行数：

- **AGENTS.md** ≤ 70 行
- **details.md** ≤ 150 行

超出时：
1. 压缩目录树注释 — 只保留关键定位信息
2. 合并短段落 — 不重复描述技术栈
3. 精简命令块 — 只列核心命令
4. 删减冗余链接 — 保留必要交叉引用

# Strict Constraints

1. **只修改 AGENTS.md 和 details.md**：不碰代码、依赖、配置文件。
2. **增量修正**：在原文基础上修正和补充，不盲目清空重写。
3. **保持统一风格**：目录树用代码块 + 注释格式，命令用 `bash` 代码块。
4. **AGENTS.md ≤ 70 行**：文档索引放 details.md，AGENTS.md 只一行引用。
5. **details.md 存放详细索引**：三端各自的 design / reference / guides / testing 文档列表 + 技术参考内容。
