# Schema Validation API Specification

## Package

`github.com/sean/apollo/api/internal/schema`

## Public Functions

```go
// Validate checks jsonData against the embedded curriculum JSON schema.
// Returns nil on success, descriptive error with field path on failure.
func Validate(jsonData []byte) error

// ValidatePoolSummary checks jsonData against the knowledge pool summary schema.
// Returns nil on success, descriptive error on failure.
func ValidatePoolSummary(jsonData []byte) error
```

## Embedded Schemas

| File | Source of Truth | Description |
|------|----------------|-------------|
| `curriculum.json` | `schemas/curriculum.json` (project root) | Curriculum structure (Topic → Modules → Lessons → Concepts) |
| `knowledge_pool_summary.json` | `schemas/knowledge_pool_summary.json` (project root) | Knowledge pool context for research sessions |

Schemas are embedded via `embed.FS` and compiled once on first use (`sync.Once`).

## Schema Files (Project Root)

| File | JSON Schema Draft | Used By |
|------|-------------------|---------|
| `schemas/curriculum.json` | 2020-12 | Claude CLI `--json-schema`, Go `Validate()` |
| `schemas/knowledge_pool_summary.json` | 2020-12 | Go `ValidatePoolSummary()` |

## Skill Prompt

| File | Purpose |
|------|---------|
| `skills/research.md` | 4-pass research pipeline prompt loaded via `--system-prompt-file` |

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/santhosh-tekuri/jsonschema/v6` | v6.0.2 | JSON Schema compilation and validation |
