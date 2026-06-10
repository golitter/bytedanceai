## Context

当前 agent 在 worktree 中工作时，Claude Code / OpenCode 运行时会在配置目录（`.claude`/`.opencode`）下生成 memory、settings、scheduled_tasks 等文件。现有 `SkillProvisioner._write_git_exclude()` 仅排除 `/.claude/skills/<name>` 这类 skill 子目录，不覆盖整个配置目录。`git add -A` 会将 agent 运行时文件一并提交。

同时，API 层 `MergeRequest.target_branch` 默认为 `"main"`，导致 agent 调 merge 时跳过 `task/{task_id}` 直接合入主线，不符合三级分支设计（agent→task→main）。

约束：
- `.git/info/exclude` 只影响 untracked 文件，已 tracked 的内容（如仓库自带的 `.claude/CLAUDE.md`）不受影响
- Worktree 的 `.git` 是一个文件（指向主仓库的 gitdir），需要解析才能找到真正的 git info 目录

## Goals / Non-Goals

**Goals:**
- agent 配置目录下的运行时文件（memory、settings 等）不会被 `git add -A` 提交
- agent merge 默认目标为 task 分支，task→main 为独立的平台操作
- exclude 逻辑归属 WorkspaceManager（workspace 级别关注点），不混在 SkillProvisioner（skill 分发关注点）中

**Non-Goals:**
- 不修改 agent 侧的行为指令（taskctl SKILL.md 等）
- 不改变现有的 branch 命名规则（`agent/{session_id}/{task_id}`、`task/{task_id}`）
- 不处理 merge conflict 的自动化解决策略
- 不排除 agent 对已 tracked 文件的修改（如仓库自带的 `.claude/CLAUDE.md` 的改动仍可提交）

## Decisions

### D1: Exclude 粒度选择整个配置目录

选择排除 `/.claude`（而非逐个排除 memory、settings 等子目录）。

**理由**：agent 运行时会生成多种文件，且列表可能随 Claude Code 版本变化。排除整个目录一劳永逸，且 `.git/info/exclude` 不影响已 tracked 的文件，仓库自带的 `.claude/CLAUDE.md` 等内容不受影响。

**备选方案**：逐个排除 `/.claude/memory/`、`/.claude/settings.json` 等 → 维护成本高，容易遗漏。

### D2: Exclude 逻辑移至 WorkspaceManager.create()

将 exclude 写入放在 `WorkspaceManager.create()` 中 worktree 创建之后，而非 SkillProvisioner.provision() 中。

**理由**：exclude 是 workspace 级别的隔离需求（所有 agent 都需要），不是 skill 分发级别的关注点。即使未来有不 provision skills 的 workspace，exclude 仍然需要。

**备选方案**：留在 provisioner 但扩大排除范围 → 职责混乱，provisioner 不应管理 workspace 级别的 git 配置。

### D3: 新增独立端点用于 task→main 合并

新增 `POST /v1/workspace/task/{task_id}/merge-to-main`，与现有的 `POST /v1/workspace/{workspace_id}/merge` 分离。

**理由**：task→main 是平台级操作（由上层系统调用），与 agent 的 workspace 级 merge 语义不同。分离端点使职责更清晰，避免 agent 误用 target_branch="main"。

**备选方案**：复用现有 merge 端点，通过 target_branch 参数区分 → agent 可能误传 target_branch="main"，安全性差。

### D4: GitOps 提供 write_exclude 公共方法

`_resolve_git_dir` + exclude 写入逻辑封装为 `GitOps.write_exclude(worktree_path, entries)` 公共方法，WorkspaceManager 调用。

**理由**：git info/exclude 的定位和写入涉及 worktree 的 `.git` 文件解析，这是 git 操作层的能力，应由 GitOps 提供。

## Risks / Trade-offs

- **[仓库自带 `.claude/` tracked 内容被误排除的风险]** → `.git/info/exclude` 只影响 untracked 文件，已 tracked 的文件不受影响，无实际风险
- **[agent 需要新增 `.claude/` 下文件到 git 的场景]** → 可通过 `git add -f` 强制添加，但通常不需要
- **[provisioner 移除 `_write_git_exclude` 后 skill 文件可能被提交]** → 由目录级排除 `/.claude` 覆盖，skill 子目录自然被排除
