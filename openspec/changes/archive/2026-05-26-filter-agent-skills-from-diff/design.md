## Context

当前 `get_diff()` 端点执行 `git diff HEAD` 和 `git ls-files --others --exclude-standard`，对输出不做任何路径过滤。SkillProvisioner 在 workspace 创建时将 render、taskctl 等大文件写入 `{config_dir}/skills/`，这些运行时产物会完整出现在 diff 展示中（render 82K 行、taskctl 82K 行），严重影响用户体验。

现有 `config-dir-exclusion` spec 通过 `.git/info/exclude` 防止 agent 文件被 `git add -A` 跟踪，但这只在 git 层面生效，不覆盖 diff 输出的过滤。

## Goals / Non-Goals

**Goals:**
- `get_diff()` 输出中不包含 provisioner 写入的 skill 文件
- 只排除当前 workspace agent 类型的 `{config_dir}/skills/{skill_name}`，不影响项目自带的 `.claude/` 文件
- skill 列表从 config.yaml manifest 动态读取，新增 skill 自动生效

**Non-Goals:**
- 不修改 `.gitignore` 或 `.git/info/exclude` 策略
- 不在前端做过滤（后端是单一防线）
- 不修改 SkillProvisioner 的目录结构

## Decisions

### 1. 使用 git pathspec 排除 tracked changes

`git diff HEAD -- . ':!.codex/skills/render' ':!.codex/skills/taskctl'`

pathspec 在 git 命令层面过滤，不产生后处理开销，且语义明确。

替代方案：
- 后处理解析 diff 输出过滤文件块 → 需要解析 unified diff 格式，复杂且脆弱
- `.gitignore` 排除 → 影响范围大，可能与项目已有 `.gitignore` 冲突

### 2. 对 untracked 文件做 Python 端路径过滤

`git ls-files --others` 不支持 pathspec 排除，因此在 Python 中过滤返回结果：

```python
skill_prefixes = [f"{config_dir}/skills/{name}/" for name in manifest]
filtered = [p for p in untracked.splitlines() if not any(p.startswith(pre) for pre in skill_prefixes)]
```

### 3. Manifest 作为 skill 名称单一数据源

读取 config.yaml 的 `skills.manifest` 键获取 skill 名称列表，与 SkillProvisioner 使用同一数据源，新增 skill 无需额外配置。

## Risks / Trade-offs

- [pathspec 只对 tracked diff 生效] → untracked 走 Python 过滤，双路径覆盖
- [config.yaml 读取需要文件系统访问] → 使用 agentend 已有的 yaml 配置加载机制
- [未来新增非 skills/ 目录的运行时产物] → 只需在过滤逻辑中添加路径前缀
