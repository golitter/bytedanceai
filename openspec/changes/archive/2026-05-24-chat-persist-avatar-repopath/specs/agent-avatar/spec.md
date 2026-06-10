## ADDED Requirements

### Requirement: Agent 头像上传
系统 SHALL 提供头像上传接口 `POST /api/agents/avatar`，接收 multipart/form-data 图片文件（仅限 jpg/png/gif/webp，最大 2MB），存储到服务器文件系统并返回可访问的 URL。

#### Scenario: 成功上传头像
- **WHEN** 上传一张 1MB 的 jpg 图片
- **THEN** 返回 HTTP 200，body 为 `{"code": 0, "data": {"avatar_url": "/uploads/avatars/xxx.jpg"}}`

#### Scenario: 文件格式不支持
- **WHEN** 上传 .svg 或 .bmp 文件
- **THEN** 返回 HTTP 400，提示不支持的文件格式

#### Scenario: 文件过大
- **WHEN** 上传超过 2MB 的图片
- **THEN** 返回 HTTP 400，提示文件过大

### Requirement: Agent 头像存储关联
系统 SHALL 在 session 或 agent 配置上存储 avatar_url 字段，前端根据此字段展示自定义头像。

#### Scenario: 有自定义头像
- **WHEN** agent 配置了 avatar_url
- **THEN** 前端使用该 URL 展示头像图片

#### Scenario: 无自定义头像
- **WHEN** agent 未配置 avatar_url
- **THEN** 前端使用 DiceBear 根据 agent name 生成头像，URL 为 `https://api.dicebear.com/7.x/initials/svg?seed=<name>`

### Requirement: Agent 名称和头像编辑
系统 SHALL 支持在创建后修改 agent 的名称和头像。前端通过点击头像或名称触发编辑弹窗。

#### Scenario: 修改名称
- **WHEN** 用户在编辑弹窗中修改 agent 名称并保存
- **THEN** 后端更新该 agent 的名称，前端即时反映变更

#### Scenario: 修改头像
- **WHEN** 用户在编辑弹窗中上传新头像并保存
- **THEN** 后端更新 avatar_url，旧头像文件被清理，前端即时展示新头像

#### Scenario: 同时修改名称和头像
- **WHEN** 用户同时修改名称和头像
- **THEN** 两者同时更新，头像变更立即生效

### Requirement: DiceBear 头像生成
前端 SHALL 在 agent 无自定义头像时，使用 DiceBear API 根据 agent name 生成确定性 initials 风格头像。

#### Scenario: 根据 name 生成头像
- **WHEN** agent 名为 "claude-code" 且无自定义头像
- **THEN** 前端渲染 DiceBear 生成的 initials 头像，同名始终生成相同头像

#### Scenario: 名称变更后头像同步
- **WHEN** agent 无自定义头像时名称被修改
- **THEN** DiceBear 生成的头像随之改变
