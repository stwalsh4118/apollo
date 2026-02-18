# Product Backlog

**PRD**: [View PRD](../prd.md)

## Backlog Items

| ID | Actor | User Story | Status | Conditions of Satisfaction (CoS) |
|----|-------|-----------|--------|----------------------------------|
| 1 | Developer | As a developer, I want a Go project scaffold with SQLite schema and migrations so that all data entities have a persistent, validated storage layer | Done | Go module initialised with dependencies; SQLite schema covers all PRD tables; migrations run idempotently; config loaded from env vars |
| 2 | Developer | As a developer, I want REST API endpoints for all curriculum entities so that the frontend can browse topics, modules, lessons, and concepts | Agreed | CRUD endpoints for topics, modules, lessons, concepts; concept references and prerequisites queryable; JSON responses match PRD schema; search endpoint with FTS5 |
| 3 | Learner | As a learner, I want to browse generated curricula in an interactive course UI so that I can read lessons and navigate between modules | Proposed | React app scaffolded with routing; course view renders module sidebar + lesson content; all 6 content section types rendered (text, code, callout, diagram, table, image); Shiki + Mermaid integrated |
| 4 | Learner | As a learner, I want to track my progress through lessons and add personal notes so that I can study methodically and capture insights | Proposed | Progress API stores per-lesson status and notes; frontend shows progress bars per module; mark-complete button works; notes saved and persisted; concept chips link to definitions |
| 5 | Developer | As a developer, I want a research skill prompt and curriculum JSON schema so that Claude Code sessions produce validated, structured curriculum output | Proposed | Research skill prompt implements 4-pass pipeline; JSON schema matches PRD section 7; validation logic rejects malformed output; schema file usable with --json-schema flag |
| 6 | Developer | As a developer, I want the Go orchestrator to spawn Claude Code CLI sessions and run single-topic research end-to-end so that a topic becomes a published curriculum in the knowledge pool | Proposed | Orchestrator spawns CLI session with correct flags; multi-turn pipeline (survey → deep dive → exercises → validation) executes; structured JSON parsed and stored in SQLite; research job status tracked via API |
| 7 | Learner | As a learner, I want to trigger research from the UI and watch progress in real-time so that I know what's happening during long research runs | Proposed | Frontend can submit research requests; progress polling shows current pass, modules planned/completed, prerequisites found; job list shows active and completed jobs; cancel button works |
| 8 | Developer | As a developer, I want the orchestrator to manage a research queue with recursive prerequisite expansion so that essential prerequisites are auto-researched up to depth 3 | Proposed | Research queue processes topics in order; essential prerequisites auto-queued; depth limit (3) enforced; parallel agent execution (configurable max); helpful/deep stored as available for expansion; expansion queue API exposed |
| 9 | Developer | As a developer, I want a connection resolver that integrates new curricula into the knowledge pool so that concepts are deduplicated and cross-referenced across topics | Proposed | Exact slug dedup works; fuzzy matching via CLI session identifies near-matches; cross-references injected bidirectionally; unresolved concepts created for missing topics; conflicts flagged for user review |
| 10 | Learner | As a learner, I want spaced repetition flashcards auto-generated from concepts so that I retain what I learn via SM-2 review sessions | Proposed | SM-2 algorithm implemented per PRD spec; concepts enter review queue on lesson completion; flashcard UI shows front/back with rating buttons; review stats tracked; mastery threshold (90 days) works |
| 11 | Learner | As a learner, I want to browse a knowledge wiki and search across all my studied content so that I can find anything by keyword | Proposed | Topic and concept index pages; concept detail page shows definition + all references; FTS5 search returns topics, lessons, and concepts; breadcrumb navigation works |
| 12 | Learner | As a learner, I want a visual concept map showing how topics and concepts connect so that I get a bird's-eye view of my knowledge | Proposed | D3.js force-directed graph renders; topics as large nodes, concepts as small nodes; edge types distinguished visually; click-to-navigate works; zoom and filter controls functional |
| 13 | Learner | As a learner, I want a dashboard showing my learning overview, review queue, and research status so that I have a single home screen for everything | Proposed | Topic cards with completion %; review queue widget with count + start button; active research jobs widget; available expansions section; concept map thumbnail |
| 14 | Developer | As a developer, I want Docker Compose deployment and production configuration so that Apollo runs as a single self-hosted service | Proposed | Dockerfile builds Go binary + React assets; docker compose up starts the service; env var configuration works; structured logging (zerolog); SQLite backup documented |

_Items are ordered by priority (highest first)._

## PBI Details

| ID | Title | Detail Document |
|----|-------|----------------|
| 1 | Project Foundation & Database Schema | [View Details](./1/prd.md) |
| 2 | Curriculum CRUD API & Search | [View Details](./2/prd.md) |
| 3 | Frontend Foundation & Course View | [View Details](./3/prd.md) |
| 4 | Learning Progress & Notes | [View Details](./4/prd.md) |
| 5 | Research Skill Prompt & Schema | [View Details](./5/prd.md) |
| 6 | Research Agent — Single Topic Pipeline | [View Details](./6/prd.md) |
| 7 | Research Progress UI | [View Details](./7/prd.md) |
| 8 | Research Orchestrator & Queue | [View Details](./8/prd.md) |
| 9 | Connection Resolver | [View Details](./9/prd.md) |
| 10 | Spaced Repetition System | [View Details](./10/prd.md) |
| 11 | Knowledge Wiki & Search UI | [View Details](./11/prd.md) |
| 12 | Concept Map Visualization | [View Details](./12/prd.md) |
| 13 | Dashboard | [View Details](./13/prd.md) |
| 14 | Deployment & Configuration | [View Details](./14/prd.md) |

## History

| Timestamp | PBI_ID | Event_Type | Details | User |
|-----------|--------|------------|---------|------|
| 20260218-090547 | ALL | Created | Initial backlog created from PRD decomposition | AI_Agent |
| 2026-02-18 09:15:55 | 1 | Status Change | Proposed -> InProgress. Started implementation. | AI_Agent |
| 2026-02-18 09:53:53 | 1 | Status Change | InProgress -> Done. All 3 tasks completed and verified. | AI_Agent |
| 2026-02-18 10:50:32 | 2 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
