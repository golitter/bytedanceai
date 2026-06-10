## ADDED Requirements

### Requirement: repoPath 校验接口（agentend）
agentend SHALL 提供 `POST /validate-repo-path` 接口，接收 `{"repo_path": "/some/path"}`，检查路径是否存在且已 git init，返回校验结果。

#### Scenario: 路径有效
- **WHEN** 传入存在的路径且该目录已 git init
- **THEN** 返回 `{"valid": true, "errors": []}`

#### Scenario: 路径不存在
- **WHEN** 传入的路径在文件系统中不存在
- **THEN** 返回 `{"valid": false, "errors": ["路径不存在: /some/path"]}`

#### Scenario: 路径存在但未 git init
- **WHEN** 传入的路径存在但不是 git 仓库
- **THEN** 返回 `{"valid": false, "errors": ["路径不是 git 仓库: /some/path"]}`

#### Scenario: 缺少 repo_path 参数
- **WHEN** 请求 body 中无 repo_path 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: repoPath 校验转发（backend）
backend SHALL 提供 `POST /api/validate-repo-path` 端点，接收前端请求后转发给 agentend 执行校验，返回校验结果。

#### Scenario: 转发校验成功
- **WHEN** 前端发送 `POST /api/validate-repo-path` 并传入有效路径
- **THEN** backend 转发给 agentend 校验，返回 agentend 的校验结果

#### Scenario: agentend 不可达
- **WHEN** agentend 服务不可用
- **THEN** backend 返回 HTTP 503，提示 agent 服务不可用

### Requirement: 发消息前校验 repoPath
前端 SHALL 在每次发送消息前先调用 repoPath 校验接口，校验通过后才发送消息，校验失败则展示错误提示。

#### Scenario: 校验通过后发送
- **WHEN** 用户点击发送，repoPath 校验返回 valid=true
- **THEN** 前端正常发送消息到 `/api/tasks/:taskId/run`

#### Scenario: 校验失败不发送
- **WHEN** 用户点击发送，repoPath 校验返回 valid=false
- **THEN** 前端展示具体错误信息（路径不存在/未 git init），不发送消息

#### Scenario: task 无 repoPath
- **WHEN** task 未设置 repoPath
- **THEN** 跳过校验，直接发送消息

### Requirement: Task 创建时设置 repoPath
前端 NewChatDialog SHALL 支持输入 repoPath，在创建 task 时传入。

#### Scenario: 创建 task 时指定 repoPath
- **WHEN** 用户在 NewChatDialog 中填写 repoPath 并创建
- **THEN** 创建 task 请求包含 repo_path 字段

#### Scenario: repoPath 可选
- **WHEN** 用户不填写 repoPath
- **THEN** 正常创建 task，repo_path 为空
