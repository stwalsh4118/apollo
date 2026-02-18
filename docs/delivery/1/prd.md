# PBI-1: Project Foundation & Database Schema

[View in Backlog](../backlog.md)

## Overview

Bootstrap the Go project, establish the SQLite database layer with all tables from the PRD data model, and set up the configuration system. This is the foundation everything else builds on — no API, no frontend, no research until the data layer exists.

## Problem Statement

Apollo has no codebase yet. We need a Go module with dependency management, a SQLite database with the full schema (topics, modules, lessons, concepts, references, prerequisites, relations, expansion queue, research jobs, learning progress, concept retention), idempotent migrations, and environment-based configuration. Without this, no other PBI can store or retrieve data.

## User Stories

- As a developer, I want a Go project scaffold with standard layout so that backend development can begin
- As a developer, I want SQLite schema migrations that create all PRD-defined tables so that curriculum and learning data has a home
- As a developer, I want configuration loaded from environment variables so that the app is configurable without code changes

## Technical Approach

- Standard Go project layout: `cmd/apollo/`, `internal/`, `migrations/`
- SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Migration system: embed SQL files, run on startup, track applied migrations
- All tables from PRD section 10: `topics`, `modules`, `lessons`, `concepts`, `concept_references`, `topic_prerequisites`, `topic_relations`, `expansion_queue`, `research_jobs`, `learning_progress`, `concept_retention`
- Configuration via environment variables with sensible defaults per PRD section 15
- Structured logging with `zerolog`
- Database connection pool and basic health check

## UX/UI Considerations

N/A — backend/infrastructure PBI.

## Acceptance Criteria

1. `go build ./cmd/apollo` produces a binary without errors
2. Binary starts, connects to SQLite, and runs migrations idempotently
3. All 11 tables from PRD section 10 exist with correct columns, types, and foreign keys
4. JSON columns store and retrieve valid JSON (verified with `json_extract()`)
5. FTS5 virtual table created for full-text search across topics, lessons, and concepts
6. Configuration loads from env vars with defaults matching PRD section 15
7. Structured JSON logging works via zerolog
8. Database file created at configured `DATABASE_PATH`

## Dependencies

- **Depends on**: None — this is the first PBI
- **External**: `modernc.org/sqlite`, `zerolog`, Go standard library

## Open Questions

- None

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 1`._
