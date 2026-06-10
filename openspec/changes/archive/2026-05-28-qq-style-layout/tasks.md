## 1. Store: 新增 activeTab 导航状态

- [x] 1.1 在 `stores/chat.ts` 的 `ChatStoreState` 新增 `activeTab: 'chat' | 'contacts' | 'admin' | 'settings'`，默认 `'chat'`
- [x] 1.2 新增 `setActiveTab` action
- [x] 1.3 导出 `useActiveTab` hook（与 `useChatNav` 平行）

## 2. 新建 IconSidebar 组件

- [x] 2.1 新建 `frontend/src/components/layout/IconSidebar.tsx`
- [x] 2.2 顶部：用户头像（DiceBear Notionists, seed=tln），hover 弹出个人卡片（名称 + 昵称 + 编辑/退出按钮，目前按钮为占位）
- [x] 2.3 中部：chat/contacts/admin 三个 NavItem，chat 为激活态（`bg-primary-soft text-primary`），其他 disabled（`opacity-40`）
- [x] 2.4 底部：settings NavItem（disabled），`mt-auto` 推到底部
- [x] 2.5 每个 NavItem 结构：icon（lucide-react） + 文字标签（9px），宽 44px 高 44px，`rounded-[10px]`
- [x] 2.6 点击 chat NavItem 调用 `setActiveTab('chat')`，disabled 项不响应

## 3. 改造 ImPage 为三栏布局

- [x] 3.1 `ImPage.tsx` 顶层 `flex h-screen` 内新增 `<IconSidebar />` 作为第一个子元素
- [x] 3.2 根据 `activeTab` 条件渲染：`chat` 时显示 ConversationList + ChatArea，其他 tab 隐藏 ConversationList 并显示占位页面
- [x] 3.3 占位页面：图标 + 标题 + "功能开发中，敬请期待" 文案

## 4. ConversationList 瘦身

- [x] 4.1 移除 `ConversationList.tsx` 顶部品牌 Header（AgentHub + "+" 按钮，约 8 行）
- [x] 4.2 搜索框行右侧追加 "+" 新建对话按钮（复用原有 `showNewChat` 状态和 `NewChatDialog`）
- [x] 4.3 搜索框 padding 调整为 `px-3 py-2`，与 session-list 的 padding 对齐

## 5. 验证

- [x] 5.1 `make run-frontend` 启动，确认三栏布局正确渲染
- [x] 5.2 点击会话条目，ChatArea 正常显示消息
- [x] 5.3 Agent 头像显示为 DiceBear Bottts 风格（不再是首字母）
- [x] 5.4 IconSidebar 用户头像 hover 弹出卡片
- [x] 5.5 点击 disabled tab 无响应，chat tab 高亮可切换
- [x] 5.6 "+" 按钮正常打开新建对话弹窗
- [x] 5.7 搜索功能正常
