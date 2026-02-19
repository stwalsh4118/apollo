# Tasks for PBI 6: Research Agent — Single Topic Pipeline

This document lists all tasks associated with PBI 6.

**Parent PBI**: [PBI 6: Research Agent — Single Topic Pipeline](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 6-1 | [Research Job Models & Constants](./6-1.md) | Proposed | Define Go types for research jobs, CLI responses, progress tracking, and configuration constants |
| 6-2 | [Research Job Repository](./6-2.md) | Proposed | CRUD operations for research_jobs table (create, get, list, update status/progress, cancel) |
| 6-3 | [Research API Handlers](./6-3.md) | Proposed | HTTP handlers for POST /api/research, GET /api/research/jobs, GET /api/research/jobs/:id, POST /api/research/jobs/:id/cancel; server wiring; API spec |
| 6-4 | [Knowledge Pool Summary Builder](./6-4.md) | Proposed | Service to query existing topics/concepts from DB and produce validated knowledge_pool_summary.json |
| 6-5 | [CLI Session Spawner](./6-5.md) | Proposed | exec.Command wrapper for Claude Code CLI with flag construction, JSON parsing, session ID extraction, and process management |
| 6-6 | [Curriculum Ingester](./6-6.md) | Proposed | Parse structured JSON output, validate against schema, store full curriculum (topic, modules, lessons, concepts, prerequisites) in SQLite |
| 6-7 | [Research Pipeline Orchestrator](./6-7.md) | Proposed | Job lifecycle management, context preparation, multi-turn 4-pass pipeline execution, progress tracking, error handling, cancellation |
| 6-8 | [E2E CoS Test](./6-8.md) | Proposed | End-to-end tests verifying all PBI 6 acceptance criteria with mocked CLI sessions |

## Dependency Graph

```
6-1 (Models & Constants)
 ├──> 6-2 (Repository)
 │     ├──> 6-3 (API Handlers)
 │     ├──> 6-6 (Curriculum Ingester)
 │     └──> 6-7 (Orchestrator)
 ├──> 6-4 (Pool Summary Builder)
 │     └──> 6-7 (Orchestrator)
 ├──> 6-5 (CLI Session Spawner)
 │     └──> 6-7 (Orchestrator)
 └──> 6-6 (Curriculum Ingester)
       └──> 6-7 (Orchestrator)
             └──> 6-8 (E2E CoS Test)
```

## Implementation Order

1. **6-1** — Models & Constants (no dependencies; foundational types for all other tasks)
2. **6-2** — Research Job Repository (depends on 6-1; data layer needed by handlers and orchestrator)
3. **6-3** — Research API Handlers (depends on 6-1, 6-2; provides HTTP interface for job creation)
4. **6-4** — Knowledge Pool Summary Builder (depends on 6-1; independent of API layer)
5. **6-5** — CLI Session Spawner (depends on 6-1; independent of API layer)
6. **6-6** — Curriculum Ingester (depends on 6-1, 6-2; needs models and DB access)
7. **6-7** — Research Pipeline Orchestrator (depends on 6-2, 6-4, 6-5, 6-6; ties everything together)
8. **6-8** — E2E CoS Test (depends on all above; validates full pipeline)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|-----------|-------------------|
| 6-1 | Simple | None |
| 6-2 | Medium | None |
| 6-3 | Medium | None |
| 6-4 | Medium | None |
| 6-5 | Complex | None |
| 6-6 | Complex | None |
| 6-7 | Complex | None |
| 6-8 | Complex | None |

## External Package Research Required

None. All tasks use stdlib (`os/exec`, `encoding/json`, `database/sql`) and existing project dependencies (`chi`, `zerolog`, `jsonschema`).
