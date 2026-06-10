## 1. parsePath 清理

- [x] 1.1 将 `parsePath` 返回值从 `(taskID, sessionID, sharedDir string, err error)` 改为 `(sessionID, sharedDir string, err error)`，移除 taskID
- [x] 1.2 更新 main() 中 parsePath 调用处，去掉 `_ = taskID`

## 2. 单文件读取支持

- [x] 2.1 修改 `cmdCommonMemory` 签名增加可选文件名参数，有参数时调用 `os.ReadFile` 读单文件，无参数保持原有批量读取逻辑；文件不存在时输出 stderr 并 exit(1)
- [x] 2.2 修改 `cmdSubMemory` 签名增加可选文件名参数，逻辑同上
- [x] 2.3 更新 main() switch 中 common-memory / sub-memory 的调用，传入 `os.Args[2:]` 中的可选文件名

## 3. write-sub-memory 改进

- [x] 3.1 添加 stdin 检测函数：使用 `os.Stdin.Stat()` 判断是否为管道输入（非 `ModeCharDevice`），是则 `io.ReadAll(os.Stdin)` 读取内容
- [x] 3.2 修改 `cmdWriteSubMemory`：stdin 有数据用 stdin 内容，否则用 `strings.Join(os.Args[3:], " ")`；两者都为空时报错 exit(1)
- [x] 3.3 添加 `atomicWriteFile` 函数：在同一目录创建临时文件（`os.CreateTemp`），写入内容后 `os.Rename` 到目标路径

## 4. 错误处理统一

- [x] 4.1 将 `cmdLs`、`cmdSummary`、`cmdCommonMemory`、`cmdSubMemory` 中的错误 return 改为 `os.Exit(1)`
- [x] 4.2 确保所有成功路径正常 return（exit 0）

## 5. help 更新

- [x] 5.1 更新 `printHelp` 中的命令说明，体现新增的可选文件名参数和 stdin 支持
