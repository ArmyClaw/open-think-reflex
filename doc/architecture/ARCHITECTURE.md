# Open-Think-Reflex Architecture Design

> **Version**: v1.0  
> **Status**: Draft  
> **Scope**: Architecture Design for v1.0 Implementation

---

## Table of Contents

1. [Overview](#1-overview)
2. [Design Principles](#2-design-principles)
3. [Language Selection](#3-language-selection)
4. [System Architecture](#4-system-architecture)
5. [Module Design](#5-module-design)
6. [CLI Interface Layer](#6-cli-interface-layer)
7. [Core Logic Layer](#7-core-logic-layer)
8. [Data Layer](#8-data-layer)
9. [AI Integration Layer](#9-ai-integration-layer)
10. [Module Communication](#10-module-communication)
11. [Extension Points](#11-extension-points)
12. [Security Design](#12-security-design)
13. [Performance Considerations](#13-performance-considerations)
14. [Deployment Architecture](#14-deployment-architecture)

---

## 1. Overview

### 1.1 Purpose

This document defines the system architecture for Open-Think-Reflex, ensuring:
- Separation of CLI UI from core logic
- Modular, testable architecture
- Extensible for future enhancements

---

## 2. Design Principles

### 2.1 Core Design Rules

1. **CLI Is A View, Not The Application** - Core logic works without CLI
2. **Everything Is A Module** - Modules communicate through interfaces
3. **Data Flow Is Unidirectional** - UI → Core → Data → AI
4. **Configuration Over Convention** - Explicit interfaces over implicit assumptions

---

## 3. Language Selection

### 3.1 Recommendation: Go

| Criteria | Score | Justification |
|----------|--------|---------------|
| CLI Ecosystem | 9/10 | Excellent (cobra, urfave/cli) |
| Terminal UI | 9/10 | tview, lipgloss, gocui |
| Single Binary | 10/10 | Native compilation |
| Async | 9/10 | Built-in goroutines |
| Cross-platform | 10/10 | Easy cross-compilation |

### 3.2 Technology Stack

| Layer | Technology | Justification |
|-------|------------|--------------|
| Language | Go 1.21+ | Performance, single binary |
| CLI Framework | urfave/cli | Simple, functional |
| Terminal UI | tview/lipgloss | Rich UI |
| Data Storage | SQLite 3.44+ | Embedded, zero-config |
| AI Integration | Anthropic SDK | Official SDK |

---

## 4. System Architecture

### 4.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│              CLI INTERFACE LAYER                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │ Terminal │  │  Input   │  │   Renderer      │  │
│  │  Handler │  │  Parser │  │   (tview)      │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                 CORE LOGIC LAYER                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │  Match   │  │ Pattern │  │   Reflex        │  │
│  │  Engine │  │ Manager │  │   Lifecycle     │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  DATA LAYER                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │  SQLite  │  │  File   │  │   Exporter      │  │
│  │ Storage  │  │  Cache  │  │                 │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│               AI INTEGRATION LAYER                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │ Provider │  │  Prompt │  │   Response      │  │
│  │  Factory │  │  Builder │  │   Parser       │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 4.2 Layer Responsibilities

| Layer | Responsibility | Dependencies |
|-------|-----------------|--------------|
| **CLI Interface** | User interaction, rendering | Core Logic (via interfaces) |
| **Core Logic** | Business rules, matching, lifecycle | Data Layer, AI Layer |
| **Data Layer** | Persistence, caching, export/import | SQLite, filesystem |
| **AI Integration** | AI provider communication | External AI APIs |

---

## 5. Module Design

### 5.1 Directory Structure

```
cmd/
├── cli/
│   ├── main.go              # Entry point
│   └── root.go             # Root command
│
internal/
├── cli/                     # CLI Interface Layer
│   ├── ui/                 # Terminal UI
│   ├── commands/           # CLI Commands
│   └── output/             # Output Formatters
│
├── core/                   # Core Logic Layer
│   ├── matcher/            # Pattern Matching
│   ├── pattern/            # Pattern Management
│   └── reflex/              # Reflex Lifecycle
│
├── data/                   # Data Layer
│   ├── sqlite/              # SQLite implementation
│   └── cache/               # Caching
│
├── ai/                      # AI Integration Layer
│   ├── provider/            # AI Providers
│   ├── prompt/              # Prompt Engineering
│   └── response/            # Response Handling
│
└── config/                 # Configuration
│
pkg/
├── contracts/               # Interfaces
└── export/                 # Export/Import
```

### 5.2 Module Dependencies

```
cmd/cli/main
    │
    ├── cli/commands (imports)
    │       ├── internal/cli/ui
    │       │       └── internal/core (interfaces)
    │       │
    │       ├── internal/data/sqlite
    │       │       └── internal/core (interfaces)
    │       │
    │       ├── internal/ai/provider
    │       │       └── internal/core (interfaces)
    │       │
    │       └── internal/config
    │
    └── internal/cli/output
            └── internal/core (interfaces)
```

---

## 6. CLI Interface Layer

### 6.1 Design Goals

1. **VIEW ONLY - NO BUSINESS LOGIC** - Pure presentation layer
2. **RENDERER PLUGGABLE** - Multiple output formats
3. **INPUT PARSING ISOLATED** - Reusable parsing logic
4. **THEME SUPPORT** - Light/dark themes

### 6.2 Three-Layer Rendering

```
Layer 1: Thought Chain (Top)
  - Horizontal tree: Root ─► Branch ─► Sub-branch
  - Confidence scores displayed
  - Selection state highlighted

Layer 2: Output (Middle)
  - AI-generated content
  - Markdown-like formatting
  - Scrollable for long content

Layer 3: Input (Bottom)
  - Prompt line
  - User input area
  - Help text
```

### 6.3 Keyboard Handling

```go
// Key bindings match requirements document
const (
    KeyUp      = rune('k')
    KeyDown    = rune('j')
    KeyRight   = rune('l')
    KeyLeft    = rune('h')
    KeyTab     = rune('\t')
    KeySpace   = rune(' ')
    KeyEnter   = rune('\r')
    KeyEscape  = rune('\x1b')
    KeyQuit   = rune('q')
    KeyHelp   = rune('h')
)
```

---

## 7. Core Logic Layer

### 7.1 Design Goals

1. **PURE BUSINESS LOGIC** - No UI dependencies
2. **INTERFACE-BASED DEPENDENCIES** - Testable with mocks
3. **THREAD-SAFE** - Concurrent access support
4. **OBSERVABLE** - Events for state changes

### 7.2 Matcher Engine

```go
import (
    "context"
    "sort"
    "sync"
    
    "github.com/armyclaw/open-think-reflex/internal/core/pattern"
    "github.com/armyclaw/open-think-reflex/pkg/contracts"
)

type Engine struct {
    strategies []MatchingStrategy  // exact, keyword, semantic, fuzzy
    cache      *Cache             // LRU cache
    storage    contracts.Storage
    scorer     Scorer
    config     EngineConfig
}

func (e *Engine) Match(ctx context.Context, query string, opts MatchOptions) ([]MatchResult, error) {
    // 1. Check cache first
    if cached := e.cache.Get(query); cached != nil {
        return cached, nil
    }
    
    // 2. Execute matching strategies in parallel
    var mu sync.Mutex
    results := make([]MatchResult, 0)
    
    var wg sync.WaitGroup
    for _, strategy := range e.strategies {
        wg.Add(1)
        go func(s MatchingStrategy) {
            defer wg.Done()
            matches, _ := s.Match(ctx, query, opts)
            
            mu.Lock()
            defer mu.Unlock()
            for _, m := range matches {
                // Calculate confidence score
                m.Confidence = e.scorer.Calculate(m, query)
                results = append(results, m)
            }
        }(strategy)
    }
    
    wg.Wait()
    
    // 3. Filter by threshold and sort by confidence
    filtered := make([]MatchResult, 0)
    for _, r := range results {
        if r.Confidence >= opts.Threshold {
            filtered = append(filtered, r)
        }
    }
    
    sort.Slice(filtered, func(i, j int) bool {
        return filtered[i].Confidence > filtered[j].Confidence
    })
    
    // 4. Limit results
    if len(filtered) > opts.Limit {
        filtered = filtered[:opts.Limit]
    }
    
    // 5. Cache results
    e.cache.Set(query, filtered)
    
    return filtered, nil
}
```

### 7.3 Matching Strategies

| Strategy | Priority | Use Case |
|----------|-----------|----------|
| Exact | P0 | Keyword exactly matches trigger |
| Keyword | P0 | Keywords extracted and matched |
| Semantic | P1 | Vector similarity (optional) |
| Fuzzy | P2 | Fuzzy matching (optional) |

---

## 8. Data Layer

### 8.1 Storage Interface

```go
type Storage interface {
    SavePattern(ctx context.Context, p *Pattern) error
    GetPattern(ctx context.Context, id string) (*Pattern, error)
    ListPatterns(ctx context.Context, opts ListOptions) ([]*Pattern, error)
    DeletePattern(ctx context.Context, id string) error
}
```

### 8.2 SQLite Schema

```sql
CREATE TABLE patterns (
    id TEXT PRIMARY KEY,
    trigger TEXT NOT NULL,
    response TEXT NOT NULL,
    strength REAL NOT NULL DEFAULT 0,
    threshold REAL NOT NULL DEFAULT 50,
    decay_rate REAL NOT NULL DEFAULT 0.01,
    decay_enabled INTEGER NOT NULL DEFAULT 1,
    connections TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    reinforcement_count INTEGER NOT NULL DEFAULT 0,
    decay_count INTEGER NOT NULL DEFAULT 0,
    tags TEXT,
    project TEXT,
    user_id TEXT
);

CREATE INDEX idx_patterns_trigger ON patterns(trigger);
CREATE INDEX idx_patterns_strength ON patterns(strength);
CREATE INDEX idx_patterns_project ON patterns(project);
```

---

## 9. AI Integration Layer

### 9.1 Provider Interface

```go
type AIProvider interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    Stream(ctx context.Context, req *GenerateRequest, output io.Writer) error
    HealthCheck(ctx context.Context) error
    Name() string
}
```

### 9.2 Supported Providers

| Provider | Status | Notes |
|----------|---------|-------|
| Claude | ✅ Default | Anthropic official SDK |
| OpenAI | ⏳ | GPT-4 support |
| Local | ⏳ | Ollama/LM Studio |

---

## 10. Module Communication

### 10.1 Communication Patterns

```
┌─────────────────────────────────────────────────────────┐
│              Communication Patterns                │
├─────────────────────────────────────────────────────────┤
│                                                  │
│  1. INTERFACE INJECTION                          │
│     Dependencies injected via constructor        │
│     Example: func NewEngine(s Storage)           │
│                                                  │
│  2. EVENT BUS (Observer Pattern)                │
│     State change events published                │
│     Example: PatternReinforced, PatternDecayed  │
│                                                  │
│  3. CONTEXT PROPAGATION                        │
│     Request-scoped values                       │
│     Cancellation support                         │
│                                                  │
│  4. DIRECT METHOD CALLS                        │
│     Synchronous operations                      │
│     Within layer communication                  │
│                                                  │
└─────────────────────────────────────────────────────────┘
```

### 10.2 Event System

```go
type EventType string

const (
    EventPatternCreated   EventType = "pattern.created"
    EventPatternMatched   EventType = "pattern.matched"
    EventPatternReached  EventType = "pattern.threshold_reached"
    EventPatternDecayed  EventType = "pattern.decayed"
    EventPatternDeleted  EventType = "pattern.deleted"
)

type Event struct {
    Type      EventType
    Timestamp time.Time
    Payload   interface{}
}
```

---

## 11. Extension Points

### 11.1 Extension Interfaces

| Extension Point | Interface | Description |
|----------------|-----------|-------------|
| **Renderer** | contracts.Renderer | Custom output formats |
| **Matcher** | contracts.Matcher | Custom matching algorithms |
| **Storage** | contracts.Storage | Custom persistence |
| **AI Provider** | contracts.AIProvider | Custom AI services |
| **Exporter** | contracts.Exporter | Custom export formats |

### 11.2 Plugin System

```go
// Plugin interface for extensions
type Plugin interface {
    Name() string
    Version() string
    Initialize(ctx context.Context, registry *PluginRegistry) error
    Shutdown(ctx context.Context) error
}

type PluginRegistry struct {
    renderers  map[string]contracts.Renderer
    matchers   map[string]contracts.Matcher
    providers  map[string]contracts.AIProvider
    exporters  map[string]contracts.Exporter
}
```

---

## 12. Security Design

### 12.1 Security Principles

```
┌─────────────────────────────────────────────────────────┐
│              Security Principles                    │
├─────────────────────────────────────────────────────────┤
│                                                  │
│  1. PRINCIPLE OF LEAST PRIVILEGE                  │
│     - Minimal permissions required                 │
│     - No root access required                     │
│     - User data isolation                       │
│                                                  │
│  2. INPUT VALIDATION                           │
│     - Sanitize all user inputs                  │
│     - Validate before processing                 │
│     - Prevent injection attacks                  │
│                                                  │
│  3. SECRETS MANAGEMENT                         │
│     - API keys from environment variables       │
│     - No hardcoded credentials                 │
│     - Config file permissions (600)             │
│                                                  │
│  4. DATA PROTECTION                           │
│     - Encrypted at rest (optional)             │
│     - Secure data disposal                      │
│     - Backup encryption                        │
│                                                  │
└─────────────────────────────────────────────────────────┘
```

### 12.2 Configuration Security

```yaml
security:
  # API keys from environment only
  apiKeys:
    anthropic: "${ANTHROPIC_API_KEY}"
    openai: "${OPENAI_API_KEY}"
  
  # Config file permissions
  configFileMode: "0600"
  
  # Audit logging
  audit:
    enabled: true
    path: ~/.otr/audit.log
```

---

## 13. Performance Considerations

### 13.1 Performance Budget

| Operation | Target | Maximum |
|-----------|--------|----------|
| Pattern matching (local) | < 50ms | < 200ms |
| Pattern matching (semantic) | < 200ms | < 500ms |
| UI render | < 16ms | < 33ms |
| AI generation | < 10s | < 30s |
| Database query | < 10ms | < 50ms |

### 13.2 Optimization Strategies

```
┌─────────────────────────────────────────────────────────┐
│              Optimization Strategies                 │
├─────────────────────────────────────────────────────────┤
│                                                  │
│  1. CACHING LAYERS                               │
│     • Match results LRU cache                     │
│     • Embeddings cache (semantic matching)      │
│     • Configuration cache                        │
│                                                  │
│  2. LAZY LOADING                                │
│     • Defer non-critical initialization        │
│     • Load patterns on demand                   │
│     • Async AI provider initialization          │
│                                                  │
│  3. CONCURRENT OPERATIONS                        │
│     • Parallel matching strategies              │
│     • Async AI calls                           │
│     • Background decay calculations              │
│                                                  │
│  4. EFFICIENT DATA STRUCTURES                   │
│     • Trie for prefix matching                  │
│     • Vector index for semantic search          │
│     • B-tree for range queries                 │
│                                                  │
└─────────────────────────────────────────────────────────┘
```

### 13.3 Memory Budget

| Resource | Budget |
|----------|--------|
| Base footprint | < 50MB |
| Per pattern | < 10KB |
| Cache (1000 patterns) | < 100MB |
| Total (10000 patterns) | < 150MB |

---

## 14. Deployment Architecture

### 14.1 Deployment Models

```
┌─────────────────────────────────────────────────────────┐
│              Deployment Models                       │
├─────────────────────────────────────────────────────────┤
│                                                  │
│  LOCAL DEVELOPMENT                                │
│  ┌─────────────────────────┐                     │
│  │ otr binary             │                     │
│  │ - SQLite storage      │                     │
│  │ - Local AI API       │                     │
│  └─────────────────────────┘                     │
│                                                  │
│  PRODUCTION (Single User)                       │
│  ┌─────────────────────────┐                     │
│  │ otr binary             │                      │
│  │ - SQLite or Redis    │                      │
│  │ - Cloud AI API       │                      │
│  └─────────────────────────┘                     │
│                                                  │
│  ENTERPRISE (Multi-user)                        │
│  ┌─────────────────────────┐                     │
│  │ API Server + Web UI   │                      │
│  │ - PostgreSQL          │                      │
│  │ - Redis cache        │                      │
│  │ - Auth service       │                      │
│  └─────────────────────────┘                     │
│                                                  │
└─────────────────────────────────────────────────────────┘
```

### 14.2 Distribution

| Method | Description |
|--------|-------------|
| **Go Binary** | `go build -o otr` |
| **Homebrew** | `brew install openclaw/tap/otr` |
| **Docker** | Multi-stage build, ~20MB image |
| **NPM** | `npm install -g @openclaw/reflex` |

### 14.3 Docker Configuration

```dockerfile
# Multi-stage build for minimal image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o otr ./cmd/cli

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/otr /usr/local/bin/
ENTRYPOINT ["otr"]
```

---

## Summary

### Architecture Highlights

| Aspect | Design |
|--------|--------|
| **Separation** | CLI UI completely isolated from core logic |
| **Modularity** | Each layer has clear responsibilities |
| **Testability** | Core logic testable without CLI |
| **Extensibility** | Plugin system for custom extensions |
| **Performance** | Caching, concurrency, efficient data structures |
| **Security** | No hardcoded secrets, input validation |

### Key Interfaces

| Interface | Purpose |
|-----------|---------|
| contracts.Matcher | Pattern matching |
| contracts.Storage | Data persistence |
| contracts.Renderer | Output formatting |
| contracts.AIProvider | AI service integration |

---

**Document Version**: v1.0  
**Created**: 2026-02-20  
**Project**: open-think-reflex
