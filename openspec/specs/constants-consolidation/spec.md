## ADDED Requirements

### Requirement: API_BASE 统一到 lib/constants.ts
系统 SHALL 在 `lib/constants.ts` 中导出唯一的 `API_BASE` 常量。`api.ts`、`DiffCard.tsx`、`AttachmentCard.tsx`、`ImageCard.tsx` 中的本地 `const API_BASE = '/api'` 声明 MUST 删除，改为从 `lib/constants.ts` import。

#### Scenario: API_BASE 只有一处定义
- **WHEN** 搜索前端代码中所有 `API_BASE` 的定义
- **THEN** 只有 `lib/constants.ts` 中有一处 `export const API_BASE = '/api'`，其余文件均为 import 引用

#### Scenario: DiffCard 的 API 请求路径不变
- **WHEN** DiffCard 加载快照数据
- **THEN** 请求路径仍为 `/api/diff-snapshots/${snapshotId}`，行为不变

### Requirement: Agent 颜色映射统一到 lib/constants.ts
系统 SHALL 在 `lib/constants.ts` 中导出唯一的 `AGENT_COLORS: Record<AgentType, string>`，值为 CSS 变量引用。`AgentAvatar.tsx` 中的 `AGENT_COLORS` 和 `AGENT_SHADOW_COLORS`、`MessageBubble.tsx` 中的 `AGENT_STRIP_COLORS` MUST 删除，改为从 `lib/constants.ts` import `AGENT_COLORS`。

#### Scenario: Agent 头像颜色不变
- **WHEN** 渲染 Claude Code Agent 的头像
- **THEN** 头像背景色仍为 `var(--agent-claude)`，boxShadow 仍为 `0 0 8px var(--agent-claude)`

#### Scenario: 消息色条颜色不变
- **WHEN** 渲染 Agent 消息的左侧色条
- **THEN** 色条颜色仍为对应 Agent 的标识色 CSS 变量（如 `var(--agent-claude)`）

### Requirement: getFileName 工具函数统一到 lib/utils.ts
系统 SHALL 在 `lib/utils.ts` 中导出 `getFileName(path: string): string` 函数。`diff-parser.ts` 中已有的 `getFileName` export MUST 删除，改为从 `lib/utils.ts` import。`AttachmentCard.tsx` 和 `ImageCard.tsx` 中的内联 `path.split('/').pop()` MUST 替换为 `getFileName(path)` 调用。

#### Scenario: 文件名提取结果不变
- **WHEN** 对 `'src/components/Button.tsx'` 调用 `getFileName`
- **THEN** 返回 `'Button.tsx'`

#### Scenario: 无路径分隔符时的降级
- **WHEN** 对 `'README.md'` 调用 `getFileName`
- **THEN** 返回 `'README.md'`

### Requirement: 魔法数字提取为命名常量
组件中的魔法数字 MUST 提取为文件顶部的命名常量。具体包括：MessageInput 中的 `200`（最大高度）、`3000`（提示持续时间）、`48`（最小高度）；MessageList 中的 `60`（滚动判定阈值）；AgentAvatar 中的 `'2s'` 和 `'1.5s'`（状态动画时长）。

#### Scenario: MessageInput 高度常量化
- **WHEN** 阅读 MessageInput 源码
- **THEN** 找到 `MAX_INPUT_HEIGHT`、`HINT_DISPLAY_DURATION`、`MIN_INPUT_HEIGHT` 命名常量，不再有裸数字 `200`、`3000`、`48`

#### Scenario: AgentAvatar 动画时长常量化
- **WHEN** 阅读 AgentAvatar 源码
- **THEN** 找到 `STATUS_READY_DURATION`、`STATUS_RUNNING_DURATION` 命名常量，不再有裸字符串 `'2s'`、`'1.5s'`

### Requirement: CHANGE_TYPE_MAP 统一
DiffFileTabs.tsx 中的 `CHANGE_TYPE_MAP` 和 `ChangeTypeBadge` 组件 SHALL 导出供 DiffCard 复用。DiffCard.tsx 文件信息条中的变更类型 badge MUST 使用 `ChangeTypeBadge` 组件，不再内联重复映射。

#### Scenario: DiffCard 复用 ChangeTypeBadge
- **WHEN** DiffCard 渲染文件信息条中的变更类型标记
- **THEN** 使用从 DiffFileTabs import 的 `ChangeTypeBadge` 组件，颜色与 Tab 栏一致
