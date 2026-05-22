# Workspace 分支清理

## 问题

`cleanup()` 只执行 `git worktree remove`，没有删除对应的 git 分支。清理后残留如下分支：

```
  agent/sess-aaa/333
  agent/sess-bbb/333333
* main
  task/333
  task/333333
```

## 分支模型

每个 task 的 workspace 会创建两类分支：

- **agent 分支**：`agent/{session_id}/{task_id}` — 每个 agent session 对应一个
- **task 分支**：`task/{task_id}` — 同一 task 共享的基础分支

## 修复方案

### 1. GitOps 新增 `branch_delete()`

`src/workspace/git_ops.py` — 执行 `git branch -D <name>` 删除指定分支。

### 2. `cleanup()` 清理 agent 分支

`src/workspace/manager.py` — worktree 移除成功后，立即删除对应的 agent 分支。

### 3. `cleanup_by_task()` 清理 task 分支

`src/workspace/manager.py` — 所有 agent worktree 清理完毕后，删除 task 分支。

### 4. recovery 孤儿清理同时删除分支

`src/workspace/recovery.py` — 移除孤儿 worktree 后，也删除其对应分支。

### 5. shutdown 使用 `cleanup_by_task()`

`src/app/main.py` — 关闭时按 task_id 分组调用 `cleanup_by_task()`，确保 task 分支被正确清理。
