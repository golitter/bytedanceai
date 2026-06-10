## 1. Skill 名称读取

- [x] 1.1 在 workspace.py 中添加函数，从 config.yaml 读取 `skills.manifest` 键返回 skill 名称列表
- [x] 1.2 复用 `get_agent_config_dir()` 获取当前 workspace 的 config_dir

## 2. Diff 生成过滤

- [x] 2.1 修改 `get_diff()` 的 `git diff HEAD` 调用，拼接 pathspec 排除 `{config_dir}/skills/{skill_name}` 路径
- [x] 2.2 修改 `get_diff()` 的 untracked 文件处理逻辑，过滤掉匹配 `{config_dir}/skills/{skill_name}/` 前缀的路径

## 3. 验证

- [x] 3.1 启动 agentend，创建 codex workspace，调用 render diff，确认 skill 文件不出现在 diff 输出中
- [x] 3.2 确认项目自带的 `.claude/` 文件变更仍正常显示
