# Open-Think-Reflex Package Structure

## Overview

This directory contains the Go source code for Open-Think-Reflex.

## Directory Structure

```
cmd/cli/                    # CLI entry point
├── main.go                # Application entry point
└── root.go                # Root command definition

internal/                   # Internal packages (not exported)
├── cli/                   # CLI Interface Layer
│   ├── ui/               # Terminal UI (tview/lipgloss)
│   ├── commands/          # CLI commands (pattern, config, etc.)
│   └── output/            # Output formatters (terminal, json, plain)
│
├── core/                  # Core Logic Layer
│   ├── matcher/           # Pattern matching engine
│   ├── pattern/           # Pattern management
│   └── reflex/            # Reflex lifecycle management
│
├── data/                  # Data Layer
│   ├── sqlite/            # SQLite storage implementation
│   └── cache/             # Caching (LRU, memory)
│
├── ai/                     # AI Integration Layer
│   ├── provider/          # AI provider implementations
│   ├── prompt/            # Prompt builders
│   └── response/          # Response parsers
│
└── config/                 # Configuration loading

pkg/                       # Public packages (exportable)
├── contracts/              # Interface definitions
├── export/                 # Export/Import utilities
└── utils/                  # Shared utilities
```

## Package Naming Conventions

- `internal/*` - Implementation details, not exported
- `pkg/contracts/*` - Interfaces and contracts
- Short, concise names preferred
- No package name in file name

## Dependencies

Run `go mod tidy` to update dependencies.
