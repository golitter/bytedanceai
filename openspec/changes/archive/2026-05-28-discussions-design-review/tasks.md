## 1. 跨文档一致性审查

- [x] 1.1 检查 4 篇文档中事件类型命名是否统一（dot vs underscore）
- [x] 1.2 检查 Workspace 路径描述是否一致（workspaces/ vs worktrees/）
- [x] 1.3 检查 Go Backend 变更范围的描述是否矛盾
- [x] 1.4 检查数据模型字段名是否跨文档统一

## 2. 代码一致性审查

- [x] 2.1 对照 contracts/schemas/ 验证 AgentRequest 新增字段是否遗漏契约更新
- [x] 2.2 对照 agentend 实际代码验证 Orchestrator 当前 5 阶段流程描述
- [x] 2.3 对照 backend StreamEvent 验证新事件类型的扩展方式
- [x] 2.4 对照 workspace/manager.py 验证 worktree 路径生成逻辑

## 3. 设计盲点识别

- [x] 3.1 标注 Conversation Layer 持久化方案缺失
- [x] 3.2 标注 RuntimeEvent 与 StreamEvent 关系未定义
- [x] 3.3 标注 Profile System 存储与 Backend SessionAgent 表的冲突
- [x] 3.4 标注 Regenerate 设计对 Backend Message 表的影响

## 4. 输出问题清单

- [x] 4.1 创建 discussions/design-review-issues.md，按 P0/P1/P2 分级记录所有问题
