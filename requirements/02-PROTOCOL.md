# Open-Think-Reflex: Protocol Specification

> **Version**: v1.0-draft  
> **Status**: Working Draft  
> **Scope**: Protocol & Interface Design (Implementation Agnostic)

---

## 1. Core Philosophy

### 1.1 Design Principles

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  1. Protocol-First Design                                           │
│     - Define interfaces, not implementations                         │
│     - Allow multiple backends (JSON, SQL, Vector, Graph)            │
│     - Language and framework agnostic                               │
│                                                                      │
│  2. Biological Inspiration                                        │
│     - Mimic human memory formation and decay                      │
│     - Strength through repetition                                  │
│     - Decay without reinforcement                                 │
│                                                                      │
│  3. Layered Architecture                                          │
│     - Protocol Layer (what and when)                             │
│     - Storage Layer (where)                                      │
│     - Interface Layer (how)                                       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Biological Model Reference

### 2.1 Human Memory System

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Human Memory System                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐                                                   │
│  │   Sensory   │──► Short-term Memory (seconds)                  │
│  │   Memory    │                                                  │
│  └──────┬──────┘                                                   │
│         │                                                          │
│         ▼                                                          │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  Working   │──►│   Long-term│──►│  Procedural │          │
│  │  Memory    │    │   Memory   │    │   Memory    │          │
│  │ (seconds)  │    │ (minutes+) │    │ (permanent) │          │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘          │
│         │                   │                   │                    │
│         ▼                   ▼                   ▼                    │
│     Attention      Consolidation      Automatic               │
│     Required       Required          Response                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Reflex Formation Process

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Reflex Formation Process                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. Stimulus (Input)                                              │
│     └── External trigger (user action, conversation)                │
│                                                                      │
│  2. Perception (Processing)                                        │
│     └── Pattern recognition and classification                     │
│                                                                      │
│  3. Memory Consolidation                                          │
│     └── Temporary storage (Working Memory)                         │
│                                                                      │
│  4. Reinforcement (Strengthening)                                  │
│     └── Repeated exposure increases synaptic weight               │
│                                                                      │
│  5. Threshold Crossing                                             │
│     └── Connection becomes automatic response                      │
│                                                                      │
│  6. Automatic Activation                                          │
│     └── Stimulus triggers response without conscious thought       │
│                                                                      │
│  7. Decay (If Not Used)                                         │
│     └── Synaptic weight decreases without reinforcement           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Key Biological Mechanisms

| Mechanism | Biological Process | System Equivalent |
|-----------|-------------------|-------------------|
| **LTP (Long-Term Potentiation)** | Repeated stimulation strengthens neural connections | Reflex strengthening |
| **Synaptic Pruning** | Weak connections are eliminated | Reflex decay |
| **Consolidation** | Short-term to long-term memory | Threshold crossing |
| **Automatic Activation** | Reflex arc bypasses brain | Pattern matching without reasoning |

---

## 3. Protocol Design

### 3.1 Core Protocol: REFLEX Protocol

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    REFLEX Protocol Layer                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Protocol Stack:                                                   │
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │                    REFLEX Protocol                        │   │
│  ├───────────────────────────────────────────────────────────────┤   │
│  │  1. Stimulus      │  Input pattern definition          │   │
│  │  2. Perception    │  Pattern recognition interface    │   │
│  │  3. Consolidation │  Memory storage interface        │   │
│  │  4. Reinforcement│  Strength update interface      │   │
│  │  5. Activation   │  Response trigger interface      │   │
│  │  6. Decay        │  Degradation interface          │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Interface Definitions

#### 3.2.1 Stimulus Interface

```typescript
interface Stimulus {
  // What triggers the reflex
  id: string;                    // Unique identifier
  type: 'text' | 'action' | 'emotion';  // Stimulus type
  content: string;              // Raw content
  metadata: {
    timestamp: number;
    source: string;             // User, System, AI
    context?: string;          // Conversation context
    tags?: string[];           // User-defined tags
  };
}

interface StimulusHandler {
  // Process incoming stimulus
  process(input: Stimulus): PatternCandidate;
  
  // Extract candidate patterns from stimulus
  extractCandidates(stimulus: Stimulus): PatternCandidate[];
}
```

#### 3.2.2 Perception Interface

```typescript
interface Pattern {
  // Extracted pattern from stimulus
  id: string;
  fingerprint: string;           // Hash of pattern content
  signature: string;            // Semantic signature
  strength: number;             // Current strength (0-100)
  threshold: number;            // Activation threshold
  metadata: {
    created: number;
    updated: number;
    reinforcementCount: number;
    decayCount: number;
  };
  response?: Response;           // Associated response
}

interface PerceptionSystem {
  // Extract patterns from stimulus
  extract(stimulus: Stimulus): PatternCandidate[];
  
  // Match input against existing patterns
  match(input: string, options?: MatchOptions): MatchResult[];
  
  // Calculate pattern similarity
  similarity(pattern1: Pattern, pattern2: Pattern): number;
}

interface MatchOptions {
  threshold?: number;           // Minimum match score
  limit?: number;              // Max results
  includeLatent?: boolean;     // Include below-threshold patterns
}
```

#### 3.2.3 Consolidation Interface

```typescript
interface Memory {
  // Long-term storage structure
  id: string;
  type: 'episodic' | 'procedural' | 'semantic';
  content: Pattern | PatternSequence;
  strength: number;            // Current strength (0-100)
  activationThreshold: number;  // Threshold to activate
  metadata: {
    created: number;
    lastActivated: number;
    accessCount: number;
    decayRate: number;         // How fast it decays
  };
  connections: Connection[];   // Linked memories
}

interface MemoryConsolidationSystem {
  // Store new pattern
  store(pattern: PatternCandidate): Pattern;
  
  // Update existing pattern
  update(id: string, updates: Partial<Pattern>): void;
  
  // Retrieve by ID
  retrieve(id: string): Pattern | null;
  
  // Retrieve by fingerprint
  retrieveByFingerprint(fingerprint: string): Pattern[];
  
  // Delete pattern
  delete(id: string): void;
}
```

#### 3.2.4 Reinforcement Interface

```typescript
interface ReinforcementSignal {
  // Signal to strengthen a pattern
  patternId: string;
  strength: number;            // Amount to increase
  type: 'positive' | 'negative' | 'neutral';
  metadata: {
    source: 'user' | 'system' | 'inferred';
    reason?: string;
    evidence?: string[];
  };
}

interface ReinforcementSystem {
  // Apply reinforcement signal
  reinforce(signal: ReinforcementSignal): void;
  
  // Batch reinforcement
  reinforceBatch(signals: ReinforcementSignal[]): void;
  
  // Check if threshold crossed
  checkThreshold(patternId: string): ActivationEvent | null;
}
```

#### 3.2.5 Activation Interface

```typescript
interface ActivationEvent {
  // Triggered when pattern crosses threshold
  patternId: string;
  triggeredBy: string;          // Input that triggered
  timestamp: number;
  automatic: boolean;          // True if automatic activation
  
  // Proposed response
  response: Response;
  
  // Alternative patterns (below threshold but nearby)
  alternatives?: MatchResult[];
}

interface ActivationSystem {
  // Check if any patterns should activate
  checkActivations(input: Stimulus): ActivationEvent[];
  
  // Manually activate a pattern
  activate(patternId: string, context: Stimulus): ActivationEvent;
  
  // Get all active patterns
  getActive(): Pattern[];
}
```

#### 3.2.6 Decay Interface

```typescript
interface DecayEvent {
  // Triggered when pattern strength drops
  patternId: string;
  previousStrength: number;
  currentStrength: number;
  triggeredBy: 'time' | 'supersession' | 'degradation';
}

interface DecaySystem {
  // Calculate decayed strength
  calculateDecay(pattern: Pattern, timeDelta: number): number;
  
  // Apply decay
  applyDecay(patternId: string): DecayEvent | null;
  
  // Batch decay calculation
  calculateAllDecay(timeDelta: number): Map<string, number>;
  
  // Prune dead patterns
  pruneDead(): string[];        // Returns pruned pattern IDs
  
  // Re-activation handling
  canReactivate(patternId: string): boolean;
}
```

### 3.3 Data Flow Protocol

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    REFLEX Protocol Data Flow                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                                                                      │
│   INPUT (Stimulus)                                                  │
│        │                                                            │
│        ▼                                                            │
│   ┌─────────────────┐                                               │
│   │  Perception     │──► Extract patterns                          │
│   │  System         │                                              │
│   └────────┬────────┘                                               │
│            │                                                         │
│            ▼                                                         │
│   ┌─────────────────┐     ┌──────────────────────────────────┐     │
│   │  Memory        │◄───│  Reinforcement                   │     │
│   │  Consolidation  │     │  System                          │     │
│   └────────┬────────┘     │  - Update strength              │     │
│            │               │  - Check threshold              │     │
│            │               └──────────────────────────────────┘     │
│            │                                                         │
│            ▼                                                         │
│   ┌─────────────────┐                                               │
│   │  Activation     │──► Trigger responses                       │
│   │  System         │                                              │
│   └────────┬────────┘                                               │
│            │                                                         │
│            ▼                                                         │
│   ┌─────────────────┐                                               │
│   │  Decay         │──► Update strength                          │
│   │  System        │                                              │
│   └─────────────────┘                                               │
│                                                                      │
│   All flows are reversible and observable                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Pattern Representation Protocol

### 4.1 Pattern Structure

```typescript
interface PatternSchema {
  // Abstract schema, implementation-agnostic
  id: string;                    // Unique identifier
  version: number;               // Schema version
  
  // What triggers this pattern
  trigger: {
    type: 'text' | 'action' | 'sequence' | 'emotion';
    signature: string;           // Hash of trigger signature
    fingerprint: string;         // Semantic fingerprint
    examples?: string[];          // Exemplar triggers
  };
  
  // What response it produces
  response: {
    type: 'action' | 'template' | 'chain';
    content: any;               // Response definition
    alternatives?: Response[];   // Alternative responses
  };
  
  // Strength tracking
  strength: {
    current: number;            // 0-100
    threshold: number;          // Activation threshold
    max: number;               // Maximum possible
    reinforcement: {
      positive: number;        // Strength gain per positive
      negative: number;         // Strength loss per negative
      decay: number;            // Per time unit decay
    };
  };
  
  // Metadata
  metadata: {
    created: number;           // Timestamp
    updated: number;            // Last modification
    version: number;           // Pattern version
    tags: string[];           // Categorization
    source: string;            // Origin of pattern
  };
}
```

### 4.2 Pattern Graph Protocol

```typescript
interface PatternGraph {
  // Graph structure of pattern relationships
  nodes: Map<string, PatternNode>;
  edges: Map<string, PatternEdge[]>;
  
  // Graph operations
  operations: {
    // Connect two patterns
    connect(sourceId: string, targetId: string, type: ConnectionType): void;
    
    // Disconnect patterns
    disconnect(sourceId: string, targetId: string): void;
    
    // Find connected patterns
    getConnected(patternId: string, options?: ConnectionOptions): PatternNode[];
    
    // Get shortest path between patterns
    getPath(fromId: string, toId: string): PatternNode[] | null;
  };
}

interface PatternNode {
  id: string;
  pattern: Pattern;
  activationLevel: number;
}

interface PatternEdge {
  source: string;
  target: string;
  type: 'causal' | 'temporal' | 'semantic';
  strength: number;              // Connection strength
  bidirectional: boolean;
}
```

---

## 5. Matching Algorithm Protocol

### 5.1 Matching Interface

```typescript
interface MatchingProtocol {
  // Strategy-agnostic matching interface
  name: string;
  version: string;
  
  // Match input against patterns
  match(input: string, patterns: Pattern[], options?: MatchOptions): MatchResult[];
  
  // Get match confidence
  calculateConfidence(input: string, pattern: Pattern): number;
  
  // Get match reason (explain why)
  explainMatch(input: string, pattern: Pattern): MatchExplanation;
}

interface MatchResult {
  patternId: string;
  score: number;               // 0-100 confidence
  matchedElement: string;       // What part matched
  explanation: string;         // Human-readable reason
}

interface MatchOptions {
  // Matching constraints
  minScore?: number;           // Minimum score to include
  maxResults?: number;         // Limit results
  includeBelowThreshold?: boolean;
  
  // Strategy selection
  strategies?: string[];      // Use specific strategies
  
  // Context
  context?: string;            // Conversation context
  userPreferences?: Map<string, number>;
}
```

### 5.2 Matching Strategies (Abstract)

```typescript
// Strategy definitions, implementations are separate
type MatchingStrategy = 
  | 'exact'      // Exact string matching
  | 'keyword'    // Keyword extraction
  | 'semantic'   // Semantic similarity
  | 'fuzzy'      // Fuzzy string matching
  | 'pattern'    // Regex pattern matching
  | 'embedding'; // Vector embedding similarity

// Strategy configuration
interface StrategyConfig {
  name: MatchingStrategy;
  weight: number;             // Importance in combined scoring
  enabled: boolean;
  params?: Record<string, any>; // Strategy-specific parameters
}
```

---

## 6. Storage Protocol

### 6.1 Storage Abstraction

```typescript
interface StorageBackend {
  // Abstract storage interface
  // Can be implemented as JSON, SQL, Vector DB, Graph DB, etc.
  
  name: string;
  version: string;
  
  // Pattern operations
  patterns: PatternStorage;
  
  // Graph operations
  graph: GraphStorage;
  
  // Index operations
  indexes: IndexStorage;
}

interface PatternStorage {
  create(pattern: Pattern): Pattern;
  read(id: string): Pattern | null;
  readByFingerprint(fingerprint: string): Pattern[];
  update(pattern: Pattern): void;
  delete(id: string): void;
  list(query?: StorageQuery): Pattern[];
}

interface GraphStorage {
  createEdge(edge: PatternEdge): void;
  readEdges(sourceId: string): PatternEdge[];
  deleteEdge(source: string, target: string): void;
  getConnected(patternId: string): string[];
}

interface IndexStorage {
  createIndex(field: string, type: IndexType): void;
  search(field: string, value: any): string[];
  reindex(): void;
}
```

### 6.2 Supported Backend Types

| Backend | Use Case | Strengths |
|---------|----------|-----------|
| **JSON** | Simple, small-scale | Human-readable, version control |
| **SQLite** | Medium-scale, relational | Queries, transactions |
| **Vector DB** | Semantic matching | Fast similarity search |
| **Graph DB** | Pattern relationships | Relationship traversal |
| **Hybrid** | Complex systems | Best of all worlds |

---

## 7. Visualization Protocol

### 7.1 Tree/Graph Display Interface

```typescript
interface VisualizationProtocol {
  // Strategy-agnostic visualization interface
  
  // Render pattern tree/graph
  render(data: PatternGraph, options?: RenderOptions): VisualOutput;
  
  // Get user selection
  getSelection(): PatternSelection | null;
  
  // Highlight paths
  highlightPaths(paths: string[][]): void;
}

interface RenderOptions {
  layout: 'tree' | 'graph' | 'radial';
  depth?: number;                // Max depth to show
  highlightActive?: boolean;
  showStrength?: boolean;
  showThreshold?: boolean;
  interactive: boolean;          // Allow user selection
}

interface VisualOutput {
  type: 'svg' | 'canvas' | 'html' | 'terminal';
  content: any;                  // Rendering result
  metadata: {
    nodeCount: number;
    edgeCount: number;
    activePath?: string[];
  };
}

interface PatternSelection {
  selectedIds: string[];
  mode: 'single' | 'multiple' | 'path';
  confirmed: boolean;
}
```

---

## 8. Protocol Versioning

### 8.1 Version Compatibility

```typescript
const PROTOCOL_VERSION = {
  major: 1,
  minor: 0,
  patch: 0,
  
  // Compatibility
  isCompatible(other: Version): boolean;
  
  // Migration
  migrate(data: any, fromVersion: Version): any;
};
```

---

## 9. Error Handling Protocol

### 9.1 Error Codes

```typescript
enum ReflexErrorCode {
  // Pattern errors
  PATTERN_NOT_FOUND = 'E001',
  PATTERN_INVALID = 'E002',
  PATTERN_DUPLICATE = 'E003',
  
  // Storage errors
  STORAGE_ERROR = 'S001',
  STORAGE_FULL = 'S002',
  
  // Matching errors
  MATCH_ERROR = 'M001',
  MATCH_TIMEOUT = 'M002',
  
  // Activation errors
  ACTIVATION_ERROR = 'A001',
  THRESHOLD_ERROR = 'A002',
  
  // Decay errors
  DECAY_ERROR = 'D001',
  PRUNE_ERROR = 'D002',
}

interface ReflexError {
  code: ReflexErrorCode;
  message: string;
  details?: Record<string, any>;
  recoverable: boolean;
  suggestion?: string;
}
```

---

## 10. Core Workflows

### 10.1 Pattern Formation Workflow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Pattern Formation Workflow                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. STIMULUS INPUT                                                │
│     └─► Extract trigger pattern                                   │
│                                                                      │
│  2. PATTERN MATCHING                                             │
│     └─► Match against existing patterns                            │
│                                                                      │
│  3. NEW PATTERN?                                                 │
│     ├─ YES ──► Create new pattern                               │
│     │           └─► Store with initial strength                  │
│     │                                                               │
│     └─ NO ──► Existing pattern                                   │
│                 └─► Apply reinforcement                         │
│                                                                      │
│  4. THRESHOLD CHECK                                             │
│     └─► If crossed ──► Emit activation event                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Response Selection Workflow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Response Selection Workflow                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. USER INPUT                                                   │
│     └─► Match against all patterns                               │
│                                                                      │
│  2. GET MATCH RESULTS                                            │
│     └─► Sort by score                                           │
│                                                                      │
│  3. FILTER BY THRESHOLD                                         │
│     └─► Get above-threshold matches                             │
│                                                                      │
│  4. BUILD PATHS                                                  │
│     └─► Construct response paths from matches                    │
│                                                                      │
│  5. DISPLAY TO USER                                             │
│     └─► Show tree/graph with alternatives                        │
│                                                                      │
│  6. USER SELECTION                                               │
│     ├─ MANUAL ──► Use selected path                            │
│     │                                                            │
│     └─ AUTO ──► Use highest score (if threshold met)           │
│                                                                      │
│  7. EXECUTE RESPONSE                                            │
│     └─► Return selected response                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.3 Decay Workflow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Decay Workflow                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  PERIODIC TRIGGER (e.g., every hour)                               │
│        │                                                           │
│        ▼                                                           │
│  FOR EACH PATTERN                                                 │
│        │                                                           │
│        ▼                                                           │
│  CALCULATE TIME DELTA                                             │
│        │                                                           │
│        ▼                                                           │
│  APPLY DECAY FORMULA                                              │
│     strength = strength × (1 - decayRate × timeDelta)              │
│        │                                                           │
│        ▼                                                           │
│  STRENGTH < THRESHOLD?                                            │
│        │                                                           │
│        ├─ YES ──► Mark as latent (below threshold)               │
│        │                                                            │
│        └─ NO ──► Continue active                               │
│                                                                      │
│  STRENGTH NEAR ZERO?                                              │
│        │                                                           │
│        └─► YES ──► Consider for pruning                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 11. Implementation Registry

### 11.1 Protocol Implementation Template

```typescript
interface ReflexImplementation {
  // Metadata
  name: string;
  version: string;
  protocolVersion: string;      // Which protocol version this implements
  
  // Protocol implementations
  perception: PerceptionSystem;
  consolidation: MemoryConsolidationSystem;
  reinforcement: ReinforcementSystem;
  activation: ActivationSystem;
  decay: DecaySystem;
  
  // Storage
  storage: StorageBackend;
  
  // Configuration
  config: ImplementationConfig;
  
  // Lifecycle
  initialize(): Promise<void>;
  shutdown(): Promise<void>;
  healthCheck(): HealthStatus;
}
```

---

## 12. Summary

### 12.1 Protocol Core Abstractions

| Layer | Interfaces | Purpose |
|--------|-------------|----------|
| **Stimulus** | Stimulus, StimulusHandler | Input handling |
| **Perception** | Pattern, PerceptionSystem | Pattern recognition |
| **Consolidation** | Memory, MemoryConsolidationSystem | Storage |
| **Reinforcement** | ReinforcementSignal, ReinforcementSystem | Strength updates |
| **Activation** | ActivationEvent, ActivationSystem | Response triggering |
| **Decay** | DecayEvent, DecaySystem | Degradation |

### 12.2 Key Design Decisions

| Decision | Approach |
|---------|----------|
| **Storage** | Abstract backend (JSON/SQL/Vector/Graph) |
| **Matching** | Strategy pattern (exact/semantic/fuzzy) |
| **Visualization** | Protocol-agnostic (SVG/Canvas/Terminal) |
| **Error Handling** | Typed error codes with suggestions |
| **Versioning** | Semantic versioning with migration support |

---

**Document Version**: v1.0-draft  
**Protocol Version**: 1.0.0  
**Status**: Working Draft - For Discussion
