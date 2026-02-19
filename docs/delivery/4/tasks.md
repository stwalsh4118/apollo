# Tasks for PBI 4: Learning Progress & Notes

This document lists all tasks associated with PBI 4.

**Parent PBI**: [PBI 4: Learning Progress & Notes](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 4-1 | [Progress models and types](./4-1.md) | Proposed | Define Go models and TypeScript types for learning progress API |
| 4-2 | [Progress repository](./4-2.md) | Proposed | Implement data access layer for learning_progress table |
| 4-3 | [Progress API handlers](./4-3.md) | Proposed | Create HTTP handlers and register routes for progress endpoints |
| 4-4 | [Frontend progress API client and hooks](./4-4.md) | Proposed | Add progress client functions, mutations, and TanStack Query hooks |
| 4-5 | [Mark Complete button and completion indicators](./4-5.md) | Proposed | Add Mark Complete button to lesson view and checkmarks in sidebar |
| 4-6 | [Module progress bars](./4-6.md) | Proposed | Add progress bars per module in the sidebar |
| 4-7 | [Personal notes textarea](./4-7.md) | Proposed | Add inline notes textarea per lesson with save functionality |
| 4-8 | [Concept chips on lessons](./4-8.md) | Proposed | Render concept badges at top of lesson linking to concept detail |
| 4-9 | [E2E CoS Test](./4-9.md) | Proposed | End-to-end tests verifying all PBI 4 acceptance criteria |

## Dependency Graph

```
4-1 (Models & Types)
 ├──► 4-2 (Repository)
 │     └──► 4-3 (API Handlers)
 │           └──► 4-4 (Frontend Client & Hooks)
 │                 ├──► 4-5 (Mark Complete + Indicators)
 │                 │     └──► 4-6 (Progress Bars)
 │                 ├──► 4-7 (Notes Textarea)
 │                 └──► 4-8 (Concept Chips)
 └──────────────────────────► 4-9 (E2E CoS Test) [after all]
```

## Implementation Order

1. **4-1**: Progress models and types — foundational data structures (no dependencies)
2. **4-2**: Progress repository — data access requires models from 4-1
3. **4-3**: Progress API handlers — HTTP layer requires repository from 4-2
4. **4-4**: Frontend progress API client and hooks — requires backend API from 4-3
5. **4-5**: Mark Complete button and completion indicators — requires frontend hooks from 4-4
6. **4-6**: Module progress bars — extends sidebar work from 4-5 (uses same progress map)
7. **4-7**: Personal notes textarea — requires frontend hooks from 4-4 (parallel with 4-5/4-6)
8. **4-8**: Concept chips on lessons — requires frontend data flow (parallel with 4-5/4-6/4-7)
9. **4-9**: E2E CoS Test — must be last, validates all acceptance criteria

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 4-1 | Simple | None |
| 4-2 | Medium | None |
| 4-3 | Medium | None |
| 4-4 | Medium | None |
| 4-5 | Medium | None |
| 4-6 | Simple | None |
| 4-7 | Medium | None |
| 4-8 | Medium | None |
| 4-9 | Complex | None |

## External Package Research Required

None — all tasks use existing project dependencies (chi, TanStack Query, Playwright, Tailwind).
