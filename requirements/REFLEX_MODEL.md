# Open-Think-Reflex 需求文档

> **版本**: v1.4  
> **状态**: 草稿  
> **核心概念**: AI反射形成与衰减系统 - 三层交互模式

---

## 1. 核心概念

### 1.1 大脑与反射的类比

```
人类神经系统              AI系统
─────────────────────────────────────────────
大脑（复杂思考）       AI大模型（通用智能）
  ↓                      ↓
脊髓反射（自动化）       代码/模式（快速响应）
```

**核心观点**：
- 大脑 = 大模型（慢速思考）
- 反射 = 代码/模式（快速执行）
- 反射形成需要重复强化
- 反射衰减因为不使用

---

## 2. 三层交互模式

### 2.1 界面布局

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  上层：思维链树展示层（从左到右展开）                                │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐            │
│  │  起点   │ ───► │ 分支A │ ───► │ 分支A1 │ ───► │ 分支A1a│    │
│  │ 用户输入 │      │ ✓高分 │      │ ✓高分 │      │ ✓高分 │    │
│  └─────────┘      └─────────┘      └─────────┘      └─────────┘    │
│                                                                      │
│  视觉反馈：                                                      │
│  ├── 已选分支：高亮显示（绿色/亮色）                             │
│  ├── 未选分支：灰色/淡化                                          │
│  └── 选中态：边框闪烁                                            │
│                                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  中层：当前输出结果层                                            │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │                                                      │     │
│  │              AI输出的结果内容...                           │     │
│  │                                                      │     │
│  └─────────────────────────────────────────────────────────────┘     │
│                                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  下层：用户输入层                                                │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │ > 用户输入命令或问题...                          [Enter] │     │
│  └─────────────────────────────────────────────────────────────┘     │
│                                                                      │
│  操作提示：[Tab]切换层 [↑↓]选择分支 [→]确认滚动 [空格]AI生成      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 交互流程

```
用户输入（下层）
    │
    ▼
匹配反射
    │
    ├── 有匹配 ──► 展示思维链树（上层）
    └── 无匹配 ──► 提示"无匹配，按[空格]AI生成"

选择与生成
    │
    ├── 选中分支 ──► 按[空格]AI生成该分支内容
    └── 无分支可选 ──► 按[空格]AI生成新内容

确认滚动
    │
    └── 按[→]确认当前方向，展开下一级

输出执行
    │
    └── 输出到AI或执行结果
```

### 2.3 各层功能

| 层级 | 功能 | 交互 |
|------|------|------|
| **上层** | 思维链树展示 | 从左到右展开，点击选择方向 |
| **中层** | 当前输出结果 | 显示AI响应或中间结果 |
| **下层** | 用户输入 | 命令行输入，[Enter]确认 |

---

## 3. 反射生命周期模型

### 3.1 七阶段模型

```
建立印象 ──► 强化 ──► 达成阈值 ──► 初步反射 ──► 加强
                │                                    │
                ▼                                    ▼
         【反射形成区】                          深度反射（永久）
                │                                    │
                ▼                                    ▼
         不用 → 降级 → 丧失                      持续强化
```

### 3.2 各阶段特征

| 阶段 | 强度 | 感知 | 可逆性 |
|------|------|------|--------|
| 建立印象 | 10% | 强 | 易 |
| 强化 | 30% | 中 | 可逆 |
| 达成阈值 | 50% | 弱 | 尚可 |
| 初步反射 | 70% | 无 | 可逆 |
| 加强 | 85% | 无 | 难 |
| 深度反射 | 100% | 无 | 极难 |
| 衰减 | ↓ | 弱 | 可逆 |

### 3.3 衰减公式

```
Strength(t) = Strength_initial × e^(-λ × t)

其中：
- λ = 衰减常数（因反射类型而异）
- t = 距离上次强化的时间
```

---

## 4. 功能需求

### 4.1 核心功能

| 功能 | 描述 | 优先级 |
|------|------|--------|
| 三层界面显示 | 上中下三层布局 | P0 |
| 思维链树展开 | 从左到右展示反射路径 | P0 |
| 方向选择 | 上下键选择分支 | P0 |
| 确认滚动 | 右键确认选择方向 | P0 |
| AI触发生成 | 无匹配时按空格生成 | P0 |

### 4.2 反射管理

| 功能 | 描述 | 优先级 |
|------|------|--------|
| 建立反射 | 从对话中提取模式 | P0 |
| 强化反射 | 多次使用增强强度 | P1 |
| 衰减机制 | 不使用自动衰减 | P2 |
| 反射存储 | 本地JSON/SQLite | P1 |

---

## 5. 交互设计

### 5.1 快捷键

| 按键 | 功能 | 适用层级 |
|------|------|----------|
| ↑/↓ | 选择分支 | 上层 |
| → | 确认滚动 | 上层 |
| ← | 返回上层 | 上层 |
| Tab | 切换层级 | 全局 |
| Enter | 确认输入 | 下层 |
| Space | 触发AI生成 | 下层/无匹配时 |
| Esc | 取消/返回 | 全局 |
| q | 退出 | 全局 |
| h | 帮助 | 全局 |

### 5.2 视觉反馈

| 状态 | 视觉表现 |
|------|----------|
| 活跃分支 | 高亮显示 |
| 选择中 | 边框闪烁 |
| 已确认 | 绿色✓ |
| 衰减中 | 黄色⚠ |
| 丧失 | 灰色✗ |

---

## 6. 数据结构

### 6.1 反射结构

```typescript
interface Pattern {
  id: string;                    // 唯一标识
  trigger: string;                 // 触发词
  response: string;               // 响应路径
  strength: number;              // 强度 0-100
  threshold: number;             // 激活阈值
  connections: string[];         // 关联反射
  metadata: {
    created: number;
    updated: number;
    reinforcementCount: number;
    decayCount: number;
  };
}
```

### 6.2 思维链树

```typescript
interface ThoughtChain {
  root: string;                  // 用户输入（根节点）
  branches: ChainBranch[];       // 分支列表
  currentDepth: number;          // 当前深度
  maxDepth: number;            // 最大深度
}

interface ChainBranch {
  id: string;
  pattern: Pattern;             // 关联反射
  next: ChainBranch[];        // 子分支
  selected: boolean;           // 是否被选择
  confirmed: boolean;          // 是否已确认滚动
}
```

---

## 7. 存储设计

### 7.1 SQLite存储（推荐）

SQLite是当前推荐的存储方式，后续可迁移到高速缓存。

```sql
-- 反射表
CREATE TABLE patterns (
  id TEXT PRIMARY KEY,
  trigger TEXT NOT NULL,
  response TEXT NOT NULL,
  strength REAL NOT NULL DEFAULT 0,
  threshold REAL NOT NULL DEFAULT 50,
  decay_rate REAL NOT NULL DEFAULT 0.01,
  decay_enabled INTEGER NOT NULL DEFAULT 1,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  reinforcement_count INTEGER NOT NULL DEFAULT 0
);

-- 连接表
CREATE TABLE connections (
  id TEXT PRIMARY KEY,
  source_pattern_id TEXT NOT NULL,
  target_pattern_id TEXT NOT NULL,
  connection_type TEXT NOT NULL,
  strength REAL NOT NULL DEFAULT 50
);

-- 索引
CREATE INDEX idx_patterns_trigger ON patterns(trigger);
CREATE INDEX idx_patterns_strength ON patterns(strength);
CREATE INDEX idx_connections_source ON connections(source_pattern_id);
```

### 7.2 存储迁移

| 存储方式 | 说明 | 适用场景 |
|----------|------|----------|
| SQLite | 单文件，好备份 | 当前推荐 |
| Redis | 高速缓存 | 生产环境 |
| PostgreSQL | 关系型数据库 | 大规模 |

---

## 8. 应用场景

### 8.1 代码编写助手

```
用户输入: "写个用户API"

上层层（思维链）:
└── 用户API
    ├── 分页 pagination ────── 85%
    │   └── pageSize, pageNum
    ├── 响应 response ────── 70%
    │   └── code, message, data
    └── 错误 error ───────── 60%
        └── errorCode, errorMsg

用户选择"分页"方向
    │
    ▼
向右滚动展开:
└── 用户API ──► 分页 pagination ──► pageSize/pageNum
```

### 8.2 问题解决助手

```
用户输入: "数据库连接失败"

上层层（思维链）:
└── 数据库连接失败
    ├── 网络问题 ────── 75%
    │   └── ping测试, 防火墙检查
    ├── 配置错误 ────── 60%
    │   └── 连接字符串, 认证信息
    └── 服务不可用 ──── 45%
        └── 端口检查, 服务状态
```

---

## 9. 衰减机制

### 9.1 衰减规则

| 反射类型 | 衰减常数(λ) | 默认周期 | 可调参数 |
|----------|---------------|----------|----------|
| 短期反射 | 0.1/天 | 7天衰减50% | 可调 |
| 中期反射 | 0.01/天 | 70天衰减50% | 可调 |
| 长期反射 | 0.001/天 | 700天衰减50% | 可调 |
| 永久反射 | 0 | 不衰减 | 不适用 |

### 9.2 衰减曲线

```
强度
100% │ 深度反射 ────────────────┐
     │                             │
 85% │ 加强 ────────┐              │
     │              │ 衰减曲线    │
 70% │ 初步反射 ──┤              │ 曲线
     │              │              │
 50% │ 阈值 ──────┤              │
     │              │              │
 30% │ 强化 ──────┤              │
     │              │              │
 10% │ 建立 ──────┘              │
     │                             │
  0% └──────────────────┴───────────→ 时间
         0     7天   14天   30天
```

---

## 10. AI集成设计

### 10.1 AI提供商抽象层

不绑定特定AI提供商，支持多种AI服务。

```typescript
interface AIProvider {
  name: string;
  version: string;
  
  complete(prompt: string): Promise<string>;
  stream(prompt: string): AsyncIterator<string>;
  embed(text: string): Promise<number[]>;
}

interface AIConfig {
  apiKey?: string;
  baseUrl?: string;
  model?: string;
  maxTokens?: number;
  temperature?: number;
}
```

### 10.2 支持的提供商

| 提供商 | 状态 | 说明 |
|--------|------|------|
| Claude | ✅ 当前推荐 | Anthropic API |
| OpenAI | ⏳ 后续 | GPT-4 API |
| 本地模型 | ⏳ 后续 | Ollama/LM Studio |

---

## 11. 导出与同步

### 11.1 MCP Protocol导出

支持导出为MCP Protocol格式，可被其他AI工具使用。

```typescript
interface MCPExport {
  name: string;
  version: string;
  patterns: MCPPattern[];
}

interface MCPPattern {
  name: string;
  trigger: string[];
  response: {
    type: 'prompt' | 'action';
    content: string;
  };
  strength: number;
}
```

### 11.2 导出命令

```bash
# 导出为MCP格式
otr export --format mcp --output patterns_mcp.json

# 导出为SQLite
otr export --format sqlite --output patterns.db
```

### 11.3 同步机制

| 方式 | 说明 | 状态 |
|------|------|------|
| 文件导出 | JSON/SQLite文件 | ✅ |
| 云同步 | Redis/PostgreSQL | ⏳ |
| P2P同步 | 局域网同步 | ⏳ |

---

## 12. 离线模式

### 12.1 离线功能

| 功能 | 说明 | 状态 |
|------|------|------|
| 反射查询 | 本地SQLite查询 | ✅ |
| 思维链展示 | 本地渲染 | ✅ |
| 衰减计算 | 本地计算 | ✅ |
| AI生成 | 需要网络 | ❌ |
| 同步 | 需要网络 | ❌ |

### 12.2 离线优先策略

```
1. 先读本地
   └── 查询本地SQLite
       └── 展示思维链

2. 后写本地
   └── 用户操作先写本地
       └── 网络恢复后同步

3. 队列同步
   └── 操作队列持久化
       └── 网络恢复后自动同步
```

---

## 13. 命令行使用

### 13.1 基本命令

```bash
# 启动交互模式
otr

# 直接输入问题
otr "写个用户API"

# 查看帮助
otr --help

# 导出反射
otr export patterns.json

# 导入反射
otr import patterns.json

# 衰减配置
otr decay status
otr decay set user_api_pagination --rate 0.05
```

### 13.2 交互模式操作

```
> otr

┌─────────────────────────────────────────┐
│ 上层：思维链                              │
│ [根] ──► [分支A] ──► [分支A1] ✓          │
│ [根] ──► [分支B] ──► [分支B1] ◉          │
└─────────────────────────────────────────┤
│ 中层：输出结果                           │
│ AI: 这是根据你的选择生成的代码...        │
└─────────────────────────────────────────┤
│ 下层：输入                               │
│ > _ [按↑↓选择,→确认,空格AI生成]          │
└─────────────────────────────────────────┘
```

---

## 14. 未来展望

### 14.1 理想状态

```
用户输入简短关键词
    │
    ▼
思维链树自动展开显示所有可能方向
    │
    ▼
用户选择方向，快速滚动迭代
    │
    ▼
AI生成高质量结果
```

### 14.2 长期目标

1. **零指令执行**
   - 从最小线索学习用户意图
   - 不需要显式指令

2. **主动协助**
   - 预测用户需求
   - 自动触发反射

3. **个性化智能**
   - 每个用户发展独特的反射集
   - AI成为个性化助手

---

**文档版本**: v1.4  
**更新日期**: 2026-02-20  
**项目**: open-think-reflex
