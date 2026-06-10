## Context

Agent 回复只有纯文本，前端 streamingContent 做字符串拼接，无法区分内容类型。AgentEnd 已有 RuleEngine + SkillProvisioning 机制（注入 system_prompt + 复制工具文件到 workspace）。前端已有 MarkdownRenderer（rehype-highlight 代码高亮）。Go 后端已是 SSE + API proxy 薄壳。

需要新增一套"输出技能"：Agent 通过特殊标记 ` ```aka_yhy ` 声明输出类型，前端解析后渲染富媒体卡片。

## Goals / Non-Goals

**Goals:**
- Agent 通过 RuleEngine 注入的 prompt 知道可以输出 5 种卡片类型
- 前端 Block Reducer 从 markdown AST 中解析 aka_yhy 块，生成 Block 数组
- 5 种卡片组件各自渲染：HTML iframe、图片、附件下载、可编辑 diff、网页预览
- AgentEnd 提供 workspace 文件读取、diff、预览 HTTP 服务
- 所有 API 经过 Go 后端 proxy 透传

**Non-Goals:**
- 不做 EventEnvelope 升级（id/task_id/agent 字段），留给 Phase 4.5
- 不做 EventLog / Replay / 断线重连
- 不做 Timeline 视图（多 Agent 协作时间线）
- 不做 Artifact DAG / versioning / lineage
- 不做 Go 后端 artifact 元数据持久化
- 不修改任何 Adapter 层代码

## Decisions

### Decision 1: 统一用 `aka_yhy` 标记而非多种 code block 语言

**选择**: 所有卡片类型统一用 ` ```aka_yhy ` 包裹，内部 `type` 字段区分类型
**替代方案**: 每种卡片用不同 code block 语言（` ```html-render `、` ```diff ` 等）
**理由**: 单一入口简化前端解析逻辑，一个正则/AST 过滤器搞定。多语言标记需要前端对每种语言单独处理，且可能与真实语言标记冲突。

### Decision 2: RuleEngine 注入 prompt 而非 workspace 扫描

**选择**: 通过 SkillRule 注入 system_prompt 告诉 Agent 可用的输出格式
**替代方案**: AgentEnd 在 stream 结束后扫描 workspace 文件变更
**理由**: Agent 主动标记意图明确（知道什么该渲染 vs 什么只是代码），支持内联内容（HTML 直接在块内），Adapter 层零改动。workspace 扫描作为后续增强。

### Decision 3: 展示型 vs 触发型的区分

**选择**: html-render/image/attachment 是展示型（数据在块内或路径），diff/preview 是触发型（纯信号，触发后端/前端操作）
**理由**: diff 和 preview 需要系统资源（git diff、HTTP 服务），不应每次都触发。由 Agent 主动声明何时需要。

### Decision 4: diff 展示全部 workspace 变更

**选择**: `type: diff` 无额外字段，AgentEnd 执行 `git diff HEAD` 展示所有变更
**替代方案**: Agent 指定 files 列表，只展示部分文件 diff
**理由**: 简化 Agent 使用成本，一个信号搞定。Agent 不需要知道自己改了哪些文件。

### Decision 5: 预览服务由 AgentEnd 启动

**选择**: AgentEnd 在 workspace 启动本地 HTTP 服务，返回 URL 给前端
**替代方案**: 前端启动预览服务（Service Worker / Blob URL）
**理由**: AgentEnd 已有 workspace 路径，直接用 Python `http.server` 或类似方案最简单。前端只需拿到 URL 渲染 iframe。

### Decision 6: 所有 API 经过 Go 后端 proxy

**选择**: 前端不直连 AgentEnd，所有 workspace file/diff/preview API 经 Go 后端透传
**替代方案**: 前端直连 AgentEnd
**理由**: 保持现有架构一致性（Go 是统一入口），前端只需知道 Go 后端地址。Go 不做任何数据存储，纯 proxy。

### Decision 7: Block Reducer 是纯函数，不依赖 Zustand

**选择**: `reduceEventToBlocks(fullText: string) → MessageBlock[]` 纯函数
**替代方案**: 逻辑直接写在 chat store 的 streamDone 里
**理由**: 纯函数可测试、可复用（后面 replay/测试都能用），和 store 解耦。

### Decision 8: DiffCard 编辑用 CodeMirror，不是行内 diff 编辑

**选择**: diff 视图只读展示，点"编辑"打开完整文件编辑器（CodeMirror）
**替代方案**: 在 diff 视图上直接编辑（行内编辑）
**理由**: 行内 diff 编辑实现复杂（需要 patch 生成、冲突处理），MVP 先做文件级编辑，后续再做行内编辑。

### Decision 8.5: Skill 可执行工具用 Go 编写，多文件结构

**选择**: diff 工具和 preview 服务作为 Go 项目，放在 `agentend/src/skills/builtin/` 目录下，各自独立目录，多文件组织。
**替代方案**: 在 AgentEnd Python 里直接实现 diff/preview 逻辑
**理由**: Go 编译为单二进制，无运行时依赖，性能好。preview 服务涉及 HTTP 静态文件服务，Go 天然适合。diff 操作涉及 git 命令执行，Go 也更合适。多文件编程，按职责拆分（main/cmd/internal），复杂度可控。

目录结构：
```
agentend/src/skills/builtin/
├── taskctl                    # 已有 skill
├── skill-diff/                # Go 项目：diff 工具
│   ├── main.go
│   ├── diff.go                # diff 逻辑
│   ├── edit.go                # 文件编辑逻辑
│   └── commit.go              # commit/revert 逻辑
└── skill-preview/             # Go 项目：preview 服务
    ├── main.go
    ├── server.go              # HTTP 服务
    ├── port.go                # 端口分配
    └── lifecycle.go           # 启停管理
```

### Decision 8.6: Agent 调用 skill 由 CLI 自动实现

**选择**: 不实现任何 skill 调用机制。Agent CLIs（Claude Code、OpenCode、Codex）自动发现和调用 workspace 里的 skill 工具。
**替代方案**: AgentEnd 实现显式的 skill 调用/注册系统
**理由**: Claude Code、OpenCode、Codex 各自有工具发现机制，通过 SkillProvisioning 将 Go 二进制复制到 workspace 后，CLI 自然识别。我们只负责：放文件 + 注入 prompt 告诉 Agent 有这些工具。调用流程完全由 CLI 自己处理。
**替代方案**: 在 diff 视图上直接编辑（行内编辑）
**理由**: 行内 diff 编辑实现复杂（需要 patch 生成、冲突处理），MVP 先做文件级编辑，后续再做行内编辑。

### Decision 9: Block marker 配置化，前端运行时同步

**选择**: block marker 放在 `config.yaml` 的 `skills.block_marker` 字段，前端通过 `GET /v1/config/block-marker` API 运行时获取
**替代方案**: 前端硬编码同名常量（简单但改一处容易忘改另一处）
**理由**: 改 marker 只需改 config.yaml，AgentEnd 和前端自动同步。前端启动时 fetch 一次缓存即可。

### Decision 10: config.yaml skills 配置结构

```yaml
skills:
  builtin_dir: "src/skills/builtin"
  block_marker: "aka_yhy"           # Agent 输出卡片的统一标记符号
  manifest: ...                      # 已有
```

AgentEnd SkillRule 从 `settings.skills.block_marker` 读取，注入到 prompt 中。AgentEnd 暴露 `GET /v1/config/block-marker` 返回该值。前端启动时调用 Go proxy `GET /api/config/block-marker`，缓存为 Block Reducer 的解析标记。

## Risks / Trade-offs

- **[Agent 可能不遵守 aka_yhy 格式]** → 前端对无法解析的 aka_yhy 块降级为普通代码块展示，不影响使用
- **[HTML iframe 安全风险]** → 使用 sandbox 属性限制权限，不允许访问 cookie/storage/top-navigation
- **[预览服务端口冲突]** → AgentEnd 动态分配端口或使用固定偏移量（基于 workspace_id hash）
- **[大文件 OOM]** → AgentEnd 文件 API 使用 FileResponse 流式传输，Go 用 io.Copy 透传，不读全量到内存
- **[Block Reducer 和现有 chat store 兼容]** → streamDone 时同时生成 blocks 和 content（content 降级用），确保向后兼容
- **[Block marker 前端缓存失效]** → 前端启动时获取一次，刷新页面重新获取，不做热更新（marker 变更极少）
