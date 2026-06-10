## Context

当前项目是三端 monorepo（frontend React + backend Go + agentend Python）。前端聊天消息仅存储在 useReducer 内存中，刷新即丢失。Agent 仅在创建时可设置名称，头像按 agent type 硬编码。Task 的 repoPath 字段在后端已存在但前端未暴露输入项，且无运行时校验。

现有架构：
- Task → 多个 Session（每个 agent 一个）→ 消息通过 SSE 流式返回
- 后端已有 Task、Session 模型（GORM）
- 前端使用 React Query 管理服务端状态，useChatStream hook 管理消息流

## Goals / Non-Goals

**Goals:**
- 消息持久化到后端数据库，重新打开 task 可加载完整历史
- Agent 头像支持自定义图片上传，未设置时 DiceBear 根据 name 生成
- Agent 名称和头像可在创建后修改
- Task 创建时支持设置 repoPath
- 每次发消息前校验 repoPath 有效性

**Non-Goals:**
- 不修改 agentend 的上下文管理机制（adapter 自己维护 history）
- 不做消息搜索、导出功能
- 不做多端消息同步（当前只有单用户场景）
- 不做消息分页加载（初期全量返回）

## Decisions

### D1: Message 存储在 backend，按 task_id 组织

**选择**: backend 新增 Message 表，外键 task_id，一个 task 的所有 session 的消息存在一起。

**理由**: 前端聊天窗口以 task 为单位，一个 task 就是一个群聊窗口，多个 agent 的消息混在一起展示。按 task_id 组织直接匹配前端消费模式。

**备选**: 按 session_id 分组。否决，因为前端需要聚合多个 session 的消息，按 session 存需要额外合并逻辑。

### D2: 消息在 SSE 流结束时批量保存

**选择**: user message 在收到请求时立即保存，agent message 在 SSE 流完成（done 事件）后保存完整内容。

**理由**: user message 是确定性的，立即持久化无风险。agent message 在流式过程中内容不断变化，完成后保存可避免中间状态写入。

**备选**: 逐事件实时保存。否决，写入过于频繁且流式中间态无消费方。

### D3: Agent 头像使用文件存储

**选择**: 后端接收图片上传，存储到本地文件系统（`uploads/avatars/`），数据库存文件路径。

**理由**: 当前项目规模小，无需对象存储服务。文件路径方案简单直接，后续可迁移到 S3。

**备选**: Base64 存数据库。否决，大头像会导致数据库膨胀。对象存储。过度设计，当前阶段不需要。

### D4: 头像 fallback 使用 DiceBear

**选择**: 使用 DiceBear API（`https://api.dicebear.com/7.x/initials/svg?seed=<name>`）生成 initials 风格头像。

**理由**: 确定性生成（同名同头像）、无需额外依赖库、SVG 格式轻量。前端在无自定义头像时使用 DiceBear URL 作为 img src。

**备选**: 前端 npm 包。需要引入额外依赖，而 API 方式零依赖。但考虑离线场景和速度，后续可切换为 npm 包。

### D5: repoPath 校验为独立步骤，先校验再发送

**选择**: 前端发消息时分两步：先 `POST /api/validate-repo-path`，成功后再调用 `/api/tasks/:taskId/run`。

**理由**: 校验和发送解耦，校验失败时可以给出明确错误提示而不会触发消息发送流程。

**备选**: 合并在 run 请求中。否决，校验失败时消息已经到达后端，错误处理更复杂。

### D6: agentend 仅提供校验接口，不存储状态

**选择**: agentend 暴露 `POST /validate-repo-path` 接口，接收 path，返回 `{valid: bool, errors: []}`。

**理由**: agentend 是无状态的执行层，校验是即时操作，不需要持久化。

## Risks / Trade-offs

- [Message 体积] → 长对话场景下全量加载可能较慢 → 初期接受，后续加消息分页
- [头像文件存储] → 单机部署无问题，多实例部署需要共享存储 → 当前单机足够，迁移时切换为对象存储
- [DiceBear 外部依赖] → 离线环境无法生成头像 → 可降级为纯色+首字母的本地方案
- [repoPath 校验延迟] → 每次发消息多一次网络请求 → 校验通常很快（<100ms），用户体验影响可忽略
