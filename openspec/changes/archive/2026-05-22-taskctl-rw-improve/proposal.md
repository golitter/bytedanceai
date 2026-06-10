## Why

taskctl 的 memory 读写存在实用性缺陷：读操作只能批量读取整个目录、无法查看单个文件；写操作的 content 参数只取单个词导致多词内容丢失；写入无原子性保证；错误退出码不一致。这些问题影响了 agent 协作场景下的实际使用。

## What Changes

- `common-memory` 和 `sub-memory` 命令增加可选的文件名参数，支持读取指定文件
- `write-sub-memory` 的 content 参数改为 join 所有剩余参数，并支持从 stdin 读取内容（stdin 优先）
- 写入操作使用临时文件 + rename 实现原子写入
- 统一所有命令的错误退出码（读失败也返回非零）
- `parsePath` 去掉未使用的 `taskID` 返回值

## Capabilities

### New Capabilities

_(无新增能力)_

### Modified Capabilities

- `taskctl-cli`: 增加单文件读取能力、stdin 写入支持、原子写入、错误处理一致性

## Impact

- 影响文件：`agentend/src/skills/builtin/taskctl/main.go`（主要改动）
- 命令行接口向后兼容：现有无参数调用行为不变
- 不影响调用方，除非调用方依赖错误退出码为 0（现在读失败会返回非零）
