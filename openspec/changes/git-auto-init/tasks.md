## 1. Agentend — Git 初始化端点

- [ ] 1.1 在 `agentend/src/api/v1/validate.py` 新增 `InitGitRepoRequest` / `InitGitRepoResponse` Pydantic Model
- [ ] 1.2 在 `validate.py` 新增 `POST /v1/init-git-repo` 端点：校验路径 → 调用 `GitOps.init_repo()` → 返回结果

## 2. Go Backend — 代理路由

- [ ] 2.1 在 `backend/pkg/agentend_client/client.go` 新增 `InitGitRepoResult` struct 和 `InitGitRepo()` 方法
- [ ] 2.2 在 `backend/internal/controller/impl/task_controller.go` 新增 `InitGitRepoReq` struct 和 `InitGitRepo()` handler
- [ ] 2.3 在 `RegisterRoutes` 中注册 `rg.POST("/init-git-repo", ctrl.InitGitRepo)`

## 3. Frontend — API 与 UI 文本

- [ ] 3.1 在 `frontend/src/lib/api.ts` 新增 `initGitRepo()` 函数
- [ ] 3.2 在 `frontend/src/lib/ui-text.ts` 新增 `UI_ACTIONS.INIT_GIT` / `UI_STATUS.INITIALIZING_GIT` / `UI_MESSAGES.GIT_INIT_*` / `UI_ERRORS.GIT_INIT_FAILED`

## 4. Frontend — RepoPathInput 确认流程

- [ ] 4.1 新增 `needsGitInit` / `confirmInput` / `initError` / `initializing` 状态
- [ ] 4.2 修改 `handleValidate`：检测"不是 git 仓库"错误 → 进入确认态而非报错
- [ ] 4.3 新增 `handleInitGit` 函数：匹配目录名 → 调用 API → 成功后标记 validated
- [ ] 4.4 JSX 新增黄色确认框（内联输入框下方）：目录名提示 + 文本输入 + 初始化按钮 + 取消按钮
