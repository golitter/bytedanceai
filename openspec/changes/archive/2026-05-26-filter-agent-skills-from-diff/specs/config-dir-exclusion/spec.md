## ADDED Requirements

### Requirement: Diff 生成层排除 skill 文件
在 `config-dir-exclusion` 已有的 git exclude 策略基础上，diff 生成层 SHALL 额外过滤 provisioner 写入的 skill 文件，确保即使文件未被 git exclude 覆盖也不会出现在 diff 输出中。

#### Scenario: git exclude 未覆盖的 skill 文件仍被 diff 过滤
- **WHEN** workspace 的 `.git/info/exclude` 未包含 agent config 目录（如 exclude 写入失败）
- **AND** `.codex/skills/render/render` 为 untracked 文件
- **THEN** `get_diff()` 返回的 diff 中 SHALL 仍不包含该文件
