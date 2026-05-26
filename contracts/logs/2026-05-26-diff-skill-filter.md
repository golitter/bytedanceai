# 变更记录：get_diff 排除 provisioned skill 文件

## 变更原因

Agent 调用 render diff 时，SkillProvisioner 写入的 skill 二进制文件（render 82K 行、taskctl 82K 行）会出现在 diff 展示中，严重影响可读性。

## 变更文件

无 schema 变更。仅修改 `agentend/src/api/v1/workspace.py` 的 `get_diff()` 内部实现。

## 对比结果

- **无契约变更**：API 接口、请求/响应类型均未改变
- `get_diff()` 返回值格式不变（PlainTextResponse），仅过滤了 skill 文件内容

## 跨端影响

- **Frontend**: 无影响，diff 渲染逻辑不变
- **Backend**: 无影响
- **AgentEnd**: `get_diff()` 输出不再包含 `{config_dir}/skills/{skill_name}/` 下的文件

## 实现细节

- 新增 `_get_skill_exclusion_prefixes()` 读取 config.yaml manifest 构建 exclusion 路径
- tracked changes: git pathspec `:!.codex/skills/render` 排除
- untracked files: Python 端前缀匹配过滤
