## MODIFIED Requirements

### Requirement: 全量色彩合规检查
系统 SHALL 对所有 `frontend/src/` 下的 TSX/CSS 文件执行色彩合规检查，确保不存在以下违规：
- 渐变（`linear-gradient`、`radial-gradient`）
- `backdrop-blur`（毛玻璃效果）
- 纯白文字（`#FFFFFF`、`text-white` 用于正文）
- 纯黑文字（`#000000`）
- Agent 标识色用于 UI 框架装饰（非身份标识场景）
- **`text-secondary` 用于渲染次要文字**（SHALL 使用 `text-text-secondary`）
- **inline `style={{ color: 'var(--text-secondary)' }}`（SHALL 使用 `text-text-secondary` 类）**

#### Scenario: 检测到 text-secondary 违规
- **WHEN** 任何 TSX 文件使用 `text-secondary` 类渲染次要文字
- **THEN** 标记为 Warning 违规，提示应使用 `text-text-secondary`

#### Scenario: 检测到 inline style 文字色违规
- **WHEN** 任何 TSX 文件使用 `style={{ color: 'var(--text-secondary)' }}` 而非 Tailwind 类
- **THEN** 标记为 Info 违规，提示改用 `text-text-secondary` 类

#### Scenario: 检测到渐变违规
- **WHEN** 任何 TSX/CSS 文件包含 `linear-gradient` 或 `radial-gradient` 属性
- **THEN** 标记为 Critical 违规，输出文件路径和行号

#### Scenario: 检测到 backdrop-blur 违规
- **WHEN** 任何 TSX/CSS 文件包含 `backdrop-blur` 且不限于固定/粘性元素
- **THEN** 标记为 Critical 违规，输出文件路径和行号

#### Scenario: 色彩全部合规
- **WHEN** 所有文件的色彩使用符合 visual-style-guide.md 定义
- **THEN** 审计报告显示色彩合规通过

### Requirement: 圆角系统合规检查
系统 SHALL 验证所有组件的 border-radius 值符合定义：按钮 6px、输入框 8px、卡片 10px、面板 12px、Agent 头像 8px、Badge 9999px，且不存在 16px 以上大圆角或 rounded-full 用于按钮/卡片。**`rounded-2xl`（16px）用于非面板/非 Badge 元素 SHALL 标记为违规。**

#### Scenario: 检测到 rounded-2xl 违规
- **WHEN** 任何非 Badge/胶囊元素使用 `rounded-2xl`（16px）
- **THEN** 标记为 Warning 违规，提示降级为 `rounded-xl`（12px）或更小

#### Scenario: 检测到超大圆角
- **WHEN** 任何组件使用大于 12px 的圆角（Badge/胶囊除外）
- **THEN** 标记为 Warning 违规

#### Scenario: 圆角全部合规
- **WHEN** 所有圆角值匹配 visual-style-guide.md 定义
- **THEN** 审计报告显示圆角合规通过

### Requirement: 边框和阴影合规检查
系统 SHALL 验证边框使用 `rgba(255,255,255,0.06)` 或对应 CSS 变量，阴影仅用于弹出菜单/下拉框场景。**非 shadcn/ui 组件（`components/ui/` 之外的组件）使用 `shadow-lg` 或 `shadow-md` SHALL 标记为违规。**

#### Scenario: 检测到非菜单阴影
- **WHEN** 非 shadcn/ui 组件使用 `shadow-lg` 或 `shadow-md`
- **THEN** 标记为 Warning 违规，提示使用背景色阶替代

#### Scenario: shadcn/ui 组件阴影豁免
- **WHEN** `components/ui/` 目录下的组件使用 `shadow-lg` 或 `shadow-md`
- **THEN** 豁免检查，不标记为违规

#### Scenario: 检测到非标准边框
- **WHEN** 组件使用实色边框（如 `border-[#333]`）而非透明度边框
- **THEN** 标记为 Warning 违规

## ADDED Requirements

### Requirement: 硬编码项目元信息 SHALL 检测并标记
系统 SHALL 检测组件中内联的项目元信息（GitHub URL、项目名称、项目描述等），标记为违规并提示集中到 `lib/constants.ts` 的 `PROJECT_META`。

#### Scenario: 检测到内联 GitHub URL
- **WHEN** 任何 TSX 文件包含 `href="https://github.com/..."` 硬编码字符串
- **THEN** 标记为 Warning 违规，提示使用 `PROJECT_META.GITHUB_URL`

#### Scenario: 检测到内联项目描述
- **WHEN** 任何 TSX 文件包含项目描述字符串（如"Multi-Agent Chat Platform"、"多 Agent 协作聊天平台"）作为 JSX 内容
- **THEN** 标记为 Warning 违规，提示使用 `PROJECT_META.DESCRIPTION`
