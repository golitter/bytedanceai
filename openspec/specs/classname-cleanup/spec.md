## ADDED Requirements

### Requirement: className 拼接统一使用 cn()
组件中条件拼接 className SHALL 使用 `cn()` 函数（来自 `@/lib/utils`），不得使用模板字面量（`` className={`...${cond}`} ``）或字符串加法。`cn()` 内部基于 clsx + tailwind-merge，可正确处理 Tailwind 类冲突。

#### Scenario: AgentAvatar 条件 className
- **WHEN** AgentAvatar 组件根据 size 拼接 className 时
- **THEN** SHALL 使用 `cn('flex items-center ...', sizeClass)` 而非模板字面量

#### Scenario: AgentHoverCard badge className
- **WHEN** AgentHoverCard 的 badge 根据类型拼接颜色类时
- **THEN** SHALL 使用 `cn('inline-flex items-center ...', badge.cls)` 而非模板字面量

#### Scenario: RuntimeStatus 状态 className
- **WHEN** RuntimeStatus 根据状态拼接背景色和文字色时
- **THEN** SHALL 使用 `cn()` 拼接 `config.bg`、`config.color` 等条件类

#### Scenario: DiffCard 条件透明度
- **WHEN** DiffCard 根据 isSettled 状态添加透明度时
- **THEN** SHALL 使用 `cn('max-h-96 overflow-auto text-xs', isSettled && 'opacity-60')` 而非模板字面量

#### Scenario: PlanCard 状态颜色
- **WHEN** PlanCard 根据任务状态拼接颜色类时
- **THEN** SHALL 使用 `cn('text-xs', statusColor[task.status] ?? 'text-muted-foreground')` 而非模板字面量

#### Scenario: Admin 页面加载图标条件动画
- **WHEN** WorkspacePage / StatisticsPage / SessionCleanupPage / AgentOverviewPage / ServiceHealthPage / DashboardPage 的加载图标根据 loading 状态添加 `animate-spin` 时
- **THEN** SHALL 使用 `cn('h-3.5 w-3.5', loading && 'animate-spin')` 而非模板字面量

#### Scenario: ServiceHealthPage 状态指示器
- **WHEN** ServiceHealthPage 的服务状态指示器根据 Running 状态添加 `animate-pulse` 时
- **THEN** SHALL 使用 `cn('h-2.5 w-2.5 rounded-full', svc.status === 'Running' && 'animate-pulse')` 而非模板字面量

#### Scenario: DiffHeader badge className
- **WHEN** DiffHeader 的快照状态 badge 拼接颜色类时
- **THEN** SHALL 使用 `cn('mr-2 inline-flex ...', BADGE_CONFIG[snapshotStatus].className)` 而非模板字面量

#### Scenario: DiffFileTabs tab className
- **WHEN** DiffFileTabs 的 tab 标签拼接样式时
- **THEN** SHALL 使用 `cn('shrink-0 rounded ...', config.className)` 而非模板字面量

#### Scenario: AgentMeta 条件字体
- **WHEN** AgentMeta 根据 mono prop 切换字体样式时
- **THEN** SHALL 使用 `cn('text-[13px] break-all', mono && 'font-mono text-xs')` 而非模板字面量
