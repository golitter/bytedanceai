## Context

taskctl 是 agent 协作场景下的共享上下文 CLI 工具，以 symlink 到 `.agent/skills/taskctl/` 目录的可执行文件形式分发，通过解析自身路径确定 task/session 上下文。

当前 memory 命令族（common-memory、sub-memory、write-sub-memory）存在以下问题：
- 读操作只能批量读取整个目录
- 写操作 `os.Args[3]` 只取单个词
- 写入无原子性保证
- 错误退出码不一致

所有改动集中在单个文件 `agentend/src/skills/builtin/taskctl/main.go`。

## Goals / Non-Goals

**Goals:**
- 支持 common-memory / sub-memory 读取单个指定文件
- write-sub-message 支持 stdin 输入 + 参数兜底
- 写入操作原子化
- 统一错误退出码
- 保持命令行接口向后兼容

**Non-Goals:**
- 不改变命令的扁平结构
- 不增加删除 memory 的命令
- 不做并发读写保护（单 agent 单 session 场景）

## Decisions

### D1: 单文件读取通过可选位置参数实现

`common-memory [file]`、`sub-memory [file]` 增加可选的第二个参数。有参数时只读该文件，无参数时读全部（现有行为不变）。

替代方案：用 `--file` flag。弃用原因：与现有 write-sub-memory 的位置参数风格不一致。

### D2: 写入内容优先从 stdin 读取

检测 stdin 是否有数据（`os.Stdin.Stat()` 检查 size 或 pipe 模式），有则读 stdin；否则 `strings.Join(os.Args[3:], " ")` 拼接剩余参数。两者都无内容时报错退出。

### D3: 原子写入用 ioutil.TempFile + rename

写入目标目录下的临时文件（`.tmp` 后缀），写完后 `os.Rename` 到最终路径。rename 在同文件系统上是原子的，保证不会出现半截文件。

### D4: 错误退出码统一

所有命令的读/写失败均 `os.Exit(1)`。help 正常退出返回 0。

## Risks / Trade-offs

- **stdin 检测在终端直接调用时可能误判** → 使用 `os.Stdin.Stat()` 的 `ModeCharDevice` 判断，终端直接调用走参数路径
- **rename 跨文件系统会失败** → 临时文件创建在同一目录下，保证同文件系统
- **向后兼容** → 所有新增参数都是可选的，现有无参数调用行为完全不变
