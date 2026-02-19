# PBI-6: Research Agent — Single Topic Pipeline

[View in Backlog](../backlog.md)

## Overview

Implement the Go orchestrator's ability to spawn a Claude Code CLI session and run the complete single-topic research pipeline end-to-end. One topic in, one published curriculum out. This covers session spawning, multi-turn execution (4 passes via `--resume`), output parsing, SQLite storage, and research job status tracking.

## Problem Statement

The research skill prompt and schema exist (PBI 5) but nothing executes them. The Go orchestrator needs to: prepare context files, spawn a Claude Code CLI session with the right flags, run 4 passes via `--resume`, parse the structured JSON output, validate it against the schema, store the resulting curriculum in SQLite, and track job status throughout. This is the core engine that turns a topic into a course.

## User Stories

- As a learner, I want to submit a topic and have a complete curriculum generated automatically (US-1)
- As a learner, I want to see research progress while agents are working (US-12)

## Technical Approach

- Research job lifecycle: `queued` → `researching` → `resolving` → `published` (or `failed`)
- Context preparation:
  - Write `knowledge_pool_summary.json` to working directory (existing topics + concepts from SQLite)
  - Prepare topic brief with depth, size limit, and boundary guidance
- Claude Code CLI spawning via `exec.Command`:
  - Pass 1 (Survey): `-p` with topic brief, `--system-prompt-file`, `--output-format json`, `--model opus`, `--allowedTools`
  - Pass 2-3: `--resume $SESSION_ID`
  - Pass 4: `--resume $SESSION_ID`, `--json-schema` for structured output
  - Extract `session_id` from each JSON response for `--resume`
  - Extract `structured_output` from final response
- Output processing:
  - Validate structured output against curriculum schema (PBI 5)
  - Store topic, modules, lessons, concepts in SQLite (using PBI 2 write endpoints or direct DB)
  - Create concept entries with flashcards
  - Store prerequisite classifications in `topic_prerequisites` and `expansion_queue`
- Research job API (PRD section 9.1):
  - `POST /api/research` — create and queue a research job
  - `GET /api/research/jobs` — list jobs
  - `GET /api/research/jobs/:id` — job status and progress
  - `POST /api/research/jobs/:id/cancel` — cancel a running job
- Progress tracking: update `research_jobs.progress` JSON field during execution (current pass, modules completed, concepts found)
- Error handling: retry individual passes on failure, mark job as `failed` with error message after max retries

## UX/UI Considerations

N/A — backend PBI. The research progress UI comes in PBI 7.

## Acceptance Criteria

1. `POST /api/research` creates a research job and returns job ID
2. Orchestrator spawns Claude Code CLI session with correct flags per PRD 6.3
3. Multi-turn pipeline executes 4 passes via `--resume`
4. Structured JSON output parsed and validated against schema
5. Curriculum stored in SQLite: topic, modules, lessons, concepts, flashcards, prerequisites
6. Research job progress updates during execution (current pass, modules, concepts)
7. `GET /api/research/jobs/:id` returns accurate status and progress
8. Failed research sessions produce a `failed` job with error details
9. Cancel endpoint stops a running research session
10. End-to-end: submit a topic, research completes, curriculum browsable via PBI 2 API

## Dependencies

- **Depends on**: PBI 5 (research skill prompt and schema), PBI 2 (storage endpoints or DB access), PBI 1 (database)
- **External**: Claude Code CLI (`claude` binary), Max plan subscription

## Open Questions

- None

## Related Tasks

[View Tasks](./tasks.md)
