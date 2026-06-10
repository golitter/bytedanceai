## ADDED Requirements

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

### Requirement: 字体系统合规检查
系统 SHALL 验证所有文字渲染使用 Geist Sans（UI 文字）或 Geist Mono（代码文字），字号不超出定义的 6 级字号体系（11/12/13/14/20px 对应的场景）。

#### Scenario: 检测到非标准字体
- **WHEN** 任何组件的 CSS 使用非 Geist 字体族
- **THEN** 标记为 Warning 违规

#### Scenario: 字体全部合规
- **WHEN** 所有字体使用匹配 visual-style-guide.md 定义
- **THEN** 审计报告显示字体合规通过

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

### Requirement: 动效合规检查
系统 SHALL 验证所有动画效果符合规范：时长 120-300ms ease-out，仅使用 transform 和 opacity，不存在弹簧物理、渐变流光、粒子效果、大幅位移等禁止动效。

#### Scenario: 检测到禁止动效
- **WHEN** 组件使用 bounce/overshoot 弹簧动画或大幅位移 keyframe
- **THEN** 标记为 Critical 违规

#### Scenario: 动效全部合规
- **WHEN** 所有动画符合 visual-style-guide.md 动效规范
- **THEN** 审计报告显示动效合规通过

### Requirement: Agent 身份系统合规检查
系统 SHALL 验证 Agent 头像、状态灯、消息色条严格遵循 visual-style-guide.md 定义的颜色、尺寸和动效规范。

#### Scenario: Agent 标识色正确应用
- **WHEN** 检查 Agent 消息组件
- **THEN** Claude 使用 `#DA7756`、OpenCode 使用 `#10B981`、Orchestrator 使用 `#EAB308`，仅用于头像和左侧色条

#### Scenario: Agent 状态灯规范
- **WHEN** 检查 Agent 状态指示器
- **THEN** 就绪为绿色脉冲、执行中为黄色旋转、离线为灰色无动效、出错为红色无动效

### Requirement: 硬编码项目元信息 SHALL 检测并标记
系统 SHALL 检测组件中内联的项目元信息（GitHub URL、项目名称、项目描述等），标记为违规并提示集中到 `lib/constants.ts` 的 `PROJECT_META`。

#### Scenario: 检测到内联 GitHub URL
- **WHEN** 任何 TSX 文件包含 `href="https://github.com/..."` 硬编码字符串
- **THEN** 标记为 Warning 违规，提示使用 `PROJECT_META.GITHUB_URL`

#### Scenario: 检测到内联项目描述
- **WHEN** 任何 TSX 文件包含项目描述字符串（如"Multi-Agent Chat Platform"、"多 Agent 协作聊天平台"）作为 JSX 内容
- **THEN** 标记为 Warning 违规，提示使用 `PROJECT_META.DESCRIPTION`
