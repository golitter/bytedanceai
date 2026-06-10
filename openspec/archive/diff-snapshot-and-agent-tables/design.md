## Context

当前 DiffCard 通过 `sessionId` 从 workspace 的 `git diff` 获取变更内容。用户 commit/revert 后 diff 为空，卡片变白。之前尝试在 Session 表加 `settled_diff`/`diff_status` 字段，但一个 session 有多条消息，每条消息可能有独立的 diff block，session 级字段无法区分。

同时 Session 表承载了 agent 信息（agent_type、agent_name、avatar_url），这些信息属于 agent 身份维度而非会话生命周期，应拆表。

## Goals / Non-Goals

**Goals:**
- 每个 diff block 有独立的 snapshot 记录，状态跨页面刷新持久化
- 同一 session 同时只有一个 pending diff，新 diff 出现时自动取消前一个
- agent 信息拆到独立表，Session 表职责更清晰
- 前端 API 响应结构不变（session_agents 通过 JOIN 透明返回）

**Non-Goals:**
- 不支持一条消息内多个 diff block（当前每条消息最多一个 diff）
- 不修改 Message 表的 agent_type/agent_name 字段（创建时快照，保持不变）
- 不做历史数据迁移（AutoMigrate 即可，存量 session 无 snapshot 无影响）

## Decisions

### 1. snapshotId 由 agent 生成 vs 前端生成

**选择：agent 生成**

diff block 标记在 agent 输出时就已确定。让 agent 生成 UUID 写入 block 内容，前端直接解析使用，无需额外请求来分配 ID。

替代方案：前端在渲染 DiffCard 时生成 snapshotId。缺点是 block 数据不自包含，需要前端维护映射关系。

### 2. snapshot 行创建时机：前端首次渲染时

DiffCard 首次渲染时调 `PUT /api/diff-snapshots/:snapshotId` 创建 pending 行（附带 workspace diff 全文）。后端此时自动取消同 session 其他 pending 行。

替代方案：后端在 stream writer 解析消息内容时创建。缺点是解析逻辑侵入 stream 写入流程，且需要从消息文本中提取 snapshotId。

### 3. diff-snapshot API 路径挂载位置

**选择：独立路径 `/api/diff-snapshots/:snapshotId`**

不挂在 `/api/session/:sessionId/` 下，因为 snapshotId 本身已全局唯一，无需 session 前缀。

### 4. session_agents 拆表策略

**选择：1:1 独立表，后端 JOIN 返回**

session_agents 以 session_id 为 UNIQUE KEY，与 Session 一一对应。后端查询 Session 时 LEFT JOIN session_agents，响应 JSON 结构不变。前端零改动。

替代方案：多对多 agent_profiles 表，多 session 共享同一 profile。当前无此需求，过度设计。

### 5. Session 表字段清理

移除：`agent_type`、`agent_name`、`avatar_url`、`settled_diff`、`diff_status`（settled_diff/diff_status 是本次之前临时加的，用新表替代）。

GORM AutoMigrate 只加列不删列。删除列需要手动 migration 或接受冗余字段暂存。考虑到不影响功能，**本轮不删列**，仅不再使用这些字段，后续清理。

## Risks / Trade-offs

- **[Agent render skill 是 Go 编译产物]** → card_diff.go 需要 `import "github.com/google/uuid"` 或用简单时间戳方案。确认 go.mod 已有依赖或选择轻量替代。
- **[AutoMigrate 不删列]** → Session 表的 agent_type/agent_name/avatar_url/settled_diff/diff_status 列物理存在但不被代码引用。不影响功能，后续手动清理。
- **[Pending snapshot 创建依赖前端]** → 如果用户从未渲染某条消息的 DiffCard（如历史消息未滚动到），该 snapshot 不会被创建。此时卡片显示 fallback 为空或跳过。
- **[CANCELLED 的 diff 内容]** → 在 pending 创建时已保存 diff 全文，cancelled 后仍可展示原始内容（只读 + 灰化）。
