## ADDED Requirements

### Requirement: Diff 输出排除 provisioned skill 文件
`get_diff()` 端点 SHALL 根据 config.yaml manifest 中的 skill 名称列表，在生成 diff 时排除当前 workspace agent 的 `{config_dir}/skills/{skill_name}` 路径下的所有文件。

#### Scenario: Codex agent 的 skill 文件不出现在 diff 中
- **WHEN** workspace 的 agent_type 为 `codex`，config.yaml manifest 包含 `render` 和 `taskctl`
- **THEN** `get_diff()` 返回的 diff 中 SHALL 不包含 `.codex/skills/render/` 和 `.codex/skills/taskctl/` 下的文件

#### Scenario: Claude agent 的 skill 文件不出现在 diff 中
- **WHEN** workspace 的 agent_type 为 `claude-code`，config.yaml manifest 包含 `render` 和 `taskctl`
- **THEN** `get_diff()` 返回的 diff 中 SHALL 不包含 `.claude/skills/render/` 和 `.claude/skills/taskctl/` 下的文件

#### Scenario: 项目自带的 .claude/ 文件正常显示
- **WHEN** 项目仓库中已有 `.claude/CLAUDE.md` 且被修改，workspace 的 agent_type 为 `codex`
- **THEN** `get_diff()` 返回的 diff 中 SHALL 仍包含 `.claude/CLAUDE.md` 的变更

### Requirement: Skill 名称从 config.yaml manifest 动态读取
`get_diff()` SHALL 从 config.yaml 的 `skills.manifest` 键读取 skill 名称列表，与 SkillProvisioner 使用同一数据源。

#### Scenario: 新增 skill 后自动被过滤
- **WHEN** config.yaml manifest 中新增了 skill `foobar`，workspace 的 agent_type 为 `codex`
- **THEN** `get_diff()` 返回的 diff 中 SHALL 不包含 `.codex/skills/foobar/` 下的文件

### Requirement: Tracked 和 untracked 文件均被过滤
`get_diff()` SHALL 对 tracked changes（`git diff HEAD`）和 untracked files（`git ls-files --others`）两条路径均应用 skill 文件过滤。

#### Scenario: 已被 git add 的 skill 文件不出现在 diff 中
- **WHEN** `.codex/skills/render/render` 已被 `git add` 跟踪，且有内容变更
- **THEN** `get_diff()` 返回的 diff 中 SHALL 不包含该文件的变更

#### Scenario: Untracked 的 skill 文件不出现在 diff 中
- **WHEN** `.codex/skills/taskctl/taskctl` 为 untracked 文件
- **THEN** `get_diff()` 返回的 diff 中 SHALL 不包含该文件
