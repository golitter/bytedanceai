## Context

当前聊天气泡中的 AgentAvatar 组件（32px 头像 + 状态点）只展示视觉标识，用户无法快速了解 Agent 的技能和身份。Agent 元数据散落在 frontend constants.ts（AGENT_NAMES/DESCRIPTIONS）、后端 SessionAgent model、agentend config.yaml 中，没有统一的展示入口。

前端技术栈：React 19 + Vite 8 + Tailwind + shadcn/ui（Radix 基础），目前只安装了 @radix-ui/react-dialog。

## Goals / Non-Goals

**Goals:**
- 聊天气泡头像 hover 时弹出极简 Agent 资料卡
- 点击卡片链接跳转 Agent 详情页（/agent/:sessionId）
- 后端提供 profile / detail API，Skills 数据先 mock
- 先交付纯 HTML demo 供审批

**Non-Goals:**
- Skills 上报接口（后续大修）
- 会话列表头像的 hover 卡片
- Agent 能力标签、简介展示（后续扩展）
- 用户自定义 Skills 上传 UI（后续扩展）

## Decisions

### D1: Hover 卡片实现方式 — @radix-ui/react-popover

**选择**: 使用 Radix Popover 组件（`@radix-ui/react-popover`），`trigger="hover"` 模式。

**替代方案**:
- 纯 CSS hover（position:absolute + opacity transition）：轻量但缺少无障碍、边界检测、焦点管理
- shadcn Popover：底层就是 Radix Popover，可直接 `npx shadcn@latest add popover` 安装

**理由**: 项目已用 Radix 体系（Dialog），Popover 保持一致。shadcn 封装提供现成的 Tailwind 样式。

### D2: 数据获取 — hover 时按需请求

**选择**: hover 触发时调用 `GET /api/sessions/:sid/profile`，用 React Query 缓存。

**替代方案**:
- 消息加载时预取所有 Agent profile：浪费请求
- 前端硬编码：Skills 数据后端已有，硬编码会不同步

**理由**: 按需请求 + React Query 缓存，hover 只在第一次触发网络请求。

### D3: 后端 Skills 数据 — 先 Mock

**选择**: GET 接口返回硬编码的 mock Skills 数据（taskctl + render），不创建 AgentSkill 表。

**替代方案**:
- 完整实现 AgentSkill model + 注册流程：用户明确表示后续大修，现在投入浪费

**理由**: 前端组件和交互是核心，数据结构后续定型。Mock 保持 API 契约，前端开发不受阻。

### D4: 详情页路由 — /agent/:sessionId

**选择**: 新增路由 `/agent/:sessionId`，使用现有 react-router。

**理由**: sessionId 是唯一标识，与后端 API 路径一致。独立页面承载完整信息。

### D5: 先 HTML Demo 再实现

**选择**: 先写一个独立的 .html 文件，包含 hover 卡片和详情页的完整视觉原型。

**理由**: 视觉交互需要快速确认，纯 HTML 可脱离项目直接在浏览器预览，减少返工。

## Risks / Trade-offs

- [Mock 数据后续需替换] → API 路径和返回结构提前定义好，后续只需改后端实现
- [Hover 延迟体验] → 设 300ms show / 200ms hide delay，避免闪烁和误触
- [Popover 定位遮挡] → Radix Popover 内置碰撞检测，卡片自动调整方向
- [Skills 数量不可控] → 卡片内限制显示 N 个 + "+M 更多"，详情页不限制
