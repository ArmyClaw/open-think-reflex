# Open-Think-Reflex Data Schema

> **Version**: v1.0  
> **Status**: Active  
> **Scope**: Data models for v1.0 implementation

---

## Table of Contents

1. [Core Models](#1-core-models)
2. [Pattern Model](#2-pattern-model)
3. [Thought Chain Model](#3-thought-chain-model)
4. [Space/Project Model](#4-spaceproject-model)
5. [Event Model](#5-event-model)
6. [Configuration Model](#6-configuration-model)
7. [JSON Schemas](#7-json-schemas)
8. [Database Schema](#8-database-schema)

---

## 1. Core Models

### 1.1 Base Types

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

// ID represents a unique identifier
type ID string

// NewID generates a new unique identifier
func NewID() ID {
    return ID(uuid.New().String())
}

// Timestamp represents a time point
type Timestamp time.Time

// Now returns the current timestamp
func Now() Timestamp {
    return Timestamp(time.Now())
}

// Confidence represents a match confidence score (0-100)
type Confidence float64

// Strength represents pattern strength (0-100)
type Strength float64

// Tags represents a list of string tags
type Tags []string
```

---

## 2. Pattern Model

### 2.1 Pattern Structure

```go
// Pattern represents a reflex pattern
type Pattern struct {
    // Identification
    ID        ID        `json:"id" db:"id"`
    CreatedAt Timestamp `json:"created_at" db:"created_at"`
    UpdatedAt Timestamp `json:"updated_at" db:"updated_at"`
    
    // Core fields
    Trigger   string    `json:"trigger" db:"trigger"`
    Response  string    `json:"response" db:"response"`
    
    // Strength management
    Strength  Strength  `json:"strength" db:"strength"`
    Threshold Strength  `json:"threshold" db:"threshold"`
    DecayRate float64   `json:"decay_rate" db:"decay_rate"`
    DecayEnabled bool   `json:"decay_enabled" db:"decay_enabled"`
    
    // Statistics
    ReinforceCnt int    `json:"reinforcement_count" db:"reinforcement_count"`
    DecayCnt     int    `json:"decay_count" db:"decay_count"`
    LastUsedAt   *Timestamp `json:"last_used_at,omitempty"`
    
    // Metadata
    Connections []string `json:"connections,omitempty" db:"connections"`
    Tags        Tags     `json:"tags,omitempty" db:"tags"`
    Project     string   `json:"project,omitempty" db:"project"`
    UserID      string   `json:"user_id,omitempty" db:"user_id"`
    
    // Soft delete
    DeletedAt   *Timestamp `json:"deleted_at,omitempty"`
}

// NewPattern creates a new pattern with defaults
func NewPattern(trigger, response string) *Pattern {
    now := Now()
    return &Pattern{
        ID:           NewID(),
        CreatedAt:    now,
        UpdatedAt:    now,
        Trigger:      trigger,
        Response:     response,
        Strength:     0.0,
        Threshold:    50.0,
        DecayRate:    0.01,
        DecayEnabled: true,
        ReinforceCnt: 0,
        DecayCnt:     0,
    }
}
```

### 2.2 Pattern Validation

```go
// Validate checks if the pattern is valid
func (p *Pattern) Validate() error {
    if err := p.ValidateTrigger(); err != nil {
        return err
    }
    if err := p.ValidateResponse(); err != nil {
        return err
    }
    if err := p.ValidateStrength(); err != nil {
        return err
    }
    return nil
}

// ValidateTrigger validates the trigger field
func (p *Pattern) ValidateTrigger() error {
    trimmed := strings.TrimSpace(p.Trigger)
    if trimmed == "" {
        return ErrValidation{Field: "trigger", Reason: "cannot be empty"}
    }
    if len(trimmed) > 500 {
        return ErrValidation{Field: "trigger", Reason: "exceeds 500 characters"}
    }
    p.Trigger = trimmed
    return nil
}

// ValidateResponse validates the response field
func (p *Pattern) ValidateResponse() error {
    if strings.TrimSpace(p.Response) == "" {
        return ErrValidation{Field: "response", Reason: "cannot be empty"}
    }
    return nil
}

// ValidateStrength validates strength and threshold
func (p *Pattern) ValidateStrength() error {
    if p.Strength < 0 || p.Strength > 100 {
        return ErrValidation{Field: "strength", Reason: "must be between 0 and 100"}
    }
    if p.Threshold < 0 || p.Threshold > 100 {
        return ErrValidation{Field: "threshold", Reason: "must be between 0 and 100"}
    }
    if p.DecayRate < 0 || p.DecayRate > 1 {
        return ErrValidation{Field: "decay_rate", Reason: "must be between 0 and 1"}
    }
    return nil
}
```

### 2.3 Pattern JSON Representation

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2026-02-20T10:00:00Z",
  "updated_at": "2026-02-20T10:00:00Z",
  "trigger": "python setup",
  "response": "python3 -m venv venv && source venv/bin/activate && pip install -r requirements.txt",
  "strength": 75.5,
  "threshold": 50.0,
  "decay_rate": 0.01,
  "decay_enabled": true,
  "reinforcement_count": 5,
  "decay_count": 2,
  "last_used_at": "2026-02-20T09:30:00Z",
  "connections": ["python-venv", "python-pip"],
  "tags": ["python", "setup", "devops"],
  "project": "my-project"
}
```

---

## 3. Thought Chain Model

### 3.1 ThoughtChain Structure

```go
// ThoughtChain represents a chain of thoughts (branch hierarchy)
type ThoughtChain struct {
    ID        ID            `json:"id"`
    Root      *ThoughtNode  `json:"root"`
    Depth     int           `json:"depth"`
    CreatedAt Timestamp      `json:"created_at"`
}

// ThoughtNode represents a single thought in the chain
type ThoughtNode struct {
    ID          ID            `json:"id"`
    Pattern     *Pattern      `json:"pattern"`
    Confidence  Confidence    `json:"confidence"`
    Children    []*ThoughtNode `json:"children,omitempty"`
    Selected    bool          `json:"selected"`
    Expanded    bool          `json:"expanded"`
    Response    string        `json:"response,omitempty"`
}

// NewThoughtChain creates a new thought chain
func NewThoughtChain(rootPattern *Pattern, confidence Confidence) *ThoughtChain {
    return &ThoughtChain{
        ID: NewID(),
        Root: &ThoughtNode{
            ID:         NewID(),
            Pattern:    rootPattern,
            Confidence: confidence,
            Selected:   true,
            Expanded:   false,
        },
        Depth:     1,
        CreatedAt: Now(),
    }
}
```

### 3.2 Tree Operations

```go
// AddChild adds a child node to a thought
func (n *ThoughtNode) AddChild(child *ThoughtNode) {
    n.Children = append(n.Children, child)
    n.Expanded = true
}

// FindByID finds a node by ID in the tree
func (tc *ThoughtChain) FindByID(id ID) *ThoughtNode {
    return tc.findByID(tc.Root, id)
}

func (n *ThoughtNode) findByID(id ID) *ThoughtNode {
    if n.ID == id {
        return n
    }
    for _, child := range n.Children {
        if found := child.findByID(id); found != nil {
            return found
        }
    }
    return nil
}

// SelectPath selects a path from root to a specific node
func (tc *ThoughtChain) SelectPath(targetID ID) bool {
    path := tc.FindPath(tc.Root.ID, targetID)
    if path == nil {
        return false
    }
    
    // Deselect all
    tc.deselectAll(tc.Root)
    
    // Select nodes on path
    for _, node := range path {
        node.Selected = true
        node.Expanded = true
    }
    
    return true
}

func (n *ThoughtNode) findByID(targetID ID) *ThoughtNode {
    if n.ID == targetID {
        return n
    }
    for _, child := range n.Children {
        if found := child.findByID(targetID); found != nil {
            return found
        }
    }
    return nil
}

func (tc *ThoughtChain) FindPath(fromID, toID ID) []*ThoughtNode {
    // BFS to find path
    return nil // Implementation...
}
```

---

## 4. Space/Project Model

### 4.1 Space Structure

```go
// Space represents an isolated namespace for patterns
type Space struct {
    ID          ID        `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description,omitempty" db:"description"`
    CreatedAt   Timestamp `json:"created_at" db:"created_at"`
    UpdatedAt   Timestamp `json:"updated_at" db:"updated_at"`
    
    // Settings
    DefaultSpace bool     `json:"default_space"`
    PatternLimit int      `json:"pattern_limit"` // 0 = unlimited
    
    // Statistics
    PatternCount int `json:"pattern_count"`
}

// DefaultSpaces returns the default spaces
func DefaultSpaces() []*Space {
    return []*Space{
        {
            ID:   "global",
            Name: "Global",
            Description: "Patterns available everywhere",
            DefaultSpace: true,
        },
        {
            ID:   "project",
            Name: "Project",
            Description: "Project-specific patterns",
        },
        {
            ID:   "personal",
            Name: "Personal",
            Description: "Personal patterns",
        },
    }
}
```

---

## 5. Event Model

### 5.1 Event Types

```go
// EventType represents the type of event
type EventType string

const (
    EventPatternCreated   EventType = "pattern.created"
    EventPatternMatched   EventType = "pattern.matched"
    EventPatternReached  EventType = "pattern.threshold_reached"
    EventPatternDecayed  EventType = "pattern.decayed"
    EventPatternDeleted  EventType = "pattern.deleted"
    EventPatternUpdated  EventType = "pattern.updated"
    EventSpaceCreated    EventType = "space.created"
    EventSpaceDeleted    EventType = "space.deleted"
    EventAIRequest       EventType = "ai.request"
    EventAIResponse      EventType = "ai.response"
    EventError           EventType = "error.occurred"
)
```

### 5.2 Event Structure

```go
// Event represents a domain event
type Event struct {
    ID        ID        `json:"id"`
    Type      EventType `json:"type"`
    Timestamp Timestamp `json:"timestamp"`
    
    // Event data
    Payload   interface{} `json:"payload"`
    
    // Metadata
    Source    string     `json:"source"`
    TraceID   string     `json:"trace_id,omitempty"`
    UserID    string     `json:"user_id,omitempty"`
}
```

### 5.3 Event Payloads

```go
// PatternCreatedPayload is the payload for pattern created events
type PatternCreatedPayload struct {
    PatternID  string `json:"pattern_id"`
    Trigger    string `json:"trigger"`
    Project    string `json:"project,omitempty"`
}

// PatternMatchedPayload is the payload for pattern matched events
type PatternMatchedPayload struct {
    PatternID   string  `json:"pattern_id"`
    Query       string  `json:"query"`
    Confidence  float64 `json:"confidence"`
    BranchTaken string  `json:"branch_taken"`
}

// PatternDecayedPayload is the payload for pattern decayed events
type PatternDecayedPayload struct {
    PatternID   string  `json:"pattern_id"`
    OldStrength float64 `json:"old_strength"`
    NewStrength float64 `json:"new_strength"`
    DecayAmount float64 `json:"decay_amount"`
}

// AIRequestPayload is the payload for AI request events
type AIRequestPayload struct {
    RequestID   string `json:"request_id"`
    Provider    string `json:"provider"`
    Model       string `json:"model"`
    PromptLen   int    `json:"prompt_length"`
    StartedAt   time.Time `json:"started_at"`
}
```

---

## 6. Configuration Model

### 6.1 Config Structure

```go
// Config represents the application configuration
type Config struct {
    // Version for config migration
    Version int `yaml:"version"`
    
    // Application settings
    App AppConfig `yaml:"app"`
    
    // Storage settings
    Storage StorageConfig `yaml:"storage"`
    
    // AI provider settings
    AI AIConfig `yaml:"ai"`
    
    // UI settings
    UI UIConfig `yaml:"ui"`
    
    // Security settings
    Security SecurityConfig `yaml:"security"`
}

// AppConfig contains application-level settings
type AppConfig struct {
    Name        string `yaml:"name"`
    Version     string `yaml:"version"`
    DataDir     string `yaml:"data_dir"`
    LogLevel    string `yaml:"log_level"`
    Profile     bool   `yaml:"profile"` // Enable profiling
}

// StorageConfig contains storage-related settings
type StorageConfig struct {
    Type     string `yaml:"type"` // "sqlite", "memory"
    Path     string `yaml:"path"`
    CacheSize int   `yaml:"cache_size"`
}

// AIConfig contains AI provider settings
type AIConfig struct {
    Provider  string       `yaml:"provider"`
    Providers ProvidersConfig `yaml:"providers"`
    DefaultModel string   `yaml:"default_model"`
    Timeout   int         `yaml:"timeout"` // seconds
    RetryMax  int         `yaml:"retry_max"`
}

// ProvidersConfig contains per-provider settings
type ProvidersConfig struct {
    Anthropic AnthropicConfig `yaml:"anthropic"`
    OpenAI    OpenAIPConfig   `yaml:"openai"`
    Local     LocalConfig     `yaml:"local"`
}

// AnthropicConfig contains Anthropic-specific settings
type AnthropicConfig struct {
    APIKey     string  `yaml:"api_key"`
    APIURL     string  `yaml:"api_url"`
    Model      string  `yaml:"model"`
    MaxTokens  int     `yaml:"max_tokens"`
    Temperature float64 `yaml:"temperature"`
}

// UIConfig contains UI-related settings
type UIConfig struct {
    Theme      string   `yaml:"theme"`
    Colors     ColorsConfig `yaml:"colors"`
    KeyMap     KeyMapConfig  `yaml:"keymap"`
    OutputMode string  `yaml:"output_mode"` // "terminal", "json", "plain"
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
    APIKeysEnvPrefix string `yaml:"api_keys_env_prefix"`
    ConfigFileMode   string `yaml:"config_file_mode"`
    AuditLog         AuditConfig `yaml:"audit_log"`
}
```

### 6.2 Default Configuration

```yaml
version: 1

app:
  name: "open-think-reflex"
  data_dir: "~/.openclaw/reflex"
  log_level: "info"
  profile: false

storage:
  type: "sqlite"
  path: "~/.openclaw/reflex/data.db"
  cache_size: 1000

ai:
  provider: "anthropic"
  default_model: "claude-sonnet-4-20250514"
  timeout: 30
  retry_max: 3
  providers:
    anthropic:
      api_url: "https://api.anthropic.com/v1"
      max_tokens: 4096
      temperature: 0.7
    openai:
      api_url: "https://api.openai.com/v1"
      model: "gpt-4"
      max_tokens: 4096
      temperature: 0.7
    local:
      api_url: "http://localhost:11434/v1"
      model: "llama2"

ui:
  theme: "dark"
  output_mode: "terminal"
  keymap:
    up: "k"
    down: "j"
    left: "h"
    right: "l"
    select: "space"
    confirm: "enter"
    cancel: "esc"
    quit: "q"
    help: "?"

security:
  api_keys_env_prefix: "OTR_"
  config_file_mode: "0600"
  audit_log:
    enabled: true
    path: "~/.openclaw/reflex/audit.log"
```

---

## 7. JSON Schemas

### 7.1 Pattern JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://open-think-reflex.schema/pattern/v1",
  "title": "Pattern",
  "description": "A reflex pattern in the Open-Think-Reflex system",
  "type": "object",
  "required": ["id", "trigger", "response", "strength", "threshold"],
  "properties": {
    "id": {
      "type": "string",
      "format": "uuid",
      "description": "Unique identifier for the pattern"
    },
    "trigger": {
      "type": "string",
      "minLength": 1,
      "maxLength": 500,
      "description": "The trigger phrase that activates this pattern"
    },
    "response": {
      "type": "string",
      "minLength": 1,
      "description": "The response content for this pattern"
    },
    "strength": {
      "type": "number",
      "minimum": 0,
      "maximum": 100,
      "description": "Current strength of the pattern (0-100)"
    },
    "threshold": {
      "type": "number",
      "minimum": 0,
      "maximum": 100,
      "description": "Threshold for activation (0-100)"
    },
    "decay_rate": {
      "type": "number",
      "minimum": 0,
      "maximum": 1,
      "default": 0.01
    },
    "decay_enabled": {
      "type": "boolean",
      "default": true
    },
    "connections": {
      "type": "array",
      "items": { "type": "string" }
    },
    "tags": {
      "type": "array",
      "items": { "type": "string" }
    },
    "project": {
      "type": "string"
    }
  }
}
```

---

## 8. Database Schema

### 8.1 SQLite Schema

```sql
-- Patterns table
CREATE TABLE IF NOT EXISTS patterns (
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
    last_used_at INTEGER,
    tags TEXT,
    project TEXT,
    user_id TEXT,
    deleted_at INTEGER
);

-- Spaces table
CREATE TABLE IF NOT EXISTS spaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    is_default INTEGER NOT NULL DEFAULT 0,
    pattern_limit INTEGER NOT NULL DEFAULT 0,
    pattern_count INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- Events table (for audit/logging)
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    payload TEXT,
    source TEXT,
    trace_id TEXT,
    pattern_id TEXT,
    user_id TEXT
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_patterns_trigger ON patterns(trigger);
CREATE INDEX IF NOT EXISTS idx_patterns_strength ON patterns(strength);
CREATE INDEX IF NOT EXISTS idx_patterns_project ON patterns(project);
CREATE INDEX IF NOT EXISTS idx_patterns_tags ON patterns(tags);
CREATE INDEX IF NOT EXISTS idx_patterns_deleted ON patterns(deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
```

### 8.2 Migrations

```go
// migrations.go

// CurrentSchemaVersion is the current schema version
const CurrentSchemaVersion = 1

// Migrations contains all schema migrations
var Migrations = []Migration{
    {
        Version: 1,
        SQL: `
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
                last_used_at INTEGER,
                tags TEXT,
                project TEXT,
                user_id TEXT,
                deleted_at INTEGER
            );
            
            CREATE TABLE spaces (
                id TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                description TEXT,
                is_default INTEGER NOT NULL DEFAULT 0,
                pattern_limit INTEGER NOT NULL DEFAULT 0,
                pattern_count INTEGER NOT NULL DEFAULT 0,
                created_at INTEGER NOT NULL,
                updated_at INTEGER NOT NULL
            );
            
            CREATE INDEX idx_patterns_trigger ON patterns(trigger);
            CREATE INDEX idx_patterns_strength ON patterns(strength);
            CREATE INDEX idx_patterns_project ON patterns(project);
        `,
    },
}
```

---

**Document Version**: v1.0  
**Created**: 2026-02-20  
**Project**: open-think-reflex
