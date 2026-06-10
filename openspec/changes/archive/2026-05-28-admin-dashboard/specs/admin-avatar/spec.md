## ADDED Requirements

### Requirement: Avatar modification entry
管理菜单顶部 SHALL 展示当前用户头像，点击后可进入头像修改。

#### Scenario: Display current avatar
- **WHEN** 用户进入管理面板
- **THEN** 管理菜单顶部显示当前用户头像

#### Scenario: Open avatar editor
- **WHEN** 用户点击头像
- **THEN** 弹出头像修改弹窗，显示当前头像和修改选项

### Requirement: Upload custom avatar
用户 SHALL 能上传自定义图片作为头像。

#### Scenario: Upload image
- **WHEN** 用户选择本地图片文件并确认上传
- **THEN** 图片上传到七牛云存储，头像更新为新图片 URL，IconSidebar 中的头像同步更新

#### Scenario: Upload fails
- **WHEN** 上传的文件格式不支持或超过大小限制
- **THEN** 显示错误提示，头像不变

### Requirement: Switch DiceBear avatar
用户 SHALL 能通过选择 DiceBear 风格和 seed 来生成头像。

#### Scenario: Change DiceBear style
- **WHEN** 用户选择不同 DiceBear 风格（如 notionists、avataaars、bottts 等）
- **THEN** 实时预览新风格的头像，确认后更新

#### Scenario: Change DiceBear seed
- **WHEN** 用户修改 seed 文本
- **THEN** 实时预览基于新 seed 生成的头像

### Requirement: Avatar API
后端 SHALL 提供头像保存和获取接口。

#### Scenario: Save avatar URL
- **WHEN** 前端发送 `PUT /api/admin/avatar` 且 body 包含 `{ url: "https://..." }`
- **THEN** 后端保存头像 URL，返回 `{ success: true }`

#### Scenario: Get avatar URL
- **WHEN** 前端请求 `GET /api/admin/avatar`
- **THEN** 返回当前头像 URL
