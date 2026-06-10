## Context

前端 ImPage 当前是两栏 flex 布局：左侧 ConversationList（280px）+ 右侧 ChatArea（flex-1）。无全局导航入口，品牌信息（AgentHub）硬编码在 ConversationList header 中。路由只有 `/` 和 `/agent/:sessionId`。

用户希望改为 QQ 风格三栏：IconSidebar（56px）+ ConversationList（280px）+ ChatArea（flex-1）。Demo 已在 `demo/qq-layout.html` 验证视觉方案。

现有技术栈：React 19 + Tailwind CSS + shadcn/ui + Zustand + TanStack React Query。AgentAvatar 已使用 DiceBear 生成头像。

## Goals / Non-Goals

**Goals:**
- 三栏布局：IconSidebar + SessionPanel + ChatArea
- IconSidebar 支持多 tab 切换（chat 激活，其他灰色占位）
- 用户头像 hover 弹出个人卡片
- ConversationList 去掉品牌 Header，"+" 按钮移到搜索框旁
- Agent 头像升级为 DiceBear Bottts 风格
- 现有聊天功能零回归

**Non-Goals:**
- 不做通讯录/管理/设置页面的真实内容（仅占位）
- 不做未读消息计数（需后端 API 支持，后续迭代）
- 不做最后消息预览替代 taskTitle（需后端 API 支持）
- 不做会话分组（今天/昨天/更早）
- 不做右键菜单（置顶/删除/标记已读）
- 不改路由结构

## Decisions

### Decision 1: Tab 状态放 Zustand store 而非组件 state

**选择**: `activeTab` 放入 `stores/chat.ts` 的 Zustand store
**替代方案**: ImPage 组件内 `useState`
**理由**: ImPage 和 IconSibling 都要读写 activeTab；未来 contacts/admin/settings 页面也需要感知当前 tab；Zustand 与现有 `currentSessionId` 保持一致。

### Decision 2: "+" 按钮移到搜索框右侧

**选择**: 在搜索栏 `.search-bar` 内追加一个 "+" icon button
**替代方案**: 移到 IconSidebar 的聊天图标上
**理由**: 搜索框旁最自然（用户寻找对话时顺便新建）；IconSidebar 空间太小，点击区域不够明确。

### Decision 3: 非 chat tab 隐藏中间面板

**选择**: 切到 contacts/admin/settings 时隐藏 ConversationList，主区域显示占位页
**替代方案**: 保持三栏，中间面板切换内容
**理由**: 通讯录/管理/设置目前没有真实内容，保持三栏只显示空列表没意义。等真实页面就绪后再决定是否保持三栏。

### Decision 4: 用户头像用 DiceBear Notionists，Agent 头像用 Bottts

**选择**: 用户用 `notionists` 风格（手绘人物），Agent 用 `bottts` 风格（机器人）
**替代方案**: 统一用 `bottts` 或 `avataaars`
**理由**: 两种风格形成视觉区分——用户是人，Agent 是机器人。与现有 agent 色彩系统兼容（每个 Agent 按名字 seed 生成独特外观）。

## File Changes

```
新建:
  frontend/src/components/layout/IconSidebar.tsx    ~100行

修改:
  frontend/src/pages/ImPage.tsx                     ~15行（两栏→三栏）
  frontend/src/components/im/ConversationList.tsx   ~10行（移除 header）
  frontend/src/stores/chat.ts                       ~8行（加 activeTab）
  frontend/src/components/chat/AgentAvatar.tsx      已改（DiceBear 7.x→9.x bottts）

不动:
  frontend/src/components/chat/ChatArea.tsx
  frontend/src/components/chat/MessageList.tsx
  frontend/src/components/chat/MessageInput.tsx
  frontend/src/components/im/ConversationItem.tsx
  frontend/src/main.tsx（路由不改）
```
