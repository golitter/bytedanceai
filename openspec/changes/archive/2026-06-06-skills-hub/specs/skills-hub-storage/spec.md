## ADDED Requirements

### Requirement: Upload and validate skill zip
Backend SHALL 接收前端上传的 .zip 文件，解压到临时目录，执行以下校验：
1. 拒绝路径含 `..` 的 entry（path traversal）
2. 拒绝符号链接指向包外
3. 拒绝绝对路径
4. 解压后总大小超过 10MB SHALL 拒绝（zip bomb 防护）
5. 解压文件数超过 100 SHALL 拒绝（资源耗尽防护）
6. 根目录必须包含 `SKILL.md`
7. `SKILL.md` 必须有 YAML frontmatter（`---` 包裹）
8. frontmatter 必须包含 `name` 字段
9. `name` 不可与已有 builtin skill 同名

校验通过后，SHALL 返回 `{valid: true, name, description, file_count, total_size}`。
校验失败 SHALL 返回 `{valid: false, errors: [...]}`。

#### Scenario: Valid zip upload
- **WHEN** 前端上传一个包含 `SKILL.md`（frontmatter 含 `name: my-skill`）的合法 .zip 文件
- **THEN** Backend 解压成功，返回 `{valid: true, name: "my-skill", description: "...", file_count: 4, total_size: 128000}`

#### Scenario: Missing SKILL.md
- **WHEN** 上传的 .zip 根目录没有 `SKILL.md`
- **THEN** Backend 返回 `{valid: false, errors: ["missing SKILL.md"]}`

#### Scenario: Missing frontmatter name
- **WHEN** `SKILL.md` 的 frontmatter 没有 `name` 字段
- **THEN** Backend 返回 `{valid: false, errors: ["missing name field"]}`

#### Scenario: Zip bomb protection
- **WHEN** 解压后总大小超过 10MB
- **THEN** Backend 返回 `{valid: false, errors: ["zip bomb: total size exceeds 10MB"]}`

#### Scenario: Path traversal rejection
- **WHEN** zip 中有 entry 路径含 `..`
- **THEN** Backend 返回 `{valid: false, errors: ["path traversal detected"]}`

#### Scenario: Builtin name conflict
- **WHEN** 上传的 skill `name` 与已有 builtin skill 同名（如 `render`）
- **THEN** Backend 返回 `{valid: false, errors: ["name conflicts with builtin skill"]}`

### Requirement: Confirm and store skill
校验通过后，前端发起确认请求。Backend SHALL 先写入 `skill_hub` 表（`builtin=false`），成功后将文件从临时目录移到 `data/skills/hub/<name>/`。若 DB 写入失败，临时文件由定时清理 job 回收。

#### Scenario: Successful confirm
- **WHEN** 前端发送 `POST /api/skills/confirm` 且 `skill_hub` 表无同名记录
- **THEN** Backend 先插入 `skill_hub` 表记录，成功后移动文件到 `data/skills/hub/<name>/`，返回成功

#### Scenario: Duplicate name on confirm
- **WHEN** 前端确认的 `name` 已存在于 `skill_hub` 表
- **THEN** Backend 返回 409 Conflict

#### Scenario: DB write failure
- **WHEN** `skill_hub` 表写入失败
- **THEN** 临时目录文件保留，不移动。Backend 返回 500 错误

### Requirement: Delete external skill from hub
Backend SHALL 支持 `DELETE /api/skills/:name` 删除 `data/skills/hub/<name>/` 目录及其文件，同时删除 `skill_hub` 表记录。仅允许删除 `builtin=false` 的记录。

#### Scenario: Delete external skill
- **WHEN** 前端发送 `DELETE /api/skills/my-skill`
- **THEN** Backend 删除 `data/skills/hub/my-skill/` 目录和 `skill_hub` 表记录

#### Scenario: Cannot delete builtin
- **WHEN** 前端发送 `DELETE /api/skills/render`（builtin skill）
- **THEN** Backend 返回 403 Forbidden

#### Scenario: Skill not found
- **WHEN** 前端删除不存在的 skill
- **THEN** Backend 返回 404

### Requirement: List all skills in hub
Backend SHALL 提供 `GET /api/skills` 返回所有 `skill_hub` 记录，按 `builtin DESC, name ASC` 排序。每条记录包含 `name`、`description`、`builtin`、`file_count`、`total_size`、`import_count`（从 `agent_skill` 表 COUNT 得出）。

#### Scenario: List skills with import count
- **WHEN** 前端调用 `GET /api/skills`
- **THEN** 返回所有 skill 记录，每条含 `import_count` 字段

### Requirement: Receive builtin skills report
Backend SHALL 提供 `POST /api/internal/builtin-skills` 接收 Agentend 上报的 builtin skills 列表。使用 UPSERT 逻辑：按 `name` 匹配，存在则更新 `description`，不存在则插入，`builtin` 强制为 `true`。

#### Scenario: First report
- **WHEN** Agentend 启动时上报 `[{name: "render", description: "..."}]`
- **THEN** Backend 插入 `skill_hub` 记录，`builtin=true`

#### Scenario: Re-report (idempotent)
- **WHEN** Agentend 重启再次上报相同列表
- **THEN** Backend 更新已有记录的 `description`，不产生重复行
