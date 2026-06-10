## Why

当前 Agent 回复只有纯文本，前端无法区分文本、代码、HTML、图片、diff 等不同内容类型。需要一套"输出技能"机制，让 Agent 主动标记输出内容的类型，前端据此渲染对应的富媒体卡片（HTML 渲染、图片展示、文件附件、可编辑 diff、网页预览）。

## What Changes

- 新增 `SkillRule`（RuleEngine），注入 system_prompt 告诉 Agent 可用 ` ```aka_yhy ` 块输出 5 种富媒体卡片
- 新增 `aka_yhy` 统一标记协议：`type: html-render | image | attachment | diff | preview`
- 前端新增 Block Reducer 纯函数，解析 markdown AST 中的 aka_yhy 块，分发到对应卡片组件
- 前端新增 5 种卡片组件：HtmlCard、ImageCard、AttachmentCard、DiffCard、PreviewCard
- AgentEnd 新增 workspace 文件读取 API、diff API、预览 HTTP 服务
- Go 后端 proxy 透传所有新增 API，不存 artifact 元数据
- 前端 ChatMessage 结构从 `content: string` 升级为 `blocks: MessageBlock[]`

## Capabilities

### New Capabilities
- `output-skill-rule`: SkillRule 注入输出技能 prompt，定义 aka_yhy 块的 5 种类型及其格式
- `block-reducer`: 前端 Block Reducer 纯函数，解析 aka_yhy 块生成 Block 数组，替代现有的 streamingContent 字符串拼接
- `output-cards`: 前端 5 种卡片组件（HtmlCard、ImageCard、AttachmentCard、DiffCard、PreviewCard）及统一的 MessageBubble block 渲染逻辑
- `workspace-file-api`: AgentEnd workspace 文件读取/写入/diff API，Go 后端 proxy 透传
- `preview-service`: AgentEnd 启动本地 HTTP 预览服务，前端通过 iframe 加载预览

### Modified Capabilities
- `skill-provisioning`: 从只注入工具技能文件扩展到同时注入输出技能 prompt
- `rule-engine`: 新增 SkillRule（低优先级，不阻塞安全规则）

## Impact

- **AgentEnd**: 新增 SkillRule（rules/builtin.py）、workspace 文件/diff API（api/v1/workspace.py）、预览 HTTP 服务
- **Go Backend**: 新增 proxy 路由（workspace file/diff/preview），纯透传
- **Frontend**: Block 类型定义、Block Reducer、5 种卡片组件、MessageBubble 改为 block 渲染、chat store 升级
- **新增依赖**: react-diff-viewer-continued（diff 渲染）
- **Adapter 层**: 零改动（aka_yhy 块作为普通 TEXT 事件传输）
