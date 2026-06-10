## Why

当前系统仅支持 Agentend 内置 Skills（builtin），用户无法上传自定义 Skill 或将其导入到指定 Agent。随着 Agent 能力扩展，需要一个统一的 **SkillsHub** 管理界面，支持外部 Skill 的上传、校验、入库、导入到 Agent、独立移除等完整生命周期，同时保留 builtin Skills 的只读展示。

## What Changes

- 新增 **SkillsHub 管理页面**（侧边栏「技能」Tab），展示 builtin + external skill 列表
- 新增 **Skill 上传流程**：前端上传 .zip → Backend 解压校验（SKILL.md、frontmatter、path traversal、zip bomb）→ 用户确认命名 → 入库到 `data/skills/hub/`
- 新增 **Skill 导入/移除**：在 Agent 详情页导入外部 Skill 到指定 session 的 worktree，支持独立移除（仅 adapter 层 agent：claude-code / opencode / codex）
- 新增 **Builtin Skills 上报**：Agentend 启动时将 builtin skills 列表上报 Backend，写入数据库（upsert 幂等）
- 新增 **Backend API**：SkillsHub CRUD + Import/Remove 端点 + Builtin 上报端点
- 新增 **Agentend API**：`GET /api/v1/skills/:agent_type` 扫描运行时 skills 目录
- **改造 Agent 详情页 Skills 区域**：增加 external 标签、导入/移除按钮、导入选择对话框
- **改造侧边栏**：新增「技能」NavTab（通讯录下方）

## Capabilities

### New Capabilities
- `skills-hub-storage`: Backend 端 SkillsHub 文件存储管理（上传解压、格式校验、zip bomb 防护、确认入库、物理删除）
- `skills-hub-api`: SkillsHub RESTful API 端点（upload、confirm、list、delete、import、remove、builtin 上报）
- `skills-hub-page`: 前端 SkillsHub 管理页面（列表展示、搜索过滤、上传对话框、删除确认）
- `skill-import-export`: Agent 级别的 Skill 导入/移除功能（物理复制到 worktree、agent_skill 表关联、前端导入选择对话框）

### Modified Capabilities
- `agent-skill-model`: Skills 数据从 mock 硬编码改为从 Agentend 实时读取（区分 builtin/external、增加 source 字段）
- `agent-detail-page`: Skills 区域增加 external 标签样式、导入按钮、移除按钮，Orchestrator 隐藏导入区域

## Impact

- **Backend (Go)**：新增 handler/service/model 层，涉及文件存储、DB 新增 `skill_hub` + `agent_skill` 两张表、路由注册
- **Agentend (Python)**：新增 skills 扫描 API、启动钩子 builtin 上报逻辑
- **Frontend (React)**：新增 SkillsHubPage / SkillUploadDialog / SkillImportDialog 组件，改造 IconSidebar（NavTab 类型）、AgentProfilePage（Skills 区域）、SkillCard（external 标签）
- **侧边栏导航**：NavTab 增加 `'skills'`，位于通讯录和管理之间
- **API 兼容性**：现有 `/api/sessions/:sid/profile` 和 `/api/sessions/:sid/detail` 的 skills 字段格式扩展（增加 `builtin`、`source` 字段），向后兼容
