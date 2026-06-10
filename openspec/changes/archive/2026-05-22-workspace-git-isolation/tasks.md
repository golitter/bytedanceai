## 1. GitOps 层：新增 write_exclude 方法

- [x] 1.1 将 `_resolve_git_dir` 逻辑从 provisioner 移至 `GitOps`，作为公共方法 `resolve_git_dir(worktree_path) -> Path | None`
- [x] 1.2 新增 `GitOps.write_exclude(worktree_path, entries)` 方法：解析 git dir，追加 entries 到 `info/exclude`（去重）
- [x] 1.3 补充 `write_exclude` 的单元测试：worktree 场景（.git 为文件）、普通 repo 场景（.git 为目录）、去重追加

## 2. WorkspaceManager 层：创建时写入 exclude + 新增 merge_task_to_main

- [x] 2.1 在 `WorkspaceManager.create()` 中，worktree 创建 + skill provisioning 之后，调用 `git_ops.write_exclude(worktree_path, ["/.claude" 或 "/.opencode"])`（根据 agent_type）
- [x] 2.2 新增 `WorkspaceManager.merge_task_to_main(repo_path, task_id)` 方法：调用 `git_ops.merge_branch(repo_path, task_branch_name(task_id), "main")`
- [x] 2.3 补充 `create` 集成测试：验证 worktree 创建后 `.git/info/exclude` 包含正确条目
- [x] 2.4 补充 `merge_task_to_main` 单元测试

## 3. API 层：修正 merge 默认值 + 新增端点

- [x] 3.1 `MergeRequest.target_branch` 默认值从 `"main"` 改为 `None`
- [x] 3.2 新增 `MergeTaskToMainRequest` 模型（含 `repo_path` 字段）
- [x] 3.3 新增 `POST /v1/workspace/task/{task_id}/merge-to-main` 端点，调用 `manager.merge_task_to_main()`
- [x] 3.4 补充 API 测试：merge 不传 target_branch 时合到 task 分支；merge-to-main 端点正常工作

## 4. 清理 SkillProvisioner

- [x] 4.1 从 `SkillProvisioner` 中移除 `_write_git_exclude` 方法及 `_resolve_git_dir` 方法
- [x] 4.2 从 `provision()` 方法中移除 `self._write_git_exclude(...)` 调用
- [x] 4.3 验证 skill 分发功能不受影响（provision 测试通过）
