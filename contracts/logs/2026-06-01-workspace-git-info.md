# 2026-06-01 — Workspace Git Info API

## 变更类型
新增 API 端点（纯查询，无数据模型变更）

## 变更内容

### 新增端点
- **Backend**: `GET /api/workspace/task/:taskId/git-info`
- **Agentend**: `GET /v1/workspace/task/{task_id}/git-info`

### 响应结构
```json
{
  "repoPath": "string",
  "branches": [
    {
      "name": "string",        // 分支名：main / task/{id} / agent/{sid}/{id}
      "headHash": "string",     // HEAD commit 短 hash
      "headMsg": "string",      // HEAD commit message
      "headAuthor": "string",   // HEAD commit author
      "headTime": "string",     // HEAD commit relative time
      "exists": true            // git ref 是否真实存在；缺失时 head 字段为空
    }
  ],
  "commits": [
    {
      "hash": "string",         // 短 hash
      "fullHash": "string",     // 完整 hash
      "msg": "string",
      "author": "string",
      "lane": "string",         // 所属分支（用于 Graph lane 分配）
      "time": "string",
      "parentHashes": ["string"] // 父 commit hashes（用于 Graph 连线）
    }
  ]
}
```

## 影响范围
- **Frontend**: `RightSidebar` → `GitGraphPanel` 组件消费此 API
- **Backend**: `WorkspaceHandler.TaskGitInfo` proxy 到 Agentend
- **Agentend**: `workspace.py` 新增 `get_task_git_info` 路由
- **Contracts**: 无 YAML schema 变更，无需 `make generate`

## 2026-06-01 补充说明

### 变更原因
Git Graph 需要区分“逻辑上应该出现的分支”和“git ref 真实存在的分支”，避免 task 分支缺失时静默消失；同时当 agent 分支 fast-forward 合并到 task 分支后，两条 ref 会指向同一 commit，前端需要保留两个 head 标签。

### 对比结果
- `branches[]` 新增可选字段 `exists: boolean`。
- `exists=false` 表示该分支是工作区记录推导出的预期分支，但当前 git ref 不存在；此时 `headHash/headMsg/headAuthor/headTime` 为空字符串。
- `exists=true` 表示该分支 ref 真实存在，可用于 Git Graph head 标记。

### task-base 语义
- 不存在 `task-base` 分支。
- `worktrees/{task_id}/task-base` 是固定 worktree 目录。
- 该目录 checkout 的分支是 `task/{task_id}`。
- Agent 分支命名仍为 `agent/{session_id}/{task_id}`，从 `task/{task_id}` 派生。

### 跨端影响
- **Frontend**: Git Graph 使用 `exists` 展示缺失分支状态；同一 commit 上可同时展示 agent/task/main 多个 head 标签。
- **Backend**: 仍仅 proxy Agentend 响应，无额外结构转换。
- **AgentEnd**: `get_task_git_info` 返回 `exists` 字段；task-base worktree 生命周期修正为目录挂载 `task/{task_id}` 分支。
- **Contracts**: 仍无 YAML schema 变更，无需 `make generate`。
