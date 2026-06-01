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
      "headTime": "string"      // HEAD commit relative time
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
