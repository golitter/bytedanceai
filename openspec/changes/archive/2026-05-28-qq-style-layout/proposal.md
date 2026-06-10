## Why

当前前端是两栏布局（会话列表 + 聊天区），缺少全局导航入口，功能扩展困难。需要改为 QQ 经典三栏布局：图标侧边栏 + 会话列表 + 聊天主区，层次分明且为后续功能（通讯录、管理页面、设置）预留扩展位。

参考 Demo：`demo/qq-layout.html`

## What Changes

- 新建 `IconSidebar` 组件（56px 窄条），顶部用户头像 + hover 个人卡片，中部功能图标（聊天/通讯录/管理），底部设置
- `ImPage.tsx` 从两栏改为三栏，IconSidebar 的 tab 切换控制中间面板内容
- `ConversationList.tsx` 移除顶部品牌 Header，"+" 新建按钮移到搜索框右侧
- `stores/chat.ts` 新增 `activeTab` 状态（`chat` | `contacts` | `admin` | `settings`）
- `AgentAvatar` 的 DiceBear 风格从 `initials` 升级为 `bottts`（已改）
- 用户头像使用 DiceBear `notionists` 风格

## Capabilities

### New Capabilities
- `icon-sidebar`: 窄条图标导航栏，顶部用户头像（hover 卡片），中部 tab 图标，底部设置
- `layout-three-column`: ImPage 三栏布局，IconSidebar tab 控制中间面板 + 右侧主区渲染

### Modified Capabilities
- `conversation-list`: 移除品牌 Header，搜索框行集成 "+" 按钮
- `chat-store`: 新增 `activeTab` 导航状态

## Impact

- **Frontend**: 新建 IconSidebar 组件，改造 ImPage 布局，瘦身 ConversationList，扩展 chat store
- **Backend**: 零改动
- **AgentEnd**: 零改动
- **新增依赖**: 无
