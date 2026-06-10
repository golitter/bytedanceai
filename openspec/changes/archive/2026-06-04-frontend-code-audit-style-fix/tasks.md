## 1. strokeWidth 统一修正

- [x] 1.1 修正 RightSidebar.tsx 中 PanelRightOpen 等 2 处 `strokeWidth={1.5}` → `strokeWidth={1.25}`
- [x] 1.2 修正 PlanReviewCard.tsx 中 PencilLine `strokeWidth={1.5}` → `strokeWidth={1.25}`
- [x] 1.3 修正 ConversationItem.tsx 中 Pin `strokeWidth={1.5}` → `strokeWidth={1.25}`
- [x] 1.4 修正 ContactsPage.tsx 中 3 处 `strokeWidth={1.5}` → `strokeWidth={1.25}`（Link/ExternalLink/Pin）
- [x] 1.5 修正 AgentProfilePage.tsx 中 Line 元素 `strokeWidth="3"` → 合理数值

## 2. Emoji → Lucide 图标替换（SkillsHubPage）

- [x] 2.1 新增 Lucide 图标 import：Shield、Package、Wrench、AlertTriangle、CheckCircle2、XCircle
- [x] 2.2 SectionLabel 组件：🛡️ → Shield，📦 → Package，替换 emoji prop 为 Lucide 图标渲染
- [x] 2.3 HubSkillCard 内置图标：⚙️ → Wrench（strokeWidth=1.25）
- [x] 2.4 HubSkillCard 外部图标：📦 → Package（strokeWidth=1.25）
- [x] 2.5 UploadDialog 拖拽区：📦 emoji → Package Lucide 图标
- [x] 2.6 UploadDialog 校验成功：✅ → CheckCircle2 图标
- [x] 2.7 UploadDialog 校验失败：❌ → XCircle 图标
- [x] 2.8 DeleteConfirmDialog：⚠️ → AlertTriangle 图标

## 3. text-white → text-primary-foreground

- [x] 3.1 SkillsHubPage 上传按钮：`text-white` → `text-primary-foreground`
- [x] 3.2 SkillsHubPage 确认入库按钮：`text-white` → `text-primary-foreground`
- [x] 3.3 AgentProfilePage 品牌色按钮：`text-white` → `text-primary-foreground`

## 4. text-red-500 → text-destructive 语义替换

- [x] 4.1 HubSkillCard 删除按钮：`text-red-500` → `text-destructive`，`border-red-500/20` → `border-destructive/20`，`bg-red-500/10` → `bg-destructive/10`
- [x] 4.2 DeleteConfirmDialog 删除按钮：`text-red-500` → `text-destructive`，`border-red-500/20` → `border-destructive/20`，`bg-red-500/10` → `bg-destructive/10`
- [x] 4.3 UploadDialog 校验失败区：`text-red-500` → `text-destructive`，`border-red-500/20` → `border-destructive/20`，`bg-red-500/5` → `bg-destructive/5`
- [x] 4.4 DeleteConfirmDialog 危险名称高亮：`text-red-500` → `text-destructive`

## 5. inline style → Tailwind 类

- [x] 5.1 SessionCleanupPage：`style={{ color: 'var(--text-secondary)' }}` → `text-text-secondary` 类（2 处）

## 6. 全量验证

- [x] 6.1 grep 验证 `strokeWidth` 零违规：`grep -rn 'strokeWidth' frontend/src/ --include='*.tsx' | grep -v '1.25'`
- [x] 6.2 grep 验证 `text-white` 零残留：`grep -rn 'text-white' frontend/src/ --include='*.tsx'`
- [x] 6.3 grep 验证 `text-red-500` / `border-red-500` / `bg-red-500` 零残留
- [x] 6.4 grep 验证装饰性 emoji 零残留
- [x] 6.5 页面功能回归确认：SkillsHubPage 上传/删除流程、AgentProfilePage 按钮样式
