## Why

当前项目中 Agent 的身份信息分散在多处（constants.ts 的 AGENT_NAMES/DESCRIPTIONS、后端 SessionAgent model、agentend config.yaml），但用户在聊天界面只能看到头像和名称，无法快速了解一个 Agent 的能力、技能和元数据。需要一个轻量的 hover 卡片和详情页来承载这些信息。

## What Changes

- 新增 **Agent Hover Card**：在聊天气泡的 Agent 头像上 hover 时弹出极简浮动卡片，展示身份（name + type + status）、Skills（name + description 截断一行）、元数据（session_id）、以及跳转详情页的链接
- 新增 **Agent 详情页**：路由 `/agent/:sessionId`，展示完整的 Agent 信息，包括 Skills 列表（不截断）、元数据（session/task/workspace/创建时间/消息数）
- 后端新增 `AgentSkill` model 和两个 GET 接口（profile / detail），Skills 上报接口先 mock 硬编码数据，后续大修
- 实现前先交付纯 HTML demo 供审批，确认后再正式实现

## Capabilities

### New Capabilities

- `agent-hover-card`: 聊天气泡头像 hover 触发的 Agent 浮动卡片，展示身份、Skills 摘要、元数据、详情跳转链接
- `agent-detail-page`: Agent 详情页路由及完整信息展示（Skills 列表、元数据、统计）
- `agent-skill-model`: 后端 AgentSkill 数据模型及 Mock API（profile / detail），为前端卡片和详情页提供数据

### Modified Capabilities

（无现有 spec 需要修改）

## Impact

- **前端**：MessageBubble 组件包裹 AgentAvatar 添加 hover 交互；新增 AgentHoverCard 组件；新增 AgentProfilePage 页面和路由；可能需要安装 `@radix-ui/react-popover`
- **后端**：新增 AgentSkill GORM model；新增 2 个 GET 接口；Skills 数据先 mock
- **Agent 端**：本次不改动，Skills 上报后续大修
- **契约层**：本次不改 contracts/schemas，API 返回结构直接在后端定义
