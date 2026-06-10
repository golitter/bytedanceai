## Context

刚提交的 agent-name-and-auto-session 改动涉及 12 个前端文件。对照 `visual-style-guide.md` 和 `development-strategy.md` 审查后发现：

- `AGENT_NAMES` 常量在 ChatArea.tsx、ConversationItem.tsx 两处重复定义
- AgentAvatar.tsx 中 `AGENT_SHADOW_COLORS` 和 `STATUS_COLORS` 用硬编码色值而非 CSS 变量
- 多个组件手写 `onMouseEnter/Leave` hover 交互（ConversationItem、ConversationList、NewChatDialog），已达 3 次重复
- `api.ts` 的 `createConversation` 访问 `detail.sessions[0]` 无空数组保护
- `Session.agent_name` 类型为必填 `string`，但旧 session 无此字段

## Goals / Non-Goals

**Goals:**
- 消除 `AGENT_NAMES` 重复，集中到 `lib/constants.ts`
- AgentAvatar 硬编码色值改用 CSS 变量
- 提取 hover 交互为 `useHoverStyle` hook
- 修复 `createConversation` 空数组保护和类型安全

**Non-Goals:**
- 不重构组件结构或路由
- 不改变功能行为
- 不涉及后端代码

## Decisions

### D1: 常量文件 `lib/constants.ts`

```ts
export const AGENT_NAMES: Record<AgentType, string> = { ... }
export const AGENT_DESCRIPTIONS: Record<AgentType, string> = { ... }
```

ChatArea、ConversationItem、NewChatDialog 统一从此导入。`AGENT_DESCRIPTIONS` 也从 NewChatDialog 移入。

### D2: AgentAvatar 使用 CSS 变量

shadow 色值从 `getComputedStyle` 读取 CSS 变量，或直接在 JS 中引用 `var(--agent-xxx)`。由于 `boxShadow` 需要拼接字符串，保持 JS 常量但改为引用 CSS 变量值：通过 `getComputedStyle(document.documentElement).getPropertyValue('--agent-claude')` 读取。

实际上更简单的做法：shadow 色值常量已经和 CSS 变量值相同（都是 #DA7756 等），问题是硬编码而非通过变量引用。在 style 属性中无法直接用 `var()` 做 boxShadow 颜色拼接，所以维持 JS 常量但确保与 CSS 变量同步。添加注释标注来源。

STATUS_COLORS 同理 — 目前无 CSS 变量对应，统一不改。

### D3: `useHoverStyle` hook

```ts
export function useHoverStyle(hoverBg = 'var(--bg-hover)', normalBg = 'transparent') {
  return {
    onMouseEnter: (e: React.MouseEvent<HTMLElement>) => { e.currentTarget.style.backgroundColor = hoverBg },
    onMouseLeave: (e: React.MouseEvent<HTMLElement>) => { e.currentTarget.style.backgroundColor = normalBg },
  }
}
```

替代 ConversationItem、ConversationList、NewChatDialog 中的手写 hover 逻辑。

### D4: `createConversation` 空数组保护

```ts
const session = detail.sessions[0]
if (!session) throw new Error('Backend failed to create session')
```

`Session.agent_name` 类型改为 `agent_name?: string`，映射处用 `s.agent_name ?? ''`。

## Risks / Trade-offs

- [useHoverStyle 适用范围] → 仅用于简单的 backgroundColor hover，不含 active/focus 状态，不覆盖所有场景
- [常量文件增长] → 目前只有 agent 相关常量，量小可控，不会变成 utils 垃圾桶
