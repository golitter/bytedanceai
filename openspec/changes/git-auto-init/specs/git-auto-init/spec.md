## ADDED Requirements

### Requirement: Git 初始化接口（agentend）
agentend SHALL 提供 `POST /v1/init-git-repo` 接口，接收 `{"repo_path": "/some/path"}`，校验路径存在且非 Git 仓库后执行 git 初始化（`git init && git add -A && git commit -m "init" && git branch -M main`），返回操作结果。

#### Scenario: 成功初始化
- **WHEN** 传入存在的目录路径且该目录不是 Git 仓库
- **THEN** 执行 git 初始化并返回 `{"success": true, "errors": []}`

#### Scenario: 路径不存在
- **WHEN** 传入的路径在文件系统中不存在
- **THEN** 返回 `{"success": false, "errors": ["路径不存在: /some/path"]}`

#### Scenario: 路径不是目录
- **WHEN** 传入的路径存在但不是目录
- **THEN** 返回 `{"success": false, "errors": ["路径不是目录: /some/path"]}`

#### Scenario: 路径已是 Git 仓库
- **WHEN** 传入的路径已经是一个 Git 仓库
- **THEN** 返回 `{"success": false, "errors": ["路径已经是 git 仓库: /some/path"]}`

#### Scenario: 初始化失败（空目录等）
- **WHEN** git 初始化命令执行失败（如空目录无法 commit）
- **THEN** 返回 `{"success": false, "errors": ["Git 初始化失败: /some/path"]}`

### Requirement: Git 初始化代理路由（backend）
backend SHALL 提供 `POST /api/init-git-repo` 端点，接收前端请求后转发给 agentend 执行初始化，返回操作结果。

#### Scenario: 转发初始化成功
- **WHEN** 前端发送 `POST /api/init-git-repo` 并传入有效路径
- **THEN** backend 转发给 agentend 执行初始化，返回操作结果

#### Scenario: agentend 不可达
- **WHEN** agentend 服务不可用
- **THEN** backend 返回 HTTP 503，提示 agent 服务不可用

#### Scenario: 缺少 repo_path 参数
- **WHEN** 请求 body 中无 repo_path 字段
- **THEN** 返回 HTTP 400，提示参数错误

### Requirement: 前端 Git 初始化确认流程
前端 RepoPathInput 组件 SHALL 在校验返回"不是 git 仓库"错误时，展示黄色确认框要求用户输入路径最后一段目录名，匹配后可调用初始化 API。

#### Scenario: 展示确认框
- **WHEN** 校验接口返回 `valid=false` 且错误包含"不是 git 仓库"
- **THEN** 显示黄色确认框，提示用户输入目录名以确认初始化

#### Scenario: 输入匹配 → 按钮激活
- **WHEN** 用户输入的文本与路径最后一段精确匹配
- **THEN** "初始化 Git" 按钮变为可点击状态

#### Scenario: 输入不匹配 → 按钮禁用
- **WHEN** 用户输入的文本与路径最后一段不匹配
- **THEN** "初始化 Git" 按钮保持禁用，提示名称不匹配

#### Scenario: 确认初始化成功
- **WHEN** 用户点击"初始化 Git"且 API 返回 `success=true`
- **THEN** 确认框消失，自动进入校验通过状态，路径显示绿色边框

#### Scenario: 确认初始化失败
- **WHEN** 用户点击"初始化 Git"且 API 返回 `success=false`
- **THEN** 确认框内显示错误信息，用户可重试

#### Scenario: 取消确认
- **WHEN** 用户点击"取消"按钮
- **THEN** 确认框消失，恢复到初始状态

#### Scenario: 初始化中禁用交互
- **WHEN** 初始化 API 请求进行中
- **THEN** 所有输入和按钮禁用，按钮显示"正在初始化 Git..."
