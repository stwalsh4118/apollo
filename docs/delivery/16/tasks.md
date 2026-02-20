# Tasks for PBI 16: Research Schema Compliance — Prompt Embedding & Output Sanitizer

This document lists all tasks associated with PBI 16.

**Parent PBI**: [PBI 16: Research Schema Compliance — Prompt Embedding & Output Sanitizer](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 16-1 | [Embed schema type definitions in research prompt](./16-1.md) | Proposed | Replace example-only File Format Reference with explicit type definitions including field constraints and "no additional fields" rules |
| 16-2 | [Implement SanitizeDir function with unit tests](./16-2.md) | Proposed | Create Go sanitizer that walks file tree, strips unknown fields, logs warnings, errors on missing required fields; includes unit tests |
| 16-3 | [Integrate sanitizer into orchestrator pipeline](./16-3.md) | Proposed | Wire SanitizeDir call between pass 4 completion and AssembleFromDir in the orchestrator |
| 16-4 | [E2E CoS test — schema compliance pipeline](./16-4.md) | Proposed | End-to-end test with fixture tree verifying sanitizer + assembler produces 0 schema errors |

## Dependency Graph

```
16-1 (prompt)         16-2 (sanitizer + tests)
  │                     │
  │                     ├──→ 16-3 (orchestrator integration)
  │                     │        │
  └─────────────────────┴────────┴──→ 16-4 (E2E CoS test)
```

## Implementation Order

1. **16-1** — No code dependencies; can be done first or in parallel with 16-2. Updates prompt only.
2. **16-2** — Core sanitizer implementation + unit tests. Foundation for remaining tasks.
3. **16-3** — Wires sanitizer into orchestrator. Depends on 16-2.
4. **16-4** — E2E validation of full pipeline. Depends on 16-2 and 16-3.

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 16-1 | Simple | None |
| 16-2 | Medium | None (uses stdlib `encoding/json` + existing `zerolog`) |
| 16-3 | Simple | None |
| 16-4 | Medium | None |

## External Package Research Required

None — all work uses existing Go stdlib and project dependencies.
