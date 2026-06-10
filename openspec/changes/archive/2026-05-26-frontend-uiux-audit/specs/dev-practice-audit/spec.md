## ADDED Requirements

### Requirement: 三层组件模型合规检查
系统 SHALL 验证所有组件遵循 Page/Smart/Dumb 三层模型：Page 负责路由和数据组合、Smart 负责状态管理和业务编排、Dumb 只接收 props 渲染 UI。不允许在 Dumb 组件中直接发起请求或访问全局状态。

#### Scenario: Dumb 组件违规引入全局状态
- **WHEN** 一个 Dumb 展示组件（如 MessageBubble、AgentAvatar）直接导入 Zustand store 或调用 useQuery
- **THEN** 标记为 Critical 违规，该组件应通过 props 接收数据

#### Scenario: 组件模型全部合规
- **WHEN** 所有组件职责清晰分层
- **THEN** 审计报告显示组件架构合规通过

### Requirement: 状态管理分类合规检查
系统 SHALL 验证状态管理严格遵循三类分类：Server State 用 TanStack Query、Global Client State 用 Zustand、Local Client State 用 useState/useReducer。派生状态不存、URL 作为廉价全局状态优先使用。

#### Scenario: Server State 误入 Zustand
- **WHEN** Zustand store 中存储了来自后端的数据（messages、tasks、sessions 等）
- **THEN** 标记为 Critical 违规，应迁移至 TanStack Query

#### Scenario: 派生状态存入 useState
- **WHEN** 组件使用 useState 存储可通过现有数据计算得出的值
- **THEN** 标记为 Warning 违规，应改为 useMemo

#### Scenario: 状态管理全部合规
- **WHEN** 所有状态分类正确
- **THEN** 审计报告显示状态管理合规通过

### Requirement: 性能模式合规检查
系统 SHALL 验证列表组件使用稳定唯一 ID 作为 key（不使用 index）、超过 100 条的列表考虑虚拟化、搜索输入使用防抖、条件渲染重组件使用 React.lazy。

#### Scenario: 列表使用 index 作为 key
- **WHEN** 任何 `.map()` 渲染使用 `key={i}` 或 `key={index}`
- **THEN** 标记为 Critical 违规，应使用稳定唯一 ID

#### Scenario: 性能模式全部合规
- **WHEN** 所有性能相关模式符合 development-strategy.md
- **THEN** 审计报告显示性能模式合规通过

### Requirement: Hook 拆分合规检查
系统 SHALL 验证自定义 Hook 遵循规范：只抽象涉及状态 + side effect + 生命周期的逻辑，纯计算函数不写成 Hook，单个 Hook 不超过 150-200 行。

#### Scenario: 纯计算函数写成 Hook
- **WHEN** 一个 Hook 只包含纯转换逻辑（如格式化、生成 ID），无 useState/useEffect/useRef
- **THEN** 标记为 Warning 违规，应提取为 src/lib/ 下的普通函数

#### Scenario: Hook 体积超标
- **WHEN** 单个 Hook 文件超过 200 行
- **THEN** 标记为 Warning 违规，建议拆分

### Requirement: TypeScript 类型使用合规检查
系统 SHALL 验证不存在 `any` 类型滥用、魔法字符串直接比较（应抽常量）、条件渲染使用提前返回而非三元嵌套。

#### Scenario: 检测到 any 类型
- **WHEN** 代码中使用 `any` 类型（非第三方库类型声明）
- **THEN** 标记为 Warning 违规

#### Scenario: 检测到三元嵌套
- **WHEN** JSX 中出现嵌套超过 2 层的三元表达式
- **THEN** 标记为 Warning 违规，应使用提前返回
