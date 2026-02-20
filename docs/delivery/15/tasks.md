# Tasks for PBI 15: Research Pipeline — File-Per-Lesson Output Structure

This document lists all tasks associated with PBI 15.

**Parent PBI**: [PBI 15: Research Pipeline — File-Per-Lesson Output Structure](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 15-1 | [Per-File Output Types & Directory Constants](./15-1.md) | Proposed | Define Go structs for intermediate file formats (TopicFile, ModuleFile) and constants for directory naming conventions |
| 15-2 | [Directory Tree Assembler](./15-2.md) | Proposed | Implement Go assembler that walks the modules/ directory tree and produces a valid CurriculumOutput |
| 15-3 | [Update Research System Prompt](./15-3.md) | Proposed | Rewrite research.md pass instructions for file-per-lesson output instead of monolithic JSON |
| 15-4 | [Orchestrator Integration](./15-4.md) | Proposed | Replace runFinalPass with runPass + assembleFromDir, remove --json-schema from Pass 4 |
| 15-5 | [E2E Conditions of Satisfaction Test](./15-5.md) | Proposed | Verify full pipeline with file-per-lesson output meets all acceptance criteria |

## Dependency Graph

```
15-1 (Types & Constants)
  │
  ▼
15-2 (Assembler)──────┐
                       │
15-3 (System Prompt)───┤
                       ▼
                 15-4 (Orchestrator Integration)
                       │
                       ▼
                 15-5 (E2E CoS Test)
```

## Implementation Order

1. **15-1** Per-File Output Types & Directory Constants — foundational types needed by assembler
2. **15-2** Directory Tree Assembler — depends on types from 15-1
3. **15-3** Update Research System Prompt — can run in parallel with 15-2, but listed after since types inform format examples in the prompt
4. **15-4** Orchestrator Integration — depends on assembler (15-2) and prompt (15-3) being ready
5. **15-5** E2E CoS Test — depends on all prior tasks

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 15-1 | Simple | None |
| 15-2 | Medium | None |
| 15-3 | Medium | None |
| 15-4 | Medium | None |
| 15-5 | Complex | None |

## External Package Research Required

None — all implementation uses existing Go standard library and project dependencies.
