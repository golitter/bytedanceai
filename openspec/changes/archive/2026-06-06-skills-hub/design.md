## Context

当前系统中的 Skills 仅支持 Agentend 内置（builtin），通过 `config.yaml` 的 `settings.skills.manifest` 配置。Agent 详情页的 Skills 展示使用后端硬编码 mock 数据。用户无法上传自定义 Skill 或按需导入到 Agent。

设计文档 `docs/design/07-skills-hub-external-skills.md` 已完成评审，demo 可视化在 `docs/payloads/skills-hub-demo.html`。本设计基于该文档的核心决策实现。

三端架构：Backend (Go) 负责存储和 API、Agentend (Python) 负责 skills 目录扫描和 builtin 上报、Frontend (React) 负责 UI 交互。

## Goals / Non-Goals

**Goals:**
- 实现完整的 external Skill 生命周期：上传 → 校验 → 入库 → 导入到 Agent → 移除
- 实现 builtin Skills 的自动上报和持久化存储
- SkillsHub 管理页面，支持搜索、上传、删除
- Agent 详情页支持导入/移除 external Skill
- 安全防护：path traversal、symlink、zip bomb、名称冲突

**Non-Goals:**
- Skill 版本管理 / 更新机制（后续迭代）
- Skill 权限控制（当前所有用户共享同一技能库）
- Skill 依赖关系管理
- 导入到 Orchestrator（设计决策明确排除）

## Decisions

### D1: Hub 存储位置 — Backend 文件系统

**选择**: `data/skills/hub/<name>/` 目录
**替代方案**: 数据库 BLOB、对象存储（S3/MinIO）
**理由**: 项目当前无对象存储依赖，文件系统简单直接。Skill 本质是文件集合（SKILL.md + scripts），文件系统天然适配。

### D2: 导入方式 — 物理复制到 worktree

**选择**: 将 hub 中的 skill 文件物理复制到 `worktree/.claude/skills/<name>/`（或 `.opencode/skills/`）
**替代方案**: 符号链接（symlink）
**理由**: 每个 session 独立 worktree，物理复制保证隔离性。删除 hub 源文件不影响已导入副本。

### D3: Builtin 上报 — Agentend 启动时 HTTP POST + Upsert

**选择**: Agentend 启动时 `POST /api/internal/builtin-skills`，Backend 使用 `INSERT ... ON DUPLICATE KEY UPDATE`
**替代方案**: Backend 直接读 Agentend 配置文件、共享数据库
**理由**: Agentend 是独立进程，通过 HTTP 上报保持松耦合。Upsert 保证幂等，重启不产生重复。

### D4: 运行时 Skills 读取 — Backend 代理到 Agentend

**选择**: `GET /api/agents/:sid` 内部调用 Agentend `GET /api/v1/skills/:agent_type?workspace_path=xxx`
**替代方案**: Backend 直接读文件系统
**理由**: Agentend 拥有 workspace_path 映射和 skills 目录结构知识，Backend 不应关心文件布局。

### D5: API 设计 — DELETE 使用路径参数

**选择**: `DELETE /api/skills/:name/sessions/:sessionId`（移除导入）
**替代方案**: `DELETE /api/skills/:name/remove` with body
**理由**: HTTP DELETE 带 request body 在部分代理/客户端中行为不一致，路径参数更可靠。

### D6: 数据模型 — 统一 `skill_hub` 表 + `agent_skill` 关联表

**选择**: 单表 `skill_hub`（builtin + external）用 `builtin` 布尔字段区分，`agent_skill` 关联 session 和 skill
**替代方案**: 分表 `skill_builtin` + `skill_external`
**理由**: 统一表简化查询逻辑，`builtin` 字段足够区分行为差异（不可删除 vs 可删除）。

## Risks / Trade-offs

**[Zip Bomb]** → 限制解压后总大小 ≤ 10MB，文件数 ≤ 100，后端强制校验

**[名称冲突]** → 前端 + 后端双重校验：external skill name 不可与 builtin 同名

**[Confirm 回滚]** → 先写 DB 后移文件。DB 写入失败时，临时文件由定时清理 job 回收

**[Builtin 上报失败]** → 指数退避重试 3 次（2s/4s/8s），仍失败则记录日志不阻塞启动

**[并发导入]** → `agent_skill` 表 UNIQUE KEY `(session_id, skill_name)` 保证幂等，前端提前检查避免冲突提示

**[运行中导入]** → 不管 Agent 状态直接复制文件。Claude Code / OpenCode 等工具在下次对话时自动发现新 skill

**[导入数量统计]** → 初版使用 `COUNT` 查询 `agent_skill` 表；数据量大时可增加 `import_count` 冗余字段优化
