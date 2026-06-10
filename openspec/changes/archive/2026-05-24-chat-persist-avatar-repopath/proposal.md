## Why

当前前端聊天记录仅存在于内存中（useReducer），刷新页面或切换 task 后消息全部丢失。Agent 创建后无法修改名称和头像，缺少个性化能力。Task 创建时无法设置 repoPath，且缺少运行时路径有效性校验，导致 agent 可能基于无效路径执行任务。

## What Changes

- 新增 Message 持久化模型，将聊天消息存储到后端数据库（按 task_id 组织），支持重新打开时加载完整历史
- 支持上传自定义图片作为 Agent 头像，未定义时使用 DiceBear 库根据 agent name 生成确定性头像
- 支持创建后修改 Agent 名称和头像
- 前端 NewChatDialog 增加 repoPath 输入项
- 新增 repoPath 校验流程：发消息前先校验路径有效性（backend 转发 agentend 执行），校验通过再发送
- agentend 新增路径校验接口，检查路径是否存在、是否已 git init

## Capabilities

### New Capabilities
- `message-persistence`: 聊天消息持久化存储与加载，按 task_id 组织，支持 user/agent 两种角色消息
- `agent-avatar`: Agent 头像管理，支持自定义图片上传和 DiceBear 自动生成
- `repopath-validation`: repoPath 运行时校验，发消息前检查路径有效性

### Modified Capabilities
- `chat-ui`: 需要集成消息历史加载、头像展示/编辑、repoPath 输入
- `task-api`: 需要支持 repoPath 设置和消息查询接口
- `session-agent-api`: 需要在发消息流程中插入 repoPath 校验步骤

## Impact

- **Backend**: 新增 Message 模型及 CRUD API，新增头像上传/存储 API，新增 repoPath 校验转发接口
- **Frontend**: 消息加载逻辑（useChatStream 改造）、头像上传/展示组件、NewChatDialog 增加 repoPath 输入、发送消息前校验流程
- **Agentend**: 新增路径校验接口（检查路径存在性 + git init 状态）
- **Dependencies**: 前端新增 DiceBear 依赖
- **API**: 新增 `GET /tasks/:id/messages`、`POST /validate-repo-path`、头像上传接口
