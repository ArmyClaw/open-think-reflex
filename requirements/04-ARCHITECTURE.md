# Open-Think-Reflex: Architecture & Implementation Views (v2.0)

> **Version**: v2.0  
> **Focus**: System Architecture & Implementation Views  
> **Perspective**: Technical & Architectural

---

## 1. System Architecture Overview

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Open-Think-Reflex Architecture                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                    CLI Interface Layer                     │       │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐│       │
│  │  │  Shell   │ │ Terminal │ │  API    │ │ Web UI ││       │
│  │  └───────────┘ └───────────┘ └───────────┘ └──────────┘│       │
│  └───────────────────────────────────────────────────────────────┘       │
│                              │                                         │
│                              ▼                                         │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                   Core Engine Layer                        │       │
│  │                                                          │       │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐│       │
│  │  │ Stimulus  │ │ Perception │ │Consolidate│ │Activate ││       │
│  │  │ Processor │ │  Engine   │ │   System  │ │  Engine ││       │
│  │  └───────────┘ └───────────┘ └───────────┘ └──────────┘│       │
│  │                                                          │       │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐            │       │
│  │  │Decay     │ │Reinforce │ │Visualize │            │       │
│  │  │ Engine   │ │ System   │ │ Engine   │            │       │
│  │  └───────────┘ └───────────┘ └───────────┘            │       │
│  │                                                          │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                              │                                         │
│                              ▼                                         │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                   Storage Layer                             │       │
│  │                                                          │       │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐│       │
│  │  │ Pattern   │ │  Graph   │ │  Index   │ │  Cache  ││       │
│  │  │  Store   │ │  Store   │ │  Store   │ │  Store  ││       │
│  │  └───────────┘ └───────────┘ └───────────┘ └──────────┘│       │
│  │                                                          │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Component Dependencies                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────┐                                                    │
│   │   CLI    │                                                    │
│   └────┬─────┘                                                    │
│        │                                                          │
│        ▼                                                          │
│   ┌──────────┐                                                    │
│   │  Router │                                                    │
│   └────┬─────┘                                                    │
│        │                                                          │
│        ▼                                                          │
│   ┌──────────┐                                                    │
│   │Controller│                                                    │
│   └────┬─────┘                                                    │
│        │                                                          │
│   ┌────┴─────────────────────────────────────────────┐             │
│   │                                          │                   │
│   ▼                                          ▼                   │
│   ┌──────────┐                    ┌──────────────┐              │
│   │Stimulus │                    │  Orchestrator │              │
│   │Processor│                    └──────┬───────┘              │
│   └────┬─────┘                           │                      │
│        │                                  │                      │
│        ▼                                  │                      │
│   ┌──────────┐                           │                      │
│   │Perception│                           │                      │
│   │  Engine  │                           │                      │
│   └────┬─────┘                           │                      │
│        │                                  │                      │
│   ┌────┴───────────────────────────────────────┐              │
│   │                                        │                  │
│   ▼                                        ▼                  │
│   ┌──────────┐                   ┌──────────────┐             │
│   │ Memory  │                   │  Decay     │             │
│   │Consolidate│                   │  Engine    │             │
│   └────┬─────┘                   └──────┬───────┘             │
│        │                                  │                    │
│   ┌────┴───────────────────────────────────────┐            │
│   │                                        │                │
│   ▼                                        ▼                │
│   ┌──────────┐                   ┌──────────────┐          │
│   │ Activation│                   │ Visualization│          │
│   │  Engine  │                   │   Engine    │          │
│   └────┬─────┘                   └──────┬───────┘          │
│        │                                  │                    │
│        └──────────────┬─────────────────┘                    │
│                       ▼                                        │
│               ┌──────────────┐                               │
│               │    CLI      │                               │
│               │   Output    │                               │
│               └──────────────┘                               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Data Models

### 2.1 Pattern Entity

```typescript
interface Pattern {
  // Identity
  id: string;
  version: number;
  
  // Trigger
  trigger: {
    type: 'text' | 'action' | 'emotion';
    signature: string;          // Hash
    fingerprint: string;        // Semantic
    examples: string[];         // Exemplars
  };
  
  // Response
  response: {
    type: 'action' | 'template' | 'chain';
    content: unknown;          // Varies by type
    alternatives: Response[];
  };
  
  // Strength tracking
  strength: {
    current: number;           // 0-100
    threshold: number;         // Activation threshold
    reinforcement: {
      positive: number;        // Strength gain
      negative: number;        // Strength loss
      decay: number;           // Per unit decay
    };
  };
  
  // Lifecycle
  lifecycle: {
    created: number;           // Timestamp
    updated: number;
    lastActivated: number;
    activationCount: number;
  };
  
  // Connections
  connections: Connection[];
  
  // Metadata
  metadata: {
    tags: string[];
    source: string;           // Origin
    userId?: string;
    privacy: 'public' | 'private';
  };
}

interface Response {
  id: string;
  type: 'action' | 'template' | 'chain';
  content: unknown;
  score?: number;             // Preference score
}

interface Connection {
  targetPatternId: string;
  type: 'causal' | 'temporal' | 'semantic';
  strength: number;            // Connection strength
}
```

### 2.2 Conversation Context

```typescript
interface ConversationContext {
  id: string;
  userId: string;
  
  // Current interaction
  currentInput: string;
  matchedPatterns: MatchResult[];
  selectedPattern?: string;
  
  // History for this conversation
  history: ConversationTurn[];
  
  // Working memory
  workingMemory: WorkingMemory;
  
  // Session metadata
  metadata: {
    startTime: number;
    turnCount: number;
    avgResponseTime: number;
  };
}

interface ConversationTurn {
  timestamp: number;
  input: string;
  response: string;
  matchedPatterns: string[];
  userFeedback?: 'positive' | 'negative' | 'neutral';
}

interface WorkingMemory {
  shortTerm: PatternCandidate[];
  longTerm: Pattern[];
  active: Pattern[];
  latent: Pattern[];
}
```

---

## 3. State Management

### 3.1 Application State

```typescript
interface ApplicationState {
  // Core engine state
  engine: {
    status: 'idle' | 'processing' | 'error';
    lastActivity: number;
    activityCount: number;
  };
  
  // Pattern state
  patterns: {
    total: number;
    active: number;
    latent: number;
    strengthDistribution: Map<number, number>;
  };
  
  // Storage state
  storage: {
    backend: string;
    size: number;
    lastSync: number;
  };
  
  // UI state
  ui: {
    currentView: 'tree' | 'list' | 'detail';
    selectedPattern?: string;
    expandedNodes: string[];
    filter?: PatternFilter;
  };
  
  // User preferences
  preferences: {
    autoThreshold: number;
    displayThreshold: number;
    decayEnabled: boolean;
    visualizationLayout: 'tree' | 'graph' | 'radial';
  };
}
```

### 3.2 State Machine

```typescript
enum PatternState {
  NEW = 'new',
  FORMING = 'forming',
  ACTIVE = 'active',
  LATENT = 'latent',
  DECAYING = 'decaying',
  LOST = 'lost',
  ARCHIVED = 'archived'
}

interface PatternStateMachine {
  // State transitions
  transitions: {
    [from: PatternState]: {
      [to: PatternState]: {
        condition: (pattern: Pattern) => boolean;
        action: (pattern: Pattern) => void;
      };
    };
  };
  
  // Transition rules
  canTransition(pattern: Pattern, newState: PatternState): boolean;
  transition(pattern: Pattern, newState: PatternState): Pattern;
}
```

---

## 4. API Design

### 4.1 REST API

```
BASE PATH: /api/v1

┌─────────────────────────────────────────────────────────────────┐
│                    Pattern Management                          │
├─────────────────────────────────────────────────────────────────┤
│                                                              │
│  GET    /patterns              # List all patterns           │
│  GET    /patterns/:id         # Get single pattern        │
│  POST   /patterns             # Create new pattern         │
│  PUT    /patterns/:id         # Update pattern            │
│  DELETE /patterns/:id         # Delete pattern            │
│  POST   /patterns/:id/reinforce  # Reinforce pattern     │
│  POST   /patterns/:id/decay  # Trigger decay            │
│                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Matching & Activation                      │
├─────────────────────────────────────────────────────────────────┤
│                                                              │
│  POST   /match                 # Match input to patterns │
│  POST   /match/batch         # Batch matching           │
│  POST   /activate/:id        # Manual activation       │
│  GET    /activations/current   # Get current activations │
│                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Visualization                             │
├─────────────────────────────────────────────────────────────────┤
│                                                              │
│  GET    /visualize/tree        # Tree visualization        │
│  GET    /visualize/graph      # Graph visualization      │
│  GET    /visualize/stats      # Statistics               │
│                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Configuration                             │
├─────────────────────────────────────────────────────────────────┤
│                                                              │
│  GET    /config               # Get configuration        │
│  PUT    /config              # Update configuration    │
│  POST   /config/reset       # Reset to defaults      │
│                                                              │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 WebSocket API

```
CONNECTION: ws://localhost:3000/ws

┌─────────────────────────────────────────────────────────────────┐
│                    Real-time Events                         │
├─────────────────────────────────────────────────────────────────┤
│                                                              │
│  OUTGOING:                                                  │
│  ─────────                                                  │
│  { event: 'pattern.matched', data: MatchResult }          │
│  { event: 'pattern.activated', data: ActivationEvent }   │
│  { event: 'pattern.decayed', data: DecayEvent }          │
│  { event: 'pattern.strength', data: { id, strength } }  │
│                                                              │
│  INCOMING:                                                 │
│  ──────────                                                 │
│  { event: 'input', data: { text: string } }            │
│  { event: 'select', data: { patternId: string } }       │
│  { event: 'feedback', data: { patternId, feedback } }    │
│  { event: 'subscribe', data: { patternIds: string[] } }  │
│                                                              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 5. Storage Architecture

### 5.1 Storage Layers

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Storage Architecture                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                    Hot Storage (Cache)                        │     │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐               │     │
│  │  │   LRU    │ │  Bloom   │ │  Count   │               │     │
│  │  │  Cache   │ │  Filter  │ │  Min-Heap│               │     │
│  │  └───────────┘ └───────────┘ └───────────┘               │     │
│  │  Access: < 1ms                                              │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                   Warm Storage (SQLite)                     │     │
│  │                                                             │     │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌──────────┐ │     │
│  │  │  Patterns │ │ Connections│ │  History │ │ Sessions │ │     │
│  │  │   Table  │ │   Table   │ │   Table   │ │  Table   │ │     │
│  │  └───────────┘ └───────────┘ └───────────┘ └──────────┘ │     │
│  │  Access: 1-10ms                                             │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                   Cold Storage (Archive)                    │     │
│  │                                                             │     │
│  │  Patterns exported to JSON files                            │     │
│  │  Archived patterns for long-term storage                     │     │
│  │  Access: seconds (on-demand load)                           │     │
│  │                                                             │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Database Schema

```sql
-- Patterns table
CREATE TABLE patterns (
  id TEXT PRIMARY KEY,
  trigger_signature TEXT NOT NULL,
  trigger_fingerprint TEXT,
  trigger_examples TEXT, -- JSON array
  response_type TEXT NOT NULL,
  response_content TEXT NOT NULL, -- JSON
  response_alternatives TEXT, -- JSON array
  
  strength_current REAL NOT NULL DEFAULT 0,
  strength_threshold REAL NOT NULL DEFAULT 50,
  reinforcement_positive REAL DEFAULT 5,
  reinforcement_negative REAL DEFAULT -5,
  reinforcement_decay REAL DEFAULT 0.1,
  
  created_at REAL NOT NULL,
  updated_at REAL NOT NULL,
  last_activated_at REAL,
  activation_count INTEGER DEFAULT 0,
  
  tags TEXT, -- JSON array
  source TEXT,
  user_id TEXT,
  privacy TEXT DEFAULT 'private',
  
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Connections table
CREATE TABLE connections (
  id TEXT PRIMARY KEY,
  source_pattern_id TEXT NOT NULL,
  target_pattern_id TEXT NOT NULL,
  connection_type TEXT NOT NULL, -- causal, temporal, semantic
  strength REAL NOT NULL DEFAULT 50,
  
  created_at REAL NOT NULL,
  
  FOREIGN KEY (source_pattern_id) REFERENCES patterns(id),
  FOREIGN KEY (target_pattern_id) REFERENCES patterns(id)
);

-- History table
CREATE TABLE history (
  id TEXT PRIMARY KEY,
  pattern_id TEXT NOT NULL,
  input TEXT NOT NULL,
  response TEXT NOT NULL,
  matched_patterns TEXT, -- JSON array
  selected_pattern_id TEXT,
  feedback TEXT, -- positive, negative, neutral
  
  created_at REAL NOT NULL,
  
  FOREIGN KEY (pattern_id) REFERENCES patterns(id)
);

-- Indexes
CREATE INDEX idx_patterns_signature ON patterns(trigger_signature);
CREATE INDEX idx_patterns_user ON patterns(user_id);
CREATE INDEX idx_patterns_strength ON patterns(strength_current);
CREATE INDEX idx_connections_source ON connections(source_pattern_id);
CREATE INDEX idx_connections_target ON connections(target_pattern_id);
```

---

## 6. Performance Architecture

### 6.1 Caching Strategy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Caching Layers                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Layer 1: In-Memory Cache (LRU)                                    │
│  ─────────────────────────────────────────────                         │
│  What's cached: Frequently accessed patterns                          │
│  Eviction: LRU                                                    │
│  Size: Configurable (default 1000)                                  │
│  TTL: No TTL (explicit eviction)                                   │
│                                                                      │
│  Layer 2: Bloom Filter                                            │
│  ───────────────────────────                                        │
│  What's filtered: Non-existent pattern lookups                     │
│  Size: 1M bits (configurable)                                    │
│  False positive rate: < 1%                                        │
│                                                                      │
│  Layer 3: Query Result Cache                                       │
│  ────────────────────────────                                      │
│  What's cached: Frequent query results                             │
│  Eviction: Time-based (1 hour)                                    │
│  Invalidation: On pattern update                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Performance Targets

| Operation | Target | 99th Percentile |
|-----------|--------|-----------------|
| Pattern match | < 50ms | < 100ms |
| Pattern activation | < 10ms | < 20ms |
| Pattern storage | < 20ms | < 50ms |
| Pattern retrieval | < 5ms | < 10ms |
| Tree visualization | < 200ms | < 500ms |
| CLI startup | < 500ms | < 1s |

---

## 7. Security Architecture

### 7.1 Threat Model

| Threat | Mitigation |
|--------|------------|
| **Pattern Injection** | Input validation, sandboxing |
| **Data Leakage** | Encryption at rest, access controls |
| **DoS Attacks** | Rate limiting, request throttling |
| **Privacy Violation** | User data isolation, anonymization |
| **Unauthorized Access** | Authentication, authorization |

### 7.2 Security Controls

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Security Layers                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                    Input Security                           │     │
│  │  - Input validation (schema-based)                        │     │
│  │  - Sanitization (XSS prevention)                        │     │
│  │  - Rate limiting (100 req/min per user)                 │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                   Data Security                            │     │
│  │  - Encryption at rest (AES-256)                          │     │
│  │  - Encryption in transit (TLS 1.3)                       │     │
│  │  - User data isolation                                   │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                              │                                     │
│                              ▼                                     │
│  ┌───────────────────────────────────────────────────────────────┐     │
│  │                   Access Security                          │     │
│  │  - Authentication (JWT)                                   │     │
│  │  - Authorization (RBAC)                                    │     │
│  │  - Audit logging                                        │     │
│  └───────────────────────────────────────────────────────────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Deployment Architecture

### 8.1 Deployment Options

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Deployment Models                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Option A: Local CLI (Single User)                                │
│  ────────────────────────────────────────                           │
│  ┌─────────────┐                                                 │
│  │  OpenThink  │                                                 │
│  │   Reflex    │                                                 │
│  │   CLI       │                                                 │
│  └─────────────┘                                                 │
│  └ SQLite DB                                                    │
│                                                                      │
│  Option B: Server (Multi-User)                                   │
│  ─────────────────────────────────                               │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐        │
│  │   Client   │──────│  API Server │──────│  SQLite   │        │
│  │   (Web/CLI)│      │   (Node)   │      │  Database │        │
│  └─────────────┘      └─────────────┘      └─────────────┘        │
│                                                                      │
│  Option C: Distributed (High Scale)                              │
│  ────────────────────────────────────────                          │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐        │
│  │  Client    │──────│  API Gateway│──────│   Redis   │        │
│  └─────────────┘      │   (Node)    │      │   Cache    │        │
│                       └──────┬──────┘      └─────────────┘        │
│                              │                                     │
│                       ┌─────┴─────┐      ┌─────────────┐        │
│                       │ PostgreSQL │      │ Vector DB  │        │
│                       │ Database  │      │ (Optional) │        │
│                       └───────────┘      └─────────────┘        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Scaling Strategy

| Component | Vertical Scaling | Horizontal Scaling |
|-----------|-----------------|-------------------|
| API Server | Up to 4CPU, 16GB | Load balancer + replicas |
| SQLite | Up to 100GB | Read replicas for queries |
| Cache | Increase size | Redis Cluster |
| Vector DB | - | Native sharding |

---

## 9. Monitoring & Observability

### 9.1 Metrics

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Key Metrics                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Business Metrics:                                                    │
│  - Patterns created per user                                         │
│  - Activation rate (%)                                               │
│  - Auto-apply rate (%)                                               │
│  - User satisfaction score                                            │
│                                                                      │
│  System Metrics:                                                     │
│  - Request latency (p50, p95, p99)                                  │
│  - Throughput (requests/second)                                       │
│  - Error rate (%)                                                    │
│  - Cache hit rate (%)                                                │
│                                                                      │
│  Resource Metrics:                                                  │
│  - CPU utilization                                                  │
│  - Memory usage                                                    │
│  - Disk I/O                                                        │
│  - Network I/O                                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Logging Strategy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Log Levels                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ERROR: Pattern failures, storage errors, security incidents         │
│  WARN:  Near-threshold patterns, degraded performance                │
│  INFO:  Pattern activations, user interactions                      │
│  DEBUG: Match scores, strength updates                             │
│  TRACE: Full request/response logging (development only)           │
│                                                                      │
│  Log Retention:                                                     │
│  - ERROR/WARN: 30 days                                             │
│  - INFO: 7 days                                                   │
│  - DEBUG/TRACE: 1 day                                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**Document Version**: v2.0  
**Focus**: Architecture & Implementation Views  
**Status**: Technical Draft for Review
