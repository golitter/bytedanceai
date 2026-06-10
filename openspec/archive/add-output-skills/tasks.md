## 1. AgentEnd: SkillRule + 输出技能 Prompt

- [x] 1.1 新建 `agentend/src/rules/skill_rule.py`，实现 SkillRule（priority=1, check=True, enforce 返回 system_prompt_append）
- [x] 1.2 在 `agentend/src/rules/builtin.py` 的 `__all__` 中导出 SkillRule
- [x] 1.3 在 `agentend/src/app/dependencies.py` 的 `create_rule_engine` 中注册 SkillRule
- [ ] 1.4 验证 SkillRule prompt 注入：启动服务，发请求确认 system_prompt_append 包含 5 种卡片说明

## 2. AgentEnd: Workspace 文件/Diff API

- [x] 2.1 在 `agentend/src/api/v1/workspace.py` 新增 `GET /files/{path}` 端点（FileResponse 流式返回，防路径穿越）
- [x] 2.2 新增 `PUT /files/{path}` 端点（写入 worktree 文件）
- [x] 2.3 新增 `GET /diff` 端点（执行 `git diff HEAD`，返回 unified diff 文本）
- [x] 2.4 新增 `POST /commit` 端点（git add -A + commit）
- [x] 2.5 新增 `POST /revert` 端点（git checkout HEAD -- .）
- [x] 2.6 在 `agentend/src/api/v1/router.py` 注册新路由
- [ ] 2.7 验证 API：curl 测试文件读取、diff、commit、revert

## 3. AgentEnd: 预览 HTTP 服务

- [x] 3.1 新建 `agentend/src/preview/server.py`，实现本地 HTTP 静态文件服务（基于 worktree 目录）
- [x] 3.2 端口动态分配（或基于 workspace_id hash 偏移），避免冲突
- [x] 3.3 在 workspace 生命周期中管理预览服务启停
- [ ] 3.4 验证：启动预览服务，浏览器访问 `http://localhost:{port}/index.html` 能正确渲染

## 4. Go Backend: Workspace API Proxy

- [x] 4.1 在 `backend/internal/handler/workspace.go` 新增文件读取 proxy（`GET /api/workspace/{id}/files/{path}`）
- [x] 4.2 新增 diff proxy（`GET /api/workspace/{id}/diff`）
- [x] 4.3 新增 commit proxy（`POST /api/workspace/{id}/commit`）
- [x] 4.4 新增 revert proxy（`POST /api/workspace/{id}/revert`）
- [x] 4.5 在 `cmd/server/main.go` 注册新路由
- [x] 4.6 所有 proxy 使用 `io.Copy` 流式传输，不使用 `io.ReadAll`

## 5. Frontend: Block 类型 + Block Reducer

- [x] 5.1 新建 `frontend/src/lib/block-types.ts`，定义 MessageBlock discriminated union（TextBlock, HtmlBlock, ImageBlock, AttachmentBlock, DiffBlock, PreviewBlock）
- [x] 5.2 新建 `frontend/src/lib/block-reducer.ts`，实现 `reduceEventToBlocks(fullText) → MessageBlock[]`
- [x] 5.3 解析逻辑：remark parse → 遍历 AST → 检测 code(language="aka_yhy") → 解析 type 字段 → 生成对应 Block
- [x] 5.4 单元测试：覆盖 5 种卡片类型解析 + 未知类型降级 + 纯文本

## 6. Frontend: 5 种卡片组件

- [x] 6.1 新建 `frontend/src/components/cards/HtmlCard.tsx`（iframe sandbox + srcdoc 渲染）
- [x] 6.2 新建 `frontend/src/components/cards/ImageCard.tsx`（img src 指向 Go proxy file API）
- [x] 6.3 新建 `frontend/src/components/cards/AttachmentCard.tsx`（文件图标 + 下载链接）
- [x] 6.4 新建 `frontend/src/components/cards/DiffCard.tsx`（react-diff-viewer 渲染 + 编辑/接受/拒绝按钮）
- [x] 6.5 新建 `frontend/src/components/cards/PreviewCard.tsx`（iframe src 指向 preview URL）
- [x] 6.6 安装依赖：`pnpm add react-diff-viewer-continued`

## 7. Frontend: Chat Store + MessageBubble 升级

- [x] 7.1 修改 `frontend/src/stores/chat.ts`：ChatMessage 结构增加 `blocks: MessageBlock[]` 字段
- [x] 7.2 修改 `streamDone` action：调用 `reduceEventToBlocks(streamingContent)` 生成 blocks 数组
- [x] 7.3 修改 `frontend/src/components/chat/MessageBubble.tsx`：遍历 `message.blocks`，按 type 渲染对应卡片组件
- [x] 7.4 降级兼容：如果 blocks 为空，fallback 到原有 `content` 字段渲染

## 8. 集成验证

- [ ] 8.1 端到端测试：发消息让 Agent 输出 html-render 块，确认前端渲染 HTML 卡片
- [ ] 8.2 端到端测试：Agent 修改文件后输出 diff 块，确认 DiffCard 显示 diff
- [ ] 8.3 端到端测试：DiffCard 编辑文件 → 保存 → 确认 worktree 文件更新
- [ ] 8.4 端到端测试：Agent 启动预览 + 输出 preview 块，确认 PreviewCard iframe 渲染
- [ ] 8.5 降级测试：Agent 不输出 aka_yhy 块时，普通文本渲染不受影响
