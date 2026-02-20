# PBI-16: Research Schema Compliance — Prompt Embedding & Output Sanitizer

[View in Backlog](../backlog.md)

## Overview

The research agent produces lesson files with extra and misnamed fields that fail JSON schema validation (`additionalProperties: false`). The first full 4-pass pipeline run on "Go Concurrency" generated high-quality content across 7 modules and 24 lessons but failed assembly with 67 schema errors. The root cause is that the agent never sees the actual curriculum schema — it only sees format examples in the system prompt, which it drifts from across passes.

## Problem Statement

**Current state:** The research prompt (`research.md`) includes a "File Format Reference" section with example JSON structures, but not the actual schema constraints. Across 4 passes of content generation, the agent improvises extra fields (`slug`, `description` on concepts; `concepts` instead of `explanation` on examples) and drops required fields (`description` on examples). The assembler validates against the strict schema and rejects the output.

**Desired state:** The research agent produces schema-compliant files on every run. A programmatic safety net catches any remaining drift before assembly, so pipeline runs never fail on schema validation.

## User Stories

- As a developer, I want the research agent to produce schema-compliant files so that the assembler validates without errors on every run

## Technical Approach

Two complementary fixes:

### 1. Schema-Aware Prompting (Primary Fix)

Embed condensed schema type definitions directly into `internal/research/prompts/research.md`. The current "File Format Reference" section shows example JSON but doesn't enforce constraints. Replace/augment it with explicit type definitions that include:
- Every required field per type
- Allowed values for enums
- Explicit "NO additional fields allowed" callouts
- Field-level descriptions where the name alone is ambiguous

Focus areas (from the 181 issues found in the first run):
- `ConceptTaught`: only `{id, name, definition, flashcard}` — no `slug`, no `description`
- `ConceptReferenced`: only `{id, defined_in}` — no `slug`, no `reason`
- `Example`: requires `{title, description, code, explanation}` — all four mandatory
- `Exercise`: requires `{type, title, instructions, success_criteria, hints, environment}`
- `ReviewQuestion`: requires `{question, answer, concepts_tested}`

### 2. Programmatic Sanitizer (Safety Net)

Add a Go function (`SanitizeDir` or similar) that walks the file tree before `AssembleFromDir` and:
- Strips fields not in the schema (extra keys on `additionalProperties: false` objects)
- Logs warnings for stripped fields (so we can track prompt drift over time)
- Returns an error for missing required fields (can't be auto-fixed)

This runs in milliseconds, costs zero API credits, and catches any residual drift the prompt improvements don't eliminate.

### Integration

The orchestrator calls `SanitizeDir(workDir)` after pass 4 completes and before `AssembleFromDir`. The sanitizer modifies files in place and logs what it changed.

## UX/UI Considerations

N/A — backend/infrastructure PBI.

## Acceptance Criteria

1. The system prompt includes explicit schema type definitions for all lesson file types (ConceptTaught, ConceptReferenced, Example, Exercise, ReviewQuestion, all ContentSection variants) with "no additional fields" constraints clearly stated.
2. A Go sanitizer function strips unknown fields from all JSON files in the work directory before assembly.
3. The sanitizer logs warnings (via zerolog) for each stripped field, including file path and field name.
4. The sanitizer returns an error for missing required fields that cannot be auto-fixed.
5. The orchestrator calls the sanitizer between pass 4 completion and `AssembleFromDir`.
6. A full pipeline run (or replayed fixture tree from the Go Concurrency run) produces 0 schema validation errors after sanitization.
7. All existing tests pass (`make check`).
8. The sanitizer has unit tests covering: extra field stripping, missing required field detection, clean file passthrough.

## Dependencies

- **Depends on**: PBI-15 (file-per-lesson output structure)
- **Blocks**: None
- **External**: None

## Open Questions

None — the approach is clear from the first full pipeline run data.

## Related Tasks

[View Tasks](./tasks.md)
