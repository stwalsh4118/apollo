# Tasks for PBI 2: Curriculum CRUD API & Search

This document lists all tasks associated with PBI 2.

**Parent PBI**: [PBI 2: Curriculum CRUD API & Search](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 2-1 | [HTTP Server & Router Foundation](./2-1.md) | Proposed | Set up chi router with middleware, wire HTTP server into main.go, add health endpoint |
| 2-2 | [Response Models & JSON Helpers](./2-2.md) | Proposed | Define Go structs for all API response types and JSON helper functions |
| 2-3 | [Topic Repository & Handlers](./2-3.md) | Proposed | Repository + handlers for topic list, detail, and full tree endpoints |
| 2-4 | [Module & Lesson Repository & Handlers](./2-4.md) | Proposed | Repository + handlers for module and lesson detail endpoints |
| 2-5 | [Concept Repository & Handlers](./2-5.md) | Proposed | Repository + handlers for concept list, detail, and references endpoints |
| 2-6 | [Curriculum Write Endpoints](./2-6.md) | Proposed | Internal write repository + handlers for research pipeline data ingestion |
| 2-7 | [Search Endpoint](./2-7.md) | Proposed | FTS5 search repository + handler for full-text search |
| 2-8 | [Knowledge Graph Endpoints](./2-8.md) | Proposed | Graph repository + handlers for topic/concept graph visualization data |
| 2-9 | [E2E CoS Test](./2-9.md) | Proposed | End-to-end verification of all PBI 2 acceptance criteria |

## Dependency Graph

```
2-1 (HTTP Server & Router)
 └──► 2-2 (Response Models & JSON Helpers)
       ├──► 2-3 (Topic Repository & Handlers)
       ├──► 2-4 (Module & Lesson Repository & Handlers)
       ├──► 2-5 (Concept Repository & Handlers)
       ├──► 2-6 (Curriculum Write Endpoints)
       ├──► 2-7 (Search Endpoint)
       └──► 2-8 (Knowledge Graph Endpoints)
             │
             ▼
           2-9 (E2E CoS Test) ◄── all of 2-3 through 2-8
```

## Implementation Order

1. **2-1** — HTTP Server & Router Foundation (no dependencies; everything else needs a running server)
2. **2-2** — Response Models & JSON Helpers (depends on 2-1; all handlers need response types)
3. **2-3** — Topic Repository & Handlers (depends on 2-2; establishes the repository pattern for others)
4. **2-4** — Module & Lesson Repository & Handlers (depends on 2-2; follows pattern from 2-3)
5. **2-5** — Concept Repository & Handlers (depends on 2-2; follows pattern from 2-3)
6. **2-6** — Curriculum Write Endpoints (depends on 2-2; write counterpart to read endpoints)
7. **2-7** — Search Endpoint (depends on 2-2 and 2-6; search_index populated by write endpoints)
8. **2-8** — Knowledge Graph Endpoints (depends on 2-2 and 2-6; graph needs data in DB)
9. **2-9** — E2E CoS Test (depends on all of 2-1 through 2-8; verifies everything together)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 2-1 | Medium | chi |
| 2-2 | Simple | None |
| 2-3 | Complex | None |
| 2-4 | Medium | None |
| 2-5 | Medium | None |
| 2-6 | Complex | None |
| 2-7 | Medium | None |
| 2-8 | Medium | None |
| 2-9 | Complex | None |

## External Package Research Required

| Package | Task | Guide Document |
|---------|------|---------------|
| `github.com/go-chi/chi/v5` | 2-1 | `2-1-chi-guide.md` |
