# Tasks for PBI 1: Project Foundation & Database Schema

This document lists all tasks associated with PBI 1.

**Parent PBI**: [PBI 1: Project Foundation & Database Schema](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 1-1 | [Go Project Scaffold, Config, and Logging](./1-1.md) | Done | Create Go module layout, startup entrypoint, environment config loader, and structured zerolog wiring. |
| 1-2 | [SQLite Schema and Migration Engine](./1-2.md) | Done | Implement embedded SQL migrations, full PRD schema, migration tracking, and startup DB health checks. |
| 1-3 | [PBI-1 End-to-End Verification](./1-3.md) | Done | Run PBI Conditions of Satisfaction verification and document evidence for idempotent migrations, JSON functions, and build health. |

## Dependency Graph

```text
1-1 -> 1-2 -> 1-3
```

## Implementation Order

1. `1-1`: establish executable scaffold and configuration primitives used by all DB setup code.
2. `1-2`: add database connection lifecycle and full schema migrations on top of the scaffold.
3. `1-3`: run acceptance checks after all foundation behavior exists.

## Complexity Ratings

| Task ID | Complexity | External Packages |
|--------|------------|-------------------|
| 1-1 | Medium | `github.com/rs/zerolog` |
| 1-2 | Complex | `modernc.org/sqlite` |
| 1-3 | Medium | None |

## External Package Research Required

| Task ID | Package | Guide File |
|--------|---------|------------|
| 1-1 | `github.com/rs/zerolog` | `docs/delivery/1/1-1-zerolog-guide.md` |
| 1-2 | `modernc.org/sqlite` | `docs/delivery/1/1-2-modernc-sqlite-guide.md` |
