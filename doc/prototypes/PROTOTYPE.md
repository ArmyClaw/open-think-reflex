# Open-Think-Reflex Prototype Design

> **Version**: v1.0  
> **Status**: Draft  
> **Scope**: CLI Interaction Prototype Design for v1.0

---

## Table of Contents

1. [Design Principles](#1-design-principles)
2. [CLI Design References](#2-cli-design-references)
3. [Three-Layer Interface](#3-three-layer-interface)
4. [Interaction Patterns](#4-interaction-patterns)
5. [Output Format](#5-output-format)
6. [Keyboard Shortcuts](#6-keyboard-shortcuts)
7. [Error Handling](#7-error-handling)
8. [Visual Design](#8-visual-design)
9. [Animation & Transitions](#9-animation--transitions)
10. [Accessibility](#10-accessibility)
11. [Prototype Specifications](#11-prototype-specifications)

---

## 1. Design Principles

### 1.1 Core Philosophy

Inspired by industry-leading CLI tools, Open-Think-Reflex follows these principles:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Open-Think-Reflex Design Philosophy                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. DO ONE THING WELL                                                  │
│     Focus on reflex-based thinking augmentation                          │
│     Don't try to be a general-purpose AI CLI                            │
│                                                                      │
│  2. COMPOSABLE COMMANDS                                                │
│     Simple inputs → predictable outputs                                  │
│     Chainable operations                                                │
│     Consistent flag patterns                                            │
│                                                                      │
│  3. PROGRESSIVE DISCOVERY                                             │
│     Simple commands for beginners                                        │
│     Advanced flags for power users                                      │
│     Helpful suggestions for corrections                                  │
│                                                                      │
│  4. CONTEXTUAL AWARENESS                                              │
│     Remember previous interactions                                       │
│     Adapt to user patterns                                              │
│     Provide relevant suggestions                                         │
│                                                                      │
│  5. GRACEFUL DEGRADATION                                              │
│     Work offline with limited features                                   │
│     Clear error messages                                                │
│     Fallback options for failures                                       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Design Inspirations

| Tool | Inspiration | Application |
|------|-------------|-------------|
| **git** | Subcommand pattern (`git <verb> <noun>`) | `otr pattern <verb> <noun>` |
| **kubectl** | Resource-oriented commands | `otr pattern create --trigger="xxx"` |
| **docker** | Container management simplicity | `otr <command> --output format` |
| **npm** | Script-friendly outputs | JSON output for automation |
| **htop** | Visual feedback | Progress bars, status indicators |
| **fzf** | Fuzzy matching | Quick selection interfaces |
| **exa/ls** | Color-coded outputs | File/branch highlighting |

---

## 2. CLI Design References

### 2.1 Command Structure Pattern

Inspired by git's successful pattern:

```
git [global options] <command> [<args>] [<flags>]
```

Open-Think-Reflex follows:

```
otr [global options] <command> [<args>] [<flags>]
```

### 2.2 Command Organization

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Command Categories                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  CORE COMMANDS                                                         │
│  ├── otr "query"                    # Main interaction                │
│  ├── otr <trigger>                  # Quick reflex execution          │
│  └── otr run <pattern>              # Execute specific pattern        │
│                                                                      │
│  PATTERN MANAGEMENT                                                    │
│  ├── otr pattern create             # Create new reflex               │
│  ├── otr pattern list               # List all reflexes               │
│  ├── otr pattern show <id>          # View reflex details             │
│  ├── otr pattern edit <id>          # Edit reflex                    │
│  ├── otr pattern delete <id>        # Delete reflex                  │
│  └── otr pattern reinforce <id>     # Manually reinforce             │
│                                                                      │
│  SPACE/PROJECT MANAGEMENT                                              │
│  ├── otr space list                 # List spaces                     │
│  ├── otr space use <name>          # Switch space                   │
│  ├── otr project list               # List projects                  │
│  ├── otr project create <name>     # Create project (v2.0)          │
│  └── otr project use <name>        # Switch project (v2.0)          │
│                                                                      │
│  EXPORT/IMPORT                                                        │
│  ├── otr export [json|sqlite]      # Export data                    │
│  └── otr import <file>              # Import data                     │
│                                                                      │
│  SYSTEM COMMANDS                                                       │
│  ├── otr init                       # Initial setup                  │
│  ├── otr config [show|set|reset]   # Configuration                  │
│  ├── otr stats                      # View statistics                 │
│  ├── otr health                     # Health check                    │
│  ├── otr backup                     # Create backup                  │
│  ├── otr decay [run|status]        # Decay management               │
│  └── otr --version / --help         # Global options                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Flag Design Pattern

Consistent with kubectl/docker:

| Flag Type | Example | Description |
|-----------|---------|-------------|
| **Global** | `--output json` | Affects all commands |
| **Command** | `--force` | Command-specific |
| **Short** | `-f` | Short form |
| **Boolean** | `--enabled=true` | Boolean flags |
| **Value** | `--threshold=50` | Value flags |
| **Array** | `--tags=api,database` | Comma-separated |

### 2.4 Exit Codes

Standard exit codes for automation:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Exit Code Standards                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  0   │ Success                                                        │
│  1   │ General error                                                  │
│  2   │ Invalid arguments                                               │
│  3   │ Resource not found (pattern, space, project)                   │
│  4   │ Permission denied                                               │
│  5   │ Configuration error                                             │
│  6   │ Network/AI service unavailable                                   │
│  7   │ Cancellation by user                                            │
│  8   │ Validation error                                                │
│  9   │ Deprecated command/option                                       │
│  130  │ Ctrl+C cancellation (standard Unix)                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Three-Layer Interface

### 3.1 Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Three-Layer Interface                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  LAYER 1: THOUGHT CHAIN (Top Layer)                          │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │   │
│  │  │  Root   │ ─►│ Branch A│ ─►│ Sub-A1 │ ─►│ Sub-A1-a│        │   │
│  │  │  85%    │   │  78%    │   │   72%  │   │   65%   │        │   │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘        │   │
│  │                                                                  │   │
│  │  Visual States:                                                  │   │
│  │  ├── Selected: Bold + Green border                              │   │
│  │  ├── Matched: Green text                                        │   │
│  │  ├── Unmatched: Gray text                                      │   │
│  │  └── Confirmed: ✓ marker                                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  LAYER 2: OUTPUT (Middle Layer)                                │   │
│  │  ┌─────────────────────────────────────────────────────────┐    │   │
│  │  │                                                         │    │   │
│  │  │              AI-generated content here...               │    │   │
│  │  │                                                         │    │   │
│  │  │  # Market Trend Analysis Framework                     │    │   │
│  │  │  ## 1. Data Sources                                   │    │   │
│  │  │  ## 2. Pattern Identification                          │    │   │
│  │  │                                                         │    │   │
│  │  └─────────────────────────────────────────────────────────┘    │   │
│  │                                                                  │   │
│  │  Scroll indicator when content overflows                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  LAYER 3: INPUT (Bottom Layer)                                │   │
│  │  ┌─────────────────────────────────────────────────────────┐    │   │
│  │  │ > _ [↓↑ select, → expand, Space execute, h help]    │    │   │
│  │  └─────────────────────────────────────────────────────────┘    │   │
│  │                                                                  │   │
│  │  Auto-suggestions based on history                             │   │
│  │  Syntax highlighting for triggers                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Layer Specifications

#### Layer 1: Thought Chain (Top)

```
Dimensions:
- Height: 4 rows maximum
- Width: 60-80% of terminal
- Position: Top of screen

Elements:
- Root node (user input/trigger)
- Branch nodes (matching reflexes)
- Sub-branches (refinement options)
- Confidence scores (0-100%)
- Selection indicator (◉/✓)

Rendering:
- Horizontal tree: Root ─► Branch ─► Sub-branch
- Vertical space: 2-4 lines per node level
- Overflow: Scrollable horizontally
```

#### Layer 2: Output (Middle)

```
Dimensions:
- Height: Auto (fit content, max 50% terminal)
- Width: Full terminal width

Elements:
- AI-generated content
- Headers (#, ##, ###)
- Lists (-, *, 1.)
- Code blocks (```)
- Tables (ASCII)
- Progress indicators

Rendering:
- Markdown-like formatting
- Syntax highlighting for code
- Word wrap at terminal width
- Pagination for long content
```

#### Layer 3: Input (Bottom)

```
Dimensions:
- Height: 2 rows
- Position: Bottom of screen

Elements:
- Prompt (> )
- User input area
- Help text (optional)
- Auto-suggestions

Rendering:
- Single-line input
- Cursor blinking
- History navigation (↑/↓)
- Tab completion
```

### 3.3 Screen Layout Example

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ otr ─ AI Thought Reflex Accelerator                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [root: market-analysis] ─► [Data Collection] ✓ 85%                       │
│                          ─► [Pattern ID] ◉ 78%                            │
│                          ─► [Risk Assessment] 72%                          │
│                          ─► [Recommendations] 65%                           │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  # Market Trend Analysis Framework                                          │
│                                                                             │
│  ## 1. Data Sources                                                        │
│     - Sales data from CRM                                                 │
│     - Customer feedback surveys                                            │
│     - Competitor pricing data                                             │
│     - Social media sentiment                                              │
│                                                                             │
│  ## 2. Pattern Identification                                             │
│     - Trend direction: Upward (+15% QoQ)                                  │
│     - Key driver: Mobile adoption                                         │
│     - Seasonality: Q4 spike expected                                      │
│                                                                             │
│  [Scroll: 2/5 ──████────]                                                │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  > market-analysis [↓↑ select, → expand, Space execute, h help]        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Interaction Patterns

### 4.1 Basic Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Interaction Flow                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  User Input                    System Response                          │
│  ────────────                  ────────────────                        │
│                                                                      │
│  $ otr "query"                ├─ Check matching patterns              │
│                             ├─ If match: Show thought chain           │
│                             └─ If no match: Prompt for AI generation  │
│                                                                      │
│  [↑/↓]                       ├─ Navigate branches                     │
│                             └─ Highlight selection                      │
│                                                                      │
│  [→]                         ├─ Expand current branch                  │
│                             └─ Show sub-branches                       │
│                                                                      │
│  [Space]                     ├─ Execute selected branch               │
│                             ├─ Copy output to clipboard               │
│                             └─ Reinforce pattern strength             │
│                                                                      │
│  [Enter]                     ├─ If no match: Confirm AI generation   │
│                             └─ Copy to clipboard                       │
│                                                                      │
│  [Esc]                       ├─ Cancel current operation              │
│                             └─ Return to previous state               │
│                                                                      │
│  [q]                         ├─ Quit interactive mode                 │
│                             └─ Return to shell                        │
│                                                                      │
│  [Ctrl+C]                    ├─ Force quit                           │
│                             └─ Exit with code 130                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Selection Patterns

#### Single Branch Selection

```
State: User has selected one branch

Visual:
┌─────────────────────────────────────────┐
│ [market-analysis] ─► [Data Collection] ✓ 85% │
│                    ─► [Pattern ID] ◉ 78%  │
│                    ─► [Risk Assessment] 72% │
└─────────────────────────────────────────┘

Action: [Space] to execute selected
```

#### Multi-Branch Expansion

```
State: User expands [Pattern ID] branch

Visual:
┌─────────────────────────────────────────┐
│ [market-analysis] ─► [Data Collection] ✓ 85% │
│                    ─► [Pattern ID] ◉ 78%  │
│                          ─► [Sales] 65%  │
│                          ─► [Customer] 62%│
│                          ─► [Competitor]58%│
│                    ─► [Risk Assessment] 72% │
└─────────────────────────────────────────┘

Action: [↑/↓] to navigate sub-branches
```

### 4.3 Command-Line Mode

Non-interactive output for scripting:

```bash
# JSON output for programmatic use
$ otr "market-analysis" --output json --no-interactive

{"trigger":"market-analysis","response":"# Market Trend Analysis...","strength":85}

# Quiet mode for pipelines
$ otr "swot" --quiet | jq '.response' > swot.json

# Specific pattern execution
$ otr run market-analysis --output markdown
```

---

## 5. Output Format

### 5.1 Default (Terminal)

Human-readable with colors and formatting:

```
$ otr "api design principles"

# API Design Principles

## RESTful Principles
- Resource-based URLs
- HTTP methods as actions
- Stateless interactions

## Best Practices
- Versioning: /api/v1/
- Error handling: Consistent formats
- Pagination: ?page=&limit=
```

### 5.2 JSON Mode

For programmatic use:

```json
{
  "trigger": "api design principles",
  "response": "# API Design Principles\n\n## RESTful Principles\n...",
  "strength": 85,
  "branches": [
    {
      "trigger": "RESTful Principles",
      "strength": 78,
      "selected": true
    }
  ],
  "metadata": {
    "createdAt": "2024-01-01T00:00:00Z",
    "lastUsedAt": "2024-01-15T00:00:00Z",
    "tags": ["api", "design"]
  }
}
```

### 5.3 Plain Mode

Stripped output for copy/paste:

```
$ otr "api design" --output plain

# API Design Principles

## RESTful Principles
- Resource-based URLs
- HTTP methods as actions

## Best Practices
- Versioning: /api/v1/
```

### 5.4 Output Options

| Option | Alias | Description |
|--------|-------|-------------|
| `--output` | `-o` | Format: terminal, json, plain, markdown |
| `--no-color` | | Disable color output |
| `--pager` | | Use pager for long output |
| `--width` | `-w` | Set output width (default: auto) |

---

## 6. Keyboard Shortcuts

### 6.1 Shortcut Reference

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Keyboard Shortcuts Reference                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  NAVIGATION                                                            │
│  ──────────                                                            │
│  ↑ / k     Move selection up                                           │
│  ↓ / j     Move selection down                                         │
│  → / l     Expand current branch                                       │
│  ← / h     Go back / Return to parent                                 │
│  gg        Jump to first option                                        │
│  G         Jump to last option                                         │
│                                                                      │
│  EXECUTION                                                            │
│  ──────────                                                            │
│  Space      Execute selected branch                                    │
│  Enter      Confirm / Copy to clipboard                                │
│  Ctrl+Enter Copy full output to clipboard                              │
│                                                                      │
│  CANCELLATION                                                          │
│  ──────────────                                                        │
│  Esc         Cancel current operation                                  │
│  q           Quit interactive mode                                     │
│  Ctrl+C      Force quit (with confirmation)                           │
│                                                                      │
│  HELP & NAVIGATION                                                    │
│  ────────────────                                                      │
│  h / ?       Show help                                                │
│  /           Search within output                                      │
│  n           Next search result                                       │
│  N           Previous search result                                    │
│                                                                      │
│  SPECIAL                                                              │
│  ────────                                                              │
│  Ctrl+L      Refresh display                                          │
│  Ctrl+R      Toggle dark/light theme                                  │
│  Ctrl+S      Save current output to file                              │
│  Ctrl+D      Toggle debug mode                                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Shortcut Design Rationale

| Shortcut | Reason | Status |
|----------|--------|--------|
| ↑/↓ | Universal navigation (vim, fzf) | Required |
| →/← | Tree expansion/contraction | Required |
| Tab | Switch between layers | Required |
| Space | Execute action (game UIs, terminals) | Required |
| Esc | Cancel (universal) | Required |
| q | Quit (vim, less, git) | Required |
| h | Help (vim, git) | Required |
| / | Search (vim, less, git) | Required |
| Ctrl+C | Force interrupt (Unix standard) | Required |

**Optional (Future Enhancement)**:
- Ctrl+L: Refresh display
- Ctrl+R: Toggle theme
- Ctrl+S: Save to file
- Ctrl+D: Debug mode

---

## 7. Error Handling

### 7.1 Error Display

Consistent error format:

```
╔═══════════════════════════════════════════════════════════════════════════╗
║  ERROR [P-003]: Pattern not found                                      ║
║                                                                           ║
║  Pattern "market-analysis" does not exist.                               ║
║                                                                           ║
║  Suggestions:                                                            ║
║    • Run "otr pattern list" to view all patterns                        ║
║    • Create new: "otr pattern create"                                    ║
║    • Search: "otr list | grep -i pattern"                               ║
║                                                                           ║
║  Type: retrieval_error                                                   ║
║  Code: P-003                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 7.2 Error Categories

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Error Categories                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  PATTERN ERRORS (P-xxx)                                                │
│  P-001 │ Pattern ID not provided                                        │
│  P-002 │ Pattern not found                                             │
│  P-003 │ Pattern already exists                                         │
│  P-004 │ Pattern validation failed                                      │
│  P-005 │ Pattern update failed                                         │
│                                                                      │
│  MATCHING ERRORS (M-xxx)                                               │
│  M-001 │ No patterns match query                                       │
│  M-002 │ Match timeout                                                │
│  M-003 │ Ambiguous match (multiple patterns)                           │
│                                                                      │
│  STORAGE ERRORS (S-xxx)                                               │
│  S-001 │ Database connection failed                                    │
│  S-002 │ Storage quota exceeded                                        │
│  S-003 │ Backup failed                                                │
│  S-004 │ Migration failed                                             │
│                                                                      │
│  AI SERVICE ERRORS (A-xxx)                                             │
│  A-001 │ AI service unavailable                                        │
│  A-002 │ AI request timeout                                           │
│  A-003 │ AI quota exceeded                                            │
│  A-004 │ AI response invalid                                          │
│                                                                      │
│  CONFIGURATION ERRORS (C-xxx)                                          │
│  C-001 │ Config file not found                                        │
│  C-002 │ Config validation failed                                     │
│  C-003 │ Invalid option value                                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Recovery Suggestions

Each error includes actionable suggestions:

| Error Type | Suggestion Pattern |
|------------|-------------------|
| Pattern not found | `• Run "otr pattern list" to view all patterns` |
| No match | `• Try more specific query` / `• Create new pattern` |
| Service unavailable | `• Check network connection` / `• Use offline mode` |
| Invalid input | `• Run with --dry-run to validate` |

---

## 8. Visual Design

### 8.1 Color Scheme

#### Light Theme

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Light Theme Colors                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Element          │ Foreground    │ Background    │ Style              │
│  ─────────────────┼───────────────┼──────────────┼────────────────    │
│  Prompt           │ Gray (#6c7373) │ Default      │ Bold               │
│  Input            │ Black (#1f2327)│ Default      │                    │
│  Selected branch  │ Green (#1a7f4d)│ Light (#ecfdf5)│ Bold + Border │
│  Matched branch   │ Green (#1a7f4d)│ Default      │                    │
│  Unmatched branch │ Gray (#9ca3af) │ Default      │                    │
│  Header           │ Blue (#1d4ed8) │ Default      │ Bold               │
│  Code block       │ Purple (#7c3aed)│ Light (#f5f3ff)│ Monospace      │
│  Error            │ Red (#dc2626)  │ Light (#fef2f2)│ Bold               │
│  Warning          │ Orange (#d97706)| Light (#fffbeb)│ Bold              │
│  Success          │ Green (#059669) │ Light (#ecfdf5)│ Bold              │
│  Info             │ Blue (#2563eb)  │ Light (#eff6ff)│ Bold              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Dark Theme

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Dark Theme Colors                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Element          │ Foreground    │ Background    │ Style              │
│  ─────────────────┼───────────────┼──────────────┼────────────────    │
│  Prompt           │ Gray (#9ca3af)│ Default      │ Bold               │
│  Input            │ White (#f9fafb)| Default     │                    │
│  Selected branch  │ Green (#34d399)│ Dark (#064e3b)│ Bold + Border │
│  Matched branch   │ Green (#34d399)│ Default      │                    │
│  Unmatched branch │ Gray (#6b7280) │ Default      │                    │
│  Header           │ Blue (#60a5fa) │ Default      │ Bold               │
│  Code block       │ Purple (#a78bfa)| Dark (#1e1b4b)│ Monospace      │
│  Error            │ Red (#f87171)  │ Dark (#7f1d1d)│ Bold               │
│  Warning          │ Orange (#fbbf24)| Dark (#713f12)│ Bold              │
│  Success          │ Green (#34d399) │ Dark (#064e3b)│ Bold              │
│  Info             │ Blue (#60a5fa) │ Dark (#1e3a8a)│ Bold              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Typography

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Typography Standards                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Headings (Markdown)                                                  │
│  #        │ 24px │ Bold     │ Terminal: ════                       │
│  ##       │ 20px │ Bold     │ Terminal: ────                        │
│  ###      │ 16px │ Bold     │ Terminal: ····                        │
│                                                                      │
│  Body Text                                                            │
│  Regular │ 14px │ Regular  │ Main content                             │
│  Code    │ 13px │ Monospace│ Code blocks, triggers                   │
│  Caption │ 12px │ Gray    │ Help text, suggestions                   │
│                                                                      │
│  Progress Bars                                                        │
│  Complete │ ████████░░ │ 8 segments │ 100-88%                         │
│  Partial  │ ██████░░░░ │ 8 segments │ 87-50%                          │
│  Low     │ ████░░░░░░ │ 8 segments │ 49-25%                          │
│  Critical│ ██░░░░░░░░ │ 8 segments │ 24-0%                           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Spacing

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Spacing Standards                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Layer Padding:                                                        │
│  ──────────────                                                        │
│  Top layer:    1 line                                                 │
│  Between lines: 0 lines (compact)                                      │
│  Bottom layer: 1 line                                                 │
│                                                                      │
│  Horizontal Margins:                                                  │
│  ────────────────────                                                  │
│  Default:     2 spaces from edge                                       │
│  Code blocks: 4 spaces indent                                         │
│  Nested:      4 spaces per level                                      │
│                                                                      │
│  Vertical Spacing:                                                    │
│  ─────────────────                                                    │
│  Sections:    1 empty line                                             │
│  Subsections: 0 empty lines                                           │
│  Lists:       0 empty lines between items                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 9. Animation & Transitions

### 9.1 Animation Principles

Inspired by fast, responsive CLIs like fzf and eza:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Animation Guidelines                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. SPEED FIRST                                                        │
│     Animations should be under 100ms                                    │
│     If >100ms, skip animation                                          │
│     Default to instant for most operations                             │
│                                                                      │
│  2. MEANINGFUL MOMENTS                                                 │
│     Only animate meaningful state changes                              │
│     • Pattern matching: Instant                                        │
│     • Branch expansion: Quick slide (50ms)                            │
│     • Loading states: Progress bar                                     │
│     • Success/Fail: Brief flash (100ms)                               │
│                                                                      │
│  3. NON-BLOCKING                                                       │
│     Animations should never block input                                 │
│     User can interrupt anytime                                         │
│                                                                      │
│  4. REDUCED MOTION                                                    │
│     Respect system `NO_MOTION` preference                               │
│     Disable animations in low-power mode                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Specific Animations

#### Loading States

```
Searching...
[████████░░░░░░░░░░] 45% • 0.3s
```

#### Progress Indicator

```
Exporting patterns...
[██████████] 100% • Complete
```

#### Branch Expansion

```
Before:  [root] ─► [branch] ─► ?
After:   [root] ─► [branch] ─► [sub-a] ✓
                                 ─► [sub-b]
                                 ─► [sub-c]
```

#### Success/Fail

```
✓ Pattern created successfully
✗ Failed to export: Network error
```

### 9.3 Transition Timing

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Transition Timing                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Operation                    │ Duration  │ Animation Type               │
│  ────────────────────────────┼───────────┼─────────────────            │
│  Branch selection             │ 0ms       │ Instant                     │
│  Branch expansion             │ 50ms      │ Quick slide                │
│  Content scroll               │ 0ms       │ Instant                     │
│  Progress bar update          │ 100ms     │ Smooth update              │
│  Error message fade-in        │ 150ms     │ Fade in                    │
│  Success flash                │ 100ms     │ Flash                       │
│  Screen refresh               │ 0ms       │ Instant                     │
│  Theme switch                 │ 200ms     │ Fade                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 10. Accessibility

### 10.1 Screen Reader Support

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Accessibility Features                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. SEMANTIC ANNOUNCEMENTS                                             │
│     Screen readers announce:                                            │
│     • Layer changes                                                    │
│     • Selection changes                                                │
│     • Error messages                                                   │
│     • Progress states                                                  │
│                                                                      │
│  2. ARIA LIVE REGIONS                                                 │
│     Dynamic content announced automatically                             │
│     Changes during typing announced                                    │
│     Suggestions announced                                              │
│                                                                      │
│  3. KEYBOARD FOCUS                                                    │
│     Clear focus indicators                                             │
│     Focus order: Input → Layers → Output                              │
│     Focus visible at all times                                         │
│                                                                      │
│  4. HIGH CONTRAST                                                     │
│     Meet WCAG AA standards                                             │
│     Minimum 4.5:1 contrast ratio                                      │
│     Bold text for emphasis                                            │
│                                                                      │
│  5. REDUCED MOTION                                                    │
│     Respect `prefers-reduced-motion`                                   │
│     Static alternatives to animations                                  │
│     Disable on request                                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Keyboard Accessibility

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Keyboard Accessibility                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  All functionality available via keyboard                              │
│                                                                      │
│  Tab Navigation:                                                      │
│  ───────────────                                                      │
│  Tab      │ Next focusable element                                     │
│  Shift+Tab│ Previous focusable element                                  │
│                                                                      │
│  Standard Shortcuts:                                                    │
│  ─────────────────                                                    │
│  Ctrl+A    │ Select all                                                │
│  Ctrl+C    │ Copy (with clipboard support)                             │
│  Ctrl+V    │ Paste                                                    │
│  Ctrl+Z    │ Undo                                                     │
│                                                                      │
│  Emergency:                                                           │
│  ────────                                                             │
│  Ctrl+Alt+D│ Disable animations (accessibility)                        │
│  Ctrl+Alt+H│ Toggle high contrast                                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 11. Prototype Specifications

### 11.1 Component Library

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    UI Component Library                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  BUTTONS                                                               │
│  ├─ Primary   │ [Space] Execute    │ Green background                  │
│  ├─ Secondary │ [Enter] Confirm     │ Gray background                  │
│  ├─ Danger   │ [Delete] Remove     │ Red background                  │
│  └─ Ghost    │ [?] Help           │ Transparent                       │
│                                                                      │
│  INPUT FIELDS                                                          │
│  ├─ Prompt   │ > user input        │ Always visible                  │
│  ├─ Search   │ / pattern search     │ With history                    │
│  └─ Command  │ otr command entry   │ With completion                │
│                                                                      │
│  LISTS                                                                 │
│  ├─ Branch   │ Selection options    │ Hover + select states           │
│  ├─ Pattern  │ Pattern list        │ Scrollable                       │
│  └─ Suggest  │ Auto-complete       │ Inline                          │
│                                                                      │
│  OVERLAYS                                                              │
│  ├─ Modal    │ Confirmations        │ Center screen                   │
│  ├─ Tooltip  │ Help text           │ Near cursor                     │
│  ├─ Toast    │ Notifications       │ Bottom right                    │
│  └─ Loading  │ Progress            │ Full overlay                    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 11.2 Terminal Compatibility

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Terminal Support                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  MINIMUM REQUIREMENTS                                                   │
│  ────────────────────                                                  │
│  • Terminal width: 80 columns minimum                                   │
│  • Terminal height: 24 lines minimum                                    │
│  • Color support: 256 colors (8-bit)                                  │
│  • Unicode: UTF-8 support                                              │
│  • Mouse: Optional (xterm reporting)                                   │
│                                                                      │
│  RECOMMENDED                                                           │
│  ──────────                                                            │
│  • Terminal width: 120 columns                                         │
│  • Terminal height: 40 lines                                           │
│  • Color support: Truecolor (24-bit)                                  │
│  • Mouse: Enabled for interactive mode                                 │
│                                                                      │
│  SUPPORTED TERMINALS                                                   │
│  ─────────────────                                                     │
│  • iTerm2 (macOS)                                                     │
│  • Windows Terminal                                                    │
│  • Alacritty                                                          │
│  • Kitty                                                              │
│  • GNOME Terminal                                                     │
│  • xterm                                                              │
│  • tmux/screen                                                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 11.3 Performance Budget

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Performance Budgets                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  INTERACTIVE MODE                                                      │
│  ─────────────────                                                     │
│  • Input response: < 50ms                                             │
│  • Branch selection: < 16ms                                            │
│  • Tree rendering: < 100ms                                             │
│  • Pattern matching: < 100ms                                           │
│                                                                      │
│  COMMAND MODE                                                          │
│  ──────────                                                            │
│  • Pattern execute: < 200ms                                            │
│  • JSON export: < 500ms                                                │
│  • List patterns: < 100ms                                               │
│  • Stats display: < 50ms                                               │
│                                                                      │
│  MEMORY                                                                │
│  ──────                                                               │
│  • Base footprint: < 50MB                                             │
│  • Per pattern: < 10KB                                               │
│  • Cache: < 100MB                                                    │
│                                                                      │
│  STORAGE                                                               │
│  ───────                                                              │
│  • Per pattern: < 10KB JSON                                           │
│  • 10,000 patterns: ~100MB                                             │
│  • SQLite database: ~100MB                                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 11.4 Mockup Examples

#### Interactive Mode Startup

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  otr v1.0.0 ─ AI Thought Reflex Accelerator                           │
│  ───────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  Welcome to Open-Think-Reflex!                                           │
│  Type a query or press [Space] for AI generation.                       │
│  Press [h] for help, [q] to quit.                                       │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  > _                                                                    │
│                                                                              │
│  [Type your query or use ↑/↓ for history]                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### With Matching Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  otr v1.0.0 ─ AI Thought Reflex Accelerator                           │
│  ───────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  [market-analysis] ─► [Data Collection] ✓ 85%                           │
│                      ─► [Pattern ID] ◉ 78%                                │
│                      ─► [Risk Assessment] 72%                             │
│                      ─► [Recommendations] 65%                              │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  # Market Trend Analysis Framework                                        │
│                                                                              │
│  ## 1. Data Sources                                                      │
│     - Internal sales data                                                 │
│     - Customer feedback surveys                                           │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  > market-analysis [↓↑ select, → expand, Space execute, h help]          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Error State

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  otr v1.0.0 ─ AI Thought Reflex Accelerator