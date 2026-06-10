## 1. 纯 HTML Demo

- [x] 1.1 编写独立 HTML 文件，包含 Agent Hover Card 视觉原型（身份 + Skills + 元数据 + 详情链接）
- [x] 1.2 在同一 HTML 文件中编写 Agent 详情页视觉原型（头部 + Skills 完整列表 + 元数据）
- [x] 1.3 交付 demo 供审批，确认视觉风格

## 2. 后端 Mock API

- [x] 2.1 实现 `GET /api/sessions/:sid/profile` 接口，返回 mock 数据（agent_name, agent_type, avatar_url, status, session_id, skills[]）
- [x] 2.2 实现 `GET /api/sessions/:sid/detail` 接口，返回 mock 数据（profile 字段 + task_id, workspace_path, created_at, message_count）
- [x] 2.3 Mock Skills 数据硬编码 taskctl 和 render 两条记录

## 3. 前端基础组件

- [x] 3.1 安装 shadcn Popover 组件（`npx shadcn@latest add popover`）
- [x] 3.2 新增 `lib/api.ts` 中的 `fetchAgentProfile` 和 `fetchAgentDetail` 函数
- [x] 3.3 新增前端类型定义（AgentProfile, AgentDetail, AgentSkill）

## 4. Agent Hover Card 组件

- [x] 4.1 实现 `AgentHoverCard` 组件：身份区（头像 + name + type + status）
- [x] 4.2 实现 Skills 展示区：name + description 截断一行，超过 N 个显示 "+M 更多"
- [x] 4.3 实现元数据区（session_id 截断）和底部详情跳转链接
- [x] 4.4 在 `MessageBubble` 的 `AgentAvatar` 上接入 Popover，设置 300ms show / 200ms hide delay

## 5. Agent 详情页

- [x] 5.1 新增 `/agent/:sessionId` 路由
- [x] 5.2 实现 `AgentProfilePage` 页面布局（返回按钮 + AgentHeader + SkillList + AgentMeta）
- [x] 5.3 实现 `SkillCard` 组件（name + 完整 description + builtin 标记 + source）
- [x] 5.4 实现 `AgentMeta` 组件（session_id, task_id, workspace, 创建时间, 消息数）

## 6. 联调与验收

- [x] 6.1 启动三端服务，验证 Hover Card 在聊天气泡上正常触发和消失
- [x] 6.2 验证点击卡片跳转到详情页，数据正确展示
- [x] 6.3 验证 Skills 溢出场景（mock 多条 Skills 数据测试截断）
