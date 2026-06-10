## 1. Contracts 契约更新

- [x] 1.1 在 `contracts/schemas/session-state.yaml` 新增 `inactive` 状态值及转换规则
- [x] 1.2 运行 `make generate` 生成三端类型文件

## 2. Backend API

- [x] 2.1 新增 `PATCH /api/sessions/:sessionId` handler，接受 `{ "status": "inactive" }` 请求体
- [x] 2.2 添加参数校验：status 值必须为 "inactive"（当前只支持这一种 PATCH 操作）
- [x] 2.3 添加 404 处理（session 不存在时返回 404）
- [x] 2.4 在 router 中注册新路由

## 3. Agentend 懒清理机制

- [x] 3.1 新增 DB 查询模块：连接 DB 并实现只读查询方法（查询 inactive session、查询每个 task 的 session 状态分布）
- [x] 3.2 新增 `inactive_cleanup` 后台任务：每 2h 轮询 DB，清理 inactive session 的 worktree，清理全 inactive 的 task
- [x] 3.3 在 `config.yaml` 新增 `cleanup_interval` 配置项，默认 7200 秒
- [x] 3.4 修改 `lifespan` startup：去掉 TTL cleanup 启动，替换为 inactive-cleanup 任务启动
- [x] 3.5 修改 `lifespan` shutdown：去掉全量 cleanup 和 session destroy，只停轮询任务
- [x] 3.6 删除 TTL 相关代码（`start_ttl_cleanup`、`stop_ttl_cleanup`、TTL 配置读取）

## 4. Frontend 停用操作

- [x] 4.1 session 列表添加"停用"按钮
- [x] 4.2 点击后调用 `PATCH /api/sessions/:sessionId` 并刷新列表
- [x] 4.3 已 inactive 的 session 显示停用状态（灰色/禁用样式）
