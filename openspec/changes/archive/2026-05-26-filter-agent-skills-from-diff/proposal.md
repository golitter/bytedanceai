## Why

Agent 调用 render diff 时，SkillProvisioner 写入的 skill 二进制文件（render 82K 行、taskctl 82K 行）会出现在 diff 展示中。这些文件是运行时产物，不应展示给用户，且体积巨大严重影响 diff 可读性。

## What Changes

- 在 `get_diff()` 中根据 config.yaml manifest 的 skill 名称列表，动态拼接 git pathspec 排除当前 workspace agent 的 `{config_dir}/skills/{skill_name}` 路径
- 对 `git ls-files --others` 返回的 untracked 文件做同样的路径过滤
- 只改 `workspace.py` 的 `get_diff()` 函数，不修改 `.gitignore`、不影响项目自带的 `.claude/` 等文件

## Capabilities

### New Capabilities

- `diff-skill-filter`: 在 workspace diff 生成时，根据 config.yaml manifest 动态排除当前 agent 的 provisioned skill 文件

### Modified Capabilities

- `config-dir-exclusion`: 补充 diff 生成层的过滤能力（现有 spec 仅覆盖 git exclude，未覆盖 diff 输出过滤）

## Impact

- `agentend/src/api/v1/workspace.py` — `get_diff()` 函数，添加 pathspec 参数和 untracked 文件过滤
- `agentend/config.yaml` — manifest 作为 skill 名称的数据源（只读，不修改）
