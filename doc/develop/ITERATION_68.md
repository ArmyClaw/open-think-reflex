# Iter 68: Space 切换逻辑

**目标**: 实现 Space 切换后 Pattern 隔离
**阶段**: Phase 8 - 项目空间
**时间**: 15 minutes

## 任务

1. 修改 UseSpace 命令 - 保存当前 Space 到配置文件
2. 配置持久化 - 使用 internal/config 保存当前 Space
3. Pattern 隔离 - Space 切换后 Pattern 操作使用对应 Space

## 实现

### 1. 修改 UseSpace 函数

更新 `cmd/cli/main.go` 中的 UseSpace 命令调用，以传入配置并保存 Space：

```go
// 修改 space use 命令
Action: func(c *cli.Context) error {
    return commands.UseSpace(storage, cfg, c.Args().First())
}
```

### 2. 修改 commands.UseSpace

更新 `internal/cli/commands/commands.go` 中的 UseSpace 函数：

- 接收配置对象
- 验证 Space 存在
- 保存当前 Space ID 到配置文件
- 输出成功消息

### 3. 添加 GetCurrentSpace/SetCurrentSpace

在 internal/config 中添加获取和设置当前 Space 的方法。

## 验收标准

- [ ] `otr space use <id>` 保存当前 Space 到配置文件
- [ ] `otr space show` 显示当前 Space
- [ ] Space 切换成功消息
- [ ] 编译通过
- [ ] 测试通过
