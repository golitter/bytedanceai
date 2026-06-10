## 1. 文字颜色修正（消除纯白文字）

- [x] 1.1 AgentAvatar.tsx — `text-white` 改为 `text-foreground`（头像首字母）
- [x] 1.2 GroupAvatar.tsx — `text-white` 改为 `text-foreground`（多 Agent 头像首字母）

## 2. 阴影移除（AgentAvatar 发光效果）

- [x] 2.1 AgentAvatar.tsx — 移除 `boxShadow: 0 0 8px ${shadowColor}`，改用同色低透明度背景圈表达层次

## 3. Agent 标识色范围限制

- [x] 3.1 PlanCard.tsx — agent badge 停止使用 AGENT_COLORS 直接着色，改用 `text-secondary` + `bg-accent` 中性配色

## 4. 动画缓动修正

- [x] 4.1 AgentAvatar.tsx — 状态灯 ready 脉冲从 `ease-in-out` 改为 `ease-out`

## 5. 字号对齐

- [x] 5.1 GroupAvatar.tsx — 首字母从 `text-[9px]` 改为 `text-[11px]`
- [x] 5.2 GroupAvatar.tsx — 成员数角标从 `text-[8px]` 改为 `text-[11px]`

## 6. 内联样式 Tailwind 化

- [x] 6.1 AgentAvatar.tsx — 将 width/height/borderRadius 内联样式迁移至 Tailwind 尺寸类（保留动态颜色内联）
- [x] 6.2 验证所有修改的组件 TypeScript 编译通过且无 ESLint 错误
