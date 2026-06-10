## 1. 数据模型 & 迁移

- [x] 1.1 创建 `skill_hub` 表（id, name, builtin, storage_path, description, file_count, total_size, uploaded_by, created_at, updated_at）
- [x] 1.2 创建 `agent_skill` 关联表（id, session_id, skill_name, agent_type, imported_at, UNIQUE KEY uk_session_skill）
- [x] 1.3 创建 Go model 文件 `internal/model/skill.go`（SkillHub + AgentSkill struct + GORM tags）

## 2. Backend — SkillsHub 存储 & 校验

- [x] 2.1 实现 `internal/service/skill_validator.go`：解压到临时目录、path traversal/symlink/绝对路径检查、zip bomb（10MB/100文件）检查
- [x] 2.2 实现校验逻辑：检查 SKILL.md 存在、解析 YAML frontmatter、提取 name/description、校验 name 不与 builtin 冲突
- [x] 2.3 实现 confirm 逻辑：先写 DB（skill_hub 表）再移文件到 `data/skills/hub/<name>/`
- [x] 2.4 实现 delete 逻辑：仅允许删除 `builtin=false` 的记录，删除文件目录 + DB 记录
- [x] 2.5 实现 list 逻辑：`GET /api/skills` 返回所有记录含 `import_count`（从 agent_skill COUNT）
- [x] 2.6 创建 `data/skills/hub/` 目录结构

## 3. Backend — API Handler & 路由

- [x] 3.1 创建 `internal/handler/skill.go`：Upload、Confirm、List、Delete handler
- [x] 3.2 实现 `POST /api/internal/builtin-skills` handler：接收 Agentend 上报，UPSERT 到 skill_hub 表（builtin=true）
- [x] 3.3 实现 `POST /api/skills/:name/import` handler：校验 agent_type、查找 session worktree、物理复制、写 agent_skill 记录
- [x] 3.4 实现 `DELETE /api/skills/:name/sessions/:sessionId` handler：删除 worktree 中文件 + agent_skill 记录
- [x] 3.5 注册所有 skills 相关路由

## 4. Backend — Agentend Client 改造

- [x] 4.1 创建 `pkg/agentend_client/skill_client.go`：调用 `GET /api/v1/skills/:agent_type?workspace_path=xxx`
- [x] 4.2 改造 `internal/handler/agent_profile.go` 的 `GetDetail`：从 mock 硬编码改为代理到 Agentend 读取 skills
- [x] 4.3 Agentend 不可达时 skills 返回空数组，不阻塞页面

## 5. Agentend — Skills 扫描 & 上报

- [x] 5.1 创建 `api/v1/skills.py`：`GET /api/v1/skills/:agent_type` 路由，扫描 workspace skills 目录，解析 SKILL.md frontmatter
- [x] 5.2 实现 builtin 判定：对比 `config.yaml` 的 `settings.skills.manifest` keys
- [x] 5.3 实现启动钩子：读取 builtin_dir，解析每个 SKILL.md，调用 `POST /api/internal/builtin-skills` 上报
- [x] 5.4 实现重试机制：指数退避 3 次（2s/4s/8s），失败记录日志不阻塞
- [x] 5.5 在 `app/main.py` 注册新路由 + 启动钩子

## 6. Frontend — 侧边栏 & 导航

- [x] 6.1 修改 `stores/navigation-store.ts`：`NavTab` 类型增加 `'skills'`
- [x] 6.2 修改 `components/layout/IconSidebar.tsx`：新增「技能」NavItem（星星图标），位于通讯录和管理之间
- [x] 6.3 修改 `pages/ImPage.tsx`：路由新增 skills tab，渲染 SkillsHubPage

## 7. Frontend — SkillsHub 管理页面

- [x] 7.1 创建 `pages/SkillsHubPage.tsx`：页面布局（标题 + 搜索 + 上传按钮 + skill 列表）
- [x] 7.2 实现 skill 列表渲染：分区显示 builtin/external、SkillCard 组件（名称、描述、标签、导入数、删除按钮）
- [x] 7.3 实现搜索过滤：输入时实时过滤 skill cards
- [x] 7.4 创建 `components/SkillUploadDialog.tsx`：上传对话框（虚线上传区、校验结果面板、名称确认、确认/取消）
- [x] 7.5 实现上传交互：选择 zip → 显示校验结果 → 命名确认 → 动画添加到列表
- [x] 7.6 实现删除交互：确认对话框 → DELETE 请求 → 动画移除 card

## 8. Frontend — Agent 详情页 Skills 改造

- [x] 8.1 创建 `components/SkillImportDialog.tsx`：导入选择对话框（列出 hub 中 external skills、勾选、已导入项灰色不可选）
- [x] 8.2 改造 `pages/AgentProfilePage.tsx` Skills 区域：增加 external 标签、导入/移除按钮
- [x] 8.3 实现导入交互：点击「导入外部 Skill」→ 弹出对话框 → 勾选确认 → 动画添加 skill cards
- [x] 8.4 实现移除交互：点击「移除」→ DELETE 请求 → 动画移除 card
- [x] 8.5 Orchestrator 类型 Agent：隐藏导入按钮区域

## 9. Frontend — SkillCard 样式 & API 层

- [x] 9.1 改造 `components/chat/SkillCard.tsx`：增加 external 标签样式（indigo 配色）
- [x] 9.2 修改 `lib/api.ts`：新增 skills 相关 API 函数（uploadSkill, confirmSkill, listSkills, deleteSkill, importSkill, removeSkill）

## 10. 契约层 & 文档

- [ ] 10.1 更新 `contracts/schemas/` 中相关 YAML（如有 skills 相关 schema 变更）
- [ ] 10.2 运行 `make generate` 生成三端类型
- [ ] 10.3 在 `contracts/logs/` 写入变更记录
