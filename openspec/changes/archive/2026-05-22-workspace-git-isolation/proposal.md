## Why

Agent 在 worktree 中运行时会在 `.claude/` 下生成 memory、settings、scheduled_tasks 等文件，`git add -A` 会将这些一并提交，合并时带入 task 分支甚至 main。当前仅排除 `/.claude/skills/<name>`，不够完整。同时 merge API 默认直接合到 main，跳过了 `task/{task_id}` 这一层，不符合"agent→task→main"的三级工作流设计。

## What Changes

- 将 agent 配置目录（`.claude`/`.opencode`）的 git exclude 逻辑从 `SkillProvisioner` 移至 `WorkspaceManager`，在 worktree 创建时写入 `/.claude` 到 `.git/info/exclude`，从 provisioner 中移除 `_write_git_exclude`
- `MergeRequest.target_branch` 默认值从 `"main"` 改为 `None`，让 agent merge 默认走到 `task/{task_id}`
- 新增 `POST /v1/workspace/task/{task_id}/merge-to-main` 平台级端点，将 task 分支合并到 main

## Capabilities

### New Capabilities

- `config-dir-exclusion`: workspace 创建时自动将 agent 配置目录（`.claude` 或 `.opencode`）写入 `.git/info/exclude`，防止 agent 运行时生成的非代码文件被提交

### Modified Capabilities

- `workspace-management`: merge 默认目标从 main 改为 task 分支；新增 `merge_task_to_main(task_id)` 方法和对应 API 端点
- `workspace-isolation`: 移除 provisioner 中的 `_write_git_exclude`，workspace 创建时的 exclude 写入改由 WorkspaceManager 负责
- `skill-provisioning`: GitOps 新增 `write_exclude(git_dir, entries)` 方法

## Impact

- `src/workspace/manager.py` — `create()` 中新增 exclude 写入；新增 `merge_task_to_main()` 方法
- `src/workspace/git_ops.py` — 新增 `write_exclude()` 方法
- `src/skills/provisioner.py` — 移除 `_write_git_exclude` 及其调用
- `src/api/v1/workspace.py` — `MergeRequest.target_branch` 默认值改为 `None`；新增 merge-to-main 端点
