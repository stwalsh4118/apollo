# Tasks for PBI 5: Research Skill Prompt & Schema

This document lists all tasks associated with PBI 5.

**Parent PBI**: [PBI 5: Research Skill Prompt & Schema](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 5-1 | [Curriculum JSON Schema](./5-1.md) | Done | Create JSON Schema for the curriculum structure matching PRD section 7 |
| 5-2 | [Knowledge Pool Summary Schema](./5-2.md) | Done | Create JSON Schema for the knowledge_pool_summary.json context file |
| 5-3 | [Research Skill Prompt](./5-3.md) | Done | Create the 4-pass research pipeline skill prompt file |
| 5-4 | [Go Schema Validation Package](./5-4.md) | Done | Create Go package to load and validate curriculum JSON against the schema |
| 5-5 | [E2E CoS Test](./5-5.md) | Done | End-to-end validation of all PBI 5 acceptance criteria |

## Dependency Graph

```
5-1 (Curriculum JSON Schema)
 │
 ├──► 5-4 (Go Schema Validation Package)
 │     │
 │     └──► 5-5 (E2E CoS Test)
 │           ▲
5-2 (Knowledge Pool Summary Schema) ──┘
 │                                     │
5-3 (Research Skill Prompt) ───────────┘
```

## Implementation Order

1. **5-1** Curriculum JSON Schema — foundational; the schema is referenced by the Go validation package and used with `--json-schema` flag
2. **5-2** Knowledge Pool Summary Schema — independent of 5-1; simple, can be done in parallel with 5-1
3. **5-3** Research Skill Prompt — independent of 5-1/5-2 (references schema conceptually but doesn't embed it); can be done in parallel with 5-1/5-2
4. **5-4** Go Schema Validation Package — depends on 5-1 (needs the schema file to embed and validate against)
5. **5-5** E2E CoS Test — depends on 5-1, 5-2, 5-3, 5-4 (validates all deliverables together)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|-----------|-------------------|
| 5-1 | Medium | None |
| 5-2 | Simple | None |
| 5-3 | Complex | None |
| 5-4 | Medium | `github.com/santhosh-tekuri/jsonschema/v6` |
| 5-5 | Medium | None |

## External Package Research Required

| Package | Task | Guide Document |
|---------|------|---------------|
| `github.com/santhosh-tekuri/jsonschema/v6` | 5-4 | `docs/delivery/5/5-4-jsonschema-guide.md` |
