# Open-Think-Reflex

> **AI输入加速器** - 通过反射机制加速人与AI的交互

## 核心概念

> **AI = 大脑**（复杂思考）  
> **反射 = 快捷指令**（快速响应）

这个项目实现了一套反射形成系统，模仿人类形成反射的过程：
- **大脑** 处理复杂推理
- **反射** 处理高频操作
- 重复使用 → 形成反射
- 长期不用 → 自动衰减

## 应用场景

```
用户输入"用户API"
    │
    ▼
┌─────────────────────────────────────┐
│ 思维链树展开                        │
│ ├── 分页 pagination ───── 85%       │
│ ├── 响应 response ───── 70%        │
│ └── 错误 error ───────── 60%       │
└─────────────────────────────────────┘
    │
    ▼
用户选择"分页" → 按空格 → AI生成完整代码
```

## 项目结构

```
open-think-reflex/
├── cmd/                    # 入口程序
│   └── otr/                # CLI 主程序
├── internal/               # 内部包 (不可被外部导入)
│   ├── ai/                 # AI 集成模块
│   │   ├── provider/       # AI Provider (Claude)
│   │   ├── prompt/         # Prompt 构建器
│   │   └── response/       # 响应解析器
│   ├── cli/                # CLI 模块
│   │   ├── commands/       # 命令实现
│   │   ├── ui/              # 终端 UI
│   │   └── output/         # 输出格式化
│   ├── config/             # 配置管理
│   ├── core/               # 核心业务
│   │   ├── matcher/        # 匹配引擎
│   │   ├── pattern/        # Pattern 管理
│   │   └── reflex/         # 反射机制
│   ├── data/              # 数据层
│   │   ├── cache/          # LRU 缓存
│   │   └── sqlite/         # SQLite 存储
│   └── ui/                 # 通用 UI
├── pkg/                    # 公共包 (可被外部导入)
│   ├── ai/                 # AI 接口定义
│   ├── contracts/         # 契约/接口
│   ├── errors/             # 错误码定义
│   ├── export/             # 导出功能
│   ├── models/             # 数据模型
│   └── utils/              # 工具函数
├── doc/                    # 项目文档
│   ├── architecture/       # 架构设计文档
│   │   ├── ARCHITECTURE.md
│   │   ├── ERROR_CODES.md
│   │   └── SCHEMA.md
│   ├── develop/            # 开发文档
│   │   ├── PROGRESS.md     # 进度追踪
│   │   ├── ITERATION_PLAN.md
│   │   └── ITERATION_*.md  # 各迭代记录
│   ├── prototypes/         # 原型设计
│   │   └── PROTOTYPE.md
│   ├── requirements/      # 需求文档
│   │   ├── REFLEX_MODEL.md
│   │   └── PROTOCOL.md
│   └── scenarios/          # 场景设计
│       └── README.md
├── otr                     # 编译后的二进制文件
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
└── README.md               # 本文件
```

## 目录说明

| 目录 | 用途 |
|------|------|
| `cmd/` | 程序入口点，包含 main 函数 |
| `internal/` | 私有代码，仅限本项目使用 |
| `internal/ai/` | Claude AI 集成 (Provider/Prompt/Response) |
| `internal/cli/` | 命令行界面实现 |
| `internal/core/` | 核心业务逻辑 (匹配/Pattern/反射) |
| `internal/data/` | 数据持久化 (SQLite/缓存) |
| `pkg/` | 公共库，可被外部项目导入 |
| `pkg/models/` | 核心数据模型定义 |
| `pkg/errors/` | 错误码系统 |
| `doc/` | 项目文档集合 |
| `doc/architecture/` | 架构、错误码、数据schema |
| `doc/develop/` | 开发进度、迭代记录 |
| `doc/requirements/` | 产品需求规格说明 |
| `doc/prototypes/` | 原型设计稿 |
| `doc/scenarios/` | 使用场景分析 |

## v1.0 目标

| 功能 | 状态 | 说明 |
|------|------|------|
| 反射形成 | ✅ | 匹配→强化→阈值→激活 |
| 快捷键触发 | ✅ | ↑↓选择 → 空格生成 |
| 本地存储 | ✅ | SQLite |
| 衰减机制 | ✅ | 自动衰减长期未用的反射 |
| AI 集成 | ✅ | Claude API 集成 |
| 流式输出 | ✅ | 打字机效果 |

## 快速开始

### Linux/macOS

```bash
# 克隆项目
git clone https://github.com/ArmyClaw/open-think-reflex.git
cd open-think-reflex

# 查看需求文档
cat doc/requirements/REFLEX_MODEL.md

# 查看架构设计
cat doc/architecture/ARCHITECTURE.md

# 查看当前进度
cat doc/develop/PROGRESS.md

# 编译
go build -o otr ./cmd/otr

# 运行
./otr --help
```

### Windows

#### 前置要求

1. **安装 Go 语言环境**
   - 下载 Go for Windows: https://go.dev/dl/
   - 选择 Windows MSI installer 或 ZIP 文件
   - 安装后打开 PowerShell 验证: `go version`

2. **安装 Git (可选)**
   - 下载 Git for Windows: https://git-scm.com/download/win
   - 或使用 Windows Terminal 自带的 Git

#### 编译步骤

```powershell
# 克隆项目
git clone https://github.com/ArmyClaw/open-think-reflex.git
cd open-think-reflex

# 编译 (PowerShell)
go build -o otr.exe .\cmd\otr

# 运行
.\otr.exe --help
```

#### 使用 WSL (推荐)

如果需要更好的体验，推荐使用 WSL:

```bash
# 在 WSL 中
git clone https://github.com/ArmyClaw/open-think-reflex.git
cd open-think-reflex
go build -o otr ./cmd/otr
./otr --help
```

### 配置文件

项目使用 YAML 格式配置文件，默认为 `~/.otr/config.yaml`

复制示例配置快速开始:

```bash
# Linux/macOS
mkdir -p ~/.otr
cp config.example.yaml ~/.otr/config.yaml

# Windows (PowerShell)
New-Item -ItemType Directory -Path $env:USERPROFILE\.otr -Force
Copy-Item config.example.yaml $env:USERPROFILE\.otr\config.yaml
```

配置项说明:

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `app.data_dir` | 数据存储目录 | `~/.otr` |
| `app.log_level` | 日志级别 (debug/info/warn/error) | `info` |
| `ai.provider` | AI 提供商 (anthropic/openai/local) | `anthropic` |
| `ai.default_model` | 默认模型 | `claude-sonnet-4-20250514` |
| `storage.type` | 存储类型 | `sqlite` |
| `storage.path` | 数据库路径 | `~/.otr/data.db` |

环境变量覆盖: 配置项可通过 `OTR_` 前缀的环境变量覆盖，如 `OTR_ANTHROPIC_API_KEY`

## 开发指南

### 迭代开发

项目采用 15 分钟小迭代开发模式：
- 每个迭代有明确的验收标准
- 详细记录在 `doc/develop/ITERATION_PLAN.md`
- 进度追踪在 `doc/develop/PROGRESS.md`

### 测试覆盖目标

| 模块 | 目标覆盖率 |
|------|-----------|
| matcher | > 80% |
| models | > 70% |
| prompt | > 80% |
| response | > 80% |
| cache | > 80% |
| config | > 70% |

## v2.0 规划

- **经验导出** - 将反射导出为AgentSkill，AI辅助润色实现经验复用
- 项目空间隔离
- 思绪整理模式
- 导出与同步
- 多人协作

## License

MIT
