## ADDED Requirements

### Requirement: SkillsHub CRUD endpoints
Backend SHALL 提供以下 RESTful API 端点：

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/skills/upload` | 上传 zip，返回校验结果 |
| POST | `/api/skills/confirm` | 确认入库 + 用户命名 |
| GET | `/api/skills` | 列出 hub 中所有 skills |
| DELETE | `/api/skills/:name` | 从 hub 删除（仅 external） |

#### Scenario: Upload returns validation result
- **WHEN** `POST /api/skills/upload` with multipart/form-data zip file
- **THEN** 返回 `{valid, name?, description?, errors?[]}`

#### Scenario: Confirm creates hub record
- **WHEN** `POST /api/skills/confirm` with `{name: "my-skill"}`
- **THEN** 创建 `skill_hub` 记录 + 移动文件，返回 200

#### Scenario: List returns all skills
- **WHEN** `GET /api/skills`
- **THEN** 返回 `[{name, description, builtin, file_count, total_size, import_count}]`

#### Scenario: Delete removes external only
- **WHEN** `DELETE /api/skills/:name` where `builtin=false`
- **THEN** 删除文件和 DB 记录，返回 200

### Requirement: Agent-level skill import endpoint
Backend SHALL 提供 `POST /api/skills/:name/import`，将 hub 中的 skill 物理复制到指定 session 的 worktree skills 目录，并在 `agent_skill` 表插入关联记录。

#### Scenario: Import to adapter agent
- **WHEN** `POST /api/skills/my-skill/import` with `{session_id: "abc-123"}`
- **AND** 该 session 的 agent_type 为 `claude-code` / `opencode` / `codex`
- **THEN** 复制 `data/skills/hub/my-skill/` 到 `worktree/.claude/skills/my-skill/`，写入 `agent_skill` 记录

#### Scenario: Reject import for orchestrator
- **WHEN** session 的 agent_type 为 `orchestrator`
- **THEN** 返回 403 Forbidden

#### Scenario: Reject duplicate import
- **WHEN** `agent_skill` 表已存在 `(session_id, skill_name)` 记录
- **THEN** 返回 409 Conflict

#### Scenario: Skill not in hub
- **WHEN** 导入的 skill name 不在 `skill_hub` 表中
- **THEN** 返回 404

### Requirement: Agent-level skill remove endpoint
Backend SHALL 提供 `DELETE /api/skills/:name/sessions/:sessionId`，从指定 session 的 worktree 删除 skill 文件和 `agent_skill` 记录。

#### Scenario: Remove imported skill
- **WHEN** `DELETE /api/skills/my-skill/sessions/abc-123`
- **THEN** 删除 `worktree/.claude/skills/my-skill/` 目录和 `agent_skill` 记录

#### Scenario: Remove non-existent import
- **WHEN** `agent_skill` 表无对应记录
- **THEN** 返回 404

### Requirement: Agentend skills scan endpoint
Agentend SHALL 提供 `GET /api/v1/skills/:agent_type?workspace_path=xxx`，扫描指定 workspace 的 skills 目录，返回 `[{name, description, builtin, source}]`。

#### Scenario: Scan Claude Code workspace
- **WHEN** `GET /api/v1/skills/claude-code?workspace_path=/worktrees/task-1/sess-abc`
- **THEN** 扫描 `/worktrees/task-1/sess-abc/.claude/skills/` 目录，解析每个子目录的 `SKILL.md` frontmatter，返回 skill 列表

#### Scenario: Builtin detection
- **WHEN** skill name 在 `config.yaml` 的 `settings.skills.manifest` 中
- **THEN** 该 skill 的 `builtin` 为 `true`，`source` 为 `"builtin"`

#### Scenario: External detection
- **WHEN** skill name 不在 manifest 中
- **THEN** 该 skill 的 `builtin` 为 `false`，`source` 为 `"hub"`

#### Scenario: Empty skills directory
- **WHEN** workspace 的 skills 目录不存在或为空
- **THEN** 返回空数组 `[]`

### Requirement: Builtin skills report from Agentend
Agentend 启动时 SHALL 读取 `config.yaml` 的 `settings.skills.manifest` 和 `skills.builtin_dir`，解析每个 SKILL.md frontmatter，调用 `POST /api/internal/builtin-skills` 上报给 Backend。

#### Scenario: Successful report
- **WHEN** Agentend 启动且 Backend 可达
- **THEN** 发送 `POST /api/internal/builtin-skills` with `[{name, description, builtin: true, source: "builtin"}]`

#### Scenario: Backend unreachable retry
- **WHEN** 上报请求失败
- **THEN** 指数退避重试 3 次（间隔 2s/4s/8s），仍失败则记录错误日志，不阻塞启动
