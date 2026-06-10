## 1. 更新 README.md

- [x] 1.1 更新「当前状态」代码块：AgentEnd ~85%、Backend ~80%、Frontend ~70%
- [x] 1.2 更新「阶段总览」表格：新增 Phase 5-7 行，标记 Phase 1-3 完成、Phase 4 待执行
- [x] 1.3 更新「Phase 依赖关系」图：Phase 4 → Phase 5 → Phase 6/7（6 和 7 可并行）
- [x] 1.4 更新「总计」预估：约 9-11 个工作日（新增 Phase 5-7）

## 2. 新增 phase5-orchestrator.md

- [x] 2.1 创建文件头部：目标、预估（3-4 天）、前置条件（Phase 4 完成）
- [x] 2.2 编写「交付标准」：curl 验证 Orchestrator 流 + 前端群聊 UI 截图标准
- [x] 2.3 编写 AgentEnd 实现步骤（Step 1-7）：config + LlmConfig → models → prompts → graph → adapter → schemas 扩展 → 注册
- [x] 2.4 编写 Go Backend Scheduler：读取 config.yaml → 按 task 数组顺序调度 Agent → SSE 透传
- [x] 2.5 编写前端群聊 UI：Agent 标签 + 颜色区分 + 协作流展示
- [x] 2.6 编写文件清单和验证流程（基于 orchestrator-plan-phase.md 的 curl 验证）

## 3. 新增 phase6-preview-deploy.md

- [x] 3.1 创建文件头部：DRAFT 标记 + 目标 + 前置条件
- [x] 3.2 列出功能范围：Artifact Manager、ArtifactCard、iframe 预览、部署卡片、Go 代理层
- [x] 3.3 列出「不做」清单：Diff 视图、版本历史、容器化部署、桌面端/移动端
- [x] 3.4 标注预估 TBD + 对照任务要求 #4 #5 的映射

## 4. 新增 phase7-demo-deliver.md

- [x] 4.1 创建文件头部：目标、预估（2 天）、前置条件（Phase 5 完成，Phase 6 可选）
- [x] 4.2 编写交付物清单（5 类）：产品设计文档、技术文档、可运行 Demo、AI 协作记录、Demo 视频
- [x] 4.3 编写 UI 打磨项：响应式适配、主题一致性、错误处理、Demo 场景脚本
- [x] 4.4 编写稳定性保障项：Agent 断连处理、超时重试、边界情况
