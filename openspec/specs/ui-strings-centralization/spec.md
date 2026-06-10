## ADDED Requirements

### Requirement: 创建 UI 字符串常量文件
项目 SHALL 创建 `src/lib/ui-text.ts` 文件，按语义分组集中管理所有中文 UI 字符串。分组 MUST 包含：
- `UI_ACTIONS` — 操作按钮文字（返回、保存、取消、删除、导出等）
- `UI_STATUS` — 状态提示文字（加载中、发送中、正在回复等）
- `UI_MESSAGES` — 消息提示文字（复制成功、发送失败、加载失败等）
- `UI_LABELS` — 区块标题和标签（Agent 信息、路径信息、Git Graph 等）
- `UI_PLACEHOLDERS` — 输入框占位符（输入消息...、搜索技能...等）

每个常量对象 SHALL 使用 `as const` 确保字面量类型安全。

#### Scenario: ui-text.ts 导出类型安全常量
- **WHEN** 组件导入 ui-text.ts 中的常量
- **THEN** TypeScript 编译通过，常量值为字面量类型（非 string）

### Requirement: 组件引用常量替代硬编码字符串
所有散落在组件 TSX 中的中文 UI 字符串 SHALL 替换为 `ui-text.ts` 中对应常量的引用。`lib/constants.ts` 中已有的字符串（如 AGENT_DESCRIPTIONS）保持不变。

#### Scenario: RightSidebar 标签引用常量
- **WHEN** RightSidebar 渲染"Agent 信息"标题
- **THEN** 使用 `UI_LABELS.AGENT_INFO` 而非硬编码 `'Agent 信息'`

#### Scenario: 错误提示引用常量
- **WHEN** use-chat-stream 显示"发送失败"提示
- **THEN** 使用 `UI_MESSAGES.SEND_FAILED` 而非硬编码 `'发送失败'`

#### Scenario: 操作按钮引用常量
- **WHEN** 组件渲染"导出"按钮
- **THEN** 使用 `UI_ACTIONS.EXPORT` 而非硬编码 `'导出'`

### Requirement: 不遗漏任何散落字符串
完成替换后，`grep -rn` 扫描 `frontend/src/components/` 和 `frontend/src/hooks/` 中的中文字符（Unicode CJK 范围），SHALL NOT 发现任何仍硬编码在组件中的 UI 字符串（注释和 console.log 除外）。

#### Scenario: 替换后无遗漏
- **WHEN** 执行 `grep -rn '[\x{4e00}-\x{9fff}]' frontend/src/components/ frontend/src/hooks/`
- **THEN** 仅有注释和 console.log 中包含中文字符，所有 UI 渲染字符串均已替换为常量引用
