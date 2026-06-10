## 1. Backend - Message 持久化

- [x] 1.1 新增 Message 模型（backend/internal/model/message.go）：id, task_id, session_id, role, content, agent_type, agent_name, created_at
- [x] 1.2 运行数据库迁移，创建 messages 表（task_id 索引）
- [x] 1.3 新增 `GET /api/tasks/:taskId/messages` 端点，返回按 created_at 升序的消息列表
- [x] 1.4 修改 `POST /api/tasks/:taskId/run` handler：收到请求时先保存 user message 到 Message 表
- [x] 1.5 修改 SSE 流处理：done 事件后保存完整 agent message 到 Message 表
- [x] 1.6 处理 SSE 中断场景：断开时保存已接收的部分 agent message

## 2. Backend - Agent 头像管理

- [x] 2.1 新增头像上传接口 `POST /api/agents/avatar`（multipart/form-data，限制 jpg/png/gif/webp，最大 2MB）
- [x] 2.2 实现文件存储逻辑：保存到 `uploads/avatars/`，数据库存路径
- [x] 2.3 新增 agent 属性更新接口（名称、头像 URL）
- [x] 2.4 新增头像静态文件服务路由

## 3. Backend - repoPath 校验

- [x] 3.1 新增 `POST /api/validate-repo-path` 端点，转发校验请求到 agentend
- [x] 3.2 处理 agentend 不可达场景（返回 503）

## 4. Agentend - repoPath 校验接口

- [x] 4.1 新增 `POST /validate-repo-path` 接口，检查路径是否存在
- [x] 4.2 检查路径是否已 git init（.git 目录是否存在）
- [x] 4.3 返回结构化校验结果 `{valid, errors[]}`

## 5. Frontend - 消息历史加载

- [x] 5.1 新增 API 调用函数 `getTaskMessages(taskId)`（frontend/src/lib/api.ts）
- [x] 5.2 修改 useChatStream hook：初始化时从后端加载历史消息填充到 state
- [x] 5.3 进入 task 页面时自动加载历史消息并滚动到底部

## 6. Frontend - Agent 头像与编辑

- [x] 6.1 安装 DiceBear 相关依赖（如需要）
- [x] 6.2 修改 AgentAvatar 组件：优先展示 avatar_url，fallback 用 DiceBear URL
- [x] 6.3 新增 agent 编辑弹窗组件（修改名称、上传头像）
- [x] 6.4 在 AgentAvatar / agent 名称上添加点击事件触发编辑弹窗
- [x] 6.5 新增头像上传 API 调用和 agent 属性更新 API 调用

## 7. Frontend - repoPath 输入与校验

- [x] 7.1 修改 NewChatDialog 组件：新增 repoPath 文本输入框
- [x] 7.2 创建 task 时将 repo_path 传入请求 body
- [x] 7.3 新增 validateRepoPath API 调用函数
- [x] 7.4 修改发送消息流程：发消息前先调用校验接口，展示 loading 状态
- [x] 7.5 校验失败时展示错误横幅，阻止消息发送
- [x] 7.6 task 无 repoPath 时跳过校验，直接发送

## 8. 契约更新

- [x] 8.1 更新 contracts/schemas/ 中相关 YAML（新增 Message schema、validate-repo-path 请求/响应 schema）
- [x] 8.2 运行 `make generate` 生成三端类型文件
- [x] 8.3 在 contracts/logs/ 写入变更记录
