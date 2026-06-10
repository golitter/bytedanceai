## ADDED Requirements

### Requirement: 相对时间格式化
系统 SHALL 提供 `formatRelativeTime(timestamp)` 函数，根据消息时间戳与当前时间的差值，返回对应的相对时间字符串。格式规则：同一天显示 "HH:mm"（如 "14:30"），昨天显示 "昨天 HH:mm"（如 "昨天 22:10"），2-7天前显示 "N天前"（如 "3天前"），今年更早显示 "M月D日 HH:mm"（如 "3月15日 14:30"），跨年显示 "YYYY年M月D日"（如 "2025年3月15日"）。

#### Scenario: 今天的消息
- **WHEN** 消息时间戳对应的日期是今天
- **THEN** 返回 "HH:mm" 格式的时间字符串（如 "14:30"）

#### Scenario: 昨天的消息
- **WHEN** 消息时间戳对应的日期是昨天
- **THEN** 返回 "昨天 HH:mm" 格式的时间字符串（如 "昨天 22:10"）

#### Scenario: 2-7天前的消息
- **WHEN** 消息时间戳对应的日期在 2 到 7 天前之间
- **THEN** 返回 "N天前" 格式的时间字符串（如 "3天前"）

#### Scenario: 今年更早的消息
- **WHEN** 消息时间戳对应的日期在今年但超过 7 天前
- **THEN** 返回 "M月D日 HH:mm" 格式的时间字符串（如 "3月15日 14:30"）

#### Scenario: 跨年的消息
- **WHEN** 消息时间戳对应的日期不在今年
- **THEN** 返回 "YYYY年M月D日" 格式的时间字符串（如 "2025年3月15日"）

### Requirement: 时间分隔线显示判断
系统 SHALL 提供 `shouldShowTimeSeparator(prevTimestamp, currentTimestamp)` 函数，判断两条相邻消息之间是否应显示时间分隔线。触发条件：首条消息（无 prevTimestamp）、两条消息间隔超过 5 分钟、两条消息不在同一天。

#### Scenario: 列表首条消息
- **WHEN** 当前消息是列表中的第一条消息（无前一条）
- **THEN** 应显示时间分隔线

#### Scenario: 消息间隔超过 5 分钟
- **WHEN** 两条相邻消息的时间间隔超过 5 分钟
- **THEN** 应在两条消息之间显示时间分隔线

#### Scenario: 消息间隔在 5 分钟内
- **WHEN** 两条相邻消息的时间间隔不超过 5 分钟且在同一天
- **THEN** 不显示时间分隔线

#### Scenario: 跨天消息
- **WHEN** 两条相邻消息不在同一天
- **THEN** 应显示时间分隔线，即使间隔不超过 5 分钟

### Requirement: TimeDivider 组件
系统 SHALL 提供 `TimeDivider` UI 组件，渲染居中的时间分隔线，显示相对时间文本，两侧带有横线装饰。使用 muted 文字颜色，12px 字号。

#### Scenario: 时间分隔线渲染
- **WHEN** TimeDivider 组件接收到 timestamp 参数
- **THEN** 渲染居中的文本，文本内容为 formatRelativeTime(timestamp) 的返回值，两侧带有分隔横线

#### Scenario: 时间分隔线样式
- **WHEN** TimeDivider 渲染
- **THEN** 使用 muted-foreground 颜色，12px 字号，上下各有 8px 内边距
