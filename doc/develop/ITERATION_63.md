# Iter 63: Export/Import 功能

## 目标
实现 Pattern 的导入导出功能，支持 JSON 格式，便于用户备份和迁移数据。

## 任务清单
- [x] 设计 Export/Import 接口
- [x] 实现 Pattern 导出为 JSON
- [x] 实现从 JSON 导入 Pattern
- [x] 添加 CLI 命令 (export/import)
- [x] 编译测试

## 验收标准
- [x] 导出命令正常工作，生成有效 JSON 文件
- [x] 导入命令正常工作，正确恢复 Pattern 数据
- [x] 编译通过
- [x] 单元测试通过

## 实现细节

### 新增文件
- `pkg/export/exporter.go`: Export/Import 功能实现
  - `Exporter` 结构体：导出 patterns 到 JSON 文件
  - `Importer` 结构体：从 JSON 文件导入 patterns
  - `ExportData` / `ImportData` 结构体：JSON 数据格式

### 修改文件
- `cmd/cli/main.go`: 添加 export/import 命令
  - `export --output <path>`: 导出 patterns
  - `export --output <path> --project <name>`: 按项目过滤导出
  - `import --input <path>`: 导入 patterns
  - `import --input <path> --force`: 覆盖已存在的 pattern

### CLI 命令
```bash
# 导出所有 patterns
otr export --output patterns.json

# 按项目导出
otr export --output patterns.json --project myproject

# 导入 patterns
otr import --input patterns.json

# 覆盖已存在的 patterns
otr import --input patterns.json --force
```

### 测试结果
- 导出 4 个 patterns 到 JSON 文件成功
- 导入新 pattern 成功
- 检测已存在的 patterns 并正确跳过
- 所有单元测试通过
