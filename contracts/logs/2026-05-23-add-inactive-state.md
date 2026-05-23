# 新增 inactive 会话状态

## 变更原因
支持用户手动停用 session，触发 agentend 懒清理机制。

## 变更文件
- `schemas/session-state.yaml`

## 对比结果
- SessionState 枚举新增 `inactive` 值
- SessionStateTransitions 新增 `inactive` 属性（终态，空转换数组）
- required 列表新增 `inactive`

## 跨端影响
- **Frontend**: `generated/session.ts` SessionState 类型新增 `"inactive"`
- **Backend**: `internal/generated/session.go` 新增 Inactive 常量
- **Agentend**: `generated/session.py` SessionState 新增 INACTIVE 值

## 契约变更
- 枚举值: +inactive
- 转换规则: +inactive（终态，不可转换到其他状态）
